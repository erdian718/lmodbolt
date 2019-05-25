/*
Copyright 2019 by ofunc

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package lmodbolt

import (
	"encoding/binary"
	"io"

	"ofunc/lua"

	"github.com/boltdb/bolt"
)

type collector interface {
	Bucket(name []byte) *bolt.Bucket
	CreateBucket(name []byte) (*bolt.Bucket, error)
	Cursor() *bolt.Cursor
	DeleteBucket(name []byte) error
	Writable() bool
}

var endian = binary.LittleEndian

func toDB(l *lua.State, i int) *bolt.DB {
	if v, ok := l.GetRaw(1).(*bolt.DB); ok {
		return v
	} else {
		panic("bolt: invalid database: " + l.ToString(i))
	}
}

func toTx(l *lua.State, i int) *bolt.Tx {
	if v, ok := l.GetRaw(1).(*bolt.Tx); ok {
		return v
	} else {
		panic("bolt: invalid tx: " + l.ToString(i))
	}
}

func toBucket(l *lua.State, i int) *bolt.Bucket {
	if v, ok := l.GetRaw(1).(*bolt.Bucket); ok {
		return v
	} else {
		panic("bolt: invalid bucket: " + l.ToString(i))
	}
}

func toCollector(l *lua.State, i int) collector {
	if v, ok := l.GetRaw(1).(collector); ok {
		return v
	} else {
		panic("bolt: invalid tx or bucket: " + l.ToString(i))
	}
}

func toWriter(l *lua.State, i int) io.Writer {
	if v, ok := l.GetRaw(1).(io.Writer); ok {
		return v
	} else {
		panic("bolt: invalid writer: " + l.ToString(i))
	}
}

func deckey(l *lua.State, xs []byte) {
	switch xs[0] {
	case 1:
		l.Push(int64(endian.Uint64(xs[1:])))
	case 2:
		l.Push(string(xs[1:]))
	default:
		panic("bolt: invalid key data")
	}
}

func enckey(l *lua.State, i int) (xs []byte) {
	switch t := l.TypeOf(i); t {
	case lua.TypeNumber:
		if v, e := l.TryInteger(i); e == nil {
			xs = make([]byte, 9)
			xs[0] = 1
			endian.PutUint64(xs[1:], uint64(v))
		} else {
			panic("bolt: invalid key type: float")
		}
	case lua.TypeString:
		v := l.ToString(i)
		xs = make([]byte, len(v)+1)
		xs[0] = 2
		copy(xs[1:], v)
	default:
		panic("bolt: invalid key type: " + t.String())
	}
	return xs
}
