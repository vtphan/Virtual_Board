package main

import (
	"log"
	"net/http"
	"virtual_board/virtual_board_server/models"

	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/courses", getCourses)
		v1.POST("/courses", addCourse)
		//v1.PUT("/courses/:courseId", updateCourse)
		//v1.DELETE("/course/:courseId", deleteCourse)

		v1.GET("/courses/:courseId/notebooks", viewNotebook)
		v1.POST("/courses/:courseId/notebooks", addNotebook)
		v1.PUT("/courses/:courseId/notebooks/:notebookId", updateNotebook)
		v1.DELETE("/courses/:courseId/notebooks/:notebookId", deleteNotebook)

	}
	router.Run()

}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCourses(c *gin.Context) {
	courses, err := models.GetCourses()
	checkErr(err)

	//TODO: Test no courses scenario
	if courses == nil {
		c.JSON(http.StatusOK, gin.H{"courses": nil})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"courses": courses})
		return
	}
}

func addCourse(c *gin.Context) {

	var json models.Course

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddCourse(json)

	if success {
		c.JSON(http.StatusCreated, success)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
}

// func updateCourse(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "updateCourse Called"})
// }

func addNotebook(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("courseId"))

	course_counters, err := models.GetCourseCountById(courseId)

	if course_counters < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	log.Printf("Here 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Here 2")
	var json models.Notebook

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Here 3")
	json.CourseId = courseId

	success, err := models.AddNotebook(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func viewNotebook(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("courseId"))

	isWhiteBoards, err := strconv.ParseBool(c.Query("whiteboards"))

	notebook, err := models.ViewNotebook(courseId, isWhiteBoards)

	checkErr(err)
	// if the name is blank we can assume nothing is found
	if notebook == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": notebook})
	}
}

func updateNotebook(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("courseId"))

	notebookId, err := strconv.Atoi(c.Param("notebookId"))

	existingNotebook, err := models.GetNotebookById(courseId, notebookId)

	if err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}

func deleteNotebook(c *gin.Context) {

	cn := c.Param("coursename")

	success, err := models.DeleteNotebook(cn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
