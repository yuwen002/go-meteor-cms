package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaLogic) GetCaptcha() (resp *types.CaptchaResp, err error) {
	// 生成验证码
	captchaID, captchaBase64, err := utils.GenerateCaptcha()
	if err != nil {
		l.Logger.Errorf("生成验证码失败: %v", err)
		return nil, err
	}

	return &types.CaptchaResp{
		CaptchaID:     captchaID,
		CaptchaBase64: captchaBase64,
	}, nil
}
