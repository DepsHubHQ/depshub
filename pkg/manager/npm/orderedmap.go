package npm

import (
	"bytes"
	"encoding/json"
)

type OrderedMap struct {
	Order    []string
	Values   map[string]string
	LineNums map[string]int
	RawLines map[string]string
}

func (o *OrderedMap) UnmarshalJSON(data []byte) error {
	// Initialize maps
	o.Values = make(map[string]string)
	o.LineNums = make(map[string]int)
	o.RawLines = make(map[string]string)

	// Pre-split data into lines for faster lookup
	lines := bytes.Split(data, []byte{'\n'})

	// Decode values
	if err := json.Unmarshal(data, &o.Values); err != nil {
		return err
	}

	// Track order and metadata
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.Token() // skip opening brace

	for dec.More() {
		pos := dec.InputOffset()
		key, _ := dec.Token()
		dec.Token() // skip value
		keyStr := key.(string)
		o.Order = append(o.Order, keyStr)

		// Calculate line number relative to the start of this block
		lineNum := 1 + bytes.Count(data[:pos], []byte{'\n'})
		o.LineNums[keyStr] = lineNum

		// Store raw line
		o.RawLines[keyStr] = string(bytes.TrimSpace(lines[lineNum]))
	}

	return nil
}
