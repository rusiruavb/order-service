package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"order_service/configs"
	"order_service/handlers"
	"order_service/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = configs.GetCollections(configs.DB, "orders")
var productCollection *mongo.Collection = configs.GetCollections(configs.DB, "products")
var orderValidator = validator.New()

func CreateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	var order models.Order
	defer cancel()

	if err := c.BodyParser(&order); err != nil {
		return handlers.SendBadRequestResponse(c, &fiber.Map{"data": err.Error()})
	}

	if validateErr := orderValidator.Struct(&order); validateErr != nil {
		return handlers.SendBadRequestResponse(c, &fiber.Map{"data": validateErr.Error()})
	}

	authToken := c.Request().Header.Peek("Authorization")
	if len(authToken) == 0 {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Authentication token is required"})
	}

	authUser, err := authenticateUser(string(authToken))
	if err != nil {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": err.Error()})
	}

	if !authUser.IsAuthorized {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Unauthorized user"})
	}

	userData, err := getUserInfo(string(authToken))
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	userId, err := primitive.ObjectIDFromHex(userData.UserId)
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	newOrder := models.Order{
		UserId:           userId,
		OrderId:          primitive.NewObjectID(),
		OrderDate:        time.Now().UTC().String(),
		OrderDiscription: order.OrderDiscription,
		OrderFee:         order.OrderFee,
		Products:         order.Products,
	}

	_, err = orderCollection.InsertOne(ctx, newOrder)
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	// Call email service
	endPoint := configs.EnvEmailService() + "/send"
	subject := "Order #" + newOrder.OrderId.Hex()
	orderPriceFloat := math.Round(newOrder.OrderFee*100) / 100
	orderPriceStr := fmt.Sprintf("%.2f", orderPriceFloat)
	body := "<h3>Hello " + userData.FirstName + " " + userData.LastName + "</h3>" + "<p>Your order created successfully.</p>" + "<p>Order ID: #" + newOrder.OrderId.Hex() + "</p>" + "<p>Total Price - Rs." + orderPriceStr + "</p>"
	emailData := map[string]string{"to": userData.Email, "subject": subject, "body": body}
	jsonData, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("📧 Sending email...")
	resp, err := http.Post(endPoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewDecoder(resp.Body)
	fmt.Println("✅ Email sent")

	return handlers.SendSuccessResponse(c, &fiber.Map{"data": newOrder})
}

func GetOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	authToken := c.Request().Header.Peek("Authorization")
	if len(string(authToken)) == 0 {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Authentication token is required"})
	}

	authUser, err := authenticateUser(string(authToken))
	if err != nil {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": err.Error()})
	}

	if !authUser.IsAuthorized {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Unauthorized user"})
	}

	userData, err := getUserInfo(string(authToken))
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	userId, err := primitive.ObjectIDFromHex(userData.UserId)
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	cursor, err := orderCollection.Find(ctx, bson.M{"userid": userId})
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	var resutls []models.OrderResponse
	err = cursor.All(ctx, &resutls)
	if err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"data": err.Error()})
	}

	for _, ord := range resutls {
		for j, prod := range ord.Products {
			var product = getProductInfo(prod.ProductId)
			ord.Products[j].ProductTitle = product.ProductTitle
			ord.Products[j].ImageURL = product.ImageURL
			ord.Products[j].Price = product.Price
			ord.Products[j].CategoryId = product.CategoryId
		}
	}

	return handlers.SendSuccessResponse(c, &fiber.Map{"data": resutls})
}

func GetOrderById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	orderId := c.Params("orderId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(orderId)

	authToken := c.Request().Header.Peek("Authorization")
	if len(string(authToken)) == 0 {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Authentication token is required"})
	}

	authUser, err := authenticateUser(string(authToken))
	if err != nil {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": err.Error()})
	}

	if !authUser.IsAuthorized {
		return handlers.SendBadAuthResponse(c, &fiber.Map{"data": "Unauthorized user"})
	}

	var order models.OrderResponse
	if err := orderCollection.FindOne(ctx, bson.M{"orderid": objId}).Decode(&order); err != nil {
		return handlers.SendErrorResponse(c, &fiber.Map{"message": err.Error()})
	}

	for i, prod := range order.Products {
		var product = getProductInfo(prod.ProductId)
		order.Products[i].ProductTitle = product.ProductTitle
		order.Products[i].ImageURL = product.ImageURL
		order.Products[i].Price = product.Price
		order.Products[i].CategoryId = product.CategoryId
	}

	return handlers.SendSuccessResponse(c, &fiber.Map{"data": order})
}

func getProductInfo(productId primitive.ObjectID) models.OrderProduct {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var product models.OrderProduct
	err := productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
	if err != nil {
		fmt.Println(err.Error())
	}
	return product
}

func authenticateUser(authToken string) (*models.UserAuth, error) {
	client := &http.Client{}
	authServiceEndpoint := configs.EnvAuthService()
	authReq, _ := http.NewRequest("GET", authServiceEndpoint, nil)
	authReq.Header.Add("Authorization", string(authToken))
	authRes, err := client.Do(authReq)
	if err != nil {
		return nil, err
	}
	defer authRes.Body.Close()
	var authUser models.UserAuth
	json.NewDecoder(authRes.Body).Decode(&authUser)
	return &authUser, nil
}

func getUserInfo(authToken string) (*models.OrderUser, error) {
	client := &http.Client{}
	userServiceEndpoint := configs.EnvUserService()
	userReq, _ := http.NewRequest("GET", userServiceEndpoint, nil)
	userReq.Header.Add("Authorization", string(authToken))
	userRes, err := client.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer userRes.Body.Close()
	var userData models.OrderUser
	json.NewDecoder(userRes.Body).Decode(&userData)
	return &userData, nil
}
