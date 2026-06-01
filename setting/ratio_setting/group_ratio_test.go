package ratio_setting

import (
	"testing"
)

func TestGroupRatioNormalizesPublicZeroButKeepsOwnerZero(t *testing.T) {
	original := GroupRatio2JSONString()
	t.Cleanup(func() {
		if err := UpdateGroupRatioByJSONString(original); err != nil {
			t.Fatalf("failed to restore group ratios: %v", err)
		}
	})

	if err := UpdateGroupRatioByJSONString(`{"default":0,"vip":0.5,"owner":0}`); err != nil {
		t.Fatalf("failed to update group ratios: %v", err)
	}

	if got := GetGroupRatio("default"); got != 1 {
		t.Fatalf("default ratio = %v, want 1", got)
	}
	if got := GetGroupRatio("vip"); got != 0.5 {
		t.Fatalf("vip ratio = %v, want 0.5", got)
	}
	if got := GetGroupRatio("owner"); got != 0 {
		t.Fatalf("owner ratio = %v, want 0", got)
	}
}

func TestGroupGroupRatioNormalizesByUsingGroup(t *testing.T) {
	original := GroupGroupRatio2JSONString()
	t.Cleanup(func() {
		if err := UpdateGroupGroupRatioByJSONString(original); err != nil {
			t.Fatalf("failed to restore group-group ratios: %v", err)
		}
	})

	if err := UpdateGroupGroupRatioByJSONString(`{"default":{"default":0,"owner":0},"owner":{"vip":0}}`); err != nil {
		t.Fatalf("failed to update group-group ratios: %v", err)
	}

	if got, ok := GetGroupGroupRatio("default", "default"); !ok || got != 1 {
		t.Fatalf("default -> default ratio = %v (ok=%v), want 1", got, ok)
	}
	if got, ok := GetGroupGroupRatio("default", "owner"); !ok || got != 0 {
		t.Fatalf("default -> owner ratio = %v (ok=%v), want 0", got, ok)
	}
	if got, ok := GetGroupGroupRatio("owner", "vip"); !ok || got != 1 {
		t.Fatalf("owner -> vip ratio = %v (ok=%v), want 1", got, ok)
	}
}
