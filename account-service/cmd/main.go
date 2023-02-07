package main

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/server"
	"acbs.com.vn/account-service/pkg/logger"
	accountService "acbs.com.vn/account-service/proto"
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("Account-Service")

	conn, err := grpc.Dial("localhost:5002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := accountService.NewAccountServiceClient(conn)

	// g := gin.Default()

	app := fiber.New()

	app.Post("/add", func(c *fiber.Ctx) error {
		var requestBody *accountService.CreateAccountReq
		err := c.BodyParser(&requestBody)

		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		res, err := client.CreateAccount(context.Background(), requestBody)
		if err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"result": fmt.Sprint(res.Name),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	})

	app.Get("/account/:a", func(c *fiber.Ctx) error {
		req := &accountService.GetAccountByNameReq{Name: c.Params("a")}

		if res, err := client.GetAccountByName(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"result": fmt.Sprint(res.Accounts),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	})

	go func() {
		log.Fatal(app.Listen(":3000"))
	}()

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
