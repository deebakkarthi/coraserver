package db

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetFreeClass(slot int, day string) []string {
	var classroom []string

	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT class_id FROM static WHERE
        slot_id = ? AND
        day = ? AND
        subject_id = "FREE"`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(slot, day)
	// Process the query results
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		classroom = append(classroom, tmp)
	}
	return classroom
}

func GetFreeSlot(class string, day string) []int {
	var slot []int
	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT slot_id FROM static WHERE
        class_id = ? AND
        day = ? AND
        subject_id = "FREE"`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(class, day)
	// Process the query results
	for rows.Next() {
		var tmp int
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		slot = append(slot, tmp)
	}
	return slot

}

func GetTimetableByDay(class string, day string) []string {
	var subject []string
	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT subject_id FROM static WHERE
        class_id = ? AND
        day = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(class, day)
	// Process the query results
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		subject = append(subject, tmp)
	}
	return subject
}

func GetAllSlot() []int {
	var slot []int
	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT id FROM slot;`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	for rows.Next() {
		var tmp int
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		slot = append(slot, tmp)
	}
	return slot
}

func GetAllClass() []string {
	var class []string
	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT UNIQUE class_id FROM static;`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		class = append(class, tmp)
	}
	return class
}

func GetAllSubject() []string {
	var subject []string
	db, err := sql.Open("mysql", "cora:@/cora")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT id FROM subject WHERE id!="FREE";`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		subject = append(subject, tmp)
	}
	return subject
}

func Booking(class string, date time.Time, slot int, faculty string, subject string) (int64, error) {
	db, err := sql.Open("mysql", "cora:@/cora?parseTime=true")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer db.Close()

	day := strings.ToUpper(date.Weekday().String()[:3])

	/*
	   INSERT INTO dynamic SELECT "A104", "2023-06-13", 1,
	   "cb.en.u4cse20613@cb.students.amrita.edu", "19CSE311" FROM dual WHERE
	   (SELECT subject_id FROM static WHERE class_id="A104" AND slot_id=1 AND
	   day="TUE")="FREE";
	*/
	stmt, err := db.Prepare(`INSERT INTO dynamic SELECT ?, ?, ?, ?, ? FROM
    dual WHERE (SELECT subject_id FROM static WHERE class_id = ? AND day = ?
    AND slot_id = ?)="FREE";`)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(class, date, slot, faculty, subject, class, day, slot)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}
