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

func lCollectorBucket(l *lua.State) int {
	col := toCollector(l, 1)
	name := enckey(l, 2)
	bucket := col.Bucket(name)
	if bucket == nil {
		if col.Writable() {
			var e error
			bucket, e = col.CreateBucket(name)
			if e != nil {
				panic(e)
			}
		} else {
			panic("bolt: bucket not found: " + l.ToString(2))
		}
	}
	l.Push(bucket)
	l.PushIndex(lua.FirstUpVal - 1)
	l.SetMetaTable(-2)
	return 1
}

func lCollectorBuckets(l *lua.State) int {
	var k, v []byte
	col := toCollector(l, 1)
	c := col.Cursor()
	if l.IsNil(2) {
		k, v = c.First()
	} else {
		k, v = c.Seek(enckey(l, 2))
	}
	l.PushClosure(func(l *lua.State) int {
		for {
			if k == nil {
				return 0
			}
			if v == nil {
				break
			}
			k, v = c.Next()
		}
		deckey(l, k)
		l.Push(col.Bucket(k))
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		k, v = c.Next()
		return 2
	}, lua.FirstUpVal-1)
	return 1
}

func lCollectorDelete(l *lua.State) int {
	e := toCollector(l, 1).DeleteBucket(enckey(l, 2))
	if e == nil {
		return 0
	} else {
		l.Push(e.Error())
		return 1
	}
}
