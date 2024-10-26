package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User is the MongoDB for the user entity
type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Username       string             `json:"username" bson:"username"`
	FirstName      string             `json:"first_name" bson:"first_name"`
	LastName       string             `json:"last_name" bson:"last_name"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	Address        string             `json:"address,omitempty" bson:"address,omitempty"`
	BirthDate      time.Time          `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	Deleted        bool               `json:"deleted" bson:"deleted"`
}

// UserEmail is the MongoDB for the user email entity
type UserEmail struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email      string             `json:"email" bson:"email"`
	AssignedAt time.Time          `json:"assigned_at" bson:"assigned_at"`
	VerifiedAt time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	RevokedAt  time.Time          `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	IsActive   bool               `json:"is_active" bson:"is_active"`
}

// UserPhoneNumber is the MongoDB for the user phone number entity
type UserPhoneNumber struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	AssignedAt  time.Time          `json:"assigned_at" bson:"assigned_at"`
	VerifiedAt  time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	RevokedAt   time.Time          `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
}

// UserEmailVerification is the MongoDB for the user email verification entity
type UserEmailVerification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserEmailID primitive.ObjectID `json:"user_email_id" bson:"user_email_id"`
	UUID        string             `json:"uuid" bson:"uuid"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt  time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	RevokedAt   time.Time          `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
}

// UserPhoneNumberVerification is the MongoDB for the user phone number verification entity
type UserPhoneNumberVerification struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UserPhoneNumberID primitive.ObjectID `json:"user_phone_number_id" bson:"user_phone_number_id"`
	VerificationCode  string             `json:"verification_code" bson:"verification_code"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt         time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt        time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	RevokedAt         time.Time          `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
}

// UserResetPassword is the MongoDB for the user password reset entity
type UserResetPassword struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	UUID       string             `json:"uuid" bson:"uuid"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt  time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	RevokedAt  time.Time          `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
}

// UserUsernameLog is the MongoDB for the user username log entity
type UserUsernameLog struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username   string             `json:"username" bson:"username"`
	AssignedAt time.Time          `json:"assigned_at" bson:"assigned_at"`
}

// UserHashedPasswordLog is the MongoDB for the user hashed password log entity
type UserHashedPasswordLog struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	AssignedAt     time.Time          `json:"assigned_at" bson:"assigned_at"`
}
