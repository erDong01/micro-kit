package base

import (
	"fmt"
	"rxt/internal/core"
)

type Service struct {
	*core.Core
	InitCore bool
	A        interface{}
}

func (s *Service) Init() *Service {
	if s.InitCore == false {
		s.Core = core.New()
		fmt.Print(s.Core)
		s.InitCore = true
	}
	return s
}
