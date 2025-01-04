package consensus

import (
	"fmt"

	"github.com/Risotto04/blockchain/blockchain"
)

type Consensus struct {
	Approvals map[string]bool
	Quorum    int
}

func selectValidator(blockHeight int, validators []string) string {
	index := blockHeight % len(validators)
	return validators[index]
}

func (c *Consensus) AddApproval(validator string) {
	c.Approvals[validator] = true
}

func (c *Consensus) HasQuorum() bool {
	count := 0
	for _, approved := range c.Approvals {
		if approved {
			count++
		}
	}
	return count >= c.Quorum
}

func finalizeBlock(block *blockchain.Block, c *Consensus, bc *blockchain.Blockchain) {
	if c.HasQuorum() {
		bc.AddBlock(block.Data)
		fmt.Println("Block finalized and added")
	} else {
		fmt.Println("Quorum not reached. Block rejected.")
	}
}
