package usecase

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/youichiro/go-slack-my-unipos/internal/models"
)

func findOrCreareMember(ctx *gin.Context, db *sql.DB, slackUserId string) (*models.Member, error) {
	// メンバーを取得する. もし存在しない場合は作成する.
	var member *models.Member
	var err error

	member, err = models.Members(qm.Where("slack_user_id = ?", slackUserId)).One(ctx, db)
	if err != nil && err.Error() == "sql: no rows in result set" {
		member = &models.Member{
			SlackUserID: slackUserId,
		}
		err = member.Insert(ctx, db, boil.Infer())
		if err != nil {
			return nil, err
		}
		log.Debug().Msg("created a new member: " + member.SlackUserID)
	}
	return member, nil
}

func CreateCardUsecase(ctx *gin.Context, db *sql.DB, senderSlackUserId string, distinationSlackUserId string, point int, message string) error {
	boil.DebugMode = true
	senderMember, err := findOrCreareMember(ctx, db, senderSlackUserId)
	if err != nil {
		return err
	}
	distinationMember, err := findOrCreareMember(ctx, db, distinationSlackUserId)
	if err != nil {
		return err
	}

	// メンバーのカードを取得する
	cards, err := models.Cards(qm.Where("sender_member_id = ?", senderMember.ID)).All(ctx, db)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	fmt.Println(cards)

	// メンバーの残pointを取得する
	remainPoint := 400
	if len(cards) > 0 {
		for _, card := range cards {
			remainPoint = remainPoint - card.Point
		}
	}
	log.Debug().Msg(fmt.Sprintf("member_id: %d, remainPoint: %d", senderMember.ID, remainPoint))

	// もしポイントが足りなかったらエラーにする
	if remainPoint-point < 0 {
		return fmt.Errorf("not enough points")
	}

	// カードを作成する
	newCard := models.Card{
		SenderMemberID:      senderMember.ID,
		DistinationMemberID: distinationMember.ID,
		Point:               point,
		Message:             message,
	}
	err = newCard.Insert(ctx, db, boil.Infer())
	if err != nil {
		return err
	}
	log.Debug().Msg("A new card is cerated!")

	return nil
}
