package kafkax

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClientConfig struct {
	Topic            string
	BooststrapServer string
	Partition        int
	Protocall        string
	TimeoutRead      time.Duration
	TimeoutWrite     time.Duration
	MaxBytes         int64
	GroupId          string
}

type KafkaClientWrapper struct {
	config *KafkaClientConfig
	conn   *kafka.Conn
	reader *kafka.Reader
	writer *kafka.Writer
	once   sync.Once
}

func (inst *KafkaClientWrapper) connect() error {
	conn, err := kafka.DialLeader(
		context.Background(),
		inst.config.Protocall,
		inst.config.BooststrapServer,
		inst.config.Topic,
		inst.config.Partition,
	)
	if err != nil {
		return err
	}

	rConfig := kafka.ReaderConfig{
		Brokers:   []string{inst.config.BooststrapServer},
		Topic:     inst.config.Topic,
		Partition: inst.config.Partition,
		MaxBytes:  int(inst.config.MaxBytes), // 10MB
	}
	if inst.config.GroupId != "" {
		rConfig.GroupID = inst.config.GroupId
	}
	r := kafka.NewReader(rConfig)

	w := &kafka.Writer{
		Addr:         kafka.TCP(inst.config.BooststrapServer),
		Topic:        inst.config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	inst.conn = conn
	inst.reader = r
	inst.writer = w
	return nil
}

func (istn *KafkaClientWrapper) doConnect() {
	for {
		err := istn.connect()
		if err != nil {
			log.Println(err)
			time.Sleep(300)
			continue
		}
		break
	}
}

func (inst *KafkaClientWrapper) Reader() *kafka.Reader {
	inst.once.Do(inst.doConnect)
	return inst.reader
}

func (inst *KafkaClientWrapper) Writer() *kafka.Writer {
	inst.once.Do(inst.doConnect)
	return inst.writer
}

func (inst *KafkaClientWrapper) Conn() *kafka.Conn {
	inst.once.Do(inst.doConnect)
	return inst.conn
}

func (inst *KafkaClientWrapper) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if inst.reader == nil {
		return kafka.Message{}, errors.New("no connection to kafka")
	}
	return inst.reader.ReadMessage(ctx)
}

func (inst *KafkaClientWrapper) WriteMessages(ctx context.Context, msg ...kafka.Message) error {
	if inst.writer == nil {
		return errors.New("no connection to kafka")
	}
	err := inst.writer.WriteMessages(context.Background(),
		msg...,
	)
	return err
}

func (inst *KafkaClientWrapper) Close() error {
	err := inst.conn.Close()
	if err != nil {
		return err
	}
	err = inst.reader.Close()

	if err != nil {
		return err
	}
	err = inst.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

var DefaultKafkaClient *KafkaClientWrapper

func NewKafkaClientWrapperWithConfig(config *KafkaClientConfig) *KafkaClientWrapper {
	return &KafkaClientWrapper{
		config: config,
	}
}
