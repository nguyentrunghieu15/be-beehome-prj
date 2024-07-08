package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupServiceResourceHandler struct {
	logger                 logwrapper.ILoggerWrapper
	groupServiceRepository mongox.Repository[model.GroupService]
}

func NewGroupServiceResourceHandler(logger logwrapper.ILoggerWrapper) *GroupServiceResourceHandler {
	return &GroupServiceResourceHandler{
		logger: logger,
		groupServiceRepository: mongox.Repository[model.GroupService]{
			Client:     mongox.DefaultClient,
			Collection: "group_service",
		},
	}
}

type GroupServiceResourceMsg struct {
	Type           string `json:"type"`
	GroupServiceId string `bson:"group_service_id"`
}

func (h *GroupServiceResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &GroupServiceResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateGroupServiceResource(parsedMsg)

	case "update":
		return h.UpdateGroupServiceResource(parsedMsg)

	case "delete":
		return h.DeleteGroupServiceResource(parsedMsg)
	}
	return nil
}

func (h *GroupServiceResourceHandler) CreateGroupServiceResource(msg *GroupServiceResourceMsg) error {
	groupService := model.GroupService{
		ID:             primitive.NewObjectID(),
		GroupServiceId: msg.GroupServiceId,
	}
	exsited, err := h.groupServiceRepository.FindOneByAtribute("group_service_id", msg.GroupServiceId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.groupServiceRepository.InsertOne(groupService)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *GroupServiceResourceHandler) UpdateGroupServiceResource(msg *GroupServiceResourceMsg) error {
	return nil
}

func (h *GroupServiceResourceHandler) DeleteGroupServiceResource(msg *GroupServiceResourceMsg) error {
	err := h.groupServiceRepository.DeleteOneByAtribute("group_service_id", msg.GroupServiceId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
