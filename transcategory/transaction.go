package transcategory

const (
	// zh 中文 en 英文
	ZhLanguage = "zh"
	EnLanguage = "en"
)

// TransactionCategory represents transaction main category types
type TransactionCategory int

var TransactionCategories = []TransactionCategory{
	TransactionCategoryCapitalSwitch,
	TransactionCategoryUserRecharge,
	TransactionCategoryUserWithdraw,
	//TransactionCategoryBankMerchantSettlement,
	TransactionCategoryCapitalRevise,
	TransactionCategoryActivity,
	//TransactionCategoryGoldReturn,
	TransactionCategoryRebate,
	//TransactionCategoryInterest,
	//TransactionCategoryTask,
	TransactionCategoryVipReward,
	TransactionCategoryRechargeBonus,
	//TransactionCategoryClub,
	//TransactionCategoryReward,
	//TransactionCategoryGuaranteeClaim,
	//TransactionCategoryProvidentFund,
	//TransactionCategoryBlindBoxLottery,
	TransactionCategoryGameBet,
	TransactionCategoryDataMigration,
}

const (
	_                                         TransactionCategory = iota
	TransactionCategoryCapitalSwitch                              // Fund Switching
	TransactionCategoryUserRecharge                               // Member deposit
	TransactionCategoryUserWithdraw                               // Member withdrawal
	TransactionCategoryBankMerchantSettlement                     // Merchant settlement
	TransactionCategoryCapitalRevise                              // Fund correction
	TransactionCategoryActivity                                   // Event
	TransactionCategoryGoldReturn                                 // Rebate
	TransactionCategoryRebate                                     // Commission
	TransactionCategoryInterest                                   // Interest
	TransactionCategoryTask                                       // Mission
	TransactionCategoryVipReward                                  // VIP reward
	TransactionCategoryRechargeBonus                              // Deposit promotion
	TransactionCategoryClub                                       // Club
	TransactionCategoryReward                                     // Reward
	TransactionCategoryGuaranteeClaim                             // Guaranteed claims
	TransactionCategoryProvidentFund                              // Provident Fund
	TransactionCategoryBlindBoxLottery                            // Blind box
	TransactionCategoryGameBet                                    // Game Bet
	TransactionCategoryDataMigration                              // Data Migration
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
		return "Fund Switching"
	case TransactionCategoryUserRecharge:
		return "Member Deposit"
	case TransactionCategoryUserWithdraw:
		return "Member Withdrawal"
	case TransactionCategoryBankMerchantSettlement:
		return "Merchant Settlement"
	case TransactionCategoryCapitalRevise:
		return "Fund Correction"
	case TransactionCategoryActivity:
		return "Event"
	case TransactionCategoryGoldReturn:
		return "Rebate"
	case TransactionCategoryRebate:
		return "Commission"
	case TransactionCategoryInterest:
		return "Interest"
	case TransactionCategoryTask:
		return "Mission"
	case TransactionCategoryVipReward:
		return "VIP Reward"
	case TransactionCategoryRechargeBonus:
		return "Deposit Promotion"
	case TransactionCategoryClub:
		return "Club"
	case TransactionCategoryReward:
		return "Reward"
	case TransactionCategoryGuaranteeClaim:
		return "Guaranteed Claims"
	case TransactionCategoryProvidentFund:
		return "Provident Fund"
	case TransactionCategoryBlindBoxLottery:
		return "Blind box"
	case TransactionCategoryGameBet:
		return "Game Bet"
	case TransactionCategoryDataMigration:
		return "Data Migration"
	default:
		return ""
	}
}

var TransactionCategoryZhMap = map[TransactionCategory]string{
	TransactionCategoryCapitalSwitch:          "资金切换",
	TransactionCategoryUserRecharge:           "会员充值",
	TransactionCategoryUserWithdraw:           "会员提现",
	TransactionCategoryBankMerchantSettlement: "银商结算",
	TransactionCategoryCapitalRevise:          "资金修正",
	TransactionCategoryActivity:               "活动",
	TransactionCategoryGoldReturn:             "返水",
	TransactionCategoryRebate:                 "返佣",
	TransactionCategoryInterest:               "利息宝",
	TransactionCategoryTask:                   "任务",
	TransactionCategoryVipReward:              "VIP奖励",
	TransactionCategoryRechargeBonus:          "充值优惠",
	TransactionCategoryClub:                   "俱乐部",
	TransactionCategoryReward:                 "奖励",
	TransactionCategoryGuaranteeClaim:         "担保理赔",
	TransactionCategoryProvidentFund:          "公积金",
	TransactionCategoryBlindBoxLottery:        "盲盒抽奖",
	TransactionCategoryGameBet:                "游戏投注",
	TransactionCategoryDataMigration:          "迁移",
}

