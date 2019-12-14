package main

import (
	"flag"
	"fmt"
	"github.com/cyrilix/robocar-base/cli"
	part2 "github.com/cyrilix/robocar-led/part"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

const (
	DefaultClientId = "robocar-led"
)

func main() {
	var mqttBroker, username, password, clientId string
	var driveModeTopic, recordTopic string

	mqttQos := InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	funcInitMqttFlags(DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)

	flag.StringVar(&driveModeTopic, "mqtt-topic-drive-mode", os.Getenv("MQTT_TOPIC_DRIVE_MODE"), "Mqtt topic that contains DriveMode value, use MQTT_TOPIC_DRIVE_MODE if args not set")
	flag.StringVar(&recordTopic, "mqtt-topic-record", os.Getenv("MQTT_TOPIC_RECORD"), "Mqtt topic that contains video recording state, use MQTT_TOPIC_RECORD if args not set")

	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := Connect(mqttBroker, username, password, clientId)
	if err != nil {
		log.Fatalf("unable to connect to mqtt bus: %v", err)
	}
	defer client.Disconnect(50)

	p := part2.NewPart(client, driveModeTopic, recordTopic)
	defer p.Stop()

	err = p.Start()
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}

func funcInitMqttFlags(defaultClientId string, mqttBroker, username, password, clientId *string, mqttQos *int, mqttRetain *bool) {
	cli.SetDefaultValueFromEnv(clientId, "MQTT_CLIENT_ID", defaultClientId)
	cli.SetDefaultValueFromEnv(mqttBroker, "MQTT_BROKER", "tcp://127.0.0.1:1883")

	flag.StringVar(mqttBroker, "mqtt-broker", *mqttBroker, "Broker Uri, use MQTT_BROKER env if arg not set")
	flag.StringVar(username, "mqtt-username", os.Getenv("MQTT_USERNAME"), "Broker Username, use MQTT_USERNAME env if arg not set")
	flag.StringVar(password, "mqtt-password", os.Getenv("MQTT_PASSWORD"), "Broker Password, MQTT_PASSWORD env if args not set")
	flag.StringVar(clientId, "mqtt-client-id", *clientId, "Mqtt client id, use MQTT_CLIENT_ID env if args not set")
	flag.IntVar(mqttQos, "mqtt-qos", *mqttQos, "Qos to pusblish message, use MQTT_QOS env if arg not set")
	flag.BoolVar(mqttRetain, "mqtt-retain", *mqttRetain, "Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set")
}

func InitIntFlag(key string, defValue int) int {
	var value int
	err := cli.SetIntDefaultValueFromEnv(&value, key, defValue)
	if err != nil {
		log.Panicf("invalid int value: %v", err)
	}
	return value
}

func Connect(uri, username, password, clientId string) (MQTT.Client, error) {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(uri)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(
		//define a function for the default message handler
		func(client MQTT.Client, msg MQTT.Message) {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		})

	//create and start a client using the above ClientOptions
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("unable to connect to mqtt bus: %v", token.Error())
	}
	return client, nil
}
