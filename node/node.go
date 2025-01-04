package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Risotto04/blockchain/blockchain"
	"github.com/Risotto04/blockchain/consensus"
	"github.com/Risotto04/blockchain/controller"
	"github.com/Risotto04/blockchain/models"
	"github.com/gin-gonic/gin"
)

func InitNodes() {
	go startNode(models.Node{NodeID: "node1", Port: "8081", Peers: []string{"http://localhost:8082", "http://localhost:8083"}})
	go startNode(models.Node{NodeID: "node2", Port: "8082", Peers: []string{"http://localhost:8081", "http://localhost:8083"}})
	go startNode(models.Node{NodeID: "node3", Port: "8083", Peers: []string{"http://localhost:8081", "http://localhost:8082"}})

	// Block main thread
	select {}
}

func startNode(config models.Node) {
	bc := blockchain.NewBlockchain()
	cons := consensus.NewConsensus(2) // Quorum = 2 approvals
	var courses []*models.Course

	router := gin.Default()

	router.GET("/blockchain", func(ctx *gin.Context) {
		controller.GetBlocks(ctx, bc)
	})

	router.POST("/block", func(ctx *gin.Context) {
		handleBlockProposal(ctx, cons, bc, config.Peers)
	})

	// router.POST("/api/blockchain", func(ctx *gin.Context) {
	// 	controller.AddBlock(ctx, &courses, bc)
	// })

	router.GET("/api/course", func(ctx *gin.Context) {
		controller.GetCourses(ctx, &courses)
	})
	router.POST("/api/course", func(ctx *gin.Context) {
		controller.AddCourse(ctx, &courses)
	})

	fmt.Printf("Starting node %s on port %s...\n", config.NodeID, config.Port)
	router.Run(":" + config.Port)
}

func broadcastBlock(block *blockchain.Block, peers []string) {
	blockBytes, _ := json.Marshal(block)
	for _, peer := range peers {
		url := peer + "/block"
		_, err := http.Post(url, "application/json", bytes.NewReader(blockBytes))
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", peer, err)
		} else {
			fmt.Printf("Block sent to %s\n", peer)
		}
	}
}

func handleBlockProposal(ctx *gin.Context, cons *consensus.Consensus, bc *blockchain.Blockchain, peers []string) {
	var block blockchain.Block
	if err := ctx.ShouldBindJSON(&block); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block data"})
		return
	}

	// Validate the block structure
	if err := blockchain.ValidateBlock(&block); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Block validation failed: %v", err)})
		return
	}

	// Simulate approval by this node
	cons.AddApproval("self")
	fmt.Println("Approval added by this node.")

	// Broadcast block to peers for further approvals
	broadcastBlock(&block, peers)

	// Check if quorum is reached and finalize the block
	if cons.HasQuorum() {
		consensus.FinalizeBlock(&block, cons, bc)
		ctx.JSON(http.StatusOK, gin.H{"message": "Block finalized and added"})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "Block approval recorded, awaiting more votes"})
	}
}
