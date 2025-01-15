# protocompiletest

## Run

Repeat

```sh
$ go test -v
```

Followed by cache clean

```sh
$ go clean -cache -testcache
```

### Example

```sh
❯ go test -v
=== RUN   TestPlainProto
!!! ACTUAL JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
!!! EXPECTED JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
--- PASS: TestPlainProto (0.00s)
PASS
ok  	github.com/bojand/protocompiletest	0.266s

❯ go clean -cache -testcache

❯ go test -v
=== RUN   TestPlainProto
!!! ACTUAL JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
!!! EXPECTED JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
--- PASS: TestPlainProto (0.00s)
PASS
ok  	github.com/bojand/protocompiletest	0.270s

❯ go clean -cache -testcache

❯ go test -v
=== RUN   TestPlainProto
!!! ACTUAL JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
    main_test.go:111:
        	Error Trace:	/Users/bojandjurkovic/dev/protocompiletest/main_test.go:111
        	Error:      	Not equal:
        	            	expected: []byte{0x8, 0x1, 0x12, 0x3, 0x34, 0x34, 0x34, 0x1a, 0x6, 0x8, 0xa0, 0xdc, 0xc9, 0xa5, 0x6, 0x22, 0x6, 0x8, 0xb0, 0xf8, 0xc9, 0xa5, 0x6}
        	            	actual  : []byte{0x22, 0x6, 0x8, 0xb0, 0xf8, 0xc9, 0xa5, 0x6, 0x8, 0x1, 0x12, 0x3, 0x34, 0x34, 0x34, 0x1a, 0x6, 0x8, 0xa0, 0xdc, 0xc9, 0xa5, 0x6}

        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,4 +1,4 @@
        	            	 ([]uint8) (len=23) {
        	            	- 00000000  08 01 12 03 34 34 34 1a  06 08 a0 dc c9 a5 06 22  |....444........"|
        	            	- 00000010  06 08 b0 f8 c9 a5 06                              |.......|
        	            	+ 00000000  22 06 08 b0 f8 c9 a5 06  08 01 12 03 34 34 34 1a  |"...........444.|
        	            	+ 00000010  06 08 a0 dc c9 a5 06                              |.......|
        	            	 }
        	Test:       	TestPlainProto
!!! EXPECTED JSON DATA:
{"version":1, "id":"444", "createdAt":"2023-07-15T10:00:00Z", "lastUpdatedAt":"2023-07-15T11:00:00Z", "orderValue":0, "lineItems":[], "revision":0}
--- FAIL: TestPlainProto (0.00s)
FAIL
exit status 1
FAIL	github.com/bojand/protocompiletest	0.272s
```

## Proto

```sh
buf generate --template=testdata/proto/buf.gen.yaml testdata/proto
```
