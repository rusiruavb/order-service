package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	OrderId          primitive.ObjectID   `json:"orderid,omitempty"`
	UserId           primitive.ObjectID   `json:"userid,omitempty"`
	OrderDate        string               `json:"orderdate,omitempty"`
	OrderDiscription string               `json:"orderdiscription,omitempty"`
	OrderFee         float64              `json:"orderfee,omitempty" validate:"required"`
	Products         []OrderProductRecord `json:"products,omitempty"`
}

type OrderProductRecord struct {
	ProductId primitive.ObjectID `json:"productid,omitempty"`
	Quantity  int                `json:"quantity,omitempty"`
	UnitPrice float64            `json:"unitprice,omitempty"`
}

type OrderResponse struct {
	OrderId          primitive.ObjectID `json:"orderid,omitempty"`
	UserId           primitive.ObjectID `json:"userid,omitempty"`
	OrderDate        string             `json:"orderdate,omitempty" validate:"required"`
	OrderDiscription string             `json:"orderdiscription,omitempty"`
	OrderFee         float64            `json:"orderfee,omitempty" validate:"required"`
	Products         []OrderProduct     `json:"products,omitempty"`
}

type OrderProduct struct {
	ProductId    primitive.ObjectID `json:"productid,omitempty"`
	CategoryId   string             `json:"categoryid,omitempty"`
	ProductTitle string             `json:"producttitle,omitempty"`
	ImageURL     string             `json:"imageurl,omitempty"`
	Price        float64            `json:"price,omitempty"`
}

type OrderUser struct {
	UserId    string `json:"_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

type UserAuth struct {
	IsAuthorized bool `json:"isAuthorized"`
}
