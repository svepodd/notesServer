package main

import (
	"flag"
	"Notes/controller/httpserver"
	"Notes/gates/storage"
	"Notes/gates/storage/list"
	"Notes/gates/storage/mp"
)

func main() {
	var st storage.Storage
	var useList bool
	flag.BoolVar(&useList, "list", false, "A boolean flag")
	flag.Parse()

	if useList {
		st = list.NewList()
	} else {
		st = mp.NewMap()
	}
	hs := httpserver.NewHttpServer(":8080", st)
	hs.Start()
}
