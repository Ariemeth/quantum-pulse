package entity

import "github.com/Ariemeth/quantum-pulse/engine/components"

// ComponentRequest presents an interface used to request a Component from an Entity.
type ComponentRequest interface {
	Name() string
	Channel() chan components.Component
}

type compRequest struct {
	name    string
	retChan chan components.Component
	done    chan interface{}
}

func (c compRequest) Name() string {
	return c.name
}

func (c *compRequest) Channel() chan components.Component {
	return c.retChan
}

func (c *compRequest) Done() {

}
