package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"gpltrans/models"
)

// CreateTrans 创建转换
//
//	@Summary		创建转换
//	@Description	创建转换
//	@Tags			trans
//	@Accept			json
//	@Produce		json
//	@Param			transModel	body		models.TransRequest		true	"transModel"
//	@Success		200			{object}	models.TransResponse	"Successful response"
//	@Router			/api/v1/trans/create [post]
func CreateTrans(c *gin.Context) {
	var transModel models.TransRequest
	err := json.NewDecoder(c.Request.Body).Decode(&transModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var result models.TransResponse
	result.Id = transModel.Id
	result.Source = transModel.Source
	result.Target = transModel.Source
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
