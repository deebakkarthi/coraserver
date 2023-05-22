package db

import (
	"database/sql"
	"log"

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
