package grt

import (
	"fmt"
	"log"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Transformer struct {
	uus            []*unstructured.Unstructured
	config         Config
	values         map[string]any
	valuesModified bool
}

type optFunc func(*Transformer) error

func WithConfigFile(file string) optFunc {
	return func(t *Transformer) error {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		return t.configure(data)
	}
}

func WithConfigString(data string) optFunc {
	return func(t *Transformer) error {
		return t.configure([]byte(data))
	}
}

func WithConfigBytes(data []byte) optFunc {
	return func(t *Transformer) error {
		return t.configure(data)
	}
}

func New(opts ...optFunc) (Transformer, error) {
	t := Transformer{}
	for _, opt := range opts {
		if err := opt(&t); err != nil {
			return Transformer{}, err
		}
	}
	return t, nil
}

func bytes2uu(buf []byte) (*unstructured.Unstructured, error) {
	obj := map[string]any{}
	err := yaml.Unmarshal(buf, obj)
	if err != nil {
		dumpYaml(buf)
		return &unstructured.Unstructured{}, fmt.Errorf("unmarshal: %v", err)
	}
	return &unstructured.Unstructured{Object: obj}, nil
}

func (t *Transformer) RegisterRaw(buf []byte) error {
	uu, err := bytes2uu(buf)
	if err != nil {
		return err
	}
	t.Register(uu)
	return nil
}

func (t *Transformer) Register(uu *unstructured.Unstructured) {
	t.uus = append(t.uus, uu)
}

func dumpYaml(s []byte) {
	for n, line := range strings.Split(string(s), "\n") {
		log.Printf("yaml[%03d]> %s", n, line)
	}
}
