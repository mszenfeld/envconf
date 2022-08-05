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

## Supported Types
`envconf` has support for the following types:
- string
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- bool
