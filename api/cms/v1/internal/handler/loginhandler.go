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

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			common.Fail(w, 40000, "参数错误")
			return
		}

		var fieldChName = map[string]string{
			"Username": "用户名",
			"Password": "密码",
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				e := errs[0] // 只取第一个错误

				field := e.Field()
				cnName := fieldChName[field]
				if cnName == "" {
					cnName = field // 没映射就用原始名字
				}

				msg := ""
				switch e.Tag() {
				case "required":
					msg = cnName + "不能为空"
				default:
					msg = "参数错误"
				}

				common.Fail(w, 40000, msg)
				return
			}

			common.Fail(w, 40000, "参数错误")
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
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
