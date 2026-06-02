package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	groupNames := make([]string, 0)
	for groupName := range ratio_setting.GetGroupRatioCopy() {
		groupNames = append(groupNames, groupName)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groupNames,
	})
}

func GetUserGroups(c *gin.Context) {
	usableGroups := make(map[string]map[string]interface{})
	userGroup := ""
	userId := c.GetInt("id")
	role := c.GetInt("role")
	if contextUserGroup := c.GetString("user_group"); contextUserGroup != "" {
		userGroup = contextUserGroup
	} else if userId > 0 {
		userGroup, _ = model.GetUserGroup(userId, false)
	}
	userUsableGroups := service.GetUserUsableGroups(userGroup)
	candidateGroups := ratio_setting.GetGroupRatioCopy()
	for groupName := range userUsableGroups {
		if _, ok := candidateGroups[groupName]; !ok {
			candidateGroups[groupName] = service.GetUserGroupRatio(userGroup, groupName)
		}
	}
	for groupName := range candidateGroups {
		if !service.CanSetTokenGroup(role, userGroup, groupName) {
			continue
		}
		// UserUsableGroups contains the groups that the user can use
		desc, ok := userUsableGroups[groupName]
		if !ok {
			desc = setting.GetUsableGroupDescription(groupName)
		}
		usableGroups[groupName] = map[string]interface{}{
			"ratio": service.GetUserGroupRatio(userGroup, groupName),
			"desc":  desc,
		}
	}
	if _, ok := userUsableGroups["auto"]; ok && service.CanSetTokenGroup(role, userGroup, "auto") {
		usableGroups["auto"] = map[string]interface{}{
			"ratio": "自动",
			"desc":  setting.GetUsableGroupDescription("auto"),
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    usableGroups,
	})
}
