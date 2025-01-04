package consensus

import (
	"fmt"

	"github.com/Risotto04/blockchain/blockchain"
)

type Consensus struct {
	Approvals map[string]bool
	Quorum    int //จำนวนโหวตที่ยอมรับ
}

func selectValidator(blockHeight int, validators []string) string {
	index := blockHeight % len(validators)
	return validators[index]
}

// Validation
func (c *Consensus) AddApproval(validator string) {
	c.Approvals[validator] = true
}

// Check validate number
func (c *Consensus) HasQuorum() bool {
	count := 0
	for _, approved := range c.Approvals {
		if approved {
			count++
		}
	}
	return count >= c.Quorum
}

// Check validation and add block
func finalizeBlock(block *blockchain.Block, c *Consensus, bc *blockchain.Blockchain) {
	if c.HasQuorum() {
		bc.AddBlock(block.Data)
		fmt.Println("Block finalized and added")
	} else {
		fmt.Println("Quorum not reached. Block rejected.")
	}
}
