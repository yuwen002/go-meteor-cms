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

// adminListHandler 获取管理员列表
// @Summary 获取管理员列表
// @Description 分页获取管理员列表，支持关键词搜索和状态过滤
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param keyword query string false "搜索关键词（用户名/昵称/邮箱）"
// @Param active query int false "过滤启用状态：1-启用，2-禁用"
// @Success 200 {object} types.AdminListResp "成功返回管理员列表"
// @Router /admin/list [get]
func adminListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
			return
		}

		l := logic.NewAdminListLogic(r.Context(), svcCtx)
		resp, err := l.AdminList(&req)
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
