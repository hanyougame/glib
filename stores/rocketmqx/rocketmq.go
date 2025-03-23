package rocketmqx

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	v2 "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"github.com/google/uuid"
	"github.com/hanyougame/glib/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

// NewConsumer 创建消费者实例
func NewConsumer(config ConsumerConfig, topic string, messageSelector consumer.MessageSelector, handler MessageHandler) (rocketmq.PushConsumer, error) {
	cc, err := rocketmq.NewPushConsumer(
		// 设置消费者组
		consumer.WithGroupName(config.GroupName),
		// 设置服务地址
		consumer.WithNsResolver(primitive.NewPassthroughResolver(config.NameServers)),
		// 设置acl权限
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: config.SecretKey,
			AccessKey: config.AccessKey,
		}),
		// 设置从起始位置开始消费
		consumer.WithConsumeFromWhere(consumer.ConsumeFromWhere(config.ConsumeFromWhere)),
		consumer.WithConsumeGoroutineNums(config.GoroutineNums),
		consumer.WithRetry(config.RetryNum),
		consumer.WithConsumeTimeout(time.Duration(config.ConsumeTimeout)*time.Second),
		consumer.WithPullBatchSize(int32(utils.Ternary(config.PullBatchSize > 32, 32, config.PullBatchSize))),
		// 设置消费模式（默认集群模式）
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithInstance(config.GroupName+"_"+uuid.New().String()),
	)
	if err != nil {
		return nil, fmt.Errorf("创建消费者失败: %w", err)
	}
	if err = cc.Subscribe(topic, messageSelector, handler); err != nil {
		return nil, fmt.Errorf("订阅主题 %s 失败: %w", topic, err)
	}
	if err = cc.Start(); err != nil {
		return nil, fmt.Errorf("启动消费者失败: %w", err)
	}
	return cc, nil
}

// StartConsumer 启动多个消费者
func StartConsumer(commonConfig RocketMQX, consumers []ConsumerConfig, handlers map[string]MessageHandler) []rocketmq.PushConsumer {
	var (
		ccs []rocketmq.PushConsumer
		mu  sync.Mutex
		wg  sync.WaitGroup
	)

	// 遍历所有消费者配置
	for _, cc := range consumers {
		cc = getConfig(commonConfig, cc)
		// 启动 WorkerNum 数量的消费者实例
		for i := 0; i < cc.WorkerNum; i++ {
			wg.Add(1) // 增加计数

			go func(index int, consumerConfig ConsumerConfig) {
				defer wg.Done() // 完成时减去计数

				// 判断是否存在处理函数
				handler, exists := handlers[consumerConfig.Handler]
				if !exists {
					logx.Infof("处理函数 %s 未找到", consumerConfig.Handler)
					return
				}

				// 确定该消费者实例的 MessageSelector
				messageSelector := getMessageSelector(consumerConfig)

				// 创建并启动消费者实例
				consumerInstance, err := NewConsumer(consumerConfig, consumerConfig.Topic, messageSelector, handler)
				if err != nil {
					logx.Errorf("创建主题 %s 的消费者失败: %v", consumerConfig.Topic, err)
					return
				}

				// 使用加锁确保 ccs 列表线程安全
				mu.Lock() // 加锁
				ccs = append(ccs, consumerInstance)
				mu.Unlock() // 解锁
			}(i, cc)
		}
	}

	// 等待所有消费者完成启动
	wg.Wait()

	return ccs
}

func getConfig(commonConfig RocketMQX, consumers ConsumerConfig) ConsumerConfig {
	consumers.ConsumeTimeout = utils.Ternary(consumers.ConsumeTimeout > 0, consumers.ConsumeTimeout, commonConfig.ConsumeTimeout)
	consumers.NameServers = utils.Ternary(len(consumers.NameServers) > 0, consumers.NameServers, commonConfig.NameServers)
	consumers.ConsumeFromWhere = utils.Ternary(consumers.ConsumeFromWhere > 0, consumers.ConsumeFromWhere, commonConfig.ConsumeFromWhere)
	consumers.SecretKey = utils.Ternary(consumers.SecretKey != "", consumers.SecretKey, commonConfig.SecretKey)
	consumers.AccessKey = utils.Ternary(consumers.AccessKey != "", consumers.AccessKey, commonConfig.AccessKey)
	consumers.GoroutineNums = utils.Ternary(consumers.GoroutineNums > 0, consumers.GoroutineNums, commonConfig.GoroutineNums)
	consumers.PullBatchSize = utils.Ternary(consumers.PullBatchSize > 0, consumers.PullBatchSize, commonConfig.PullBatchSize)
	consumers.RetryNum = utils.Ternary(consumers.RetryNum > 0, consumers.RetryNum, commonConfig.RetryNum)
	consumers.WorkerNum = utils.Ternary(consumers.WorkerNum > 0, consumers.WorkerNum, commonConfig.WorkerNum)
	return consumers
}

