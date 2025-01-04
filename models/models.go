package models

type Student struct {
	StudentID   string
	StudentName string
}
type Score struct {
	Student Student
	Point   uint
}

type Course struct {
	CourseCode string
	CourseName string
	Semester   string
	Credits    int
	Score      []Score
}

type Node struct {
	NodeID string
	Port   string
	Peers  []string
}
