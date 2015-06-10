package gosubscriber

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/delectable/goworker"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// Configurable variables
var (
	appListKey     = "resquebus_apps"
	appSingleKey   = "resquebus_app"
	specialPrepend = "bus_special_value_"
	logger, _      = seelog.LoggerFromWriterWithMinLevel(os.Stdout, seelog.InfoLvl)
)

// Special Values are for Matcher criteria. For usage example, see example.go
// For Ruby driver implementation, see:
// https://github.com/taskrabbit/resque-bus/blob/master/lib/resque_bus/matcher.rb
type specialValues struct {
	Key     string
	Blank   string
	Nil     string
	Present string
	Empty   string
	Value   string
}

var (
	SpecialValues = specialValues{
		Key:     fmt.Sprintf("%s%s", specialPrepend, "key"),
		Blank:   fmt.Sprintf("%s%s", specialPrepend, "blank"),
		Nil:     fmt.Sprintf("%s%s", specialPrepend, "nil"),
		Present: fmt.Sprintf("%s%s", specialPrepend, "present"),
		Empty:   fmt.Sprintf("%s%s", specialPrepend, "empty"),
		Value:   fmt.Sprintf("%s%s", specialPrepend, "value"),
	}
)

// This is will be the serialized subscription value, stored in redis for the
// Driver to parse
type subscriptionValueStruct struct {
	QueueName string            `json:"queue_name"`
	Key       string            `json:"key"`
	Class     string            `json:"class"`
	Matcher   map[string]string `json:"matcher"`
}

// Takes a function, returns its path
func getFunctionPath(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Takes a function, returns its name (sans path)
func getFunctionName(i interface{}) string {
	slice := strings.Split(getFunctionPath(i), ".")
	return slice[len(slice)-1]
}

// Builds the Hash key under which the serialized subscription value is stored
func buildSubscriptionKey(functionPath, busEventType string) string {
	slice := strings.Split(functionPath, ".")
	pathPrefix := strings.Join(slice[:len(slice)-1], ".")
	functionName := slice[len(slice)-1]

	return fmt.Sprintf(
		"%s.__resquebussubscriber__%s__%s",
		pathPrefix,
		functionName,
		busEventType,
	)
}

// Turns a subscriberFunc into a goworker.workerFunc for registering in Goworker
func wrapSubscriber(subscriber subscriberFunc, functionPath string) func(string, ...interface{}) error {
	return func(queue string, args ...interface{}) error {
		logger.Debugf(
			"Bus Subscriber '%s' activated for queue '%s' and event '%s'",
			functionPath,
			queue,
			args[0].(map[string]interface{})["bus_event_type"],
		)
		return subscriber(args[0].(map[string]interface{}))
	}
}

func buildRedisKey(application string) string {
	return fmt.Sprintf(
		"%s%s:%s",
		goworker.Namespace(),
		appSingleKey,
		application,
	)
}

// Performs all necessary functions to subscribe a function to a ResquBus event
func Subscribe(application string, queueName string, subscriber subscriberFunc, matcher map[string]string) {
	// Set a default bus_event_type
	if matcher["bus_event_type"] == "" {
		if matcher == nil {
			matcher = make(map[string]string)
		}
		matcher["bus_event_type"] = getFunctionName(subscriber)
	}

	// Build the Redis Key to store this application's subscriptions
	redisKey := buildRedisKey(application)

	functionPath := getFunctionPath(subscriber)
	subscriptionKey := buildSubscriptionKey(functionPath, matcher["bus_event_type"])

	className := fmt.Sprintf("%s-%s", functionPath, matcher["bus_event_type"])

	logger.Info("Subscribing.")
	logger.Info("  Application:      ", application)
	logger.Info("  Function:         ", functionPath)
	logger.Info("  Queue Name:       ", queueName)
	logger.Info("  Matcher:          ", matcher)
	logger.Info("  ClassName:        ", className)
	logger.Info("  Redis Key:        ", redisKey)
	logger.Info("  Subscription Key: ", subscriptionKey)

	// Initialize goworker, so we can use its Redis connection
	if err := goworker.Init(); err != nil {
		logger.Critical("ERROR IN INIT:", err)
	}
	defer goworker.Close() // tear down goworker once we're done here

	conn, err := goworker.GetConn() // pull a conn from goworker's connection pool

	if err != nil {
		logger.Critical("ERROR GETTING GOWORKER CONNECTION:", err)
	}

	// subscriptionValue contains this subscription's configuration, stored as a
	// serialized string in the 'redisKey' Redis hash under 'subscriptionKey'
	subscriptionValue := subscriptionValueStruct{
		QueueName: queueName,
		Key:       subscriptionKey,
		Class:     className,
		Matcher:   matcher,
	}

	serializedSubscriptionValue, err := json.Marshal(subscriptionValue)

	if err != nil {
		logger.Critical("ERROR SERIALIZING: ", err)
	} else {
		logger.Debug("  Serialized Subscription Value: ", string(serializedSubscriptionValue))
	}

	// Store the subscription configuration (pipelined)
	err = conn.Send("HSET", redisKey, subscriptionKey, serializedSubscriptionValue)

	if err != nil {
		logger.Critical("ERROR IN HSET:", err)
	}

	// Ensure that this application is registered, so the ResqueBus driver can
	// find our subscriptions (pipelined)
	err = conn.Send("SADD", fmt.Sprintf("%s%s", goworker.Namespace(), appListKey), application)

	conn.Flush() // Finalize the pipeline

	goworker.PutConn(conn) // return the Redis connection back to Goworker

	// Register the subscriberFunc with goworker, wrapped as a goworker.workerFunc
	goworker.Register(className, wrapSubscriber(subscriber, functionPath))
}

// Removes an entire application from ResqueBus
func Unsubscribe(application string) {
	// Initialize goworker, so we can use its Redis connection
	if err := goworker.Init(); err != nil {
		logger.Critical("ERROR IN INIT:", err)
	}
	defer goworker.Close()          // tear down goworker once we're done here
	conn, err := goworker.GetConn() // pull a conn from goworker's connection pool

	err = conn.Send("SREM", fmt.Sprintf("%s%s", goworker.Namespace(), appListKey), application)
	if err != nil {
		logger.Critical("ERROR IN SREM:", err)
	}

	err = conn.Send("DEL", buildRedisKey(application))
	if err != nil {
		logger.Critical("ERROR IN HSET:", err)
	}

	conn.Flush() // Finalize the pipeline

	goworker.PutConn(conn) // return the Redis connection back to Goworker
}

func Work() error {
	logger.Info("Working.")
	return goworker.Work()
}
