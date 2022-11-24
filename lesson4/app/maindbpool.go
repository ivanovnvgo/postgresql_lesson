package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"time"
)

func main() {
	ctx := context.Background()
	// Строка для подключения к базе данных
	url := "postgres://admin:secret4All@localhost:5432/exam_session"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}
	// Pool соединений обязательно ограничивать сверху,
	// так как иначе есть потенциальная опасность превысить лимит соединений с базой.
	cfg.MaxConns = 8
	cfg.MinConns = 4
	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute
	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour
	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute
	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	// pgx предоставляет набор адаптеров для популярных логеров
	// это позволяет организовать сбор ошибок при работе с базой
	// @see: https://github.com/jackc/pgx/tree/master/log
	// cfg.ConnConfig = zerologadapter.NewLogger(logger)
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	var greeting string
	err = dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Type dpool: %T\n", dbpool) // Type dpool: *pgxpool.Pool
	fmt.Printf("Type cfg: %T\n", cfg)      // Type cfg: *pgxpool.Config
	// Где в этом коде *pgx.Conn ?
	fmt.Println(cfg)
	fmt.Println(greeting)
}