func (t TransactionCategory) StringByLanguage(language string) string {
	switch language {
	case EnLanguage:
		// 英文内容
		return t.String()
	case ZhLanguage:
		// 获取中文内容
		return TransactionCategoryZhMap[t]
	default:
		// 默认英文内容
		return t.String()
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
			//TransactionSubCategoryUalaTransfer, TransactionSubCategoryMercadoPagoTransfer, TransactionSubCategoryUPDAYTransfer, TransactionSubCategoryRechargeOnline,
			TransactionSubCategoryRechargeOnline,
		}
	case TransactionCategoryUserWithdraw:
		return []TransactionSubCategory{
			TransactionSubCategoryWithdrawFrozen, TransactionSubCategoryWithdrawReject,
			//TransactionSubCategoryWithdrawDefrost, TransactionSubCategoryWithdrawSucceed,
		}
	//case TransactionCategoryBankMerchantSettlement:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryBankMerchantTransfer, TransactionSubCategoryBankMerchantAddAmount, TransactionSubCategoryBankMerchantRecharge, // 银商充值
	//		TransactionSubCategoryBankMerchantGiveUser, TransactionSubCategoryBankMerchantSubtractAmount,
	//	}
	case TransactionCategoryCapitalRevise:
		return []TransactionSubCategory{
			//TransactionSubCategoryManualAddAmount, TransactionSubCategoryManualAddOrder, TransactionSubCategoryBalanceRevise,
			//TransactionSubCategoryManualAddRewardAmount, TransactionSubCategoryManualSubtractAmount, TransactionSubCategoryDeductAll, // 扣除全部资产
			//TransactionSubCategorySurDeduct, TransactionSubCategoryManualPullBack, TransactionSubCategoryDeductExcessProfit,
			TransactionSubCategoryManualAddAmount, TransactionSubCategoryManualAddOrder,
			TransactionSubCategoryManualSubtractAmount, TransactionSubCategoryDeductAll,
			TransactionSubCategoryBloggerSalary,
		}
	case TransactionCategoryActivity:
		return []TransactionSubCategory{
			//TransactionSubCategoryAgentActivity, TransactionSubCategoryLuckyBetActivity, TransactionSubCategoryInvestActivity,
			//TransactionSubCategoryNewcomerRewardActivity, TransactionSubCategoryBenefitActivity, TransactionSubCategoryPromoteActivity,
			//TransactionSubCategoryFeedbackRewardActivity, TransactionSubCategoryRedPacketActivity, TransactionSubCategoryBetActivity,
			//TransactionSubCategoryLotteryAssistanceActivity, TransactionSubCategoryRankActivity, TransactionSubCategoryCustomizeActivity,
			//TransactionSubCategoryBargainActivity, TransactionSubCategorySpinActivity, TransactionSubCategoryChannelRewardActivity,
			//TransactionSubCategoryWordCollectionActivity, TransactionSubCategoryQuizActivity, TransactionSubCategoryRechargeActivity,
			TransactionSubCategoryAgentActivity, TransactionSubCategorySignInActivity, TransactionSubCategoryRechargeActivity, TransactionSubCategoryBetActivity, TransactionSubCategoryPromoteActivity, TransactionSubCategoryNewcomerProfitTask, TransactionSubCategoryBenefitActivity, TransactionSubCategoryRedPacketActivity, TransactionSubCategoryBargainActivity, TransactionSubCategorySpinActivity, TransactionSubCategoryLoginActivity, TransactionSubCategoryReturnUserActivity,
		}
	//case TransactionCategoryGoldReturn:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryServiceChargeReceive,
	//	}
	case TransactionCategoryRebate:
		return []TransactionSubCategory{
			//TransactionSubCategoryRebateSend,
			TransactionSubCategoryRebateReceive,
		}
	//case TransactionCategoryInterest:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryInterestProfit, TransactionSubCategoryHallToInterest,
	//		TransactionSubCategoryInterestManualPullback, TransactionSubCategoryInterestToHall,
	//	}
	//case TransactionCategoryTask:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryDailyTask, TransactionSubCategoryWeekLyTask, TransactionSubCategoryNewcomerProfitTask,
	//		TransactionSubCategoryAlivenessBox, TransactionSubCategorySecretTask,
	//	}
	case TransactionCategoryVipReward:
		return []TransactionSubCategory{
			TransactionSubCategoryVipMonthlyReward, TransactionSubCategoryVipDailyReward,
			TransactionSubCategoryVipWeeklyReward, TransactionSubCategoryVipUpgradeReward,
		}
	case TransactionCategoryRechargeBonus:
		return []TransactionSubCategory{
			TransactionSubCategoryRechargeBonus,
		}
	//case TransactionCategoryClub:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryManualPullbackClub, TransactionSubCategoryHallToClub, TransactionSubCategoryClubToHall,
	//	}
	//case TransactionCategoryReward:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryGiveUpReward, TransactionSubCategoryTransferOut, TransactionSubCategoryTransferIn,
	//	}
	//case TransactionCategoryGuaranteeClaim:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryClaimFreeze, TransactionSubCategoryClaimDefrost, TransactionSubCategoryClaimFee,
	//		TransactionSubCategoryClaimScoresIncrease, TransactionSubCategoryClaimScoresDecrease,
	//	}
	//case TransactionCategoryProvidentFund:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryRechargeBonus,
	//	}
	//case TransactionCategoryBlindBoxLottery:
	//	return []TransactionSubCategory{
	//		TransactionSubCategoryBlindBoxLotteryDeduct, TransactionSubCategoryBlindBoxLotteryReward,
	//	}
	case TransactionCategoryGameBet:
		return []TransactionSubCategory{
			TransactionSubCategoryGameBet,
			TransactionSubCategoryGameBetSettlement,
			TransactionSubCategoryGameBetCancel,
			TransactionSubCategoryGameBetRollback,
		}
	default:
		return []TransactionSubCategory{}
	}
}

