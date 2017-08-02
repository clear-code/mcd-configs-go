# mcd-go

A library to read values from [MCD (Misson Control Desktop)](https://developer.mozilla.org/en-US/docs/MCD,_Mission_Control_Desktop_AKA_AutoConfig) configuration files for Firefox addons.
This is strongly designed to implement [native messaging host](https://developer.mozilla.org/en-US/Add-ons/WebExtensions/Native_messaging) applications.

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
  configs, err := mcd.Load()
  if err != nil {
    log.Fatal(err)
  }

  response := &ReadConfigsResponse{}

  foo, err := mcd.GetStringValue(configs, "extensions.myaddon@example.com.foo")
  if err == nil { response.Foo = foo }
  bazz, err := mcd.GetIntegerValue(configs, "extensions.myaddon@example.com.bar")
  if err == nil { response.Bar = bar }
  bazz, err := mcd.GetBooleanValue(configs, "extensions.myaddon@example.com.bazz")
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
 * This just loads the `failover.jsc` file in the default Firefox profile placed under `%AppData%\Mozilla\Profiles\*.default`.
   * If you actually use different profile, this doesn't detect it.

## License

MPL 2.0
