package route

import (
	"os"
	"game-service/internal/adapters/inbound"
	"game-service/internal/adapters/outbound"
	"game-service/internal/app/service"
	"game-service/pkg/database"
	"game-service/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	// Initialize MongoDB connection
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	bonusGameCollection := database.GetCollection(databaseName, "bonus_games")
	gameResultCollection := database.GetCollection(databaseName, "game_results")
	gameSessionCollection := database.GetCollection(databaseName, "game_sessions")
	playerCollection := database.GetCollection(databaseName, "players")

	// Create repositories and services
	bonusGameRepo := outbound.NewBonusGameRepository(bonusGameCollection)
	bonusGameService := service.NewBonusGameService(bonusGameRepo)

	gameResultRepo := outbound.NewGameResultRepository(gameResultCollection)
	gameResultService := service.NewGameResultService(gameResultRepo)

	gameSessionRepo := outbound.NewGameSessionRepository(gameSessionCollection)
	gameSessionService := service.NewGameSessionService(gameSessionRepo)

	playerRepo := outbound.NewPlayerRepository(playerCollection)
	playerService := service.NewPlayerService(playerRepo)


	symbolService := service.NewSymbolService()
	paylineService := service.NewPaylineService()
	betService := service.NewBetService()
	reelService := service.NewReelService()


	// Create handlers
	bonusGameHandler := inbound.NewBonusGameHandler(bonusGameService)
	gameResultHandler := inbound.NewGameResultHandler(gameResultService)
	gameSessionHandler := inbound.NewGameSessionHandler(gameSessionService)
	playerHandler := inbound.NewPlayerHandler(playerService)
	symbolHandler := inbound.NewSymbolHandler(symbolService)
	paylineHandler := inbound.NewPaylineHandler(paylineService)
	betHandler := inbound.NewBetHandler(betService)
	reelHandler := inbound.NewReelHandler(reelService)

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

	e.POST("/bonus_game", bonusGameHandler.CreateBonusGame)
	e.GET("/bonus_game/:id", bonusGameHandler.GetBonusGameByID)
	e.PUT("/bonus_game/:id", bonusGameHandler.UpdateBonusGame)
	e.DELETE("/bonus_game/:id", bonusGameHandler.DeleteBonusGame)

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
	
	e.GET("/symbols", symbolHandler.GetSymbols)
	e.GET("/paylines", paylineHandler.GetPaylines)
	e.POST("/calculate_winnings", paylineHandler.CalculateWinnings)
	e.GET("/bets", betHandler.GetBets)
	e.GET("/reels", reelHandler.GetReel)

	return e
}
