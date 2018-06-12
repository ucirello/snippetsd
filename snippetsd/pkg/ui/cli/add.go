package cli

import (
	"fmt"

	"cirello.io/snippetsd/pkg/errors"
	"cirello.io/snippetsd/pkg/models/user"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) addUser() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add a user",
		Action: func(ctx *cli.Context) error {
			u, err := user.NewFromEmail(ctx.Args().First())
			if err != nil {
				return errors.E(ctx, err, "cannot create user from email")
			}
			if _, err := user.Add(c.db, u); err != nil {
				return errors.E(ctx, err, "cannot store the new user")
			}

			fmt.Println(u, "added")
			return nil
		},
	}
}
