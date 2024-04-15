// Package web -----------------------------
// @file      : user_test.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-28 11:50
// -------------------------------------------
package web

import (
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123!",
				  "confirmPassword": "hcj123!"
				}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "用户已存在",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123!",
				  "confirmPassword": "hcj123!"
				}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱或者手机号码冲突",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail",
				  "password": "hcj123!",
				  "password": "hcj123!"
				}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不对",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123!",
				}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "两次输入的密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123!",
				  "confirmPassword": "hcj123"
				}`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123",
				  "confirmPassword": "hcj123"
				}`,
			wantCode: http.StatusOK,
			wantBody: "至少包含字母、数字、特殊字符，1-16位",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(),
					gomock.Any()).Return(errors.New("随便一个 error"))
				return usersvc
			},
			reqBody: `{
				  "email": "hcjjj@foxmail.com",
				  "password": "hcj123!",
				  "confirmPassword": "hcj123!"
				}`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 注册路由
			server := gin.Default()
			// SignUp 没用到 codeSvc
			h := NewUserHandler(tc.mock(ctrl), nil, nil)
			h.RegisterRoutes(server)

			// 构造请求
			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			// 设置数据是JSON格式
			req.Header.Set("Content-type", "application/json")
			//t.Log(req)
			require.NoError(t, err)
			// 验证响应
			resp := httptest.NewRecorder()
			//t.Log(resp)

			// 这是 HTTP 请求进入 GIN 框架的入口
			// 这样子调用的时候 GIN 就会处理这个请求，响应会写回到 resp
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}
