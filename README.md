<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/time" width="720"></p>

# time

[![ci](https://github.com/go-composites/time/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/time/actions/workflows/ci.yml)

A `Time` composite (plus a `Duration` subpackage) for **Composition-Oriented
Programming**. `Time` wraps a Go `time.Time` (Ruby's `Time` as the reference)
and exposes its **fallible operations as `Result` values** — so a malformed
input (the canonical example being a bad `Parse`) is a *value*, never a panic
and never `nil`.

`Time` is deterministic by construction: there is **no `Now()`**. Instances are
built only from explicit values (`FromUnix`, `Parse`), which keeps behaviour
reproducible across architectures and test runs.

```golang
parsed := Time.Parse(time.RFC3339, value)
if parsed.HasError() {
    fmt.Println(parsed.Error().Message())
} else {
    fmt.Println(parsed.Payload().(Time.Interface).ToGoString())
}
```

`Time` and `Duration` follow the org's Null-Object / never-nil invariant
(enforced by the `nonnil` CI analyzer): `Time.Null()` and `Duration.Null()`
satisfy the same `Interface` and report `IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/time@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"
    "time"

    Duration "github.com/go-composites/time/src/duration"
    Time "github.com/go-composites/time/src"
)

func main() {
    epoch := Time.FromUnix(0)
    later := Time.FromUnix(3600)

    // Comparisons are plain Go bools.
    fmt.Println(epoch.Before(later)) // true

    // Sub yields the Duration between two Times.
    fmt.Println(later.Sub(epoch).ToGoString()) // 1h0m0s

    // Add returns a Result — shifting an instant by a Duration.
    shifted := epoch.Add(Duration.FromSeconds(90))
    fmt.Println(shifted.Payload().(Time.Interface).ToGoString())

    // Parsing is fallible: a malformed value is a value, not a panic.
    bad := Time.Parse(time.RFC3339, "not-a-time")
    fmt.Println("has error:", bad.HasError()) // true
}
```

```bash
$ task build
```

## API

### Time (`github.com/go-composites/time/src`, package `Time`)

Constructors

- `FromUnix(sec int64) Interface` — an instant from a Unix timestamp (UTC).
- `Parse(layout, value string) Result.Interface` — a `Result` whose payload is a
  `Time`, or an `Error.New(...)` when `value` does not match `layout`.
- `Null() Interface` — the Null-Object `Time` (`IsNull() == true`).

Conversions

- `ToUnix() int64`, `Format(layout string) string`, `ToGoString() string`
  (RFC3339).

Time zones (the IANA tz database is embedded via `_ "time/tzdata"`, so zone
lookups work on every architecture and inside toolchain-less CI containers
without a system tzdata install)

- `In(name string) Result.Interface` — a `Result` whose payload is a new `Time`
  denoting the **same instant** in the IANA location `name` (e.g.
  `"Europe/Paris"`), or an `Error.New(...)` when `name` is an unknown zone.
  Because only the location changes, `Before`/`After`/`Equal` are preserved.
- `Zone() string` — the zone abbreviation in effect at the instant (e.g.
  `"UTC"`, `"CET"`).
- `UTC() Interface` — the same instant with its location set to UTC.

Comparisons (each returns `bool`)

- `Before(other)` / `After(other)` / `Equal(other)`.

Arithmetic

- `Add(Duration.Interface) Result.Interface` — a `Result` whose payload is a new
  `Time` shifted forward by the duration.
- `Sub(other Interface) Duration.Interface` — the `Duration` between two instants.

Null-Object

- `IsNull() bool`.

### Duration (`github.com/go-composites/time/src/duration`, package `Duration`)

Constructors

- `FromSeconds(sec int64) Interface`.
- `Parse(s string) Result.Interface` — a `Result` whose payload is a `Duration`,
  or an `Error.New(...)` when `s` is malformed (see `time.ParseDuration`).
- `Null() Interface` — the Null-Object `Duration` (`IsNull() == true`).

Methods

- `ToSeconds() int64`, `ToGoString() string`.
- `Add(other) Interface` / `Sub(other) Interface` — arithmetic.
- `Equal(other) bool`, `IsNull() bool`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
