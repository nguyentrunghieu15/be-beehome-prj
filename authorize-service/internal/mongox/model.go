package mongox


type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserId string `bson:"user_id"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserId string `bson:"user_id"`
}

type Provider struct {
	ID           primitive.ObjectID `bson:"_id"`
	ProviderId string `bson:"provider_id"`
	UserId string `bson:"user_id"`
}

type Hire struct {
	ID           primitive.ObjectID `bson:"_id"`
	HireId string `bson:"hire_id"`
	ProviderId string `bson:"provider_id"`
	UserId string `bson:"user_id"`
}

type GroupService struct {
	ID           primitive.ObjectID `bson:"_id"`
	GroupServiceId string `bson:"group_service_id"`
}

type Service struct {
	ID           primitive.ObjectID `bson:"_id"`
	GroupServiceId string `bson:"group_service_id"`
	ServiceId string `bson:"service_id"`
}

type SocialMedia struct {
	ID           primitive.ObjectID `bson:"_id"`
	SocialMediaId string `bson:"social_media_id"`
	ProviderId string `bson:"provider_id"`
}

type PaymentMethod struct {
	ID           primitive.ObjectID `bson:"_id"`
	PaymentMethodId string `bson:"payment_method_id"`
	ProviderId string `bson:"provider_id"`
}


