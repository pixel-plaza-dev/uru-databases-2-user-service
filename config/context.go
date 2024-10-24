package config

import "time"

const (
	// MongoDbConnectionCtxTimeout is the timeout for the MongoDB connection
	MongoDbConnectionCtxTimeout = 60 * time.Second

	// MongoDbQueryCtxTimeout is the MongoDB query context
	MongoDbQueryCtxTimeout = 15 * time.Second

	// MongoDbTransactionCtxTimeout is the MongoDB transaction context
	MongoDbTransactionCtxTimeout = 60 * time.Second
)
