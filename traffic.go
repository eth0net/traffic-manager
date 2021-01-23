package main

import (
	"github.com/EngoEngine/engo"
	"github.com/raziel2244/traffic-manager/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          800,
		Height:         600,
		StandardInputs: true,
	}

	engo.Run(opts, &scenes.Scene{})
}
