package grt

import (
	"reflect"
	"testing"
)

func TestTransformer_ApplyValuesToValues(t *testing.T) {
	tests := []struct {
		name     string
		values   map[string]any
		maxDepth int
		want     map[string]any
		wantErr  bool
	}{
		{
			name: "nop",
			values: map[string]any{
				"user": "admin",
				"pass": "s3cr3t",
			},
			want: map[string]any{
				"user": "admin",
				"pass": "s3cr3t",
			},
		},
		{
			name: "basic replace",
			values: map[string]any{
				"user": "admin",
				"pass": "s3cr3t",
				"url":  "https://{{.user}}:{{.pass}}@example.com",
			},
			want: map[string]any{
				"user": "admin",
				"pass": "s3cr3t",
				"url":  "https://admin:s3cr3t@example.com",
			},
		},
		{
			name: "complex replace",
			values: map[string]any{
				"user":   "admin",
				"pass":   "s3cr3t",
				"domain": "example.com",
				"db": map[string]any{
					"host": "dbhost.{{.domain}}",
					"port": 5432,
					"url":  "postgresql://{{.user}}:{{.pass}}@{{.db.host}}:{{.db.port}}",
					"x": map[string]any{
						"y": map[string]any{
							"z": map[string]any{
								"foohost": "{{.db.host}}",
							},
						},
					},
				},
			},
			want: map[string]any{
				"user":   "admin",
				"pass":   "s3cr3t",
				"domain": "example.com",
				"db": map[string]any{
					"host": "dbhost.example.com",
					"port": 5432,
					"url":  "postgresql://admin:s3cr3t@dbhost.example.com:5432",
					"x": map[string]any{
						"y": map[string]any{
							"z": map[string]any{
								"foohost": "dbhost.example.com",
							},
						},
					},
				},
			},
		},
		{
			name: "nested",
			values: map[string]any{
				"a": "*",
				"b": "{{.a}}",
				"c": "{{.b}}",
				"d": "{{.c}}",
				"e": "{{.d}}",
				"f": "{{.e}}",
			},
			want: map[string]any{
				"a": "*",
				"b": "*",
				"c": "*",
				"d": "*",
				"e": "*",
				"f": "*",
			},
		},
		{
			name:    "exceed maxDeep",
			values:  map[string]any{"a": "{{.a}}b"},
			wantErr: true,
		},
		{
			name:    "parse error",
			values:  map[string]any{"a": "{{.a}"},
			wantErr: true,
		},
		{
			name:    "no value error",
			values:  map[string]any{"a": "{{.b}}"},
			wantErr: true,
		},
		{
			name:    "exec error",
			values:  map[string]any{"a": "{{div 1 0}}"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{values: tt.values}
			if err := tr.ApplyValuesToValues(); (err != nil) != tt.wantErr {
				t.Errorf("Transformer.ApplyValuesToValues() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.wantErr {
					return
				}
			}
			if !reflect.DeepEqual(tr.values, tt.want) {
				t.Errorf("Transformer.ApplyValuesToValues()\ngot:   %#v\nwant:  %#v\n", tr.values, tt.want)
			}
		})
	}
}
