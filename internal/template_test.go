package grt

import (
	"os"
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestTransformer_TemplateTransform(t *testing.T) {
	type fields struct {
		config Config
		values map[string]any
	}
	type args struct {
		repl      Replacement
		inputFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "simple replace",
			fields: fields{values: map[string]any{
				"hostname": "myhost.example.com",
			}},
			args: args{
				inputFile: "testdata/template/input-1.yaml",
			},
			want: []string{"Welcome to myhost.example.com"},
		},
		{
			name: "replace in secret",
			fields: fields{values: map[string]any{
				"password": "GeneralReplacementsTransformer",
				"fooToken": "foo-bar-baz-123",
			}},
			args: args{
				inputFile: "testdata/template/input-2.yaml",
			},
			want: []string{
				"PASSWORD: R2VuZXJhbFJlcGxhY2VtZW50c1RyYW5zZm9ybWVy",
				"FOO_TOKEN: Zm9vLWJhci1iYXotMTIz", // "foo-bar-baz-123" (replaced)
				"BAR_TOKEN: W1suZm9vVG9rZW5dXQ==", // "[[.fooToken]]"   (not replaced)
			},
		},
		{
			name: "use slim-sprig funcs",
			fields: fields{values: map[string]any{
				"password": "GeneralReplacementsTransformer",
			}},
			args: args{
				inputFile: "testdata/template/input-3.yaml",
			},
			want: []string{"<i>2022/03/25</i>", "<b>ABCABCABC</b>"},
		},
		{
			name: "default delims",
			fields: fields{values: map[string]any{
				"var": "abc",
			}},
			args: args{
				inputFile: "testdata/template/input-4.yaml",
				repl:      Replacement{},
			},
			want: []string{"a=abc", "b=[[.var]]", "c=((.var))"},
		},
		{
			name: "global delims",
			fields: fields{
				config: Config{Delimiters: &[2]string{"[[", "]]"}},
				values: map[string]any{
					"var": "abc",
				}},
			args: args{
				inputFile: "testdata/template/input-4.yaml",
			},
			want: []string{"a={{.var}}", "b=abc", "c=((.var))"},
		},
		{
			name: "global delims with local overwrite",
			fields: fields{
				config: Config{Delimiters: &[2]string{"[[", "]]"}},
				values: map[string]any{
					"var": "abc",
				}},
			args: args{
				inputFile: "testdata/template/input-4.yaml",
				repl:      Replacement{Delimiters: &[2]string{"((", "))"}},
			},
			want: []string{"a={{.var}}", "b=[[.var]]", "c=abc"},
		},
		{
			name: "replace in secret with local delims",
			fields: fields{values: map[string]any{
				"fooToken": "foo-bar-baz-123",
			}},
			args: args{
				inputFile: "testdata/template/input-2.yaml",
				repl:      Replacement{Delimiters: &[2]string{"[[", "]]"}},
			},
			want: []string{
				"FOO_TOKEN: e3suZm9vVG9rZW59fQ==", // "{{.fooToken}}"   (not replaced)
				"BAR_TOKEN: Zm9vLWJhci1iYXotMTIz", // "foo-bar-baz-123" (replaced)
			},
		},
		{
			name:    "broken secret date type error",
			args:    args{inputFile: "testdata/template/input-5.yaml"},
			wantErr: true,
		},
		{
			name:    "broken secret no data",
			args:    args{inputFile: "testdata/template/input-6.yaml"},
			wantErr: true,
		},
		{
			name:    "broken secret invalid base64",
			args:    args{inputFile: "testdata/template/input-7.yaml"},
			wantErr: true,
		},
		{
			name:    "broken template in secret",
			args:    args{inputFile: "testdata/template/input-8.yaml"},
			wantErr: true,
		},
		{
			name:    "execute error in template in secret",
			args:    args{inputFile: "testdata/template/input-9.yaml"},
			wantErr: true,
		},
		{
			name:    "parse error",
			args:    args{inputFile: "testdata/template/input-10.yaml"},
			wantErr: true,
		},
		{
			name:    "execute error",
			args:    args{inputFile: "testdata/template/input-11.yaml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Transformer{values: tt.fields.values, config: tt.fields.config}
			raw, err := os.ReadFile(tt.args.inputFile)
			if err != nil {
				t.Fatalf("ReadFile(): error = %v", err)
			}
			uu, _ := bytes2uu(raw)
			if err != nil {
				t.Fatalf("bytes2uu(): error = %v", err)
			}
			if err := tr.TemplateTransform(uu, tt.args.repl); (err != nil) != tt.wantErr {
				t.Fatalf("Transformer.TemplateTransform() error = %v, wantErr %v", err, tt.wantErr)
			}
			output, err := yaml.Marshal(uu.Object)
			if err != nil {
				t.Fatalf("Marshal(): error = %v", err)
			}
			result := string(output)
			failed := false
			for n, s := range tt.want {
				if !strings.Contains(result, s) {
					failed = true
					t.Errorf("Transformer.TemplateTransform() check func %v failed", n+1)
				}
			}
			if failed {
				t.Errorf("output:\n---\n%v---\n", result)
			}
		})
	}
}
