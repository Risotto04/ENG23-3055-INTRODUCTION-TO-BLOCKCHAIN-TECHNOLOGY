package controller

import (
	"net/http"

	"github.com/Risotto04/blockchain/blockchain"
	"github.com/Risotto04/blockchain/models"
	"github.com/gin-gonic/gin"
)

func AddCourse(c *gin.Context, courses *[]*models.Course) {
	var body models.Course

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	*courses = append(*courses, &body)
	c.JSON(http.StatusCreated, gin.H{"message": "Course added successfully"})
}

func AddBlock(c *gin.Context, courses *[]*models.Course, bc *blockchain.Blockchain) {
	if len(*courses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No courses to add to the blockchain"})
		return
	}

	bc.AddBlock(*courses)
	*courses = nil
	c.JSON(http.StatusCreated, gin.H{"message": "Block added successfully"})
}

func GetBlocks(c *gin.Context, bc *blockchain.Blockchain) {
	blocks := bc.GetBlocks() // Assuming `GetBlocks` is a method in your `Blockchain` struct that returns all blocks.
	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
	})
}
func GetCourses(c *gin.Context, courses *[]*models.Course) {
	c.JSON(http.StatusCreated, gin.H{"courses": *courses})

}