// TransactionSubCategory represents transaction sub-category types
type TransactionSubCategory int

// Capital Switch Subcategories
const (
	SubCategoryCapitalSwitch         TransactionSubCategory = 1000
	TransactionSubCategoryWithdrawal                        = iota + SubCategoryCapitalSwitch // Transfer Out
	TransactionSubCategoryDeposits                                                            // Transfer In
)

// User Recharge Subcategories
const (
	SubCategoryUserRecharge                   TransactionSubCategory = 2000
	TransactionSubCategoryUalaTransfer                               = iota + SubCategoryUserRecharge // Uala Transfer
	TransactionSubCategoryMercadoPagoTransfer                                                         // Mercado Pago Transfer
	TransactionSubCategoryUPDAYTransfer                                                               // UPDAY Transfer
	TransactionSubCategoryRechargeOnline                                                              // Online Recharge
)

// User Withdrawal Subcategories
const (
	SubCategoryUserWithdraw               TransactionSubCategory = 3000
	TransactionSubCategoryWithdrawFrozen                         = iota + SubCategoryUserWithdraw // Withdrawal Freeze
	TransactionSubCategoryWithdrawReject                                                          // Withdrawal Rejection
	TransactionSubCategoryWithdrawDefrost                                                         // Withdrawal Unfreeze
	TransactionSubCategoryWithdrawSucceed                                                         // Withdrawal Success
)

// Bank Merchant Settlement Subcategories
const (
	SubCategoryBankMerchantSettlement                TransactionSubCategory = 4000
	TransactionSubCategoryBankMerchantTransfer                              = iota + SubCategoryBankMerchantSettlement // Transfer to Others
	TransactionSubCategoryBankMerchantAddAmount                                                                        // Merchant Addition
	TransactionSubCategoryBankMerchantRecharge                                                                         // Merchant Recharge
	TransactionSubCategoryBankMerchantGiveUser                                                                         // Merchant User Gift
	TransactionSubCategoryBankMerchantSubtractAmount                                                                   // Merchant Deduction
)

