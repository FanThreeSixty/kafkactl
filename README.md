# kafkactl

A Swiss Army CLI tool for producing, consuming and otherwise interacting with your Kafka Cluster.

```sh
kafkactl is a command line script to interact with (a) Kafka client(s)/cluster(s)

Usage:
  kafkactl [command]

Available Commands:
  consume     consume commands for Kafka
  describe    command for describing resources in Kafka.
  get         command for getting resources in Kafka.
  health      Check health of configured Kafka Cluster
  help        Help about any command
  manage      command for managing resources in Kafka.
  produce     producer commands for Kafka
  version     Display version information for Kafkactl

Flags:
  -b, --brokers stringSlice   Space delimited list of Kafka Brokers to connect to. (default [127.0.0.1:9092])
  -e, --environment string    Change configuration based on Environment.
  -h, --help                  help for kafkactl
  -v, --verbose               Enable verbose debugging of script.

Use "kafkactl [command] --help" for more information about a command.
```

## Setup

kafkactl utilizes the github.com/spf13/viper package for configuring the target Kafka cluster to interact with.  Because of this we are able to use three different ways of configuring the kafkactl script.

**FILE** based configuration is the lowest priority of the three.

*Location*: $HOME/.kafkactl/config.yaml, /etc/kafkactl/config.yaml, $PWD/config.yaml

Below is an example of the configuration one can have while using the **FILE** configuration:

```yaml
default:
  brokers:
    - localhost:9092

development:
  brokers:
    - development-01:9092
    - development-02:9092
```

Notice there is a 'brokers' at the top level, as well as 'brokers' underneath a dictionary key value.  This is controlled by using the '-e|--environment' inline argument.  If no '-e|--environment' is used, the default 'brokers' value will be used.  If a '-e|--environment' argument is passed in, it will look for said dictionary key value match and load in the brokers as it finds a match.  This is meant to allow for multiple configurations under one file without having to point to different configs.

```
# examples:
$ kafkactl health
{"brokers":[{"endpoint":"localhost:9092","endpoint_addr":"localhost","healthy":false,"message":""}],"healthy":true}
$ kafkactl -e local health
{"brokers":[{"endpoint":"localhost:9092","endpoint_addr":"localhost","healthy":false,"message":""}],"healthy":true}
```

**ENV** based configuration is the second priority and takes precedence over **FILE**.

**INLINE** arguments are the highest priority and take precedence over all other configurations.

## Building

```sh
cd $GOPATH
go get github.com/FanThreeSixty/kafkactl
go install github.com/FanThreeSixty/kafkactl
```
or

```sh
git clone github.com/FanThreeSixty/kafkactl
go build -o kafkactl .
```

**Building for Linux:**

```sh
git clone github.com/FanThreeSixty/kafkactl
GOOS=linux go build -o kafkactl .
```

# Commands

**List Available Topics:**

```sh
$ kafkactl get topics
TOPIC                      	NUM OF PARTITIONS
ACTIVITY_FAILURE           	1
...
```

**List Consumer Groups:**

```sh
$ kafkactl get groups
GROUP                                  	MEMBERS 	COORDINATOR.ID	COORDINATOR.ADDR
consumerGroup     	                    consumer	1             	dev-kafka-01.local.com:9092
```

**Manage Consumer Groups:**

*Fast Forward Offset (catch consumer group up to most recent message in Kafka):*

```sh
kafkactl manage offsets (fastforward|rewind) -g consumerGroup -t TOPIC
```


**Querying the Health of a Cluster:**

```sh
$ kafkactl health cluster
{"brokers":[{"endpoint":"dev-kafka-01.local.com:9092","endpoint_addr":"dev-kafka-01.local.com:9092","healthy":true,"message":""}],"healthy":true}

# or use JQ....
$ kafkactl health cluster | jq ''
{
  "brokers": [
    {
      "endpoint": "dev-kafka-01.local.com:9092",
      "endpoint_addr": "dev-kafka-01.local.com:9092",
      "healthy": true,
      "message": ""
    }
  ],
  "healthy": true
}
```

**Querying Kafka Content Metrics:**

```sh
$ kafkactl get metrics | jq ''
{
  "status": "ok",
  "message": "",
  "groups": {
    "consumerGroup": {
      "name": "consumerGroup",
      "members": "consumer",
      "coordinator_id": 1,
      "coordinator_addr": "dev-kafka-01.local.com",
      "topics": [
        {
          "name": "ACTIVITY_FAILURE",
          "partitions": [
            {
              "id": 0,
              "leader_addr": "dev-kafka-01.local.com:9092",
              "leader_id": 1,
              "offsets": {
                "group": 2,
                "log_end": 2,
                "lag": 0
              }
            }
          ]
        }
      ]
    }
  }
}
     
```

# Resources

* https://github.com/spf13/viper
* https://github.com/spf13/cobra
* https://github.com/Shopify/sarama
* https://github.com/fgeller/kt
