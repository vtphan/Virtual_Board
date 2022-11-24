package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"virtual_board/virtual_board_server/models"

	"github.com/gin-gonic/gin"
)

type config struct {
	Server     string
	Coursename string
}

func main() {
	// EncodeJson() - TODO: Remove this
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./templates")

	v1 := router.Group("/api/v1")
	{
		v1.GET("/courses", fetchCourses)
		v1.GET("/courses/view", viewCourses)
		v1.GET("/courses/:courseName", fetchCourseByName)
		v1.POST("/courses", addCourse)
		//v1.PUT("/courses/:courseId", updateCourse)
		// v1.DELETE("/courses/:courseName", deleteCourse)

		v1.GET("/courses/:courseName/notebooks", fetchAllNotebookByCourse)
		v1.GET("/courses/:courseName/notebooks/view", viewNotebooks)
		v1.GET("/courses/:courseName/notebooks/:notebookName", fetchNotebookByName)
		v1.POST("/courses/:courseName/notebooks", addNotebook)
		v1.PUT("/courses/:courseName/notebooks/:notebookName", updateNotebook)
		v1.DELETE("/courses/:courseName/notebooks/:notebookName", deleteNotebook)
		v1.DELETE("/courses/:courseName", deleteCourse)
	}
	router.Run()
}

func fetchCourses(c *gin.Context) {
	courses, err := models.GetCourses()
	if err != nil {
		logError("fetchCourses: during GetCourses", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"courses": courses})
	}
}

func fetchCourseByName(c *gin.Context) {
	courseName := c.Param("courseName")
	course, err := models.GetCourseByName(courseName)
	if err != nil {
		logError("fetchCourseByName: during GetCourseByName", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"course": course})
}

func addCourse(c *gin.Context) {
	var json models.Course
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course, err := models.GetCourseByName(json.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			logError("fetchCourseByName: during GetCourseByName", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if course != nil {
		//Course exists
		c.JSON(http.StatusConflict, gin.H{"error": "Course already exists"})
		return
	}

	res, err := models.AddCourse(json)

	if res != nil {
		id, _ := res.LastInsertId()
		c.JSON(http.StatusCreated, gin.H{"id": id})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
}

func fetchAllNotebookByCourse(c *gin.Context) {
	courseName := c.Param("courseName")
	course, err := models.GetCourseByName(courseName)
	if err != nil {
		logError("fetchAllNotebookByCourse: during GetCourseByName", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	notebooks, err := models.GetAllNotebooks(course.Id)
	if err != nil {
		logError("fetchAllNotebookByCourse: during GetAllNotebooks", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notebooks": notebooks})
}

func fetchNotebookByName(c *gin.Context) {
	notebookName := c.Param("notebookName")
	courseName := c.Param("courseName")
	course, err := models.GetCourseByName(courseName)
	if err != nil {
		logError("fetchNotebookByName: during GetCourseByName", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	notebook, err := models.GetNotebookByName(course.Id, notebookName)

	if err != nil {
		logError("fetchNotebookByName: during GetNotebookByName", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notebook does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(notebook.Content))
}

func addNotebook(c *gin.Context) {
	courseName := c.Param("courseName")
	course, err := models.GetCourseByName(courseName)

	//TODO: revisit for 404
	if err != nil {
		logError("addNotebook: during GetCourseByName", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var json models.Notebook
	if err := c.ShouldBindJSON(&json); err != nil {
		logError("addNotebook: during ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	json.CourseId = course.Id
	success, err := models.AddNotebook(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func updateNotebook(c *gin.Context) {
	notebookName := c.Param("notebookName")
	courseName := c.Param("courseName")
	course, err := models.GetCourseByName(courseName)

	existingNotebook, err := models.GetNotebook(course.Id, notebookName)

	if err != nil {
		fmt.Println("err1", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var requestNotebook models.Notebook
	if err := c.ShouldBindJSON(&requestNotebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.UpdateNotebook(requestNotebook, existingNotebook)

	if success {
		c.Status(http.StatusNoContent)
	} else {
		fmt.Println("err2", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func deleteNotebook(c *gin.Context) {

	cn := c.Param("courseName")
	course, err := models.GetCourseByName(cn)
	success, err := models.DeleteNotebook(course.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
func deleteCourse(c *gin.Context) {

	cn := c.Param("courseName")
	success, err := models.DeleteCourse(cn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func logError(message string, err error) {
	log.Print("At: {} \n Error: {}", message, err)
}

func viewCourses(c *gin.Context) {
	baseUrl := "http://localhost:8080/api/v1/courses/"
	courses, err := models.GetCourses()
	if err != nil {
		logError("fetchCourses: during GetCourses", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	courseNames := populateCourseNames(courses)
	c.HTML(http.StatusOK, "course.html", gin.H{
		"title":       "Displays courses",
		"baseUrl":     baseUrl,
		"courseNames": courseNames,
	})
}

func viewNotebooks(c *gin.Context) {
	courseName := c.Param("courseName")
	baseUrl := "http://localhost:8080/api/v1/courses/" + courseName + "/notebooks/"
	course, err := models.GetCourseByName(courseName)
	if err != nil {
		logError("fetchAllNotebookByCourse: during GetCourseByName", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	notebooks, err := models.GetAllNotebooks(course.Id)

	notebook := populateNotebookNames(notebooks)
	c.HTML(http.StatusOK, "notebook.html", gin.H{
		"title":     "Displays notebooks",
		"baseUrl":   baseUrl,
		"notebooks": notebook,
		"course":    courseToString(course),
	})
}

func populateNotebookNames(notebooks []models.Notebook) string {
	s := make([]string, 0, len(notebooks))
	for _, notebook := range notebooks {
		bytes, err := json.Marshal(notebook)
		if err != nil {
			panic(err)
		}
		s = append(s, string(bytes))
	}
	output := strings.Join(s, `&`)
	return output
}

func populateCourseNames(courses []models.Course) string {
	s := make([]string, 0, len(courses))
	for _, course := range courses {
		bytes, err := json.Marshal(course)
		if err != nil {
			panic(err)
		}
		s = append(s, string(bytes))
	}
	output := strings.Join(s, `&`)
	return output
}

func courseToString(course *models.Course) string {
	bytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
