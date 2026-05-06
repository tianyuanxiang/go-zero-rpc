// Package casbin 提供基于Casbin的RBAC权限控制初始化和辅助工具函数。
// 使用 gorm-adapter/v3 将策略数据持久化存储到MySQL数据库。
package casbin

import (
	"fmt"

	casbinv2 "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbin 初始化Casbin执行器（Enforcer）。
//
// 使用MySQL作为策略存储后端，通过gorm-adapter连接。
// 模型文件定义了RBAC策略的匹配规则。
//
// 参数：
//   - dataSource  : MySQL连接字符串
//   - modelPath   : Casbin RBAC模型配置文件路径（如：etc/rbac_model.conf）
//
// 返回：
//   - *casbinv2.Enforcer : 初始化完成的Casbin执行器
//   - error              : 初始化失败时的错误信息
func NewCasbin(db *gorm.DB, modelPath string) (*casbinv2.Enforcer, error) {
	// 初始化gorm-adapter，自动创建casbin_rule表（如不存在）
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("casbin: 初始化gorm-adapter失败: %w", err)
	}

	// 加载模型配置文件并创建执行器
	enforcer, err := casbinv2.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin: 创建Enforcer失败: %w", err)
	}

	// 从数据库加载策略到内存
	if err = enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("casbin: 加载策略失败: %w", err)
	}

	return enforcer, nil
}

// AddPolicyForRole 为角色添加API访问策略。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码（casbin中的subject）
//   - path     : API路径（如：/api/system/user）
//   - method   : HTTP方法（如：GET、POST、PUT、DELETE）
//
// 返回：
//   - bool  : true表示新增成功，false表示策略已存在
//   - error : 添加失败时的错误信息
func AddPolicyForRole(enforcer *casbinv2.Enforcer, roleCode, path, method string) (bool, error) {
	return enforcer.AddPolicy(roleCode, path, method)
}

// RemovePolicyForRole 删除角色的某条API访问策略。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码
//   - path     : API路径
//   - method   : HTTP方法
//
// 返回：
//   - bool  : true表示删除成功，false表示策略不存在
//   - error : 删除失败时的错误信息
func RemovePolicyForRole(enforcer *casbinv2.Enforcer, roleCode, path, method string) (bool, error) {
	return enforcer.RemovePolicy(roleCode, path, method)
}

// RemoveAllPoliciesForRole 删除指定角色的所有API访问策略。
//
// 通常在删除角色或重新分配角色权限前调用。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码
//
// 返回：
//   - error : 删除失败时的错误信息
func RemoveAllPoliciesForRole(enforcer *casbinv2.Enforcer, roleCode string) error {
	_, err := enforcer.RemoveFilteredPolicy(0, roleCode)
	return err
}

// AddRolePolicies 批量为角色添加API访问策略。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码
//   - rules    : 策略规则列表，每条规则为 [path, method]
//
// 返回：
//   - error : 批量添加失败时的错误信息
func AddRolePolicies(enforcer *casbinv2.Enforcer, roleCode string, rules [][]string) error {
	// 先删除该角色的所有旧策略，再批量写入新策略（全量覆盖模式）
	if err := RemoveAllPoliciesForRole(enforcer, roleCode); err != nil {
		return fmt.Errorf("casbin: 清除旧策略失败: %w", err)
	}

	if len(rules) == 0 {
		return nil
	}

	// 构建完整的策略规则（在最前面插入roleCode）
	fullRules := make([][]string, 0, len(rules))
	for _, rule := range rules {
		if len(rule) != 2 {
			continue
		}
		fullRules = append(fullRules, []string{roleCode, rule[0], rule[1]})
	}

	if len(fullRules) == 0 {
		return nil
	}

	_, err := enforcer.AddPolicies(fullRules)
	if err != nil {
		return fmt.Errorf("casbin: 批量添加策略失败: %w", err)
	}

	return nil
}

// CheckPermission 检查指定角色是否对某API有访问权限。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码
//   - path     : API路径
//   - method   : HTTP方法
//
// 返回：
//   - bool  : true表示有权限，false表示无权限
//   - error : 检查过程出错时的错误信息
func CheckPermission(enforcer *casbinv2.Enforcer, roleCode, path, method string) (bool, error) {
	return enforcer.Enforce(roleCode, path, method)
}

// GetRolePolicies 获取指定角色的所有API访问策略。
//
// 参数：
//   - enforcer : Casbin执行器
//   - roleCode : 角色编码
//
// 返回：
//   - [][]string : 策略列表，每条策略为 [sub, obj, act] 格式
//   - error      : 查询失败时的错误信息
func GetRolePolicies(enforcer *casbinv2.Enforcer, roleCode string) ([][]string, error) {
	return enforcer.GetFilteredPolicy(0, roleCode)
}

// ReloadPolicy 重新从数据库加载策略（用于多实例部署时同步策略变更）。
//
// 参数：
//   - enforcer : Casbin执行器
//
// 返回：
//   - error : 加载失败时的错误信息
func ReloadPolicy(enforcer *casbinv2.Enforcer) error {
	return enforcer.LoadPolicy()
}
