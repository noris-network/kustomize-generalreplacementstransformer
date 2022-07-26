package grt

import (
	"bytes"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (t Transformer) SubstringTransform(uu *unstructured.Unstructured, repl Replacement) error {

	replacements := []string{}
	for k, v := range t.values {
		value, ok := v.(string)
		if ok {
			replacements = append(replacements, k, value)
		}
	}

	replacer := strings.NewReplacer(replacements...)

	data, err := yaml.Marshal(uu.Object)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	bb := bytes.Buffer{}
	_, err = replacer.WriteString(&bb, string(data))
	if err != nil {
		return fmt.Errorf("replace: %v", err)
	}
	uuout, err := bytes2uu(bb.Bytes())
	if err != nil {
		return fmt.Errorf("bytes2uu: %v", err)
	}
	uu.Object = uuout.Object
	return nil
}
