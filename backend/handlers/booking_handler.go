package handlers

import (
	"net/http"
	"sort"
	"time"

	"coworking-backend/config"
	"coworking-backend/models"

	"github.com/gin-gonic/gin"
)

type bookingInput struct {
	SpaceID   uint      `json:"space_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

// CreateBooking books a space after checking it isn't already taken in that slot.
func CreateBooking(c *gin.Context) {
	var in bookingInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !in.EndTime.After(in.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "l'orario di fine deve essere dopo l'orario di inizio"})
		return
	}

	var space models.Space
	if err := config.DB.First(&space, in.SpaceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "spazio non trovato"})
		return
	}

	// Overlap check: any confirmed booking on the same space whose interval intersects.
	var conflicts int64
	config.DB.Model(&models.Booking{}).
		Where("space_id = ? AND status = ? AND start_time < ? AND end_time > ?",
			in.SpaceID, "confirmed", in.EndTime, in.StartTime).
		Count(&conflicts)

	if conflicts > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "lo spazio è già prenotato in questa fascia oraria"})
		return
	}

	userID, _ := c.Get("user_id")
	booking := models.Booking{
		UserID:    userID.(uint),
		SpaceID:   in.SpaceID,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Status:    "confirmed",
	}
	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile creare la prenotazione"})
		return
	}

	config.DB.Preload("Space").First(&booking, booking.ID)
	c.JSON(http.StatusCreated, booking)
}

// ListMyBookings returns all bookings for the authenticated user.
func ListMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var bookings []models.Booking
	if err := config.DB.Preload("Space").
		Where("user_id = ?", userID).
		Order("start_time desc").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile recuperare le prenotazioni"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Streak calcola per quanti giorni consecutivi (fino a oggi o ieri) l'utente
// ha almeno una prenotazione confermata: una piccola gamification per chi
// frequenta lo spazio con costanza.
func Streak(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var bookings []models.Booking
	if err := config.DB.
		Where("user_id = ? AND status = ?", userID, "confirmed").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile calcolare lo streak"})
		return
	}

	days := map[string]bool{}
	for _, b := range bookings {
		days[b.StartTime.Format("2006-01-02")] = true
	}

	dayKeys := make([]string, 0, len(days))
	for d := range days {
		dayKeys = append(dayKeys, d)
	}
	sort.Strings(dayKeys)

	today := time.Now().Truncate(24 * time.Hour)
	streak := 0
	cursor := today

	// Se oggi non ha ancora prenotato, lo streak riparte da ieri: non lo penalizziamo
	// finché la giornata non è ancora finita.
	if !days[cursor.Format("2006-01-02")] {
		cursor = cursor.AddDate(0, 0, -1)
	}

	for days[cursor.Format("2006-01-02")] {
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}

	c.JSON(http.StatusOK, gin.H{"streak_days": streak})
}

// CancelBooking marks a booking (owned by the caller) as cancelled.
func CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var booking models.Booking
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prenotazione non trovata"})
		return
	}

	booking.Status = "cancelled"
	if err := config.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile annullare la prenotazione"})
		return
	}
	c.JSON(http.StatusOK, booking)
}
