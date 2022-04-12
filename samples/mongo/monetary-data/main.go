package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const tax string = "0.1"

type Product struct {
	ID        primitive.ObjectID   `bson:"_id"`
	Name      string               `json:"name,omitempty" bson:"name"`
	Price     primitive.Decimal128 `json:"price,omitempty" bson:"price"`
	Tax       primitive.Decimal128 `json:"tax,omitempty" bson:"tax"`
	Total     primitive.Decimal128 `json:"total,omitempty" bson:"total"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

var collection *mongo.Collection

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")
	collection = client.Database("testing").Collection("products")
}

func main() {
	// Creating a product
	item := Product{
		ID:        primitive.NewObjectID(),
		Name:      "Orange Soda",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setting monetary data to decimal
	price, err := decimal.NewFromString("0.2")
	if err != nil {
		panic(err)
	}

	taxDecimal, err := decimal.NewFromString(tax)
	if err != nil {
		panic(err)
	}

	fmt.Println("Price in Decimal", price, "Tax in Decimal:", taxDecimal)

	// Calculate taxes of the product
	productTaxDecimal := taxDecimal.Mul(price)
	fmt.Println("Product Taxes in Decimal:", productTaxDecimal, "Expected Result:", 0.02)

	// Calculate total (price + productTax)
	totalDecimal := productTaxDecimal.Add(price)
	fmt.Println("Product Total:", totalDecimal, "Expected Result:", 0.22)

	// Convert decimal.Decimal to primitive.Decimal128
	mongoPriceDecimal, err := primitive.ParseDecimal128(price.String())
	if err != nil {
		panic(err)
	}

	mongoProductTaxesDecimal, err := primitive.ParseDecimal128(productTaxDecimal.String())
	if err != nil {
		panic(err)
	}

	mongoTotalDecimal, err := primitive.ParseDecimal128(totalDecimal.String())
	if err != nil {
		panic(err)
	}

	item.Price = mongoPriceDecimal
	item.Tax = mongoProductTaxesDecimal
	item.Total = mongoTotalDecimal

	_, err = collection.InsertOne(context.TODO(), item)
	if err != nil {
		panic(err)
	}

	cur, err := collection.Find(context.TODO(), bson.M{"total": bson.M{"$type": "decimal"}})
	if err != nil {
		panic(err)
	}

	var results []Product
	err = cur.All(context.TODO(), &results)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", results)

}
