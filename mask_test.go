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

		// TODO: how to compare map[string]interface{} ?

		//{
		//	name: "SUCCESS: Generic json string with default mask function",
		//	args: args{
		//		jsonString: genericJSONString,
		//		config:     defaultMaskConfig,
		//	},
		//	want:    `{"array":[2,3,1],"boolean":true,"color":"#******","null":null,"number":123,"object":{"a":"b","c":"d","e":"f"},"string":"挨*******************"}`,
		//	wantErr: false,
		//},
		//{
		//	name: "SUCCESS: Generic json string with custom mask function",
		//	args: args{
		//		jsonString: genericJSONString,
		//		config: &MaskConfig{
		//			Callback: func(s string) string {
		//				return string([]rune(s)[:1]) + "$$$"
		//			},
		//		},
		//	},
		//	want:    `{"array":[1,2,3],"boolean":true,"color":"#***","null":null,"number":123,"object":{"a":"b***","c":"d***","e":"f***"},"string":"挨***"}`,
		//	wantErr: false,
		//},
		//{
		//	name: "SUCCESS: Nested object json string with default mask function",
		//	args: args{
		//		jsonString: nestedObjectJSONString,
		//		config:     defaultMaskConfig,
		//	},
		//	want:    `{"object":{"a":{"b":{"c":{"d":"*"}}}}}`,
		//	wantErr: false,
		//},
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
