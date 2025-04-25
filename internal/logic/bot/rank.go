package bot

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"sort"
)

func (s *sBot) Rank(ctx context.Context, address string) (rank []model.RankBot, selfRank int, selfAmount int64, err error) {
	list := make([]*entity.Bot, 0)
	err = dao.Bot.Ctx(ctx).OrderDesc(dao.Bot.Columns().Amount).Limit(0, consts.RankLimit).Scan(&list)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Amount > list[j].Amount
	})
	for i, v := range list {
		if i == consts.RankLimit {
			break
		}
		if v.Address == address {
			selfRank = i + 1
			selfAmount = v.Amount
		}
		rank = append(rank, model.RankBot{
			Rank:    i + 1,
			Address: v.Address,
			Amount:  v.Amount,
		})
	}
	if selfAmount == 0 {
		bot, err := s.FindByAddress(ctx, address)
		if err != nil {
			return rank, selfRank, selfAmount, err
		}
		selfAmount = bot.Amount
	}
	if selfRank == 0 {
		selfRank = 999
	}
	return
}
