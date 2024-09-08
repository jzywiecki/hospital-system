package main

import (
	"encoding/json"
	"log"
	"server/server/internal/connections"
	"server/server/internal/errors"
	"server/server/internal/types"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/google/uuid"
)

func main() {
	doctorId := uuid.New().String()
	log.Printf("Doctor ID: %s", doctorId)
	conn := connections.ConnectToRabbit()
	defer conn.Close()

	ch := connections.CreateChannel(conn)
	defer ch.Close()

	connections.DeclareExchange(ch)
	q := connections.DeclareQueue(ch, doctorId)
	connections.BindQueue(ch, q.Name, "doctor."+doctorId, types.ExchangeName)

	connections.DeclareLogExchange(ch)
	qLog := connections.DeclareQueue(ch, "")
	connections.BindQueue(ch, qLog.Name, "", types.LogExchangeName)

	go func() {
		msgs, err := ch.Consume(
			q.Name,
			doctorId,
			true,
			false,
			false,
			false,
			nil,
		)
		errors.FailOnError(err, "Failed to register a consumer")

		for d := range msgs {
			messageUnmarshalled := types.Examination{}

			err := json.Unmarshal(d.Body, &messageUnmarshalled)
			errors.FailOnError(err, "Failed to unmarshal message")

			log.Printf(" [x] Got %s for %s results from examination!", messageUnmarshalled.ExaminationType, messageUnmarshalled.PatientsName)
		}
	}()

	go func() {
		msgs, err := ch.Consume(
			qLog.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		errors.FailOnError(err, "Failed to register a consumer")
		for d := range msgs {
			log.Println(string(d.Body))
		}
	}()

	for {
		examination := types.RandomExamination()
		examination.Sender = doctorId
		time.Sleep(5 * time.Second)

		examinationMarshalled, err := json.Marshal(examination)
		errors.FailOnError(err, "Failed to marshal message")

		err = ch.Publish(
			types.ExchangeName,
			"technician."+examination.ExaminationType,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        examinationMarshalled,
			},
		)
		errors.FailOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s for %s to examination", examination.ExaminationType, examination.PatientsName)
	}
}
