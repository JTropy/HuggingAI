package service

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
)

func GetUserUsableGroups(userGroup string) map[string]string {
	userGroup = strings.TrimSpace(userGroup)
	configuredGroups := setting.GetUserUsableGroupsCopy()
	if userGroup == "" {
		return configuredGroups
	}

	groupsCopy := make(map[string]string)
	if desc, ok := configuredGroups["default"]; ok {
		groupsCopy["default"] = desc
	}
	if userGroup != "" {
		specialSettings, b := ratio_setting.GetGroupRatioSetting().GroupSpecialUsableGroup.Get(userGroup)
		if b {
			// 处理特殊可用分组
			for specialGroup, desc := range specialSettings {
				if strings.HasPrefix(specialGroup, "-:") {
					// 移除分组
					groupToRemove := strings.TrimPrefix(specialGroup, "-:")
					delete(groupsCopy, groupToRemove)
				} else if strings.HasPrefix(specialGroup, "+:") {
					// 添加分组
					groupToAdd := strings.TrimPrefix(specialGroup, "+:")
					if desc == "" {
						desc = setting.GetUsableGroupDescription(groupToAdd)
					}
					groupsCopy[groupToAdd] = desc
				} else {
					// 直接添加分组
					if desc == "" {
						desc = setting.GetUsableGroupDescription(specialGroup)
					}
					groupsCopy[specialGroup] = desc
				}
			}
		}
		// 用户自己的分组始终可用，但全局 UserUsableGroups 里的其它分组不再自动授权给所有用户。
		if _, ok := groupsCopy[userGroup]; !ok {
			desc := setting.GetUsableGroupDescription(userGroup)
			if desc == userGroup {
				desc = "用户分组"
			}
			groupsCopy[userGroup] = desc
		}
	}
	return groupsCopy
}

func GroupInUserUsableGroups(userGroup, groupName string) bool {
	_, ok := GetUserUsableGroups(userGroup)[groupName]
	return ok
}

func CanUseTokenGroup(isAdmin bool, userGroup, groupName string) bool {
	userGroup = strings.TrimSpace(userGroup)
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return true
	}
	if userGroup == "" && !isAdmin {
		return false
	}
	if ratio_setting.IsOwnerGroup(groupName) {
		return isAdmin && ratio_setting.ContainsGroupRatio(groupName)
	}
	if groupName == "auto" {
		return isAdmin || len(GetUserAutoGroup(userGroup)) > 0
	}
	if isAdmin && ratio_setting.ContainsGroupRatio(groupName) {
		return true
	}
	if groupName != "auto" && !ratio_setting.ContainsGroupRatio(groupName) {
		return userGroup != "" && groupName == userGroup && GroupInUserUsableGroups(userGroup, groupName)
	}
	return GroupInUserUsableGroups(userGroup, groupName)
}

func CanSetTokenGroup(role int, userGroup, groupName string) bool {
	return CanUseTokenGroup(role >= common.RoleAdminUser, userGroup, groupName)
}

// GetUserAutoGroup 根据用户分组获取自动分组设置
func GetUserAutoGroup(userGroup string) []string {
	groups := GetUserUsableGroups(userGroup)
	autoGroups := make([]string, 0)
	for _, group := range setting.GetAutoGroups() {
		if _, ok := groups[group]; ok {
			autoGroups = append(autoGroups, group)
		}
	}
	return autoGroups
}

// GetUserGroupRatio 获取用户使用某个分组的倍率
// userGroup 用户分组
// group 需要获取倍率的分组
func GetUserGroupRatio(userGroup, group string) float64 {
	ratio, ok := ratio_setting.GetGroupGroupRatio(userGroup, group)
	if ok {
		return ratio
	}
	if group != "" && !ratio_setting.ContainsGroupRatio(group) {
		return 1
	}
	return ratio_setting.GetGroupRatio(group)
}
