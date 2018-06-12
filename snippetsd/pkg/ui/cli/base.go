package cli // import "cirello.io/snippetsd/pkg/ui/cli"

import (
	"log"
	"os"
	"sort"
	"strings"

	"cirello.io/snippetsd/pkg/errors"
	"cirello.io/snippetsd/pkg/models/snippet"
	"cirello.io/snippetsd/pkg/models/user"
	"github.com/jmoiron/sqlx"
	"gopkg.in/urfave/cli.v1"
)

type commands struct {
	db *sqlx.DB
}

func (c *commands) bootstrap(ctx *cli.Context) error {
	if err := snippet.NewRepository(c.db).Bootstrap(); err != nil {
		return errors.E(ctx, err, "failed when bootstrapping snippets")
	}
	if err := user.NewRepository(c.db).Bootstrap(); err != nil {
		return errors.E(ctx, err, "failed when bootstrapping users")
	}

	return nil
}

// Run executes the application in CLI mode
func Run(db *sqlx.DB) {
	app := cli.NewApp()
	app.Name = "snippetsd"
	app.Usage = "snippets server"
	app.Version = "0.0.1"

	cmds := &commands{
		db: db,
	}
	app.Before = cmds.bootstrap
	app.Commands = []cli.Command{
		cmds.addUser(),
		cmds.httpMode(),
	}
	sort.Slice(app.Commands, func(i, j int) bool {
		return strings.Compare(app.Commands[i].Name, app.Commands[j].Name) < 0
	})
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
