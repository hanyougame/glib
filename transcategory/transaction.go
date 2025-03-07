package transcategory

// TransactionCategory 账变大类
type TransactionCategory int

var TransactionCategories = []TransactionCategory{
	TransactionCategoryCapitalSwitch,
	TransactionCategoryUserRecharge,
	TransactionCategoryUserWithdraw,
	TransactionCategoryBankMerchantSettlement,
	TransactionCategoryCapitalRevise,
	TransactionCategoryActivity,
	TransactionCategoryGoldReturn,
	TransactionCategoryRebate,
	TransactionCategoryInterest,
	TransactionCategoryTask,
	TransactionCategoryVipReward,
	TransactionCategoryRechargeBonus,
	TransactionCategoryClub,
	TransactionCategoryReward,
	TransactionCategoryGuaranteeClaim,
	TransactionCategoryProvidentFund,
	TransactionCategoryBlindBoxLottery,
	TransactionCategoryGameBet,
}

const (
	_                                         TransactionCategory = iota
	TransactionCategoryCapitalSwitch                              // 资金切换
	TransactionCategoryUserRecharge                               // 会员充值
	TransactionCategoryUserWithdraw                               // 会员提现
	TransactionCategoryBankMerchantSettlement                     // 银商结算
	TransactionCategoryCapitalRevise                              // 资金修正
	TransactionCategoryActivity                                   // 活动
	TransactionCategoryGoldReturn                                 // 返水
	TransactionCategoryRebate                                     // 返佣
	TransactionCategoryInterest                                   // 利息宝
	TransactionCategoryTask                                       // 任务
	TransactionCategoryVipReward                                  // VIP奖励
	TransactionCategoryRechargeBonus                              // 充值优惠
	TransactionCategoryClub                                       // 俱乐部
	TransactionCategoryReward                                     // 奖励
	TransactionCategoryGuaranteeClaim                             // 担保理赔
	TransactionCategoryProvidentFund                              // 公积金
	TransactionCategoryBlindBoxLottery                            // 盲盒抽奖
	TransactionCategoryGameBet                                    // 游戏投注
)

func (t TransactionCategory) Int() int {
	return int(t)
}

func (t TransactionCategory) Int64() int64 {
	return int64(t)
}

func (t TransactionCategory) String() string {
	switch t {
	case TransactionCategoryCapitalSwitch:
		return "资金切换"
	case TransactionCategoryUserRecharge:
		return "会员充值"
	case TransactionCategoryUserWithdraw:
		return "会员提现"
	case TransactionCategoryBankMerchantSettlement:
		return "银商结算"
	case TransactionCategoryCapitalRevise:
		return "资金修正"
	case TransactionCategoryActivity:
		return "活动"
	case TransactionCategoryGoldReturn:
		return "返水"
	case TransactionCategoryRebate:
		return "返佣"
	case TransactionCategoryInterest:
		return "利息宝"
	case TransactionCategoryTask:
		return "任务"
	case TransactionCategoryVipReward:
		return "VIP奖励"
	case TransactionCategoryRechargeBonus:
		return "充值优惠"
	case TransactionCategoryClub:
		return "俱乐部"
	case TransactionCategoryReward:
		return "奖励"
	case TransactionCategoryGuaranteeClaim:
		return "担保理赔"
	case TransactionCategoryProvidentFund:
		return "公积金"
	case TransactionCategoryBlindBoxLottery:
		return "盲盒抽奖"
	case TransactionCategoryGameBet:
		return "游戏投注"
	default:
		return ""
	}
}

