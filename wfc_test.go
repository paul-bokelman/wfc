package wfc

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type ConstraintsSet struct {
	Name       string           `yaml:"name"`
	Resolution int              `yaml:"resolution"`
	Cells      Constraints  `yaml:"cells"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func TestWFC(t *testing.T) {
    data, err := os.ReadFile("assets/road/constraints.yaml")
    check(err)

	constraints := ConstraintsSet{}
    
	yaml_err := yaml.Unmarshal([]byte(data), &constraints)
	check(yaml_err)

	wave := WaveFunctionCollapse{
		SetName: constraints.Name,
		Constraints: constraints.Cells,
		Resolution: constraints.Resolution,
		Size: 10,
	}

	wave.Build()
	err = wave.Collapse()

	if err != nil {
		t.Errorf("Failed to collapse: %v", err.Error())
	}

	err = wave.Save()

	if err != nil {
		t.Errorf("Failed to save: %v", err.Error())
	}

	fmt.Println("Success")
}