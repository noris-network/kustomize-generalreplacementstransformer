package grt

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []optFunc
	}
	tests := []struct {
		name    string
		args    args
		want    Transformer
		wantErr bool
	}{
		{
			name: "from file",
			args: args{[]optFunc{WithConfigFile("testdata/new/config-a1.yaml")}},
			want: Transformer{
				config: Config{Values: map[string]interface{}{"a": 1}},
			},
		},
		{
			name:    "from non existing file",
			args:    args{[]optFunc{WithConfigFile("testdata/new/config-404.yaml")}},
			wantErr: true,
		},
		{
			name: "from string",
			args: args{[]optFunc{WithConfigString(`{values: {a: 1}}`)}},
			want: Transformer{
				config: Config{Values: map[string]interface{}{"a": 1}},
			},
		},
		{
			name: "from empty string",
			args: args{[]optFunc{WithConfigString(``)}},
			want: Transformer{
				config: Config{Values: map[string]interface{}{}},
			},
		},
		{
			name: "from bytes",
			args: args{[]optFunc{WithConfigBytes([]byte(`{values: {a: 1}}`))}},
			want: Transformer{
				config: Config{Values: map[string]interface{}{"a": 1}},
			},
		},
		{
			name:    "broken config",
			args:    args{[]optFunc{WithConfigBytes([]byte(`{values: {a: 1`))}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %#v,\n want %#v", got, tt.want)
			}
		})
	}
}
