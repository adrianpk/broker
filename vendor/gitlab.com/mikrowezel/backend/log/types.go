package log

import (
	"github.com/rs/zerolog"
)

var (
	cfg config
)

// Logger is structured leveled logger
type Logger struct {
	// Level of min logging
	Level int
	// Version
	Version string
	// Revision
	Revision string
	// DebugLog logger
	StdLog zerolog.Logger
	// ErrorLog logger
	ErrLog zerolog.Logger
	// Dynamic fields
	dynafields []interface{}
}

type config struct {
	// Name
	name string
	// Level of min logging
	level int
	// Static fields
	stfields []interface{}
	// configured
	configured bool
}

// It is advisable to define keys, not as plain strings,
// but as package unexported types to avoid name collisions.
// In this way those keys are not shared shared outside the package
// Context can store custom user and service specific data types.
// It gets it using interfaces, ergo: we lose type safety.
// If someone using the same public key (an int or a string i.e.)
// stores the context value using another type instead of
// a *logger object service will just continue along
// until it try to use a type assertion to get the now overwrited logger
// and then it will panic.
type contextKey string

// String is human readable representation of a context key.
func (c contextKey) String() string {
	return "hc-" + string(c)
}

// setup name and static fields.
// Each new instance of logger will always append these
// key-value pairs to the output and name if it is not empty.
// These values cannot be modified after they are configured.
func setup(name string, stfields []interface{}) {
	if cfg.configured {
		return
	}
	cfg.name = name
	cfg.stfields = append(cfg.stfields, stfields...)
	cfg.configured = true
}

// SetDyna fields.
// The receiver instance will always append these
// key-value pairs to the output.
func (l *Logger) SetDyna(dynafields ...interface{}) {
	l.dynafields = make([]interface{}, 2)
	l.dynafields = append(l.dynafields, dynafields...)
}

// AddDyna fields.
// The receiver instance will always append these
// key-value pairs to the output.
func (l *Logger) AddDyna(key, value interface{}) {
	l.dynafields = append(l.dynafields, []interface{}{key, value})
}

// ResetDyna fields.
// Remove dynamic fields.
func (l *Logger) ResetDyna() {
	l.dynafields = make([]interface{}, 2)
}
