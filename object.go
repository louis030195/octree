package octree

import "github.com/The-Tensox/protometry"

// Object stores data and collider about the object
type Object struct {
	data   interface{}
	bounds protometry.Box
}

// NewObject is a Object constructor with collider for ease of use
func NewObject(data interface{}, bounds protometry.Box) *Object {
	return &Object{data: data, bounds: bounds}
}

func NewObjectCube(data interface{}, x, y, z, size float64) *Object {
	return NewObject(data, *protometry.NewBoxOfSize(*protometry.NewVectorN(x, y, z), size))
}

func (o *Object) Equal(object Object) bool {
	return o.data == object.data && o.bounds.Equal(object.bounds)

}
