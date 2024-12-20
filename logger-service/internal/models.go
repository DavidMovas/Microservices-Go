package internal

import "time"

type Log struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
