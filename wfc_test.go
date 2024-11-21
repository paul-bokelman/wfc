package wfc

import (
	"fmt"
	"testing"
)

func TestWFC(t *testing.T) {
	wave := WaveFunctionCollapse{}
	err := wave.Initialize(WaveFunctionCollapseConfig{AssetsPath: "assets/road", OutputPath: "compiled.png", GridSize: 10})

	if err != nil {
		t.Errorf("Failed to initialize wfc")
	}

	err =  wave.Collapse()

	if err != nil {
		t.Errorf("Failed to collapse: %v", err.Error())
	}

	err = wave.Save()

	if err != nil {
		t.Errorf("Failed to save: %v", err.Error())
	}

	fmt.Println("Success")
}