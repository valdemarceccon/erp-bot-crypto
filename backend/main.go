package main

import (
	"errors"
	"os"
	"time"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
	// jwt "github.com/appleboy/gin-jwt/v2"`
)

type loginRequest struct {
	Username string `json:"username" bind:"username"`
	Password string `json:"password" bind:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type AuthController struct {
	userRepo repository.UserRepository
}

type UserClaim struct {
	jwt.RegisteredClaims
	Username string
	Id       uint32
}

func (ac *AuthController) loginHandler(c *gin.Context) {
	var loginVal loginRequest
	if err := c.ShouldBindWith(&loginVal, binding.JSON); err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err)
		return
	}
	user, err := ac.userRepo.ValidateUser(loginVal.Username, loginVal.Password)

	if err != nil {
		c.Error(errors.New("invalid user or password"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Username: user.Username,
		Id:       user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(time.Duration(2).Hours()))),
		},
	})

	signedToken, err := token.SignedString([]byte("somesecret"))

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: signedToken})
}

func main() {
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	r := gin.Default()

	if port == "" {
		port = "8000"
	}

	if jwtSecret == "" {
		jwtSecret = "some-secret"
	}

	userRepo := repository.NewUserRepositoryInMemory()
	userRepo.Create(&model.User{
		Name:     "valdemar",
		Username: "valdemar",
		Password: "123456",
		Email:    "valdemar.ceccon@gmail.com",
		Telegram: "qualquer coisa",
	})

	authControler := AuthController{
		userRepo: userRepo,
	}

	r.POST("/login", authControler.loginHandler)

	// r.POST("/login", func(c *gin.Context) {
	//  v := c.Request.Header.Get("Authorization")
	// 	v, err := ioutil.ReadAll(c.Request.Body)

	// 	if err != nil {
	// 		c.Error(err)
	// 	}

	// 	c.JSON(200, string(v))
	// })

	authGroup := r.Group("/auth")

	{
		authGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pog",
			})
		})

	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

type Resp struct {
	Message string `json:"message"`
}
