package jsonmask

import (
	"encoding/json"
	"strconv"
	"strings"
)

var (
	// MaskAllFunc is function which masks the whole string
	MaskAllFunc = func(s string) string {
		return strings.Repeat("*", len(s))
	}

	// MaskWithoutFirstCharFunc is function which masks the string except the first character
	MaskWithoutFirstCharFunc = func(s string) string {
		return string([]rune(s)[:1]) + strings.Repeat("*", len(s)-1)
	}

	defaultMaskFunc   = MaskWithoutFirstCharFunc
	defaultMaskConfig = &MaskConfig{
		Callback: defaultMaskFunc,
	}
)

type MaskConfig struct {
	Callback   func(s string) string
	SkipFields []string
}

// Mask masks the given json string with config
func Mask(jsonString string, config ...*MaskConfig) (string, error) {

	before := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonString), &before); err != nil {
		return "", err
	}

	cfg := defaultMaskConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.Callback == nil {
		cfg.Callback = defaultMaskFunc
	}

	skipFieldMap := make(map[string]bool, len(cfg.SkipFields))
	for _, skipField := range cfg.SkipFields {
		skipFieldMap[skipField] = true
	}

	after := make(map[string]interface{})
	for k, v := range before {
		mask(k, v, after, cfg.Callback, skipFieldMap)
	}

	b, err := json.Marshal(after)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func mask(key string, val interface{}, obj map[string]interface{}, callback func(s string) string, skipFieldMap map[string]bool) {

	switch value := val.(type) {

	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, v := range value {
			mask(k, v, ret, callback, skipFieldMap)
		}
		obj[key] = ret

	case []interface{}:
		ret := make(map[string]interface{})
		for k, v := range value {
			i := strconv.Itoa(k)
			mask(i, v, ret, callback, skipFieldMap)
		}
		slice := make([]interface{}, 0)
		for _, v := range ret {
			slice = append(slice, v)
		}
		obj[key] = slice

	case string:
		if _, doSkip := skipFieldMap[key]; doSkip {
			// skip masking and set the original value
			obj[key] = value
		} else {
			// mask and set the modified value
			obj[key] = callback(value)
		}

	default:
		obj[key] = val
	}
}
