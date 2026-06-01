package ratio_setting

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/types"
)

const OwnerGroupName = "owner"

var defaultGroupRatio = map[string]float64{
	"default": 1,
	"vip":     1,
	"svip":    1,
	"owner":   0,
}

var groupRatioMap = types.NewRWMap[string, float64]()

var defaultGroupGroupRatio = map[string]map[string]float64{
	"vip": {
		"edit_this": 0.9,
	},
}

var groupGroupRatioMap = types.NewRWMap[string, map[string]float64]()

var defaultGroupSpecialUsableGroup = map[string]map[string]string{
	"vip": {
		"append_1":   "vip_special_group_1",
		"-:remove_1": "vip_removed_group_1",
	},
}

type GroupRatioSetting struct {
	GroupRatio              *types.RWMap[string, float64]            `json:"group_ratio"`
	GroupGroupRatio         *types.RWMap[string, map[string]float64] `json:"group_group_ratio"`
	GroupSpecialUsableGroup *types.RWMap[string, map[string]string]  `json:"group_special_usable_group"`
}

var groupRatioSetting GroupRatioSetting

func init() {
	groupSpecialUsableGroup := types.NewRWMap[string, map[string]string]()
	groupSpecialUsableGroup.AddAll(defaultGroupSpecialUsableGroup)

	groupRatioMap.AddAll(defaultGroupRatio)
	groupGroupRatioMap.AddAll(defaultGroupGroupRatio)

	groupRatioSetting = GroupRatioSetting{
		GroupSpecialUsableGroup: groupSpecialUsableGroup,
		GroupRatio:              groupRatioMap,
		GroupGroupRatio:         groupGroupRatioMap,
	}

	config.GlobalConfig.Register("group_ratio_setting", &groupRatioSetting)
}

func GetGroupRatioSetting() *GroupRatioSetting {
	if groupRatioSetting.GroupSpecialUsableGroup == nil {
		groupRatioSetting.GroupSpecialUsableGroup = types.NewRWMap[string, map[string]string]()
		groupRatioSetting.GroupSpecialUsableGroup.AddAll(defaultGroupSpecialUsableGroup)
	}
	return &groupRatioSetting
}

func IsOwnerGroup(name string) bool {
	return strings.EqualFold(strings.TrimSpace(name), OwnerGroupName)
}

func normalizeGroupRatioValue(groupName string, ratio float64) float64 {
	if ratio == 0 && !IsOwnerGroup(groupName) {
		return 1
	}
	return ratio
}

func normalizeGroupRatioMap(ratios map[string]float64) map[string]float64 {
	normalized := make(map[string]float64, len(ratios))
	for name, ratio := range ratios {
		normalized[name] = normalizeGroupRatioValue(name, ratio)
	}
	return normalized
}

func normalizeGroupGroupRatioMap(ratios map[string]map[string]float64) map[string]map[string]float64 {
	normalized := make(map[string]map[string]float64, len(ratios))
	for userGroup, overrides := range ratios {
		normalized[userGroup] = make(map[string]float64, len(overrides))
		for usingGroup, ratio := range overrides {
			normalized[userGroup][usingGroup] = normalizeGroupRatioValue(usingGroup, ratio)
		}
	}
	return normalized
}

func NormalizeGroupRatioJSONString(jsonStr string) (string, error) {
	checkGroupRatio := make(map[string]float64)
	err := common.Unmarshal([]byte(jsonStr), &checkGroupRatio)
	if err != nil {
		return "", err
	}
	for name, ratio := range checkGroupRatio {
		if ratio < 0 {
			return "", errors.New("group ratio must be not less than 0: " + name)
		}
	}
	jsonBytes, err := common.Marshal(normalizeGroupRatioMap(checkGroupRatio))
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func NormalizeGroupGroupRatioJSONString(jsonStr string) (string, error) {
	checkGroupGroupRatio := make(map[string]map[string]float64)
	err := common.Unmarshal([]byte(jsonStr), &checkGroupGroupRatio)
	if err != nil {
		return "", err
	}
	for userGroup, overrides := range checkGroupGroupRatio {
		for usingGroup, ratio := range overrides {
			if ratio < 0 {
				return "", errors.New("group-group ratio must be not less than 0: " + userGroup + " -> " + usingGroup)
			}
		}
	}
	jsonBytes, err := common.Marshal(normalizeGroupGroupRatioMap(checkGroupGroupRatio))
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func GetGroupRatioCopy() map[string]float64 {
	return normalizeGroupRatioMap(groupRatioMap.ReadAll())
}

func ContainsGroupRatio(name string) bool {
	_, ok := groupRatioMap.Get(name)
	return ok
}

func GroupRatio2JSONString() string {
	jsonBytes, err := common.Marshal(GetGroupRatioCopy())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func UpdateGroupRatioByJSONString(jsonStr string) error {
	normalizedJSON, err := NormalizeGroupRatioJSONString(jsonStr)
	if err != nil {
		return err
	}
	return types.LoadFromJsonString(groupRatioMap, normalizedJSON)
}

func GetGroupRatio(name string) float64 {
	ratio, ok := groupRatioMap.Get(name)
	if !ok {
		common.SysLog("group ratio not found: " + name)
		return 1
	}
	return normalizeGroupRatioValue(name, ratio)
}

func GetGroupGroupRatio(userGroup, usingGroup string) (float64, bool) {
	gp, ok := groupGroupRatioMap.Get(userGroup)
	if !ok {
		return -1, false
	}
	ratio, ok := gp[usingGroup]
	if !ok {
		return -1, false
	}
	return normalizeGroupRatioValue(usingGroup, ratio), true
}

func GroupGroupRatio2JSONString() string {
	jsonBytes, err := common.Marshal(normalizeGroupGroupRatioMap(groupGroupRatioMap.ReadAll()))
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func UpdateGroupGroupRatioByJSONString(jsonStr string) error {
	normalizedJSON, err := NormalizeGroupGroupRatioJSONString(jsonStr)
	if err != nil {
		return err
	}
	return types.LoadFromJsonString(groupGroupRatioMap, normalizedJSON)
}

func CheckGroupRatio(jsonStr string) error {
	checkGroupRatio := make(map[string]float64)
	err := common.Unmarshal([]byte(jsonStr), &checkGroupRatio)
	if err != nil {
		return err
	}
	for name, ratio := range checkGroupRatio {
		if ratio < 0 {
			return errors.New("group ratio must be not less than 0: " + name)
		}
	}
	return nil
}
