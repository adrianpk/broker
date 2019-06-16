package mapper

// NewMessageMapper creates a new MessageMapper
func NewMessageMapper() BaseMessageMapper {
	return &StaticMapper{}
}
