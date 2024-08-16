package main

import (
	"basic-go/webook/internal/events"
	"basic-go/webook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
