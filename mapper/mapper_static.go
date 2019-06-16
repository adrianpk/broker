package mapper

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mikrowezel/backend/broker"
)

// MapMessage decodes different serialized entities into a Message using the appropriate decoder.
// TODO: Make types pluggable.
func (sm *StaticMapper) MapMessage(messageTypeID string, serialized interface{}) (broker.BaseMessage, error) {
	var bm broker.BaseMessage

	// There wil bm a case for each broker BaseMessage concrete implementation.
	// TODO: Make types pluggable.
	switch messageTypeID {
	case "message":
		bm = &broker.Message{}

	case "text":
		bm = &broker.Text{}

	default:
		return nil, fmt.Errorf("unknown message type '%s'", messageTypeID)
	}

	switch s := serialized.(type) {
	case []byte:
		err := json.Unmarshal(s, bm)
		if err != nil {
			return nil, fmt.Errorf("message %s cannot bm unmarshalled: %s", messageTypeID, err)
		}
	default:
		cfg := mapstructure.DecoderConfig{
			Result:  bm,
			TagName: "json",
		}
		dec, err := mapstructure.NewDecoder(&cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize a decoder for message %s. %s", messageTypeID, err)
		}

		err = dec.Decode(s)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal message %s: %s", messageTypeID, err)
		}
	}

	return bm, nil
}
