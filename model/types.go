package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User is the model for the user entity
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Address   string             `json:"address" bson:"address"`
	BirthDate time.Time          `json:"birth_date" bson:"birth_date"`
	Deleted   bool               `json:"deleted" bson:"deleted"`
}

// UserEmail is the model for the user email entity
type UserEmail struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email      string             `json:"email" bson:"email"`
	AssignedAt time.Time          `json:"assigned" bson:"assigned"`
	VerifiedAt time.Time          `json:"verified" bson:"verified"`
	RevokedAt  time.Time          `json:"revoked" bson:"revoked"`
	IsActive   bool               `json:"is_active" bson:"is_active"`
}

// UserPhoneNumber is the model for the user phone number entity
type UserPhoneNumber struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	AssignedAt  time.Time          `json:"assigned" bson:"assigned"`
	VerifiedAt  time.Time          `json:"verified" bson:"verified"`
	RevokedAt   time.Time          `json:"revoked" bson:"revoked"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
}

// UserEmailVerification is the model for the user email verification entity
type UserEmailVerification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserEmailID primitive.ObjectID `json:"user_email_id" bson:"user_email_id"`
	UUID        string             `json:"uuid" bson:"uuid"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt  time.Time          `json:"verified_at" bson:"verified_at"`
	RevokedAt   time.Time          `json:"revoked_at" bson:"revoked_at"`
}

// UserPhoneNumberVerification is the model for the user phone number verification entity
type UserPhoneNumberVerification struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UserPhoneNumberID primitive.ObjectID `json:"user_phone_number_id" bson:"user_phone_number_id"`
	VerificationCode  string             `json:"verification_code" bson:"verification_code"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt         time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt        time.Time          `json:"verified_at" bson:"verified_at"`
	RevokedAt         time.Time          `json:"revoked_at" bson:"revoked_at"`
}

// UserPasswordReset is the model for the user password reset entity
type UserPasswordReset struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	UUID       string             `json:"uuid" bson:"uuid"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt  time.Time          `json:"expires_at" bson:"expires_at"`
	VerifiedAt time.Time          `json:"verified_at" bson:"verified_at"`
	RevokedAt  time.Time          `json:"revoked_at" bson:"revoked_at"`
}
