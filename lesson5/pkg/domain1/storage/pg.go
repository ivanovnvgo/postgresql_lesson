package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	dbpool *pgxpool.Pool
}

func NewPG(dbpool *pgxpool.Pool) *PG {
	return &PG{dbpool}
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
//func Search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]FullNameSearchDuplicate, error) {
func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]FullNameSearchDuplicate, error) {
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
	rows, err := s.dbpool.Query(ctx, sql, pattern, limit)
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
