package rocketmqx

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RocketMQX struct {
	NameServers  []string         `json:"name_servers"`
	ConsumerList []ConsumerConfig `json:"consumer_list"`
}

// ConsumerConfig 配置结构体
type ConsumerConfig struct {
	GroupName                  string   `json:"group_name"`                                 // 消费者组名称
	Topic                      string   `json:"topic"`                                      // 消费的 Topic
	NameServers                []string `json:"name_servers,optional"`                      // NameServer 地址列表
	WorkerNum                  int      `json:"worker_num,default=5"`                       // 消费者消费线程数
	Handler                    string   `json:"handler"`                                    // 消费方法
	RetryNum                   int      `json:"retry_num,default=5"`                        // 重试次数
	GoroutineNums              int      `json:"goroutine_nums,default=10"`                  // 最大线程数
	ConsumeFromWhere           int      `json:"consume_from_where,default=1"`               // 消费位置 0 从末尾消费 1 从开头消费
	ConsumeMessageBatchMaxSize int      `json:"consume_message_batch_max_size,default=500"` // 每次拉取消息的最大条数
	PullBatchSize              int      `json:"pull_batch_size,default=1024"`               // 每次拉取的消息体最大字节数
	InstanceName               string   `json:"instance_name,optional"`
	MessageSelector            struct {
		Expression string `json:"expression,default=*"`
		Type       string `json:"type,default=TAG"`
	} `json:"message_selector,omitempty"` // 客户端标识符 要保证每个消费者组中的唯一
	ConsumeTimeout int64 `json:"consume_timeout,default=60"`
}

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
