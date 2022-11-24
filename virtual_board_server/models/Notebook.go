package models

import (
	"fmt"
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
		fmt.Println("AddNotebook", err)
		return false, err
	}
	fmt.Println("Adding notebook")
	stmt, err := ad.Prepare("INSERT INTO notebook (course_id, notebookname,content) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("AddNotebook Prepare", err)
		return false, err
	}

	defer stmt.Close()
	_, err = stmt.Exec(newNotebook.CourseId, newNotebook.NotebookName, newNotebook.Content)

	if err != nil {
		fmt.Println("AddNotebook Exec", err)
		return false, err
	}

	ad.Commit()

	return true, nil
}

func GetAllNotebooks(course_id int) ([]Notebook, error) {
	rows, err := db.Query("SELECT * from notebook WHERE course_id = ? AND (DATE(created_at)=(SELECT date('now')) OR DATE(updated_at)=(SELECT date('now')))", course_id)
	if err != nil {
		fmt.Println("Error during select on NB table!!", err)
		return nil, err
	}

	defer rows.Close()

	notebooks := make([]Notebook, 0)

	for rows.Next() {
		singleNotebook := Notebook{}
		sqlErr := rows.Scan(&singleNotebook.Id, &singleNotebook.CourseId, &singleNotebook.NotebookName, &singleNotebook.Content, &singleNotebook.CreatedAt, &singleNotebook.UpdatedAt)
		if sqlErr != nil {
			log.Println("error here!!", err)
			return nil, err
		}
		notebooks = append(notebooks, singleNotebook)

	}
	return notebooks, err
}

func DeleteNotebook(cn int) (bool, error) {

	tx, err := db.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare("DELETE from notebook where course_id = ?")

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
		stmt, err := tx.Prepare("UPDATE notebook SET content = ?, updated_at = ? WHERE course_id= ? AND notebookname = ?")
		if err != nil {
			return false, err
		}

		defer stmt.Close()
		_, err = stmt.Exec(requestNotebook.Content, time.Now().UTC().Format(time.RFC3339), existingNotebook.CourseId, existingNotebook.NotebookName)
		if err != nil {
			return false, err
		}
		tx.Commit()
	}
	return true, nil
}

func GetNotebook(course_id int, notebookname string) (Notebook, error) {
	singleNotebook := Notebook{}

	err := db.QueryRow("SELECT * FROM notebook where course_id = ? AND notebookname = ?", course_id, notebookname).Scan(&singleNotebook.Id, &singleNotebook.CourseId, &singleNotebook.NotebookName, &singleNotebook.Content, &singleNotebook.CreatedAt, &singleNotebook.UpdatedAt)
	if err != nil {
		fmt.Println("GetNotebook err", err)
	}
	return singleNotebook, err
}

func GetNotebookByName(course_id int, notebookname string) (*Notebook, error) {
	singleNotebook := Notebook{}
	if err := db.QueryRow("SELECT * FROM notebook where course_id=? AND notebookname=?", course_id, notebookname).Scan(&singleNotebook.Id, &singleNotebook.CourseId, &singleNotebook.NotebookName, &singleNotebook.Content, &singleNotebook.CreatedAt, &singleNotebook.UpdatedAt); err != nil {
		return nil, err
	}
	return &singleNotebook, nil
}
