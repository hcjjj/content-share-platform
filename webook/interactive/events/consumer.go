package events

import (
	"basic-go/webook/interactive/repository"
	"basic-go/webook/pkg/logger"
	"basic-go/webook/pkg/samarax"
	"context"
	"time"

	"github.com/IBM/sarama"
)

const TopicReadEvent = "article_read"

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l logger.LoggerV1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

//func (i *InteractiveReadEventConsumer) Start() error {
//	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
//	if err != nil {
//		return err
//	}
//	go func() {
//		er := cg.Consume(context.Background(),
//			[]string{TopicReadEvent},
//			samarax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
//		if er != nil {
//			i.l.Error("退出消费", logger.Error(er))
//		}
//	}()
//	return err
//}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			samarax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}

type ReadEvent struct {
	Aid int64
	Uid int64
}

//func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
//	events []ReadEvent) error {
//	bizs := make([]string, 0, len(events))
//	bizIds := make([]int64, 0, len(events))
//	for _, evt := range events {
//		bizs = append(bizs, "article")
//		bizIds = append(bizIds, evt.Aid)
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
//}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "article", event.Aid)
}
