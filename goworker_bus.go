package goworker_bus

import (
	"encoding/json"
	"fmt"
	"github.com/benmanns/goworker"
	"reflect"
	"runtime"
	"strings"
)

var (
	appListKey     = "resquebus_apps"
	appSingleKey   = "resquebus_app"
	specialPrepend = "bus_special_value_"
)

type subscriptionValueStruct struct {
	QueueName string            `json:"queue_name"`
	Key       string            `json:"key"`
	Class     string            `json:"class"`
	Matcher   map[string]string `json:"matcher"`
}

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

func getFunctionPath(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func getFunctionName(i interface{}) string {
	slice := strings.Split(getFunctionPath(i), ".")
	return slice[len(slice)-1]
}

func buildSubscriptionKey(functionPath, busEventType string) string {
	slice := strings.Split(functionPath, ".")
	pathPrefix := strings.Join(slice[:len(slice)-1], ".")
	functionName := slice[len(slice)-1]

	return fmt.Sprintf("%s.__resquebussubscriber__%s__%s", pathPrefix, functionName, busEventType)
}

// Subscribe registers a goworker worker function and
// subscribes it to a given bus event type
func Subscribe(application string, queueName string, callback subscriberFunc, matcher map[string]string) {
	if matcher["bus_event_type"] == "" {
		if matcher == nil {
			matcher = make(map[string]string)
		}
		matcher["bus_event_type"] = getFunctionName(callback)
	}

	redisKey := fmt.Sprintf("%s%s:%s", goworker.Namespace(), appSingleKey, application)
	subscriptionKey := buildSubscriptionKey(getFunctionPath(callback), matcher["bus_event_type"])

	fmt.Println("SUBSCRIBING")
	fmt.Println("  Application:", application)
	fmt.Println("  Queue Name: ", queueName)
	fmt.Println("  Matcher:    ", matcher)
	fmt.Println("  Redis Key:  ", redisKey)
	fmt.Println("  Subscription Key:  ", subscriptionKey)

	if err := goworker.Init(); err != nil {
		fmt.Println("ERROR IN INIT:", err)
	}
	defer goworker.Close()
	conn, err := goworker.GetConn()

	if err == nil {
		subscriptionValue := subscriptionValueStruct{
			QueueName: queueName,
			Key:       subscriptionKey,
			Class:     queueName,
			Matcher:   matcher,
		}

		serializedSubscriptionValue, err := json.Marshal(subscriptionValue)

		if err != nil {
			fmt.Println("ERROR IN SERIALIZATION:", err)
		} else {
			fmt.Println("serializedSubscriptionValue", string(serializedSubscriptionValue))
		}

		err = conn.Send("HSET", redisKey, subscriptionKey, serializedSubscriptionValue)

		if err != nil {
			fmt.Println("ERROR IN SEND:", err)
		}
		// ResqueBus.redis.sadd(self.class.appListKey, app_key)
		err = conn.Send("SADD", fmt.Sprintf("%s%s", goworker.Namespace(), appListKey), application)

		conn.Flush()
		goworker.PutConn(conn)

		goworker.Register(queueName, wrapSubscriber(callback))
	} else {
		fmt.Println("ERROR IN GETCONN:", err)
	}
}

func wrapSubscriber(subscriber subscriberFunc) func(string, ...interface{}) error {
	return func(queue string, args ...interface{}) error {
		fmt.Printf("Bus Subscriber %s activated for queue %s\n", getFunctionName(subscriber), queue)
		return subscriber(args)
	}
}
