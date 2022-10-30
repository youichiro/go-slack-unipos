package usecase

import (
	"github.com/youichiro/go-slack-my-unipos/internal/gateway"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func OpenSlackModalUsecase(token string, triggerID string) error {
	err := gateway.SlackOpenModal(token, triggerID)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}
	return nil
}
