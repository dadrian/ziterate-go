ZIterate, in Go
===============

This core ZMap iteration logic, reimplemented in Go. It is not as fast, but the
UintGroupIterator using the NextUint() function is pretty fast. You can pair
this with [zmap/go-iptree](https://github.com/zmap/go-iptree) to get the same
allowlisting functionality as in ZMap.

This is not a 1:1 replacement, the exact logic for allowlisting and blocklisting
in ZMap would need to be reimplemented, as well as the seed generation.

If I ever get this to feature complete, I'll move it to the ZMap organization.
