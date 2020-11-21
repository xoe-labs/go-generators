package policy

import (
	"context"
)

// Policer knows how to make decisions on access policy
type Policer interface {
	// Can answers the question wether a policeable actor can perform an action given the current json encoded domain model
	Can(ctx context.Context, p Policeable, action string, data []byte) bool
}

// Policeable is an actor that can be policed
type Policeable interface {
	// User knows how to access the user
	User() string
	// ElevationToken knows how to access an optional elevation token
	ElevationToken() string
}
