package database

import (
	"fmt"

	"github.com/caddyserver/caddy"
	"github.com/nextdhcp/nextdhcp/core/dhcpserver"
	"github.com/nextdhcp/nextdhcp/core/lease"
)

func init() {
	caddy.RegisterPlugin("database", caddy.Plugin{
		ServerType: "dhcpv4",
		Action:     parseDatabaseDirective,
	})
}

func parseDatabaseDirective(c *caddy.Controller) error {
	if !c.Next() {
		return c.ArgErr()
	}

	if !c.NextArg() {
		return c.ArgErr()
	}
	driverName := c.Val()

	var options = make(map[string][]string)
	remaining := c.RemainingArgs()
	if len(remaining) > 0 {
		options["__args__"] = remaining
	}

	for c.NextBlock() {
		options[c.Val()] = c.RemainingArgs()
	}

	if c.Next() {
		return c.ArgErr()
	}

	fmt.Println("opening database: ", driverName)
	db, err := lease.Open(driverName, options)
	if err != nil {
		return err
	}

	dhcpserver.GetConfig(c).Database = db

	return nil
}
