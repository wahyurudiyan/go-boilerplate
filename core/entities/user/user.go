package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        int64      `db:"id"`
	Role      string     `db:"role"`
	Email     string     `db:"email"`
	UniqueId  string     `db:"unique_id"`
	Fullname  string     `db:"fullname"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toMongoDocument converts User to UserMongoDocument
func (user User) ToMongoDocument() UserMongoDocument {
	doc := UserMongoDocument{
		Role:      user.Role,
		Email:     user.Email,
		UniqueId:  user.UniqueId,
		Fullname:  user.Fullname,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	// If Id is set, try to convert it to ObjectID
	if user.Id > 0 {
		// In MongoDB migration scenarios, you might store the original SQL ID
		// or convert it to ObjectID as needed
	}

	return doc
}

type UserMongoDocument struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Role      string             `bson:"role"`
	Email     string             `bson:"email"`
	UniqueId  string             `bson:"unique_id"`
	Fullname  string             `bson:"fullname"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

// toUserEntity converts UserMongoDocument to User
func (doc UserMongoDocument) ToUserEntity() User {
	return User{
		// Convert the ObjectID to int64 if needed for compatibility
		// In a real scenario, you might want to handle this ID conversion differently
		Id:        int64(doc.Id.Timestamp().Unix()), // Simple conversion example
		Role:      doc.Role,
		Email:     doc.Email,
		UniqueId:  doc.UniqueId,
		Fullname:  doc.Fullname,
		Username:  doc.Username,
		Password:  doc.Password,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		DeletedAt: doc.DeletedAt,
	}
}

// toUserEntities converts a slice of UserMongoDocument to a slice of User
func MongoDocsToUserEntities(docs []UserMongoDocument) []User {
	users := make([]User, len(docs))
	for i, doc := range docs {
		users[i] = doc.ToUserEntity()
	}
	return users
}
