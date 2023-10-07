# Dynamic Remote IP matcher for Caddy

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/lanrat/graphs/caddy-dynamic-remoteip) [![Go Reference](https://pkg.go.dev/badge/github.com/lanrat/caddy-dynamic-remoteip.svg)](https://pkg.go.dev/github.com/lanrat/caddy-dynamic-remoteip) [![GoReportCard example](https://goreportcard.com/badge/github.com/lanrat/caddy-dynamic-remoteip)](https://goreportcard.com/report/github.com/lanrat/caddy-dynamic-remoteip) ![GitHub](https://img.shields.io/github/license/lanrat/caddy-dynamic-remoteip)

The `dynamic_remote_ip` module is a clone of the `remote_ip` matcher with one key difference: instead of providing IP ranges upfront, you specify an `IPRangeSource`. This allows IP ranges to be dynamically loaded per request.

This module is based on the [caddy-dynamic-clientip](https://github.com/tuzzmaniandevil/caddy-dynamic-clientip) module by [tuzzmaniandevil](https://github.com/tuzzmaniandevil).

## Installation

Build Caddy using [xcaddy](https://github.com/caddyserver/xcaddy):

```shell
xcaddy build --with github.com/lanrat/caddy-dynamic-remoteip
```

## Usage

```caddyfile
:8880 {
    @denied dynamic_remote_ip my_dynamic_provider
    abort @denied

    reverse_proxy localhost:8080
}
```

Example using the built-in static provider (But why though?)

```caddyfile
:8880 {
    @denied dynamic_remote_ip static 12.34.56.0/24 1200:ab00::/32
    abort @denied

    reverse_proxy localhost:8080
}
```

## Development

Before diving into development, make sure to follow the [Extending Caddy](https://caddyserver.com/docs/extending-caddy#extending-caddy) guide. This ensures you're familiar with the Caddy development process and that your environment is set up correctly.

To run Caddy with this module:

```shell
xcaddy run
```

## License

The project is licensed under the [Apache License](LICENSE).
