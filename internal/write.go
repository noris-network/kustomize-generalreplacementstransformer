package grt

import (
	"fmt"
	"io"

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
			if (replKind == "" || replKind == uuKind) && (replName == "" || replName == uuName) {
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
