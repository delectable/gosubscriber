package goworker_bus

type subscriberFunc func(...interface{}) error
