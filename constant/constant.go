package constant

const (
	// 用户状态
	UserStatusActive   = 1
	UserStatusInactive = 0

	// 用户角色
	RoleAdmin = "admin"
	RoleUser  = "user"

	// 任务状态
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"

	// 任务优先级
	TaskPriorityLow    = 1
	TaskPriorityMedium = 2
	TaskPriorityHigh   = 3

	// 缓存键前缀
	CacheKeyUser = "user:"
	CacheKeyTask = "task:"

	// API响应码
	Success              = 20000
	InvalidParams        = 40000
	Unauthorized         = 40001
	UserStatusAbnormal   = 40002
	InsufficientPrivileges = 40003
	NotFound             = 40004
	TooManyRequests      = 42901
	InternalServerError  = 50000
)