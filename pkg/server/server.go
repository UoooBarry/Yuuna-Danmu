package server

import "uooobarry/yuuna-danmu/pkg/live"

type Server interface {
	Start(port int) error
	Stop() error
	Dispatch(event live.Event)
	IsRunning() bool
}
