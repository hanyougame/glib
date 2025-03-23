package rocketmqx

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	rmqclient "github.com/apache/rocketmq-clients/golang/v5"
)

type RocketMQX struct {
	NameServers      []string         `json:"name_servers"`
	ConsumerList     []ConsumerConfig `json:"consumer_list"`
	PullConsumerList []ConsumerConfig `json:"pull_consumer_list"`
	WorkerNum        int              `json:"worker_num,default=5"`         // 消费者消费线程数
	RetryNum         int              `json:"retry_num,default=5"`          // 重试次数
	GoroutineNums    int              `json:"goroutine_nums,default=10"`    // 最大线程数
	ConsumeFromWhere int              `json:"consume_from_where,default=1"` // 消费位置 0 从末尾消费 1 从开头消费
	PullBatchSize    int              `json:"pull_batch_size,default=32"`
	ConsumeTimeout   int64            `json:"consume_timeout,default=60"`
	SecretKey        string           `json:"secret_key,optional"`
	AccessKey        string           `json:"access_key,optional"`
}

// ConsumerConfig 配置结构体
type ConsumerConfig struct {
	GroupName        string   `json:"group_name"`                  // 消费者组名称
	Topic            string   `json:"topic"`                       // 消费的 Topic
	NameServers      []string `json:"name_servers,optional"`       // NameServer 地址列表
	WorkerNum        int      `json:"worker_num,optional"`         // 消费者消费线程数
	Handler          string   `json:"handler"`                     // 消费方法
	RetryNum         int      `json:"retry_num,optional"`          // 重试次数
	GoroutineNums    int      `json:"goroutine_nums,optional"`     // 最大线程数
	ConsumeFromWhere int      `json:"consume_from_where,optional"` // 消费位置 0 从末尾消费 1 从开头消费
	PullBatchSize    int      `json:"pull_batch_size,optional"`
	MessageSelector  struct {
		Expression string `json:"expression,default=*"`
		Type       string `json:"type,default=TAG"`
	} `json:"message_selector,omitempty"`
	ConsumeTimeout        int64  `json:"consume_timeout,default=60"`
	SecretKey             string `json:"secret_key,optional"`
	AccessKey             string `json:"access_key,optional"`
	AwaitDuration         int64  `json:"await_duration,default=5"`
	PullConsumerSleepTime int64  `json:"pull_consumer_sleep_time,default=2"`
}

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
type PullMessageHandler func(ctx context.Context, msgs ...*rmqclient.MessageView) (consumer.ConsumeResult, []string, error)
