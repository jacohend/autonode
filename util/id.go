package util

import "github.com/oklog/ulid/v2"

func BytesToUlid(idbytes []byte) (id ulid.ULID) {
	id = ulid.ULID{}
	id.UnmarshalBinary(idbytes)
	return
}
