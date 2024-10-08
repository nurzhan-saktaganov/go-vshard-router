## Unreleased

BUG FIXES:

* RouterCallImpl: fix decoding response from storage_ref (partially #42)
* RouterCallImpl: fix decoding response from storage_map (partially #42)
* BucketDiscovery: check res for nil
* BucketStat: decode bsInfo by ptr
* ReplicaCall: fix decoding response (#42)
* ReplicaCall: fix ignoring timeout while waiting for future.Get()

FEATURES:

* Support new Sprintf-like logging interface (#48)
* DiscoveryTimeout by default is 1 minute (zero DiscoveryTimeout is not allowed #60)
* All discovering logs has new prefix [DISCOVERY]
* Introduce Replicaset.CallAsync, it is usefull to send concurrent requests to replicasets;
	additionally, CallAsync provides new interface to interact with replicaset without cons of interface of ReplicaCall

REFACTOR:

* resolve issue #38: simplify DiscoveryAllBuckets and remove suspicious if
* resolve issue #46: drastically simplify RouterMapCallRWImpl and added tests with real tnt
* Use typed nil pointers instead of memory allocation for EmptyMetrics and emptyLogger structs
* resolve issue #44: remove bucketCount field from struct Replicaset
* rename startCronDiscovery to cronDiscovery and make it panic-tolerant
* BucketStat: split into bucketStatAsync and bucketStatWait parts
* BucketDiscovery: do not spawn goroutines, just use futures in the single goroutine
* BucketResolve: make it alias for BucketDiscovery

TESTS:

* New test for RouterCallImpl (and fix the old one)
* New tnt tests for discovery logic
* New tnt tests for RouterMapCallRWImpl
* New tnt tests for topology logic
* Big CI update
  * 2 sections for CI: static checks and tests
  * integration tests run on ci with Tarantool cluster on vshard
  * implemented luacheck for static checks
* New tnt tests for ReplicaCall
* New tnt tests for CallAsync

EXAMPLES:
* customer go mod fixed 
* add customer example listen addr log

## 0.0.12

BUG FIXES:

* RouterCallImpl: fix decoding the response from vshard.storage.call
* RouterCallImpl: do not return nil error when StorageCallAssertError has happened
* BucketStat: always returns non-nil err, fixed
* DiscoveryAllBuckets returns nil even if errGr.Wait() returns err, fixed
* DiscoveryHandleBuckets: misusage of atomics, fixed
* race when accessing to idToReplicaset, fixed: idToReplicaset is immutable object now
* RouterMapCallRWImpl: fix misusage of refID atomic
* RouterMapCallRWImpl: decode bucketCount into 32 bit integer instead of 16 bit
* RouterMapCallRWImpl: fix concurrent access to idToResult map
* BucketDiscovery: fix possible concurrent access to resultRs and err vars
* RouterMapCallRWImpl: compare totalBucketCount against r.cfg.TotalBucketCount
* issue #39: fixed concurrent access to routeMap: use consistent view (immutable object) + atomics

FEATURES:

* Added etcd v2 topology provider implementation (#16)
* Add TopologyController mock for testing improve
* Add linter job (#33)
* New test framework with real tarantools

REFACTOR:

* Refactored docs (add QuickStart doc) and that library base on vhsard router
* Several linters are enabled because they are usefull
* Ignore .tmp files
* Refactored provider creation test caused by golang-ci lint (#33)
* Router implements directly TopologyController, no proxy object is used now
* Makefile refactored for racetests
* Tests coverage up 22% -> 33%


## 0.0.11

BUG FIXES:

* Fix buckets discovery (it doesn't freeze now)

FEATURES:

* BucketStat has become a public method (#21)
* Add golang-ci logic

REFACTOR:

* WSL lint providers fix
* Lint refactor with spaces
* Split tests for shadow and not vshard module
* Update Makefile with cover & lint
* Add more tests for providers
* TopologyController now is an interface
* Pool struct has been replaced with the Pooler interface, which has improved coverage (#21)

## 0.0.10

BUG FIXES:

* fix empty and multiple tnt procedure responses

FEATURES:

* start write CHANGELOG file
* use TopologyProvider instead direct config topology; if us uses old schema just change it to static provider
* add go report card

REFACTOR:

* refactored place for rs method
