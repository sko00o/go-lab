package database

import (
	"time"
)

func TimeTruncate(p *People) {
	p.CreatedAt = p.CreatedAt.UTC().Truncate(time.Second)
	p.UpdatedAt = p.UpdatedAt.UTC().Truncate(time.Second)
}
