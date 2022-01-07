package dbase

const (
	studentCQuery = `CREATE TABLE IF NOT EXISTS STUDENT( 
		student_id SERIAL PRIMARY KEY,
		student_name TEXT NOT NULL
	)
	`

	categoryCQuery = `CREATE TABLE IF NOT EXISTS CATEGORY( 
		category_id SERIAL PRIMARY KEY,
		category_name TEXT UNIQUE NOT NULL
	)
	`

	mentorCQuery = `CREATE TABLE IF NOT EXISTS MENTOR( 
		mentor_id SERIAL PRIMARY KEY,
		mentor_name TEXT NOT NULL,
		category_id int NOT NULL,
		CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(category_id)
	)
	`

	relationCQuery = `CREATE TABLE IF NOT EXISTS RELATION( 
		relation_id SERIAL PRIMARY KEY,
		student_id int UNIQUE NOT NULL,
		mentor_id int NOT NULL, 
		CONSTRAINT fk_student FOREIGN KEY(student_id) REFERENCES student(student_id),
		CONSTRAINT fk_mentor FOREIGN KEY(mentor_id) REFERENCES mentor(mentor_id)
	)
  	`

	studentIQuery = `INSERT INTO student (student_name)
	VALUES ($1) 
	RETURNING student_id
	`

	categoryIQuery = `INSERT INTO category (category_name)
	VALUES ($1) 
	RETURNING category_id
	`

	mentorIQuery = `INSERT INTO mentor(mentor_name, category_id)
	VALUES ($1, $2)
	RETURNING mentor_id
	`

	relationIQuery = `INSERT INTO relation (student_id, mentor_id)
	VALUES ($1, $2)
	ON CONFLICT (student_id)
	DO NOTHING
	RETURNING relation_id
	`

	getCategoryIdQuery = `SELECT category_id 
	FROM CATEGORY
	WHERE category_name=$1
	`

	getStudentQuery = `SELECT s.student_id, s.student_name
	FROM relation r
	INNER JOIN student s
	ON s.student_id = r.student_id
	WHERE mentor_id=$1
	`

	getMentorQuery = `SELECT mentor_id,mentor_name
	FROM mentor
	WHERE category_id=$1`
)
