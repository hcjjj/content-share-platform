package main

import (
	"basic-go/webook/internal/events"

	"github.com/gin-gonic/gin"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
