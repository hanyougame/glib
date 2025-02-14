package rocketmqx

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/hanyougame/glib/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

// NewConsumer 创建消费者实例
func NewConsumer(options []consumer.Option, topic string, messageSelector consumer.MessageSelector, handler MessageHandler) (rocketmq.PushConsumer, error) {
	cc, err := rocketmq.NewPushConsumer(options...)
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
func StartConsumer(nameServers []string, consumers []ConsumerConfig, handlers map[string]MessageHandler) []rocketmq.PushConsumer {
	var (
		ccs []rocketmq.PushConsumer
		mu  sync.Mutex
		wg  sync.WaitGroup
	)

	// 遍历所有消费者配置
	for _, cc := range consumers {
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

				// 如果配置中没有 NameServers，则使用传入的 nameServers
				if len(consumerConfig.NameServers) == 0 {
					consumerConfig.NameServers = nameServers
				}

				// 确定该消费者实例的 MessageSelector
				messageSelector := getMessageSelector(consumerConfig)
				options := getOptions(consumerConfig)

				// 创建并启动消费者实例
				consumerInstance, err := NewConsumer(options, consumerConfig.Topic, messageSelector, handler)
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

// getOptions 获取消费者设置
func getOptions(config ConsumerConfig) []consumer.Option {
	list := []consumer.Option{
		consumer.WithGroupName(config.GroupName),
		consumer.WithNameServer(config.NameServers),
		consumer.WithInstance(fmt.Sprintf("%s_%s", config.InstanceName, uuid.NewString())),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromWhere(config.ConsumeFromWhere)),
		consumer.WithConsumeMessageBatchMaxSize(config.ConsumeMessageBatchMaxSize),
		consumer.WithPullBatchSize(int32(config.PullBatchSize)),
		consumer.WithConsumeGoroutineNums(config.GoroutineNums),
		consumer.WithRetry(config.RetryNum),
		consumer.WithConsumeTimeout(time.Duration(config.ConsumeTimeout) * time.Second),
	}
	if config.AccessKey != "" && config.SecretKey != "" {
		list = append(list, consumer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}))
	}
	return list
}
