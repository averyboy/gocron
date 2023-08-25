package container

import (
	"golang.org/x/exp/constraints"
)

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer Signed & Unsigned
type Integer interface {
	Signed | Unsigned
}

// Float float32 & float64
type Float interface {
	~float32 | ~float64
}

// String string & byte
type String interface {
	~string | ~byte
}

// InnerType all golang base type
type InnerType interface {
	constraints.Float | constraints.Integer | String
}
