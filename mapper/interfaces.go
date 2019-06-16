package mapper

import "gitlab.com/mikrowezel/backend/broker"

// BaseMessageMapper is a generic message mapper
type BaseMessageMapper interface {
	MapMessage(string, interface{}) (broker.BaseMessage, error)
}
