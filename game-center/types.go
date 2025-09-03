package game_center

type (
	CommonResp[T any] struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    T      `json:"data"`
		TraceID string `json:"trace_id"`
	}

	RegisterAccountReq struct {
		UserID   string `json:"user_id,optional" validate:"required,max=32,alphaNumUnderscore"`
		Nickname string `json:"nickname,optional" validate:"max=32"`
		Avatar   string `json:"avatar,optional" validate:"max=512"`
	}

	RegisterAccountResp struct {
		OpenID string `json:"openid"`
	}

	GetGameLinkReq struct {
		GameID   int64  `json:"game_id,optional" validate:"required"`
		OpenID   int64  `json:"openid,string,optional" validate:"required"`
		Balance  string `json:"balance,optional" validate:"balance"`
		TaskID   string `json:"task_id,optional" validate:"required,max=64"` // 转入订单号
		Language string `json:"language,optional,default=en"`
		Extend   Extend `json:"extend,optional"`
	}
	Extend struct {
		DeviceOS string `json:"device_os,optional" validate:"max=32"`
		DeviceID string `json:"device_id,optional" validate:"max=64"`
		ClientIP string `json:"client_ip,optional" validate:"max=45"`
	}
	GetGameLinkResp struct {
		Url            string `json:"url"`             // 游戏地址
		TaskID         string `json:"task_id"`         // 商户转入订单号
		OrderNo        string `json:"order_no"`        // 订单号
		TransferStatus int    `json:"transfer_status"` // 转账状态 0 处理中 1 成功 2 失败 3 未处理
		OpenMethod     int    `json:"open_method"`     // 打开方式 1 iframe 2 新标签页
	}

	GetGameListReq struct {
		PlatformID int64 `form:"platform_id,optional"`
	}

	GetGameListResp struct {
		List []GameInfo `json:"list"`
	}
	GameInfo struct {
		GameID          int64    `json:"game_id"`
		GameName        string   `json:"game_name"`
		Icon            string   `json:"icon"`
		PlatformID      int64    `json:"platform_id"`
		Countries       []string `json:"countries"`
		Languages       []string `json:"languages"`
		Currencies      []string `json:"currencies"`
		CategoryCode    string   `json:"category_code"`
		ThumbCover      string   `json:"thumb_cover" gorm:"thumb_cover"`           // 缩略图
		HorizontalCover string   `json:"horizontal_cover" gorm:"horizontal_cover"` // 垂直封面
		VerticalCover   string   `json:"vertical_cover" gorm:"vertical_cover"`     // 水平封面
		ScreenType      int      `json:"screen_type" gorm:"screen_type"`           // 横竖屏 1都支持 2 竖屏 3横屏
		Status          int      `json:"status" gorm:"status"`                     // 状态 1启用 2禁用
		IsTrialPlay     int      `json:"is_trial_play" gorm:"is_trial_play"`       // 是否允许试玩 1 是 2否
	}

	GetBetListReq struct {
		StartTime      int64  `form:"start_time,optional" validate:"required"`
		EndTime        int64  `form:"end_time,optional" validate:"required,gtfield=StartTime"`
		MerchantUserID string `form:"merchant_user_id,optional"`
		OpenID         int64  `form:"openid,string,optional"`
		Page           int    `form:"page,optional" validate:"required"`
		PageSize       int    `form:"page_size,default=1000" validate:"max=2000"`
	}

	GetBetListResp struct {
		Rows []BetInfo `json:"rows"`
	}

	BetInfo struct {
		SummaryID          int64  `json:"summary_id"`                    // 汇总id
		CurrencyCode       string `json:"currency_code"`                 // 币种
		RoundID            string `json:"round_id"`                      // 牌局id
		OpenID             string `json:"open_id" gorm:"column:user_id"` // open id
		MerchantID         int64  `json:"merchant_id"`                   // 商户id
		MerchantUserID     string `json:"merchant_user_id"`              // 商户用户id
		GameID             int64  `json:"game_id"`                       // 游戏id
		PlatformID         int64  `json:"platform_id"`                   // 厂商id
		BetCount           int64  `json:"bet_count"`                     // 投注笔数
		BetAmount          string `json:"bet_amount"`                    // 投注金额
		CancelBetCount     int64  `json:"cancel_bet_count"`              // 取消投注次数
		CancelBetAmount    string `json:"cancel_bet_amount"`             // 取消投注金额
		SettleCount        int64  `json:"settle_count"`                  // 结算次数
		SettleAmount       string `json:"settle_amount"`                 // 结算金额
		CancelSettleCount  int64  `json:"cancel_settle_count"`           // 取消结算次数
		CancelSettleAmount string `json:"cancel_settle_amount"`          // 取消结算金额
		UpdatedAt          int64  `json:"updated_at"`                    // 更新时间
		CreatedAt          int64  `json:"created_at"`                    // 创建时间
		BetAt              int64  `json:"bet_at"`                        // 下注时间
		SettledAt          int64  `json:"settled_at"`                    // 结算时间
	}

	GetGamePlatformListResp struct {
		List []GamePlatformInfo `json:"list"`
	}

	GamePlatformInfo struct {
		PlatformID   int64    `json:"platform_id"`
		PlatformName string   `json:"platform_name"`
		Cover        string   `json:"cover"`
		Countries    []string `json:"countries"`
		Languages    []string `json:"languages"`
		Status       int      `json:"status"` // 状态 1 启用 2 禁用
	}

	CheckTransferStatusReq struct {
		OpenID  int64    `json:"openid,string,optional" validate:"required"`
		TaskIds []string `json:"task_ids,optional" validate:"required,max=100"` // 下游订单号
	}

	CheckTransferStatusResp struct {
		Orders []TransferStatus `json:"orders"`
	}

	TransferStatus struct {
		TaskID       string `json:"task_id"`       // 商户转入订单号
		OrderNo      string `json:"order_no"`      // 中台订单号
		Status       int    `json:"status"`        // 状态 1待处理 2 处理中 3 成功 4 失败
		TransferType int    `json:"transfer_type"` // 转账类型 1 转入 2 转出
		Amount       string `json:"amount"`        // 金额
		CreatedAt    int64  `json:"created_at"`    // 创建时间
		UpdatedAt    int64  `json:"updated_at"`    // 更新时间
	}

	TransferOutReq struct {
		OpenID     int64  `json:"openid,string,optional" validate:"required"`  // 用户id
		TaskID     string `json:"task_id,optional" validate:"required,max=64"` // 商户转入订单号
		PlatformID int64  `json:"platform_id,optional" validate:"required"`    // 平台ID
	}

	TransferOutResp struct {
		OrderNo        string `json:"order_no"`
		TaskID         string `json:"task_id"`         // 商户转入订单号
		Amount         string `json:"amount"`          // 金额
		TransferStatus int    `json:"transfer_status"` // 转账状态 0 处理中 1 成功 2失败
	}
	GetBalanceReq struct {
		OpenID      int64   `json:"openid,string,optional" validate:"required"`
		PlatformIds []int64 `json:"platform_ids,optional" validate:"required"`
	}

	GetBalanceResp struct {
		List []GetBalanceInfo `json:"list"`
	}

	GetBalanceInfo struct {
		PlatformID int64  `json:"platform_id"`
		Balance    string `json:"balance"`
	}
)
