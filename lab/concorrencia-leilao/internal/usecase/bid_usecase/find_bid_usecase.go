package bid_usecase

import (
	"context"

	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
)

func (bu *BidUseCase) FindBidByAuctionId(
	ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError) {
	bidList, err := bu.BidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	var bidOutputDTOs []BidOutputDTO
	for _, value := range bidList {
		bidOutputDTOs = append(bidOutputDTOs, BidOutputDTO{
			Id:        value.Id,
			UserId:    value.UserId,
			AuctionId: value.AuctionId,
			Amount:    value.Amount,
			Timestamp: value.Timestamp,
		})
	}

	return bidOutputDTOs, nil
}

func (bu *BidUseCase) FindWinnigBidByAuctionId(
	ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError) {

	bidEntity, err := bu.BidRepository.FindWinnigBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	bidOutputDTO := &BidOutputDTO{
		Id:        bidEntity.Id,
		UserId:    bidEntity.UserId,
		AuctionId: bidEntity.AuctionId,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}

	return bidOutputDTO, nil
}
