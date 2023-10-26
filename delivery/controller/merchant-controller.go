package controller

import (
	"github.com/gin-gonic/gin"
	"mnc-test/delivery/middleware"
	"mnc-test/model"
	"mnc-test/usecase"
)

type MerchantController struct {
	merchantUC usecase.MerchantUsecase
	gin        *gin.Engine
}

func (m MerchantController) CreateNewMerchant(c *gin.Context) {
	var merchant model.Merchant

	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	if err := m.merchantUC.CreateNewMerchant(merchant); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"Message": "Successfully create new merchant"})
}

func (m MerchantController) FindAllMerchants(c *gin.Context) {
	var merchants []model.Merchant

	merchants, err := m.merchantUC.FindAllMerchant()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"Data": merchants})
}

func (m MerchantController) Route() {
	merchantGroup := m.gin.Group("/app/merchants")
	{
		merchantGroup.GET("/list", middleware.AuthMiddleware(), m.FindAllMerchants)
		merchantGroup.POST("/create", middleware.AuthMiddleware(), m.CreateNewMerchant)
	}
}

func NewMerchantController(merchantuc usecase.MerchantUsecase, g *gin.Engine) *MerchantController {
	return &MerchantController{
		merchantUC: merchantuc,
		gin:        g,
	}
}