// SubCategory 获取子类
func (t TransactionCategory) SubCategory() []TransactionSubCategory {
	switch t {
	case TransactionCategoryCapitalSwitch:
		return []TransactionSubCategory{
			TransactionSubCategoryWithdrawal, TransactionSubCategoryDeposits,
		}
	case TransactionCategoryUserRecharge:
		return []TransactionSubCategory{
			TransactionSubCategoryUalaTransfer, TransactionSubCategoryMercadoPagoTransfer, TransactionSubCategoryUPDAYTransfer, TransactionSubCategoryRechargeOnline,
		}
	case TransactionCategoryUserWithdraw:
		return []TransactionSubCategory{
			TransactionSubCategoryWithdrawFrozen, TransactionSubCategoryWithdrawReject,
			TransactionSubCategoryWithdrawDefrost, TransactionSubCategoryWithdrawSucceed,
		}
	case TransactionCategoryBankMerchantSettlement:
		return []TransactionSubCategory{
			TransactionSubCategoryBankMerchantTransfer, TransactionSubCategoryBankMerchantAddAmount, TransactionSubCategoryBankMerchantRecharge, // 银商充值
			TransactionSubCategoryBankMerchantGiveUser, TransactionSubCategoryBankMerchantSubtractAmount,
		}
	case TransactionCategoryCapitalRevise:
		return []TransactionSubCategory{
			TransactionSubCategoryManualAddAmount, TransactionSubCategoryManualAddOrder, TransactionSubCategoryBalanceRevise,
			TransactionSubCategoryManualAddRewardAmount, TransactionSubCategoryManualSubtractAmount, TransactionSubCategoryDeductAll, // 扣除全部资产
			TransactionSubCategorySurDeduct, TransactionSubCategoryManualPullBack, TransactionSubCategoryDeductExcessProfit,
		}
	case TransactionCategoryActivity:
		return []TransactionSubCategory{
			TransactionSubCategoryAgentActivity, TransactionSubCategoryLuckyBetActivity, TransactionSubCategoryInvestActivity,
			TransactionSubCategoryNewcomerRewardActivity, TransactionSubCategoryBenefitActivity, TransactionSubCategoryPromoteActivity,
			TransactionSubCategoryFeedbackRewardActivity, TransactionSubCategoryRedPacketActivity, TransactionSubCategoryBetActivity,
			TransactionSubCategoryLotteryAssistanceActivity, TransactionSubCategoryRankActivity, TransactionSubCategoryCustomizeActivity,
			TransactionSubCategoryBargainActivity, TransactionSubCategorySpinActivity, TransactionSubCategoryChannelRewardActivity,
			TransactionSubCategoryWordCollectionActivity, TransactionSubCategoryQuizActivity,
		}
	case TransactionCategoryGoldReturn:
		return []TransactionSubCategory{
			TransactionSubCategoryServiceChargeReceive,
		}
	case TransactionCategoryRebate:
		return []TransactionSubCategory{
			TransactionSubCategoryRebateSend, TransactionSubCategoryRebateReceive,
		}
	case TransactionCategoryInterest:
		return []TransactionSubCategory{
			TransactionSubCategoryInterestProfit, TransactionSubCategoryHallToInterest,
			TransactionSubCategoryInterestManualPullback, TransactionSubCategoryInterestToHall,
		}
	case TransactionCategoryTask:
		return []TransactionSubCategory{
			TransactionSubCategoryDailyTask, TransactionSubCategoryWeekLyTask, TransactionSubCategoryNewcomerProfitTask,
			TransactionSubCategoryAlivenessBox, TransactionSubCategorySecretTask,
		}
	case TransactionCategoryVipReward:
		return []TransactionSubCategory{
			TransactionSubCategoryVipMonthlyReward, TransactionSubCategoryVipDailyReward,
			TransactionSubCategoryVipWeeklyReward, TransactionSubCategoryVipUpgradeReward,
		}
	case TransactionCategoryRechargeBonus:
		return []TransactionSubCategory{
			TransactionSubCategoryRechargeBonus,
		}
	case TransactionCategoryClub:
		return []TransactionSubCategory{
			TransactionSubCategoryManualPullbackClub, TransactionSubCategoryHallToClub, TransactionSubCategoryClubToHall,
		}
	case TransactionCategoryReward:
		return []TransactionSubCategory{
			TransactionSubCategoryGiveUpReward, TransactionSubCategoryTransferOut, TransactionSubCategoryTransferIn,
		}
	case TransactionCategoryGuaranteeClaim:
		return []TransactionSubCategory{
			TransactionSubCategoryClaimFreeze, TransactionSubCategoryClaimDefrost, TransactionSubCategoryClaimFee,
			TransactionSubCategoryClaimScoresIncrease, TransactionSubCategoryClaimScoresDecrease,
		}
	case TransactionCategoryProvidentFund:
		return []TransactionSubCategory{
			TransactionSubCategoryRechargeBonus,
		}
	case TransactionCategoryBlindBoxLottery:
		return []TransactionSubCategory{
			TransactionSubCategoryBlindBoxLotteryDeduct, TransactionSubCategoryBlindBoxLotteryReward,
		}
	case TransactionCategoryGameBet:
		return []TransactionSubCategory{
			TransactionSubCategoryGameBet,
			TransactionSubCategoryGameBetSettlement,
			TransactionSubCategoryGameBetCancel,
		}
	default:
		return []TransactionSubCategory{}
	}
}

