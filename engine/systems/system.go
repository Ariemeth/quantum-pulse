// Package systems provides the various systems used by the engine
package systems

// System represents the behavior all engine systems should implement.
type System interface {
	SystemType() string
}
