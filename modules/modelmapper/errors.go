package modelmapper

import "errors"

var ErrUnsupportedPK = errors.New("unsupported primary key")

var ErrUnsuportedIDRequired = errors.New("model mapper id required")
