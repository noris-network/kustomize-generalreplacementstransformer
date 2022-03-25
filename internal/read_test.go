package grt

import (
	"io"
	"strings"
	"testing"
)

func TestTransformer_ReadStream(t *testing.T) {
	type fields struct {
		values map[string]any
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantObjs int
	}{
		{
			name: "empty input",
			args: args{r: strings.NewReader("")},
		},
		{
			name:     "one object",
			args:     args{r: strings.NewReader("a: 1")},
			wantObjs: 1,
		},
		{
			name:     "one object with extra seperators",
			args:     args{r: strings.NewReader("---\na: 1\n---")},
			wantObjs: 1,
		},
		{
			name:     "one object with extra seperators and empty lines",
			args:     args{r: strings.NewReader("\n\n---\na: 1\n---\n\n")},
			wantObjs: 1,
		},
		{
			name:     "two objects",
			args:     args{r: strings.NewReader("a: 1\n---\nb: 2")},
			wantObjs: 2,
		},
		{
			name: "very long lines",
			args: args{r: strings.NewReader("" +
				"a: " + strings.Repeat("1", maxLineLength-4) +
				"\n---\n" +
				"b: " + strings.Repeat("2", maxLineLength-4),
			)},
			wantObjs: 2,
		},
		{
			name:    "no yaml",
			args:    args{r: strings.NewReader("a=1")},
			wantErr: true,
		},
		{
			name:    "two objects no yaml",
			args:    args{r: strings.NewReader("a=1\n---\nb=2")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{
				values: tt.fields.values,
			}
			if err := tr.ReadStream(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Transformer.ReadStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tr.uus) != tt.wantObjs {
				t.Errorf("Transformer.ReadStream() objects = %v, want %v", len(tr.uus), tt.wantObjs)
			}
		})
	}
}
