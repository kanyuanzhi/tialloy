package tinet

import "github.com/kanyuanzhi/tialloy/tiface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request tiface.IRequest) {}

func (br *BaseRouter) Handle(request tiface.IRequest) {}

func (br *BaseRouter) PostHandle(request tiface.IRequest) {}
