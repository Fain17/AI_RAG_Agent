package routes

import (
	_ "github.com/fain17/rag-backend/docs"

	handlers "github.com/fain17/rag-backend/api/handlers"
	"github.com/fain17/rag-backend/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(queries *db.Queries) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	//Swagger Routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CRUD + search routes
	r.POST("/files", handlers.UploadHandler(queries))
	r.GET("/files/:id", handlers.GetHandler(queries))
	r.PUT("/files/:id", handlers.UpdateHandler(queries))
	r.DELETE("/files/:id", handlers.DeleteHandler(queries))
	r.POST("/search", handlers.SearchHandler(queries))

	return r
}
