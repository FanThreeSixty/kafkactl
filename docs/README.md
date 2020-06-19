# Helpful Kafka Commands

The following are helpful Kafka Commands that aren't yet associated or built into the kafkactl command line utility.

## Commands

### Topic Commands

Alter partitions on an existing topic:

```sh
./kafka-topics --zookeeper zookeeper:port --topic topic --alter --partitions 5
```

## Resources

https://github.com/fgeller/kt
