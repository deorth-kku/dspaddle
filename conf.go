package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path"

	v1 "encoding/json"
)

type (
	List[T any] []T
	StringList  = List[string]
)

type unmarshaler[T any] interface {
	*T
	json.UnmarshalerFrom
}

func UnmarshalV1[T any, P unmarshaler[T]](data []byte, dst P) error {
	return json.Unmarshal(data, dst, v1.DefaultOptionsV1())
}

func MarshalV1[T json.MarshalerTo](v T) ([]byte, error) {
	return json.Marshal(v, v1.DefaultOptionsV1())
}

func (sl *List[T]) UnmarshalJSON(data []byte) error {
	return UnmarshalV1(data, sl)
}

func (sl *List[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case 'n':
		*sl = nil
		return dec.SkipValue()
	case '[':
		return json.UnmarshalDecode(dec, (*[]T)(sl))
	default:
		*sl = make(List[T], 1)
		return json.UnmarshalDecode(dec, &((*sl)[0]))
	}
}

func (sl List[T]) MarshalJSON() ([]byte, error) {
	return MarshalV1(sl)
}

func (sl List[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if len(sl) == 1 {
		return json.MarshalEncode(enc, sl[0])
	}
	return json.MarshalEncode(enc, []T(sl))
}

type LogConfig struct {
	File  string `json:"file,omitzero"`
	Level string `json:"level,omitzero"`
}

type KeysConfig struct {
	Left  StringList `json:"left,omitzero"`
	Right StringList `json:"right,omitzero"`
	Mode  Mode       `json:"mode,omitzero"`
}

type Mode uint

const (
	ModeSendInput = iota
	ModePostMessage
)

type Config struct {
	Log  LogConfig  `json:"log,omitzero"`
	Keys KeysConfig `json:"keys,omitzero"`
}

func GetConfig() (*Config, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot find exe path err: %w", err)
	}
	dir := path.Dir(exe)
	conf := path.Join(dir, "dspaddle.json")
	f, err := os.Open(conf)
	if err != nil {
		return nil, fmt.Errorf("fail when open config file %s, err: %w", conf, err)
	}
	defer f.Close()
	dec := jsontext.NewDecoder(f)
	c := new(Config)
	return c, json.UnmarshalDecode(dec, c)
}
