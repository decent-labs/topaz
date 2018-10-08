package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	conn := "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	_, err = db.Exec(`create table users (
		"id" uuid primary key,
		"name" varchar(255),
		"interval" interval,
		"created_at" timestamp default now(),
		"flushed_at" timestamp
	);`)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`create table flushes (
		"id" uuid primary key,
		"user_id" uuid references users(id),
		"hash" varchar(255),
		"created_at" timestamp default now()
	);`)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`create table objects (
		"id" uuid primary key,
		"user_id" uuid references users(id),
		"flush_id" uuid references flushes(id),
		"hash" varchar(255),
		"created_at" timestamp default now()
	);`)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`insert into users (
		id,
		name,
		interval
	) values (
		'966D1035-73A1-4E57-9D1F-95920524689B',
		'parker',
		'10'
	);`)
}
