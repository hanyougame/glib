package core

import (
	"sync"

	sdkactivity "go.temporal.io/sdk/activity"
	sdkworkflow "go.temporal.io/sdk/workflow"
)

// Registry 管理工作流和活动的注册
type Registry struct {
	mu         sync.RWMutex
	workflows  map[string]interface{}
	activities map[string]interface{}
	debug      bool // 控制是否打印调试信息
}

// Module 表示一个可注册的模块
type Module interface {
	RegisterWorkflows(r *Registry)
	RegisterActivities(r *Registry)
	Name() string
}

// WorkerRegistry 用于 Worker 注册的接口
type WorkerRegistry interface {
	RegisterWorkflow(workflow interface{})
	RegisterWorkflowWithOptions(workflow interface{}, options sdkworkflow.RegisterOptions)
	RegisterActivity(activity interface{})
	RegisterActivityWithOptions(activity interface{}, options sdkactivity.RegisterOptions)
}

var globalRegistry = NewRegistry()

// NewRegistry 创建一个新的注册表
func NewRegistry() *Registry {
	return &Registry{
		workflows:  make(map[string]interface{}),
		activities: make(map[string]interface{}),
		debug:      false, // 默认关闭调试
	}
}

// SetDebug 设置调试模式
func (r *Registry) SetDebug(debug bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.debug = debug
}

// RegisterWorkflow 注册工作流
func (r *Registry) RegisterWorkflow(name string, workflow interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.workflows[name] = workflow

	if r.debug {
		// 仅在调试模式下打印
		println("Registered workflow:", name)
	}
}

// RegisterActivity 注册活动
func (r *Registry) RegisterActivity(name string, activity interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.activities[name] = activity

	if r.debug {
		println("Registered activity:", name)
	}
}

// ApplyTo 将所有注册项应用到 Worker
func (r *Registry) ApplyTo(w WorkerRegistry) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, workflow := range r.workflows {
		w.RegisterWorkflow(workflow)
	}

	for _, activity := range r.activities {
		w.RegisterActivity(activity)
	}
}

// GetWorkflows 获取所有已注册的工作流
func (r *Registry) GetWorkflows() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]interface{}, len(r.workflows))
	for k, v := range r.workflows {
		result[k] = v
	}
	return result
}

// GetActivities 获取所有已注册的活动
func (r *Registry) GetActivities() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]interface{}, len(r.activities))
	for k, v := range r.activities {
		result[k] = v
	}
	return result
}

// Count 返回已注册的数量
func (r *Registry) Count() (workflows int, activities int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.workflows), len(r.activities)
}

// Clear 清空注册表
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.workflows = make(map[string]interface{})
	r.activities = make(map[string]interface{})
}

// ========================= 全局函数 =========================

// RegisterWorkflow 全局注册工作流
func RegisterWorkflow(name string, workflow interface{}) {
	globalRegistry.RegisterWorkflow(name, workflow)
}

// RegisterActivity 全局注册活动
func RegisterActivity(name string, activity interface{}) {
	globalRegistry.RegisterActivity(name, activity)
}

// RegisterModule 注册模块
func RegisterModule(module Module) {
	if globalRegistry.debug {
		println("Registering module:", module.Name())
	}
	module.RegisterWorkflows(globalRegistry)
	module.RegisterActivities(globalRegistry)
}

// GetGlobal 获取全局注册表
func GetGlobal() *Registry {
	return globalRegistry
}

// SetGlobalDebug 设置全局调试模式
func SetGlobalDebug(debug bool) {
	globalRegistry.SetDebug(debug)
}
