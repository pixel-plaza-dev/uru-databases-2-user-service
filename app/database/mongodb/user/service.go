package user

import (
	"context"
	"errors"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb"
	commonmongodbuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb/model/user"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Database struct {
	database    *mongo.Database
	collections *map[string]*commonmongodb.Collection
	client      *mongo.Client
	authClient  pbauth.AuthClient
}

// NewDatabase creates a new MongoDB user database handler
func NewDatabase(
	client *mongo.Client,
	databaseName string,
	authClient pbauth.AuthClient,
) (database *Database, err error) {
	// Get the user service database
	userServiceDb := client.Database(databaseName)

	// Create map of collections
	collections := make(map[string]*commonmongodb.Collection)

	for _, collection := range []*commonmongodb.Collection{
		UserCollection,
		UserEmailCollection,
		UserPhoneNumberCollection,
		UserUsernameLogCollection,
		UserHashedPasswordLogCollection,
	} {
		// Create the collection
		collections[collection.Name] = collection
		if _, err = collection.CreateCollection(userServiceDb); err != nil {
			return nil, err
		}
	}

	return &Database{
		client: client, database: userServiceDb, collections: &collections,
		authClient: authClient,
	}, nil
}

// Database returns the MongoDB users database
func (d *Database) Database() *mongo.Database {
	return d.database
}

// GetCollection returns a collection
func (d *Database) GetCollection(collection *commonmongodb.Collection) *mongo.Collection {
	return d.database.Collection(collection.Name)
}

// NewUserUsernameLog creates a new user username log object
func (d *Database) NewUserUsernameLog(
	userId *primitive.ObjectID,
	username string,
) commonmongodbuser.UserUsernameLog {
	return commonmongodbuser.UserUsernameLog{
		ID:         primitive.NewObjectID(),
		UserID:     *userId,
		Username:   username,
		AssignedAt: time.Now(),
	}
}

// NewUserHashedPasswordLog creates a new user hashed password log object
func (d *Database) NewUserHashedPasswordLog(
	userId *primitive.ObjectID,
	hashedPassword string,
) commonmongodbuser.UserHashedPasswordLog {
	return commonmongodbuser.UserHashedPasswordLog{
		ID:             primitive.NewObjectID(),
		UserID:         *userId,
		HashedPassword: hashedPassword,
		AssignedAt:     time.Now(),
	}
}

// NewUserEmail creates a new user email object
func (d *Database) NewUserEmail(
	userId *primitive.ObjectID,
	email string,
) commonmongodbuser.UserEmail {
	return commonmongodbuser.UserEmail{
		ID:         primitive.NewObjectID(),
		UserID:     *userId,
		Email:      email,
		AssignedAt: time.Now(),
	}
}

// NewUserPhoneNumber creates a new user phone number object
func (d *Database) NewUserPhoneNumber(
	userId *primitive.ObjectID,
	phoneNumber string,
) commonmongodbuser.UserPhoneNumber {
	return commonmongodbuser.UserPhoneNumber{
		ID:          primitive.NewObjectID(),
		UserID:      *userId,
		PhoneNumber: phoneNumber,
		AssignedAt:  time.Now(),
	}
}

// CreateUserHashedPasswordLog creates a new user hashed password log
func (d *Database) CreateUserHashedPasswordLog(
	ctx context.Context,
	userId *primitive.ObjectID,
	hashedPassword string,
) error {
	// Create the UserHashedPasswordLog object
	userHashedPasswordLog := d.NewUserHashedPasswordLog(
		userId,
		hashedPassword,
	)

	// Insert user hashed password log
	_, err := d.GetCollection(UserHashedPasswordLogCollection).InsertOne(
		ctx,
		userHashedPasswordLog,
	)
	return err
}

// CreateUserUsernameLog creates a new user username log and inserts it into the database
func (d *Database) CreateUserUsernameLog(
	ctx context.Context,
	userId *primitive.ObjectID,
	username string,
) error {
	// Create the UserUsernameLog object
	userUsernameLog := d.NewUserUsernameLog(userId, username)

	// Insert user username log
	_, err := d.GetCollection(UserUsernameLogCollection).InsertOne(
		ctx,
		userUsernameLog,
	)
	return err
}

