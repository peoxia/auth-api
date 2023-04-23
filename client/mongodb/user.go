package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/peoxia/auth-api/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user in the database
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	GoogleID string             `bson:"googleid"`
	FullName string             `bson:"fullname"`
	Email    string             `bson:"email"`
	Phone    string             `bson:"phone"`
}

// FindUserByID finds a user in the database by their ID
func (c *Client) FindUserByEmail(ctx context.Context, email string) (*auth.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := c.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	return &auth.Profile{
		ID:    user.GoogleID,
		Name:  user.FullName,
		Email: user.Email,
	}, nil
}

func (c *Client) UpsertUser(ctx context.Context, user auth.Profile) error {
	filter := bson.M{"googleid": user.ID}
	update := bson.M{
		"$set": bson.M{
			"fullname": user.Name,
			"email":    user.Email,
		},
	}

	// Set the upsert option to true
	opt := options.Update().SetUpsert(true)

	// Perform an upsert operation
	_, err := c.collection.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from the database
func (c *Client) DeleteUser(ctx context.Context, email string) error {
	filter := bson.M{"email": email}

	// Delete the user that matches the email filter
	res, err := c.collection.DeleteOne(context.Background(), filter, options.Delete())

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no user with email %s found", email)
	}

	return nil
}
