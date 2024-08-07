package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
)

type UserResourceHandler struct {
	logger             logwrapper.ILoggerWrapper
	userRepository     mongox.Repository[model.User]
	providerRepository mongox.Repository[model.Provider]
}

func NewUserResourceHandler(logger logwrapper.ILoggerWrapper) *UserResourceHandler {
	return &UserResourceHandler{
		logger: logger,
		userRepository: mongox.Repository[model.User]{
			Client:     mongox.DefaultClient,
			Collection: "users",
		},
		providerRepository: mongox.Repository[model.Provider]{
			Client:     mongox.DefaultClient,
			Collection: "providers",
		},
	}
}

type UserResourceMsg struct {
	Type       string `json:"type"`
	UserId     string `json:"user_id"`
	ProviderId string `json:"provider_id"`
	Role       string `json:"role"`
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
		h.UpdateUserResource(parsedMsg)

	case "delete":
		h.DeleteUserResource(parsedMsg)
	}
	return nil
}

func (h *UserResourceHandler) CreateUserResource(msg *UserResourceMsg) error {
	user := model.User{
		ProviderId: msg.ProviderId,
		UserId:     msg.UserId,
		Role:       msg.Role,
	}

	exsited, err := h.userRepository.FindOneByAtribute("user_id", msg.UserId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.userRepository.InsertOne(user)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *UserResourceHandler) UpdateUserResource(msg *UserResourceMsg) error {
	return nil
}

func (h *UserResourceHandler) DeleteUserResource(msg *UserResourceMsg) error {
	err := h.userRepository.DeleteOneByAtribute("user_id", msg.UserId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	h.providerRepository.DeleteOneByAtribute("user_id", msg.UserId)

	return nil
}
