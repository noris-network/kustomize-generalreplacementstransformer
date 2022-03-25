package grt

import (
	"os"
	"reflect"
	"testing"
)

func TestTransformer_ScanForValues(t *testing.T) {
	type args struct {
		configFile string
		inputFile  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    map[string]any
	}{
		{
			name: "plain and default values",
			args: args{
				configFile: "testdata/scan-for-values/config-1.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"namespace-with-default": "default-namespace",
				"namespace-no-default":   "",
				"namespace":              "myspace",
				"mySecretName":           "GeneralReplacementsTransformer",
			},
		},
		{
			name: "slurp configmap",
			args: args{
				configFile: "testdata/scan-for-values/config-2.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"db": map[string]string{
					"host": "mydbhost",
					"name": "myexampledb",
				},
			},
		},
		{
			name: "slurp and splat configmap",
			args: args{
				configFile: "testdata/scan-for-values/config-3.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"host": "mydbhost",
				"name": "myexampledb",
			},
		},
		{
			name: "broken configmap",
			args: args{
				configFile: "testdata/scan-for-values/config-3.yaml",
				inputFile:  "testdata/scan-for-values/input-2.yaml",
			},
			wantErr: true,
		},
		{
			name: "broken configmap, no data",
			args: args{
				configFile: "testdata/scan-for-values/config-3.yaml",
				inputFile:  "testdata/scan-for-values/input-3.yaml",
			},
			wantErr: true,
		},
		{
			name: "slurp secret",
			args: args{
				configFile: "testdata/scan-for-values/config-4.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"secret": map[string]string{
					"PASSWORD":  "GeneralReplacementsTransformer",
					"FOO_TOKEN": "foo-bar-baz-123",
				},
			},
		},
		{
			name: "slurp and splat secret",
			args: args{
				configFile: "testdata/scan-for-values/config-5.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"PASSWORD":  "GeneralReplacementsTransformer",
				"FOO_TOKEN": "foo-bar-baz-123",
			},
		},
		{
			name: "empty name",
			args: args{
				configFile: "testdata/scan-for-values/config-7.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			wantErr: true,
		},
		{
			name: "empty resource selector",
			args: args{
				configFile: "testdata/scan-for-values/config-8.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			wantErr: true,
		},
		{
			name: "advanced selection",
			args: args{
				configFile: "testdata/scan-for-values/config-9.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{
				"exporterImage": "myexporter:v0.4.2",
			},
		},
		{
			name: "broken advanced selection",
			args: args{
				configFile: "testdata/scan-for-values/config-10.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			wantErr: true,
		},
		{
			name: "no singe value",
			args: args{
				configFile: "testdata/scan-for-values/config-12.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			wantErr: true,
		},
		{
			name: "no match",
			args: args{
				configFile: "testdata/scan-for-values/config-13.yaml",
				inputFile:  "testdata/scan-for-values/input-1.yaml",
			},
			want: map[string]interface{}{"exporterImage": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := New(WithConfigFile(tt.args.configFile))
			if err != nil {
				t.Fatalf("New(): error = %v", err)
			}
			f, err := os.Open(tt.args.inputFile)
			if err != nil {
				t.Fatalf("Open(): error = %v", err)
			}
			defer f.Close()
			if err = tr.ReadStream(f); err != nil {
				t.Fatalf("ReadStream(): error = %v", err)
			}
			if err := tr.ScanForValues(); (err != nil) != tt.wantErr {
				t.Fatalf("Transformer.ScanForValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tr.values, tt.want) {
				t.Errorf("Transformer.ScanForValues()\nvalues: %#v\nwant:   %#v\n", tr.values, tt.want)
			}
		})
	}
}
