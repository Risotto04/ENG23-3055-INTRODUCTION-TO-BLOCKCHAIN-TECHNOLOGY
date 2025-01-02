package main

import (
	"fmt"
	"time"

	"github.com/Risotto04/blockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	transcript1 := blockchain.Transcript{
		StudentID:   "B6419936",
		StudentName: "Teerachot Sonnok",
		Faculty:     "Engineering",
		Major:       "Computer Engineering",
		University:  "Khon Kaen University",
		Courses: []blockchain.Course{
			{CourseCode: "ENG101", CourseName: "English for Communication 1", Semester: "1/2023", Credits: 3, Grade: "A"},
			{CourseCode: "CPE102", CourseName: "Computer Programming I", Semester: "1/2023", Credits: 3, Grade: "B+"},
			{CourseCode: "MAT101", CourseName: "Calculus I", Semester: "1/2023", Credits: 4, Grade: "A"},
			{CourseCode: "PHY101", CourseName: "Physics I", Semester: "1/2023", Credits: 3, Grade: "B"},
			{CourseCode: "PHY102", CourseName: "Physics Laboratory I", Semester: "1/2023", Credits: 1, Grade: "A"},
		},
		Gpa:       3.75,
		IssueDate: time.Now(),
	}
	transcript2 := blockchain.Transcript{
		StudentID:   "B6419937",
		StudentName: "Teerachot Sonnok",
		Faculty:     "Engineering",
		Major:       "Computer Engineering",
		University:  "Khon Kaen University",
		Courses: []blockchain.Course{
			{CourseCode: "ENG101_", CourseName: "English for Communication 1", Semester: "1/2023", Credits: 3, Grade: "A"},
			{CourseCode: "CPE102", CourseName: "Computer Programming I", Semester: "1/2023", Credits: 3, Grade: "B+"},
			{CourseCode: "MAT101", CourseName: "Calculus I", Semester: "1/2023", Credits: 4, Grade: "A"},
			{CourseCode: "PHY101", CourseName: "Physics I", Semester: "1/2023", Credits: 3, Grade: "B"},
			{CourseCode: "PHY102", CourseName: "Physics Laboratory I", Semester: "1/2023", Credits: 1, Grade: "A"},
		},
		Gpa:       3.75,
		IssueDate: time.Now(),
	}

	bc.AddBlock(transcript1)
	bc.AddBlock(transcript2)

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data.StudentID)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Merkle Root: %x\n", block.MerkleRoot)
		fmt.Println("----------------------------------")

	}

	bc.GetBlock(1).Data.StudentID = "B6477777"              //invalid block hash
	bc.GetBlock(1).Data.Courses[0].CourseCode = "Physics I" //invalid Merkle Root

	if err := bc.Validate(); err != nil {
		fmt.Println("Blockchain validation failed:", err)
	} else {
		fmt.Println("Blockchain is valid!")
	}

}
