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
		Size: 10,
	}

	wave.Build()
	err = wave.Collapse()

	if err != nil {
		fmt.Printf("Failed to collapse: %v", err.Error())
	}

	err = wave.Save()

	if err != nil {
		fmt.Printf("Failed to save: %v", err.Error())
	}

	fmt.Printf("Success")
}