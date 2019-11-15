# Easylog

A simple logger for Go, support output level.

[![Build Status](https://github.com/winterssy/easylog/workflows/CI/badge.svg)](https://github.com/winterssy/easylog/actions) [![Go Report Card](https://goreportcard.com/badge/winterssy/easylog)](https://goreportcard.com/report/winterssy/easylog) [![GoDoc](https://godoc.org/github.com/winterssy/easylog?status.svg)](https://godoc.org/github.com/winterssy/easylog) [![License](https://img.shields.io/github/license/winterssy/easylog.svg)](LICENSE)

## Install

```sh
go get -u github.com/winterssy/easylog
```

## Usage
```go
import "github.com/winterssy/easylog"
```

## Example
```go
package main

import "github.com/winterssy/easylog"

func main() {
	easylog.SetLevel(easylog.Ldebug)
	easylog.Debug("hello world")
	easylog.Info("hello world")
	easylog.Warn("hello world")
	easylog.Error("hello world")
}
```

## License
MIT.
