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
	"errors"
	"fmt"

	"ofunc/lua"

	"github.com/boltdb/bolt"
)

func metaDB(l *lua.State, mtx int) int {
	l.NewTable(0, 4)
	idx := l.AbsIndex(-1)

	l.Push("close")
	l.Push(lDBClose)
	l.SetTableRaw(idx)

	l.Push("view")
	l.PushClosure(lDBView, mtx)
	l.SetTableRaw(idx)

	l.Push("update")
	l.PushClosure(lDBUpdate, mtx)
	l.SetTableRaw(idx)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(idx)

	return idx
}

func lDBClose(l *lua.State) int {
	db := toDB(l, 1)
	path := db.Path()
	mutex.Lock()
	defer mutex.Unlock()
	if x, ok := cache[path]; ok {
		if x.ref > 1 {
			x.ref -= 1
			return 0
		} else {
			if e := x.db.Close(); e == nil {
				delete(cache, path)
				return 0
			} else {
				l.Push(e.Error())
				return 1
			}
		}
	} else {
		l.Push(bolt.ErrDatabaseNotOpen.Error())
		return 1
	}
}

func lDBView(l *lua.State) int {
	e := toDB(l, 1).View(func(tx *bolt.Tx) error {
		l.PushIndex(2)
		l.Push(tx)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		if msg := l.PCall(1, 1, false); msg == nil {
			if l.IsNil(-1) {
				return nil
			} else {
				return errors.New(l.ToString(-1))
			}
		} else {
			return fmt.Errorf("%v", msg)
		}
	})
	if e == nil {
		return 0
	} else {
		l.Push(e.Error())
		return 1
	}
}

func lDBUpdate(l *lua.State) int {
	e := toDB(l, 1).Update(func(tx *bolt.Tx) error {
		l.PushIndex(2)
		l.Push(tx)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		if msg := l.PCall(1, 1, false); msg == nil {
			if l.IsNil(-1) {
				return nil
			} else {
				return errors.New(l.ToString(-1))
			}
		} else {
			return fmt.Errorf("%v", msg)
		}
	})
	if e == nil {
		return 0
	} else {
		l.Push(e.Error())
		return 1
	}
}
