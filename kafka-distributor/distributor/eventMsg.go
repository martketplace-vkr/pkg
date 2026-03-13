package distributor

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/tidwall/gjson"
)

type EventMsg struct {
	*sarama.ConsumerMessage
}

func NewMsg(data *sarama.ConsumerMessage) *EventMsg {
	return &EventMsg{
		ConsumerMessage: data,
	}
}

func Unmarshal[T any](event *EventMsg, model *T) (err error) {
	err = json.Unmarshal(event.Value, model)
	if err != nil {
		return err
	}

	return err
}

func (e EventMsg) GetFieldByName(path string) (result string) {
	return gjson.GetBytes(e.Value, path).String()
}

func (e EventMsg) GetEventType() (result string) {
	return string(e.Key)
}
