# kafkactl Tutorial

This tutorial is intended to get one started using kafkactl for its basic purposes:

* downloading and configuring two clusters
* producing to a topic
* consuming a topic
* consuming a topic with a consumer group
* monitoring offset/lag of a consumer group
* exporting metrics
* fast-forwarding a consumer group's offset on a particular topic

<a name='#setup'></a>
## Setting up the Environment

First thing is first, download the release that is appropriate for your system from the [release page for kafkactl](https://github.com/sporting-innovations/kafkactl/releases).  

> ie:
>  
> ```sh
> curl https://github.com/sporting-innovations/kafkactl/releases/download/v0.6.0/kafkactl-0.6.0-darwin -o /usr/local/bin/kafkactl
> chmod +x /usr/local/bin/kafkactl
> ```

Once downloaded you can either use the command line arguments each request, or you can create a configuration file under:

* ~/.kafkactl/config.yaml
* /etc/kafkactl/config.yaml
* $PWD/config.yaml

Here is an example of a configuration file:

```sh
version: V0_10_2_0
brokers:
  - localhost:9092
groups:
  kafkactl:
    - KAFKA_LOADTESTING
    - TRIGGER_LOCATION

external:
  version: V0_10_2_0
  brokers:
    - external.com:9092
  groups:
    consumer_group_1:
      - cg_1_topic_1
    consumer_group_2:
      - cg_2_topic_1
```

From the above configuration we are now able to hit two Kafka clusters.  One, by default (being the localhost:9092) and the other by passing in an argument for environment ``-e external`` (being the external.com:9092)