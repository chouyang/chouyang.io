package main

import (
	"chouyang.io/src"
)

func main() {
	turbine := src.DefaultEngine()
	turbine.Migrate()

	turbine.Ignite()
}
