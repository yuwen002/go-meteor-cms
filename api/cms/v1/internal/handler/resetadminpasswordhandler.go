// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"errors"
	"net/http"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/logic"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func resetAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetAdminPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
			return
		}

		l := logic.NewResetAdminPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ResetAdminPassword(&req)
		if err != nil {
			var be *common.BizError
			if errors.As(err, &be) {
				common.Fail(w, be.Code, be.Msg)
				return
			}
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}
		common.Ok(w, resp)
	}
}
