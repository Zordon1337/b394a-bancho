package main

import (
	"retsu/cho"
	frontend "retsu/frontend/src"
	"retsu/shared/db"
)

func main() {
	db.InitDatabase()
	go cho.Bancho()
	frontend.Frontend()
}
