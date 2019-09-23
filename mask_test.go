package jsonmask

import (
	"testing"
)

const (
	simpleJSONString = `
{
  "string": "挨拶はHello World"
}
`
	genericJSONString = `
{
  "array": [
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

func TestMask(t *testing.T) {
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
			name: "SUCCESS: simple json string",
			args: args{
				jsonString: simpleJSONString,
			},
			want:    `{"string":"挨*******************"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: simple json string with config: custom callback",
			args: args{
				jsonString: simpleJSONString,
				config: &MaskConfig{
					Callback: func(s string) string {
						return string([]rune(s)[:1]) + "$$$"
					},
				},
			},
			want:    `{"string":"挨$$$"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: simple json string with config: skip fields",
			args: args{
				jsonString: simpleJSONString,
				config: &MaskConfig{
					SkipFields: []string{"string"},
				},
			},
			want:    `{"string":"挨拶はHello World"}`,
			wantErr: false,
		},

		{
			name: "SUCCESS: Generic json string with default mask function",
			args: args{
				jsonString: genericJSONString,
				config:     defaultMaskConfig,
			},
			want:    `{"array":[3],"boolean":true,"color":"#******","null":null,"number":123,"object":{"a":"b","c":"d","e":"f"},"string":"挨*******************"}`,
			wantErr: false,
		},
		{
			name: "SUCCESS: Generic json string with config",
			args: args{
				jsonString: genericJSONString,
				config: &MaskConfig{
					Callback: func(s string) string {
						return string([]rune(s)[:1]) + "$$$"
					},
					SkipFields: []string{"color"},
				},
			},
			want:    `{"array":[3],"boolean":true,"color":"#82b92c","null":null,"number":123,"object":{"a":"b$$$","c":"d$$$","e":"f$$$"},"string":"挨$$$"}`,
			wantErr: false,
		},
		{
			// TODO: future task
			name: "SUCCESS: Nested object json string is not masked",
			args: args{
				jsonString: nestedObjectJSONString,
				config:     defaultMaskConfig,
			},
			want:    `{"object":{"a":{"b":{"c":{"d":"e"}}}}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			if tt.args.config == nil {
				got, err = Mask(tt.args.jsonString)
			} else {
				got, err = Mask(tt.args.jsonString, tt.args.config)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Mask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
