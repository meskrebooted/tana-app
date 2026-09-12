package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"coworking-backend/config"
	"coworking-backend/models"

	"github.com/gin-gonic/gin"
)

func ListSpaces(c *gin.Context) {
	var spaces []models.Space
	if err := config.DB.Order("id").Find(&spaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile recuperare gli spazi"})
		return
	}
	c.JSON(http.StatusOK, spaces)
}

type spaceInput struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=desk room booth"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
	Location string `json:"location"`
}

func CreateSpace(c *gin.Context) {
	var in spaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	space := models.Space{
		Name:     in.Name,
		Type:     in.Type,
		Capacity: in.Capacity,
		Location: in.Location,
	}
	if err := config.DB.Create(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile creare lo spazio"})
		return
	}
	c.JSON(http.StatusCreated, space)
}

// RandomAvailableSpace picks, at random, a space with no confirmed booking
// covering the next hour — comodo per chi non vuole scegliere e vuole solo
// sedersi da qualche parte adesso.
func RandomAvailableSpace(c *gin.Context) {
	now := time.Now()
	soon := now.Add(time.Hour)

	var spaces []models.Space
	if err := config.DB.Find(&spaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile recuperare gli spazi"})
		return
	}

	var free []models.Space
	for _, s := range spaces {
		var conflicts int64
		config.DB.Model(&models.Booking{}).
			Where("space_id = ? AND status = ? AND start_time < ? AND end_time > ?",
				s.ID, "confirmed", soon, now).
			Count(&conflicts)
		if conflicts == 0 {
			free = append(free, s)
		}
	}

	if len(free) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "nessuno spazio libero nella prossima ora, oggi va di moda la prenotazione"})
		return
	}

	pick := free[rand.Intn(len(free))]
	c.JSON(http.StatusOK, pick)
}

func DeleteSpace(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Space{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile eliminare lo spazio"})
		return
	}
	c.Status(http.StatusNoContent)
}
