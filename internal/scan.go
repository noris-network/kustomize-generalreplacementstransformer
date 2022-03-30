package grt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (t *Transformer) ScanForValues() (err error) {

	t.values = t.config.Values

	defer func() {
		if err != nil {
			t.values = nil
		}
	}()

	for _, sel := range t.config.SelectValues {
		if sel.Name == "" {
			return errors.New("select values: name must not be empty")
		}
		t.values[sel.Name] = sel.Default
		found := false
		for _, uu := range t.uus {
			uuKind := uu.GetKind()
			uuName := uu.GetName()
			if strings.ContainsAny(strings.TrimRight(sel.Resource.Name, "*"), "*") {
				return errors.New("select values: no prefix or infix wildcards allowed")
			}

			nameMatches, err := nameMatch(uuName, sel.Resource.Name)
			if err != nil {
				return fmt.Errorf("select value: %v", err)
			}

			if sel.Resource.Kind == "" || sel.Resource.Name == "" || sel.Resource.FieldPath == "" {
				return fmt.Errorf("select value %q: resource kind, name and filepath must not be empty", sel.Name)
			}
			if uuKind != sel.Resource.Kind || !nameMatches {
				continue
			}
			if (uuKind == "ConfigMap" || uuKind == "Secret") && sel.Resource.FieldPath == "data.*" {
				value, ok, err := unstructured.NestedStringMap(uu.Object, "data")
				if err != nil {
					return fmt.Errorf("nestedStringMap: %v", err)
				}
				if !ok {
					return fmt.Errorf("nestedStringMap: %v/%v: data not found\n", uuKind, uuName)
				}
				if uuKind == "Secret" {
					for k, v := range value {
						b, _ := base64.StdEncoding.DecodeString(v)
						value[k] = string(b)
					}
				}
				if sel.Splat {
					for k, v := range value {
						fmt.Printf("## GeneralReplacementsTransformer: %v = %q\n", k, v)
						t.values[k] = v
					}
					delete(t.values, sel.Name)
				} else {
					for k, v := range value {
						fmt.Printf("## GeneralReplacementsTransformer: %v.%v = %q\n", sel.Name, k, v)
					}
					t.values[sel.Name] = value
				}
				found = true
			} else {
				uuYaml, _ := yaml.Marshal(uu.Object)
				node, _ := yaml.Parse(string(uuYaml))
				node, err := node.Pipe(yaml.Lookup(strings.Split(sel.Resource.FieldPath, ".")...))
				if err != nil {
					return fmt.Errorf("pipe: %v", err)
				}
				if node.IsNilOrEmpty() {
					fmt.Printf("## GeneralReplacementsTransformer: %v/%v: %q not found\n", uuKind, uuName, sel.Resource.FieldPath)
					continue
				}
				if len(node.Content()) > 0 {
					return fmt.Errorf("%v/%v: %q: returns %v nodes\n", uuKind, uuName, sel.Resource.FieldPath, len(node.Content()))
				}
				value, err := node.String()
				if err != nil {
					return fmt.Errorf("node2string: %v", err)
				}
				value = strings.TrimRight(value, "\"\n")
				value = strings.TrimLeft(value, `"`)
				if uuKind == "Secret" {
					plain, _ := base64.StdEncoding.DecodeString(value)
					t.values[sel.Name] = string(plain)
				} else {
					t.values[sel.Name] = value
				}
				fmt.Printf("## GeneralReplacementsTransformer: %v = %q\n", sel.Name, t.values[sel.Name])
				found = true
			}
		}
		if !found {
			fmt.Printf("## GeneralReplacementsTransformer: select value %q not found\n", sel.Name)
		}
	}

	return nil
}
