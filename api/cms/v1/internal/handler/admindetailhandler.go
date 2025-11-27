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

func adminDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, 40000, "参数错误")
			return
		}

		l := logic.NewAdminDetailLogic(r.Context(), svcCtx)
		resp, err := l.AdminDetail(&req)
		if err != nil {
			var be *common.BizError
			if errors.As(err, &be) {
				common.Fail(w, be.Code, be.Msg)
				return
			}
			common.Fail(w, common.ErrAdminUserNotFound, common.GetErrorMessage(common.ErrAdminUserNotFound))
			return
		}
		common.Ok(w, resp)
	}
}
