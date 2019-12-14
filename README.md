# robocar-led

Microservice part to manage leds
     
## Usage
`rc-led <OPTIONS>`

  -mqtt-broker string
        Broker Uri, use MQTT_BROKER env if arg not set (default "tcp://127.0.0.1:1883")
  -mqtt-client-id string
        Mqtt client id, use MQTT_CLIENT_ID env if args not set (default "robocar-led")
  -mqtt-password string
        Broker Password, MQTT_PASSWORD env if args not set
  -mqtt-qos int
        Qos to pusblish message, use MQTT_QOS env if arg not set
  -mqtt-retain
        Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set
  -mqtt-topic-drive-mode string
        Mqtt topic that contains DriveMode value, use MQTT_TOPIC_DRIVE_MODE if args not set
  -mqtt-topic-record string
        Mqtt topic that contains video recording state, use MQTT_TOPIC_RECORD if args not set
  -mqtt-username string
        Broker Username, use MQTT_USERNAME env if arg not set

## Docker build

```bash
export DOCKER_CLI_EXPERIMENTAL=enabled
docker buildx build . --platform linux/amd64,linux/arm/7,linux/arm64 -t cyrilix/robocar-led
```
