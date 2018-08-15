package main

import (
	"errors"
)

type entityData struct {
	shareCount uint64
}

func (e *entityData) SetData(d string, t finDataType) error {
	switch t {
	case finDataSharesOutstanding:
		e.shareCount = normalizeNumber(d)
		if e.shareCount <= 0 {
			return errors.New("Not the share count data")
		}
	}
	return nil
}
