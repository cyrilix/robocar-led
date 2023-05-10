package main

import (
	"flag"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-led/pkg/part"
	"go.uber.org/zap"
	"log"
	"os"
)

const (
	DefaultClientId = "robocar-led"
)

func main() {
	var mqttBroker, username, password, clientId string
	var driveModeTopic, recordTopic, speedZoneTopic, throttleTopic string
	var enableSpeedZoneMode bool

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	cli.InitMqttFlags(DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)

	flag.StringVar(&driveModeTopic, "mqtt-topic-drive-mode", os.Getenv("MQTT_TOPIC_DRIVE_MODE"), "Mqtt topic that contains DriveMode value, use MQTT_TOPIC_DRIVE_MODE if args not set")
	flag.StringVar(&recordTopic, "mqtt-topic-record", os.Getenv("MQTT_TOPIC_RECORD"), "Mqtt topic that contains video recording state, use MQTT_TOPIC_RECORD if args not set")
	flag.StringVar(&speedZoneTopic, "mqtt-topic-speed-zone", os.Getenv("MQTT_TOPIC_SPEED_ZONE"), "Mqtt topic that contains speed zone, use MQTT_TOPIC_SPEED_ZONE if args not set")
	flag.StringVar(&throttleTopic, "mqtt-topic-throttle", os.Getenv("MQTT_TOPIC_THROTTLE"), "Mqtt topic that contains throttle, use MQTT_TOPIC_THROTTLE if args not set")
	flag.BoolVar(&enableSpeedZoneMode, "enable-speedzone-mode", false, "Enable speed-zone mode")

	logLevel := zap.LevelFlag("log", zap.InfoLevel, "log level")
	flag.Parse()

	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(*logLevel)
	lgr, err := config.Build()
	if err != nil {
		log.Fatalf("unable to init logger: %v", err)
	}
	defer func() {
		if err := lgr.Sync(); err != nil {
			log.Printf("unable to Sync logger: %v\n", err)
		}
	}()
	zap.ReplaceGlobals(lgr)

	client, err := cli.Connect(mqttBroker, username, password, clientId)
	if err != nil {
		zap.S().Fatalf("unable to connect to mqtt bus: %v", err)
	}
	defer client.Disconnect(50)

	mode := part.LedModeBrake
	if enableSpeedZoneMode {
		mode = part.LedModeSpeedZone
	}
	p := part.NewPart(client, driveModeTopic, recordTopic, speedZoneTopic, throttleTopic, mode)
	defer p.Stop()

	cli.HandleExit(p)

	err = p.Start()
	if err != nil {
		zap.S().Fatalf("unable to start service: %v", err)
	}
}