// TransactionSubCategory 账变小类
type TransactionSubCategory int

// 资金切换子类
const (
	SubCategoryCapitalSwitch         TransactionSubCategory = 1000
	TransactionSubCategoryWithdrawal                        = iota + SubCategoryCapitalSwitch // 转出
	TransactionSubCategoryDeposits                                                            // 转入

)

// 会员充值子类
const (
	SubCategoryUserRecharge                   TransactionSubCategory = 2000
	TransactionSubCategoryUalaTransfer                               = iota + SubCategoryUserRecharge // Uala转账
	TransactionSubCategoryMercadoPagoTransfer                                                         // Mercado Pago转账
	TransactionSubCategoryUPDAYTransfer                                                               // UPDAY转
	TransactionSubCategoryRechargeOnline                                                              // 在线充值
)

// 会员提现子类
const (
	SubCategoryUserWithdraw               TransactionSubCategory = 3000
	TransactionSubCategoryWithdrawFrozen                         = iota + SubCategoryUserWithdraw // 提现冻结
	TransactionSubCategoryWithdrawReject                                                          // 提现拒绝
	TransactionSubCategoryWithdrawDefrost                                                         // 提现解冻
	TransactionSubCategoryWithdrawSucceed                                                         // 提现成功
)

// 银商结算子类
const (
	SubCategoryBankMerchantSettlement                TransactionSubCategory = 4000
	TransactionSubCategoryBankMerchantTransfer                              = iota + SubCategoryBankMerchantSettlement // 转账给他人
	TransactionSubCategoryBankMerchantAddAmount                                                                        // 银商加款
	TransactionSubCategoryBankMerchantRecharge                                                                         // 银商充值
	TransactionSubCategoryBankMerchantGiveUser                                                                         // 银商赠送会员金额
	TransactionSubCategoryBankMerchantSubtractAmount                                                                   // 银商扣款
)

// 资金修正子类
const (
	SubCategoryCapitalRevise                    TransactionSubCategory = 5000
	TransactionSubCategoryManualAddAmount                              = iota + SubCategoryCapitalRevise // 人工加款
	TransactionSubCategoryManualAddOrder                                                                 // 手动补单
	TransactionSubCategoryBalanceRevise                                                                  // 修正负数余额
	TransactionSubCategoryManualAddRewardAmount                                                          // 奖励手动加款
	TransactionSubCategoryManualSubtractAmount                                                           // 手动扣款
	TransactionSubCategoryDeductAll                                                                      // 扣除全部资产
	TransactionSubCategorySurDeduct                                                                      // 追加扣除
	TransactionSubCategoryManualPullBack                                                                 // 人工拉回
	TransactionSubCategoryDeductExcessProfit                                                             // 扣除超额盈利
)

