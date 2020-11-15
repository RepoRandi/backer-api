package transaction

import (
	"backer-api/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if campaign.User.ID != input.User.ID {
		return []Transaction{}, errors.New("Not owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
