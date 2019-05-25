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
	"ofunc/lua"
)

func metaTx(l *lua.State, mbucket int) int {
	l.NewTable(0, 4)
	idx := l.AbsIndex(-1)

	l.Push("bucket")
	l.PushClosure(lTxBucket, mbucket)
	l.SetTableRaw(-3)

	l.Push("buckets")
	l.PushClosure(lTxBuckets, mbucket)
	l.SetTableRaw(-3)

	l.Push("delete")
	l.Push(lTxDelete)
	l.SetTableRaw(-3)

	l.Push("writeto")
	l.Push(lTxWriteTo)
	l.SetTableRaw(-3)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(-3)

	return idx
}

func lTxBucket(l *lua.State) int {
	tx := toTx(l, 1)
	sname := l.ToString(2)
	bname := []byte(sname)
	bucket := tx.Bucket(bname)
	if bucket == nil {
		if tx.Writable() {
			var e error
			bucket, e = tx.CreateBucket(bname)
			if e != nil {
				panic(e)
			}
		} else {
			panic("bolt: bucket not found: " + sname)
		}
	}
	l.Push(bucket)
	l.PushIndex(lua.FirstUpVal - 1)
	l.SetMetaTable(-2)
	return 1
}

func lTxBuckets(l *lua.State) int {
	var k []byte
	tx := toTx(l, 1)
	c := tx.Cursor()
	if l.IsNil(2) {
		k, _ = c.First()
	} else {
		k, _ = c.Seek([]byte(l.ToString(2)))
	}
	l.PushClosure(func(l *lua.State) int {
		if k == nil {
			return 0
		}
		l.Push(string(k))
		l.Push(tx.Bucket(k))
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		k, _ = c.Next()
		return 2
	}, lua.FirstUpVal-1)
	return 1
}

func lTxDelete(l *lua.State) int {
	e := toTx(l, 1).DeleteBucket([]byte(l.ToString(2)))
	if e == nil {
		return 0
	} else {
		l.Push(e.Error())
		return 1
	}
}

func lTxWriteTo(l *lua.State) int {
	n, e := toTx(l, 1).WriteTo(toWriter(l, 2))
	l.Push(n)
	if e == nil {
		return 1
	} else {
		l.Push(e.Error())
		return 2
	}
}
