package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserId     string             `bson:"user_id"`
	ProviderId string             `bson:"provider_id"`
	Role       string             `bson:"role"`
}

type Provider struct {
	ID         primitive.ObjectID `bson:"_id"`
	ProviderId string             `bson:"provider_id"`
	UserId     string             `bson:"user_id"`
}

type Hire struct {
	ID         primitive.ObjectID `bson:"_id"`
	HireId     string             `bson:"hire_id"`
	ProviderId string             `bson:"provider_id"`
	UserId     string             `bson:"user_id"`
}

type Review struct {
	ID         primitive.ObjectID `bson:"_id"`
	ReviewId   string             `bson:"review_id"`
	HireId     string             `bson:"hire_id"`
	ProviderId string             `bson:"provider_id"`
	UserId     string             `bson:"user_id"`
}

type GroupService struct {
	ID             primitive.ObjectID `bson:"_id"`
	GroupServiceId string             `bson:"group_service_id"`
}

type Service struct {
	ID             primitive.ObjectID `bson:"_id"`
	GroupServiceId string             `bson:"group_service_id"`
	ServiceId      string             `bson:"service_id"`
}

type SocialMedia struct {
	ID            primitive.ObjectID `bson:"_id"`
	SocialMediaId string             `bson:"social_media_id"`
	ProviderId    string             `bson:"provider_id"`
}

type PaymentMethod struct {
	ID              primitive.ObjectID `bson:"_id"`
	PaymentMethodId string             `bson:"payment_method_id"`
	ProviderId      string             `bson:"provider_id"`
}