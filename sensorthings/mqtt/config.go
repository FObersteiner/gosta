package mqtt

import (
	"github.com/FObersteiner/gosta-server/sensorthings/models"
)

// CreateTopics creates the pre-defined MQTT Topics
func CreateTopics(prefix string) []models.Topic {
	topics := []models.Topic{
		{
			Path:    prefix + "/#",
			Handler: MainMqttHandler,
		},
	}

	return topics
}
