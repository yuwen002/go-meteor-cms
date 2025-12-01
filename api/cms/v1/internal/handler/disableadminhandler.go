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

func disableAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DisableAdminReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
			return
		}

		l := logic.NewDisableAdminLogic(r.Context(), svcCtx)
		resp, err := l.DisableAdmin(&req)
		if err != nil {
			// 检查错误类型并返回对应的错误码和消息
			var bizErr *common.BizError
			if errors.As(err, &bizErr) {
				common.Fail(w, bizErr.Code, bizErr.Msg)
			} else {
				common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			}
			return
		}
		common.Ok(w, resp)
	}
}
