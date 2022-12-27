build:
	go build -o bin/consumer sconsumer.go

run: build
	KAFKA_PEERS=localhost:9092 KAFKA_TOPIC=t3 KAFKA_VERSION=2.7.0 KAFKA_GROUP=orders  ./bin/consumer
