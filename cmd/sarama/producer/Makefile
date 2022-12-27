build:
	go build -o bin/producer producer.go

run: build
	KAFKA_PEERS=localhost:9092 KAFKA_TOPIC=t3 KAFKA_PARTITION=1 KAFKA_GROUP=orders ./bin/producer

post:
	curl -X POST  http://localhost:9000/publish/curlTestData
