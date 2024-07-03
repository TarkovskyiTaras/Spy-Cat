package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"spycat/internal/missions"
)

func (a *API) GetTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	ctx := context.Background()
	target, err := a.missions.GetTargetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, target)
}

func (a *API) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	var input missions.TargetUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	target, err := a.missions.GetTargetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if target.Complete {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update a completed target"})
		return
	}

	mission, err := a.missions.GetMissionByID(ctx, target.MissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.Complete {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update a target of a completed mission"})
		return
	}

	target, err = a.missions.UpdateTarget(ctx, id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}

func (a *API) AddTarget(c *gin.Context) {
	idStr := c.Param("mission_id")
	missionID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	ctx := context.Background()
	mission, err := a.missions.GetMissionByID(ctx, missionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.Complete {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add a target to a completed mission"})
		return
	}

	targets, err := a.missions.GetTargetsByMissionID(ctx, missionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(targets) >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A mission cannot have more than 3 targets"})
		return
	}

	var targetDTO missions.TargetCreateDTO
	if err := c.ShouldBindJSON(&targetDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := a.missions.AddTarget(ctx, missionID, &targetDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

func (a *API) DeleteTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	ctx := context.Background()
	target, err := a.missions.GetTargetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if target.Complete {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete a completed target"})
		return
	}

	if err := a.missions.DeleteTarget(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

func (a *API) GetTargetsByMissionID(c *gin.Context) {
	missionIDStr := c.Param("mission_id")
	missionID, err := uuid.Parse(missionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission ID"})
		return
	}

	ctx := context.Background()
	targets, err := a.missions.GetTargetsByMissionID(ctx, missionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, targets)
}
