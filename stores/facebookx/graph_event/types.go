package graph_event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

type EventName string

const (
	EventNameRegister       EventName = "CompleteRegistration"  // 注册（标准事件）
	EventNameLogin          EventName = "H5LoginEvent"          // 登录（自定义事件）
	EventNameWithdraw       EventName = "H5WithdrawEvent"       // 提现（自定义事件）
	EventNameRecharge       EventName = "Purchase"              // 充值（标准事件）
	EventNameRechargeFirst  EventName = "AddToCart"             // 首充（标准事件）
	EventNameRechargeSecond EventName = "H5SecondRechargeEvent" // 二充（自定义事件）
	EventNameRechargeThird  EventName = "H5ThirdRechargeEvent"  // 三充（自定义事件）
	EventNameRechargeFourth EventName = "H5FourthRechargeEvent" // 四充（自定义事件）
)

var EventNames = []EventName{
	EventNameRegister,
	EventNameLogin,
	EventNameWithdraw,
	EventNameRecharge,
	EventNameRechargeFirst,
	EventNameRechargeSecond,
	EventNameRechargeThird,
	EventNameRechargeFourth,
}

type ActionSource string

const (
	ActionSourceEmail             ActionSource = "email"              // 转化通过邮件发生。
	ActionSourceWebsite           ActionSource = "website"            // 转化在网站上完成。
	ActionSourceApp               ActionSource = "app"                // 转化在移动应用上完成。
	ActionSourcePhoneCall         ActionSource = "phone_call"         // 转化通过电话完成。
	ActionSourceChat              ActionSource = "chat"               // 转化通过消息应用、短信或在线消息功能完成。
	ActionSourcePhysicalStore     ActionSource = "physical_store"     // 转化在实体店面对面完成。
	ActionSourceSystemGenerated   ActionSource = "system_generated"   // 转化自动发生，例如将订阅续费设置为每月自动支付。
	ActionSourceBusinessMessaging ActionSource = "business_messaging" // 转化是通过 Messenger、Instagram 或 WhatsApp 的直达广告产生的。
	ActionSourceOther             ActionSource = "other"              // 转化通过本文未列出的其他方式完成。

)

type UserGender string

const (
	UserGenderMale   UserGender = "m" // 男性
	UserGenderFemale UserGender = "f" // 女性
)

type EventsSendRequest struct {
	Data []*EventSendRequestData `json:"data"`
}

type EventSendRequestData struct {
	//事件名称
	//必要
	//标准事件或自定义事件的名称
	EventName EventName `json:"event_name"`
	//事件时间
	//必要
	//Unix 时间戳（以秒为单位），表示事件的实际发生时间。指定时间可能早于您将事件发送至 Facebook 的时间。其目的在于实现批处理和服务器性能优化。您必须以格林威治标准时间 (GMT) 时区格式发送该日期。
	//event_time 最长可比您将事件发送至 Facebook 的时间早 7 天。如果 data 中的任何 event_time 早于 7 天前，我们会为整个请求返回错误，而且不会处理任何事件。
	EventTime int64 `json:"event_time"`
	//用户信息参数
	//必要
	//包含客户信息数据的映射。
	UserData GraphEventUserData `json:"user_data"`
	//标准参数
	//非必要
	//包含事件其他企业数据的映射。
	CustomData GraphEventCustomData `json:"custom_data"`
	//非必要
	//表示不应将此事件用于广告投放优化的标记。设置为 true 时，该事件仅可用作归因条件。
	OptOut bool `json:"opt_out"`
	//非必要
	//此编号可以是广告主选定的任何唯一字符串。可以是订单号或交易编号
	EventID string `json:"event_id"`
	//必要
	//您可通过此字段指定转化发生的位置。了解事件发生的位置有助于确保将广告投放至正确的受众。使用转化 API，即表示您同意：就您所知 action_source 参数准确无误
	ActionSource ActionSource `json:"action_source"`
	//非必要。
	//事件发生的浏览器网址。该网址应与已认证的网域一致。
	//注意：event_source_url 是使用转化 API 分享网站事件的必要项。
	EventSourceUrl string `json:"event_source_url"`
}

