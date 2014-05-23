
# go-debug

 Conditional debug logging for Go libraries.

## Installation

```
$ go get github.com/visionmedia/go-debug
```

## Example

```go
package main

import . "github.com/visionmedia/go-debug"
import "time"

var debug DebugFunction = Debug("single")

func main() {
  for {
    debug("sending mail")
    debug("send email to %s", "tobi@segment.io")
    debug("send email to %s", "loki@segment.io")
    debug("send email to %s", "jane@segment.io")
    time.Sleep(500 * time.Millisecond)
  }
}
```

outputs:

```
    43us single - sending mail
    31us single - send email to tobi@segment.io
     4us single - send email to loki@segment.io
     3us single - send email to jane@segment.io
   500ms single - sending mail
    49us single - send email to tobi@segment.io
    11us single - send email to loki@segment.io
    10us single - send email to jane@segment.io
   501ms single - sending mail
    60us single - send email to tobi@segment.io
    16us single - send email to loki@segment.io
    11us single - send email to jane@segment.io
   500ms single - sending mail
    71us single - send email to tobi@segment.io
    28us single - send email to loki@segment.io
    29us single - send email to jane@segment.io
```

## What?

 The basic premise is that every library should have some form of debug logging,
 ideally enabled without touching code. When disabled a no-op function is returned,
 which Go can easily execute 100m ops/s on a MBP retina, in other words it's negligable for most code paths.

 Executables often support `--verbose` flags for conditional logging, but
 libraries typically either require altering your code to enable logging,
 or simply omit logging all together. go-debug allows conditional logging
 to be enabled via the __DEBUG__ environment variable, where one or more
 patterns may be specified.

 For example suppose your application has several models and you want
 to output logs for users only, you might use `DEBUG=models:user`. In contrast
 if you wanted to see what all database activity was you might use `DEBUG=models:*`,
 or if you're love being swamped with logs: `DEBUG=*`.

 The name given _should_ be the package name, however you can use whatever you like.

 A nanosecond delta is also displayed in the log output to help identify timing issues
 or potential bottlenecks.

# License

MIT