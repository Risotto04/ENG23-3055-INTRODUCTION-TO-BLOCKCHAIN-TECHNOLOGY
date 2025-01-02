package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     time.Time
	Data          Transcript
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

type Blockchain struct {
	Blocks []*Block
}

type Transcript struct {
	StudentID   string
	StudentName string
	Faculty     string
	Major       string
	University  string
	Courses     []Course
	Gpa         float64
	IssueDate   time.Time
}

type Course struct {
	CourseCode string
	CourseName string
	Semester   string
	Credits    int
	Grade      string
}

func hashCourse(course Course) []byte {
	courseBytes, _ := json.Marshal(course)
	hash := sha256.Sum256(courseBytes)
	return hash[:]
}

func computeMerkleRoot(courses []Course) []byte {
	if len(courses) == 0 {
		return []byte{}
	}

	var hashes [][]byte
	for _, course := range courses {
		hashes = append(hashes, hashCourse(course))
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

func NewBlock(data Transcript, prevBlockHash []byte) *Block {
	merkleRoot := computeMerkleRoot(data.Courses)

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

func (bc *Blockchain) AddBlock(data Transcript) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
	genesisTranscript := Transcript{
		StudentID:   "B6499999",
		StudentName: "Genesis",
		Faculty:     "None",
		Major:       "None",
		University:  "Blockchain University",
		Courses:     []Course{},
		Gpa:         0.0,
		IssueDate:   time.Now(),
	}
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
	expectedMerkleRoot := computeMerkleRoot(block.Data.Courses)
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

func calculateBlockHash(block *Block) []byte {
	dataBytes, _ := json.Marshal(block.Data)

	timestamp := []byte(block.Timestamp.Format(time.RFC3339))

	headers := bytes.Join([][]byte{block.PrevBlockHash, dataBytes, block.MerkleRoot, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	return hash[:]
}
