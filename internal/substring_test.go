package grt

import (
	"os"
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestTransformer_ReplaceTransform(t *testing.T) {
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
				"Welcome":  "<b>Welcome</b>",
			}},
			args: args{
				inputFile: "testdata/replace/input-1.yaml",
			},
			want: []string{"<b>Welcome</b> to myhost.example.com"},
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
			if err := tr.SubstringTransform(uu, tt.args.repl); (err != nil) != tt.wantErr {
				t.Fatalf("Transformer.SubstringTransform() error = %v, wantErr %v", err, tt.wantErr)
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
					t.Errorf("Transformer.SubstringTransform() check func %v failed", n+1)
				}
			}
			if failed {
				t.Errorf("output:\n---\n%v---\n", result)
			}
		})
	}
}
