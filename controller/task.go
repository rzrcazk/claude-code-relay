package controller

import (
	"claude-scheduler/constant"
	"claude-scheduler/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	ScheduleAt  string `json:"schedule_at"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	ScheduleAt  string `json:"schedule_at"`
}

func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	// 解析调度时间
	var scheduleAt time.Time
	if req.ScheduleAt != "" {
		var err error
		scheduleAt, err = time.Parse("2006-01-02 15:04:05", req.ScheduleAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "时间格式错误，请使用 YYYY-MM-DD HH:MM:SS 格式",
				"code":  constant.InvalidParams,
			})
			return
		}
	} else {
		scheduleAt = time.Now()
	}

	// 设置默认优先级
	if req.Priority == 0 {
		req.Priority = constant.TaskPriorityLow
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		UserID:      userID,
		ScheduleAt:  scheduleAt,
		Status:      constant.TaskStatusPending,
	}

	if err := model.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建任务失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"code":    constant.Success,
		"data": gin.H{
			"task": task,
		},
	})
}

func GetTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 普通用户只能查看自己的任务
	var userID uint
	user := c.MustGet("user").(*model.User)
	if user.Role != constant.RoleAdmin {
		userID = user.ID
	}

	tasks, total, err := model.GetTasks(page, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取任务列表失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data": gin.H{
			"tasks": tasks,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "任务ID格式错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	task, err := model.GetTaskById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "任务不存在",
			"code":  constant.NotFound,
		})
		return
	}

	// 检查权限：普通用户只能查看自己的任务
	user := c.MustGet("user").(*model.User)
	if user.Role != constant.RoleAdmin && task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  constant.InsufficientPrivileges,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data": gin.H{
			"task": task,
		},
	})
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "任务ID格式错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	task, err := model.GetTaskById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "任务不存在",
			"code":  constant.NotFound,
		})
		return
	}

	// 检查权限：普通用户只能修改自己的任务
	user := c.MustGet("user").(*model.User)
	if user.Role != constant.RoleAdmin && task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  constant.InsufficientPrivileges,
		})
		return
	}

	// 更新字段
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
		if req.Status == constant.TaskStatusCompleted {
			now := time.Now()
			task.CompletedAt = &now
		}
	}
	if req.Priority > 0 {
		task.Priority = req.Priority
	}
	if req.ScheduleAt != "" {
		scheduleAt, err := time.Parse("2006-01-02 15:04:05", req.ScheduleAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "时间格式错误，请使用 YYYY-MM-DD HH:MM:SS 格式",
				"code":  constant.InvalidParams,
			})
			return
		}
		task.ScheduleAt = scheduleAt
	}

	if err := model.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新任务失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"code":    constant.Success,
		"data": gin.H{
			"task": task,
		},
	})
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "任务ID格式错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	task, err := model.GetTaskById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "任务不存在",
			"code":  constant.NotFound,
		})
		return
	}

	// 检查权限：普通用户只能删除自己的任务
	user := c.MustGet("user").(*model.User)
	if user.Role != constant.RoleAdmin && task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  constant.InsufficientPrivileges,
		})
		return
	}

	if err := model.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除任务失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
		"code":    constant.Success,
	})
}