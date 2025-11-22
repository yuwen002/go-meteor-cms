// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/logic"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/internal/common"
)

func GetCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.GetCaptcha()
		if err != nil {
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}
		common.Ok(w, resp)
	}
}
