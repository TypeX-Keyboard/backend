package bot

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/shopspring/decimal"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/util/redis"
	"keyboard-api-go/internal/util/str"
	"math"
	"strings"
)

func (s *sBot) ActiveTask(ctx context.Context) error {
	keys := redis.New().MatchKey(ctx, fmt.Sprintf("%s:*", consts.RedisActivePrefix))
	activeAddress := make([]string, 0)
	for _, key := range keys {
		activeAddress = append(activeAddress, strings.ReplaceAll(key, fmt.Sprintf("%s:", consts.RedisActivePrefix), ""))
	}
	slice := str.SplitSlice(activeAddress, 1000)
	var rowNum int64 = 0
	gmlock.Lock(consts.PowerLock)
	defer gmlock.Unlock(consts.PowerLock)
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, activeAddress := range slice {
			_, err := dao.Bot.Ctx(ctx).WhereNotIn(dao.Bot.Columns().Address, activeAddress).Data(g.Map{
				dao.Bot.Columns().Active: 0,
			}).Update()
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			update, err := dao.Friend.Ctx(ctx).WhereNotIn(dao.Friend.Columns().FriendAddress, activeAddress).Data(g.Map{
				dao.Friend.Columns().Active: 0,
			}).Update()
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			affected, err := update.RowsAffected()
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			if affected > 0 {
				rowNum += affected
			}
		}
		if rowNum > 0 {
			counts, err := s.ActiveFriendCounts(ctx)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			for _, count := range counts {
				if count.FriendCount == count.ActiveCount {
					continue
				}
				_, err := dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, count.Address).Data(g.Map{
					dao.Bot.Columns().Buff: 100 + (consts.BuffStep * count.ActiveCount),
				}).Update()
				if err != nil {
					g.Log().Error(ctx, err)
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (s *sBot) ActiveFriendCounts(ctx context.Context) ([]model.ActiveFriends, error) {
	activeFriends := make([]model.ActiveFriends, 0)
	all, err := dao.Friend.Ctx(ctx).Fields("address,count(*) friendCount,sum(active) activeCount").Group("address").All()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	err = all.Structs(&activeFriends)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return activeFriends, nil
}

func (s *sBot) MiningTask(ctx context.Context) error {
	gmlock.Lock(consts.PowerLock)
	defer gmlock.Unlock(consts.PowerLock)
	eveOf1m := s.EveOf1M()
	_, err := dao.Bot.Ctx(ctx).WhereGT(dao.Bot.Columns().Power, 0).Data(g.Map{
		dao.Bot.Columns().Power:  gdb.Raw("power - 1"),
		dao.Bot.Columns().Amount: gdb.Raw(fmt.Sprintf("amount+(%d*buff/100)", eveOf1m)),
	}).Update()
	return err
}

func (s *sBot) EveOf1M() int64 {
	eve := decimal.NewFromFloat(math.Pow10(consts.EVEDecimal) * consts.MiningRateOfH)
	eveOf1M := eve.Div(decimal.NewFromInt(60)).Round(0)
	return eveOf1M.IntPart()
}
