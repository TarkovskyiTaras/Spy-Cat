package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"spycat/internal/missions"
)

func (a *API) CreateMission(c *gin.Context) {
	var input missions.MissionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Targets) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A mission must have at least one target"})
		return
	}

	mission, err := a.missions.CreateMission(context.Background(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mission)
}

func (a *API) GetMission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission ID"})
		return
	}

	ctx := context.Background()
	mission, err := a.missions.GetMissionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	c.JSON(http.StatusOK, mission)
}

func (a *API) GetAllMissions(c *gin.Context) {
	ctx := context.Background()
	allMissions, err := a.missions.GetAllMissions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allMissions)
}

func (a *API) UpdateMission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission ID"})
		return
	}

	ctx := context.Background()
	mission, err := a.missions.GetMissionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.Complete {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update a completed mission"})
		return
	}

	var input missions.MissionUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mission, err = a.missions.UpdateMission(ctx, id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission successfully updated"})
}

func (a *API) DeleteMission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission ID"})
		return
	}

	ctx := context.Background()
	mission, err := a.missions.GetMissionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.CatID != uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete a mission assigned to a cat"})
		return
	}

	if err := a.missions.DeleteMission(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission deleted"})
}
