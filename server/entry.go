package main

import (
	"retsu/cho"
	frontend "retsu/frontend/src"
)

func main() {
	go cho.Bancho()
	frontend.Frontend()
}
