# Argus Monitor Data API for Go

Go interface for the Argus Monitor Data API based on the [reference C++ implementation](https://github.com/argotronic/argus_data_api)

## Usage

Requires Argus Monitor 6.0.01+

In Argus Monitor:

  - Open Settings (F2)
  - Stability
  - Enable Argus Data Monitor API
  - OK and restart Argus Monitor

### Telegraf Plugin

This library comes with a Telegraf execd compatible plugin

```
# Note: How you put this binary in your PATH is up to you
go build -o telegraf-argus.exe github.com/relvacode/cmd/telegraf-argus
```

Configure Telegraf with

```toml
[[inputs.execd]]
command = ["telegraf-argus.exe"]
signal = "STDIN"
data_format = "influx"
```