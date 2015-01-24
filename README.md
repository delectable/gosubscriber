gosubscriber
============
gosubscriber is a [Resque Bus](https://github.com/taskrabbit/resque-bus)-compatible, Go-based background subscriber build on top of Benjamin Manns' [goworker](https://github.com/benmanns/goworker). It allows you to publish Resque Bus events in Ruby (or any other language), and perform background tasks in Go.

Gosubscriber subscribers can run alongside Ruby/Node.js ResqueBus workers and subscribe to the exact same events, so a single action can prompte any number of orchestrated actions across any number of servers.

## Installation

To install gosubscriber, use

```sh
go get github.com/delectable/gosubscriber
```

to install the package, and then from your worker

```go
import "github.com/delectable/gosubscriber"
```

## Getting Started

To create a worker, write a function matching the signature

```go
func(args map[string]interface{}) error
```

and subscribe it to an event using:

```go
gosubscriber.Subscribe("my_application", "my_queue", "my_event", mySubscriber)
```

Here is a simple subscriber that subscribes to the event `testEventOne` and prints its arguments:

```go
package main

import (
  "fmt"
  "github.com/delectable/gosubscriber"
)

func testSubscriber(args map[string]interface{}) error {
  fmt.Printf("testSubscriber running with args: %v\n", args)
  return nil
}

func init() {
  gosubscriber.Subscribe("example_application", "example_queue", testSubscriber, map[string]string{
    "bus_event_type": "testEventOne",
  })
}

func main() {
  if err := gosubscriber.Work(); err != nil {
    fmt.Println("Error:", err)
  }
}
```

Here is a slightly more complex subscriber that subscribes to the event `testEventTwo` and prints its arguments, but this example requires the argument `required` to be present:

```go
package main

import (
  "fmt"
  "github.com/delectable/gosubscriber"
)

func testSubscriber(args map[string]interface{}) error {
  fmt.Printf("testSubscriber running with args: %v\n", args)
  return nil
}

func init() {
  gosubscriber.Subscribe("example_application", "example_queue", testSubscriber, map[string]string{
    "required":       gosubscriber.SpecialValues.Present,
    "bus_event_type": "testEventTwo",
  })
}

func main() {
  if err := gosubscriber.Work(); err != nil {
    fmt.Println("Error:", err)
  }
}

```

gosubscriber subscribers receive the arguments sent over ResqueBus as a single map of interfaces. To use them as parameters to other functions, use Go type assertions to convert them into usable types.

``` go
// where doSomething expects (int64, string)
func testSubscriber(args map[string]interface{}) error {
  idNum, ok := args["id"].(json.Number)
  if !ok {
    return errorInvalidParam
  }
  id, err := idNum.Int64()
  if err != nil {
    return errorInvalidParam
  }

  name, ok := args["name"].(string)
  if !ok {
    return errorInvalidParam
  }
  doSomething(id, name)
  return nil
}
```

For testing, it's helpful to use IRB to publish events (note that a ResqueBus Driver must be running)

``` ruby
ResqueBus.publish(:testEventOne)
ResqueBus.publish(:testEventTwo, {required: 1})
```

For information on [configuration/flags](https://github.com/benmanns/goworker#flags), [signal handling](https://github.com/benmanns/goworker#signal-handling-in-goworker), and [failure modes](https://github.com/benmanns/goworker#failure-modes), see the [goworker readme](https://github.com/benmanns/goworker)
