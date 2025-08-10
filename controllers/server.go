package controllers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Server struct {
	DB   *pgxpool.Pool
	amqp *amqp.Channel
	q    amqp.Queue
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{DB: db}
}

func (s *Server) InitRabbit(ctx context.Context, amqpURL, queueName string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	s.amqp = ch
	s.q = q
	return nil
}

type TrxEvent struct {
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
	Source        string `json:"source"`
	Timestamp     string `json:"timestamp"`
}

func (s *Server) PublishTrxEvent(ctx context.Context, ev TrxEvent) error {
	if s.amqp == nil {
		return nil
	}
	body, _ := json.Marshal(ev)
	return s.amqp.PublishWithContext(
		ctx,
		"",
		s.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
			Type:         "trx.event",
		},
	)
}
