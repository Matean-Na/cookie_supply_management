package main

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/core/connect"
	"cookie_supply_management/core/database"
	"cookie_supply_management/core/server"
	"flag"
	"fmt"

	"strings"
)

type flags []string

func (f *flags) String() string {
	return strings.Join(*f, ", ")
}

func (f *flags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var sources flags

	flag.Var(&sources, "s", "Specify the source for seeding")
	flag.Parse()

	conf, err := config.Load()
	if err != nil {
		fmt.Printf("err config.Load() %s\n", err)
		return
	}

	dbase, err := database.Connect(conf.Database)
	if err != nil {
		fmt.Printf("err db.Connect() %s\n", err)
		return
	}
	connect.DB = dbase

	server.Runner(sources)
}