// Capital Revision Subcategories
const (
	SubCategoryCapitalRevise                    TransactionSubCategory = 5000
	TransactionSubCategoryManualAddAmount                              = iota + SubCategoryCapitalRevise // Manual Addition
	TransactionSubCategoryManualAddOrder                                                                 // Manual Order Adjustment
	TransactionSubCategoryBalanceRevise                                                                  // Negative Balance Correction
	TransactionSubCategoryManualAddRewardAmount                                                          // Reward Manual Addition
	TransactionSubCategoryManualSubtractAmount                                                           // Manual Deduction
	TransactionSubCategoryDeductAll                                                                      // Full Asset Deduction
	TransactionSubCategorySurDeduct                                                                      // Additional Deduction
	TransactionSubCategoryManualPullBack                                                                 // Manual Recall
	TransactionSubCategoryDeductExcessProfit                                                             // Excess Profit Deduction
	TransactionSubCategoryBloggerSalary                                                                  // Blogger Salary
)

// Activity Subcategories
const (
	SubCategoryActivity                             TransactionSubCategory = 6000
	TransactionSubCategoryAgentActivity                                    = iota + SubCategoryActivity // Agent Activity
	TransactionSubCategoryLuckyBetActivity                                                              // Lucky Bet Activity
	TransactionSubCategoryInvestActivity                                                                // Investment Activity
	TransactionSubCategoryNewcomerRewardActivity                                                        // Newcomer Bonus
	TransactionSubCategoryBenefitActivity                                                               // Relief Fund Activity
	TransactionSubCategoryPromoteActivity                                                               // Promotion Activity
	TransactionSubCategoryFeedbackRewardActivity                                                        // Feedback Reward Activity
	TransactionSubCategoryRedPacketActivity                                                             // Red Packet Activity
	TransactionSubCategoryBetActivity                                                                   // Betting Activity
	TransactionSubCategoryLotteryAssistanceActivity                                                     // Lottery Assistance
	TransactionSubCategoryRankActivity                                                                  // Ranking Activity
	TransactionSubCategoryCustomizeActivity                                                             // Custom Activity
	TransactionSubCategoryBargainActivity                                                               // Bargain Activity
	TransactionSubCategorySpinActivity                                                                  // Spin Activity
	TransactionSubCategoryChannelRewardActivity                                                         // Channel Reward
	TransactionSubCategoryWordCollectionActivity                                                        // Word Collection
	TransactionSubCategoryQuizActivity                                                                  // Quiz
	TransactionSubCategoryRechargeActivity                                                              // 充值活动
	TransactionSubCategorySignInActivity                                                                // 签到活动// Activity
	TransactionSubCategoryLoginActivity                                                                 // Login Activity 登录活动
	TransactionSubCategoryReturnUserActivity                                                            // return user 用户召回活动
)

// Gold Return Subcategories
const (
	SubCategoryServiceCharge                   TransactionSubCategory = 7000
	TransactionSubCategoryServiceChargeReceive                        = iota + SubCategoryServiceCharge // Service Charge Receipt
)

// Rebate Subcategories
const (
	SubCategoryRebate                   TransactionSubCategory = 8000
	TransactionSubCategoryRebateSend                           = iota + SubCategoryRebate // Rebate Distribution
	TransactionSubCategoryRebateReceive                                                   // Rebate Receipt
)

// Interest Subcategories
const (
	SubCategoryInterest                          TransactionSubCategory = 9000
	TransactionSubCategoryInterestProfit                                = iota + SubCategoryInterest // Interest Earnings
	TransactionSubCategoryHallToInterest                                                             // Hall to Interest Transfer
	TransactionSubCategoryInterestManualPullback                                                     // Manual Recall - Interest
	TransactionSubCategoryInterestToHall                                                             // Interest to Hall Transfer
)

