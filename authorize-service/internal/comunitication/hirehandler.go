package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
)

type HireResourceHandler struct {
	logger         logwrapper.ILoggerWrapper
	hireRepository mongox.Repository[model.Hire]
}

func NewHireResourceHandler(logger logwrapper.ILoggerWrapper) *HireResourceHandler {
	return &HireResourceHandler{
		logger: logger,
		hireRepository: mongox.Repository[model.Hire]{
			Client:     mongox.DefaultClient,
			Collection: "hires",
		},
	}
}

type HireResourceMsg struct {
	Type       string `json:"type"`
	HireId     string `json:"hire_id"`
	ProviderId string `json:"provider_id"`
	UserId     string `json:"user_id"`
}

func (h *HireResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &HireResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateHireResource(parsedMsg)

	case "update":
		return h.UpdateHireResource(parsedMsg)

	case "delete":
		return h.DeleteHireResource(parsedMsg)
	}
	return nil
}

func (h *HireResourceHandler) CreateHireResource(msg *HireResourceMsg) error {
	hire := model.Hire{
		HireId:     msg.HireId,
		ProviderId: msg.ProviderId,
		UserId:     msg.UserId,
	}
	hireExsited, err := h.hireRepository.FindOneByAtribute("hire_id", msg.HireId)
	if err == nil && hireExsited != nil {
		return nil
	}

	err = h.hireRepository.InsertOne(hire)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *HireResourceHandler) UpdateHireResource(msg *HireResourceMsg) error {
	return nil
}

func (h *HireResourceHandler) DeleteHireResource(msg *HireResourceMsg) error {
	err := h.hireRepository.DeleteOneByAtribute("hire_id", msg.HireId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
