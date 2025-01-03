package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Risotto04/blockchain/models"
)

type Block struct {
	Timestamp     time.Time
	Data          []*models.Course
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

type Blockchain struct {
	Blocks []*Block
}

func hashCourse(course models.Course) []byte {
	courseBytes, _ := json.Marshal(course)
	hash := sha256.Sum256(courseBytes)
	return hash[:]
}

func computeMerkleRoot(courses []*models.Course) []byte {
	if len(courses) == 0 {
		return []byte{}
	}

	var hashes [][]byte
	for _, course := range courses {
		hashes = append(hashes, hashCourse(*course))
	}

	for len(hashes) > 1 {
		var newHashes [][]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {

				combined := append(hashes[i], hashes[i+1]...)
				newHash := sha256.Sum256(combined)
				newHashes = append(newHashes, newHash[:])
			} else {
				newHash := sha256.Sum256(hashes[i])
				newHashes = append(newHashes, newHash[:])
			}
		}
		hashes = newHashes
	}

	return hashes[0]
}

func (b *Block) SetHash() {
	dataBytes, _ := json.Marshal(b.Data)

	timestamp := []byte(b.Timestamp.Format(time.RFC3339))

	headers := bytes.Join([][]byte{b.PrevBlockHash, dataBytes, b.MerkleRoot, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data []*models.Course, prevBlockHash []byte) *Block {
	merkleRoot := computeMerkleRoot(data)

	block := &Block{
		Timestamp:     time.Now(),
		Data:          data,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    merkleRoot,
		Hash:          []byte{},
	}
	block.SetHash()
	return block
}

func (bc *Blockchain) AddBlock(data []*models.Course) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
	genesisTranscript := []*models.Course{
		{
			CourseCode: "Genesis",
			CourseName: "Blockchain",
			Semester:   "2024/1",
			Credits:    4,
			Score: []models.Score{
				{Student: models.Student{StudentID: "G6410000", StudentName: "Genesis"}, Point: 0},
			},
		}}
	return NewBlock(genesisTranscript, []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) Validate() error {
	for i, block := range bc.Blocks {
		if err := ValidateBlock(block); err != nil {
			return fmt.Errorf("block %d is invalid: %v", i, err)
		}

		if i > 0 && !bytes.Equal(block.PrevBlockHash, bc.Blocks[i-1].Hash) {
			return fmt.Errorf("block %d has an invalid PrevBlockHash", i)
		}
	}
	return nil
}

func ValidateBlock(block *Block) error {
	expectedMerkleRoot := computeMerkleRoot(block.Data)
	if !bytes.Equal(block.MerkleRoot, expectedMerkleRoot) {
		return fmt.Errorf("invalid Merkle Root")
	}

	expectedHash := calculateBlockHash(block)
	if !bytes.Equal(block.Hash, expectedHash) {
		return fmt.Errorf("invalid block hash")
	}

	return nil
}
func (bc *Blockchain) GetBlock(index int) *Block {
	if index < 0 || index >= len(bc.Blocks) {
		return nil
	}
	return bc.Blocks[index]
}
func (bc *Blockchain) GetBlocks() []*Block {
	return bc.Blocks
}
func calculateBlockHash(block *Block) []byte {
	dataBytes, _ := json.Marshal(block.Data)

	timestamp := []byte(block.Timestamp.Format(time.RFC3339))

	headers := bytes.Join([][]byte{block.PrevBlockHash, dataBytes, block.MerkleRoot, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	return hash[:]
}
