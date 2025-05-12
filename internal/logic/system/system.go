package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"keyboard-api-go/internal/service"
	"time"
)

func init() {
	service.RegisterSystem(New())
}

func New() service.ISystem {
	return &sSystem{}
}

type sSystem struct {
}

func (s *sSystem) GetVersionSetting(ctx context.Context, os string) (model.Version, error) {
	var version model.Version
	setting, err := s.GetSystemSetting(ctx, s.getVersionSettingKey(os))
	if err != nil || setting == nil {
		return model.Version{}, err
	}
	err = json.Unmarshal([]byte(setting.Value), &version)
	if err != nil {
		return model.Version{}, err
	}
	return version, nil
}

func (s *sSystem) getVersionSettingKey(os string) string {
	return fmt.Sprintf("%s%s", consts.HockeyAppVersionPrefix, os)
}

func (s *sSystem) GetSystemSetting(ctx context.Context, settingKey string) (e *entity.SystemSetting, err error) {
	// Query the system settings table
	err = dao.SystemSetting.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: time.Hour * 24,
		Name:     settingKey,
		Force:    false,
	}).Where(dao.SystemSetting.Columns().SettingKey, settingKey).Scan(&e)
	return
}
