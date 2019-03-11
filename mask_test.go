package go_json_mask

import (
	"testing"
)

const (
	genericJSONString = `
{
  "array": [
    1,
    2,
    3
  ],
  "boolean": true,
  "color": "#82b92c",
  "null": null,
  "number": 123,
  "object": {
    "a": "b",
    "c": "d",
    "e": "f"
  },
  "string": "挨拶はHello World"
}
`
	nestedObjectJSONString = `
{
  "object": {
    "a": {
      "b": {
        "c": {
          "d": "e"
        }
      }
    }
  }
}
`
)

func TestMaskWithFunc(t *testing.T) {
	type args struct {
		jsonString string
		maskFunc   func(s string) string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "SUCCESS: Generic json string with default mask function",
			args: args{
				jsonString: genericJSONString,
				maskFunc:   defaultMaskFunc,
			},
			want:    `{"array":[1,2,3],"boolean":true,"color":"*******","null":null,"number":123,"object":{"a":"*","c":"*","e":"*"},"string":"********************"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: Generic json string with custom mask function",
			args: args{
				jsonString: genericJSONString,
				maskFunc: func(s string) string {
					return string([]rune(s)[:1]) + "***"
				},
			},
			want:    `{"array":[1,2,3],"boolean":true,"color":"#***","null":null,"number":123,"object":{"a":"b***","c":"d***","e":"f***"},"string":"挨***"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: Nested object json string with default mask function",
			args: args{
				jsonString: nestedObjectJSONString,
				maskFunc:   defaultMaskFunc,
			},
			want:    `{"object":{"a":{"b":{"c":{"d":"*"}}}}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaskWithFunc(tt.args.jsonString, tt.args.maskFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskWithFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MaskWithFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
