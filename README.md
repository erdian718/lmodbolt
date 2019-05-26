# lmodbolt

[boltdb/bolt](https://github.com/boltdb/bolt) bindings for [Lua](https://github.com/ofunc/lua).

## Usage

```go
package main

import (
	"ofunc/lmodbolt"
	"ofunc/lua/util"
)

func main() {
	l := util.NewState()
	l.Preload("bolt", lmodbolt.Open)
	util.Run(l, "main.lua")
}
```

```lua
local bolt = require 'bolt'

local db = bolt.open('test.db')
db:update(function(tx)
	local demo = tx:bucket('demo')
	demo:set('foo', {
		id = 123;
		name = 'test';
		valid = true;
	})

	local foo = tx:bucket('demo'):get('foo')
	print(foo.id, foo.name, foo.valid)
end)
db:close()
```

## Dependencies

* [ofunc/lua](https://github.com/ofunc/lua)
* [ofunc/lmodmsgpack](https://github.com/ofunc/lmodmsgpack)
* [boltdb/bolt](https://github.com/boltdb/bolt)

## Documentation

### bolt.open(path)

Creates and opens a database at the given `path`.
If the file does not exist then it will be created automatically.

### db:close()

Releases all database resources.
All transactions must be closed before closing the database.

### db:view(f)

Executes the function `f` within the context of a managed read-only transaction.
Any error in function `f` is returned from the `view` method.

### db:update(f)

Executes the function `f` within the context of a read-write managed transaction.
If no error in the function `f` then the transaction is committed.
Otherwise the entire transaction is rolled back.
Any error in function `f` is returned from the `update` method.

### tx:writeto(w)

Writes the entire database to a writer.

### c:bucket(key)

Retrieves a nested bucket by `key`.
If the bucket does not exist then it will be created automatically.
`c` is a tx or a bucket.

### c:buckets([key])

Returns an iterator that can traverse over all key/bucket pairs (start from `key`) in a bucket in sorted order.
`c` is a tx or a bucket.

### c:delbucket(key)

Deletes a bucket at the given `key`.
`c` is a tx or a bucket.

### bucket:get(key)

Retrieves the value for a `key` in the bucket.

### bucket:set(key, value)

Sets the `value` for the `key` in the bucket.
If the `key` exist then it's previous `value` will be overwritten.

### bucket:pairs([key])

Returns an iterator that can traverse over all key/value pairs (start from `key`) in a bucket in sorted order.
