package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONStringSlice is a []string that stores itself as a JSON TEXT value in
// SQLite (e.g. ["forge","go"]). Use it for any []string column — see S4.
type JSONStringSlice []string

// Value implements driver.Valuer. A nil or empty slice is stored as "[]".
func (s JSONStringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, fmt.Errorf("JSONStringSlice marshal: %w", err)
	}
	return string(b), nil
}

// Scan implements sql.Scanner. Accepts a TEXT column containing a JSON array,
// or NULL (treated as empty slice).
func (s *JSONStringSlice) Scan(src any) error {
	if src == nil {
		*s = JSONStringSlice{}
		return nil
	}
	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("JSONStringSlice: unsupported type %T", src)
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return fmt.Errorf("JSONStringSlice unmarshal: %w", err)
	}
	*s = out
	return nil
}
