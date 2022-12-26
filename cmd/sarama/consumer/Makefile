build:
	go build -o consumer sconsumer.go

run: build
	KAFKA_PEERS=localhost:9092 KAFKA_TOPIC=orders KAFKA_VERSION=2.7.0 KAFKA_GROUP=orders  ./consumer
