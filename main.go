package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type handler struct {
	DSN string // e.g. "bugzilla:secret@tcp(auroradb.dev.unee-t.com:3306)/bugzilla?multiStatements=true&sql_mode=TRADITIONAL"
	db  *sql.DB
}

func main() {

	h, err := New()
	if err != nil {
		log.WithError(err).Fatal("error setting configuration")
		return
	}

	defer h.db.Close()

	addr := ":" + os.Getenv("PORT")
	app := mux.NewRouter()
	app.HandleFunc("/", h.ping).Methods("GET")
	if err := http.ListenAndServe(addr, app); err != nil {
		log.WithError(err).Fatal("error listening")
	}

}

// New setups the configuration assuming various parameters have been setup in the AWS account
func New() (h handler, err error) {

	h = handler{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:3306)/test?multiStatements=true&sql_mode=TRADITIONAL",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASS"),
			os.Getenv("MYSQL_HOST"),
		),
	}

	h.db, err = sql.Open("mysql", h.DSN)
	if err != nil {
		log.WithError(err).Fatal("error opening database")
		return
	}

	return

}

func (h handler) ping(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping()
	if err != nil {
		log.WithError(err).Error("failed to ping database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "OK")
}