package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// test in MySQL
var (
	driver = "mysql"
	dsn    = "root:toor@tcp(localhost:13306)/demo_test?parseTime=True"
)

func TestMySQLSkipLocked(t *testing.T) {
	t.Run("test for "+driver, func(t *testing.T) {
		assert := require.New(t)
		ctx := context.Background()

		// connect database
		db, err := sql.Open(driver, dsn)
		assert.NoError(err)
		defer func() {
			assert.NoError(db.Close())
		}()

		// verify connection
		assert.NoError(db.PingContext(ctx))

		// reset table
		_, err = db.ExecContext(ctx, `
			DROP TABLE IF EXISTS seats;
		`)
		assert.NoError(err)
		_, err = db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS seats (
			  seat_no INT PRIMARY KEY,
			  booked ENUM('YES', 'NO') DEFAULT 'NO'
			);
		`)
		assert.NoError(err)

		// insert 100 seats
		stmtIns, err := db.PrepareContext(ctx, `
			INSERT INTO seats (
				seat_no
			) VALUES (
				?
			);
		`)
		defer func() {
			assert.NoError(stmtIns.Close())
		}()
		for i := 1; i <= 100; i++ {
			_, err := stmtIns.ExecContext(ctx, i)
			assert.NoError(err)
		}

		t.Run("transaction for skip locked", func(t *testing.T) {
			assert := require.New(t)

			// someone check available seats
			tx, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx.Rollback())
			}()

			rows, err := tx.QueryContext(ctx, `
				SELECT seat_no 
				FROM seats
				WHERE seat_no BETWEEN ? AND ? 
				AND booked = 'NO'
				FOR UPDATE SKIP	LOCKED;
			`, 1, 2)
			assert.NoError(err)

			var seatNos []int
			for rows.Next() {
				var seatNo int
				assert.NoError(rows.Scan(&seatNo))
				seatNos = append(seatNos, seatNo)
			}
			assert.Equal([]int{1, 2}, seatNos)

			// another one check available seats
			tx2, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx2.Rollback())
			}()

			timeoutCtx, _ := context.WithTimeout(ctx, time.Second)
			rows, err = tx2.QueryContext(timeoutCtx, `
				SELECT seat_no
				FROM seats
				WHERE seat_no BETWEEN ? AND ?
				AND booked = 'NO'
				FOR UPDATE SKIP LOCKED;
			`, 3, 5)
			assert.NoError(err)

			seatNos = []int{}
			for rows.Next() {
				var seatNo int
				assert.NoError(rows.Scan(&seatNo))
				seatNos = append(seatNos, seatNo)
			}
			// failed here, see ISSUES.md
			assert.Equal([]int{3, 4, 5}, seatNos)
		})

	})
}
