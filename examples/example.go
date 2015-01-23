package main

import (
	"fmt"
	"github.com/benmanns/goworker"
	"github.com/delectable/goworker-bus"
)

func testSubscriber(args ...interface{}) error {
	fmt.Printf("From %v\n", args)
	return nil
}

func testWorker(queue string, args ...interface{}) error {
	fmt.Printf("From %s, %v\n", queue, args)
	return nil
}

func main() {
	goworker_bus.Subscribe("delectaroutes", "test", testSubscriber, map[string]string{
		// "arg":            goworker_bus.SpecialValues.Present,
		"bus_event_type": "test",
	})
	goworker.Work()
}
