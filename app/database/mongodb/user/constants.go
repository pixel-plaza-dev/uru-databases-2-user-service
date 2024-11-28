package user

import (
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb"
)

const (
	// UriKey is the key of the MongoDB host
	UriKey = "USER_SERVICE_MONGODB_HOST"

	// DbNameKey is the key of the MongoDB database name
	DbNameKey = "USER_SERVICE_MONGODB_NAME"
)

var (
	// userCollectionSingleFieldIndex is the single field indexes for the user collection
	userCollectionSingleFieldIndex = []*commonmongodb.SingleFieldIndex{
		commonmongodb.NewSingleFieldIndex(
			commonmongodb.FieldIndex{
				Name:  "username",
				Order: commonmongodb.Ascending,
			}, true,
		),
		commonmongodb.NewSingleFieldIndex(
			commonmongodb.FieldIndex{
				Name:  "uuid",
				Order: commonmongodb.Ascending,
			}, true,
		),
	}

	// UserCollection is the users collection in MongoDB
	UserCollection = commonmongodb.NewCollection(
		"User",
		&userCollectionSingleFieldIndex,
		nil,
	)

	// UserSharedIdentifierCollection is the user shared identifiers collection in MongoDB
	UserSharedIdentifierCollection = commonmongodb.NewCollection(
		"UserSharedIdentifier",
		nil,
		nil,
	)

	// UserEmailCollection is the user emails collection in MongoDB
	UserEmailCollection = commonmongodb.NewCollection("UserEmail", nil, nil)

	// UserPhoneNumberCollection is the user phone numbers collection in MongoDB
	UserPhoneNumberCollection = commonmongodb.NewCollection(
		"UserPhoneNumber",
		nil,
		nil,
	)

	// UserUsernameLogCollection is the user username log collection in MongoDB
	UserUsernameLogCollection = commonmongodb.NewCollection(
		"UserUsernameLog",
		nil,
		nil,
	)

	// UserHashedPasswordLogCollection is the user hashed password log collection in MongoDB
	UserHashedPasswordLogCollection = commonmongodb.NewCollection(
		"UserHashedPasswordLog",
		nil,
		nil,
	)
)
