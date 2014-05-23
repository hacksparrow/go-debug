
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
    79us 15us     multiple:a - doing stuff
    44us 55us     multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   101ms 101ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   101ms 101ms    multiple:b - doing stuff
   101ms 101ms    multiple:b - doing stuff
   101ms 101ms    multiple:b - doing stuff
   101ms 101ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
    92ms 1s       multiple:a - doing stuff
     8ms 101ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
   100ms 100ms    multiple:b - doing stuff
```

Two deltas are displayed, the left-most delta is relative to the previous debug call of any name, followed by a delta specific to that debug function. These may be useful to identify timing issues and potential bottlenecks.

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

# License

MIT