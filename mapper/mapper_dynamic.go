package mapper

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mikrowezel/backend/broker"
)

// DynamicMapper is a dynamic mapper struct.
type DynamicMapper struct {
	typeMap map[string]reflect.Type
}

// NewDynamicMapper returns a new dynamic BaseMessageMapper.
func NewDynamicMapper() BaseMessageMapper {
	return &DynamicMapper{
		typeMap: make(map[string]reflect.Type),
	}
}

// MapMessage maps broker messages to structs.
func (e *DynamicMapper) MapMessage(messageTypeID string, serialized interface{}) (broker.BaseMessage, error) {
	eType, ok := e.typeMap[messageTypeID]
	if !ok {
		return nil, fmt.Errorf("no mapping configured for message %s", messageTypeID)
	}

	instance := reflect.New(eType)
	ifc := instance.Interface()

	message, ok := ifc.(broker.BaseMessage)
	if !ok {
		return nil, fmt.Errorf("type %s does not implement the Message interface", ifc)
	}

	switch s := serialized.(type) {

	case []byte:
		err := json.Unmarshal(s, message)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal message %s. %s", messageTypeID, err)
		}

	default:
		cfg := mapstructure.DecoderConfig{
			Result:  message,
			TagName: "json",
		}
		dec, err := mapstructure.NewDecoder(&cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize decoder for message %s. %s", messageTypeID, err)
		}

		err = dec.Decode(s)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal message %s. %s", messageTypeID, err)
		}
	}

	return message, nil
}

// RegMapping let register a mapping.
func (e *DynamicMapper) RegMapping(messageType reflect.Type) error {
	instance := reflect.New(messageType).Interface()
	message, ok := instance.(broker.Message)

	if !ok {
		return fmt.Errorf("type %s does not implement the Message interface", instance)
	}

	e.typeMap[message.TypeID()] = messageType
	return nil
}
