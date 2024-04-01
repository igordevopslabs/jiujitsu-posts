package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/models"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /post

// @Summary     Create user
// @Description Create a new user to consume API
// @ID          create-user
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       post body models.User true "User create"
// @Success     200 {object} models.User
// @Router      /singup [post]
func Singup(c *gin.Context) {
	//pegar email/password a partir da requisição
	var body struct {
		Email    string
		Password string
	}

	if c.ShouldBindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to process body req",
		})

		return
	}

	//encryptar a senha (gerar hash)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate pasword hash",
		})

		return
	}

	//Criar o user no banco
	user := models.User{Email: body.Email, Password: string(hash)}
	result := config.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	//Responder o evento de criação do user
	c.JSON(200, gin.H{
		"User": user,
	})
}

// @Summary     Login user
// @Description Login a user to consume API
// @ID          login-user
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       post body models.User true "User log\in"
// @Success     200 {object} models.User
// @Router      /login [post]
func Login(c *gin.Context) {
	//pegar email e password atraves da requisição
	var body struct {
		Email    string
		Password string
	}

	if c.ShouldBindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to process body req",
		})

		return
	}

	//Selecionar o email correto do banco
	var user models.User
	config.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to find user",
		})

		return
	}

	//Comparar a password informada com o hash da password salva
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Failed to comparing password or password doesn't match",
		})

		return
	}

	//gerar o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to create token",
		})

		return
	}
	//devOlver o Token para ser consumido.
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
