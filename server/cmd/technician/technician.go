package main

import (
	"encoding/json"
	"log"
	"os"
	"server/server/internal/connections"
	"server/server/internal/errors"
	"server/server/internal/types"
	"strings"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [injury_type] [injury_type]", os.Args[0])
		os.Exit(0)
	}

	conn := connections.ConnectToRabbit()
	defer conn.Close()

	ch := connections.CreateChannel(conn)
	defer ch.Close()

	connections.DeclareExchange(ch)

	connections.DeclareLogExchange(ch)
	qLog := connections.DeclareQueue(ch, "")
	connections.BindQueue(ch, qLog.Name, "", types.LogExchangeName)

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

	for _, injuryType := range os.Args[1:] {
		technicianId := uuid.New().String()
		injuryType = strings.ToLower(injuryType)

		q := connections.DeclareQueue(ch, injuryType)
		connections.BindQueue(ch, q.Name, "technician."+injuryType, types.ExchangeName)

		msgs, err := ch.Consume(
			q.Name,
			technicianId,
			true,
			false,
			false,
			false,
			nil,
		)
		errors.FailOnError(err, "Failed to register a consumer")

		go func() {
			for d := range msgs {
				var recievedMessage types.Examination

				err := json.Unmarshal(d.Body, &recievedMessage)
				errors.FailOnError(err, "Failed to unmarshal message")

				log.Printf(" [x] Got %s for %s to examine!", recievedMessage.ExaminationType, recievedMessage.PatientsName)

				messageToSend := types.Examination{
					ExaminationType: recievedMessage.ExaminationType,
					PatientsName:    recievedMessage.PatientsName,
					Sender:          technicianId,
					IsSenderADoctor: false,
				}

				messageToSendMarshalled, err := json.Marshal(messageToSend)
				errors.FailOnError(err, "Failed to marshal message")

				err = ch.Publish(
					types.ExchangeName,
					"doctor."+recievedMessage.Sender,
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(messageToSendMarshalled),
					},
				)
				errors.FailOnError(err, "Failed to publish a message")
				log.Printf(" [x] Sent %s results for %s to doctor!", messageToSend.ExaminationType, messageToSend.PatientsName)
			}
		}()
	}

	// keep alive the main thread
	var forever chan struct{}
	<-forever
}
