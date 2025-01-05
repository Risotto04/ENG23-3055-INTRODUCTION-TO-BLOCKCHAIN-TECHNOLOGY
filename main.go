package main

import (
	// "github.com/Risotto04/blockchain/blockchain"
	// "github.com/Risotto04/blockchain/controller"
	// "github.com/Risotto04/blockchain/models"
	"github.com/Risotto04/blockchain/models"
	"github.com/Risotto04/blockchain/node"
	// "github.com/gin-gonic/gin"
)

func main() {
	// router := gin.Default()
	// bc := blockchain.NewBlockchain()
	// var courses []*models.Course
	// router.GET("/api/course", func(ctx *gin.Context) {
	// 	controller.GetCourses(ctx, &courses)
	// })
	// router.POST("/api/course", func(ctx *gin.Context) {
	// 	controller.AddCourse(ctx, &courses)
	// })
	// router.GET("/api/blockchain", func(ctx *gin.Context) {
	// 	controller.GetBlocks(ctx, bc)
	// })
	// router.POST("/api/blockchain", func(ctx *gin.Context) {
	// 	controller.AddBlock(ctx, &courses, bc)
	// })

	

	// router.Run()
	// node.InitNodes()

	node.StartNode(models.Node{NodeID: "node1", Port: "8081", Peers: []string{"http://localhost:8082", "http://localhost:8083"}})
	node.StartNode(models.Node{NodeID: "node2", Port: "8082", Peers: []string{"http://localhost:8081", "http://localhost:8083"}})
	node.StartNode(models.Node{NodeID: "node3", Port: "8083", Peers: []string{"http://localhost:8081", "http://localhost:8082"}})
}
