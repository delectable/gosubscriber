package main

import (
	"fmt"
	"github.com/delectable/goworker-bus"
)

func testSubscriber(args map[string]interface{}) error {
	fmt.Printf("testSubscriber running with args: %v\n", args)
	return nil
}

func main() {
	// Subscribing to testEventOne
	goworker_bus.Subscribe("delectaroutes", "test", testSubscriber, map[string]string{
		"bus_event_type": "testEventOne",
	})

	// Publishing to testEventOne from Ruby:
	// ResqueBus.publish(:testEventOne)

	// Subscribing to testEventTwo, this time requiring the "required" argument
	// to be present
	goworker_bus.Subscribe("delectaroutes", "test", testSubscriber, map[string]string{
		"required":       goworker_bus.SpecialValues.Present,
		"bus_event_type": "testEventTwo",
	})

	// Publishing to testEventTwo from Ruby and testing the required argument
	// ResqueBus.publish(:testEventTwo, {required: 1})   // works
	// ResqueBus.publish(:testEventTwo, {blarg: 1})      // doesn't work

	// Firing up the worker
	goworker_bus.Work()
}
