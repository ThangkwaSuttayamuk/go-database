package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Supplier_id int    `json:"supplier_id"`
	Time        string `json:"time"`
	Description string `json:"description"`
	Category_id int    `json:"category_id"`
}

type Category struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type Supplier struct {
	ID   int
	Name string
}

type DayOff struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	Uid  int    `json:"uid"`
}

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var db *sql.DB

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	db = sdb

	// Check the connection to make sure
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/product/:id", getProduct)
	app.Post("/product", createProductHandler)
	app.Put("/product/:id", updateProductHandle)
	app.Delete("/product/:id", deleteProductHandler)
	app.Get("/product", getProductsHandler)
	app.Post("/user", createUserHandler)
	app.Post("/dayoff/:uid", addDayOffHandler)
	app.Delete("/dayoff/:uid", deleteDayOffHandle)
	app.Get("/dayoff/:uid", getDayOffsHandler)
	app.Get("/users", getUsersHandler)
	app.Post("/category", createCategoryHandler)
	app.Get("/category", getCategoriesHandler)

	// Start Fiber and Socket.IO
	app.Listen(":8080")

	// app.Listen(":8080")

	// fmt.Println("Successfully connected!")

	// err = createProduct(&Product{Name: "Go product", Price: 220})

	// product, err := getProductById(2)

	// for choice 1
	// err = updateProduct(1, &Product{Name: "New name", Price: 310})

	//for choice 2
	// product, err := updateProduct(5, &Product{Name: "New name id 3", Price: 320})

	// err = deleteProduct(1)

	// product, err := getProducts()

	// err = addProductAndSupplier("CPF", "Potato Chips", 25)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Create Successful !")

	// fmt.Println("Create Successful !", product)
}

func getUsersHandler(c *fiber.Ctx) error {
	users, err := getUsers()

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(users)
}

func getCategoriesHandler(c *fiber.Ctx) error {
	categories, err := getCategories()

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(categories)
}

func getProduct(c *fiber.Ctx) error {
	// Convert the "id" parameter from the request to an integer
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	// Call the getProductById function to retrieve the product
	product, err := getProductById(id)

	if err != nil {
		// Check if the error is due to no rows being found
		if err.Error() == fmt.Sprintf("no product found with id %d", id) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		// Handle other possible errors
		return c.Status(fiber.StatusInternalServerError).SendString("An error occurred while retrieving the product")
	}

	// If the product is found, return it as a JSON response
	return c.JSON(product)
}

func createProductHandler(c *fiber.Ctx) error {
	product := new(Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := createProduct(product)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendString("Create Product Successfully.")
}

func createCategoryHandler(c *fiber.Ctx) error {
	category := new(Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := createCategory(category)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendString("Create Category Successfully.")
}

func createUserHandler(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := createUser(user)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendString("Create User Successfully.")
}

func addDayOffHandler(c *fiber.Ctx) error {
	dayOff := new(DayOff)

	if err := c.BodyParser(dayOff); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	uid, err := strconv.Atoi(c.Params("uid"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = addDayOff(uid, dayOff)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendString("Add Dayoff Successfully.")
}

func deleteDayOffHandle(c *fiber.Ctx) error {
	dayOff := new(DayOff)

	// Parse JSON body into DayOff struct
	if err := c.BodyParser(dayOff); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Parse the uid parameter
	uid, err := strconv.Atoi(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	// Attempt to delete the day off
	err = deleteDayOff(uid, dayOff)
	if err != nil {
		if err.Error() == "no record found to delete" {
			return c.Status(fiber.StatusNotFound).SendString("No matching record found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete day off")
	}

	return c.SendString("Day off deleted successfully.")
}

func updateProductHandle(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	product := new(Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	updateProduct, err := updateProduct(id, product)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(updateProduct)
}

func deleteProductHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = deleteProduct(id)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func getProductsHandler(c *fiber.Ctx) error {
	products, err := getProducts()

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(products)
}

func getDayOffsHandler(c *fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	dayOffs, err := getDayOffs(uid)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(dayOffs)
}
