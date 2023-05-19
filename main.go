package main

import (
	"flag"

	"github.com/agonist/hotel-reservation/api"
	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dburi = "host=localhost user=admin password=supersecret dbname=ticketing port=5432 sslmode=disable"

var config = fiber.Config{
	JSONEncoder: sonic.Marshal,
	JSONDecoder: sonic.Unmarshal,
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client := setupDb()

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewPgUserStore(client))
	orgHandler := api.NewOrgHandler(db.NewPgOrgStore(client))
	eventHandler := api.NewEventHandler(db.NewPgEventStore(client))

	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleListUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	apiV1.Post("/organization", orgHandler.HandlePost)
	apiV1.Get("/organization", orgHandler.HandleList)
	apiV1.Get("/organization/:id", orgHandler.HandleGet)
	apiV1.Delete("/organization/:id", orgHandler.HandleDelete)
	apiV1.Put("/organization/:id", orgHandler.HandlePut)

	apiV1.Post("/event", eventHandler.HandlePost)
	apiV1.Get("/event", eventHandler.HandleList)
	apiV1.Get("/event/:id", eventHandler.HandleGet)
	apiV1.Delete("/event/:id", eventHandler.HandleDelete)
	apiV1.Put("/event/:id", eventHandler.HandlePut)

	app.Listen(*listenAddr)
}

func setupDb() *gorm.DB {
	// dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dburi), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the db")
	}
	err = db.AutoMigrate(&types.User{})
	if err != nil {
		panic("failed ti run migrations")
	}

	return db
}
