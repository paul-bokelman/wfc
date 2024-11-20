package main

import (
	"fmt"
	"os"
	wfc "wfc/lib"

	"gopkg.in/yaml.v3"
)

type ConstraintsSet struct {
	Name       string           `yaml:"name"`
	Resolution int              `yaml:"resolution"`
	Cells      wfc.Constraints  `yaml:"cells"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    data, err := os.ReadFile("assets/road/constraints.yaml")
    check(err)

	constraints := ConstraintsSet{}
    
	yaml_err := yaml.Unmarshal([]byte(data), &constraints)
	check(yaml_err)

	wave := wfc.WaveFunctionCollapse{
		SetName: constraints.Name,
		Constraints: constraints.Cells,
		Resolution: constraints.Resolution,
		Size: 20,
	}

	wave.Build()
	collapseError := wave.Collapse()

	if collapseError != nil {
		fmt.Printf("Failed to collapse: %v", collapseError.Error())
	}

	saveError := wave.Save()

	if saveError != nil {
		fmt.Printf("Failed to save: %v", saveError.Error())
	}

	fmt.Printf("Success")
}