// 活动子类
const (
	SubCategoryActivity                             TransactionSubCategory = 6000
	TransactionSubCategoryAgentActivity                                    = iota + SubCategoryActivity // 代理活动
	TransactionSubCategoryLuckyBetActivity                                                              // 幸运注单活动
	TransactionSubCategoryInvestActivity                                                                // 投资活动
	TransactionSubCategoryNewcomerRewardActivity                                                        // 新人彩金
	TransactionSubCategoryBenefitActivity                                                               // 救援金活动
	TransactionSubCategoryPromoteActivity                                                               // 推广活动
	TransactionSubCategoryFeedbackRewardActivity                                                        // 有奖反馈活动
	TransactionSubCategoryRedPacketActivity                                                             // 红包活动
	TransactionSubCategoryBetActivity                                                                   // 打码活动
	TransactionSubCategoryLotteryAssistanceActivity                                                     // 抽奖助力
	TransactionSubCategoryRankActivity                                                                  // 排行榜活动
	TransactionSubCategoryCustomizeActivity                                                             // 自定义活动
	TransactionSubCategoryBargainActivity                                                               // 砍一刀
	TransactionSubCategorySpinActivity                                                                  // 转盘活动
	TransactionSubCategoryChannelRewardActivity                                                         // 渠道奖励
	TransactionSubCategoryWordCollectionActivity                                                        // 集字活动
	TransactionSubCategoryQuizActivity                                                                  // 竞猜活动
	TransactionSubCategoryRechargeActivity                                                              // 充值活动
	TransactionSubCategorySignInActivity                                                                // 签到活动
)

// 返水子类
const (
	SubCategoryServiceCharge                   TransactionSubCategory = 7000
	TransactionSubCategoryServiceChargeReceive                        = iota + SubCategoryServiceCharge // 返水领取
)

// 返佣子类
const (
	SubCategoryRebate                   TransactionSubCategory = 8000
	TransactionSubCategoryRebateSend                           = iota + SubCategoryRebate // 发放佣金
	TransactionSubCategoryRebateReceive                                                   // 领取佣金
)

// 利息宝子类
const (
	SubCategoryInterest                          TransactionSubCategory = 9000
	TransactionSubCategoryInterestProfit                                = iota + SubCategoryInterest // 利息宝收益
	TransactionSubCategoryHallToInterest                                                             // 大厅转入利息宝
	TransactionSubCategoryInterestManualPullback                                                     // 手动拉回-利息宝
	TransactionSubCategoryInterestToHall                                                             // 利息宝转到大厅
)

// 任务子类
const (
	SubCategoryTask                          TransactionSubCategory = 10000
	TransactionSubCategoryDailyTask                                 = iota + SubCategoryTask // 每日任务
	TransactionSubCategoryWeekLyTask                                                         // 每周任务
	TransactionSubCategoryNewcomerProfitTask                                                 // 新人福利
	TransactionSubCategoryAlivenessBox                                                       // 活跃度宝箱
	TransactionSubCategorySecretTask                                                         // 神秘任务
)

// VIP奖励子类
const (
	SubCategoryVipReward                   TransactionSubCategory = 11000
	TransactionSubCategoryVipMonthlyReward                        = iota + SubCategoryVipReward // VIP月奖金
	TransactionSubCategoryVipDailyReward                                                        // VIP日奖励
	TransactionSubCategoryVipWeeklyReward                                                       // VIP周奖金
	TransactionSubCategoryVipUpgradeReward                                                      // VIP晋级奖金
)

// 充值奖励子类
const (
	SubCategoryRechargeBonus            TransactionSubCategory = 12000
	TransactionSubCategoryRechargeBonus                        = iota + SubCategoryRechargeBonus // 优惠
)

// 俱乐部子类
const (
	SubCategoryClub                          TransactionSubCategory = 13000
	TransactionSubCategoryManualPullbackClub                        = iota + SubCategoryClub // 手动拉回-俱乐部
	TransactionSubCategoryHallToClub                                                         // 大厅转入俱乐部
	TransactionSubCategoryClubToHall                                                         // 俱乐部转到大厅
)

