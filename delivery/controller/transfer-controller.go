package controller

import (
	"github.com/gin-gonic/gin"
	"mnc-test/delivery/middleware"
	"mnc-test/model"
	"mnc-test/usecase"
)

type TransferController struct {
	tfUsecase usecase.TransferUsecase
	gin       *gin.Engine
}

func (t TransferController) TransferAccNumbToAccNumb(c *gin.Context) {
	var tf model.TransferRequest

	if err := c.ShouldBindJSON(&tf); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	if err := t.tfUsecase.MakeTransferAccNumbToAccNumb(tf.TransactionID, tf.SenderAccountNumber, tf.ReceiverAccountNumber, tf.Amount); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"Message": "Successfully", "Data": tf})
}

func (t TransferController) GetIncomingMoneyHandler(c *gin.Context) {
	var tf model.TransferHistoryIncome

	if err := c.ShouldBindJSON(&tf); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	tfHistoryIncome, err := t.tfUsecase.GetIncomingMoney(tf.CustomerID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// Create a slice to hold the response data
	responseData := make([]model.TransferHistoryIncomeResponse, len(tfHistoryIncome))

	//for ranggeeee for give a response hehe
	for i, transfer := range tfHistoryIncome {
		responseData[i] = model.TransferHistoryIncomeResponse{
			ID:                    transfer.ID,
			TransactionID:         transfer.TransactionID,
			SenderAccountNumber:   transfer.SenderAccountNumber,
			ReceiverAccountNumber: transfer.ReceiverAccountNumber,
			Amount:                transfer.Amount,
			TransferTime:          transfer.TransferTimeStamp.Format("2006-01-02 15:04:05"), //beautiful layout :))
			SenderBankName:        transfer.SenderBankName,
			SenderBankId:          transfer.SenderBankId,
			ReceiverBankName:      transfer.ReceiverBankName,
			ReceiverBankId:        transfer.ReceiverBankId,
		}
	}

	c.JSON(200, gin.H{"Data": responseData})
}

func (t TransferController) GetOutcomeMoneyHandler(c *gin.Context) {
	var tf model.TransferHistoryIncome

	if err := c.ShouldBindJSON(&tf); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	tfHistoryIncome, err := t.tfUsecase.GetOutcomeMoney(tf.CustomerID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// Create a slice to hold the response data
	responseData := make([]model.TransferHistoryIncomeResponse, len(tfHistoryIncome))

	//for ranggeeee for give a response hehe
	for i, transfer := range tfHistoryIncome {
		responseData[i] = model.TransferHistoryIncomeResponse{
			ID:                    transfer.ID,
			TransactionID:         transfer.TransactionID,
			SenderAccountNumber:   transfer.SenderAccountNumber,
			ReceiverAccountNumber: transfer.ReceiverAccountNumber,
			Amount:                transfer.Amount,
			TransferTime:          transfer.TransferTimeStamp.Format("2006-01-02 15:04:05"), //beautiful layout :))
			SenderBankName:        transfer.SenderBankName,
			SenderBankId:          transfer.SenderBankId,
			ReceiverBankName:      transfer.ReceiverBankName,
			ReceiverBankId:        transfer.ReceiverBankId,
		}
	}

	c.JSON(200, gin.H{"Data": responseData})
}

func (t TransferController) Route() {
	tfGroup := t.gin.Group("/app/transfer")
	{
		tfGroup.POST("/create/account", middleware.AuthMiddleware(), t.TransferAccNumbToAccNumb)
		tfGroup.POST("/list/income", middleware.AuthMiddleware(), t.GetIncomingMoneyHandler)
		tfGroup.POST("/list/outcome", middleware.AuthMiddleware(), t.GetOutcomeMoneyHandler)
	}
}
func NewTransferController(tfUsecase usecase.TransferUsecase, g *gin.Engine) *TransferController {
	return &TransferController{
		tfUsecase: tfUsecase,
		gin:       g,
	}
}

//func (t TransferController) TransferFromPhoneNumbToPhoneNumbHandler(c *gin.Context) {
//	var tf model.TransferRequestPhoneNumb
//
//	if err := c.ShouldBindJSON(&tf); err != nil {
//		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
//		return
//	}
//
//	//generate uuid in here
//	tf.TransactionID = helper.GenerateUUID()
//
//	if err := t.tfUsecase.MakeTransferPhoneNumbToPhoneNumb(tf.TransactionID, tf.SenderPhoneNumber, tf.ReceiverPhoneNumber, tf.Amount); err != nil {
//		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
//		return
//	}
//
//	c.JSON(200, gin.H{"Message": "Successfully", "Data": tf})
//}