// Task Subcategories
const (
	SubCategoryTask                          TransactionSubCategory = 10000
	TransactionSubCategoryDailyTask                                 = iota + SubCategoryTask // Daily Task
	TransactionSubCategoryWeekLyTask                                                         // Weekly Task
	TransactionSubCategoryNewcomerProfitTask                                                 // Newcomer Benefits
	TransactionSubCategoryAlivenessBox                                                       // Activity Box
	TransactionSubCategorySecretTask                                                         // Secret Task
)

// VIP Reward Subcategories
const (
	SubCategoryVipReward                   TransactionSubCategory = 11000
	TransactionSubCategoryVipMonthlyReward                        = iota + SubCategoryVipReward // VIP Monthly Reward
	TransactionSubCategoryVipDailyReward                                                        // VIP Daily Reward
	TransactionSubCategoryVipWeeklyReward                                                       // VIP Weekly Reward
	TransactionSubCategoryVipUpgradeReward                                                      // VIP Upgrade Reward
)

// Recharge Bonus Subcategories
const (
	SubCategoryRechargeBonus            TransactionSubCategory = 12000
	TransactionSubCategoryRechargeBonus                        = iota + SubCategoryRechargeBonus // Bonus
)

// Club Subcategories
const (
	SubCategoryClub                          TransactionSubCategory = 13000
	TransactionSubCategoryManualPullbackClub                        = iota + SubCategoryClub // Manual Recall - Club
	TransactionSubCategoryHallToClub                                                         // Hall to Club Transfer
	TransactionSubCategoryClubToHall                                                         // Club to Hall Transfer
)

// Reward Subcategories
const (
	SubCategoryReward                  TransactionSubCategory = 14000
	TransactionSubCategoryGiveUpReward                        = iota + SubCategoryReward // Reward Waiver
	TransactionSubCategoryTransferOut                                                    // Transfer Out
	TransactionSubCategoryTransferIn                                                     // Transfer In
)

// Guarantee Claim Subcategories
const (
	SubCategoryGuaranteeClaim                 TransactionSubCategory = 15000
	TransactionSubCategoryClaimFreeze                                = iota + SubCategoryGuaranteeClaim // Claim Freeze
	TransactionSubCategoryClaimDefrost                                                                  // Claim Unfreeze
	TransactionSubCategoryClaimFee                                                                      // Claim Fee
	TransactionSubCategoryClaimScoresIncrease                                                           // Score Increase
	TransactionSubCategoryClaimScoresDecrease                                                           // Score Decrease
)

// Provident Fund Subcategories
const (
	SubCategoryProvidentFund                   TransactionSubCategory = 16000
	TransactionSubCategoryProvidentFundReceive                        = iota + SubCategoryProvidentFund // Provident Fund Receipt
)

// Blind Box Lottery Subcategories
const (
	SubCategoryBlindBoxLottery                  TransactionSubCategory = 17000
	TransactionSubCategoryBlindBoxLotteryDeduct                        = iota + SubCategoryBlindBoxLottery // Blind Box Deduction
	TransactionSubCategoryBlindBoxLotteryReward                                                            // Blind Box Reward
)

