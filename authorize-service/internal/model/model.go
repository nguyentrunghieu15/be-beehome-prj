package model

type User struct {
	UserId     string `bson:"user_id"`
	ProviderId string `bson:"provider_id"`
	Role       string `bson:"role"`
}

type Provider struct {
	ProviderId string `bson:"provider_id"`
	UserId     string `bson:"user_id"`
}

type Hire struct {
	HireId     string `bson:"hire_id"`
	ProviderId string `bson:"provider_id"`
	UserId     string `bson:"user_id"`
}

type Review struct {
	ReviewId   string `bson:"review_id"`
	HireId     string `bson:"hire_id"`
	ProviderId string `bson:"provider_id"`
	UserId     string `bson:"user_id"`
}

type GroupService struct {
	GroupServiceId string `bson:"group_service_id"`
}

type Service struct {
	GroupServiceId string `bson:"group_service_id"`
	ServiceId      string `bson:"service_id"`
}

type SocialMedia struct {
	SocialMediaId string `bson:"social_media_id"`
	ProviderId    string `bson:"provider_id"`
}

type PaymentMethod struct {
	PaymentMethodId string `bson:"payment_method_id"`
	ProviderId      string `bson:"provider_id"`
}
