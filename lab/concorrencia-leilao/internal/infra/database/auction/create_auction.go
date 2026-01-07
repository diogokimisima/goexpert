package auction

import (
	"context"
	"os"
	"time"

	"github.com/diogokimisima/fullcycle-auction/configuration/logger"
	"github.com/diogokimisima/fullcycle-auction/internal/entity/auction_entity"
	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	TimeStamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auction *auction_entity.Auction) *internal_error.InternalError {

	auctionMongo := &AuctionEntityMongo{
		id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		TimeStamp:   auction.TimeStamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionMongo)
	if err != nil {
		return internal_error.NewInternalServerError(
			"error trying to create auction: " + err.Error())
	}

	go func() {
		time.Sleep(getAuctionInterval())
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		filter := bson.M{"_id": auction.Id}

		_, err := ar.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("error trying to update auction status: ", err)
			return
		}
	}()

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInverval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInverval)

	if err != nil {
		return time.Minute * 5
	}

	return duration
}
