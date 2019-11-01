package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type People struct {
	ID        int64      `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Name string `db:"name"`
	Age  int    `db:"age"`
}

func Create(db sqlx.Ext, p *People) (err error) {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	baseSQL := `
		INSERT INTO peoples (
			created_at, updated_at, deleted_at, name, age
		) VALUES (
			:created_at, :updated_at, :deleted_at, :name, :age
		)`

	if db.DriverName() == "postgres" {
		query, args, err := db.BindNamed(baseSQL+` RETURNING id`, p)
		if err != nil {
			return err
		}
		err = sqlx.Get(db, &p.ID, query, args...)
		return err
	}

	res, err := sqlx.NamedExec(db, baseSQL, p)
	if err != nil {
		return err
	}
	p.ID, err = res.LastInsertId()

	return
}

func QueryByName(db sqlx.Ext, name string) (p People, err error) {
	rows, err := sqlx.NamedQuery(db, `
		SELECT * 
		FROM peoples
		WHERE name = :name
		LIMIT 1`,
		People{Name: name})
	if err != nil {
		return p, err
	}

	for rows.Next() {
		if err = rows.StructScan(&p); err != nil {
			return
		}
		break
	}
	return
}

func QueryByAgeRange(db sqlx.Ext, min, max int) (ps []People, err error) {
	rows, err := sqlx.NamedQuery(db, `
		SELECT * 
		FROM peoples
		WHERE age BETWEEN :min AND :max
		`,
		map[string]interface{}{
			"min": min,
			"max": max,
		},
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p People
		if err = rows.StructScan(&p); err != nil {
			return
		}
		ps = append(ps, p)
	}
	return
}

func QueryByAgeRangeLock(db sqlx.Ext, min, max int) (ps []People, err error) {
	rows, err := sqlx.NamedQuery(db, `
		SELECT * 
		FROM peoples
		WHERE age BETWEEN :min AND :max
		FOR UPDATE SKIP LOCKED`,
		map[string]interface{}{
			"min": min,
			"max": max,
		},
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p People
		if err = rows.StructScan(&p); err != nil {
			return
		}
		ps = append(ps, p)
	}
	return
}

func QueryByAgeRangeLimitLock(db sqlx.Ext, min, max int, limit int) (ps []People, err error) {
	rows, err := sqlx.NamedQuery(db, `
		SELECT p.* 
		FROM peoples p
		WHERE age BETWEEN :min AND :max 
		LIMIT :limit
		FOR UPDATE OF p SKIP LOCKED`,
		map[string]interface{}{
			"min":   min,
			"max":   max,
			"limit": limit,
		},
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p People
		if err = rows.StructScan(&p); err != nil {
			return
		}
		ps = append(ps, p)
	}
	return
}
