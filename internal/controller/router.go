package controller

import (
	"marketplace/internal/models/domain"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (ctrl *Controller) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", ctrl.ping)
	r.GET("/health", ctrl.healthCheck)
	authG := r.Group("/auth")
	{
		authG.POST("/sign-up", ctrl.SignUp)
		authG.POST("/sign-in", ctrl.SignIn)
		authG.GET("/refresh", ctrl.RefreshTokenPair)
	}
	apiV1G := r.Group("/api/v1", ctrl.checkUserAuthentication)
	adminG := apiV1G.Group("/admin", ctrl.checkRole(domain.AdminRole))
	{
		adminG.PUT("/users/:id/role", ctrl.SetUserRoleHandler)
	}
	shopkeeperG := apiV1G.Group("", ctrl.checkRole(domain.AdminRole, domain.ShopkeperRole))
	{
		shopkeeperG.POST("/products", ctrl.CreateProductHandler)
		shopkeeperG.PUT("/products/:id", ctrl.UpdateProductHandler)
		shopkeeperG.DELETE("/products/:id", ctrl.DeleteProductHandler)
		shopkeeperG.PUT("/shops/:id", ctrl.UpdateShopHandler)
		shopkeeperG.DELETE("/shops/:id", ctrl.DeleteShopHandler)
	}
	{
		apiV1G.GET("/products/:id", ctrl.GetProductByIDHandler)
		apiV1G.GET("/products", ctrl.ListProductsHandler)
		apiV1G.POST("/shops", ctrl.CreateShopHandler)
		apiV1G.GET("/shops/:id", ctrl.GetShopByIDHandler)
		apiV1G.GET("/shops", ctrl.ListShopsHandler)
		apiV1G.POST("/orders", ctrl.CreateOrderHandler)
		apiV1G.GET("/orders/:id", ctrl.GetOrderHandler)
	}
	return r
}
