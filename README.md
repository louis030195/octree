
# octree-go

## Test

```bash
go test -v .
```

## Benchmark

```bash
go test -benchmem -run . -bench . -benchtime 0.2s
```

> BenchmarkOctreeNode_Insert-8         231           1656908 ns/op          944112 B/op      20360 allocs/op

>BenchmarkOctreeNode_Search-8         800           3247648 ns/op         1188731 B/op      39624 allocs/op

>PASS
ok      github.com/The-Tensox/octree    9.121s
