package dhcpmain

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/caddyserver/caddy"
	"github.com/nextdhcp/nextdhcp/core/lease"
)

var (
	conf       string
	serverType = "dhcpv4"
	listLeases bool
)

func init() {
	caddy.DefaultConfigFile = "Dhcpfile"
	caddy.Quiet = false

	flag.StringVar(&conf, "conf", "", "Dhcpfile to load (default \""+caddy.DefaultConfigFile+"\")")
	flag.BoolVar(&listLeases, "leases", false, "list leases from the database file(s) listed as arguments")

	caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(configLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))

	caddy.AppName = "NextDHCP"
	caddy.AppVersion = "v0.1.0"
}

// Run start NextDHCP and blocks until the server stopped
func Run() {
	flag.Parse()

	// -leases does not start the server
	if listLeases {
		ListLeases()
		return // unreached
	}

	caddy.TrapSignals()

	dhcpfile, err := caddy.LoadCaddyfile(serverType)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := caddy.Start(dhcpfile)
	if err != nil {
		log.Fatal(err)
	}

	instance.Wait()
}

func ListLeases() {
	for _, f := range flag.Args() {
		db := lease.MustOpen("bolt", map[string][]string{"file": {f}})
		leases, err := db.Leases(context.Background())
		if err != nil {
			panic(err)
		}
		for _, l := range leases {
			fmt.Printf("%s\t%s\t%s\n", l.Address, l.HwAddr, l.Hostname)
		}
	}
}

func configLoader(serverType string) (caddy.Input, error) {
	if conf == "" {
		return nil, nil
	}

	if conf == "stdin" || conf == "-" {
		return caddy.CaddyfileFromPipe(os.Stdin, serverType)
	}

	file, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}

	return caddy.CaddyfileInput{
		Contents:       file,
		Filepath:       conf,
		ServerTypeName: serverType,
	}, nil
}

func defaultLoader(serverType string) (caddy.Input, error) {
	conf = caddy.DefaultConfigFile
	return configLoader(serverType)
}
