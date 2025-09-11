package models

import (
	"database/sql"
	"errors"
	"log"
	"server_archetype/defs"
	"time"

	m "github.com/SteveMCWin/archetype-common/models"
	"golang.org/x/crypto/bcrypt"
)

func (Db *DataBase) ReadUser(user_id int) (m.User, error) {
	user := m.User{}

	err := Db.Data.QueryRow("select id, username, password, email, date_created, tests_started, tests_completed, all_time_avg_wpm, all_time_avg_acc from users where id = ?", user_id).Scan(
		&user.Id,
		&user.UserName,
		&user.Password,
		&user.Email,
		&user.DateCreated,
		&user.TestsStarted,
		&user.TestsCompleted,
		&user.AllTimeAvgWPM,
		&user.AllTimeAvgACC,
	)

	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (Db *DataBase) AuthUser(email, password string) (m.User, error) {

	user := m.User{}

	err := Db.Data.QueryRow("select id, username, password, email, date_created, tests_started, tests_completed, all_time_avg_wpm, all_time_avg_acc from users where email = ?", email).Scan(
		&user.Id,
		&user.UserName,
		&user.Password,
		&user.Email,
		&user.DateCreated,
		&user.TestsStarted,
		&user.TestsCompleted,
		&user.AllTimeAvgWPM,
		&user.AllTimeAvgACC,
	)

	if err != nil {
		return m.User{}, err
	}

	// check if password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("The provided password isn't correct")
		return m.User{}, err
	}

	return user, nil
}

func (Db *DataBase) CreateUser(user *m.User) (uint64, error) {
	if user.Email == "" {
		return defs.NO_ID, errors.New("Cannot store a user without their email")
	}

	// used to check if the user already has an account
	err := Db.Data.QueryRow("select id from users where email = ?", user.Email).Scan(&user.Id)

	if err == nil {
		return defs.NO_ID, errors.New("ERROR: user already has an account")
	}

	statement := "insert into users (username, password, email, date_created, tests_started, tests_completed, all_time_avg_wpm, all_time_avg_acc) values (?, ?, ?, ?, ?, ?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err = Db.Data.Prepare(statement)
	if err != nil {
		return defs.NO_ID, err
	}

	defer stmt.Close()

	encrypted_pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return defs.NO_ID, err
	}

	err = stmt.QueryRow(
		user.UserName,
		string(encrypted_pass),
		user.Email,
		time.Now(),
		defs.DEFAULT_NUM_TESTS,
		defs.DEFAULT_NUM_TESTS,
		defs.DEFAULT_WPM,
		defs.DEFAULT_ACC,
		).Scan(&user.Id)

	if err != nil {
		return defs.NO_ID, err
	}

	return user.Id, nil

}

func (Db *DataBase) UpdateUserCredentials(user m.User) error {
	statement := "UPDATE users SET username = ?, password = ?, email = ? WHERE id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.UserName, user.Password, user.Email, user.Id)

	return err
}

func (Db *DataBase) UpdateUserStats(user m.User) error {
	statement := "UPDATE users SET tests_started = ?, tests_completed = ?, all_time_avg_wpm = ?, all_time_avg_acc = ? WHERE id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.TestsStarted, user.TestsCompleted, user.AllTimeAvgWPM, user.AllTimeAvgACC, user.Id)

	return err
}

func (Db *DataBase) DeleteUser(user_id int) error {
	statement := "DELETE from users where id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user_id)

	return err
}
