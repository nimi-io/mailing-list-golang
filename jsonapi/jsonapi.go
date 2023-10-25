package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailinlist/mdb"

	_ "github.com/mattn/go-sqlite3"

	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func fromJson[T any](body io.Reader, target T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	json.Unmarshal(buf.Bytes(), &target)

}

func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) {
	setJsonHeader(w)

	data, serverErr := withData()

	if serverErr != nil {
		w.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Print(err)
			return //data, err
		}
		w.Write(serverErrJson)
		return
	}

	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}

	w.Write(dataJson)
}

func returnErr(w http.ResponseWriter, err error, code int) {
	returnJson(w, func() (interface{}, error) {
		errMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errMessage, nil
	})

}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			returnErr(w, nil, 405)
			return
		}

		entry := mdb.MailEntry{}
		fromJson(r.Body, &entry)
		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Println("JSON CreateEmail:", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})

}

func GetEmai(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			returnErr(w, nil, 405)
			return
		}

		entry := mdb.MailEntry{}
		fromJson(r.Body, &entry)
		
		returnJson(w, func() (interface{}, error) {
			log.Println("JSON GetEmail:", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})

}

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "PUT" {
			returnErr(w, nil, 405)
			return
		}

		entry := mdb.MailEntry{}
		fromJson(r.Body, &entry)

		if err := mdb.UpdateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Println("JSON CreateEmail:", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})

}

func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			returnErr(w, nil, 405)
			return
		}

		entry := mdb.MailEntry{}
		fromJson(r.Body, &entry)
		if err := mdb.DeleteEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Println("JSON DeleteEmail:", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})

}

func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			returnErr(w, nil, 405)
			return
		}

		queryOptions := mdb.GetBatchQueryParams{}
		fromJson(r.Body, &queryOptions)

		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("Page and Count fields are required"), 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Println("JSON GetBatchEmail:", queryOptions)
			return mdb.GetEmailBatch(db, queryOptions)
		})

	})

}
func Serve(db *sql.DB, bind string) {

	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmai(db))
	http.Handle("/email/get_batch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))
	log.Printf("Json APi Server Listening on %v\n",bind)

	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("Json server error %v", err)
	}
}
