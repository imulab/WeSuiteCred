package sqlitedb

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMigrate(t *testing.T) {
	logger := zerolog.Nop()

	db, err := NewMemory()
	require.NoError(t, err)

	if !assert.NoError(t, Migrate(db, &logger)) {
		return
	}

	type Result struct {
		Name string `bun:"name"`
	}
	results := make([]Result, 0)

	if !assert.NoError(t, db.NewSelect().
		Column("name").
		Table("sqlite_master").
		Where("type = ?", "table").
		Scan(context.Background(), &results)) {
		return
	}

	names := lo.Map(results, func(r Result, _ int) string { return r.Name })

	assert.Contains(t, names, "suite_ticket")
	assert.Contains(t, names, "corp_authz")
}
