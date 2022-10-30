package gateway

import (
	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func SlackOpenModal(token string, triggerID string, modalViewRequest slack.ModalViewRequest) error {
	api := slack.New(token)
	_, err := api.OpenView(triggerID, modalViewRequest)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}
	return nil
}
