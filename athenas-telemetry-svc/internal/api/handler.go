package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

func InitHandlers(router *gin.Engine) {
	router.POST("/v1/telemetry", HandleRecordTelemetry)
}

// ------------------------------------------------------------------------------------------------

func HandleRecordTelemetry(ctx *gin.Context) {
	var data specifications.Telemetry

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: do something with data

	ctx.Status(http.StatusAccepted)
}
