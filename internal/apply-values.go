package grt

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const maxDepth = 4

func (t *Transformer) ApplyValuesToValues() (err error) {
	depth := 0
	for {
		t.valuesModified = false
		if err := t.walk(&t.values); err != nil {
			return fmt.Errorf("apply values: %v", err)
		}
		if !t.valuesModified {
			break
		}
		depth++
		if depth > maxDepth {
			return fmt.Errorf("apply values: max depth (%v) reached", maxDepth)
		}
	}

	return nil
}

func (t *Transformer) walk(value *map[string]any) error {
	for k, v := range *value {
		if kv, ok := v.(map[string]any); ok {
			t.walk(&kv)
			continue
		}
		newValue, changed, err := t.replace(v)
		if err != nil {
			return fmt.Errorf("walk: key=%q: %v", k, err)
		}
		if changed {
			t.valuesModified = true
			(*value)[k] = newValue
		}
	}
	return nil
}

func (t *Transformer) replace(value any) (string, bool, error) {
	valstr, ok := value.(string)
	if !ok {
		return "", false, nil
	}
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(valstr)
	if err != nil {
		return "", false, fmt.Errorf("parse: %v", err)
	}
	bb := bytes.Buffer{}
	err = tmpl.Execute(&bb, t.values)
	if err != nil {
		return "", false, fmt.Errorf("exec: %v", err)
	}
	newValue := bb.String()
	if strings.Contains(newValue, "<no value>") {
		return "", false, fmt.Errorf("no value: %q", valstr)
	}
	return newValue, newValue != value, nil
}
