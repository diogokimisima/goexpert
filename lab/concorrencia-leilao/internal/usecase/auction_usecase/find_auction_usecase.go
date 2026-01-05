package auction_usecase

import (
	"context"

	"github.com/diogokimisima/fullcycle-auction/internal/entity/auction_entity"
	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
)

func (au *AuctionUseCase) FindAuctionById(
	ctx context.Context,
	id string) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.auctionRepository.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}

	auctionOutput := &AuctionOutputDTO{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		TimeStamp:   auctionEntity.TimeStamp,
	}

	return auctionOutput, nil
}

func (au *AuctionUseCase) FindAuctions(
	ctx context.Context,
	status AuctionStatus,
	category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {

	auctionEntities, err := au.auctionRepository.FindAuctions(ctx,
		auction_entity.AuctionStatus(status),
		category,
		productName)
	if err != nil {
		return nil, err
	}

	var auctionOutputs []AuctionOutputDTO
	for _, value := range auctionEntities {
		auctionOutput := AuctionOutputDTO{
			Id:          value.Id,
			ProductName: value.ProductName,
			Category:    value.Category,
			Description: value.Description,
			Condition:   ProductCondition(value.Condition),
			Status:      AuctionStatus(value.Status),
			TimeStamp:   value.TimeStamp,
		}
		auctionOutputs = append(auctionOutputs, auctionOutput)
	}

	return auctionOutputs, nil
}
