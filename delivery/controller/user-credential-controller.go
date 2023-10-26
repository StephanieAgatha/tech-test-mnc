package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mnc-test/model"
	"mnc-test/usecase"
	"time"
)

type UserCredentialController struct {
	userCredUc usecase.UserCredentialUsecase
	gin        *gin.Engine
	redisC     *redis.Client
	log        *zap.Logger
}

func (u UserCredentialController) Register(c *gin.Context) {
	var userCred model.UserCredentials

	if err := c.ShouldBindJSON(&userCred); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Bad JSON Format"})
		return
	}

	if err := u.userCredUc.Register(userCred); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	//log goes hereeee
	if u.log != nil {
		u.log.Info("New Customer Has Been Created",
			zap.String("Custormer Name", userCred.Name),
			zap.String("Customer Email", userCred.Email))
	} else {
		fmt.Println("Logger is not initialized")
	}

	c.JSON(200, gin.H{"Message": "Successfully Register"})
}

func (u UserCredentialController) Login(c *gin.Context) {
	var userCred model.UserCredentials

	if err := c.ShouldBindJSON(&userCred); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Bad JSON Format"})
		return
	}

	_, err := u.userCredUc.FindUserEMail(userCred.Email)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	userToken, err := u.userCredUc.Login(userCred)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	// save email dan token ke redis
	ctx := context.Background()

	err = u.redisC.Set(ctx, "user:email:"+userCred.Email, userToken, 24*time.Hour).Err()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": "Failed to save data to Redis"})
		return
	}

	//log
	if u.log != nil {
		u.log.Info("Customer has been logged in",
			zap.String("Customer Email", userCred.Email))
	} else {
		fmt.Println("Logger is not initialized")
	}

	c.JSON(200, gin.H{"Data": userToken})

}

func (u UserCredentialController) Logout(c *gin.Context) {
	ctx := context.Background()

	var data struct {
		CustEmail string `json:"customer_email"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Data"})
		return
	}

	//delete token nya
	err := u.redisC.Del(ctx, "user:email:"+data.CustEmail).Err()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": "Failed to logout"})
		return
	}

	//set log
	if u.log != nil {
		u.log.Info("Customer has been logged out",
			zap.String("Customer Email", data.CustEmail))
	} else {
		fmt.Println("Logger is not initialized")
	}

	c.JSON(200, gin.H{"Message": "Logout successful"})
}

func (u UserCredentialController) Route() {
	authGroup := u.gin.Group("/auth")
	{
		authGroup.POST("/register", u.Register)
		authGroup.POST("/login", u.Login)
		authGroup.POST("/logout", u.Logout)
	}
}

func NewUserCredentialController(uc usecase.UserCredentialUsecase, g *gin.Engine, redisC *redis.Client, log *zap.Logger) *UserCredentialController {
	return &UserCredentialController{
		userCredUc: uc,
		gin:        g,
		redisC:     redisC,
		log:        log,
	}
}
