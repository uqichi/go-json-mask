package jsonmask

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

func TestMaskWithConfig(t *testing.T) {
	type args struct {
		jsonString string
		config     *MaskConfig
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
				config:     defaultMaskConfig,
			},
			want:    `{"array":[2,3,1],"boolean":true,"color":"#******","null":null,"number":123,"object":{"a":"b","c":"d","e":"f"},"string":"挨*******************"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: Generic json string with custom mask function",
			args: args{
				jsonString: genericJSONString,
				config: &MaskConfig{
					Callback: func(s string) string {
						return string([]rune(s)[:1]) + "$$$"
					},
				},
			},
			want:    `{"array":[1,2,3],"boolean":true,"color":"#***","null":null,"number":123,"object":{"a":"b***","c":"d***","e":"f***"},"string":"挨***"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: Nested object json string with default mask function",
			args: args{
				jsonString: nestedObjectJSONString,
				config:     defaultMaskConfig,
			},
			want:    `{"object":{"a":{"b":{"c":{"d":"*"}}}}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaskWithConfig(tt.args.jsonString, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MaskWithConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
