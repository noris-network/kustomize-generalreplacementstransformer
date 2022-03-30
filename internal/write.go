package grt

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (t Transformer) WriteStream(w io.Writer) error {

	for i, uu := range t.uus {

		if i > 0 {
			fmt.Fprintln(w, "---")
		}

		uuKind := uu.GetKind()
		uuName := uu.GetName()

		for _, repl := range t.config.Replacements {
			replKind := repl.Resource.Kind
			replName := repl.Resource.Name
			nameMatches, err := nameMatch(uuName, replName)
			if err != nil {
				return fmt.Errorf("replacements: %v", err)
			}
			if (replKind == "" || replKind == uuKind) && (replName == "" || nameMatches) {
				switch repl.Type {
				case "template":
					err := t.TemplateTransform(uu, repl)
					if err != nil {
						return fmt.Errorf("transformation failed: %v", err)
					}
				default:
					return fmt.Errorf("replacement type %q not defined", repl.Type)
				}
			}
		}

		data, err := yaml.Marshal(uu.Object)
		if err != nil {
			return err
		}

		fmt.Fprint(w, string(data))
	}

	return nil
}

func nameMatch(name, wildcard string) (bool, error) {
	if name == wildcard {
		return true, nil
	}
	if strings.ContainsAny(strings.TrimRight(wildcard, "*"), "*") {
		return false, errors.New("no prefix or infix wildcards allowed")
	}
	if strings.HasSuffix(wildcard, "*") &&
		strings.HasPrefix(name, strings.TrimRight(wildcard, "*")) {
		return true, nil
	}
	return false, nil
}
