package communication

import (
	"os"
	"time"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/kafkax"
)

const (
	TOPIC_RESOURCE_USER          = "user-resource"
	TOPIC_RESOURCE_PROVIDER      = "provider-resource"
	TOPIC_RESOURCE_SERVICE       = "service-resource"
	TOPIC_RESOURCE_GSERVICE      = "group-service-resource"
	TOPIC_RESOURCE_HIRE          = "hire-resource"
	TOPIC_RESOURCE_SOCIALMEDIA   = "social-media-resource"
	TOPIC_RESOURCE_PAYMENTMETHOD = "payment-method-resource"
	TOPIC_RESOURCE_REVIEW        = "review-resource"
)

var UserResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_USER,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var ProviderResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_PROVIDER,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var ServiceResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_SERVICE,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var GroupServiceResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_GSERVICE,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var HireResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_HIRE,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var ReviewResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_REVIEW,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var PaymentMethodResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_PAYMENTMETHOD,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)

var SocialMediaResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
	&kafkax.KafkaClientConfig{
		Topic:            TOPIC_RESOURCE_SOCIALMEDIA,
		BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Protocall:        "tcp",
		MaxBytes:         10e6,
		TimeoutRead:      time.Second,
		TimeoutWrite:     time.Second,
	},
)
