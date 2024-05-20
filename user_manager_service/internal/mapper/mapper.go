package mapper

import (
	"time"

	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
)

func ConvertUserToUserInfor(user datasource.User) (*userapi.UserInfor, error) {
	// Handle potential conversion errors (e.g., time format)
	createdAt := user.CreatedAt.Format(time.RFC3339Nano)
	updatedAt := user.UpdatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	if !user.DeletedAt.Time.IsZero() {
		deletedAt = user.DeletedAt.Time.Format(time.RFC3339Nano)
	}

	// Combine first and last name for UserInfor.Name
	name := user.FirstName + " " + user.LastName

	return &userapi.UserInfor{
		Id:        user.ID.String(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      name,
		Status:    user.Status,
	}, nil
}
