package communication

import (
	"encoding/json"

	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewResourceHandler struct {
	logger           logwrapper.ILoggerWrapper
	reviewRepository mongox.Repository[model.Review]
}

func NewReviewResourceHandler(logger logwrapper.ILoggerWrapper) *ReviewResourceHandler {
	return &ReviewResourceHandler{
		logger: logger,
		reviewRepository: mongox.Repository[model.Review]{
			Client:     mongox.DefaultClient,
			Collection: "review",
		},
	}
}

type ReviewResourceMsg struct {
	Type       string `json:"type"`
	ReviewId   string `json:"review_id"`
	HireId     string `json:"hire_id"`
	ProviderId string `json:"provider_id"`
	UserId     string `json:"user_id"`
}

func (h *ReviewResourceHandler) Router(msg kafka.Message) error {
	parsedMsg := &ReviewResourceMsg{}
	err := json.Unmarshal(msg.Value, parsedMsg)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	switch parsedMsg.Type {
	case "create":
		return h.CreateReviewResource(parsedMsg)

	case "update":
		return h.UpdateReviewResource(parsedMsg)

	case "delete":
		return h.DeleteReviewResource(parsedMsg)
	}
	return nil
}

func (h *ReviewResourceHandler) CreateReviewResource(msg *ReviewResourceMsg) error {
	review := model.Review{
		ID:         primitive.NewObjectID(),
		ReviewId:   msg.ReviewId,
		HireId:     msg.HireId,
		ProviderId: msg.ProviderId,
		UserId:     msg.UserId,
	}

	exsited, err := h.reviewRepository.FindOneByAtribute("review_id", msg.ReviewId)
	if err == nil && exsited != nil {
		return nil
	}

	err = h.reviewRepository.InsertOne(review)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	return nil
}

func (h *ReviewResourceHandler) UpdateReviewResource(msg *ReviewResourceMsg) error {

	return nil
}

func (h *ReviewResourceHandler) DeleteReviewResource(msg *ReviewResourceMsg) error {
	err := h.reviewRepository.DeleteOneByAtribute("review_id", msg.ReviewId)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return nil
}
