package service

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"time"
)

type TaskScheduler struct {
	running bool
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		running: false,
	}
}

func (ts *TaskScheduler) Start() {
	if ts.running {
		return
	}

	ts.running = true
	common.SysLog("Task scheduler started")

	// 启动定时任务处理器
	go ts.processScheduledTasks()
}

func (ts *TaskScheduler) Stop() {
	ts.running = false
	common.SysLog("Task scheduler stopped")
}

func (ts *TaskScheduler) processScheduledTasks() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for ts.running {
		select {
		case <-ticker.C:
			ts.executePendingTasks()
		}
	}
}

func (ts *TaskScheduler) executePendingTasks() {
	// 获取待执行的任务
	tasks, err := model.GetTasksByStatus(constant.TaskStatusPending)
	if err != nil {
		common.SysError("Failed to get pending tasks: " + err.Error())
		return
	}

	now := time.Now()
	for _, task := range tasks {
		// 检查是否到了执行时间
		if task.ScheduleAt.Before(now) || task.ScheduleAt.Equal(now) {
			go ts.executeTask(task)
		}
	}
}

func (ts *TaskScheduler) executeTask(task model.Task) {
	common.SysLog("Executing task: " + task.Title)

	// 更新任务状态为运行中
	err := model.UpdateTaskStatus(task.ID, constant.TaskStatusRunning)
	if err != nil {
		common.SysError("Failed to update task status: " + err.Error())
		return
	}

	// 模拟任务执行
	time.Sleep(2 * time.Second)

	// 这里可以根据具体的任务类型执行不同的逻辑
	// 目前只是简单的模拟执行完成

	// 更新任务状态为已完成
	err = model.UpdateTaskStatus(task.ID, constant.TaskStatusCompleted)
	if err != nil {
		common.SysError("Failed to complete task: " + err.Error())
		// 如果更新失败，标记为失败状态
		model.UpdateTaskStatus(task.ID, constant.TaskStatusFailed)
		return
	}

	common.SysLog("Task completed: " + task.Title)
}
