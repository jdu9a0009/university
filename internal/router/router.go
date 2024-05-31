package router

import (
	"project/foundation/web"
	"project/internal/auth"

	"github.com/redis/go-redis/v9"

	"project/internal/middleware"
	"project/internal/pkg/repository/postgresql"

	"project/internal/repository/postgres/department"
	"project/internal/repository/postgres/district"
	"project/internal/repository/postgres/position"
	"project/internal/repository/postgres/region"
	"project/internal/repository/postgres/republic"
	"project/internal/repository/postgres/user"

	auth_controller "project/internal/controller/http/v1/auth"
	department_controller "project/internal/controller/http/v1/department"
	district_controller "project/internal/controller/http/v1/district"
	position_controller "project/internal/controller/http/v1/position"
	region_controller "project/internal/controller/http/v1/region"
	republic_controller "project/internal/controller/http/v1/republic"
	user_controller "project/internal/controller/http/v1/user"
)

type Router struct {
	*web.App
	postgresDB         *postgresql.Database
	redisDB            *redis.Client
	port               string
	auth               *auth.Auth
	fileServerBasePath string
}

func NewRouter(
	app *web.App,
	postgresDB *postgresql.Database,
	redisDB *redis.Client,
	port string,
	auth *auth.Auth,
	fileServerBasePath string,
) *Router {
	return &Router{
		app,
		postgresDB,
		redisDB,
		port,
		auth,
		fileServerBasePath,
	}
}

func (r Router) Init() error {

	// repositories:
	// - postgresql
	userPostgres := user.NewRepository(r.postgresDB)
	republicPostgres := republic.NewRepository(r.postgresDB)
	departmentProgres := department.NewRepository(r.postgresDB)
	positionProgres := position.NewRepository(r.postgresDB)
	regionProgres := region.NewRepository(r.postgresDB)
	districtProgres := district.NewRepository(r.postgresDB)

	// controller
	userController := user_controller.NewController(userPostgres)
	republicController := republic_controller.NewController(republicPostgres)
	authController := auth_controller.NewController(userPostgres)
	departmentController := department_controller.NewController(departmentProgres)
	positionController := position_controller.NewController(positionProgres)
	regionController := region_controller.NewController(regionProgres)
	districtController := district_controller.NewController(districtProgres)

	// #auth
	r.Post("/api/v1/sign-in", authController.SignIn)

	// #user
	r.Get("/api/v1/user/list", userController.GetList, middleware.Authenticate(r.auth, auth.RoleAdmin))
	r.Get("/api/v1/user/:id", userController.GetDetailById, middleware.Authenticate(r.auth, auth.RoleAdmin))
	r.Post("/api/v1/user/create", userController.Create, middleware.Authenticate(r.auth, auth.RoleAdmin))
	r.Put("/api/v1/user/:id", userController.UpdateAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
	r.Patch("/api/v1/user/:id", userController.UpdateColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
	r.Delete("/api/v1/user/:id", userController.Delete, middleware.Authenticate(r.auth, auth.RoleAdmin))

	// #republic
	r.Get("/api/v1/republic/list", republicController.GetList, middleware.Authenticate(r.auth))
	r.Get("/api/v1/republic/:id", republicController.GetRepublicDetailById, middleware.Authenticate(r.auth))
	r.Post("/api/v1/republic/create", republicController.CreateRepublic, middleware.Authenticate(r.auth))
	r.Put("/api/v1/republic/:id", republicController.UpdateRepublicAll, middleware.Authenticate(r.auth))
	r.Patch("/api/v1/republic/:id", republicController.UpdateRepublicColumns, middleware.Authenticate(r.auth))
	r.Delete("/api/v1/republic/:id", republicController.DeleteRepublic, middleware.Authenticate(r.auth))

	// #department
	r.Get("/api/v1/department/list", departmentController.GetList, middleware.Authenticate(r.auth))
	r.Get("/api/v1/department/:id", departmentController.GetDetailById, middleware.Authenticate(r.auth))
	r.Post("/api/v1/department/create", departmentController.Create, middleware.Authenticate(r.auth))
	r.Put("/api/v1/department/:id", departmentController.UpdateAll, middleware.Authenticate(r.auth))
	r.Patch("/api/v1/department/:id", departmentController.UpdateColumns, middleware.Authenticate(r.auth))
	r.Delete("/api/v1/department/:id", departmentController.Delete, middleware.Authenticate(r.auth))

	// #position
	r.Get("/api/v1/position/list", positionController.GetList, middleware.Authenticate(r.auth))
	r.Get("/api/v1/position/:id", positionController.GetDetailById, middleware.Authenticate(r.auth))
	r.Post("/api/v1/position/create", positionController.Create, middleware.Authenticate(r.auth))
	r.Put("/api/v1/position/:id", positionController.UpdateAll, middleware.Authenticate(r.auth))
	r.Patch("/api/v1/position/:id", positionController.UpdateColumns, middleware.Authenticate(r.auth))
	r.Delete("/api/v1/position/:id", positionController.Delete, middleware.Authenticate(r.auth))

	// #region
	r.Get("/api/v1/region/list", regionController.GetRegionList, middleware.Authenticate(r.auth))
	r.Get("/api/v1/region/:id", regionController.GetRegionDetailById, middleware.Authenticate(r.auth))
	r.Post("/api/v1/region/create", regionController.CreateRegion, middleware.Authenticate(r.auth))
	r.Put("/api/v1/region/:id", regionController.UpdateRegionAll, middleware.Authenticate(r.auth))
	r.Patch("/api/v1/region/:id", regionController.UpdateRegionColumns, middleware.Authenticate(r.auth))
	r.Delete("/api/v1/region/:id", regionController.DeleteRegion, middleware.Authenticate(r.auth))
	r.Get("/api/v1/region/list/by/republic/:republic_id", regionController.GetRegionByRepublicIDList, middleware.Authenticate(r.auth))

	// #district
	r.Get("/api/v1/district/list", districtController.GetList, middleware.Authenticate(r.auth))
	r.Get("/api/v1/district/:id", districtController.GetDetailById, middleware.Authenticate(r.auth))
	r.Post("/api/v1/district/create", districtController.Create, middleware.Authenticate(r.auth))
	r.Put("/api/v1/district/:id", districtController.UpdateAll, middleware.Authenticate(r.auth))
	r.Patch("/api/v1/district/:id", districtController.UpdateColumns, middleware.Authenticate(r.auth))
	r.Delete("/api/v1/district/:id", districtController.Delete, middleware.Authenticate(r.auth))
	r.Get("/api/v1/district/list/by/region/:region_id", districtController.GetListByRegionID, middleware.Authenticate(r.auth))

	return r.Run(r.port)
}
