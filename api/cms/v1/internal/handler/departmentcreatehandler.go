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

func departmentCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DepartmentCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
			return
		}

		l := logic.NewDepartmentCreateLogic(r.Context(), svcCtx)
		resp, err := l.DepartmentCreate(&req)
		if err != nil {
			// 如果是业务错误，使用业务错误码，否则使用内部服务器错误
			var bizErr *common.BizError
			if errors.As(err, &bizErr) {
				common.Fail(w, bizErr.Code, bizErr.Msg)
			}
			return
		}
		common.Ok(w, resp)
	}
}
