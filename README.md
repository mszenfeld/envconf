# envconf
Go library for loading configuration from the environment variables. 

## Installation
```
go get github.com/mszenfeld/envconf
```

## Getting Started

To start working with `envconf` define a structure for your config and use the loader:

```sh
export HOST="localhost"
export PORT="1337"
export DEBUG="true"
```
```go
package main

import (
  "log"

  "github.com/mszenfeld/envconf"
)

type Config struct {
  Host  string
  Port  int
  Debug bool
}

func main() {
  var c Config

  l := envconf.NewLoader()
  if err := l.Load(&c); err != nil {
    log.Fatal(err)
  }
}
```
`envconf` also supports environment variable prefixes:
```sh
export APP_HOST="localhost"
export APP_PORT="1337"
export APP_DEBUG="true"
```
```go
package main

import (
  "log"

  "github.com/mszenfeld/envconf"
)

type Config struct {
  Host  string
  Port  int
  Debug bool
}

func main() {
  var c Config

  l := envconf.NewLoader()
  l.SetPrefix("app")
  if err := l.Load(&c); err != nil {
    log.Fatal(err)
  }
}
```

## Tags

`envconf` allows you to use the following tags:
- `env`
- `default`
- `required`

```go
type Config struct {
  Host  string `default:"localhost"`
  Port  int `required:"true"`
  Debug bool `env:"ENABLE_DEBUG" default:"false"`
}
```

### `env`

This tag overrides the default environment variable name of the struct field.

### `default`

The tag value will be used as a field's value if the environment variable associated with the field does not exist.

### `required`

The available values for this tag are only "true" and "false". If field 
is marked as required and environment variable does not exist, loader 
will return an error. By default, fields are optional.

## Supported Types
`envconf` has support for the following types:
- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
