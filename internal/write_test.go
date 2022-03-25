package grt

import (
	"bytes"
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestTransformer_WriteStream(t *testing.T) {
	type fields struct {
		uus    []*unstructured.Unstructured
		config Config
		values map[string]any
	}
	type args struct {
		inputFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantF   string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{inputFile: "testdata/write/input-1.yaml"},
			fields: fields{
				values: map[string]any{"a": 1, "b": 2},
				config: Config{Replacements: []Replacement{{Type: "template"}}},
			},
			wantF: "testdata/write/output-1.yaml",
		},
		{
			name: "unknown type",
			args: args{inputFile: "testdata/write/input-1.yaml"},
			fields: fields{
				values: map[string]any{"a": 1, "b": 2},
				config: Config{Replacements: []Replacement{{Type: "no-such-type"}}},
			},
			wantF:   "testdata/write/output-1.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Transformer{
				values: tt.fields.values,
				config: tt.fields.config,
			}
			r, err := os.Open(tt.args.inputFile)
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}
			defer r.Close()
			tr.ReadStream(r)
			w := &bytes.Buffer{}
			if err := tr.WriteStream(w); (err != nil) != tt.wantErr {
				t.Fatalf("Transformer.WriteStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			want, err := os.ReadFile(tt.wantF)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			if gotW := w.String(); gotW != string(want) {
				t.Errorf("Transformer.WriteStream()\ngot:\n%v,\nwant:\n%v", gotW, string(want))
			}
		})
	}
}
