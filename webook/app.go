package main

import (
	"basic-go/webook/internal/events"

	"github.com/robfig/cron/v3"

	"github.com/gin-gonic/gin"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
