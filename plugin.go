package main

import (
	"fmt"
)

// Plugin defines the Device farm plugin parameters.
type Plugin struct {
	Key          string
	Secret       string
	Region       string
	YamlVerified bool

	TestProject string
	RunName     string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	fmt.Println("Begin ")

	fmt.Println("End")

	return nil
}
