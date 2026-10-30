// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/sysconfig"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ConfigList 查询系统配置，支持按键/名称模糊搜索
func (l *ConfigLogic) ConfigList(req *types.ConfigListReq) (resp *types.ConfigListResp, err error) {
	query := l.svcCtx.EntClient.SysConfig.Query()
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		query = query.Where(sysconfig.Or(
			sysconfig.ConfigKeyContainsFold(kw),
			sysconfig.NameContainsFold(kw),
		))
	}

	rows, err := query.
		Order(ent.Asc(sysconfig.FieldConfigKey)).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list config failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	list := make([]types.ConfigItem, 0, len(rows))
	for _, c := range rows {
		list = append(list, types.ConfigItem{
			ID:          c.ID,
			ConfigKey:   c.ConfigKey,
			ConfigValue: c.ConfigValue,
			Name:        c.Name,
			Remark:      c.Remark,
			IsSystem:    c.IsSystem,
			CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		})
	}
	return &types.ConfigListResp{List: list}, nil
}

// ConfigCreate 新增配置项，键全局唯一
func (l *ConfigLogic) ConfigCreate(req *types.ConfigCreateReq) error {
	key := strings.TrimSpace(req.ConfigKey)
	name := strings.TrimSpace(req.Name)
	if key == "" || name == "" {
		return common.NewBizError(common.ErrInvalidParams)
	}
	exists, err := l.svcCtx.EntClient.SysConfig.
		Query().
		Where(sysconfig.ConfigKeyEQ(key)).
		Exist(l.ctx)
	if err != nil {
		logx.Errorf("check config key failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if exists {
		return common.NewBizError(common.ErrConfigKeyExists)
	}

	if _, err := l.svcCtx.EntClient.SysConfig.
		Create().
		SetConfigKey(key).
		SetConfigValue(req.ConfigValue).
		SetName(name).
		SetRemark(req.Remark).
		SetIsSystem(false).
		Save(l.ctx); err != nil {
		logx.Errorf("create config failed: %v", err)
		return common.NewBizError(common.ErrConfigCreateFail)
	}
	return nil
}

// ConfigUpdate 更新配置项，仅覆盖传入的字段；系统内置项仅允许改值
func (l *ConfigLogic) ConfigUpdate(req *types.ConfigUpdateReq) error {
	c, err := l.svcCtx.EntClient.SysConfig.Get(l.ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return common.NewBizError(common.ErrConfigNotFound)
		}
		logx.Errorf("get config failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}

	update := c.Update()
	if req.ConfigValue != "" {
		update = update.SetConfigValue(req.ConfigValue)
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		update = update.SetName(name)
	}
	if req.Remark != "" {
		update = update.SetRemark(req.Remark)
	}
	if _, err := update.Save(l.ctx); err != nil {
		logx.Errorf("update config failed: %v", err)
		return common.NewBizError(common.ErrConfigUpdateFail)
	}
	return nil
}

// ConfigDelete 删除配置项；系统内置配置禁止删除，避免误删导致功能异常
func (l *ConfigLogic) ConfigDelete(req *types.ConfigDetailReq) error {
	c, err := l.svcCtx.EntClient.SysConfig.Get(l.ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return common.NewBizError(common.ErrConfigNotFound)
		}
		logx.Errorf("get config failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if c.IsSystem {
		return common.NewBizError(common.ErrConfigDeleteSystem)
	}
	if err := l.svcCtx.EntClient.SysConfig.DeleteOneID(req.ID).Exec(l.ctx); err != nil {
		logx.Errorf("delete config failed: %v", err)
		return common.NewBizError(common.ErrConfigUpdateFail)
	}
	return nil
}
