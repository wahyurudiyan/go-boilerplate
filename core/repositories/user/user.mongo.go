package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	userEnt "github.com/wahyurudiyan/go-boilerplate/core/entities/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure userMongoRepositoryImpl implements IUserRepository interface
var _ IUserRepository = (*userMongoRepositoryImpl)(nil)

// userMongoRepositoryImpl implements the IUserRepository interface for MongoDB
type userMongoRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserMongoRepository creates a new instance of IUserRepository for MongoDB
func NewUserMongoRepository(db *mongo.Database, collectionName string) IUserRepository {
	collection := db.Collection(collectionName)

	// Create indexes for efficient queries
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"deleted_at": nil}),
		},
		{
			Keys:    bson.D{{Key: "unique_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"deleted_at": nil}),
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}

	// Create indexes in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection.Indexes().CreateMany(ctx, indexModels)
	}()

	return &userMongoRepositoryImpl{
		collection: collection,
	}
}

// SaveUser inserts a single user into MongoDB
func (r *userMongoRepositoryImpl) SaveUser(ctx context.Context, user userEnt.User) error {
	doc := user.ToMongoDocument()

	// Ensure times are set
	now := time.Now()
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = now
	}
	if doc.UpdatedAt.IsZero() {
		doc.UpdatedAt = now
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with this email or unique_id already exists: %w", err)
		}
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// SaveUsers inserts multiple users into MongoDB
func (r *userMongoRepositoryImpl) SaveUsers(ctx context.Context, users []userEnt.User) error {
	if len(users) == 0 {
		return nil
	}

	docs := make([]interface{}, len(users))
	now := time.Now()

	for i, user := range users {
		doc := user.ToMongoDocument()

		// Ensure times are set
		if doc.CreatedAt.IsZero() {
			doc.CreatedAt = now
		}
		if doc.UpdatedAt.IsZero() {
			doc.UpdatedAt = now
		}

		docs[i] = doc
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("one or more users with duplicate email or unique_id: %w", err)
		}
		return fmt.Errorf("failed to save users: %w", err)
	}
	return nil
}

// UpdateUser updates an existing user in MongoDB
func (r *userMongoRepositoryImpl) UpdateUser(ctx context.Context, user userEnt.User) error {
	// Set update time
	user.UpdatedAt = time.Now()

	// Convert to mongo document
	doc := user.ToMongoDocument()

	// In MongoDB, we'll use the UniqueId for updates if Id is not provided
	var filter bson.M
	if user.Id > 0 {
		filter = bson.M{"_id": doc.Id, "deleted_at": nil}
	} else if user.UniqueId != "" {
		filter = bson.M{"unique_id": user.UniqueId, "deleted_at": nil}
	} else {
		return errors.New("user id or unique_id is required for update")
	}

	update := bson.M{
		"$set": bson.M{
			"role":       user.Role,
			"email":      user.Email,
			"unique_id":  user.UniqueId,
			"fullname":   user.Fullname,
			"username":   user.Username,
			"Password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with this email or unique_id already exists: %w", err)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUserById performs a soft delete by ID
func (r *userMongoRepositoryImpl) DeleteUserById(ctx context.Context, id int64) error {
	// Since MongoDB uses ObjectID and our interface uses int64,
	// we need a strategy to handle this. Here we assume UniqueId is set based on ID
	// A better approach would depend on your ID mapping strategy
	filter := bson.M{"_id": primitive.ObjectID{}, "deleted_at": nil} // Need to convert id to ObjectID

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// DeleteUserByEmail performs a soft delete by email
func (r *userMongoRepositoryImpl) DeleteUserByEmail(ctx context.Context, email string) error {
	filter := bson.M{"email": email, "deleted_at": nil}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete user by email: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with email %s not found", email)
	}

	return nil
}

// DeleteUserByUniqueId performs a soft delete by unique ID
func (r *userMongoRepositoryImpl) DeleteUserByUniqueId(ctx context.Context, uniqueId string) error {
	filter := bson.M{"unique_id": uniqueId, "deleted_at": nil}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete user by unique id: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with unique id %s not found", uniqueId)
	}

	return nil
}

// RetrieveAllUser retrieves all users with pagination
func (r *userMongoRepositoryImpl) RetrieveAllUser(ctx context.Context, offset, limit int) ([]userEnt.User, error) {
	filter := bson.M{"deleted_at": nil}

	// Set options for pagination
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "_id", Value: 1}}) // Sort by _id

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all users: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []userEnt.UserMongoDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return userEnt.MongoDocsToUserEntities(docs), nil
}

// RetrieveUserById retrieves a user by ID
func (r *userMongoRepositoryImpl) RetrieveUserById(ctx context.Context, id int64) (userEnt.User, error) {
	// Similar to DeleteUserById, we need a strategy to map int64 to ObjectID
	filter := bson.M{"_id": primitive.ObjectID{}, "deleted_at": nil} // Need to convert id to ObjectID

	var doc userEnt.UserMongoDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userEnt.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by id: %w", err)
	}

	return doc.ToUserEntity(), nil
}

// RetrieveUserByIds retrieves users by IDs
func (r *userMongoRepositoryImpl) RetrieveUserByIds(ctx context.Context, ids []int64) ([]userEnt.User, error) {
	if len(ids) == 0 {
		return []userEnt.User{}, nil
	}

	// Convert int64 IDs to ObjectIDs (or another strategy)
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	// Conversion would go here

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by ids: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []userEnt.UserMongoDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return userEnt.MongoDocsToUserEntities(docs), nil
}

// RetrieveUserByEmail retrieves a user by email
func (r *userMongoRepositoryImpl) RetrieveUserByEmail(ctx context.Context, email string) (userEnt.User, error) {
	filter := bson.M{"email": email, "deleted_at": nil}

	var doc userEnt.UserMongoDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userEnt.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	return doc.ToUserEntity(), nil
}

// RetrieveUserByEmails retrieves users by emails
func (r *userMongoRepositoryImpl) RetrieveUserByEmails(ctx context.Context, emails []string) ([]userEnt.User, error) {
	if len(emails) == 0 {
		return []userEnt.User{}, nil
	}

	filter := bson.M{
		"email":      bson.M{"$in": emails},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by emails: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []userEnt.UserMongoDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return userEnt.MongoDocsToUserEntities(docs), nil
}

// RetrieveUserByUniqueId retrieves a user by unique ID
func (r *userMongoRepositoryImpl) RetrieveUserByUniqueId(ctx context.Context, uniqueId string) (userEnt.User, error) {
	filter := bson.M{"unique_id": uniqueId, "deleted_at": nil}

	var doc userEnt.UserMongoDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userEnt.User{}, fmt.Errorf("user with unique id %s not found", uniqueId)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by unique id: %w", err)
	}

	return doc.ToUserEntity(), nil
}

// RetrieveUserByUniqueIds retrieves users by unique IDs
func (r *userMongoRepositoryImpl) RetrieveUserByUniqueIds(ctx context.Context, uniqueIds []string) ([]userEnt.User, error) {
	if len(uniqueIds) == 0 {
		return []userEnt.User{}, nil
	}

	filter := bson.M{
		"unique_id":  bson.M{"$in": uniqueIds},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by unique ids: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []userEnt.UserMongoDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return userEnt.MongoDocsToUserEntities(docs), nil
}
