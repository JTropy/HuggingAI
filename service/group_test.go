package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
)

func TestCanSetTokenGroupAllowsOwnGroupWithoutRatio(t *testing.T) {
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
	if err := setting.UpdateUserUsableGroupsByJSONString(`{}`); err != nil {
		t.Fatalf("failed to set user usable groups: %v", err)
	}

	if !CanSetTokenGroup(common.RoleCommonUser, "default", "default") {
		t.Fatal("ordinary users should be able to set their own group as a token group")
	}
	if CanSetTokenGroup(common.RoleCommonUser, "default", "missing") {
		t.Fatal("ordinary users should not be able to set unrelated missing groups")
	}
	if ratio := GetUserGroupRatio("default", "default"); ratio != 1 {
		t.Fatalf("own missing group ratio = %v, want 1", ratio)
	}
}

func TestDefaultUserCannotSetVipWhenVipIsOnlyRatioGroup(t *testing.T) {
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

	if err := ratio_setting.UpdateGroupRatioByJSONString(`{"default":1,"vip":1}`); err != nil {
		t.Fatalf("failed to set group ratios: %v", err)
	}
	if err := setting.UpdateUserUsableGroupsByJSONString(`{"default":"默认分组"}`); err != nil {
		t.Fatalf("failed to set user usable groups: %v", err)
	}

	if CanSetTokenGroup(common.RoleCommonUser, "default", "vip") {
		t.Fatal("ordinary default users should not be able to set vip tokens")
	}
	if !CanSetTokenGroup(common.RoleCommonUser, "default", "default") {
		t.Fatal("ordinary default users should be able to set default tokens")
	}
}

func TestDefaultUserCannotSetVipWhenVipIsGloballyConfigured(t *testing.T) {
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

	if err := ratio_setting.UpdateGroupRatioByJSONString(`{"default":1,"vip":1}`); err != nil {
		t.Fatalf("failed to set group ratios: %v", err)
	}
	if err := setting.UpdateUserUsableGroupsByJSONString(`{"default":"默认分组","vip":"vip分组"}`); err != nil {
		t.Fatalf("failed to set user usable groups: %v", err)
	}

	defaultGroups := GetUserUsableGroups("default")
	if _, ok := defaultGroups["vip"]; ok {
		t.Fatalf("vip should not be available to default users: %#v", defaultGroups)
	}
	if CanSetTokenGroup(common.RoleCommonUser, "default", "vip") {
		t.Fatal("ordinary default users should not be able to set vip tokens")
	}

	vipGroups := GetUserUsableGroups("vip")
	if _, ok := vipGroups["vip"]; !ok {
		t.Fatalf("vip users should keep access to their own group: %#v", vipGroups)
	}
	if _, ok := vipGroups["default"]; !ok {
		t.Fatalf("vip users should keep access to default when it is configured: %#v", vipGroups)
	}
}
