package models

import (
	"database/sql"
	"errors"
	"server_archetype/defs"
	"strconv"

	m "github.com/SteveMCWin/archetype-common/models"
)

func (Db *DataBase) ReadQuote(quote_id int) (m.Quote, error) {

	quote := m.Quote{Id: quote_id}

	err := Db.Data.QueryRow("select source, quote, len from quotes where id = ?", quote_id).Scan(
		&quote.Source,
		&quote.Quote,
		&quote.Length,
	)

	if err != nil {
		return m.Quote{}, err
	}

	return quote, nil
}

func (Db *DataBase) RandomQuote(quote_len m.QuoteLen) (m.Quote, error) {
	quote := m.Quote{Length: quote_len}

	err := Db.Data.QueryRow("select id, source, quote from quotes where len = ? order by random() limit 1", quote_len).Scan(
		&quote.Id,
		&quote.Source,
		&quote.Quote,
	)

	if err != nil {
		return m.Quote{}, err
	}

	return quote, nil
}

func (Db *DataBase) CreateQuote(quote m.Quote) (uint64, error) {
	err := Db.Data.QueryRow("select id from quotes where quote = ?", quote.Quote).Scan(&quote.Id)

	if err == nil {
		return defs.NO_ID, errors.New("ERROR: quote already exists, it's id is " + strconv.Itoa(quote.Id))
	}

	statement := "insert into quotes (source, quote, len) values (?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err = Db.Data.Prepare(statement)
	if err != nil {
		return defs.NO_ID, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
		quote.Source,
		quote.Quote,
		quote.Length,
	).Scan(&quote.Id)

	if err != nil {
		return defs.NO_ID, err
	}

	return uint64(quote.Id), nil
}

func (Db *DataBase) UpdateQuote(quote m.Quote) error {
	statement := "UPDATE quotes SET source = ?, quote = ?, len = ? WHERE id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(quote.Source, quote.Quote, quote.Length, quote.Id)

	return err
}

func (Db *DataBase) DeleteQuote(quote_id int) error {
	statement := "DELETE from quotes where id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(quote_id)

	return err
}
