package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"objectswaterfall.com/application/handlers"
	"objectswaterfall.com/data"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}
	err = data.InitDbConnection()
	if err != nil {
		panic(err)
	}
	engine := gin.Default()
	// mux := http.NewServeMux()
	// logsHub := &hubs.LogsHub{}
	// signalrServer, err := signalr.NewServer(context.Background(), signalr.SimpleHubFactory(logsHub))
	// if err != nil {
	// 	panic(err)
	// }
	// signalrServer.MapHTTP(signalr.WithHTTPServeMux(mux), "/logsConnection")

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.POST("/add", handlers.Add)
	engine.POST("/start", handlers.Start)
	engine.GET("/stop", handlers.Stop)
	engine.POST("/seed", handlers.Seed)
	engine.GET("/getWorkers", handlers.GetWorkers)
	engine.GET("/getRunningWorkers", handlers.GetRunningWorkers)
	engine.GET("/logsWs", handlers.WebSocketHandler)
	//engine.Any("/logsConnection", gin.WrapH(mux))

	engine.Run(":8888")
}
