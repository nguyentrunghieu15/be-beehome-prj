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
}

type ClientWrapper struct {
	config *MomgoClientConfig
	client *mongo.Client
	database *mongo.Database	
	once sync.Once
}
  
func (istn *ClientWrapper) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), istn.config.Timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v@%v",istn.config.UserName,istn.config.Password,istn.config.Address)))
	if err != nil { return err }
	istn.client = client
	return nil
}

func (istn *ClientWrapper) doConnect() {
	for {
		err:= istn.connect()
		if err != nil {
			time.Sleep(300)
			continue
		}
		break
	}
}

func (istn *ClientWrapper) Client()  (*mongo.Client,error) {
	istn.once.Do(istn.doConnect)
	return istn.client ,nil
}

func (istn *ClientWrapper) Db()  (*mongo.Database,error) {
	
}

type Repository[T any] struct{
	Client *ClientWrapper
	Collection string
}


func (istn *Repository[T]) FindOneByAtribute(name string, value interface{}) (*T,error) {
	filter := bson.D{{name, value}}
	result  = new(T)
	err = istn.Client.Db().Collection(istn.Collection).(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result ,nil
}
