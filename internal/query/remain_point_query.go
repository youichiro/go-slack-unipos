package query

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/youichiro/go-slack-my-unipos/internal/models"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func RemainPointQuery(ctx *gin.Context, db *sql.DB, memberID int) (int, error) {
	boil.DebugMode = true

	maxWeekPoint, _ := strconv.Atoi(os.Getenv("MAX_WEEK_POINT"))
	sumPoint, err := models.Cards(
		qm.Select("point"),
		qm.Where("sender_member_id = ?", memberID),
	).Count(ctx, db)
	if err != nil {
		util.Logger.Error(err.Error())
		return 0, err
	}
	remainPoint := maxWeekPoint - int(sumPoint)
	return remainPoint, nil
}
