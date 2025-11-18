package mq_model

import (
	"os"
)

var (
	// TopicUserLoginNotify 用户登录通知
	TopicUserLoginNotify = "user_login_notify"
	// TopicUserRechargeNotify 用户充值通知
	TopicUserRechargeNotify = "user_recharge_notify"
	// TopicUserWithdrawNotify 用户提现通知
	TopicUserWithdrawNotify = "user_withdraw_notify"
	// TopicUserBetSettlementNotify 用户投注结算通知
	TopicUserBetSettlementNotify = "user_bet_settlement_notify"
	// TopicUserBetStatNotify 投注统计专用 会接收到一笔投注记录 更新的通知 UserBetSettlementNotify
	TopicUserBetStatNotify = "user_bet_stat_notify"
	// TopicUserPromotionBonusNotify 用户优惠活动领取通知
	TopicUserPromotionBonusNotify = "user_promotion_bonus_notify"
	// TopicUserRegisterNotify 用户注册通知
	TopicUserRegisterNotify = "user_register_notify"
	// TopicPublishMqtt 发送mqtt消息
	TopicPublishMqtt = "publish_mqtt"
	// TopicUserRechargeLogNotify 用户充值日志
	TopicUserRechargeLogNotify = "user_recharge_log_notify"
	// TopicUserRechargeTimeoutNotify 用户充值超时通知
	TopicUserRechargeTimeoutNotify = "user_recharge_timeout_notify"
	// TopicUserSingleGameBetNotify  用户单一钱包游戏投注相关通知·
	TopicUserSingleGameBetNotify = "user_single_game_bet_notify"
	// TopicUserSingleGameBetSettleNotify  用户单一钱包游戏结算相关通知
	TopicUserSingleGameBetSettleNotify = "user_single_game_bet_settle_notify"
	// TopicUserBetNotify 用户游戏投注通知
	TopicUserBetNotify = "user_bet_notify"
	// TopicUserRedPacketNotify 用户红包活动通知
	TopicUserRedPacketNotify = "user_red_packet_notify"
	// TopicUserRedPacketCountdownNotify 用户红包活动倒计时通知
	TopicUserRedPacketCountdownNotify = "user_red_packet_countdown_notify"
	// TopicUserRecallNotify 用户召回通知
	TopicUserRecallNotify = "user_recall_notify"
	// TopicUserProfileAuthNotify 用户资料认证
	TopicUserProfileAuthNotify = "user_profile_auth_notify"
	// TopicUserProfileUpdateNotify 用户资料更新
	TopicUserProfileUpdateNotify = "user_profile_update_notify"
	// TopicPvReportNotify pv上报
	TopicPvReportNotify = "pv_report_notify"
	// TopicUserBalanceChange 用户余额变动
	TopicUserBalanceChange = "user_balance_change_%d"

	TopicUserDepositTripartite = "user_dep_wdl_change_topic"
)

func UpdateTopicPrefix(prefixes ...string) (prefix string) {
	if len(prefixes) > 0 && prefixes[0] != "" {
		prefix = prefixes[0]
	}
	if prefix == "" {
		prefix = os.Getenv("ROCKETMQ_TOPIC_PREFIX")
	}

	TopicUserLoginNotify = prefix + TopicUserLoginNotify
	TopicUserRechargeNotify = prefix + TopicUserRechargeNotify
	TopicUserWithdrawNotify = prefix + TopicUserWithdrawNotify
	TopicUserBetSettlementNotify = prefix + TopicUserBetSettlementNotify
	TopicUserBetStatNotify = prefix + TopicUserBetStatNotify
	TopicUserPromotionBonusNotify = prefix + TopicUserPromotionBonusNotify
	TopicUserRegisterNotify = prefix + TopicUserRegisterNotify
	TopicPublishMqtt = prefix + TopicPublishMqtt
	TopicUserRechargeLogNotify = prefix + TopicUserRechargeLogNotify
	TopicUserRechargeTimeoutNotify = prefix + TopicUserRechargeTimeoutNotify
	TopicUserSingleGameBetNotify = prefix + TopicUserSingleGameBetNotify
	TopicUserSingleGameBetSettleNotify = prefix + TopicUserSingleGameBetSettleNotify
	TopicUserBetNotify = prefix + TopicUserBetNotify
	TopicUserRedPacketNotify = prefix + TopicUserRedPacketNotify
	TopicUserRedPacketCountdownNotify = prefix + TopicUserRedPacketCountdownNotify
	TopicUserRecallNotify = prefix + TopicUserRecallNotify
	TopicUserProfileAuthNotify = prefix + TopicUserProfileAuthNotify
	TopicUserProfileUpdateNotify = prefix + TopicUserProfileUpdateNotify
	TopicPvReportNotify = prefix + TopicPvReportNotify
	TopicUserBalanceChange = prefix + TopicUserBalanceChange
	TopicUserDepositTripartite = prefix + TopicUserDepositTripartite

	return
}
