package main

import (
	"os"
	"log"

	"github.com/joho/godotenv"

	"server_archetype/mail"
	"server_archetype/models"
	"server_archetype/handlers"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load the .env")
	}

	domain := os.Getenv("DOMAIN")

	if domain == "" {
		log.Fatal("Couldn't load .env variables: domain")
	}

	jwt_key := os.Getenv("JWT_KEY")

	if jwt_key == "" {
		log.Fatal("Couldn't load .env variables: jwt_key")
	}

	mail_pass := os.Getenv("GMAIL_APP_PASS")
	mail_sender := os.Getenv("MAIL_SENDER")

	if mail_pass == "" || mail_sender == "" {
		log.Fatal("Couldn't load .env variables: mail")
	}

	mail.InitMail(mail_pass, mail_sender)


	db := &models.DataBase{}
    db.InitDatabase()

	// MostEloLb := &models.LeaderBoard{}
	// MostEloLb.InitLeaderBoard()
	// MostEloLb.RunLeaderBoard(db)

	handler := handlers.SetUpRouter(domain, jwt_key, db)

	handler.Run(":5000")
}