// Game Bet Subcategories
const (
	SubCategoryGameBetSettlement            TransactionSubCategory = 18000
	TransactionSubCategoryGameBet                                  = iota + SubCategoryGameBetSettlement // Game Bet
	TransactionSubCategoryGameBetSettlement                                                              // Bet Settlement
	TransactionSubCategoryGameBetCancel                                                                  // Bet Cancellation
	TransactionSubCategoryGameBetRollback                                                                // Bet Rollback
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
		return "Transfer Out"
	case TransactionSubCategoryDeposits:
		return "Transfer In"
	case TransactionSubCategoryUalaTransfer:
		return "Uala Transfer"
	case TransactionSubCategoryMercadoPagoTransfer:
		return "Mercado Pago Transfer"
	case TransactionSubCategoryUPDAYTransfer:
		return "UPDAY Transfer"
	case TransactionSubCategoryWithdrawFrozen:
		return "Withdrawals&Deductions"
	case TransactionSubCategoryWithdrawReject:
		return "Withdrawal Rejection"
	case TransactionSubCategoryWithdrawDefrost:
		return "Withdrawal Fail"
	case TransactionSubCategoryWithdrawSucceed:
		return "Withdrawal Success"
	case TransactionSubCategoryBankMerchantTransfer:
		return "Transfer To Others"
	case TransactionSubCategoryBankMerchantAddAmount:
		return "Merchant Addition"
	case TransactionSubCategoryBankMerchantRecharge:
		return "Merchant Recharge"
	case TransactionSubCategoryBankMerchantGiveUser:
		return "Merchant User Gift"
	case TransactionSubCategoryBankMerchantSubtractAmount:
		return "Merchant Deduction"
	case TransactionSubCategoryManualAddAmount:
		return "Manual Addition"
	case TransactionSubCategoryManualAddOrder:
		return "Manual Order Filling"
	case TransactionSubCategoryBalanceRevise:
		return "Negative Balance Correction"
	case TransactionSubCategoryManualAddRewardAmount:
		return "Reward Manual Addition"
	case TransactionSubCategoryManualSubtractAmount:
		return "Manual Deduction"
	case TransactionSubCategoryDeductAll:
		return "Deduct All Assets"
	case TransactionSubCategorySurDeduct:
		return "Additional Deduction"
	case TransactionSubCategoryManualPullBack:
		return "Manual Recall"
	case TransactionSubCategoryDeductExcessProfit:
		return "Excess Profit Deduction"
	case TransactionSubCategoryAgentActivity:
		return "Agent Event"
	case TransactionSubCategoryLuckyBetActivity:
		return "Lucky Bet Activity"
	case TransactionSubCategoryInvestActivity:
		return "Investment Activity"
	case TransactionSubCategoryNewcomerRewardActivity:
		return "Newcomer Bonus"
	case TransactionSubCategoryBenefitActivity:
		return "Daily Loss Recovery"
	case TransactionSubCategoryPromoteActivity:
		return "Promotion Event"
	case TransactionSubCategoryFeedbackRewardActivity:
		return "Feedback Reward Activity"
	case TransactionSubCategoryRedPacketActivity:
		return "Red Envelope Event"
	case TransactionSubCategoryBetActivity:
		return "Bet Event"
	case TransactionSubCategoryLotteryAssistanceActivity:
		return "Lottery Assistance"
	case TransactionSubCategoryRankActivity:
		return "Ranking Activity"
	case TransactionSubCategoryCustomizeActivity:
		return "Custom Activity"
	case TransactionSubCategoryBargainActivity:
		return "Lucky draw"
	case TransactionSubCategorySpinActivity:
		return "Fortune Wheel Event"
	case TransactionSubCategoryChannelRewardActivity:
		return "Channel Reward"
	case TransactionSubCategoryWordCollectionActivity:
		return "Word Collection"
	case TransactionSubCategoryQuizActivity:
		return "Quiz Activity"
	case TransactionSubCategoryServiceChargeReceive:
		return "Service Charge Receipt"
	case TransactionSubCategoryRebateSend:
		return "Commission Distribution"
	case TransactionSubCategoryRebateReceive:
		return "Claim Commission"
	case TransactionSubCategoryInterestProfit:
		return "Interest Earnings"
	case TransactionSubCategoryHallToInterest:
		return "Hall to Interest Transfer"
	case TransactionSubCategoryInterestManualPullback:
		return "Manual Recall - Interest"
	case TransactionSubCategoryInterestToHall:
		return "Interest to Hall Transfer"
	case TransactionSubCategoryDailyTask:
		return "Daily Task"
	case TransactionSubCategoryWeekLyTask:
		return "Weekly Task"
	case TransactionSubCategoryNewcomerProfitTask:
		return "New Player Bonus"
	case TransactionSubCategoryAlivenessBox:
		return "Activity Box"
	case TransactionSubCategorySecretTask:
		return "Secret Task"
	case TransactionSubCategoryVipMonthlyReward:
		return "VIP Monthly Bonus"
	case TransactionSubCategoryVipDailyReward:
		return "VIP Daily Bonus"
	case TransactionSubCategoryVipWeeklyReward:
		return "VIP Weekly Bonus"
	case TransactionSubCategoryVipUpgradeReward:
		return "VIP Promotion Bonus"
	case TransactionSubCategoryRechargeBonus:
		return "Offers"
	case TransactionSubCategoryManualPullbackClub:
		return "Manual Recall - Club"
	case TransactionSubCategoryHallToClub:
		return "Hall to Club Transfer"
	case TransactionSubCategoryClubToHall:
		return "Club to Hall Transfer"
	case TransactionSubCategoryGiveUpReward:
		return "Reward Waiver"
	case TransactionSubCategoryTransferOut:
		return "Transfer Out"
	case TransactionSubCategoryTransferIn:
		return "Transfer In"
	case TransactionSubCategoryClaimFreeze:
		return "Claim Freeze"
	case TransactionSubCategoryClaimDefrost:
		return "Claim Unfreeze"
	case TransactionSubCategoryClaimFee:
		return "Claim Fee"
	case TransactionSubCategoryClaimScoresIncrease:
		return "Score Increase"
	case TransactionSubCategoryClaimScoresDecrease:
		return "Score Decrease"
	case TransactionSubCategoryProvidentFundReceive:
		return "Provident Fund Receipt"
	case TransactionSubCategoryBlindBoxLotteryDeduct:
		return "Blind Box Deduction"
	case TransactionSubCategoryBlindBoxLotteryReward:
		return "Blind Box Reward"
	case TransactionSubCategoryGameBetSettlement:
		return "Bet Settlement"
	case TransactionSubCategoryRechargeOnline:
		return "Online Deposit"
	case TransactionSubCategoryGameBet:
		return "Game Bet"
	case TransactionSubCategoryGameBetCancel:
		return "Game Bet Cancel"
	case TransactionSubCategoryGameBetRollback:
		return "Game Bet Roll-Back"
	case TransactionSubCategoryRechargeActivity:
		return "Deposit Event" // 充值活动
	case TransactionSubCategorySignInActivity:
		return "Check-In Event" // 签到活动
	case TransactionSubCategoryLoginActivity:
		return "Login Activity" // 登录活动
	case TransactionSubCategoryBloggerSalary:
		return "Blogger Salary" // 博主工资
	case TransactionSubCategoryReturnUserActivity:
		return "Returning User Event" // 召回惠东
	default:
		return ""
	}
}

