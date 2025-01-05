package consensus

import (
	"fmt"

	"github.com/Risotto04/blockchain/blockchain"
)

type Consensus struct {
	Approvals map[string]bool
	Quorum    int
}

// Initialize a new Consensus instance
func NewConsensus(quorum int) *Consensus {
	return &Consensus{
		Approvals: make(map[string]bool),
		Quorum:    quorum,
	}
}

//Round-Robin
func SelectValidator(blockHeight int, validators []string) string {
	index := blockHeight % len(validators)
	return validators[index]
}

// AddApproval records approval from a validator
func (c *Consensus) AddApproval(validator string) {
	c.Approvals[validator] = true
}

// HasQuorum checks if enough approvals have been received
func (c *Consensus) HasQuorum() bool {
	count := 0
	for _, approved := range c.Approvals {
		if approved {
			count++
		}
		fmt.Println(count)
	}
	return count >= c.Quorum
}

// FinalizeBlock validates and adds the block if quorum is reached
func FinalizeBlock(block *blockchain.Block, c *Consensus, bc *blockchain.Blockchain) {
	if c.HasQuorum() {
		bc.AddBlock(block.Data)
		fmt.Println("Block finalized and added")
	} else {
		fmt.Println("Quorum not reached. Block rejected.")
	}
}
