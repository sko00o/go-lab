package database

import (
	"database/sql"
)

type Seat struct {
	SeatNo int64
	Booked bool
}

func Create(db sql.DB, s *Seat) (err error) {
	_, err = db.Exec(`
		INSERT INTO seats (seat_no, booked) VALUES ($1,$2)
		`, s.SeatNo, s.Booked)
	return
}

func GetBySeatNo(db sql.DB, seatNo int64) (s Seat, err error) {
	db.QueryRow(`
		SELECT seat_no, booked FROM seats WHERE seat_no = $1
		`, seatNo).Scan(&s.SeatNo, &s.Booked)
	return
}

func GetSeats(db sql.DB, startNo, endNo int64, booked bool) ([]Seat, error) {
	rows, err := db.Query(`
		SELECT seat_no, booked FROM seats WHERE seat_no BETWEEN $1 AND $2 AND booked = $3
		`, startNo, endNo, booked)
	if err != nil {
		return nil, err
	}

	var ss []Seat
	for rows.Next() {
		var s Seat
		if err := rows.Scan(&s.SeatNo, &s.Booked); err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, nil
}
