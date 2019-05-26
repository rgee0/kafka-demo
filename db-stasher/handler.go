package function

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

var db *sql.DB

type inference struct {
	Score float64 `json:"score"`
	Name  string  `json:"name"`
}

type payload struct {
	URL        string      `json:"url"`
	Inferences []inference `json:"inferences"`
	Category   string      `json:"category"`
}

// init establishes a persistent connection to the remote database
// the function will panic if it cannot establish a link and the
// container will restart / go into a crash/back-off loop
func init() {
	password := string(getSecret("db-password"))
	user := string(getSecret("db-username"))
	host := string(getSecret("db-host"))

	dbName := os.Getenv("postgres_db")
	port := os.Getenv("postgres_port")
	sslmode := os.Getenv("postgres_sslmode")

	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=" + sslmode

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error

	dbErr := db.Ping()
	if dbErr != nil {
		return handler.Response{
			Body:       []byte("DB went away"),
			StatusCode: http.StatusInternalServerError,
		}, dbErr
	}

	//unmarshall body
	d := payload{}
	if err := json.Unmarshal(req.Body, &d); err != nil {
		log.Fatalf("Unable to unmarshal object for: %s", string(req.Body))
	}

	analysis, err := json.Marshal(d.Inferences)
	if err != nil {
		return handler.Response{
			Body:       []byte("string conversion of the analysis failed."),
			StatusCode: http.StatusOK,
		}, err
	}

	insertErr := insertMessage(d.Category, d.URL, string(analysis))
	if insertErr != nil {
		log.Printf("%s\n", insertErr.Error())
	}

	return handler.Response{
		Body:       []byte("value inserted"),
		StatusCode: http.StatusOK,
	}, err
}

func insertMessage(category string, URL string, analysisJSON string) error {
	res, err := db.Query(`insert into pictures (category, url, analysis, last_seen) values ($1, $2, $3, now()) ON CONFLICT (url) DO UPDATE SET last_seen = now();`,
		category, URL, analysisJSON)

	if err == nil {
		defer res.Close()
	}

	return err
}

func getSecret(name string) []byte {
	mounts := []string{"/var/openfaas/secrets/", "/run/secrets/"}
	var b []byte
	var err error
	for _, m := range mounts {
		if b, err = ioutil.ReadFile(m + name); err == nil {
			return b
		}
	}
	log.Fatal(err)
	return nil
}
