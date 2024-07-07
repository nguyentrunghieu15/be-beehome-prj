package cerbosx

type PrincipalInfo struct {
	ID         string `bson:"id"`
	UserId     string `bson:"user_id"`
	ProviderId string `bson:"provider_id"`
	Role       string `bson:"role"`
}
