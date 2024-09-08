package connections

import (
	"log"
	"math/rand"
	"server/server/internal/errors"
	"server/server/internal/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbit() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	errors.FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	errors.FailOnError(err, "Failed to open a channel")
	return ch
}

func DeclareExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		types.ExchangeName,
		"topic",
		false,
		true,
		false,
		false,
		nil,
	)
	errors.FailOnError(err, "Failed to declare an exchange")

	log.Printf("Declared exchange %s", types.ExchangeName)
}

func DeclareLogExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		types.LogExchangeName,
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)
	errors.FailOnError(err, "Failed to declare an exchange")

	log.Printf("Declared exchange %s", types.LogExchangeName)
}

func DeclareQueue(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	errors.FailOnError(err, "Failed to declare a queue")
	return q
}

func BindQueue(ch *amqp.Channel, name string, key string, exchangeName string) {
	err := ch.QueueBind(
		name,
		key,
		exchangeName,
		false,
		nil,
	)

	errors.FailOnError(err, "Failed to bind a queue")
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
