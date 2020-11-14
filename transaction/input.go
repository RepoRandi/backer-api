package transaction

type GetCampaignTransactionsInput struct {
	ID int `uri:"id" binding:"required"`
}
