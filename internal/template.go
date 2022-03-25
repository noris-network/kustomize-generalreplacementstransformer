package grt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (t Transformer) TemplateTransform(uu *unstructured.Unstructured, repl Replacement) error {

	tmpl := template.New("").Funcs(sprig.TxtFuncMap())

	if t.config.Delimiters != nil {
		tmpl = tmpl.Delims(t.config.Delimiters[0], t.config.Delimiters[1])
	}
	if repl.Delimiters != nil {
		tmpl = tmpl.Delims(repl.Delimiters[0], repl.Delimiters[1])
	}

	if uu.GetKind() == "Secret" {
		data, ok, err := unstructured.NestedStringMap(uu.Object, "data")
		if err != nil {
			return fmt.Errorf("nestedStringMap: %v", err)
		}
		if !ok {
			return fmt.Errorf("%v/%v: %q not found\n", uu.GetKind(), uu.GetName(), "data")
		}
		for k, v := range data {
			plain, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return fmt.Errorf("DecodeString(): %v", err)
			}
			tmpl, err = tmpl.Parse(string(plain))
			if err != nil {
				return fmt.Errorf("parse: %v", err)
			}
			bb := bytes.Buffer{}
			err = tmpl.Execute(&bb, t.values)
			if err != nil {
				return fmt.Errorf("execute: %v", err)
			}
			data[k] = base64.StdEncoding.EncodeToString(bb.Bytes())
		}
		uu.Object["data"] = data
		return nil
	} else {
		data, err := yaml.Marshal(uu.Object)
		if err != nil {
			return fmt.Errorf("Marshal: %v", err)
		}
		tmpl, err = tmpl.Parse(string(data))
		if err != nil {
			return fmt.Errorf("parse: %v", err)
		}
		bb := bytes.Buffer{}
		err = tmpl.Execute(&bb, t.values)
		if err != nil {
			return fmt.Errorf("execute: %v", err)
		}
		uuout, err := bytes2uu(bb.Bytes())
		if err != nil {
			return fmt.Errorf("bytes2uu: %v", err)
		}
		uu.Object = uuout.Object
		return nil
	}
}