// 奖励子类
const (
	SubCategoryReward                  TransactionSubCategory = 14000
	TransactionSubCategoryGiveUpReward                        = iota + SubCategoryReward // 放弃奖励
	TransactionSubCategoryTransferOut                                                    // 转出
	TransactionSubCategoryTransferIn                                                     // 转入
)

// 担保理赔子类
const (
	SubCategoryGuaranteeClaim                 TransactionSubCategory = 15000
	TransactionSubCategoryClaimFreeze                                = iota + SubCategoryGuaranteeClaim // 理赔冻结
	TransactionSubCategoryClaimDefrost                                                                  // 理赔解冻
	TransactionSubCategoryClaimFee                                                                      // 理赔手续费
	TransactionSubCategoryClaimScoresIncrease                                                           // 理赔上分
	TransactionSubCategoryClaimScoresDecrease                                                           // 理赔下分
)

// 公积金子类
const (
	SubCategoryProvidentFund                   TransactionSubCategory = 16000
	TransactionSubCategoryProvidentFundReceive                        = iota + SubCategoryProvidentFund // 公积金领取
)

// 盲盒抽奖子类
const (
	SubCategoryBlindBoxLottery                  TransactionSubCategory = 17000
	TransactionSubCategoryBlindBoxLotteryDeduct                        = iota + SubCategoryBlindBoxLottery // 盲盒抽奖扣除
	TransactionSubCategoryBlindBoxLotteryReward                                                            // 盲盒抽奖奖励
)

// SubCategoryGameBetSettlement 游戏投注子类
const (
	SubCategoryGameBetSettlement            TransactionSubCategory = 18000
	TransactionSubCategoryGameBet                                  = iota + SubCategoryGameBetSettlement // 游戏投注
	TransactionSubCategoryGameBetSettlement                                                              // 游戏投注结算
	TransactionSubCategoryGameBetCancel                                                                  // 游戏撤单
)

func (t TransactionSubCategory) Int() int {
	return int(t)
}

func (t TransactionSubCategory) Int64() int64 {
	return int64(t)
}

