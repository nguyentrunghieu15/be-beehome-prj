package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
)

type ProviderResourceHandler struct {
	logger                  logwrapper.ILoggerWrapper
	providerRepository      mongox.Repository[model.Provider]
	userRepository          mongox.Repository[model.User]
	hireRepository          mongox.Repository[model.Hire]
	socialMediaRepository   mongox.Repository[model.SocialMedia]
	paymentMethodRepository mongox.Repository[model.PaymentMethod]
}

func NewProviderResourceHandler(logger logwrapper.ILoggerWrapper) *ProviderResourceHandler {
	return &ProviderResourceHandler{
		logger: logger,
		providerRepository: mongox.Repository[model.Provider]{
			Client:     mongox.DefaultClient,
			Collection: "providers",
		},
		userRepository: mongox.Repository[model.User]{
			Client:     mongox.DefaultClient,
			Collection: "users",
		},
		hireRepository: mongox.Repository[model.Hire]{
			Client:     mongox.DefaultClient,
			Collection: "hires",
		},
		socialMediaRepository: mongox.Repository[model.SocialMedia]{
			Client:     mongox.DefaultClient,
			Collection: "social_medias",
		},
		paymentMethodRepository: mongox.Repository[model.PaymentMethod]{
			Client:     mongox.DefaultClient,
			Collection: "payment_methods",
		},
	}
}

type ProviderResourceMsg struct {
	Type       string `json:"type"`
	ProviderId string `json:"provider_id"`
	UserId     string `json:"user_id"`
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
		ProviderId: msg.ProviderId,
		UserId:     msg.UserId,
	}

	exsited, err := h.providerRepository.FindOneByAtribute("provider_id", msg.ProviderId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.providerRepository.InsertOne(provider)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	h.userRepository.UpdateOneByFilterForOneAtribute(map[string]interface{}{
		"user_id": msg.UserId,
	}, "provider_id", msg.ProviderId)

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
	err = h.socialMediaRepository.DeleteOneByAtribute("provider_id", msg.ProviderId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	err = h.paymentMethodRepository.DeleteOneByAtribute("provider_id", msg.ProviderId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	err = h.hireRepository.DeleteOneByAtribute("provider_id", msg.ProviderId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}
