package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialMediaResourceHandler struct {
	logger                logwrapper.ILoggerWrapper
	socialMediaRepository mongox.Repository[model.SocialMedia]
}

func NewSocialMediaResourceHandler(logger logwrapper.ILoggerWrapper) *SocialMediaResourceHandler {
	return &SocialMediaResourceHandler{
		logger: logger,
		socialMediaRepository: mongox.Repository[model.SocialMedia]{
			Client:     mongox.DefaultClient,
			Collection: "social_media",
		},
	}
}

type SocialMediaResourceMsg struct {
	Type          string `json:"type"`
	SocialMediaId string `json:"social_media_id"`
	ProviderId    string `json:"provider_id"`
}

func (h *SocialMediaResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &SocialMediaResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateSocialMediaResource(parsedMsg)

	case "update":
		return h.UpdateSocialMediaResource(parsedMsg)

	case "delete":
		return h.DeleteSocialMediaResource(parsedMsg)
	}
	return nil
}

func (h *SocialMediaResourceHandler) CreateSocialMediaResource(msg *SocialMediaResourceMsg) error {
	socialMedia := model.SocialMedia{
		ID:            primitive.NewObjectID(),
		SocialMediaId: msg.SocialMediaId,
		ProviderId:    msg.ProviderId,
	}

	exsited, err := h.socialMediaRepository.FindOneByAtribute("social_media_id", msg.SocialMediaId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.socialMediaRepository.InsertOne(socialMedia)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *SocialMediaResourceHandler) UpdateSocialMediaResource(msg *SocialMediaResourceMsg) error {

	return nil
}

func (h *SocialMediaResourceHandler) DeleteSocialMediaResource(msg *SocialMediaResourceMsg) error {
	err := h.socialMediaRepository.DeleteOneByAtribute("social_media_id", msg.SocialMediaId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
