// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/logic"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/internal/common"
)

func departmentTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewDepartmentTreeLogic(r.Context(), svcCtx)
		resp, err := l.DepartmentTree()
		if err != nil {
			common.Fail(w, common.ErrDepartmentListFail, err.Error())
			return
		}
		common.Ok(w, resp)
	}
}
