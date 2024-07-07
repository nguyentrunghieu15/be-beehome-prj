package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProviderResourceHandler struct {
	logger             logwrapper.ILoggerWrapper
	providerRepository mongox.Repository[model.Provider]
}

func NewProviderResourceHandler(logger logwrapper.ILoggerWrapper) *ProviderResourceHandler {
	return &ProviderResourceHandler{
		logger: logger,
		providerRepository: mongox.Repository[model.Provider]{
			Client:     mongox.DefaultClient,
			Collection: "provider",
		},
	}
}

type ProviderResourceMsg struct {
	Type       string `json:"type"`
	ProviderId string `bson:"provider_id"`
	UserId     string `bson:"user_id"`
}

func (h *ProviderResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &ProviderResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateProviderResource(parsedMsg)

	case "update":
		return h.UpdateProviderResource(parsedMsg)

	case "delete":
		return h.DeleteProviderResource(parsedMsg)
	}
	return nil
}

func (h *ProviderResourceHandler) CreateProviderResource(msg *ProviderResourceMsg) error {
	provider := model.Provider{
		ID:         primitive.NewObjectID(),
		ProviderId: msg.ProviderId,
		UserId:     msg.UserId,
	}

	err := h.providerRepository.InsertOne(provider)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *ProviderResourceHandler) UpdateProviderResource(msg *ProviderResourceMsg) error {
	return nil
}

func (h *ProviderResourceHandler) DeleteProviderResource(msg *ProviderResourceMsg) error {
	err := h.providerRepository.DeleteOneByAtribute("provider_id", msg.ProviderId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
