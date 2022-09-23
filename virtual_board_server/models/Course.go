package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Course struct {
	Id        int       `json:"id"`
	Name      string    `json:"coursename"`
	CreatedAt time.Time `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
	UpdatedAt time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"UTC"`
}

func AddCourse(newCourse Course) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO courses (coursename) VALUES (?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newCourse.Name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetCourseCountById(course_id int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses where id = ?", course_id).Scan(&count)
	return count, err
}

func GetCourses() ([]Course, error) {
	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	courses := make([]Course, 0)

	for rows.Next() {
		singleCourse := Course{}
		err = rows.Scan(&singleCourse.Id, &singleCourse.Name, &singleCourse.CreatedAt, &singleCourse.UpdatedAt)

		if err != nil {
			return nil, err
		}

		courses = append(courses, singleCourse)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return courses, err

}
