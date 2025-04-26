package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain/okx"
	"time"
)

type sBot struct {
}

func init() {
	service.RegisterBot(New())
}

func New() service.IBot {
	return &sBot{}
}

func (s *sBot) Entity2BotInfo(ctx context.Context, bot *entity.Bot) *model.BotInfo {
	if bot == nil {
		return nil
	}
	power := decimal.NewFromInt(int64(bot.Power))
	h := decimal.NewFromInt(60)
	_, selfRank, _, err := s.Rank(ctx, bot.Address)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return &model.BotInfo{
		Address:        bot.Address,
		Amount:         bot.Amount,
		Decimals:       consts.EVEDecimal,
		TypeRate:       bot.TypeRate,
		TypeLimit:      consts.TypeLimit,
		TypeCount:      bot.TypeCount,
		Active:         bot.Active,
		Power:          power.Div(h).Round(1).InexactFloat64(),
		PowerLimit:     bot.PowerLimit / 60,
		AutoAccept:     bot.AutoAccept,
		Buff:           bot.Buff,
		BuffLimit:      consts.BuffLimit,
		SelfRank:       selfRank,
		MiningRateOfH:  consts.MiningRateOfH,
		ConsumptionOfH: consts.PowerOfOne / 60,
	}
}

func (s *sBot) FindByAddress(ctx context.Context, address string) (bot *entity.Bot, err error) {
	if err := dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, address).Scan(&bot); err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	if bot == nil {
		return nil, errors.New("bot not found")
	}
	return bot, nil
}

func (s *sBot) CreateBot(ctx context.Context, address, deviceId string) (bot *entity.Bot, err error) {
	if bot, err := s.FindByAddress(ctx, address); err == nil {
		return bot, nil
	}
	bot = &entity.Bot{
		Address:    address,
		DeviceId:   deviceId,
		Amount:     0,
		TypeRate:   0,
		TypeCount:  0,
		Active:     1,
		Power:      consts.PowerLimit,
		PowerLimit: consts.PowerLimit,
		AutoAccept: 0,
		Buff:       100,
	}
	result, err := dao.Bot.Ctx(ctx).Save(bot)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	bot.Id = id
	byAddress, err := okx.Okx().AllTokenBalancesByAddress(ctx, "501", address, "")
	if err != nil {
		g.Log().Error(ctx, err)
	}
	tokenAssets := make([]okx.TokenAsset, 0)
	if byAddress != nil && len(byAddress.Data) > 0 {
		tokenAssets = byAddress.Data[0].TokenAssets
	}
	if len(tokenAssets) > 0 {
		_, err = service.Wallet().SyncHold(ctx, address, tokenAssets)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}
	err = cache.GetCache().Set(ctx, fmt.Sprintf("%s:%s", consts.RedisActivePrefix, address), 1, time.Hour)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return bot, nil
}

func (s *sBot) addBuff(ctx context.Context, bot entity.Bot) (err error) {
	if bot.Buff >= consts.BuffLimit {
		return nil
	}
	bot.Buff += consts.BuffStep
	if bot.Buff > consts.BuffLimit {
		bot.Buff = consts.BuffLimit
	}
	_, err = dao.Bot.Ctx(ctx).Data(bot).Where(dao.Bot.Columns().Id, bot.Id).Data(g.Map{
		dao.Bot.Columns().Buff: bot.Buff,
	}).Update()
	return err
}

func (s *sBot) subBuff(ctx context.Context, bot entity.Bot) (err error) {
	if bot.Buff <= 100 {
		return nil
	}
	bot.Buff -= consts.BuffStep
	if bot.Buff < 100 {
		bot.Buff = 100
	}
	_, err = dao.Bot.Ctx(ctx).Data(bot).Where(dao.Bot.Columns().Id, bot.Id).Data(g.Map{
		dao.Bot.Columns().Buff: bot.Buff,
	}).Update()
	return err
}
