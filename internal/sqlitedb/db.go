package sqlitedb

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

const (
	defaultDsn = "file:/var/WeSuiteCred/main.db?cache=shared&mode=rwc&_journal_mode=WAL"
	memoryDsn  = "file::memory:?mode=rwc&_journal_mode=WAL"
)

func New() (*bun.DB, error) {
	return newDb(defaultDsn)
}

func NewMemory() (*bun.DB, error) {
	db, err := newDb(memoryDsn)
	if err != nil {
		return nil, err
	}

	db.DB.SetMaxIdleConns(1000)
	db.DB.SetConnMaxLifetime(0)

	return db, nil
}

func newDb(dsn string) (*bun.DB, error) {
	sqlDb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlDb, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.FromEnv("WSC_DEBUG"),
	))

	return db, nil
}