// InsertUserEmail inserts a user email into the database
func (d *Database) InsertUserEmail(
	ctx context.Context,
	userEmail *commonmongodbuser.UserEmail,
) error {
	_, err := d.GetCollection(UserEmailCollection).InsertOne(ctx, userEmail)
	return err
}

// CreateUserEmail creates a new user email and inserts it into the database
func (d *Database) CreateUserEmail(
	ctx context.Context,
	userId *primitive.ObjectID,
	email string,
) error {
	// Create the UserEmail object
	userEmail := d.NewUserEmail(userId, email)

	// Insert user email
	err := d.InsertUserEmail(ctx, &userEmail)
	return err
}

// InsertUserPhoneNumber inserts a user phone number into the database
func (d *Database) InsertUserPhoneNumber(
	ctx context.Context,
	userPhoneNumber *commonmongodbuser.UserPhoneNumber,
) error {
	_, err := d.GetCollection(UserPhoneNumberCollection).InsertOne(
		ctx,
		userPhoneNumber,
	)
	return err
}

// CreateUserPhoneNumber creates a new user phone number and inserts it into the database
func (d *Database) CreateUserPhoneNumber(
	ctx context.Context,
	userId *primitive.ObjectID,
	phoneNumber string,
) error {
	// Create the UserPhoneNumber object
	userPhoneNumber := d.NewUserPhoneNumber(userId, phoneNumber)

	// Insert user phone number
	err := d.InsertUserPhoneNumber(ctx, &userPhoneNumber)
	return err
}

// InsertUser inserts a user into the database
func (d *Database) InsertUser(
	user *commonmongodbuser.User,
	userEmail *commonmongodbuser.UserEmail,
	userPhoneNumber *commonmongodbuser.UserPhoneNumber,
) error {
	// Run the transaction
	err := commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Insert user
			if _, err := d.GetCollection(UserCollection).InsertOne(
				sc,
				user,
			); err != nil {
				return err
			}

			// Insert user email
			if err := d.InsertUserEmail(sc, userEmail); err != nil {
				return err
			}

			// Insert user phone number
			if err := d.InsertUserPhoneNumber(
				sc,
				userPhoneNumber,
			); err != nil {
				return err
			}

			// Create a new user hashed password log
			if err := d.CreateUserHashedPasswordLog(
				sc,
				&user.ID,
				user.HashedPassword,
			); err != nil {
				return err
			}

			// Create a new user username log
			err := d.CreateUserUsernameLog(sc, &user.ID, user.Username)
			return err
		},
	)
	return err
}

// FindUser finds a user
func (d *Database) FindUser(
	ctx context.Context,
	filter interface{},
	projection interface{},
	sort interface{},
) (*commonmongodbuser.User, error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Add not deleted filter
	filter = bson.M{
		"$and": []interface{}{
			filter,
			bson.M{"deleted_at": bson.M{"$exists": false}},
		},
	}

	// Initialize the user variable
	user := &commonmongodbuser.User{}

	// Find the user
	err := d.GetCollection(UserCollection).FindOne(
		ctx,
		filter,
		findOptions,
	).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds a user by username
func (d *Database) FindUserByUsername(
	ctx context.Context,
	username string,
	projection interface{},
	sort interface{},
) (user *commonmongodbuser.User, err error) {
	// Check if the username is empty
	if username == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Find the user
	return d.FindUser(ctx, bson.M{"username": username}, projection, sort)
}

// FindUserByUserId finds a user by the user ID
func (d *Database) FindUserByUserId(
	ctx context.Context,
	userId string,
	projection interface{},
	sort interface{},
) (user *commonmongodbuser.User, err error) {
	// Check if the user ID is empty
	if userId == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return nil, err
	}

	// Find the user
	return d.FindUser(ctx, bson.M{"_id": *userObjectId}, projection, sort)
}

// GetUserHashedPassword gets the user's hashed password
func (d *Database) GetUserHashedPassword(
	ctx context.Context,
	username string,
) (user *commonmongodbuser.User, err error) {
	// Check if the username is empty
	if username == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Find the user
	return d.FindUserByUsername(
		ctx,
		username,
		bson.M{"_id": 1, "hashed_password": 1, "uuid": 1},
		nil,
	)
}

