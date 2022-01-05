package dbase

const (
	studentIQuery = `CREATE TABLE IF NOT EXISTS STUDENT( 
		student_id SERIAL PRIMARY KEY,
		student_name TEXT NOT NULL
	)
	`

	categoryIQuery = `CREATE TABLE IF NOT EXISTS CATEGORY( 
		category_id SERIAL PRIMARY KEY,
		category_name TEXT UNIQUE NOT NULL
	)
	`

	mentorIQuery = `CREATE TABLE IF NOT EXISTS MENTOR( 
		mentor_id SERIAL PRIMARY KEY,
		mentor_name TEXT NOT NULL,
		category_id int NOT NULL,
		CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(category_id)
	)
	`

	relationIQuery = `CREATE TABLE IF NOT EXISTS RELATION( 
		relation_id SERIAL PRIMARY KEY,
		student_id int UNIQUE NOT NULL,
		mentor_id int NOT NULL, 
		CONSTRAINT fk_student FOREIGN KEY(student_id) REFERENCES student(student_id),
		CONSTRAINT fk_mentor FOREIGN KEY(mentor_id) REFERENCES mentor(mentor_id)
	)
  	`
)
