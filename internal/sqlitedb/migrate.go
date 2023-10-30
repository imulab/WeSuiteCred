package sqlitedb

import (
	"context"
	"embed"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"time"
)

var (
	//go:embed scripts/*.sql
	migrationFiles embed.FS
)

// Migrate runs database migrations with the scripts embedded in the project.
func Migrate(db *bun.DB, logger *zerolog.Logger) error {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migrationFiles); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := migrator.Init(ctx); err != nil {
		return err
	}

	g, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	logger.Info().
		Str("group", g.String()).
		Strs("details", lo.Map(g.Migrations, func(m migrate.Migration, _ int) string { return m.String() })).
		Msg("Completed migration.")

	return nil
}
