package route

import (
	"game-service/internal/adapters/inbound"
	"game-service/internal/adapters/outbound"
	"game-service/internal/app/service"
	"game-service/pkg/database"
	"game-service/pkg/utils"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	// Initialize MongoDB connection
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	gameResultCollection := database.GetCollection(databaseName, "game_results")
	gameSessionCollection := database.GetCollection(databaseName, "game_sessions")
	playerCollection := database.GetCollection(databaseName, "players")
	symbolsCollection := database.GetCollection(databaseName, "symbols")
	configsCollection := database.GetCollection(databaseName, "configs")
	featuresCollection := database.GetCollection(databaseName, "features")
	paylinesCollection := database.GetCollection(databaseName, "paylines")
	reelsCollection := database.GetCollection(databaseName, "reels")

	// Create repositories and services

	gameResultRepo := outbound.NewGameResultRepository(gameResultCollection)
	gameResultService := service.NewGameResultService(gameResultRepo)

	gameSessionRepo := outbound.NewGameSessionRepository(gameSessionCollection)
	gameSessionService := service.NewGameSessionService(gameSessionRepo)

	playerRepo := outbound.NewPlayerRepository(playerCollection)
	playerService := service.NewPlayerService(playerRepo)

	symbolsRepo := outbound.NewSymbolsRepository(symbolsCollection)
	symbolsService := service.NewSymbolsService(symbolsRepo)

	configsRepo := outbound.NewConfigsRepository(configsCollection)
	configsService := service.NewConfigsService(configsRepo)

	featuresRepo := outbound.NewFeaturesRepository(featuresCollection)
	featuresService := service.NewFeaturesService(featuresRepo)

	paylinesRepo := outbound.NewPaylinesRepository(paylinesCollection)
	paylinesService := service.NewPaylinesService(paylinesRepo)

	reelsRepo := outbound.NewReelsRepository(reelsCollection)
	reelsService := service.NewReelsService(reelsRepo)

	// Backup
	mongoDB := symbolsCollection.Database() // Lấy đối tượng *mongo.Database
	backupRepo := outbound.NewBackupRepository()
	backupService := service.NewBackupService(backupRepo, mongoDB)
	backupHandler := inbound.NewBackupHandler(backupService)

	// Create handlers
	gameResultHandler := inbound.NewGameResultHandler(gameResultService)
	gameSessionHandler := inbound.NewGameSessionHandler(gameSessionService)
	playerHandler := inbound.NewPlayerHandler(playerService)
	symbolHandler := inbound.NewSymbolsHandler(symbolsService)
	paylineHandler := inbound.NewPaylinesHandler(paylinesService)
	reelHandler := inbound.NewReelsHandler(reelsService)
	configsHandler := inbound.NewConfigsHandler(configsService)
	featureHandler := inbound.NewFeaturesHandler(featuresService)

	// Set up Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Use(utils.CorsMiddleware())

	e.POST("/game_result", gameResultHandler.CreateGameResult)
	e.GET("/game_result/:id", gameResultHandler.GetGameResultByID)
	e.PUT("/game_result/:id", gameResultHandler.UpdateGameResult)
	e.DELETE("/game_result/:id", gameResultHandler.DeleteGameResult)
	e.GET("/game_results/:session_id", gameResultHandler.GetGameResultsBySessionID)

	e.POST("/game_session", gameSessionHandler.CreateGameSession)
	e.GET("/game_session/:id", gameSessionHandler.GetGameSessionByID)
	e.PUT("/game_session/:id", gameSessionHandler.UpdateGameSession)
	e.DELETE("/game_session/:id", gameSessionHandler.DeleteGameSession)
	e.GET("/game_sessions/:player_id", gameSessionHandler.GetGameSessionsByPlayerID)

	e.POST("/player", playerHandler.CreatePlayer)
	e.GET("/player/:id", playerHandler.GetPlayerByID)
	e.PUT("/player/:id", playerHandler.UpdatePlayer)
	e.DELETE("/player/:id", playerHandler.DeletePlayer)
	e.GET("/players", playerHandler.GetAllPlayers)
	e.PUT("/player/:id/balance", playerHandler.UpdateBalance)

	e.GET("/symbols/:game_name", symbolHandler.GetSymbols)
	e.GET("/paylines/:game_name", paylineHandler.GetPaylines)
	// e.POST("/calculate_winnings", paylineHandler.CalculateWinnings)
	e.GET("/reels/:game_name", reelHandler.GetReels)
	e.GET("/configs/:game_name", configsHandler.GetConfig)
	e.GET("/features/:game_name", featureHandler.GetFeature)

	// Backup API
	e.POST("/backup", backupHandler.BackupAll)

	return e
}
