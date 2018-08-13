package main

type entityData struct {
	shareCount uint64
}

func (e *entityData) SetData(d string, t finDataType) {
	switch t {
	case finDataSharesOutstanding:
		e.shareCount = normalizeNumber(d)
	}
}
