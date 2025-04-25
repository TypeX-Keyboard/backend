package bot

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/util/cache"
	"time"
)

func (s *sBot) Active(ctx context.Context, address string) (err error) {
	_, err = dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, address).Data(g.Map{
		dao.Bot.Columns().Active: 1,
	}).Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	_, err = dao.Friend.Ctx(ctx).Where(dao.Friend.Columns().FriendAddress, address).Data(g.Map{
		dao.Bot.Columns().Active: 1,
	}).Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	err = cache.GetCache().Set(ctx, fmt.Sprintf("%s:%s", consts.RedisActivePrefix, address), 1, time.Hour)
	return err
}
