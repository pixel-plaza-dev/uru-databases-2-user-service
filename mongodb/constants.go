package mongodb

import (
	commonMongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
)

var (
	// userCollectionSingleFieldIndex is the single field indexes for the user collection
	userCollectionSingleFieldIndex = []*commonMongodb.SingleFieldIndex{
		commonMongodb.NewSingleFieldIndex(commonMongodb.FieldIndex{Name: "username", Order: commonMongodb.Ascending}, true)}

	// UserCollection is the users collection in MongoDB
	UserCollection = commonMongodb.NewCollection("User", &userCollectionSingleFieldIndex, nil)

	// UserEmailCollection is the user emails collection in MongoDB
	UserEmailCollection = commonMongodb.NewCollection("UserEmail", nil, nil)

	// UserPhoneNumberCollection is the user phone numbers collection in MongoDB
	UserPhoneNumberCollection = commonMongodb.NewCollection("UserPhoneNumber", nil, nil)

	// UserUsernameLogCollection is the user username log collection in MongoDB
	UserUsernameLogCollection = commonMongodb.NewCollection("UserUsernameLog", nil, nil)

	// UserHashedPasswordLogCollection is the user hashed password log collection in MongoDB
	UserHashedPasswordLogCollection = commonMongodb.NewCollection("UserHashedPasswordLog", nil, nil)
)
