package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/gin-gonic/gin"
)

func TestGetUserGroupsIncludesOwnGroupWithoutRatio(t *testing.T) {
	originalRatio := ratio_setting.GroupRatio2JSONString()
	originalUsable := setting.UserUsableGroups2JSONString()
	t.Cleanup(func() {
		if err := ratio_setting.UpdateGroupRatioByJSONString(originalRatio); err != nil {
			t.Fatalf("failed to restore group ratios: %v", err)
		}
		if err := setting.UpdateUserUsableGroupsByJSONString(originalUsable); err != nil {
			t.Fatalf("failed to restore user usable groups: %v", err)
		}
	})

	if err := ratio_setting.UpdateGroupRatioByJSONString(`{"vip":1}`); err != nil {
		t.Fatalf("failed to set group ratios: %v", err)
	}
	if err := setting.UpdateUserUsableGroupsByJSONString(`{"default":"默认分组","vip":"vip分组"}`); err != nil {
		t.Fatalf("failed to set user usable groups: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/groups", func(c *gin.Context) {
		c.Set("id", 1)
		c.Set("role", common.RoleCommonUser)
		c.Set("user_group", "default")
		GetUserGroups(c)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/groups", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var body struct {
		Success bool                              `json:"success"`
		Data    map[string]map[string]interface{} `json:"data"`
	}
	if err := common.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !body.Success {
		t.Fatal("response success = false")
	}
	if _, ok := body.Data["default"]; !ok {
		t.Fatalf("default group missing from response: %#v", body.Data)
	}
	if _, ok := body.Data["vip"]; ok {
		t.Fatalf("vip should not be available to default users: %#v", body.Data)
	}
}
