package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"sort"
	"strings"
)

func (s *sBot) UpdateAutoAccept(ctx context.Context, address string, autoAccept bool) (err error) {
	_, err = dao.Bot.Ctx(ctx).Where(dao.Bot.Columns().Address, address).Data(g.Map{
		dao.Bot.Columns().AutoAccept: autoAccept,
	}).Update()
	return err
}

func (s *sBot) FriendList(ctx context.Context, address string) (res []model.FriendInfo, err error) {
	list := make([]entity.Friend, 0)
	err = dao.Friend.Ctx(ctx).Where(dao.Friend.Columns().Address, address).OrderDesc(dao.Friend.Columns().Active).Scan(&list)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res = make([]model.FriendInfo, 0)
	for _, friend := range list {
		res = append(res, model.FriendInfo{
			Address: friend.FriendAddress,
			Active:  friend.Active == 1,
		})
	}
	return
}

func (s *sBot) AddFriend(ctx context.Context, address, friendAddress string) (bool, error) {
	params := []string{
		address, friendAddress,
	}
	sort.Strings(params)
	paramStr := strings.Join(params, "-")
	gmlock.Lock(fmt.Sprintf(consts.AddFriendLock, paramStr))
	defer gmlock.Unlock(fmt.Sprintf(consts.AddFriendLock, paramStr))
	bot, err := s.FindByAddress(ctx, address)
	if err != nil {
		g.Log().Error(ctx, err)
		return false, gerror.NewCode(gcode.New(consts.UserNotFoundErr, "This address has not been registered with TypeX Keyboard", nil))
	}
	friend, err := s.FindByAddress(ctx, friendAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return false, gerror.NewCode(gcode.New(consts.UserNotFoundErr, "This address has not been registered with TypeX Keyboard", nil))
	}
	if bot.DeviceId == friend.DeviceId {
		return false, gerror.NewCode(gcode.New(consts.CanNotAddFriendSelfErr, "Can't add yourself as a friend", nil))
	}
	if !s.verify(ctx, address) {
		return false, gerror.NewCode(gcode.New(consts.FriendLimitErr, "The number of friends has reached the upper limit", nil))
	}
	if !s.verify(ctx, friendAddress) {
		return false, gerror.NewCode(gcode.New(consts.FriendLimitErr, "The number of friends on the other party has reached the upper limit", nil))
	}
	status := consts.FriendPending
	if friend.AutoAccept == 1 {
		status = consts.FriendAccept
	}
	isAutoAccept := false
	count, err := dao.Friend.Ctx(ctx).Where(g.Map{
		dao.Friend.Columns().Address:       address,
		dao.Friend.Columns().FriendAddress: friendAddress,
	}).WhereOr(g.Map{
		dao.Friend.Columns().Address:       friendAddress,
		dao.Friend.Columns().FriendAddress: address,
	}).Count()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		g.Log().Error(ctx, err)
		return false, gerror.NewCode(gcode.CodeInternalError, "Friend addition failed. Check the network and retry.")
	}
	if count > 0 {
		return false, gerror.NewCode(gcode.New(consts.FriendAddedErr, "Friend request has been sent", nil))
	}
	count, err = dao.FriendMsg.Ctx(ctx).Where(g.Map{
		dao.FriendMsg.Columns().Address:       address,
		dao.FriendMsg.Columns().FriendAddress: friendAddress,
		dao.FriendMsg.Columns().Status:        consts.FriendPending,
	}).WhereOr(g.Map{
		dao.FriendMsg.Columns().Address:       friendAddress,
		dao.FriendMsg.Columns().FriendAddress: address,
		dao.FriendMsg.Columns().Status:        consts.FriendPending,
	}).Count()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		g.Log().Error(ctx, err)
		return false, gerror.NewCode(gcode.CodeInternalError, "Friend addition failed. Check the network and retry.")
	}
	if count > 0 {
		return false, gerror.NewCode(gcode.New(consts.FriendAddedErr, "Friend request has been sent", nil))
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.FriendMsg.Ctx(ctx).Data(entity.FriendMsg{
			Address:       address,
			FriendAddress: friendAddress,
			Status:        status,
		}).Insert()
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		if status == consts.FriendAccept {
			isAutoAccept = true
			err = s.addFriendCallback(ctx, *bot, *friend)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return isAutoAccept, gerror.NewCode(gcode.CodeInternalError, "Friend addition failed. Check the network and retry.")
	}
	return isAutoAccept, nil
}

