package models

import (
	"database/sql"
	"errors"
	"math/rand"
	"server_archetype/defs"
	"strconv"
	"strings"

	m "github.com/SteveMCWin/archetype-common/models"
)

func (Db *DataBase) ReadQuote(quote_id uint64) (m.Quote, error) {

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

func (Db *DataBase) RandomQuote(req m.QuoteRequest) (m.Quote, error) {
	if req.Length == m.QUOTE_ANY_SIZE {
		req.Length = m.QuoteLen(rand.Intn(int(m.QUOTE_ANY_SIZE)))
	}

	quote := m.Quote{Length: req.Length}

	err := Db.Data.QueryRow("select id, source, quote from quotes where len = ? order by random() limit 1", req.Length).Scan(
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
	quote.Quote = strings.TrimSpace(quote.Quote)

	err := Db.Data.QueryRow("select id from quotes where quote = ?", quote.Quote).Scan(&quote.Id)

	if err == nil {
		return defs.NO_ID, errors.New("ERROR: quote already exists, it's id is " + strconv.FormatUint(quote.Id, 10))
	}

	quote.Length = m.CalcQuoteLen(quote.Quote)
	if quote.Length <= 1 {
		return defs.NO_ID, errors.New("ERROR: quote is too short!")
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
