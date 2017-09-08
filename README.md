# mcd-go

A library to read values from [MCD (Mission Control Desktop)](https://developer.mozilla.org/en-US/docs/MCD,_Mission_Control_Desktop_AKA_AutoConfig) configuration files for Firefox addons.
This is strongly designed to implement [native messaging host](https://developer.mozilla.org/en-US/Add-ons/WebExtensions/Native_messaging) applications.

This is a workaround alternative of missing [`storage.managed`](https://developer.mozilla.org/en-US/Add-ons/WebExtensions/API/storage/managed). See also the [bug 1230802](https://bugzilla.mozilla.org/show_bug.cgi?id=1230802). After the API is landed, you should migrate to it.

## Usage

```go
import (
  "github.com/clear-code/mcd-go"

  "log"
  "encoding/json"
  "github.com/lhside/chrome-go"
)

type ReadConfigsResponse struct {
  Foo  string `json:"foo,omitempty"`
  Bar  int64  `json:"bar,omitempty"`
  Bazz bool   `json:"bazz,omitempty"`
}

func ReadConfigs() {
  configs, err := mcd.New()
  if err != nil {
    log.Fatal(err)
  }

  response := &ReadConfigsResponse{}

  foo, err := configs.GetStringValue("extensions.myaddon@example.com.foo")
  if err == nil { response.Foo = foo }
  bar, err := configs.GetIntegerValue("extensions.myaddon@example.com.bar")
  if err == nil { response.Bar = bar }
  bazz, err := configs.GetBooleanValue("extensions.myaddon@example.com.bazz")
  if err == nil { response.Bazz = bazz }

  body, err := json.Marshal(response)
  if err != nil {
    log.Fatal(err)
  }
  err = chrome.Post(body, os.Stdout)
  if err != nil {
    log.Fatal(err)
  }
}
```

## Restrictions

 * This just loads the first `*.cfg` file under the Firefox's directory.
   * If there are multiple files, others are simply ignored even if one of them is actually used.
   * Even if the `*.cfg` file is not used, this always loads it.
 * This just loads the `failover.jsc` file in the default Firefox profile placed under `%AppData%\Mozilla\Profiles\*.default`, as the remote configuration file. In other words, this doesn't fetch actual remote configuration file specified via `autoadmin.global_config_url`.
   * If you actually use different profile, this doesn't detect it.
 * Loaded configuration files are parsed in a sandbox.
   * Only limited directives are available.
   * User values in the profile are not accessible.
   * `Components.utils.import()` and other internal features don't work. If your configuration file depends on such internal technologies, this cannot load values correctly.

## License

MPL 2.0
