package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
)

var kafkaBrokers = os.Getenv("KAFKA_PEERS")
var kafkaTopic = os.Getenv("KAFKA_TOPIC")
var kafkaPartition = os.Getenv("KAFKA_PARTITION")
var partition int32 = -1
var globalProducer sarama.SyncProducer

func main() {
	fmt.Println("booting up producer")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	p, err := strconv.Atoi(kafkaPartition)

	if err != nil {
		fmt.Println("Failed to convert KAFKA_PARTITION to Int32")
		panic(err)
	}

	partition = int32(p)

	producer, err := sarama.NewSyncProducer(strings.Split(kafkaBrokers, ","), config)
	if err != nil {
		fmt.Printf("Failed to open Kafka producer: %s", err)
		panic(err)
	}

	globalProducer = producer

	defer func() {
		fmt.Println("Closing Kafka producer...")
		if err := globalProducer.Close(); err != nil {
			fmt.Printf("Failed to close Kafka producer cleanly: %s", err)
			panic(err)
		}
	}()

	router := httprouter.New()

	router.POST("/publish/:message", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		submit(w, r, p)
	})

	fmt.Println("Running... http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}

type Employee struct {
	Name  string
	Age   int
	Email string
}

func submit(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {

	messageValue := p.ByName("message")
	message := &sarama.ProducerMessage{Topic: kafkaTopic, Partition: partition}
	e := Employee{Name: "najam awan", Email: "najamsk@gmail.com", Age: 33}
	eStr, _ := json.Marshal(e)
	// message.Value = sarama.StringEncoder(messageValue)
	message.Value = sarama.StringEncoder(eStr)
	message.Key = sarama.StringEncoder("1119")

	fmt.Println("Received message: " + messageValue)
	fmt.Println("msg to kafka:", eStr)

	partition, offset, err := globalProducer.SendMessage(message)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to Kafka", err)
	}

	fmt.Printf("publish success! topic=%s\tpartition=%d\toffset=%d\n", kafkaTopic, partition, offset)

}