func (s *sBot) FriendMsgList(ctx context.Context, address string, isSelf bool, page, size int) (res []model.FriendMsg, total int, err error) {
	list := make([]*entity.FriendMsg, 0)
	field := dao.FriendMsg.Columns().FriendAddress
	if isSelf {
		field = dao.FriendMsg.Columns().Address
	}
	db := dao.FriendMsg.Ctx(ctx).Where(field, address).Order(dao.FriendMsg.Columns().Status)
	if size > 0 {
		db = db.Page(page, size)
	}
	err = db.ScanAndCount(&list, &total, true)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res = make([]model.FriendMsg, 0)
	for _, msg := range list {
		temp := model.FriendMsg{
			Id:      msg.Id,
			Address: msg.Address,
			Status:  msg.Status,
		}
		if isSelf {
			temp.Address = msg.FriendAddress
		}
		res = append(res, temp)
	}
	return
}

func (s *sBot) HandleFriendMsg(ctx context.Context, id int64, address string, accept bool) (err error) {
	var msg *entity.FriendMsg
	err = dao.FriendMsg.Ctx(ctx).Where(dao.FriendMsg.Columns().Id, id).Scan(&msg)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if msg.FriendAddress != address {
		return errors.New("no right to operate")
	}
	if msg.Status != consts.FriendPending {
		return errors.New("the friend request has been processed")
	}
	params := []string{
		msg.Address, msg.FriendAddress,
	}
	sort.Strings(params)
	paramStr := strings.Join(params, "-")
	gmlock.Lock(fmt.Sprintf(consts.AddFriendLock, paramStr))
	defer gmlock.Unlock(fmt.Sprintf(consts.AddFriendLock, paramStr))
	status := consts.FriendReject
	if accept {
		status = consts.FriendAccept
	}
	bot, err := s.FindByAddress(ctx, msg.Address)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	friend, err := s.FindByAddress(ctx, msg.FriendAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	defer func(status *int, bot, friend *entity.Bot) {
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err = dao.FriendMsg.Ctx(ctx).Where(dao.FriendMsg.Columns().Id, id).Data(g.Map{
				dao.FriendMsg.Columns().Status: *status,
			}).Update()
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			if accept {
				err = s.addFriendCallback(ctx, *bot, *friend)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}(&status, bot, friend)
	if bot.DeviceId == friend.DeviceId {
		status = consts.FriendReject
		return errors.New("You can't add yourself as a friend")
	}
	if !s.verify(ctx, msg.Address) {
		status = consts.FriendReject
		return errors.New("The number of friends has reached the limit")
	}
	if !s.verify(ctx, msg.FriendAddress) {
		status = consts.FriendReject
		return errors.New("The maximum number of friends has been reached")
	}
	return err
}

func (s *sBot) DelFriend(ctx context.Context, address, friendAddress string) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Friend.Ctx(ctx).Where(
			dao.Friend.Ctx(ctx).Builder().Where(dao.Friend.Columns().Address, address).Where(dao.Friend.Columns().FriendAddress, friendAddress),
		).WhereOr(
			dao.Friend.Ctx(ctx).Builder().Where(dao.Friend.Columns().Address, friendAddress).Where(dao.Friend.Columns().FriendAddress, address),
		).Delete()
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		bot, err := s.FindByAddress(ctx, address)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		friend, err := s.FindByAddress(ctx, friendAddress)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		err = s.subBuff(ctx, *bot)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		err = s.subBuff(ctx, *friend)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		return nil
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	return nil
}

func (s *sBot) addFriendCallback(ctx context.Context, bot, friend entity.Bot) error {
	record := entity.Friend{
		Address:       bot.Address,
		FriendAddress: friend.Address,
		Active:        friend.Active,
	}
	friendRecord := entity.Friend{
		Address:       friend.Address,
		FriendAddress: bot.Address,
		Active:        bot.Active,
	}
	_, err := dao.Friend.Ctx(ctx).Data([]entity.Friend{record, friendRecord}).Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	err = s.addBuff(ctx, bot)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	err = s.addBuff(ctx, friend)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	return nil
}

func (s *sBot) verify(ctx context.Context, address string) bool {
	count, err := dao.Friend.Ctx(ctx).Where(dao.Friend.Columns().Address, address).Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if count >= consts.FriendLimit {
		return false
	}
	return true
}
