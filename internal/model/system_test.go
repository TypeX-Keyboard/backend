package model

import (
	"encoding/json"
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model/entity"
	"testing"
)

func TestSeedVersion(t *testing.T) {
	ctx := gctx.GetInitCtx()
	data := make([]entity.SystemSetting, 0)
	version := Version{
		CurrentVersion: "v1.0.0",
		LatestVersion:  "v1.0.0",
		DownloadURL:    "",
		ForceUpdate:    false,
		ReleaseNotes:   []string{"fix bug"},
	}
	b, _ := json.Marshal(version)
	data = append(data, entity.SystemSetting{
		SettingKey: fmt.Sprintf("%s%s", consts.HockeyAppVersionPrefix, "Android"),
		Value:      string(b),
	})
	data = append(data, entity.SystemSetting{
		SettingKey: fmt.Sprintf("%s%s", consts.HockeyAppVersionPrefix, "Ios"),
		Value:      string(b),
	})
	_, err := dao.SystemSetting.Ctx(ctx).Data(data).Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
}
