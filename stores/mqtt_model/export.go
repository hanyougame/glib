package mqtt_model

// ExportResult 导出结果通知
type ExportResult struct {
	FileName   string `json:"file_name"`   // 文件名
	FailReason string `json:"fail_reason"` // 失败原因
	Status     int64  `json:"status"`      // 状态 1成功 2失败
}
