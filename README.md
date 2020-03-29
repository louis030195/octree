
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
// Create an Octree of region 1;1;1 to 4;4;4
o := NewOctree(protometry.NewBox(1, 1, 1, 4, 4, 4))
// Out of bounds
err := o.Insert(NewPoint(10, 10, 10, 1))
// Insert 2 at 3;3;3
err = o.Insert(NewPoint(3, 3, 3, 2))
// Insert 3 at 3;3;4
err = o.Insert(NewPoint(3, 3, 4, 3))
// Insert 4 at 3;4;4
err = o.Insert(NewPoint(3, 4, 4, 4))
// 3 points
points := o.Get(1, 1, 1, 4, 4, 4)
// Point 3;3;3 of value 2
point := o.Get(3, 3, 3)

o = NewOctree(protometry.NewBox(0, 0, 0, 10, 10, 10))
p1 := protometry.NewBox(0, 0, 0, 0, 1, 0) // Line collider
o.Insert(*NewPointCollide(1, *p1, *p1.GetCenter()))
p2 := protometry.NewBox(0, 2, 0, 0, 3, 0)
o.Insert(*NewPointCollide(2, *p2, *p2.GetCenter()))
p3 := protometry.NewBox(0, 4, 0, 0, 5, 0)
o.Insert(*NewPointCollide(3, *p3, *p3.GetCenter()))

// Cast a ray toward up from 0;0;0 of length 10
points := *o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 10) // 3 points
```

## Test

```bash
go test -v .
```

## Benchmark

```bash
# XXX will skip tests
go test -benchmem -run XXX -bench . -benchtime 0.2s
```

## Roadmap

- [ ] Improve performance, more complexity checks / benchmarks
- [ ] Take decision either to use ...float64 or protometry.VectorN in args (since there is no overloading in Go ;))
- [ ] point.OnCollisionEnter(myCallback), point.OnCollisionStay(myCallback2), point.OnCollisionExit(myCallback3) ...
- [ ] Trigger collision callback on move, insert, remove
- [ ] Implement "cubecast", fix raycast
- [ ] Implement shrink, merge
- [ ] Tree vizualisation ?

## References

- [Github storpipfugl/pykdtree: Fast kd-tree implementation in Python](https://github.com/storpipfugl/pykdtree)
- [Github Rust rust-octree](https://github.com/ybyygu/rust-octree)
- [Github JS sparse-octree](https://github.com/vanruesc/sparse-octree)
- [Github Distributed adaptive octree construction, 2:1 balancing & partitioning based on space filling curves](https://github.com/paralab/Dendro-5.01)
- [Github UnityOctree](https://github.com/Nition/UnityOctree)
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