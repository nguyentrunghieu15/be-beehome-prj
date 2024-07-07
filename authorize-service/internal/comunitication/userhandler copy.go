package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
)

type UserResourceHandler struct {
	logger         logwrapper.ILoggerWrapper
	userRepository mongox.Repository[model.User]
}

func NewUserResourceHandler(logger logwrapper.ILoggerWrapper) *UserResourceHandler {
	return &UserResourceHandler{
		logger: logger,
		userRepository: mongox.Repository[model.User]{
			Client:     mongox.DefaultClient,
			Collection: "user",
		},
	}
}

type UserResourceMsg struct {
	Type       string `json:"type"`
	UserId     string `bson:"user_id"`
	ProviderId string `bson:"provider_id"`
	Role       string `bson:"role"`
}

func (h *UserResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &UserResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		h.CreateUserResource(parsedMsg)

	case "update":
		h.CreateUserResource(parsedMsg)

	case "delete":
		h.CreateUserResource(parsedMsg)
	}
	return nil
}

func (h *UserResourceHandler) CreateUserResource(msg *UserResourceMsg) error {
	return nil
}

func (h *UserResourceHandler) UpdateUserResource(msg *UserResourceMsg) error {
	return nil
}

func (h *UserResourceHandler) DeleteUserResource(msg *UserResourceMsg) error {
	return nil
}
