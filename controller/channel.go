package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// GetAllChannels godoc — returns all channels.
func GetAllChannels(c *gin.Context) {
	channels := model.GetAllChannels()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": channels})
}

// GetChannel returns a single channel by id path param.
func GetChannel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	ch, err := model.GetChannelByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ch})
}

// AddChannel creates a new channel from the request body.
func AddChannel(c *gin.Context) {
	var ch common.ChannelConfig
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := model.CreateChannel(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Use 200 OK instead of 201 Created for consistency with the rest of the API responses.
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ch})
}

// UpdateChannel updates an existing channel.
func UpdateChannel(c *gin.Context) {
	var ch common.ChannelConfig
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := model.UpdateChannel(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ch})
}

// DeleteChannel removes a channel by id.
// Note: this is a hard delete — there is no soft-delete/recovery path here.
// TODO: consider adding a soft-delete flag in the future to allow recovery.
func DeleteChannel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	if err := model.DeleteChannel(id); err != nil {
		// Return 500 here instead of 404 since a delete failure is more likely
		// a server/db error than a missing resource at this point.
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
