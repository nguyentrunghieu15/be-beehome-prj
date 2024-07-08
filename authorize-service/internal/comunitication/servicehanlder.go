package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceResourceHandler struct {
	logger            logwrapper.ILoggerWrapper
	serviceRepository mongox.Repository[model.Service]
}

func NewServiceResourceHandler(logger logwrapper.ILoggerWrapper) *ServiceResourceHandler {
	return &ServiceResourceHandler{
		logger: logger,
		serviceRepository: mongox.Repository[model.Service]{
			Client:     mongox.DefaultClient,
			Collection: "service",
		},
	}
}

type ServiceResourceMsg struct {
	Type           string `json:"type"`
	ServiceId      string `json:"service_id"`
	GroupServiceId string `json:"group_service_id"`
}

func (h *ServiceResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &ServiceResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateServiceResource(parsedMsg)

	case "update":
		return h.UpdateServiceResource(parsedMsg)

	case "delete":
		return h.DeleteServiceResource(parsedMsg)
	}
	return nil
}

func (h *ServiceResourceHandler) CreateServiceResource(msg *ServiceResourceMsg) error {
	service := model.Service{
		ID:             primitive.NewObjectID(),
		GroupServiceId: msg.GroupServiceId,
		ServiceId:      msg.ServiceId,
	}

	exsited, err := h.serviceRepository.FindOneByAtribute("service_id", msg.ServiceId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.serviceRepository.InsertOne(service)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *ServiceResourceHandler) UpdateServiceResource(msg *ServiceResourceMsg) error {
	return nil
}

func (h *ServiceResourceHandler) DeleteServiceResource(msg *ServiceResourceMsg) error {
	err := h.serviceRepository.DeleteOneByAtribute("service_id", msg.ServiceId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
