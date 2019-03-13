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
		return string([]rune(s)[:1]) + "*****"
	}

	defaultMaskFunc = MaskFunc
)

// Mask masks the given json string
func Mask(jsonString string) (string, error) {
	return MaskWithFunc(jsonString, defaultMaskFunc)
}

// MaskWithFunc masks the given json string using the given callback function
func MaskWithFunc(jsonString string, maskFunc func(s string) string) (string, error) {

	before := make(map[string]interface{})
	after := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonString), &before); err != nil {
		return "", err
	}

	mask(before, after, maskFunc)

	b, err := json.Marshal(after)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func mask(before, after map[string]interface{}, callback func(s string) string) {
	for k, v := range before {
		recursiveMask(k, v, after, callback)
	}
}

func recursiveMask(key string, val interface{}, obj map[string]interface{}, callback func(s string) string) {

	switch value := val.(type) {

	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, v := range value {
			recursiveMask(k, v, ret, callback)
		}
		obj[key] = ret

	case []interface{}:
		ret := make(map[string]interface{})
		for k, v := range value {
			i := strconv.Itoa(k)
			recursiveMask(i, v, ret, callback)
		}
		slice := make([]interface{}, 0)
		for _, v := range ret {
			slice = append(slice, v)
		}
		obj[key] = slice

	case string:
		// mask the string value
		obj[key] = callback(value)

	default:
		obj[key] = val
	}
}
