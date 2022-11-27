package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

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

// Раздел SELECT PostgreSQL

type (
	FirstName string
	LastName  string
)
type FullNameSearchDuplicate struct {
	FirstName FirstName
	LastName  LastName
}

// search ищет всех сотрудников со схожими фамилиями.
// Из функции возвращается список FullNameSearchDuplicate, отсортированный по FirstName.
// Размер возвращаемого списка ограничен значением limit.
func search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]FullNameSearchDuplicate, error) {
	const sql = `
	select
	first_name,
	last_name
	from students
	where last_name like $1
	order by first_name asc
	limit $2;
	`
	pattern := prefix + "%"
	rows, err := dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()
	// В слайс hints будут собраны все строки, полученные из базы
	var hints []FullNameSearchDuplicate
	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint FullNameSearchDuplicate
		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.FirstName, &hint.LastName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		hints = append(hints, hint)
	}
	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}
	return hints, nil
}

// Раздел нагрузочный тест PostgreSQL

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64
	attacker := func(stopAt time.Time) {
		for {
			_, err := search(ctx, dbpool, "alex", 5)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)
			if time.Now().After(stopAt) {
				return
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(threads)
	startAt := time.Now()
	stopAt := startAt.Add(duration)
	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}
	wg.Wait()
	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
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
	hints, err := search(ctx, dbpoolSearch, "Pet", limit)
	if err != nil {
		log.Fatal(err)
	}
	for _, hint := range hints {
		fmt.Println(hint.LastName, hint.FirstName)
	}

	// Раздел нагрузочный тест PostgreSQL

	cfg.MaxConns = 16
	cfg.MinConns = 16
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	duration := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Println("start attack")
	res := attack(ctx, duration, threads, dbpool)
	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)

}
