package mcd

import (
  "errors"
  "os"
  "syscall"
  "unicode/utf16"
  "unsafe"
  "path/filepath"
  "io/ioutil"
  "github.com/robertkrimen/otto"
)

type Configs struct {
  vm *otto.Otto
}

func (c Configs) GetStringValue(key string) (stringValue string, err error) {
  value, err := c.vm.Run("getPref('" + key + "')")
  stringValue, err = value.ToString();
  if stringValue == "undefined" {
    return "", errors.New("unknown pref: " + key)
  }
  return stringValue, nil
}

func (c Configs) GetIntegerValue(key string) (integerValue int64, err error) {
  value, err := c.vm.Run("getPref('" + key + "')")
  stringValue, err := value.ToString();
  if stringValue == "undefined" {
    return 0, errors.New("unknown pref: " + key)
  }
  integerValue, err = value.ToInteger();
  if err != nil {
    return 0, err
  }
  return integerValue, nil
}

func (c Configs) GetBooleanValue(key string) (booleanValue bool, err error) {
  value, err := c.vm.Run("getPref('" + key + "')")
  stringValue, err := value.ToString();
  if stringValue == "undefined" {
    return false, errors.New("unknown pref: " + key)
  }
  booleanValue, err = value.ToBoolean();
  if err != nil {
    return false, err
  }
  return booleanValue, nil
}

func New() (configs Configs, err error) {
  local := ReadLocalConfigs()
  remote := ReadRemoteConfigs()

  vm := otto.New()
  configs = Configs{vm: vm}

  vm.Set("getenv", func(call otto.FunctionCall) otto.Value {
    name := call.Argument(0).String()
    result, _ := vm.ToValue(os.ExpandEnv("${" + name + "}"))
    return result
  })
  // See also https://dxr.mozilla.org/mozilla-central/source/extensions/pref/autoconfig/src/prefcalls.js
  _, err = vm.Run(`
    var $$defaultPrefs = {};
    var $$prefs = {};
    function pref(key, value) {
      $$prefs[key] = value;
    }
    function defaultPref(key, value) {
      $$defaultPrefs[key] = value;
    }
    function lockPref(key, value) {
      delete $$prefs[key];
      $$defaultPrefs[key] = value;
    }
    function clearPref(key) {
      delete $$prefs[key];
    }
    function getPref(key) {
      if (key in $$prefs)
        return $$prefs[key];
      if (key in $$defaultPrefs)
        return $$defaultPrefs[key];
      return undefined;
    }
    function unlockPref(key) {
    }
    var Components = {
      classes: {},
      interfaces: {},
      utils: {}
    };
  ` + local + "\n" + remote)


  return configs, err
}

func ReadLocalConfigs() (configs string) {
  path, err := GetLocalConfigPath()
  if err != nil {
    return ""
  }
  buffer, err := ioutil.ReadFile(path)
  if err != nil {
    return ""
  }
  return string(buffer)
}

func GetLocalConfigPath() (path string, err error) {
  exePath, err := GetPathToRunningApp()
  if err != nil {
    return "", err
  }
  //TODO: We should detect the effective file.
  // Currently we return the first one always.
  pattern := filepath.Join(filepath.Dir(exePath), "*.cfg")
  path, err = GetFirstMatchedFile(pattern)
  return
}

func GetFirstMatchedFile(pattern string) (path string, err error) {
  possibleFiles, err := filepath.Glob(pattern)
  if err != nil {
    return "", err
  }
  if len(possibleFiles) == 0 {
    return "", errors.New("no match")
  }
  return possibleFiles[0], nil
}

const PROCESS_VM_READ = 1 << 4

func GetPathToRunningApp() (path string, err error) {
  parentId := os.Getppid()
  inheritHandle := false
  processHandle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, inheritHandle, uint32(parentId))
  defer syscall.CloseHandle(processHandle)
  if err != nil {
    return "", err
  }
  getModuleFileNameEx := syscall.MustLoadDLL("psapi.dll").MustFindProc("GetModuleFileNameExW")
  buffer := make([]uint16, syscall.MAX_PATH)
  bufferSize := uint32(len(buffer))
  rawLength, _, err := getModuleFileNameEx.Call(uintptr(processHandle), 0, uintptr(unsafe.Pointer(&buffer[0])), uintptr(bufferSize))
  length := uint32(rawLength)
  if length == 0 {
    return "", errors.New("failed to get the path of the application")
  }
  return string(utf16.Decode(buffer[0:length])), nil
}

func ReadRemoteConfigs() (configs string) {
  // codes to read failover.jsc in the profile
  path, err := GetFailoverJscPath()
  if err != nil {
    return ""
  }
  buffer, err := ioutil.ReadFile(path)
  if err != nil {
    return ""
  }
  return string(buffer)
}

func GetFailoverJscPath() (path string, err error) {
  //TODO: We should detect the actually used profile directory.
  // Currently we return the default profile.
  pattern := os.ExpandEnv(`${AppData}\Mozilla\Firefox\Profiles\*.default\failover.jsc`)
  path, err = GetFirstMatchedFile(pattern)
  if path != "" {
    return
  }
  pattern = os.ExpandEnv(`${AppData}\Mozilla\Firefox\Profiles\*\failover.jsc`)
  path, err = GetFirstMatchedFile(pattern)
  return
}

