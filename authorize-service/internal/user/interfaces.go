package user

type UserInfor struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedAt string `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
	DeletedBy string `json:"deletedBy"`
	DeletedAt string `json:"deletedAt"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AddCardRequest struct {
	Card Card `json:"card"`
}

type Card struct {
	CardNumber string `json:"cardNumber"`
	OwnerName  string `json:"ownerName"`
	BankName   string `json:"bankName"`
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

type ChangeNameRequest struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ListUsersResponse struct {
	Users []UserInfor `json:"users"`
}

type UserServiceBlockUserBody struct {
	Description string `json:"description"`
}

type UserServiceUpdateUserBody struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
}