// GetUsernameByUserId gets the username by the user ID
func (d *Database) GetUsernameByUserId(
	ctx context.Context,
	userId string,
) (username string, err error) {
	// Find the user
	user, err := d.FindUserByUserId(ctx, userId, bson.M{"username": 1}, nil)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

// GetUserIdByUsername gets the user ID by the username
func (d *Database) GetUserIdByUsername(
	ctx context.Context,
	username string,
) (userId string, err error) {
	// Find the user
	user, err := d.FindUserByUsername(ctx, username, bson.M{"_id": 1}, nil)
	if err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

// UsernameExists checks if the username exists
func (d *Database) UsernameExists(
	ctx context.Context,
	username string,
) (exists bool, err error) {
	// Find the user
	user, err := d.FindUserByUsername(ctx, username, bson.M{"_id": 1}, nil)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		return false, err
	}
	return user != nil, nil
}

// UpdateUserByUserId updates a user by the user ID
func (d *Database) UpdateUserByUserId(
	ctx context.Context,
	userId string,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return nil, err
	}

	// Create the filter
	filter := bson.M{"_id": *userObjectId}

	// Update the user
	result, err = d.GetCollection(UserCollection).UpdateOne(
		ctx,
		filter,
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUserUsername updates the user username
func (d *Database) UpdateUserUsername(userId string, username string) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Update the user username
			if _, err = d.GetCollection(UserCollection).UpdateOne(
				sc,
				bson.M{"_id": *userObjectId},
				bson.M{"username": username},
			); err != nil {
				return err
			}

			// Create a new user username log
			err = d.CreateUserUsernameLog(sc, userObjectId, username)
			return err
		},
	)
	return err
}

// UpdateUserPassword updates the user password
func (d *Database) UpdateUserPassword(
	grpcCtx context.Context,
	userId string,
	hashedPassword string,
) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Update the user password
			if _, err = d.GetCollection(UserCollection).UpdateOne(
				sc,
				bson.M{"_id": *userObjectId},
				bson.M{"hashed_password": hashedPassword},
			); err != nil {
				return err
			}

			// Create a new user hashed password log
			if err = d.CreateUserHashedPasswordLog(
				sc,
				userObjectId,
				hashedPassword,
			); err != nil {
				return err
			}

			// Revoke all user's refresh tokens
			_, err = d.authClient.RevokeRefreshTokens(
				grpcCtx,
				&emptypb.Empty{},
			)

			return nil
		},
	)
	return err
}

// UpdateUserProfile updates a user
func (d *Database) UpdateUserProfile(
	ctx context.Context,
	userId string,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	return d.UpdateUserByUserId(ctx, userId, update)
}

// GetUserProfile gets the user's profile
func (d *Database) GetUserProfile(
	ctx context.Context,
	username string,
) (user *commonmongodbuser.User, err error) {
	return d.FindUserByUsername(
		ctx,
		username,
		bson.M{"first_name": 1, "last_name": 1, "birthdate": 1},
		nil,
	)
}

// FindUserPhoneNumber finds a user's phone number
func (d *Database) FindUserPhoneNumber(
	ctx context.Context,
	filter interface{},
	projection interface{},
	sort interface{},
) (*commonmongodbuser.UserPhoneNumber, error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Initialize the userPhoneNumber variable
	userPhoneNumber := &commonmongodbuser.UserPhoneNumber{}

	// Find the user's phone number
	err := d.GetCollection(UserPhoneNumberCollection).FindOne(
		ctx,
		filter,
		findOptions,
	).Decode(userPhoneNumber)
	if err != nil {
		return nil, err
	}
	return userPhoneNumber, nil
}

// GetUserPhoneNumber gets the user's phone number
func (d *Database) GetUserPhoneNumber(
	ctx context.Context,
	userId string,
) (phoneNumber string, err error) {
	// Check if the user ID is empty
	if userId == "" {
		return "", mongo.ErrNoDocuments
	}

	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return "", err
	}

	// Create the find options with the most recent document based on the given field
	sort := bson.M{"assigned_at": -1}

	// Find the user's phone number
	var userPhoneNumber *commonmongodbuser.UserPhoneNumber
	userPhoneNumber, err = d.FindUserPhoneNumber(
		ctx,
		bson.M{"user_id": userObjectId},
		bson.M{"phone_number": 1},
		sort,
	)
	if err != nil {
		return "", err
	}
	return userPhoneNumber.PhoneNumber, nil
}

