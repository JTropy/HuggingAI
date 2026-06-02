package setting

import "testing"

func TestDefaultUserUsableGroupsDoNotGrantVip(t *testing.T) {
	groups := GetUserUsableGroupsCopy()
	if _, ok := groups["default"]; !ok {
		t.Fatalf("default group missing from default user usable groups: %#v", groups)
	}
	if _, ok := groups["vip"]; ok {
		t.Fatalf("vip should not be globally selectable by default: %#v", groups)
	}
}

func TestUpdateUserUsableGroupsInvalidPreservesExisting(t *testing.T) {
	original := UserUsableGroups2JSONString()
	t.Cleanup(func() {
		if err := UpdateUserUsableGroupsByJSONString(original); err != nil {
			t.Fatalf("failed to restore user usable groups: %v", err)
		}
	})

	if err := UpdateUserUsableGroupsByJSONString(`{"default":"Default group"}`); err != nil {
		t.Fatalf("failed to set user usable groups: %v", err)
	}
	if err := UpdateUserUsableGroupsByJSONString(`{"broken"`); err == nil {
		t.Fatal("expected invalid JSON to fail")
	}

	groups := GetUserUsableGroupsCopy()
	if got := groups["default"]; got != "Default group" {
		t.Fatalf("default group description = %q, want %q", got, "Default group")
	}
	if len(groups) != 1 {
		t.Fatalf("groups length = %d, want 1", len(groups))
	}
}
