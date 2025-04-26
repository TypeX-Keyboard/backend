package system

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/system/v1"
)

func (c *ControllerV1) Settings(ctx context.Context, req *v1.SettingsReq) (res *v1.SettingsRes, err error) {
	res = &v1.SettingsRes{
		Timezone:           model.SettingConfig.Timezone,
		WorkReportInterval: model.SettingConfig.WorkReportInterval,
	}
	val := ctx.Value(consts.OS)
	if val == nil {
		return res, gerror.NewCode(gcode.CodeMissingParameter)
	}
	os := val.(string)
	version, err := service.System().GetVersionSetting(ctx, os)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.Version = version
	res.MiningRateOfH = consts.MiningRateOfH
	return res, gerror.NewCode(gcode.CodeOK)
}
