gosubscriber
============
gosubscriber is a [Resque Bus](https://github.com/taskrabbit/resque-bus)-compatible, Go-based background subscriber build on top of Benjamin Manns' [goworker](https://github.com/benmanns/goworker). It allows you to publish Resque Bus events in Ruby (or any other language), and perform background tasks in Go.

Gosubscriber workers can run alongside Ruby/Node.js ResqueBus subscribers and subscribe to the exact same events, so that a single action can prompte any number of orchestrated actions across any number of servers.
