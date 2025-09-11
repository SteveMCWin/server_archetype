package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	m "github.com/SteveMCWin/archetype-common/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"server_archetype/mail"
	"server_archetype/models"
)

var Domain string
var JwtKey string

func SetUpRouter(domain, jwt_key string, db *models.DataBase) *gin.Engine {

	Domain = domain
	JwtKey = jwt_key

	router := gin.Default()

	api := router.Group("/api")

	{
		api.GET("/ping", HandleGetPing())

		api.POST("/signup", HandlePostSignup(db))
		api.GET("/verify", HandleGetVerify(db))
		api.POST("/login", HandlePostLogin(db))

		api.GET("/profile", HandleGetProfile(db))
		api.PUT("/profile/credentials", HandlePutProfileCredentials(db))
		api.POST("/profile/email", HandlePostProfileEmail(db))
		api.GET("/profile/update_email", HandleGetUpdateEmail(db))
		api.PUT("/profile/stats", HandlePutProfileStats(db))
		api.DELETE("/profile")

		api.GET("/quote", HandleGetQuote(db))
		api.GET("/words")

		api.GET("/leaderboard")
	}

	return router
}

func HandleGetPing() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}
}

func HandlePostSignup(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_email := c.PostForm("email")
		user_password := c.PostForm("password")
		user_name := c.PostForm("username")

		if user_email == "" || user_password == "" || user_name == "" {
			log.Println("Credentials are empty!!")
			log.Println("user_email:", user_email)
			log.Println("user_password:", user_password)
			log.Println("user_name:", user_name)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		log.Println("User's email:", user_email)
		if email_exists := db.EmailExists(user_email); email_exists == true {
			log.Println("You already have an account!")
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		token_val := CreateToken(user_email, user_password, user_name)

		new_mail := &mail.Mail{
			Recievers:    []string{user_email},
			Subject:      "Signup Verification",
			TempaltePath: "./templates/mail_register.html",
			ExtLink:      Domain + "/api/verify?token=" + strconv.Itoa(token_val) + "&email=" + user_email} // NOTE: the domain mustn't end with a '/'

		err := mail.SendMailHtml(new_mail)
		if err != nil {
			log.Println("FAILLLLED TO SEND MAILLLL")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func HandleGetVerify(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		token_str := c.Query("token")
		email := c.Query("email")

		token, err := strconv.Atoi(token_str)
		if err != nil {
			c.String(http.StatusInternalServerError, "ERROR: couldn't parse url parameters: "+err.Error())
			return
		}

		user_data, ok := signupTokens[token]
		if !ok {
			c.String(http.StatusOK, "This token has expired. The expiration time is 10 minutes. Please try signing up again.")
			return
		}

		if email != user_data.UserMail {
			c.String(http.StatusInternalServerError, "The email read from the url doesn't match the email the user provided in the CLI")
			return
		}

		new_user := m.User{
			Email:    user_data.UserMail,
			Password: user_data.Password,
			UserName: user_data.UserName,
		}

		_, err = db.CreateUser(&new_user)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating user: "+err.Error())
			return
		}

		delete(signupTokens, token)

		c.String(http.StatusOK, "Successfully created user :D You can now log in though your terminal!")
	}
}

func HandlePostLogin(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		user_email := c.PostForm("email")
		user_password := c.PostForm("password")

		if user_email == "" || user_password == "" {
			log.Println("Credentials are empty!!")
			log.Println("user_email:", user_email)
			log.Println("user_password:", user_password)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		user, err := db.AuthUser(user_email, user_password)
		if err != nil {
			log.Println("Error reading user data from database")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		// claims := jwt.MapClaims {
		// 	"sub": user.Id,
		// }

		claims := jwt.RegisteredClaims{
			Subject: strconv.FormatUint(user.Id, 10),
		}

		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := t.SignedString([]byte(JwtKey))
		if err != nil {
			log.Println("Error creating a JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			// "user":  user,
			"token": token,
		})
	}
}

func HandleGetProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			log.Println("No JWT token provided?!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user_id, err := getUserIdFromJWT(auth_header)
		if err != nil {
			log.Println("Error getting user_id from JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		user, err := db.ReadUser(user_id)

		if err != nil {
			log.Println("Error loading user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"profile": user,
		})

	}
}

func getUserIdFromJWT(auth_header string) (uint64, error) {

	jwt_string := strings.TrimPrefix(auth_header, "Bearer ")

	token, err := verifyJWT(jwt_string)
	if err != nil {
		log.Println("Invalid JWT:", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Couldn't get claims from jwt")
		return 0, errors.New("Error getting jwt claims")
	}

	user_id_string, err := claims.GetSubject()
	if err != nil {
		log.Println("Couldn't get sub from jwt:", err)
		return 0, err
	}

	user_id, err := strconv.ParseUint(user_id_string, 10, 64)
	if err != nil {
		log.Println("Error converting user_id_string to user_id:", err)
		return 0, err
	}

	return user_id, nil
}

func verifyJWT(jwt_string string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwt_string, func(token *jwt.Token) (any, error) {
		return []byte(JwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func HandlePutProfileCredentials(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_name := c.PostForm("username")
		user_password := c.PostForm("password")

		if user_name == "" || user_password == "" {
			log.Println("Credentials are empty!!")
			log.Println("user_user_name:", user_name)
			log.Println("user_password:", user_password)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			log.Println("No JWT token provided?!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user_id, err := getUserIdFromJWT(auth_header)
		if err != nil {
			log.Println("Error getting user_id from JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		user, err := db.ReadUser(user_id)
		if err != nil {
			log.Println("Error getting old user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		err = db.UpdateUserCredentials(&user)
		if err != nil {
			log.Println("Error updating user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

	}
}

func HandlePostProfileEmail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		new_email := c.PostForm("email")

		if new_email == "" {
			log.Println("Credentials are empty!!")
			log.Println("user_email:", new_email)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		log.Println("User's email:", new_email)
		if email_exists := db.EmailExists(new_email); email_exists == true {
			log.Println("You already have an account!")
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		token_val := CreateToken(new_email, "", "")

		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			log.Println("No JWT token provided?!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user_id, err := getUserIdFromJWT(auth_header)
		if err != nil {
			log.Println("Error getting user_id from JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		new_mail := &mail.Mail{
			Recievers:    []string{new_email},
			Subject:      "Signup Verification",
			TempaltePath: "./templates/mail_register.html",
			ExtLink:      Domain + "/api/profile/update_email?token=" + strconv.Itoa(token_val) + "&email=" + new_email + "&user_id="+strconv.FormatUint(user_id, 10)} // NOTE: the domain mustn't end with a '/'

		err = mail.SendMailHtml(new_mail)
		if err != nil {
			log.Println("FAILLLLED TO SEND MAILLLL")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func HandleGetUpdateEmail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		token_str := c.Query("token")
		email := c.Query("email")
		user_id, err := strconv.Atoi(c.Query("user_id"))
		if err != nil {
			c.String(http.StatusInternalServerError, "ERROR: couldn't parse user id from usl")
			return
		}

		token, err := strconv.Atoi(token_str)
		if err != nil {
			c.String(http.StatusInternalServerError, "ERROR: couldn't parse url parameters: "+err.Error())
			return
		}

		user_data, ok := signupTokens[token]
		if !ok {
			c.String(http.StatusOK, "This token has expired. The expiration time is 10 minutes. Please try signing up again.")
			return
		}

		if email != user_data.UserMail {
			c.String(http.StatusInternalServerError, "The email read from the url doesn't match the email the user provided in the CLI")
			return
		}

		db.UpdateUserEmail(user_id, email)

		delete(signupTokens, token)

		c.String(http.StatusOK, "Successfully updated email!")
	}
}

func HandlePutProfileStats(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user_data m.User
		if err := c.BindJSON(&user_data); err != nil {
			log.Println("Error binding json:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			log.Println("No JWT token provided?!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user_id, err := getUserIdFromJWT(auth_header)
		if err != nil {
			log.Println("Error getting user_id from JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		user_data.Id = user_id

		err = db.UpdateUserStats(&user_data)
		if err != nil {
			log.Println("Error updating user stats:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

	}
}

func HandleDeleteProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			log.Println("No JWT token provided?!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user_id, err := getUserIdFromJWT(auth_header)
		if err != nil {
			log.Println("Error getting user_id from JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		err = db.DeleteUser(user_id)
		if err != nil {
			log.Println("Error deleting the user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func HandleGetQuote(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
