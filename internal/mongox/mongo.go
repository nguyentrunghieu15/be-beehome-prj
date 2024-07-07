package mongox

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MomgoClientConfig struct {
	DatabaseName string
	Username     string
	Password     string
	Address      string
	Uri          string
	Timeout      time.Duration
}

type ClientWrapper struct {
	config   *MomgoClientConfig
	client   *mongo.Client
	database *mongo.Database
	once     sync.Once
}

var DefaultClient *ClientWrapper

func NewClientMongoWrapperWithConfig(config *MomgoClientConfig) *ClientWrapper {
	return &ClientWrapper{
		config: config,
	}
}

func (istn *ClientWrapper) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), istn.config.Timeout)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%v:%v@%v", istn.config.Username, istn.config.Password, istn.config.Address)),
	)
	if err != nil {
		return err
	}
	istn.client = client
	return nil
}

func (istn *ClientWrapper) doConnect() {
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

func (istn *ClientWrapper) Client() *mongo.Client {
	istn.once.Do(istn.doConnect)
	return istn.client
}

func (istn *ClientWrapper) Db() *mongo.Database {
	istn.once.Do(istn.doConnect)
	istn.database = istn.client.Database(istn.config.DatabaseName)
	return istn.database
}

type Repository[T any] struct {
	Client     *ClientWrapper
	Collection string
}

func (istn *Repository[T]) InsertOne(value interface{}) error {
	_, err := istn.Client.Db().Collection(istn.Collection).InsertOne(context.Background(), value)
	return err
}

func (istn *Repository[T]) FindOneByAtribute(name string, value interface{}) (*T, error) {
	filter := bson.D{{Key: name, Value: value}}
	result := new(T)
	err := istn.Client.Db().Collection(istn.Collection).FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (istn *Repository[T]) FindAllByAtribute(name string, value interface{}) (*[]T, error) {
	filter := bson.D{{Key: name, Value: value}}
	results := make([]T, 0)
	cursor, err := istn.Client.Db().Collection(istn.Collection).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	// Unpacks the cursor into a slice
	if err = cursor.All(context.TODO(), results); err != nil {
		return nil, err
	}
	return &results, nil
}

func (istn *Repository[T]) DeleteOneByAtribute(name string, value interface{}) error {
	filter := bson.D{{Key: name, Value: value}}
	_, err := istn.Client.Db().Collection(istn.Collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
