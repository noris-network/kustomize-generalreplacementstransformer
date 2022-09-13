package grt

import (
	"bytes"
	"os"
	"testing"
)

func TestTransformer_WriteStream(t *testing.T) {
	type fields struct {
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
		{
			name: "wildcard",
			args: args{inputFile: "testdata/write/input-2.yaml"},
			fields: fields{
				config: Config{
					Replacements: []Replacement{
						{Type: "template",
							Resource: Resource{
								Kind: "ConfigMap",
								Name: "mydbconfig-*",
							},
						},
					},
				},
				values: map[string]any{"c": "c"},
			},
			wantF: "testdata/write/output-2.yaml",
		},
		{
			name: "infix wildcard",
			args: args{inputFile: "testdata/write/input-2.yaml"},
			fields: fields{
				config: Config{
					Replacements: []Replacement{
						{Type: "template", Resource: Resource{Name: "my*config-*"}},
					},
				},
			},
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

func Test_prefixMatch(t *testing.T) {
	type args struct {
		name     string
		wildcard string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "no match",
			args: args{name: "foo", wildcard: "foo-*"},
			want: false,
		},
		{
			name: "match",
			args: args{name: "foo-bar", wildcard: "foo-*"},
			want: true,
		},
		{
			name: "exact match",
			args: args{name: "foo-bar", wildcard: "foo-bar"},
			want: true,
		},
		{
			name:    "infix error",
			args:    args{name: "foo-bar", wildcard: "fo*ar"},
			wantErr: true,
		},
		{
			name:    "prefix error",
			args:    args{name: "foo-bar", wildcard: "*bar"},
			wantErr: true,
		},
		{
			name:    "error",
			args:    args{name: "foo-bar", wildcard: "*bar*"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nameMatch(tt.args.name, tt.args.wildcard)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefixMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("prefixMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
