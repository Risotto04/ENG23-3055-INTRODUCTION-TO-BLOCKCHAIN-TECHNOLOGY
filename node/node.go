package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Risotto04/blockchain/blockchain"
	"github.com/Risotto04/blockchain/controller"
	"github.com/Risotto04/blockchain/models"
	"github.com/gin-gonic/gin"
)

// Run nodes
func InitNodes() {
	go startNode(models.Node{NodeID: "node1", Port: "8081", Peers: []string{"http://localhost:8082", "http://localhost:8083"}})
	go startNode(models.Node{NodeID: "node2", Port: "8082", Peers: []string{"http://localhost:8081", "http://localhost:8083"}})
	go startNode(models.Node{NodeID: "node3", Port: "8083", Peers: []string{"http://localhost:8081", "http://localhost:8082"}})

	// Block main thread
	select {}
}

// Init node
func startNode(config models.Node) {
	bc := blockchain.NewBlockchain()

	router := gin.Default()

	router.GET("/blockchain", func(ctx *gin.Context) {
		// Send the node's blockchain to the requester
		controller.GetBlocks(ctx, bc)
	})

	router.POST("/block", func(ctx *gin.Context) {
		// Handle incoming block proposals
		handleBlockProposal(ctx, bc, config.Peers)
	})

	fmt.Printf("Starting node %s on port %s...\n", config.NodeID, config.Port)
	router.Run(":" + config.Port)

}

// Broadcast block
func broadcastBlock(block *blockchain.Block, peers []string) {
	blockBytes, _ := json.Marshal(block)
	for _, peer := range peers {
		url := peer + "/validate"
		_, err := http.Post(url, "application/json", bytes.NewReader(blockBytes))
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", peer, err)
		} else {
			fmt.Printf("Block sent to %s\n", peer)
		}
	}
}

// Handle incoming block
func handleBlockProposal(w http.ResponseWriter, r *http.Request) {
	var block blockchain.Block
	err := json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		http.Error(w, "Invalid block data", http.StatusBadRequest)
		return
	}

	// Validate the block
	if err := blockchain.ValidateBlock(&block); err != nil {
		http.Error(w, fmt.Sprintf("Block validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Add approval and check consensus
	//แดงๆ งงๆ
	consensus
	if consensus.hasQuorum() {
		finalizeBlock(&block, consensus, blockchain)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Block approved")
}
