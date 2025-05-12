// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
)

type (
	ISystem interface {
		GetVersionSetting(ctx context.Context, os string) (model.Version, error)
		GetSystemSetting(ctx context.Context, settingKey string) (e *entity.SystemSetting, err error)
	}
)

var (
	localSystem ISystem
)

func System() ISystem {
	if localSystem == nil {
		panic("implement not found for interface ISystem, forgot register?")
	}
	return localSystem
}

func RegisterSystem(i ISystem) {
	localSystem = i
}
