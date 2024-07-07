package communication

import (
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
)

var ProviderResourceKafka *kafkax.KafkaClientWrapper
var UserResourceKafka *kafkax.KafkaClientWrapper
var ServiceResourceKafka *kafkax.KafkaClientWrapper
var GroupServiceResourceKafka *kafkax.KafkaClientWrapper
var HireResourceKafka *kafkax.KafkaClientWrapper
var SocialMediaResourceKafka *kafkax.KafkaClientWrapper
var PaymentMethodResourceKafka *kafkax.KafkaClientWrapper