// GraphEventUserData 用户信息数据，至少提供以下其中一项参数
type GraphEventUserData struct {
	//邮箱
	//必须进行哈希处理。去除所有前导和尾部空格。将所有字符转换为小写形式
	Email string `json:"em"`
	//电话
	//必须进行哈希处理。
	//移除符号、字母及任何前导零。手机号必须包含用于匹配的国家/地区代码（例如，美国手机号前面必须带有数字 1）。始终在客户的手机号中加入国家/地区代码，即便所有数据都来自同一国家/地区亦是如此。
	Phone string `json:"ph"`
	//姓氏
	//必须进行哈希处理。
	//推荐使用罗马字母 a-z 字符。仅限小写字母，且不可包含标点符号。若使用特殊字符，则须按 UTF-8 格式对文本进行编码
	FirstName string `json:"fn"`
	// 同FirstName
	LastName string `json:"ln"`
	//出生日期
	//必须进行哈希处理。
	//YYYYMMDD格式，如：19970216
	Birthday string `json:"db"`
	//性别
	//必须进行哈希处理。
	//f 表示女性，m 表示男性
	Gender UserGender `json:"ge"`
	//外部编号
	//推荐进行哈希处理。
	//可以是广告主提供的任何唯一编号，如会员编号、用户编号和外部 Cookie 编号。您可以为给定事件发送一或多个外部编号。
	//如果是通过其他渠道发送外部编号，此编号的格式应与通过转化 API 发送时的格式相同。
	ExternalID []string `json:"external_id"`
	//IP地址
	//无需进行哈希处理。
	//必须是有效的 IPV4 或 IPV6 地址。若用户启用了 IPV6，请优先使用 IPV6
	IP string `json:"client_ip_address"`
	//点击编号
	//无需进行哈希处理。
	//Facebook 点击编号值存储在您网域下的 _fbc 浏览器 Cookie 中
	FBC string `json:"fbc"`
	//浏览器编号
	//无需进行哈希处理。
	//Facebook 浏览器编号值存储在您网域下的 _fbp 浏览器 Cookie 中。
	FBP string `json:"fbp"`
	//订阅编号
	//无需进行哈希处理。
	//此交易中用户的订阅编号；类似于单件商品的订单编号。
	SubscriptionID string `json:"subscription_id"`

	//客户端用户代理程序
	//无需进行哈希处理。
	//与事件对应的浏览器的用户代理程序。对于使用转化 API 分享的网站事件，必须填写 client_user_agent 参数。
	//请注意：此信息会自动添加到通过浏览器发送的事件中，但对于通过服务器发送的事件，您必须手动配置该信息。
	ClientUserAgent string `json:"client_user_agent"`
}

// GraphEventCustomData 标准参数
type GraphEventCustomData struct {
	//币种
	//对于购买事件为必要参数。
	//指定 value 的货币（如果适用）。货币必须是有效的 ISO 4217 三位数货币代码。
	Currency string `json:"currency"`
	//数值
	//对于购买事件或使用价值优化的事件必备。
	//与事件相关的一个数值。该数值须代表一个货币金额。
	Value float64 `json:"value"`
	//与事件相关的内容编号，例如 AddToCart 事件中商品的商品 SKU
	//如果 content_type 为 product，则您的内容编号必须为包含单个字符串值的数组。
	//如果该参数不是 product，该数组可以包含任意数量的字符串值。
	ContentIds []string `json:"content_ids"`
}

func (reqs *EventsSendRequest) verify() error {
	if len(reqs.Data) == 0 {
		return fmt.Errorf("data must not be empty")
	}
	for _, req := range reqs.Data {
		if stringx.NotEmpty(req.UserData.Email) {
			req.UserData.Email = hashData(req.UserData.Email)
		}
		if stringx.NotEmpty(req.UserData.FirstName) {
			req.UserData.FirstName = hashData(req.UserData.FirstName)
		}
		if stringx.NotEmpty(req.UserData.LastName) {
			req.UserData.LastName = hashData(req.UserData.LastName)
		}
		if stringx.NotEmpty(req.UserData.Phone) {
			req.UserData.Phone = hashData(req.UserData.Phone)
		}

		var undefinedEvent = true
		for _, eventName := range EventNames {
			if eventName == req.EventName {
				undefinedEvent = false
				break
			}
		}
		if undefinedEvent {
			return fmt.Errorf("undefined event name: %s", req.EventName)
		}
	}

	return nil
}

func hashData(data string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(data))))
	return hex.EncodeToString(hash[:])
}
