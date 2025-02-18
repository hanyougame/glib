package rocketmqx

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/hanyougame/glib/utils"
	"github.com/zeromicro/go-zero/core/logx"
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
	consumers.ConsumeTimeout = utils.Ternary(commonConfig.ConsumeTimeout > 0, commonConfig.ConsumeTimeout, consumers.ConsumeTimeout)
	consumers.NameServers = utils.Ternary(len(commonConfig.NameServers) > 0, commonConfig.NameServers, consumers.NameServers)
	consumers.ConsumeFromWhere = utils.Ternary(commonConfig.ConsumeFromWhere > 0, commonConfig.ConsumeFromWhere, consumers.ConsumeFromWhere)
	consumers.SecretKey = utils.Ternary(commonConfig.SecretKey != "", commonConfig.SecretKey, consumers.SecretKey)
	consumers.AccessKey = utils.Ternary(commonConfig.AccessKey != "", commonConfig.AccessKey, consumers.AccessKey)
	consumers.GoroutineNums = utils.Ternary(commonConfig.GoroutineNums > 0, commonConfig.GoroutineNums, consumers.GoroutineNums)
	consumers.PullBatchSize = utils.Ternary(commonConfig.PullBatchSize > 0, commonConfig.PullBatchSize, consumers.PullBatchSize)
	consumers.RetryNum = utils.Ternary(commonConfig.RetryNum > 0, commonConfig.RetryNum, consumers.RetryNum)
	consumers.WorkerNum = utils.Ternary(commonConfig.WorkerNum > 0, commonConfig.WorkerNum, consumers.WorkerNum)
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
