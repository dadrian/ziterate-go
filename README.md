ZIterate, in Go
===============

This is core ZMap-style iteration logic, reimplemented in Go. It is not as fast,
but the UintGroupIterator using the NextUint() function is pretty fast.

This is not a 1:1 replacement. It includes allowlist, blocklist, and sharding
support similar to ziterate CLI, but does not try to exactly match ZMap's
address ordering or blocklist semantics. Using the same seed and flags between
the ziterate in the ZMap repository, and the Go implementation, will **not**
result in the same outputs.

Usage
----

Install or run the CLI with Go:

```sh
go install github.com/zmap/ziterate/cmd/ziterate@latest
ziterate --help
```

Iterate over the public IPv4 space:

```sh
ziterate
```

Iterate over one or more CIDR allowlist entries:

```sh
ziterate 10.0.0.0/24 192.0.2.10
```

Emit IP and port pairs:

```sh
ziterate --target-ports 80,443,8000-8002 10.0.0.0/24
```

Use allowlist and blocklist files:

```sh
ziterate --allowlist-file allow.txt --blocklist-file block.txt
```

Generate repeatable output with a seed and cap the number of targets:

```sh
ziterate --seed 12345 --max-targets 1000 10.0.0.0/16
```

Split output across shards. Sharding requires a seed so each shard uses the same
base ordering:

```sh
ziterate --seed 12345 --shards 4 --shard 0 10.0.0.0/16
ziterate --seed 12345 --shards 4 --shard 1 10.0.0.0/16
```

Examples
--------

Pick the smallest built-in ZMap group that can iterate over a target set:

```go
package main

import (
	"fmt"

	"github.com/zmap/ziterate"
)

func main() {
	group, err := ziterate.SmallestZMapGroupFor(1 << 32)
	if err != nil {
		panic(err)
	}
	fmt.Println(group.P)
}
```

Use the fast uint64 iterator when the selected group is small enough:

```go
package main

import (
	"crypto/rand"
	"fmt"

	"github.com/zmap/ziterate"
)

func main() {
	group, err := ziterate.SmallestZMapGroupFor(1 << 32)
	if err != nil {
		panic(err)
	}
	it, err := ziterate.UintGroupIteratorFromGroup(group, rand.Reader)
	if err != nil {
		panic(err)
	}
	for x := it.NextUint(); x != 0; x = it.NextUint() {
		fmt.Println(x)
	}
}
```

Use the big.Int iterator for larger groups:

```go
package main

import (
	"crypto/rand"
	"fmt"

	"github.com/zmap/ziterate"
)

func main() {
	group := ziterate.ZMapGroups[len(ziterate.ZMapGroups)-1]
	it, err := ziterate.BigIntGroupIteratorFromGroup(group, rand.Reader)
	if err != nil {
		panic(err)
	}
	for x := it.NextBigInt(); x != nil; x = it.NextBigInt() {
		fmt.Println(x)
	}
}
```
