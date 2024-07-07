package communication

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
	"github.com/segmentio/kafka-go"
)

type UserResourceHandler struct {
	logger        logwrapper.ILoggerWrapper
	proRepository datasource.IProviderRepo
}

func NewUserResourceHandler(logger logwrapper.ILoggerWrapper,
	proRepository datasource.IProviderRepo) *UserResourceHandler {
	return &UserResourceHandler{
		logger:        logger,
		proRepository: proRepository,
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
	pro, err := h.proRepository.FindOneByUserId(uuid.MustParse(msg.UserId))
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	err = h.proRepository.DeleteOneById(pro.ID)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
