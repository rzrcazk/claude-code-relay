package model

type ApiLog struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Method     string `json:"method" gorm:"not null"`
	Path       string `json:"path" gorm:"not null"`
	StatusCode int    `json:"status_code"`
	UserID     uint   `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	RequestID  string `json:"request_id" gorm:"index"`
	Duration   int64  `json:"duration"` // 毫秒
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (a *ApiLog) TableName() string {
	return "api_logs"
}

func CreateApiLog(log *ApiLog) error {
	log.ID = 0
	return DB.Create(log).Error
}

func GetApiLogs(page, limit int) ([]ApiLog, int64, error) {
	var logs []ApiLog
	var total int64

	err := DB.Model(&ApiLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = DB.Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func GetApiLogsByUser(userID uint, page, limit int) ([]ApiLog, int64, error) {
	var logs []ApiLog
	var total int64

	query := DB.Model(&ApiLog{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
