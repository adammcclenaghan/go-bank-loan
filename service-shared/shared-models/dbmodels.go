package shared_models

import "go.mongodb.org/mongo-driver/bson/primitive"

//ApplicationEntry represents an entry in the database for a loan application
type ApplicationEntry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"application_id"`
	Status    Status             `bson:"status,omitempty" json:"status"`
	FirstName string             `bson:"firstname,omitempty" json:"first_name"`
	LastName  string             `bson:"lastname,omitempty" json:"last_name"`
}
