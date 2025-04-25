package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model/entity"
)

func (s *sBot) SubmitWork(ctx context.Context, address string, typeCount int) (err error) {
	gmlock.Lock(fmt.Sprintf(consts.SubmitWorkLock, address))
	defer gmlock.Unlock(fmt.Sprintf(consts.SubmitWorkLock, address))
	var bot *entity.Bot
	err = dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, address).Scan(&bot)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if bot == nil {
		return errors.New("bot not found")
	}
	if _, err := dao.WorkRecord.Ctx(ctx).Data(entity.WorkRecord{Address: address, TypeCount: typeCount}).Insert(); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	total := typeCount + bot.TypeRate
	power := bot.Power + (total/consts.TypeLimit)*consts.PowerOfOne // For every 100 characters, 1 point of electricity is added, and 1 point of electricity is divided into 60 parts
	if power > bot.PowerLimit {
		power = bot.PowerLimit
	}
	_, err = dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, address).Data(g.Map{
		dao.Bot.Columns().TypeRate:  total % consts.TypeLimit,
		dao.Bot.Columns().TypeCount: bot.TypeCount + int64(typeCount),
		dao.Bot.Columns().Power:     power,
	}).Update()
	return err
}
