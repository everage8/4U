package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Seq      uint               `bson:"seq" json:"id"`
	Login    string             `bson:"login" json:"login"`
	Password string             `bson:"password" json:"-"`
	Role     string             `bson:"role" json:"role"`
}

type AdminPublic struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

func (a *Admin) ToPublic() AdminPublic {
	return AdminPublic{
		ID:    a.Seq,
		Login: a.Login,
		Role:  a.Role,
	}
}
