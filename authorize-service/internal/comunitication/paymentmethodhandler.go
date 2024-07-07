package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethodResourceHandler struct {
	logger                  logwrapper.ILoggerWrapper
	paymentMethodRepository mongox.Repository[model.PaymentMethod]
}

func NewPaymentMethodResourceHandler(logger logwrapper.ILoggerWrapper) *PaymentMethodResourceHandler {
	return &PaymentMethodResourceHandler{
		logger: logger,
		paymentMethodRepository: mongox.Repository[model.PaymentMethod]{
			Client:     mongox.DefaultClient,
			Collection: "payment_method",
		},
	}
}

type PaymentMethodResourceMsg struct {
	Type            string `json:"type"`
	PaymentMethodId string `bson:"payment_method_id"`
	ProviderId      string `bson:"provider_id"`
}

func (h *PaymentMethodResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &PaymentMethodResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreatePaymentMethodResource(parsedMsg)

	case "update":
		return h.UpdatePaymentMethodResource(parsedMsg)

	case "delete":
		return h.DeletePaymentMethodResource(parsedMsg)
	}
	return nil
}

func (h *PaymentMethodResourceHandler) CreatePaymentMethodResource(msg *PaymentMethodResourceMsg) error {
	paymentMethod := model.PaymentMethod{
		ID:              primitive.NewObjectID(),
		PaymentMethodId: msg.PaymentMethodId,
		ProviderId:      msg.ProviderId,
	}

	err := h.paymentMethodRepository.InsertOne(paymentMethod)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *PaymentMethodResourceHandler) UpdatePaymentMethodResource(msg *PaymentMethodResourceMsg) error {

	return nil
}

func (h *PaymentMethodResourceHandler) DeletePaymentMethodResource(msg *PaymentMethodResourceMsg) error {
	err := h.paymentMethodRepository.DeleteOneByAtribute("payment_method_id", msg.PaymentMethodId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
