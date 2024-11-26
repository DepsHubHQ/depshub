package npm

import (
	"bytes"
	"encoding/json"
)

// OrderedMap preserves the order of map keys
type OrderedMap struct {
	Order  []string
	Values map[string]string
}

func (o *OrderedMap) UnmarshalJSON(data []byte) error {
	o.Values = make(map[string]string)
	// Use a temporary map for decoding
	var tmp map[string]string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// Decode again to get keys in order
	dec := json.NewDecoder(bytes.NewReader(data))
	// Read opening brace
	_, err := dec.Token()
	if err != nil {
		return err
	}

	// Read key-value pairs in order
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return err
		}
		// Skip the value, we already have it in our map
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return err
		}
		keyStr := key.(string)
		o.Order = append(o.Order, keyStr)
		o.Values[keyStr] = tmp[keyStr]
	}
	return nil
}
