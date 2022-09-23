package models

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Notebook struct {
	Id           int       `json:"id"`
	CourseId     int       `json:"course_id"`
	NotebookName string    `json:"notebookname"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
	UpdatedAt    time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"UTC"`
}

func AddNotebook(newNotebook Notebook) (bool, error) {
	ad, err := db.Begin()
	if err != nil {
		log.Println("error here!!")
		return false, err
	}

	stmt, err := ad.Prepare("INSERT INTO notebook (course_id, notebookname,content) VALUES (?, ?, ?)")

	if err != nil {
		log.Println("error here!!")
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newNotebook.CourseId, newNotebook.NotebookName, newNotebook.Content)

	if err != nil {
		return false, err
	}

	ad.Commit()

	return true, nil
}

func ViewNotebook(course_id int, isWhiteBoards bool) ([]Notebook, error) {
	var rows1 *sql.Rows
	var err error
	if isWhiteBoards {
		rows1, err = db.Query("SELECT * from notebook WHERE course_id = ? and date(created_at) = ?", course_id, time.Now().UTC().Format("2006-01-02"))
	} else {
		rows1, err = db.Query("SELECT * from notebook WHERE course_id = ?", course_id)
	}

	defer rows1.Close()
	notebooks := make([]Notebook, 0)

	for rows1.Next() {
		singleNotebook := Notebook{}
		sqlErr := rows1.Scan(&singleNotebook.Id, &singleNotebook.CourseId, &singleNotebook.NotebookName, &singleNotebook.Content, &singleNotebook.CreatedAt, &singleNotebook.UpdatedAt)
		if sqlErr != nil {
			return nil, err
		}
		notebooks = append(notebooks, singleNotebook)

	}

	err = rows1.Err()

	if err != nil {
		return nil, err
	}
	return notebooks, err
}

func DeleteNotebook(cn string) (bool, error) {

	tx, err := db.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare("DELETE from notebook where coursename = ?")

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
func UpdateNotebook(requestNotebook Notebook, existingNotebook Notebook) (bool, error) {

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	if existingNotebook.Content != requestNotebook.Content {
		stmt, err := tx.Prepare("UPDATE notebook SET content = ?, updated_at = ? WHERE course_id = ? AND id = ?")
		if err != nil {
			return false, err
		}

		defer stmt.Close()
		_, err = stmt.Exec(requestNotebook.Content, time.Now().UTC().Format(time.RFC3339), existingNotebook.CourseId, existingNotebook.Id)
		if err != nil {
			return false, err
		}
		tx.Commit()
	}
	return true, nil
}

func GetNotebookById(course_id int, id int) (Notebook, error) {
	singleNotebook := Notebook{}
	err := db.QueryRow("SELECT * FROM notebook where id = ? AND course_id = ?", id, course_id).Scan(&singleNotebook.Id, &singleNotebook.CourseId, &singleNotebook.NotebookName, &singleNotebook.Content, &singleNotebook.CreatedAt, &singleNotebook.UpdatedAt)
	return singleNotebook, err
}
