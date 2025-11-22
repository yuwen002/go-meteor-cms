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
			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
			return
		}

		var fieldChName = map[string]string{
			"Username":     "用户名",
			"PasswordHash": "密码",
			"Email":        "邮箱",
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
				case "min":
					msg = cnName + "长度不能少于" + e.Param() + "位"
				case "email":
					msg = cnName + "格式不正确"
				default:
					msg = "参数错误"
				}

				common.Fail(w, common.ErrInvalidParams, msg)
				return
			}

			common.Fail(w, common.ErrInvalidParams, common.GetErrorMessage(common.ErrInvalidParams))
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
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}

		common.Ok(w, resp)
	}
}