func (t TransactionSubCategory) String() string {
	switch t {
	case TransactionSubCategoryWithdrawal:
		return "转出"
	case TransactionSubCategoryDeposits:
		return "转入"
	case TransactionSubCategoryUalaTransfer:
		return "Uala转账"
	case TransactionSubCategoryMercadoPagoTransfer:
		return "Mercado Pago转账"
	case TransactionSubCategoryUPDAYTransfer:
		return "UPDAY转账"
	case TransactionSubCategoryWithdrawFrozen:
		return "提现扣款"
	case TransactionSubCategoryWithdrawReject:
		return "提现失败"
	case TransactionSubCategoryWithdrawDefrost:
		return "提现解冻"
	case TransactionSubCategoryWithdrawSucceed:
		return "提现成功"
	case TransactionSubCategoryBankMerchantTransfer:
		return "转账给他人"
	case TransactionSubCategoryBankMerchantAddAmount:
		return "银商加款"
	case TransactionSubCategoryBankMerchantRecharge:
		return "银商充值"
	case TransactionSubCategoryBankMerchantGiveUser:
		return "银商赠送会员金额"
	case TransactionSubCategoryBankMerchantSubtractAmount:
		return "银商扣款"
	case TransactionSubCategoryManualAddAmount:
		return "人工加款"
	case TransactionSubCategoryManualAddOrder:
		return "手动补单"
	case TransactionSubCategoryBalanceRevise:
		return "修正负数余额"
	case TransactionSubCategoryManualAddRewardAmount:
		return "奖励手动加款"
	case TransactionSubCategoryManualSubtractAmount:
		return "手动扣款"
	case TransactionSubCategoryDeductAll:
		return "扣除全部资产"
	case TransactionSubCategorySurDeduct:
		return "追加扣除"
	case TransactionSubCategoryManualPullBack:
		return "人工拉回"
	case TransactionSubCategoryDeductExcessProfit:
		return "扣除超额盈利"
	case TransactionSubCategoryAgentActivity:
		return "代理活动"
	case TransactionSubCategoryLuckyBetActivity:
		return "幸运注单活动"
	case TransactionSubCategoryInvestActivity:
		return "投资活动"
	case TransactionSubCategoryNewcomerRewardActivity:
		return "新人彩金"
	case TransactionSubCategoryBenefitActivity:
		return "救援金活动"
	case TransactionSubCategoryPromoteActivity:
		return "推广活动"
	case TransactionSubCategoryFeedbackRewardActivity:
		return "有奖反馈活动"
	case TransactionSubCategoryRedPacketActivity:
		return "红包活动"
	case TransactionSubCategoryBetActivity:
		return "打码活动"
	case TransactionSubCategoryLotteryAssistanceActivity:
		return "抽奖助力"
	case TransactionSubCategoryRankActivity:
		return "排行榜活动"
	case TransactionSubCategoryCustomizeActivity:
		return "自定义活动"
	case TransactionSubCategoryBargainActivity:
		return "砍一刀"
	case TransactionSubCategorySpinActivity:
		return "转盘活动"
	case TransactionSubCategoryChannelRewardActivity:
		return "渠道奖励"
	case TransactionSubCategoryWordCollectionActivity:
		return "集字活动"
	case TransactionSubCategoryQuizActivity:
		return "竞猜活动"
	case TransactionSubCategoryServiceChargeReceive:
		return "返水领取"
	case TransactionSubCategoryRebateSend:
		return "发放佣金"
	case TransactionSubCategoryRebateReceive:
		return "领取佣金"
	case TransactionSubCategoryInterestProfit:
		return "利息宝收益"
	case TransactionSubCategoryHallToInterest:
		return "大厅转入利息宝"
	case TransactionSubCategoryInterestManualPullback:
		return "手动拉回-利息宝"
	case TransactionSubCategoryInterestToHall:
		return "利息宝转到大厅"
	case TransactionSubCategoryDailyTask:
		return "每日任务"
	case TransactionSubCategoryWeekLyTask:
		return "每周任务"
	case TransactionSubCategoryNewcomerProfitTask:
		return "新人福利"
	case TransactionSubCategoryAlivenessBox:
		return "活跃度宝箱"
	case TransactionSubCategorySecretTask:
		return "神秘任务"
	case TransactionSubCategoryVipMonthlyReward:
		return "VIP月奖金"
	case TransactionSubCategoryVipDailyReward:
		return "VIP日奖励"
	case TransactionSubCategoryVipWeeklyReward:
		return "VIP周奖金"
	case TransactionSubCategoryVipUpgradeReward:
		return "VIP晋级奖金"
	case TransactionSubCategoryRechargeBonus:
		return "优惠"
	case TransactionSubCategoryManualPullbackClub:
		return "手动拉回-俱乐部"
	case TransactionSubCategoryHallToClub:
		return "大厅转入俱乐部"
	case TransactionSubCategoryClubToHall:
		return "俱乐部转到大厅"
	case TransactionSubCategoryGiveUpReward:
		return "放弃奖励"
	case TransactionSubCategoryTransferOut:
		return "转出"
	case TransactionSubCategoryTransferIn:
		return "转入"
	case TransactionSubCategoryClaimFreeze:
		return "理赔冻结"
	case TransactionSubCategoryClaimDefrost:
		return "理赔解冻"
	case TransactionSubCategoryClaimFee:
		return "理赔手续费"
	case TransactionSubCategoryClaimScoresIncrease:
		return "理赔上分"
	case TransactionSubCategoryClaimScoresDecrease:
		return "理赔下分"
	case TransactionSubCategoryProvidentFundReceive:
		return "公积金领取"
	case TransactionSubCategoryBlindBoxLotteryDeduct:
		return "盲盒抽奖扣除"
	case TransactionSubCategoryBlindBoxLotteryReward:
		return "盲盒抽奖奖励"
	case TransactionSubCategoryGameBetSettlement:
		return "游戏投注结算"
	case TransactionSubCategoryRechargeOnline:
		return "在线充值"

	default:
		return ""
	}
}
