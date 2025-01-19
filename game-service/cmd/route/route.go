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
	betOptionCollection := database.GetCollection(databaseName, "bet_options")
	bonusGameCollection := database.GetCollection(databaseName, "bonus_games")
	gameResultCollection := database.GetCollection(databaseName, "game_results")
	gameSessionCollection := database.GetCollection(databaseName, "game_sessions")
	playerCollection := database.GetCollection(databaseName, "players")
	reelCollection := database.GetCollection(databaseName, "reels")
	symbolCollection := database.GetCollection(databaseName, "symbols")

	// Create repositories and services
	betOptionRepo := outbound.NewBetOptionRepository(betOptionCollection)
	betOptionService := service.NewBetService(betOptionRepo)

	bonusGameRepo := outbound.NewBonusGameRepository(bonusGameCollection)
	bonusGameService := service.NewBonusGameService(bonusGameRepo)

	gameResultRepo := outbound.NewGameResultRepository(gameResultCollection)
	gameResultService := service.NewGameResultService(gameResultRepo)

	gameSessionRepo := outbound.NewGameSessionRepository(gameSessionCollection)
	gameSessionService := service.NewGameSessionService(gameSessionRepo)

	playerRepo := outbound.NewPlayerRepository(playerCollection)
	playerService := service.NewPlayerService(playerRepo)

	reelRepo := outbound.NewReelRepository(reelCollection)
	reelService := service.NewReelService(reelRepo)

	symbolRepo := outbound.NewSymbolRepository(symbolCollection)
	symbolService := service.NewSymbolService(symbolRepo)

	// Create handlers
	betOptionHandler := inbound.NewBetHandler(betOptionService)
	bonusGameHandler := inbound.NewBonusGameHandler(bonusGameService)
	gameResultHandler := inbound.NewGameResultHandler(gameResultService)
	gameSessionHandler := inbound.NewGameSessionHandler(gameSessionService)
	playerHandler := inbound.NewPlayerHandler(playerService)
	reelHandler := inbound.NewReelHandler(reelService)
	symbolHandler := inbound.NewSymbolHandler(symbolService)

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

	e.POST("/bet_option", betOptionHandler.CreateBet)
	e.GET("/bet_option/:id", betOptionHandler.GetBetByID)
	e.PUT("/bet_option/:id", betOptionHandler.UpdateBet)
	e.DELETE("/bet_option/:id", betOptionHandler.DeleteBet)
	
	e.POST("/bonus_game", bonusGameHandler.CreateBonusGame)
	e.GET("/bonus_game/:id", bonusGameHandler.GetBonusGameByID)
	e.PUT("/bonus_game/:id", bonusGameHandler.UpdateBonusGame)
	e.DELETE("/bonus_game/:id", bonusGameHandler.DeleteBonusGame)

	e.POST("/game_result", gameResultHandler.CreateGameResult)
	e.GET("/game_result/:id", gameResultHandler.GetGameResultByID)
	e.PUT("/game_result/:id", gameResultHandler.UpdateGameResult)
	e.DELETE("/game_result/:id", gameResultHandler.DeleteGameResult)

	e.POST("/game_session", gameSessionHandler.CreateGameSession)
	e.GET("/game_session/:id", gameSessionHandler.GetGameSessionByID)
	e.PUT("/game_session/:id", gameSessionHandler.UpdateGameSession)
	e.DELETE("/game_session/:id", gameSessionHandler.DeleteGameSession)

	e.POST("/player", playerHandler.CreatePlayer)
	e.GET("/player/:id", playerHandler.GetPlayerByID)
	e.PUT("/player/:id", playerHandler.UpdatePlayer)
	e.DELETE("/player/:id", playerHandler.DeletePlayer)

	e.POST("/reel", reelHandler.CreateReel)
	e.GET("/reel/:id", reelHandler.GetReelByID)
	e.PUT("/reel/:id", reelHandler.UpdateReel)
	e.DELETE("/reel/:id", reelHandler.DeleteReel)

	e.POST("/symbol", symbolHandler.CreateSymbol)
	e.GET("/symbol/:id", symbolHandler.GetSymbolByID)
	// e.PUT("/symbol/:id", symbolHandler.UpdateBet)
	// e.DELETE("/symbol/:id", symbolHandler.Dele)

	return e
}
