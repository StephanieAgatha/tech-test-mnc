package controller

import (
	"github.com/gin-gonic/gin"
	"mnc-test/model"
	"mnc-test/usecase"
)

type UserCredentialController struct {
	userCredUc usecase.UserCredentialUsecase
	gin        *gin.Engine
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

	c.JSON(200, gin.H{"Data": userToken})
}

func (u UserCredentialController) Logout(c *gin.Context) {
	// Extract the user's credentials or token from the request context
	// This could be done in a variety of ways depending on how you handle authentication
	var userCred model.UserCredentials
	if err := c.ShouldBindJSON(&userCred); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Bad JSON Format"})
		return
	}

	// Call the Logout method of the UserCredentialUsecase
	// This is just a placeholder and should be replaced with your actual implementation
	err := u.userCredUc.Logout(userCred)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	// Return a response to the client
	c.JSON(200, gin.H{"Message": "Successfully logged out"})
}

func (u UserCredentialController) Route() {
	authGroup := u.gin.Group("/auth")
	{
		authGroup.POST("/register", u.Register)
		authGroup.POST("/login", u.Login)
		authGroup.POST("/logout", u.Logout)
	}
}

func NewUserCredentialController(uc usecase.UserCredentialUsecase, g *gin.Engine) *UserCredentialController {
	return &UserCredentialController{
		userCredUc: uc,
		gin:        g,
	}
}
