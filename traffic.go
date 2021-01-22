package main

import "github.com/EngoEngine/engo"

type myScene struct{}

// Type identifies the scene type.
func (*myScene) Type() string {
	return "myGame"
}

// Preload is called before loading assets,
// allowing them to be registered / queued.
func (*myScene) Preload() {}

// Setup is called before the main loop starts,
// allowing entities and systems to be added.
func (*myScene) Setup(engo.Updater) {}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  800,
		Height: 600,
	}
	engo.Run(opts, &myScene{})
}
