package mdb

import (
	"database/sql"
	"log"
	"time"
)

type MailEntry struct {
	Id          int64
	Email       string
	ConfirmedAt *time.Time
	OptOut      bool
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS email (
		id INTEGER PRIMARY KEY  ,
		email TEXT NOT NULL,
		confirmed_at INTEGER,
		opt_out INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

}

func emailEntryFromRow(row *sql.Rows) (*MailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool

	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t := time.Unix(confirmedAt, 0)
	return &MailEntry{
		Id:          id,
		Email:       email,
		ConfirmedAt: &t,
		OptOut:      optOut,
	}, err

}

func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`INSERT INTO email (email, confirmeat, opt_out) VALUES (?,0,false)`, email)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetEmail(db *sql.DB, email string) (*MailEntry, error) {
	row, err := db.Query(`SELECT * FROM emails WHERE email = ?`, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		return emailEntryFromRow(row)
	}
	return nil, nil
}

func UpdateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`UPDATE email SET confirmed_at = ?, opt_out = ? WHERE email = ?`, /*confirmedAt, optOut,*/ email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`DELETE FROM email WHERE email = ?`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type GetBatchQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatch(db *sql.DB, params GetBatchQueryParams) ([]MailEntry, error) {
	var empty []MailEntry

	rows, err := db.Query(`SELECT * FROM email ORDER BY id DESC LIMIT ? OFFSET ?`, params.Count, (params.Page-1)*params.Count)

	if err != nil {
		log.Println(err)
		return empty, err
	}
	defer rows.Close()

	emails := make([]MailEntry, 0, params.Count)

	for rows.Next() {
		email, err := emailEntryFromRow(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		emails = append(emails, *email)

	}
	return emails, nil
}
