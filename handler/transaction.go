package handler

import (
	"backer-api/helper"
	"backer-api/transaction"
	"backer-api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatCampaignTransactions(transactions)
	response := helper.APIResponse("Campaign's transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)

	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatUserTransactions(transactions)
	response := helper.APIResponse("User's transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidatorError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)

	if err != nil {
		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatTransaction(newTransaction)
	response := helper.APIResponse("Success to create transaction", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.ProcessPayment(input)

	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
