package integration

import (
	"basic-go/webook/internal/integration/startup"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service/sms"
	smsmocks "basic-go/webook/internal/service/sms/mocks"
	"errors"
	"testing"
	"time"

	"github.com/ecodeclub/ekit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type AsyncSMSTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *AsyncSMSTestSuite) SetupSuite() {
	s.db = startup.InitDB()
}

func (s *AsyncSMSTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE table `async_sms`")
}

func (s *AsyncSMSTestSuite) TestSend() {
	t := s.T()
	testCases := []struct {
		name string

		// 虽然是集成测试，但是也不想真的发短信，所以用 mock
		mock func(ctrl *gomock.Controller) sms.Service

		tplId   string
		args    []string
		numbers []string

		wantErr error
	}{
		{
			name: "异步",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				return svc
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := startup.InitAsyncSmsService(tc.mock(ctrl))
			err := svc.Send(context.Background(), tc.tplId, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func (s *AsyncSMSTestSuite) TestAsyncCycle() {
	now := time.Now()
	testCases := []struct {
		name string
		// 虽然是集成测试，但是也不想真的发短信，所以用 mock
		mock func(ctrl *gomock.Controller) sms.Service
		// 准备数据
		before func(t *testing.T)
		after  func(t *testing.T)
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123",
					[]string{"123456"}, []string{"15212345678"}).
					Return(nil)
				return svc
			},
			before: func(t *testing.T) {
				// 准备一条数据
				err := s.db.Create(&dao.AsyncSms{
					Id: 1,
					Config: sqlx.JsonColumn[dao.SmsConfig]{
						Val: dao.SmsConfig{
							TplId:   "123",
							Args:    []string{"123456"},
							Numbers: []string{"15212345678"},
						},
						Valid: true,
					},
					RetryMax: 3,
					Status:   0,
					Ctime:    now.Add(-time.Minute * 2).UnixMilli(),
					Utime:    now.Add(-time.Minute * 2).UnixMilli(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				var as dao.AsyncSms
				err := s.db.Where("id=?", 1).First(&as).Error
				assert.NoError(t, err)
				assert.Equal(t, uint8(2), as.Status)
			},
		},
		{
			name: "发送失败，标记为失败",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123",
					[]string{"123456"}, []string{"15212345678"}).
					Return(errors.New("模拟失败"))
				return svc
			},
			before: func(t *testing.T) {
				// 准备一条数据
				err := s.db.Create(&dao.AsyncSms{
					Id: 2,
					Config: sqlx.JsonColumn[dao.SmsConfig]{
						Val: dao.SmsConfig{
							TplId:   "123",
							Args:    []string{"123456"},
							Numbers: []string{"15212345678"},
						},
						Valid: true,
					},
					RetryMax: 3,
					RetryCnt: 2,
					Status:   0,
					Ctime:    now.Add(-time.Minute * 2).UnixMilli(),
					Utime:    now.Add(-time.Minute * 2).UnixMilli(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				var as dao.AsyncSms
				err := s.db.Where("id=?", 2).First(&as).Error
				assert.NoError(t, err)
				assert.Equal(t, uint8(1), as.Status)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tc.before(t)
			svc := startup.InitAsyncSmsService(tc.mock(ctrl))
			defer tc.after(t)
			svc.AsyncSend()
		})
	}
}

func TestAsyncSmsService(t *testing.T) {
	suite.Run(t, &AsyncSMSTestSuite{})
}
