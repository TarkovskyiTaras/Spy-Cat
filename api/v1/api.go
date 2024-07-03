package v1

import (
	"spycat/internal/cats"
	"spycat/internal/missions"
	"spycat/logger"

	"github.com/gin-gonic/gin"
)

type API struct {
	cats     *cats.CatService
	missions *missions.MissionService
}

func NewAPI(catService *cats.CatService, missionService *missions.MissionService) *API {
	return &API{
		cats:     catService,
		missions: missionService,
	}
}

func (a *API) SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(logger.Logger())

	router.POST("/cat", a.CreateCat)
	router.GET("/cats/all", a.GetCats)
	router.GET("/cat/:id", a.GetCat)
	router.PUT("/cat/:id", a.UpdateCat)
	router.DELETE("/cat/:id", a.DeleteCat)

	router.POST("/mission", a.CreateMission)
	router.GET("/mission/:id", a.GetMission)
	router.GET("/missions", a.GetAllMissions)
	router.PUT("/mission/:id", a.UpdateMission)
	router.DELETE("/mission/:id", a.DeleteMission)

	router.POST("/target/:mission_id", a.AddTarget)
	router.GET("/target/:mission_id/all", a.GetTargetsByMissionID)
	router.GET("/targets/:id", a.GetTarget)
	router.PUT("/targets/:id", a.UpdateTarget)
	router.DELETE("/targets/:id", a.DeleteTarget)

	return router
}
