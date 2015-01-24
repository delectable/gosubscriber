package goworker_bus

type subscriberFunc func(map[string]interface{}) error
