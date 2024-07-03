package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"spycat/internal/cats"
)

func (a *API) CreateCat(c *gin.Context) {
	var input cats.CatCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	breeds, err := cats.FetchValidBreeds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cat breeds"})
		return
	}

	if !cats.IsValidBreed(input.Breed, breeds) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat breed"})
		return
	}

	cat, err := a.cats.CreateCat(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (a *API) GetCats(c *gin.Context) {
	allCats, err := a.cats.GetAllCats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allCats)
}

func (a *API) GetCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat ID"})
		return
	}

	cat, err := a.cats.GetCatByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (a *API) UpdateCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat ID"})
		return
	}

	var updateCat cats.CatUpdateDTO
	if err := c.ShouldBindJSON(&updateCat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := a.cats.UpdateCat(id, &updateCat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (a *API) DeleteCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat ID"})
		return
	}

	if err := a.cats.DeleteCat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cat deleted"})
}
