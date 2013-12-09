# logger - Logging framework for go

Simple logging framework for go. It supports different logging levels and
has inheritance of that level to child loggers.

## Installation
### Using `go get`

    $ go get github.com/AlexanderThaller/logger

This will clone the source of *logger* to
`$GOROOT/src/pkg/github.com/op/go-logging`

After that you can use *logger* by importing
`github.com/AlexanderThaller/logger` into your *go* application.

## Example
```go
package main

import (
  "github.com/AlexanderThaller/logger"
)

func init() {
  // Set application logger to Info (default is Notice)
  logger.SetLevel("application", logger.Info)
}

func main() {
  l := logger.NewLogger("application.main")

  logger.NoticeM(l, "Starting")

  // This message will be shown because in the init function we set the
  // level of the "application" logger to Info
  logger.InfoM(l, "Some infos")

  // This message will be hidden because it is below the logger level
  // Info
  logger.DebugM(l, "Debug Values")

  logger.SetLevel(l, logger.Debug)

  // This message will be shown because we set the level of the "application.main"
  // logger to Debug
  logger.DebugM(l, "More debug values")
}
```