// getMessageSelector 获取 MessageSelector
func getMessageSelector(config ConsumerConfig) consumer.MessageSelector {
	// 使用辅助函数确保配置项的默认值
	config.MessageSelector.Type = utils.Ternary(config.MessageSelector.Type == "", "TAG", config.MessageSelector.Type)
	config.MessageSelector.Expression = utils.Ternary(config.MessageSelector.Expression == "", "*", config.MessageSelector.Expression)
	return consumer.MessageSelector{
		Type:       consumer.ExpressionType(config.MessageSelector.Type),
		Expression: config.MessageSelector.Expression,
	}
}

func NewPullConsumer(commonConfig RocketMQX, config ConsumerConfig, handler PullMessageHandler) {
	_ = os.Setenv("mq.consoleAppender.enabled", "true")
	rmq_client.ResetLogger()
	config = getConfig(commonConfig, config)
	sleepTime := time.Duration(config.PullConsumerSleepTime) * time.Second
	simpleConsumer, err := rmq_client.NewSimpleConsumer(&rmq_client.Config{
		Endpoint:      utils.Ternary(len(config.NameServers) > 0, config.NameServers[0], ""),
		ConsumerGroup: config.GroupName,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    config.AccessKey,
			AccessSecret: config.SecretKey,
		},
	},
		rmq_client.WithAwaitDuration(time.Duration(config.AwaitDuration)*time.Second),
		rmq_client.WithSubscriptionExpressions(map[string]*rmq_client.FilterExpression{
			config.Topic: rmq_client.SUB_ALL,
		}),
	)
	if err != nil {
		logx.Errorf("初始化消费者失败，原因为：%s", err.Error())
		return
	}

	if err = simpleConsumer.Start(); err != nil {
		logx.Errorf("启动消费者失败，原因为：%s", err.Error())
		return
	}

	defer simpleConsumer.GracefulStop()
	for {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.ConsumeTimeout)*time.Second)
		mvs, err := simpleConsumer.Receive(ctx, int32(config.PullBatchSize), time.Duration(config.ConsumeTimeout)*time.Second)
		if err != nil {
			cancel()
			if strings.Contains(err.Error(), v2.Code_name[int32(v2.Code_MESSAGE_NOT_FOUND)]) {
				time.Sleep(sleepTime)
				continue
			}
			logx.Errorf("拉取消息失败，topic:%s,原因为:%s", config.Topic, err.Error())
			time.Sleep(sleepTime)
			continue
		}
		// ack message
		res, msgIDList, err := handler(ctx, mvs...)
		if err != nil {
			cancel()
			logx.Errorf("处理消息失败,topic:%s,原因为：%s", config.Topic, err.Error())
			time.Sleep(sleepTime)
			continue
		}
		// 如果全部成功
		if res == consumer.ConsumeSuccess {
			for _, mv := range mvs {
				if err = simpleConsumer.Ack(ctx, mv); err != nil {
					logx.Errorf("ack message failed, reason: %s, msgID:%s", err.Error(), mv.GetMessageId())
					continue
				}
			}
			cancel()
			continue
		} else {
			for _, mv := range mvs {
				if !slices.Contains(msgIDList, mv.GetMessageId()) {
					continue
				}
				if err = simpleConsumer.Ack(ctx, mv); err != nil {
					logx.Errorf("ack message failed, reason: %s, msgID:%s", err.Error(), mv.GetMessageId())
					continue
				}
			}
		}
		cancel()
		time.Sleep(sleepTime)
	}

}
