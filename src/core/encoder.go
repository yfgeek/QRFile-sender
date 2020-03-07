package core

import (
	"encoding/ascii85"
	"encoding/base64"
)

type Encoder interface {
	Encode([]byte) (string, error)
}

type Base64Encoder struct{}
type ASCII85Encoder struct{}

func (Base64Encoder) Encode(in []byte) (string, error){
	return base64.StdEncoding.EncodeToString(in), nil
}

func (ASCII85Encoder) Encode(in []byte) (string, error){
	var dst []byte
	ascii85.Encode(dst,in)
	return string(dst), nil
}