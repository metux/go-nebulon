package main

// simple Nebulon server program

import (
	"flag"
	"fmt"
	"os"

	"github.com/metux/go-nebulon/webapi/servers"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Perseus Nebulon server\n\n")
		flag.PrintDefaults()
		os.Exit(127)
	}

	var flag_conffile string
	var flag_help bool
	var flag_serverid string

	flag.BoolVar(&flag_help, "help", false, "help")
	flag.StringVar(&flag_conffile, "conf", "", "config file name")
	flag.StringVar(&flag_serverid, "server", "", "server config section")
	flag.Parse()

	if flag_help {
		flag.Usage()
	}

	if flag_conffile == "" {
		fmt.Fprintf(os.Stderr, "%s: missing -conf flag\n\n", os.Args[0])
		flag.Usage()
	}

	if flag_serverid == "" {
		fmt.Fprintf(os.Stderr, "%s: missing -server flag\n\n", os.Args[0])
		flag.Usage()
	}

	server, err := servers.BootServer(flag_conffile, flag_serverid)
	if err != nil {
		panic(err)
	}
	server.Serve()
}
