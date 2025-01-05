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
	go StartNode(models.Node{NodeID: "node1", Port: "8081", Peers: []string{"http://localhost:8082", "http://localhost:8083"}})
	go StartNode(models.Node{NodeID: "node2", Port: "8082", Peers: []string{"http://localhost:8081", "http://localhost:8083"}})
	go StartNode(models.Node{NodeID: "node3", Port: "8083", Peers: []string{"http://localhost:8081", "http://localhost:8082"}})

	// Block main thread
	select {}
}

var seenBlocks = make(map[string]bool)

func StartNode(config models.Node) {
	bc := blockchain.NewBlockchain()
	cons := consensus.NewConsensus(3)
	
	validators := []string{"node1", "node2", "node3"}
	var courses []*models.Course

	router := gin.Default()

	router.GET("/blockchain", func(ctx *gin.Context) {
		controller.GetBlocks(ctx, bc)
	})

	router.POST("/block", func(ctx *gin.Context) {
		handleBlockProposal(ctx, cons, bc, config.Peers, config.NodeID)
	})

	router.POST("/api/createBlock", func(ctx *gin.Context) {
		createBlock(bc, validators, config.NodeID, cons, config.Peers)
	})


	router.GET("/api/course", func(ctx *gin.Context) {
		controller.GetCourses(ctx, &courses)
	})
	router.POST("/api/course", func(ctx *gin.Context) {
		controller.AddCourse(ctx, &courses)
	})

	fmt.Printf("Starting node %s on port %s...\n", config.NodeID, config.Port)
	router.Run(":" + config.Port)


	// var transactionPool int
	// go func () {
	// 	for {
	// 		if transactionPool >= 5 {
	// 			createBlock(bc, validators, config.NodeID)
	// 		}
	// 		transactionPool += 1
	// 		fmt.Println(transactionPool)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

}

func broadcastBlock(block *blockchain.Block, peers []string, approvals map[string]bool) {
	blockHash := string(block.Hash)

    // Check if seen
    if seenBlocks[blockHash] {
        fmt.Println("Block already broadcasted, skipping.")
        return
    }
    seenBlocks[blockHash] = true

	// blockBytes, err := json.Marshal(block)
	// if err != nil { 
	// 	fmt.Printf("Error marshalling block: %v\n", err)
	// 	return
	// }
	// approvalsBytes, err := json.Marshal(approvals)
	// if err != nil {
	// 	fmt.Printf("Error marshalling approvals: %v\n", err)
	// 	return
	// }


	payload := map[string]interface{}{
		"block":     block,
		"approvals": approvals,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling payload: %v\n", err)
		return
	}

	// Broadcast to all peers
	for _, peer := range peers {
		url := peer + "/block"
		resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", peer, err)
		} else {
			fmt.Printf("Block sent to %s\n", peer)
			resp.Body.Close()
		}
	}
}

func handleBlockProposal(ctx *gin.Context, cons *consensus.Consensus, bc *blockchain.Blockchain, peers []string, nodeId string) {

	var payload struct {
		Block     blockchain.Block    `json:"block"`
		Approvals map[string]bool `json:"approvals"`
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block data"})
		return
	}

	block := payload.Block
	approvals := payload.Approvals

	// Validate the block structure
	if err := blockchain.ValidateBlock(&block); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Block validation failed: %v", err)})
		return
	}

	for approver, approved := range approvals {
		if approved {
			cons.AddApproval(approver)
		}
	}
	fmt.Printf("Updated approvals: %v\n", cons.Approvals)
	

	//Mark approval
	cons.AddApproval(nodeId)
	fmt.Println("Approval added by ", nodeId)
	fmt.Printf("Current approvals: %v\n", cons.Approvals)

	// Broadcast block to peers
	fmt.Println(seenBlocks[string(block.Hash)])
	broadcastBlock(&block, peers, cons.Approvals)

	// Check if quorum is reached and finalize the block
	fmt.Println("qourum", cons.HasQuorum())
	if cons.HasQuorum() {
		consensus.FinalizeBlock(&block, cons, bc)
		ctx.JSON(http.StatusOK, gin.H{"message": "Block finalized and added"})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "Block approval recorded, awaiting more votes"})
	}
}



func createBlock(bc *blockchain.Blockchain, validators []string, nodeId string, cons *consensus.Consensus, peer []string) {
	blockHeight := len(bc.Blocks)
	primaryValidator := consensus.SelectValidator(blockHeight, validators)
	fmt.Printf("Selected Validator for block %d: %s\n", blockHeight, primaryValidator)

	//Simulate that node get transaction from pool
	transactions := []*models.Course{
		{
			CourseCode: "CS101",
			CourseName: "Computer Science",
			Semester:   "2025/1",
			Credits:    3,
			Score:      []models.Score{{Student: models.Student{StudentID: "S12345", StudentName: "Alice"}, Point: 85}},
		},
	}
	newBlock := blockchain.NewBlock(transactions, bc.Blocks[blockHeight-1].Hash)

	if primaryValidator == nodeId {
		cons.AddApproval(nodeId)
		fmt.Println("Approval added by ", nodeId)
		broadcastBlock(newBlock, peer, cons.Approvals)
	}
}