// UpdateUserPhoneNumber updates the user's phone number
func (d *Database) UpdateUserPhoneNumber(
	userId string,
	phoneNumber string,
) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Revoke the user's phone number
			if _, err = d.GetCollection(UserPhoneNumberCollection).UpdateOne(
				sc,
				bson.M{
					"user_id":    *userObjectId,
					"revoked_at": bson.M{"$exists": false},
				},
				bson.M{"revoked_at": time.Now()},
			); err != nil {
				return err
			}

			// Create a new user phone number
			if err = d.CreateUserPhoneNumber(
				sc,
				userObjectId,
				phoneNumber,
			); err != nil {
				return err
			}

			return nil
		},
	)
	return err
}

// DeleteUser deletes a user
func (d *Database) DeleteUser(grpcCtx context.Context, userId string) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Update the user deleted at field
			if _, err = d.GetCollection(UserCollection).UpdateOne(
				sc,
				bson.M{"_id": *userObjectId},
				bson.M{"deleted_at": time.Now()},
			); err != nil {
				return err
			}

			// Revoke all user's refresh tokens
			_, err = d.authClient.RevokeRefreshTokens(
				grpcCtx,
				&emptypb.Empty{},
			)

			return nil
		},
	)
	return err
}

// AddUserEmail adds an email to a user
func (d *Database) AddUserEmail(
	ctx context.Context,
	userId string,
	email string,
) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Check if the user email already exists
			_, err = d.FindUserEmail(
				ctx,
				bson.M{
					"user_id":    *userObjectId,
					"email":      email,
					"revoked_at": bson.M{"$exists": false},
				},
				bson.M{"_id": 1},
				nil,
			)
			if err == nil || !errors.Is(mongo.ErrNoDocuments, err) {
				return EmailAlreadyExistsError
			}

			// Create the new user email
			err = d.CreateUserEmail(ctx, userObjectId, email)
			return err
		},
	)
	return err
}

// UpdateUserPrimaryEmail updates the user's primary email
func (d *Database) UpdateUserPrimaryEmail(userId string, email string) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Update the current user's primary email as not primary
			if _, err = d.GetCollection(UserEmailCollection).UpdateOne(
				sc,
				bson.M{
					"user_id":    *userObjectId,
					"is_primary": true,
					"revoked_at": bson.M{"$exists": false},
				},
				bson.M{"$set": bson.M{"is_primary": false}},
			); err != nil {
				return err
			}

			// Update the new user's primary email
			if _, err = d.GetCollection(UserEmailCollection).UpdateOne(
				sc,
				bson.M{
					"user_id":    *userObjectId,
					"email":      email,
					"revoked_at": bson.M{"$exists": false},
				},
				bson.M{"$set": bson.M{"is_primary": true}},
			); err != nil {
				return err
			}

			return nil
		},
	)
	return err
}

// DeleteUserEmail deletes an email from a user
func (d *Database) DeleteUserEmail(
	ctx context.Context,
	userId string,
	email string,
) error {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Revoke the user's email
	_, err = d.GetCollection(UserEmailCollection).UpdateOne(
		ctx,
		bson.M{
			"user_id":    *userObjectId,
			"email":      email,
			"is_primary": false,
			"revoked_at": bson.M{"$exists": false},
		},
		bson.M{"revoked_at": time.Now()},
	)
	return err
}

