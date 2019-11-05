package database

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
)

// test in MySQL
var (
	driver = "mysql"
	dsn    = "root:toor@tcp(localhost:13306)/demo_test?parseTime=True"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	exitCode := m.Run()

	if err = db.Close(); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func TestMySQLSkipLocked(t *testing.T) {
	assert := require.New(t)
	ctx := context.Background()
	var err error

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
			WHERE seat_no IN (? ,?)
			AND booked = 'NO'
			FOR UPDATE;
		`, 3, 4)
		assert.NoError(err)

		var seatNos []int
		for rows.Next() {
			var seatNo int
			assert.NoError(rows.Scan(&seatNo))
			seatNos = append(seatNos, seatNo)
		}
		assert.Equal([]int{3, 4}, seatNos)

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
		`, 1, 5)
		assert.NoError(err)

		seatNos = []int{}
		for rows.Next() {
			var seatNo int
			assert.NoError(rows.Scan(&seatNo))
			seatNos = append(seatNos, seatNo)
		}

		assert.Equal([]int{1, 2, 5}, seatNos)
	})
}

func TestMultiTableSkipLocked(t *testing.T) {
	assert := require.New(t)
	ctx := context.Background()
	var err error

	// reset tables
	_, err = db.ExecContext(ctx, `
		DROP TABLE IF EXISTS seat_rows;
	`)
	assert.NoError(err)
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS seat_rows ( 
			row_no INT PRIMARY KEY, 
			cost DECIMAL 
		);
	`)
	assert.NoError(err)
	_, err = db.ExecContext(ctx, `
		DROP TABLE IF EXISTS seats;
	`)
	assert.NoError(err)
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS seats (
			seat_no INT NOT NULL,
			row_no INT NOT NULL,
			booked ENUM('YES', 'NO') DEFAULT 'NO',
			PRIMARY KEY (seat_no, row_no)
		);
	`)
	assert.NoError(err)

	// insert 20 stadium rows with 100 seats/row
	stmtIns, err := db.PrepareContext(ctx, `
		INSERT INTO seats (
			seat_no, row_no
		) VALUES (
			?, ?
		);
	`)
	defer func() {
		assert.NoError(stmtIns.Close())
	}()
	for i := 1; i <= 100; i++ {
		for j := 1; j <= 20; j++ {
			_, err := stmtIns.ExecContext(ctx, i, j)
			assert.NoError(err)
		}
	}

	// add pricing information for rows
	stmtPrice, err := db.PrepareContext(ctx, `
		INSERT INTO seat_rows (
			row_no, cost
		) VALUES (
			?, ?
		);
	`)
	defer func() {
		assert.NoError(stmtPrice.Close())
	}()
	for n := 1; n <= 20; n++ {
		_, err := stmtPrice.ExecContext(ctx, n, 100-n*2)
		assert.NoError(err)
	}

	t.Run("lock rows only in seats", func(t *testing.T) {
		assert := require.New(t)

		// lock seats
		tx, err := db.Begin()
		assert.NoError(err)
		defer func() {
			assert.NoError(tx.Rollback())
		}()

		rows, err := tx.QueryContext(ctx, `
			SELECT seat_no, row_no, cost
			FROM seats s
				JOIN seat_rows sr USING (row_no)
			WHERE seat_no IN (3, 4)
			  	AND sr.row_no IN (5, 6)
			FOR UPDATE OF s SKIP LOCKED;`)
		assert.NoError(err)

		var results [][3]int
		for rows.Next() {
			var res [3]int
			assert.NoError(rows.Scan(&res[0], &res[1], &res[2]))
			results = append(results, res)
		}

		assert.Equal([][3]int{
			{3, 5, 100 - 5*2},
			{3, 6, 100 - 6*2},
			{4, 5, 100 - 5*2},
			{4, 6, 100 - 6*2},
		}, results)

		t.Run("check seats", func(t *testing.T) {
			assert := require.New(t)

			tx, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx.Rollback())
			}()

			rows, err = tx.QueryContext(ctx, `
			SELECT seat_no, row_no, cost
			FROM seats s
				JOIN seat_rows sr USING (row_no)
			WHERE seat_no IN (2,3,4,5)
			  	AND sr.row_no IN (5, 6)
			FOR UPDATE OF s SKIP LOCKED;`)
			assert.NoError(err)

			results = [][3]int{}
			for rows.Next() {
				var res [3]int
				assert.NoError(rows.Scan(&res[0], &res[1], &res[2]))
				results = append(results, res)
			}

			assert.Equal([][3]int{
				{2, 5, 100 - 5*2},
				{2, 6, 100 - 6*2},
				{5, 5, 100 - 5*2},
				{5, 6, 100 - 6*2},
			}, results)
		})

		t.Run("check rows", func(t *testing.T) {
			assert := require.New(t)

			tx, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx.Rollback())
			}()

			rows, err = tx.QueryContext(ctx, `
			SELECT seat_no, row_no, cost
			FROM seats s
				JOIN seat_rows sr USING (row_no)
			WHERE seat_no IN (3,4)
			  	AND sr.row_no IN (4,5,6,7)
			FOR UPDATE OF s SKIP LOCKED;`)
			assert.NoError(err)

			results = [][3]int{}
			for rows.Next() {
				var res [3]int
				assert.NoError(rows.Scan(&res[0], &res[1], &res[2]))
				results = append(results, res)
			}

			assert.Equal([][3]int{
				{3, 4, 100 - 4*2},
				{3, 7, 100 - 7*2},
				{4, 4, 100 - 4*2},
				{4, 7, 100 - 7*2},
			}, results)
		})
	})

	t.Run("lock only the tables you want", func(t *testing.T) {
		assert := require.New(t)

		// someone lock 10 rows
		tx, err := db.Begin()
		assert.NoError(err)
		defer func() {
			assert.NoError(tx.Rollback())
		}()

		_, err = tx.QueryContext(ctx, `
			SELECT *
			FROM seat_rows
			WHERE row_no >= 10
			FOR UPDATE;`)
		assert.NoError(err)

		t.Run("check available rows", func(t *testing.T) {
			assert := require.New(t)

			tx, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx.Rollback())
			}()

			rows, err := tx.QueryContext(ctx, `
			SELECT seat_no, row_no, cost
			FROM seats
				JOIN seat_rows USING (row_no)
			WHERE seat_no IN (3, 4)
			  	AND seat_rows.row_no IN (5, 6)
			  	AND booked = 'NO'
			FOR UPDATE OF seats SKIP LOCKED
			FOR SHARE OF seat_rows;`)
			assert.NoError(err)

			var results [][3]int
			for rows.Next() {
				var res [3]int
				assert.NoError(rows.Scan(&res[0], &res[1], &res[2]))
				results = append(results, res)
			}

			assert.Equal([][3]int{
				{3, 5, 100 - 5*2},
				{3, 6, 100 - 6*2},
				{4, 5, 100 - 5*2},
				{4, 6, 100 - 6*2},
			}, results)
		})

		t.Run("attempt to acquire lock, fail immediately if not possible", func(t *testing.T) {
			assert := require.New(t)

			tx, err := db.Begin()
			assert.NoError(err)
			defer func() {
				assert.NoError(tx.Rollback())
			}()

			_, err = tx.QueryContext(ctx, `
			SELECT seat_no, row_no, cost
			FROM seats
				JOIN seat_rows USING (row_no)
			WHERE seat_no IN (3, 4)
			  	AND seat_rows.row_no IN (12)
			  	AND booked = 'NO'
			FOR UPDATE OF seats SKIP LOCKED
			FOR SHARE OF seat_rows NOWAIT;`)
			assert.Error(err)
			assert.Equal(err.(*mysql.MySQLError).Number, uint16(3572))
		})
	})
}
