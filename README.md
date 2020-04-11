
# octree

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3aa076e74fce4e80af0e694116444410)](https://app.codacy.com/gh/The-Tensox/octree?utm_source=github.com&utm_medium=referral&utm_content=The-Tensox/octree&utm_campaign=Badge_Grade_Dashboard)
[![Build Status](https://img.shields.io/circleci/project/The-Tensox/octree/master.svg)](https://circleci.com/gh/The-Tensox/octree)

This is a work in progress, API may change a little bit and current implementations may not match the ideal complexity shown in the papers, some functions have just currently a naive non-optimal implementation.

## Installation

```bash
go get -u github.com/The-Tensox/octree
```

## Usage

```go
o := NewOctree(protometry.NewBoxMinMax(1, 1, 1, 4, 4, 4))
myObj := NewObjectCube(0, 2, 2, 3, 0.5)
ok := o.Insert(*myObj)
colliders := o.GetColliding(*protometry.NewBoxMinMax(0, 0, 0, 0.9, 2.9, 0.9))
myObj = o.Move(*myObj, 0, 0, 0, 2, 2, 2) // Using bounds
myObj = o.Move(*myObj, 3, 3, 3) // Using position, assume cube of side 1
```

## Test

```bash
go test
```

## Benchmark

```bash
# XXX will skip tests
go test -benchmem -run XXX -bench . -benchtime 0.2s
```

THEN

|Name   |   Runs   |   time   |   Bytes   |   Allocs   |
|:-----:|:--------:|:--------:|:---------:|:----------:|
|BenchmarkNode_Insert-8   |   43530   |   5192 ns/op   |   1658 B/op  |   23 allocs/op   |
|BenchmarkNode_GetColliding-8   |   10000   |   21816 ns/op   |   8748 B/op   |   47 allocs/op   |
|BenchmarkNode_Remove-8   |   25642   |   9252 ns/op   |   2344 B/op   |   21 allocs/op   |
|BenchmarkNode_Move-8   |   13378   |   17764 ns/op   |   3791 B/op   |   45 allocs/op   |

 ---
NOW

|Name   |   Runs   |   time   |   Bytes   |   Allocs   |
|:-----:|:--------:|:--------:|:---------:|:----------:|
|BenchmarkNode_Insert-8   |   15242   |   74134 ns/op   |   28540 B/op  |   920 allocs/op   |
|BenchmarkNode_GetColliding-8   |   5514   |   360566 ns/op   |   128144 B/op   |   4098 allocs/op   |
|BenchmarkNode_Remove-8   |   10000   |   146253 ns/op   |   52386 B/op   |   1678 allocs/op   |
|BenchmarkNode_Move-8   |   7546   |   235199 ns/op   |   88335 B/op   |   2870 allocs/op   |

## Roadmap

- [ ] Improve performance, more complexity checks / benchmarks
- [ ] Tree vizualisation ?

## References

- [Github storpipfugl/pykdtree: Fast kd-tree implementation in Python](https://github.com/storpipfugl/pykdtree)
- [Github Rust rust-octree](https://github.com/ybyygu/rust-octree)
- [Github JS sparse-octree](https://github.com/vanruesc/sparse-octree)
- [Github Distributed adaptive octree construction, 2:1 balancing & partitioning based on space filling curves](https://github.com/paralab/Dendro-5.01)
- [Github UnityOctree](https://github.com/Nition/UnityOctree)
- [Book: Real-Time Collision Detection - Christer Ericson](https://www.amazon.com/exec/obidos/tg/detail/-/1558607323?tag=realtimecolli-20)
- [AN OVERVIEW OF QUADTREES, OCTREES, AND RELATED HIERARCHICAL DATA STRUCTURES](https://www.cs.umd.edu/~hjs/pubs/Samettfcgc88-ocr.pdf)
- [Efficient Sparse Voxel Octrees](https://research.nvidia.com/publication/efficient-sparse-voxel-octrees)
- [An Efficient Parametric Algorithm for Octree Traversal](http://wscg.zcu.cz/wscg2000/Papers_2000/X31.pdf)
- [O-CNN: Octree-based Convolutional Neural Networks for 3D ShapeAnalysis](https://wang-ps.github.io/O-CNN_files/CNN3D.pdf)
- Behley, J.; Steinhage, V.; Cremers, A. B. Efficient Radius Neighbor Search in
    Three-Dimensional Point Clouds. In 2015 IEEE International Conference on
    Robotics and Automation (ICRA); 2015; pp 3625–3630.
- [scipy.spatial.cKDTree — SciPy Reference Guide](https://docs.scipy.org/doc/scipy/reference/generated/scipy.spatial.cKDTree.html)
- Milinda Fernando, David Neilsen, Hyun Lim, Eric Hirschmann, Hari Sundar, ”Massively Parallel Simulations of Binary Black Hole Intermediate-Mass-Ratio Inspirals” SIAM Journal on Scientific Computing 2019. [Paper](https://doi.org/10.1137/18M1196972)
- Milinda Fernando, David Neilsen, Hari Sundar, ”A scalable framework for Adaptive Computational General Relativity on Heterogeneous Clusters”, (ACM International Conference on Supercomputing, ICS’19)
- Milinda Fernando, Dmitry Duplyakin, and Hari Sundar. 2017. ”Machine and Application Aware Partitioning for Adaptive Mesh Refinement Applications”. In Proceedings of the 26th International Symposium on High-Performance Parallel and Distributed Computing (HPDC ’17). ACM, New York, NY, USA, 231-242. [Paper](https://doi.org/10.1145/3078597.3078610)