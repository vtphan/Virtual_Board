package models

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Course struct {
	Id        int       `json:"id"`
	Name      string    `json:"coursename"`
	CreatedAt time.Time `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
	UpdatedAt time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"UTC"`
}

func GetCourses() ([]Course, error) {
	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		logError("GetCourses: during Query", err)
		return nil, err
	}

	defer rows.Close()
	courses := make([]Course, 0)

	for rows.Next() {
		singleCourse := Course{}
		err = rows.Scan(&singleCourse.Id, &singleCourse.Name, &singleCourse.CreatedAt, &singleCourse.UpdatedAt)

		if err != nil {
			logError("GetCourses: during Scan", err)
			return nil, err
		}

		courses = append(courses, singleCourse)
	}
	err = rows.Err()

	if err != nil {
		logError("GetCourses: during mapping", err)
		return nil, err
	}
	return courses, err
}

func GetCourseByName(course_name string) (*Course, error) {
	singleCourse := Course{}
	if err := db.QueryRow("SELECT * FROM courses WHERE coursename=?", course_name).Scan(&singleCourse.Id, &singleCourse.Name, &singleCourse.CreatedAt, &singleCourse.UpdatedAt); err != nil {
		return nil, err
	}
	return &singleCourse, nil
}

func AddCourse(newCourse Course) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		logError("AddCourse: during Begin", err)
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO courses (coursename) VALUES (?)")

	if err != nil {
		logError("AddCourse: during Prepare", err)
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(newCourse.Name)

	if err != nil {
		logError("AddCourse: during Exec", err)
		return nil, err
	}

	tx.Commit()
	return res, nil
}

func DeleteCourse(cn string) (bool, error) {

	tx, err := db.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare("DELETE FROM courses WHERE coursename = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cn)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func logError(message string, err error) {
	log.Printf(message, err)
}
