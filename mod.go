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

// boltdb/bolt bindings for Lua.
package lmodbolt

import (
	"path/filepath"
	"sync"

	"ofunc/lua"

	"github.com/boltdb/bolt"
)

type info struct {
	ref int
	db  *bolt.DB
}

var cache = map[string]*info{}
var mutex sync.Mutex

// Open opens the module.
func Open(l *lua.State) int {
	mbucket := metaBucket(l)
	mtx := metaTx(l, mbucket)
	mdb := metaDB(l, mtx)
	l.NewTable(0, 2)

	l.Push("version")
	l.Push("0.0.1")
	l.SetTableRaw(-3)

	l.Push("open")
	l.PushClosure(lOpen, mdb)
	l.SetTableRaw(-3)

	return 1
}

func lOpen(l *lua.State) int {
	path := l.ToString(1)
	if p, e := filepath.Abs(path); e == nil {
		path = p
	}
	mutex.Lock()
	defer mutex.Unlock()
	if x, ok := cache[path]; ok {
		x.ref += 1
		l.Push(x.db)
	} else {
		db, e := bolt.Open(path, 0666, nil)
		if e != nil {
			l.Push(nil)
			l.Push(e.Error())
			return 2
		}
		x = &info{
			db:  db,
			ref: 1,
		}
		cache[db.Path()] = x
		l.Push(db)
	}
	l.PushIndex(lua.FirstUpVal - 1)
	l.SetMetaTable(-2)
	return 1
}
