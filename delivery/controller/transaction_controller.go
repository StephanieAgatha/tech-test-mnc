package controller

import (
	"github.com/gin-gonic/gin"
	"mnc-test/delivery/middleware"
	"mnc-test/model"
	"mnc-test/usecase"
	"time"
)

type TransactionController struct {
	txUc usecase.TransactionUsecase
	g    *gin.Engine
}

func (t TransactionController) MakePaymentController(c *gin.Context) {
	var tx model.Transaction

	if err := c.ShouldBindJSON(&tx); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	if err := t.txUc.MakePayment(&tx); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	//custom response disini
	response := model.TransactionResponse{
		MerchantName: tx.MerchantName,
		Amount:       tx.Amount,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(200, gin.H{"Message": "Successfully", "Data": response})
}

func (t TransactionController) GetCustTransactionByIDHandler(c *gin.Context) {
	var custid struct {
		CustID int `json:"customer_id"`
	}

	if err := c.ShouldBindJSON(&custid); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Wrong JSON Format"})
		return
	}

	txs, err := t.txUc.GetCustomerTransaction(custid.CustID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	if len(txs) == 0 {
		c.JSON(404, gin.H{"Error": "No transactions found"})
		return
	}

	// Convert txs to TransactionResponse
	var txResponses []model.TransactionResponse
	for _, tx := range txs {
		txResponses = append(txResponses, model.TransactionResponse{
			MerchantName: tx.MerchantName,
			Amount:       tx.Amount,
			CreatedAt:    tx.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(200, gin.H{"Data": txResponses})
}

func (t TransactionController) Route() {
	txGroup := t.g.Group("/app/transaction")
	{
		txGroup.POST("/create", middleware.AuthMiddleware(), t.MakePaymentController)
		txGroup.POST("/list", middleware.AuthMiddleware(), t.GetCustTransactionByIDHandler)
	}
}

func NewTransactionController(txuc usecase.TransactionUsecase, g *gin.Engine) *TransactionController {
	return &TransactionController{
		txUc: txuc,
		g:    g,
	}
}
