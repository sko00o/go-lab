package transaction

import (
	"time"
)

func TimeTruncate(p *People) {
	p.CreatedAt = p.CreatedAt.UTC().Round(time.Second)
	p.UpdatedAt = p.UpdatedAt.UTC().Round(time.Second)
}
