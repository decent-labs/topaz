package flush

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

// Message is posted to the ethereum service
type Message struct {
	Address string
	Hash    string
}

// TXResp is returned from the ethereum service
type TXResp struct {
	TX string
}

var sh *shell.Shell
var db *sql.DB

func requestHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("starting flush service handler")

	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(
	// 		w,
	// 		fmt.Sprintf("error reading flush service request body: %s", err.Error()),
	// 		http.StatusBadRequest,
	// 	)
	// 	return
	// }

	// userID := string(b)

	// f := new(models.Flush)

	// if err := f.CreateFlush(database.Manager); err != nil {
	// 	log.Printf("couldn't create new flush: %s", err.Error())
	// 	return
	// }

	// defer db.Close()

	// log.Println("BATCH: BEGINNING 'flush()'.")

	// stmt := fmt.Sprintf("insert into flushes (app_id, created_at) values ('%d', now()) returning id;", id)

	// rows, err := db.Query(stmt)
	// if err != nil {
	// 	log.Printf("couldn't execute new flush statement: %s", err.Error())
	// 	return
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var id int

	// 	err = rows.Scan(&id)
	// 	if err != nil {
	// 		log.Printf("couldn't scan flush id: %s", err.Error())
	// 		continue
	// 	}

	// 	log.Printf("created flush id: %d", id)
	// }

	flushStmt := fmt.Sprintf(
		"insert into flushes (user_id) values ('%s') returning id, created_at;",
		userID,
	)

	flushRows, err := db.Query(flushStmt)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error adding row to flushes table: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer flushRows.Close()

	var flushID string
	var flushCreatedAt string

	for flushRows.Next() {
		err = flushRows.Scan(&flushID, &flushCreatedAt)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error scanning row into flush metadata strings: %s", err.Error()),
				http.StatusInternalServerError,
			)
			continue
		}
	}

	objStmt := fmt.Sprintf("select id, hash from objects where user_id = '%s' AND flush_id is null;",
		userID,
	)

	objRows, err := db.Query(objStmt)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error selecting objects to flush: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer objRows.Close()

	dir, err := sh.NewObject("unixfs-dir")
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error creating new ipfs emtpy directory: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	for objRows.Next() {
		var id string
		var hash string

		err = objRows.Scan(&id, &hash)
		if err != nil {
			log.Printf("error scanning object data into variables: %s", err.Error())
			continue
		}

		log.Printf("flushing object with hash: %s.", hash)

		dir, err = sh.PatchLink(dir, hash, hash, true)
		if err != nil {
			log.Printf("error patching ipfs directory with new hash: %s", err.Error())
			continue
		}

		log.Printf("new directory hash: %s", dir)

		// TODO: update ALL of the objects AT ONCE when the process is finished

		objUpdateStmt := fmt.Sprintf("update objects set flush_id = '%s' where id = '%s'", flushID, id)

		_, err := db.Exec(objUpdateStmt)
		if err != nil {
			log.Printf("error updating object with new flush id: %s", err.Error())
			continue
		}
	}

	flushUpdateStmt := fmt.Sprintf("update flushes set hash = '%s' where id = '%s'", dir, flushID)

	_, err = db.Exec(flushUpdateStmt)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error updating flush with directory hash: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	userUpdateStmt := fmt.Sprintf("update users set flushed_at = '%s' where id = '%s'",
		flushCreatedAt,
		userID,
	)

	_, err = db.Exec(userUpdateStmt)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error updating user with last flushed_at time: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	url := fmt.Sprintf("http://%s:%s/store", os.Getenv("ETHEREUM_HOST"), os.Getenv("ETHEREUM_PORT"))

	m := Message{os.Getenv("ETHEREUM_ADDRESS"), dir}

	b, err = json.Marshal(m)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("could not create JSON out of Message: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	read := bytes.NewReader(b)
	resp, err := http.Post(url, "application/octet-stream", read)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed posting data to ethereum service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer resp.Body.Close()

	txresp := new(TXResp)
	err = json.NewDecoder(resp.Body).Decode(txresp)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error decoding ethereum service tx response: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// TODO: save tx into a database now

	log.Println("finished with flush service handler")
}

func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(shConn)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	_db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer _db.Close()

	err = _db.Ping()
	if err != nil {
		log.Fatalf("couldn't ping database: %s", err.Error())
	}

	db = _db

	http.HandleFunc("/", requestHandler)

	log.Println("wake up, flush...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}