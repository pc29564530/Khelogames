package handlers

import (
	"khelogames/api/transactions"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HandlersServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
	txStore    *transactions.SQLStore
	r2Client   *s3.Client
}

func NewHandlerServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, txStore *transactions.SQLStore, r2Client *s3.Client) *HandlersServer {
	return &HandlersServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, txStore: txStore, r2Client: r2Client}
}
