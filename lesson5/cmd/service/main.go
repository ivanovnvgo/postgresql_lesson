package main

import (
	"context"
	"fmt"
	"log"

	"lesson5/pkg/domain1/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Раздел INSERT PostgreSQL

type (
	GroupID   int
	StudentID int
	CourseID  int
	GradeID   int
)

type Group struct {
	GroupName string
}

type Student struct {
	FirstName   string
	LastName    string
	Condition   bool
	AverageMark float64
	GroupId     int
}

type Course struct {
	Title         string
	NumberOfHours int
}

type Grade struct {
	Grade     float64
	StudentId int
	CourseId  int
}

func insertGroup(ctx context.Context, dbpool *pgxpool.Pool, group Group) (GroupID, error) {
	const sql = `
insert into groups (group_name) values
($1)
returning id;
`
	// При insert разумно использовать метод dbpool.Exec,
	// который не требует возврата данных из запроса.
	// В данном случае после вставки строки мы получаем её идентификатор.
	// Идентификатор вставленной строки может быть использован
	// в интерфейсе приложения.

	var id GroupID
	err := dbpool.QueryRow(ctx, sql,
		// Параметры должны передаваться в том порядке,
		// в котором перечислены столбцы в SQL запросе.
		group.GroupName,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert groups: %w", err)
	}
	return id, nil
}

func insertStudent(ctx context.Context, dbpool *pgxpool.Pool, student Student) (StudentID, error) {
	const sql = `
insert into students (first_name, last_name, condition, average_mark, group_id) values
($1, $2, $3, $4, $5)
returning id;
`
	var id StudentID
	err := dbpool.QueryRow(ctx, sql,
		student.FirstName,
		student.LastName,
		student.Condition,
		student.AverageMark,
		student.GroupId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert students: %w", err)
	}
	return id, nil
}

func insertCourse(ctx context.Context, dbpool *pgxpool.Pool, course Course) (CourseID, error) {
	const sql = `
insert into courses (title, number_of_hours) values
($1, $2)
returning id;
`
	var id CourseID
	err := dbpool.QueryRow(ctx, sql,
		course.Title,
		course.NumberOfHours,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert courses: %w", err)
	}
	return id, nil
}

func insertGrade(ctx context.Context, dbpool *pgxpool.Pool, grade Grade) (StudentID, error) {
	const sql = `
insert into grades (grade, course_id, student_id) values
($1, $2, $3)
returning student_id;
`
	var id StudentID
	err := dbpool.QueryRow(ctx, sql,
		grade.Grade,
		grade.CourseId,
		grade.StudentId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert grades: %w", err)
	}
	return id, nil
}

func main() {
	// Подключение к БД
	ctx := context.Background()
	url := "postgres://admin:secret4All@localhost:5432/exam_session"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	// Раздел INSERT PostgreSQL
	dbpoolGroup, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpoolGroup.Close()
	group := Group{
		GroupName: "122-5",
	}
	idGroup, err := insertGroup(ctx, dbpoolGroup, group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(idGroup)

	dbpoolStudent, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpoolStudent.Close()
	student := Student{
		FirstName:   "Viktor",
		LastName:    "Petrov",
		Condition:   true,
		AverageMark: 4.9,
		GroupId:     5,
	}
	idStudent, err := insertStudent(ctx, dbpoolStudent, student)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(idStudent)

	dbpoolCourse, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpoolCourse.Close()
	course := Course{
		Title:         "C++",
		NumberOfHours: 107,
	}
	idCourse, err := insertCourse(ctx, dbpoolCourse, course)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(idCourse)

	dbpoolGrade, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpoolGrade.Close()
	grade := Grade{
		Grade:     4.2,
		CourseId:  4,
		StudentId: 4,
	}
	idStudent, err = insertGrade(ctx, dbpoolGrade, grade)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(idStudent)

	// Раздел SELECT PostgreSQL

	dbpoolSearch, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpoolSearch.Close()
	limit := 7
	hints, err := storage.Search(ctx, dbpoolSearch, "Pet", limit)
	if err != nil {
		log.Fatal(err)
	}
	for _, hint := range hints {
		fmt.Println(hint.LastName, hint.FirstName)
	}
}
