package post

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" swaggertype:"primitive,string"`
	Content string             `bson:"content" json:"content"`
	Author  string             `bson:"author" json:"author"`
}
