package server

import (
	"avito-tech-internship/service"
	"avito-tech-internship/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	router *gin.Engine
	db     *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	s := &Server{
		router: gin.Default(),
		db:     db,
	}
	s.setupRouter()
	return s
}

func (s *Server) setupRouter() {
	repository := storage.NewPostgresRepository(s.db)
	appService := service.NewService(repository)
	httpHandler := NewHandler(appService)

	teams := s.router.Group("/team")
	{
		teams.POST("/add", httpHandler.createNewTeam)
		teams.GET("/get/:teamName", httpHandler.GetTeamByName)
	}

	users := s.router.Group("/users")
	{
		users.POST("/addNew", httpHandler.AddNewUser)
		users.GET("/users/:id", httpHandler.GetUser)
		users.POST("/setIsActive", httpHandler.SetUserActive)
	}

	s.router.GET("/health", func(c *gin.Context) {
		if err := s.db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})
}

func (s *Server) Start() {
	s.router.Run(":8080")
}
