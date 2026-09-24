package delivery

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	boilerplatedelivery "backend/delivery/boilerplate"
	boletadelivery "backend/delivery/boleta"
	productodelivery "backend/delivery/producto"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	boilerplatedelivery.NewModule(db).RegisterRoutes(api)
	productodelivery.NewModule(db).RegisterRoutes(api)
	boletadelivery.NewModule(db).RegisterRoutes(api)

	return r
}
