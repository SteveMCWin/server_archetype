package handlers

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"server_archetype/defs"
	"time"
)

type UserCreationData struct {
	UserMail string
	Password string
	UserName string
}

var signupTokens map[int]UserCreationData

func init() {
	signupTokens = make(map[int]UserCreationData)
}

func CreateToken(user_mail, user_password, user_name string) int {
	var token_val int

	for {
		token_val = generateCode()
		if _, exists := signupTokens[token_val]; exists == false {
			break
		}
	}

	signupTokens[token_val] = UserCreationData{
		UserMail: user_mail,
		Password: user_password,
		UserName: user_name,
	}

	timer := time.NewTimer(defs.MAIL_VALIDATION_TIME)

	go func(val int) {
		<-timer.C
		if _, exists := signupTokens[val]; exists == true {
			delete(signupTokens, val)
		}
	}(token_val)

	return token_val
}

func generateCode() int {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return int(binary.BigEndian.Uint32(b)%1000000)
}
