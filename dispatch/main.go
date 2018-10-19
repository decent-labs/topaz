package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	m "github.com/decentorganization/topaz/dispatch/models"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Find all users who are due for a flush and call the 'flush' service for them.
func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	i, err = strconv.Atoi(os.Getenv("DISPATCH_RATE"))
	if err != nil {
		log.Fatalf("missing valid DISPATCH_RATE environment variable: %s", err.Error())
	}

	log.Println("wake up, dispatch...")

	setupData(db)

	for range time.Tick(time.Duration(i) * time.Second) {
		log.Println("tick")

		// stmt := `
		// 	select distinct u.id
		// 	from users u
		// 		inner join objects o on o.user_id = u.id
		// 	where ((u.flushed_at is null) or (now() - u.flushed_at >= u.interval))
		// 		and (o.flush_id is null)
		// `

		// rows, err := db.Query(stmt)
		// if err != nil {
		// 	log.Printf("error executing user-finding query: %s", err.Error())
		// 	return
		// }
		// defer rows.Close()

		// var apps []App
		// var users []User
		// var flushes []Flush
		// db.Find(&apps)
		// log.Println(apps)

		// for rows.Next() {
		// 	var id string

		// 	err = rows.Scan(&id)
		// 	if err != nil {
		// 		log.Printf("error scanning row into user id var: %s", err.Error())
		// 		continue
		// 	}

		// 	url := fmt.Sprintf("http://%s:%s", os.Getenv("FLUSH_HOST"), os.Getenv("FLUSH_PORT"))
		// 	sr := strings.NewReader(id)
		// 	_, err = http.Post(url, "application/octet-stream", sr)
		// 	if err != nil {
		// 		log.Printf("error dispatching user id '%s' to flush service: %s", id, err.Error())
		// 		continue
		// 	}
		// }

		log.Println("tock")
	}
}

func setupData(db *gorm.DB) {
	parker := m.User{Name: "Parker"}
	adam := m.User{Name: "Adam"}
	nate := m.User{Name: "Nate"}
	db.Create(&parker)
	db.Create(&adam)
	db.Create(&nate)

	latestParkerFlushTime, _ := time.Parse("2006-01-02 15:04:05", "2018-10-18 09:30:30")
	secondLatestParkerFlushTime, _ := time.Parse("2006-01-02 15:04:05", "2018-10-18 09:30:00")
	latestNateFlushTime, _ := time.Parse("2006-01-02 15:04:05", "2018-10-18 09:30:00")

	appParker1 := m.App{Name: "parker app 1", User: parker, Interval: "30 seconds", LastFlushed: &latestParkerFlushTime}
	appNate1 := m.App{Name: "nate app 1", User: nate, Interval: "25 seconds", LastFlushed: &latestNateFlushTime}
	appAdam1 := m.App{Name: "adam app 1", User: adam, Interval: "35 seconds"}
	db.Create(&appParker1)
	db.Create(&appAdam1)
	db.Create(&appNate1)

	flushParker1 := m.Flush{Transaction: "0x0", App: appParker1, DirectoryHash: "0x0"}
	flushParker1.CreatedAt = latestParkerFlushTime
	flushParker2 := m.Flush{Transaction: "0x0", App: appParker1, DirectoryHash: "0x0"}
	flushParker2.CreatedAt = secondLatestParkerFlushTime
	flushNate1 := m.Flush{Transaction: "0x0", App: appNate1, DirectoryHash: "0x0"}
	flushNate1.CreatedAt = latestNateFlushTime
	db.Create(&flushParker1)
	db.Create(&flushParker2)
	db.Create(&flushNate1)

	objParkerFlush1 := m.Object{Hash: "0x0", App: appParker1, Flush: flushParker1}
	objParkerFlush2 := m.Object{Hash: "0x0", App: appParker1, Flush: flushParker2}
	objParkerUnFlush1 := m.Object{Hash: "0x0", App: appParker1}
	objParkerUnFlush2 := m.Object{Hash: "0x0", App: appParker1}
	objAdamUnFlush1 := m.Object{Hash: "0x0", App: appAdam1}
	objAdamUnFlush2 := m.Object{Hash: "0x0", App: appAdam1}
	objNateFlush1 := m.Object{Hash: "0x0", App: appNate1, Flush: flushNate1}
	db.Create(&objParkerFlush1)
	db.Create(&objParkerFlush2)
	db.Create(&objParkerUnFlush1)
	db.Create(&objParkerUnFlush2)
	db.Create(&objAdamUnFlush1)
	db.Create(&objAdamUnFlush2)
	db.Create(&objNateFlush1)
}
