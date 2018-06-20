package cli

import (
	"net"
	"net/http"

	"cirello.io/snippetsd/pkg/errors"
	"cirello.io/snippetsd/pkg/ui/web"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) httpMode() cli.Command {
	return cli.Command{
		Name:        "http",
		Aliases:     []string{"serve"},
		Usage:       "http mode",
		Description: "starts snippets web server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "host",
				Value:  "localhost",
				EnvVar: "HOST",
			},
			cli.StringFlag{
				Name:   "port",
				Value:  "8080",
				EnvVar: "PORT",
			},
		},
		Action: func(ctx *cli.Context) error {
			addr := net.JoinHostPort(ctx.String("host"), ctx.String("port"))
			l, err := net.Listen("tcp", addr)
			if err != nil {
				return errors.E(ctx, err, "cannot bind port")
			}
			err = http.Serve(l, web.New(c.db))
			return errors.E(ctx, err)
		},
	}
}
