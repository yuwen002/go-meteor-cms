// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/logic"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func registerAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			// 判断是否是 validator.ValidationErrors
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				// 拼接提示信息
				msg := ""
				for _, e := range errs {
					switch e.Tag() {
					case "required":
						msg = e.Field() + "不能为空"
					case "min":
						msg = e.Field() + "长度不能少于" + e.Param() + "位"
					case "email":
						msg = e.Field() + "格式不正确"
					default:
						msg = "参数错误"
					}
					break // 只返回第一个错误
				}
				common.Fail(w, 40000, msg)
				return
			}
			common.Fail(w, 40000, "参数错误")
			return
		}

		l := logic.NewRegisterAdminLogic(r.Context(), svcCtx)
		resp, err := l.RegisterAdmin(&req)
		if err != nil {
			var be *common.BizError
			if errors.As(err, &be) {
				common.Fail(w, be.Code, be.Msg)
				return
			}
			common.Fail(w, 50000, "系统错误，请稍后重试")
			return
		}

		common.Ok(w, resp)
	}
}