// FindUserEmail finds a user's email
func (d *Database) FindUserEmail(
	ctx context.Context,
	filter interface{},
	projection interface{},
	sort interface{},
) (*commonmongodbuser.UserEmail, error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Initialize the userEmail variable
	userEmail := &commonmongodbuser.UserEmail{}

	// Find the user's email
	err := d.GetCollection(UserEmailCollection).FindOne(
		ctx,
		filter,
		findOptions,
	).Decode(userEmail)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// FindUserEmailByEmail finds a user's email by email
func (d *Database) FindUserEmailByEmail(
	ctx context.Context,
	userId primitive.ObjectID,
	email string,
	projection interface{},
	sort interface{},
) (userEmail *commonmongodbuser.UserEmail, err error) {
	// Check if the user email already exists
	userEmail, err = d.FindUserEmail(
		ctx,
		bson.M{
			"user_id":    userId,
			"email":      email,
			"revoked_at": bson.M{"$exists": false},
		},
		projection,
		sort,
	)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// FindUserEmailPrimaryEmail finds a user's primary email
func (d *Database) FindUserEmailPrimaryEmail(
	ctx context.Context,
	userId primitive.ObjectID,
	projection interface{},
	sort interface{},
) (userEmail *commonmongodbuser.UserEmail, err error) {
	// Find the user's primary email
	userEmail, err = d.FindUserEmail(
		ctx,
		bson.M{
			"user_id":    userId,
			"is_primary": true,
			"revoked_at": bson.M{"$exists": false},
		},
		projection,
		sort,
	)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// GetUserPrimaryEmail gets the user's primary email
func (d *Database) GetUserPrimaryEmail(
	ctx context.Context,
	userId string,
) (email string, err error) {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return "", err
	}

	// Find the user's primary email
	var userEmail *commonmongodbuser.UserEmail
	userEmail, err = d.FindUserEmailPrimaryEmail(
		ctx,
		*userObjectId,
		bson.M{"email": 1},
		nil,
	)
	if err != nil {
		return "", err
	}
	return userEmail.Email, nil
}

// FindUserActiveEmails finds the user's active emails
func (d *Database) FindUserActiveEmails(
	ctx context.Context,
	userId string,
	projection interface{},
	sort interface{},
) (emails []*commonmongodbuser.UserEmail, err error) {
	// Convert the user ID to an object ID
	userObjectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return nil, err
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOptions(projection, sort, 0, 0)

	// Find the user's active emails
	cur, err := d.GetCollection(UserEmailCollection).Find(
		ctx,
		bson.M{
			"user_id":    *userObjectId,
			"revoked_at": bson.M{"$exists": false},
		},
		findOptions,
	)
	if err != nil {
		return nil, err
	}

	// Iterate through the cursor
	for cur.Next(ctx) {
		// Decode the user email
		var userEmail commonmongodbuser.UserEmail
		if err = cur.Decode(&userEmail); err != nil {
			return nil, err
		}
		emails = append(emails, &userEmail)
	}

	return emails, nil
}

// GetUserActiveEmails gets the user's active emails
func (d *Database) GetUserActiveEmails(
	ctx context.Context,
	userId string,
) (userActiveEmails []string, err error) {
	// Find the user's active emails
	var emails []*commonmongodbuser.UserEmail
	emails, err = d.FindUserActiveEmails(ctx, userId, bson.M{"email": 1}, nil)
	if err != nil {
		return nil, err
	}

	// Iterate through the emails
	for _, email := range emails {
		userActiveEmails = append(userActiveEmails, email.Email)
	}

	return userActiveEmails, nil
}

// UserEmailExists checks if the user's email exists
func (d *Database) UserEmailExists(
	ctx context.Context,
	userId primitive.ObjectID,
	email string,
) (userEmailId string, err error) {
	// Check if the user email already exists
	userEmail, err := d.FindUserEmail(
		ctx,
		bson.M{
			"user_id":    userId,
			"email":      email,
			"revoked_at": bson.M{"$exists": false},
		},
		bson.M{"_id": 1},
		nil,
	)
	if err != nil {
		return "", err
	}
	return userEmail.ID.Hex(), nil
}

// GetMyProfile gets the user's profile
func (d *Database) GetMyProfile(userId string) (
	user *commonmongodbuser.User,
	userActiveEmails *[]string,
	userPhoneNumber string,
	err error,
) {
	// Run the transaction
	var activeEmails []string
	err = commonmongodb.CreateTransaction(
		d.client, func(sc mongo.SessionContext) error {
			// Get the full user profile
			user, err = d.FindUserByUserId(
				sc, userId, bson.M{
					"username":   1,
					"first_name": 1,
					"last_name":  1,
					"birthdate":  1,
					"joined_at":  1,
				}, nil,
			)

			// Get the user's active emails
			activeEmails, err = d.GetUserActiveEmails(sc, userId)

			// Get the user's phone number
			userPhoneNumber, err = d.GetUserPhoneNumber(sc, userId)

			return nil
		},
	)

	// Check if the transaction failed
	if err != nil {
		return nil, nil, "", err
	}

	return user, &activeEmails, userPhoneNumber, nil
}
