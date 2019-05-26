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
	"ofunc/lmodmsgpack"
	"ofunc/lua"
)

func metaBucket(l *lua.State) int {
	l.NewTable(0, 8)
	idx := l.AbsIndex(-1)

	l.Push("bucket")
	l.PushClosure(lCollectorBucket, idx)
	l.SetTableRaw(idx)

	l.Push("buckets")
	l.PushClosure(lCollectorBuckets, idx)
	l.SetTableRaw(idx)

	l.Push("delbucket")
	l.Push(lCollectorDelBucket)
	l.SetTableRaw(idx)

	l.Push("get")
	l.Push(lBucketGet)
	l.SetTableRaw(idx)

	l.Push("set")
	l.Push(lBucketSet)
	l.SetTableRaw(idx)

	l.Push("pairs")
	l.Push(lBucketPairs)
	l.SetTableRaw(idx)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(idx)

	return idx
}

func lBucketGet(l *lua.State) int {
	xs := toBucket(l, 1).Get(enckey(l, 2))
	if xs == nil {
		return 0
	}
	_, e := lmodmsgpack.DecodeBytes(l, xs)
	if e != nil {
		panic(e.Error())
	}
	return 1
}

func lBucketSet(l *lua.State) int {
	var e error
	b := toBucket(l, 1)
	k := enckey(l, 2)
	if l.IsNil(3) {
		e = b.Delete(k)
	} else {
		e = b.Put(k, lmodmsgpack.EncodeBytes(l, 3))
	}
	if e != nil {
		panic(e.Error())
	}
	return 0
}

func lBucketPairs(l *lua.State) int {
	var k, v []byte
	c := toBucket(l, 1).Cursor()
	if l.IsNil(2) {
		k, v = c.First()
	} else {
		k, v = c.Seek(enckey(l, 2))
	}
	l.Push(func(l *lua.State) int {
		for {
			if k == nil {
				return 0
			}
			if v != nil {
				break
			}
			k, v = c.Next()
		}
		deckey(l, k)
		if _, e := lmodmsgpack.DecodeBytes(l, v); e != nil {
			panic(e.Error())
		}
		k, v = c.Next()
		return 2
	})
	return 1
}
