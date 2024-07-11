package communication

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/datasource"
	"github.com/segmentio/kafka-go"
)

func UserResourceHandler(kafka.Message) {

}

type ProviderResourceMsg struct {
	Type       string `json:"type"`
	UserId     string `json:"user_id"`
	ProviderId string `json:"provider_id"`
}

func NewProviderResourceHandler(userRepo datasource.IUserRepo,
	bannedAccountRepo datasource.IBannedAccountsRepo,
	cardRepo datasource.ICardRepo,
	logger logwrapper.ILoggerWrapper) *ProviderResourceHandler {
	return &ProviderResourceHandler{
		userRepo:          userRepo,
		bannedAccountRepo: bannedAccountRepo,
		cardRepo:          cardRepo,
		logger:            logger,
	}
}

type ProviderResourceHandler struct {
	userRepo          datasource.IUserRepo
	bannedAccountRepo datasource.IBannedAccountsRepo
	cardRepo          datasource.ICardRepo
	logger            logwrapper.ILoggerWrapper
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
		h.CreateProviderResource(parsedMsg)

	case "update":
		h.UpdateProviderResource(parsedMsg)

	case "delete":
		h.DeleteProviderResource(parsedMsg)
	}
	return nil
}

func (h *ProviderResourceHandler) CreateProviderResource(msg *ProviderResourceMsg) error {
	h.userRepo.UpdateOneById(uuid.MustParse(msg.UserId), map[string]interface{}{
		"provider_id": msg.ProviderId,
	})
	return nil
}

func (h *ProviderResourceHandler) UpdateProviderResource(msg *ProviderResourceMsg) error {
	return nil
}

func (h *ProviderResourceHandler) DeleteProviderResource(msg *ProviderResourceMsg) error {
	h.userRepo.UpdateOneById(uuid.MustParse(msg.UserId), map[string]interface{}{
		"provider_id": "",
	})
	return nil
}
