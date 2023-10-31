package show

import (
	"absurdlab.io/WeSuiteCred/internal/corp"
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

func Command() *cli.Command {
	conf := new(config)

	return &cli.Command{
		Name:   "show",
		Usage:  "Show corporation client credentials and permissions",
		Flags:  conf.flags(),
		Action: func(c *cli.Context) error { return runQuery(c.Context, conf) },
	}
}

func runQuery(ctx context.Context, conf *config) error {
	db, err := sqlitedb.New()
	if err != nil {
		return err
	}

	records := make([]corp.Authorization, 0)

	n, err := db.NewSelect().
		Model(&records).
		WhereOr("corp_name LIKE ?", "%"+conf.Query+"%").
		WhereOr("corp_id = ?", conf.Query).
		ScanAndCount(ctx)
	if err != nil {
		return err
	}

	if n == 0 {
		fmt.Println("No records found.")
		return nil
	}

	for i, record := range records {
		fmt.Printf("Record #%d\n", i+1)
		fmt.Printf("\tName: %s\n", record.CorpName)
		fmt.Printf("\tCorp id: %s\n", record.CorpID)
		fmt.Printf("\tCorp secret: %s\n", record.PermanentCode)
		fmt.Printf("\tPermissions: %s\n", strings.Join(record.Permissions.Unwrap().AppPermissions, " "))
		fmt.Println()
	}

	return nil
}
