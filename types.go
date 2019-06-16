package broker

// BaseMessage is a base broker message interface.
type BaseMessage interface {
	TypeID() string
}

// Message is a BaseMessage referecence implementation.
type Message struct {
}

// Text is a BaseMessage referecence implementation.
// for sending plain text.
type Text struct {
}

// TypeID returns Message type id.
func (m *Message) TypeID() string{
	return "message"
}

// TypeID returns Text type id.
func (t *Text) TypeID() string{
	return "text"
}