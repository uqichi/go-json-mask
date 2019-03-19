package jsonmask

import (
	"encoding/json"
	"strconv"
	"strings"
)

var (
	// MaskFunc is function which masks the string with its length same
	MaskFunc = func(s string) string {
		return strings.Repeat("*", len(s))
	}

	// MaskWithoutFirstOneCharFunc is function which masks the string except the first character
	MaskWithoutFirstOneCharFunc = func(s string) string {
		return string([]rune(s)[:1]) + strings.Repeat("*", len(s)-1)
	}

	defaultMaskFunc = MaskWithoutFirstOneCharFunc
)

// Mask masks the given json string
func Mask(jsonString string, skipFields []string) (string, error) {
	return MaskWithFunc(jsonString, defaultMaskFunc, skipFields)
}

// MaskWithFunc masks the given json string using the given callback function
func MaskWithFunc(jsonString string, maskFunc func(s string) string, skipFields []string) (string, error) {

	before := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonString), &before); err != nil {
		return "", err
	}

	skipFieldMap := make(map[string]bool, len(skipFields))
	for _, skipField := range skipFields {
		skipFieldMap[skipField] = true
	}

	after := make(map[string]interface{})
	for k, v := range before {
		mask(k, v, after, maskFunc, skipFieldMap)
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
			// skip masking
			obj[key] = value
		} else {
			// mask the string value
			obj[key] = callback(value)
		}

	default:
		obj[key] = val
	}
}
