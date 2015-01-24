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
	application, queue := "example_application", "example_queue"

	// Cleaning out any stale subscriptions
	gosubscriber.Unsubscribe(application)

	// Subscribing to testEventOne
	gosubscriber.Subscribe(application, queue, testSubscriber, map[string]string{
		"bus_event_type": "testEventOne",
	})

	// Publishing to testEventOne from Ruby:
	// ResqueBus.publish(:testEventOne)

	// Subscribing to testEventTwo, this time requiring the "required" argument
	// to be present
	gosubscriber.Subscribe(application, queue, testSubscriber, map[string]string{
		"required":       gosubscriber.SpecialValues.Present,
		"bus_event_type": "testEventTwo",
	})

	// Publishing to testEventTwo from Ruby and testing the required argument
	// ResqueBus.publish(:testEventTwo, {required: 1})   // works
	// ResqueBus.publish(:testEventTwo, {blarg: 1})      // doesn't work
}

func main() {
	// Firing up the worker
	if err := gosubscriber.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
