package textmap

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type TextMap interface {
	LoadJson(reader io.Reader) error
	Of(key string) TextEntry
}

type TextEntry string

type _TextMap struct {
	table map[string]string
}

func New() TextMap {
	return new(_TextMap)
}

func (tm *_TextMap) LoadJson(reader io.Reader) error {
	table := make(map[string]string)

	if err := json.NewDecoder(reader).Decode(&table); err != nil {
		return err
	}

	tm.table = table

	return nil
}

func (tm *_TextMap) Of(key string) TextEntry {
	return TextEntry(tm.table[key])
}

func (te TextEntry) _set(key, format string, value any) TextEntry {
	from := fmt.Sprintf("{{%s}}", key)
	to := fmt.Sprintf(format, value)
	return TextEntry(strings.ReplaceAll(string(te), from, to))
}

func (te TextEntry) Set(key string, value any) TextEntry {
	return te._set(key, "%s", value)
}

func (te TextEntry) SetInt(key string, value any) TextEntry {
	return te._set(key, "%d", value)
}

func (te TextEntry) String() string {
	return string(te)
}
