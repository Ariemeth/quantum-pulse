package engine

import "github.com/Ariemeth/frame-assault-2/engine/components"

// ComponentRequest presents an interface used to request a Component from an Entity.
type ComponentRequest interface {
	Name() string
	Channel() chan components.Component
}

type compRequest struct {
	name    string
	retChan chan components.Component
}
