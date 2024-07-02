package mongox

import (
	"time"
	"mongo"
)


type MomgoClientConfig struct {
	DatabaseName string
	Username string
	Password string
	Address string
	Uri string
	Timeout time.Duration
}filter := bson.D{{"name", "Bagels N Buns"}}

type ClientWrapper struct {
	config *MomgoClientConfig
	client *mongo.Client
	database *mongo.Database	
	once sync.Once
}
  
func (istn *ClientWrapper) connect(){

}

func (istn *ClientWrapper) Client()  (*mongo.Client,error) {

}

func (istn *ClientWrapper) Db()  (*mongo.Database,error) {
	
}

type Repository[T any] struct{
	Client ClientWrapper
	Collection string
}


func (istn *Repository[T]) FindOneByAtribute(name string, value interface{}) (*T,error) {
	filter := bson.D{{name, value}}

	result  = new(T)
	err = istn.Client.Db().Collection(istn.Collection).(context.TODO(), filter).Decode(&result)
	// Prints a message if no documents are matched or if any
	// other errors occur during the operation
	if err != nil {
		return nil, err
	}
	return result ,nil
}

const DefaultClientWrapper := new(ClientWrapper)