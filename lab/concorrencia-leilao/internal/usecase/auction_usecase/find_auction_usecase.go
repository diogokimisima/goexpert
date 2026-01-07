package auction_usecase

import (
	"context"

	"github.com/diogokimisima/fullcycle-auction/internal/entity/auction_entity"
	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
	"github.com/diogokimisima/fullcycle-auction/internal/usecase/bid_usecase"
)

func (au *AuctionUseCase) FindAuctionById(
	ctx context.Context,
	id string) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.auctionRepositoryInterface.FindAuctionById(ctx, id)
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

	auctionEntities, err := au.auctionRepositoryInterface.FindAuctions(ctx,
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

func (au *AuctionUseCase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auction, err := au.auctionRepositoryInterface.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	auctionOutputDTO := AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		TimeStamp:   auction.TimeStamp,
	}

	bidWinnig, err := au.bidRepositoryInterface.FindWinnigBidByAuctionId(ctx, auction.Id)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutputDTO,
			Bid:     nil,
		}, nil
	}

	bidOutputDTO := &bid_usecase.BidOutputDTO{
		Id:        bidWinnig.Id,
		UserId:    bidWinnig.UserId,
		AuctionId: bidWinnig.AuctionId,
		Amount:    bidWinnig.Amount,
		Timestamp: bidWinnig.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		Bid:     bidOutputDTO,
	}, nil
}
