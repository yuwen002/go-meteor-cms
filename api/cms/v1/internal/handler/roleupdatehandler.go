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

func roleUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRoleUpdateLogic(r.Context(), svcCtx)
		resp, err := l.RoleUpdate(&req)
		if err != nil {
			var bizErr *common.BizError
			if errors.As(err, &bizErr) {
				common.Fail(w, bizErr.Code, bizErr.Msg)
			} else {
				common.Fail(w, common.ErrInternalServer, err.Error())
			}
			return
		}
		common.Ok(w, resp)
	}
}
