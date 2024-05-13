package singletonmanager

import (
	"reflect"
	"sync"
)

type IInstances interface {
	Init() IInstances
}

type ISingletonManager interface {
	GetInstance(IInstances) interface{}
	RegisterInstances(...IInstances)
}

type SingletonManager struct {
	resource map[string]interface{}
	sync.Mutex
}

func (s *SingletonManager) RegisterInstances(i ...IInstances) {
	s.Lock()
	defer s.Unlock()
	for _, k := range i {
		name := reflect.TypeOf(k).String()
		s.resource[name[1:]] = k.Init()
	}
}

func (s *SingletonManager) GetInstance(k IInstances) interface{} {
	s.Lock()
	defer s.Unlock()
	name := reflect.TypeOf(k).String()
	if v, ok := s.resource[name[1:]]; ok {
		return v
	}
	return nil
}
