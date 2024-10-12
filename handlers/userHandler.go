package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "calorie-tracker/models" // Adjust this import path as necessary
)

var userCollection *mongo.Collection

// SetUpUserCollection initializes the user collection
func SetUpUserCollection(client *mongo.Client) {
    userCollection = client.Database("your_database_name").Collection("users")
}

// Register creates a new user account
func Register(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Check if the email already exists
    existingUser := &models.User{}
    err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(existingUser)
    if err == nil {
        return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
    }
    user.Password = string(hashedPassword)

    // Insert the user into the database
    _, err = userCollection.InsertOne(context.TODO(), user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
    }

    return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login authenticates a user
func Login(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Find the user in the database
    var dbUser models.User
    err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&dbUser)
    if err != nil {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }

    // Compare hashed passwords
    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }

    // Create JWT token
    token, err := createJWT(dbUser.ID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})
    }

    return c.JSON(fiber.Map{"token": token})
}

// createJWT generates a JWT token for the user
func createJWT(userID primitive.ObjectID) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := []byte("your_secret_key") // Replace with your actual secret key

    return token.SignedString(secret)
}
