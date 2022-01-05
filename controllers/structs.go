package controllers

type RegsiterDetail struct {
	Name     string `json:"name"`
	Choice   string `json:"choice"`
	Category string `json:"category_name"`
}

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DuoIdentity struct {
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}
