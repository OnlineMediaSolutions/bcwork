package export

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

func (e *ExportModule) ExportCSV(ctx context.Context, srcs []json.RawMessage) ([]byte, error) {
	var (
		head header
		buf  bytes.Buffer
	)

	csvwriter := csv.NewWriter(&buf)
	headerData := srcs[0]

	// getting header
	err := json.Unmarshal(headerData, &head)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal %v to header: %w", string(headerData), err)
	}

	// writing header and first row
	err = csvwriter.Write(head.keys)
	if err != nil {
		return nil, fmt.Errorf("cannot write headers to buffer: %w", err)
	}
	err = csvwriter.Write(head.values)
	if err != nil {
		return nil, fmt.Errorf("cannot write first row to buffer: %w", err)
	}

	// writing subsequent rows
	for i := 1; i < len(srcs); i++ {
		var temp map[string]interface{}
		err := json.Unmarshal(srcs[i], &temp)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal %v to map[string]interface{}", srcs[i])
		}

		if len(temp) != len(head.keys) {
			return nil, fmt.Errorf(
				"cannot process different objects: len(headers) [%v] != len(row) [%v]",
				len(head.keys), len(temp),
			)
		}

		row := make([]string, 0, len(head.keys))
		for _, key := range head.keys {
			var cell string
			v, ok := temp[key]
			if ok {
				cell = fmt.Sprint(v)
			}

			row = append(row, cell)
		}

		err = csvwriter.Write(row)
		if err != nil {
			return nil, fmt.Errorf("cannot write row #%v to buffer: %w", i+1, err)
		}
	}

	csvwriter.Flush()

	return buf.Bytes(), nil
}

type header struct {
	keys   []string
	values []string // values of first row
}

// UnmarshalJSON Unmarshalling simple structs (key:value) like
// `{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`
func (h *header) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	// must open with a delim token '{'
	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{', received [%v]", t)
	}

	var (
		tempKeys   []string
		tempValues []string
	)

	for decoder.More() {
		keyRaw, err := decoder.Token()
		if err != nil {
			return err
		}

		valueRaw, err := decoder.Token()
		if err != nil {
			return err
		}

		key, ok := keyRaw.(string)
		if !ok {
			return fmt.Errorf("cannot cast %v to string", keyRaw)
		}

		tempKeys = append(tempKeys, key)
		tempValues = append(tempValues, fmt.Sprint(valueRaw))
	}

	// must close with a delim token '}'
	t, err = decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object open with '}', received [%v]", t)
	}

	h.keys = tempKeys
	h.values = tempValues

	return nil
}
