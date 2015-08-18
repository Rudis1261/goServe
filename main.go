package main

import (
	"flag"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
	"fmt"
	"time"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

var PATH = flag.String("path", "public", "Public directory")
var PORT = flag.String("port", "3000", "On which port the HTTP server listens")
var MySQL = flag.String("dsn", "root:root@tcp(mysql:3306)/tvtracker", "MySQL host in DSN form (user:password@tcp(127.0.0.1:3306)/hello)")
var output string


func main() {

	flag.Parse()

	// Prepare SQL
	conn, err := sql.Open("mysql", *MySQL)
	checkErr(err)

	// Get the DATA from the DB
	data, err := getData(conn)
	checkErr(err)

	// Reasign the data
	output = data
	log.Printf("%s", output)
	log.Printf("%s", data)

	// Defer closing the SQL
	defer conn.Close()

	log.Printf("goServe Started, listening on PORT: %s", *PORT)
	requestHandlers := alice.New(loggingHandler, recoverHandler, imgHandler)
	http.Handle("/data", requestHandlers.ThenFunc(dataHandler))
	http.Handle("/favicon.ico", requestHandlers.ThenFunc(favIcoHandler))
	http.Handle("/", requestHandlers.ThenFunc(indexHandler))
	http.ListenAndServe(":"+*PORT, nil)
}


func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, *PATH+"/index.html")
}


func favIcoHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, *PATH+"/favicon.ico")
}


func dataHandler(rw http.ResponseWriter, req *http.Request) {
	if len(output) > 0 {
		fmt.Fprintf(rw, output)
		return
	}
	http.Error(rw, http.StatusText(404), 404)
}


func imgHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		reqUrl := req.URL.String()
		if len(reqUrl) > 0 {
			if _, err := os.Stat(*PATH + reqUrl); os.IsNotExist(err) {
				next.ServeHTTP(rw, req)
				return
			}
		}
		http.ServeFile(rw, req, *PATH+reqUrl)
		return
	}

	return http.HandlerFunc(fn)
}


func loggingHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		end := time.Now()
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), end.Sub(start))
	}

	return http.HandlerFunc(fn)
}


func recoverHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %+v", err)
				http.Error(rw, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}


// Generic Error handler
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}


func getData(db *sql.DB) (string, error) {
	var (
		id string
		name string
	)
	var data [][]string

	// Do the select and iterate through the ROWs
	rows, err := db.Query("SELECT `id`, `seriesname` FROM tv")
	checkErr(err)

	// Loop through each row, check for an error and handle it
	// Otherwise append it to the array
	for rows.Next() {
		err := rows.Scan(&id, &name)
		checkErr(err)
		result := make([]string, 2)
		result[0] = id
		result[1] = name
		data = append(data, result)
	}

	// Close the connection and defer closing
	checkErr(rows.Err())
	defer rows.Close()

	// Unmarshal the JSON and return the string
	jsonString, err := json.Marshal(data)
	checkErr(err)
	return string(jsonString), nil
}
