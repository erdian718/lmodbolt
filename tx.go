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
	l.NewTable(0, 8)
	idx := l.AbsIndex(-1)

	l.Push("bucket")
	l.PushClosure(lCollectorBucket, mbucket)
	l.SetTableRaw(-3)

	l.Push("buckets")
	l.PushClosure(lCollectorBuckets, mbucket)
	l.SetTableRaw(-3)

	l.Push("delete")
	l.Push(lCollectorDelete)
	l.SetTableRaw(-3)

	l.Push("writeto")
	l.Push(lTxWriteTo)
	l.SetTableRaw(-3)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(-3)

	return idx
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
