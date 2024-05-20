package mapper

import (
	"time"

	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
)

func ConvertUserToUserInfor(user *datasource.User) (*userapi.UserInfor, error) {
	// Handle potential conversion errors (e.g., time format)
	createdAt := user.CreatedAt.Format(time.RFC3339Nano)
	updatedAt := ""
	deletedAt := ""

	if !user.UpdatedAt.IsZero() {
		deletedAt = user.UpdatedAt.Format(time.RFC3339Nano)
	}

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

func ConvertListUserToListUserInfor(users []*datasource.User) ([]*userapi.UserInfor, error) {
	userInfos := make([]*userapi.UserInfor, 0, len(users))
	for _, user := range users {
		userInfo, err := ConvertUserToUserInfor(user)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}
	return userInfos, nil
}
