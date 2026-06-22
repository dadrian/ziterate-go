ZIterate, in Go
===============

This core ZMap iteration logic, reimplemented in Go. It is not as fast, but the
UintGroupIterator using the NextUint() function is pretty fast. You can pair
this with [zmap/go-iptree](https://github.com/zmap/go-iptree) to get the same
allowlisting functionality as in ZMap.

This is not a 1:1 replacement, the exact logic for allowlisting and blocklisting
in ZMap would need to be reimplemented, as well as the seed generation.

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