var TransactionSubCategoryMap = map[TransactionSubCategory]string{
	TransactionSubCategoryWithdrawal:                 "转入",
	TransactionSubCategoryDeposits:                   "转出",
	TransactionSubCategoryUalaTransfer:               "Uala转账",
	TransactionSubCategoryMercadoPagoTransfer:        "MercadoPago转账",
	TransactionSubCategoryUPDAYTransfer:              "UPDAY转账",
	TransactionSubCategoryWithdrawFrozen:             "提现扣款",
	TransactionSubCategoryWithdrawReject:             "提现失败",
	TransactionSubCategoryWithdrawDefrost:            "提现失败",
	TransactionSubCategoryWithdrawSucceed:            "提现成功",
	TransactionSubCategoryBankMerchantTransfer:       "转账给他人",
	TransactionSubCategoryBankMerchantAddAmount:      "银商加款",
	TransactionSubCategoryBankMerchantRecharge:       "银商充值",
	TransactionSubCategoryBankMerchantGiveUser:       "银商赠送会员金额",
	TransactionSubCategoryBankMerchantSubtractAmount: "银商扣款",
	TransactionSubCategoryManualAddAmount:            "人工加款",
	TransactionSubCategoryManualAddOrder:             "手动补单",
	TransactionSubCategoryBalanceRevise:              "修正负数余额",
	TransactionSubCategoryManualAddRewardAmount:      "人工加款",
	TransactionSubCategoryManualSubtractAmount:       "手动扣款",
	TransactionSubCategoryDeductAll:                  "扣除全部资产",
	TransactionSubCategorySurDeduct:                  "追缴扣除",
	TransactionSubCategoryManualPullBack:             "返还资金",
	TransactionSubCategoryDeductExcessProfit:         "扣除超额盈利",
	TransactionSubCategoryAgentActivity:              "代理活动",
	TransactionSubCategoryLuckyBetActivity:           "幸运注单活动",
	TransactionSubCategoryInvestActivity:             "投资活动",
	TransactionSubCategoryNewcomerRewardActivity:     "新人彩金",
	TransactionSubCategoryBenefitActivity:            "救援金活动",
	TransactionSubCategoryPromoteActivity:            "推广活动",
	TransactionSubCategoryFeedbackRewardActivity:     "有奖反馈活动",
	TransactionSubCategoryRedPacketActivity:          "红包活动",
	TransactionSubCategoryBetActivity:                "打码活动",
	TransactionSubCategoryLotteryAssistanceActivity:  "抽奖助力活动",
	TransactionSubCategoryRankActivity:               "排行榜活动",
	TransactionSubCategoryCustomizeActivity:          "自定义活动",
	TransactionSubCategoryBargainActivity:            "砍一刀",
	TransactionSubCategorySpinActivity:               "转盘活动",
	TransactionSubCategoryChannelRewardActivity:      "渠道奖励",
	TransactionSubCategoryWordCollectionActivity:     "集字活动",
	TransactionSubCategoryQuizActivity:               "竞猜活动",
	TransactionSubCategoryServiceChargeReceive:       "返水领取",
	TransactionSubCategoryRebateSend:                 "发放佣金",
	TransactionSubCategoryRebateReceive:              "领取佣金",
	TransactionSubCategoryInterestProfit:             "利息宝收益",
	TransactionSubCategoryHallToInterest:             "大厅转入利息宝",
	TransactionSubCategoryInterestManualPullback:     "手动拉回-利息宝",
	TransactionSubCategoryInterestToHall:             "利息宝转到大厅",
	TransactionSubCategoryDailyTask:                  "每日任务",
	TransactionSubCategoryWeekLyTask:                 "每周任务",
	TransactionSubCategoryNewcomerProfitTask:         "新人彩金",
	TransactionSubCategoryAlivenessBox:               "活跃度宝箱",
	TransactionSubCategorySecretTask:                 "神秘任务",
	TransactionSubCategoryVipMonthlyReward:           "VIP月奖金",
	TransactionSubCategoryVipDailyReward:             "VIP日奖励",
	TransactionSubCategoryVipWeeklyReward:            "VIP周奖金",
	TransactionSubCategoryVipUpgradeReward:           "VIP晋级奖金",
	TransactionSubCategoryRechargeBonus:              "优惠",
	TransactionSubCategoryManualPullbackClub:         "手动拉回-俱乐部",
	TransactionSubCategoryHallToClub:                 "大厅转入俱乐部",
	TransactionSubCategoryClubToHall:                 "俱乐部转到大厅",
	TransactionSubCategoryGiveUpReward:               "放弃奖励",
	TransactionSubCategoryTransferOut:                "转入",
	TransactionSubCategoryTransferIn:                 "转出",
	TransactionSubCategoryClaimFreeze:                "理赔冻结",
	TransactionSubCategoryClaimDefrost:               "理赔解冻",
	TransactionSubCategoryClaimFee:                   "理赔手续费",
	TransactionSubCategoryClaimScoresIncrease:        "理赔上分",
	TransactionSubCategoryClaimScoresDecrease:        "理赔下分",
	TransactionSubCategoryProvidentFundReceive:       "公积金领取",
	TransactionSubCategoryBlindBoxLotteryDeduct:      "盲盒抽奖扣除",
	TransactionSubCategoryBlindBoxLotteryReward:      "盲盒抽奖奖励",
	TransactionSubCategoryGameBetSettlement:          "投注派奖",
	TransactionSubCategoryRechargeOnline:             "在线充值",
	TransactionSubCategoryGameBet:                    "游戏投注",
	TransactionSubCategoryGameBetCancel:              "投注取消",
	TransactionSubCategoryGameBetRollback:            "投注回滚",
	TransactionSubCategoryRechargeActivity:           "充值活动", // 充值活动
	TransactionSubCategorySignInActivity:             "签到活动", // 签到活动
	TransactionSubCategoryLoginActivity:              "登录活动", // 登录活动
	TransactionSubCategoryBloggerSalary:              "博主工资",
	TransactionSubCategoryReturnUserActivity:         "用户召回活动",
}

func (t TransactionSubCategory) StringByLanguage(language string) string {
	switch language {
	case EnLanguage:
		// 英文内容
		return t.String()
	case ZhLanguage:
		// 中文内容获取
		return TransactionSubCategoryMap[t]
	default:
		// 默认英文内容
		return t.String()
	}
}
