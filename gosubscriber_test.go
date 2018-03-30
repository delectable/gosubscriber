package gosubscriber

import (
	"fmt"
	"testing"
	"time"
)

func testSubscriber(args map[string]interface{}) error {
	fmt.Printf("testSubscriber running with args: %v\n", args)
	return nil
}

/*
TODO figure out how to explicitly make the test succeed after 30 seconds of "working"

func TestSubscribing(t *testing.T) {
	Subscribe("test_application", "test_queue", testSubscriber, map[string]string{
		"bus_event_type": "test_event",
	})

	timer := time.NewTimer(30 * time.Second)
	go func() {
		<-timer.C
		t.SkipNow()
	}()

	err := Work()

	if err != nil {
		t.Error("Unable to subscribe")
	}
}
*/
