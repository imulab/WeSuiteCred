package suite

import "github.com/uptrace/bun"

// Ticket models the suite_ticket table which stores the suite ticket issued periodically from WeChat.
type Ticket struct {
	bun.BaseModel `bun:"table:suite_ticket"`

	ID     int64  `bun:"id,pk,autoincrement"`
	Ticket string `bun:"ticket,notnull"`
}
