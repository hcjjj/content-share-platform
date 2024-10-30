package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service

	// v1 的字段
	// 当前服务商下标
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

// 每次都从头开始轮询，绝大多数请求会在 svcs[0] 就成功，负载不均衡

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		// 输出日志，做好监控
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}

// 起始下标轮询
// 并且出错也轮询

func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 取下一个节点为起始节点
	// 不让每次都从 0 开始
	// 原子操作是轻量级并发工具
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 要迭代 length
	for i := idx; i < idx+length; i++ {
		// 取余数来计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是被取消，后者是超时
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}
