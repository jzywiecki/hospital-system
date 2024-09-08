package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"server/server/internal/connections"
	"server/server/internal/errors"
	"server/server/internal/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn := connections.ConnectToRabbit()
	defer conn.Close()

	ch := connections.CreateChannel(conn)
	defer ch.Close()

	connections.DeclareExchange(ch)
	q := connections.DeclareQueue(ch, "admin")
	connections.BindQueue(ch, "admin", "#", types.ExchangeName)

	ch2 := connections.CreateChannel(conn)
	defer ch2.Close()

	connections.DeclareLogExchange(ch2)

	go func() {
		msgs, err := ch.Consume(
			q.Name, 
			"",     
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

			if messageUnmarshalled.IsSenderADoctor {
				log.Printf(" [x] Doctor send %s for %s to examination!", messageUnmarshalled.ExaminationType, messageUnmarshalled.PatientsName)
			} else {
				log.Printf(" [x] Technician send %s results for %s examination!", messageUnmarshalled.ExaminationType, messageUnmarshalled.PatientsName)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// read input from the user
		text := scanner.Text()

		// publish the input to the exchange
		err := ch2.Publish(
			types.LogExchangeName,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(text),
			},
		)
		errors.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", text)

	}
}
