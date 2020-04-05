package octree

import (
	"fmt"
	"github.com/The-Tensox/protometry"
)

// Object stores data and bounds
type Object struct {
	Data   interface{}
	Bounds protometry.Box
}

// NewObject is a Object constructor with bounds for ease of use
func NewObject(data interface{}, bounds protometry.Box) *Object {
	return &Object{Data: data, Bounds: bounds}
}

func NewObjectCube(data interface{}, x, y, z, size float64) *Object {
	return NewObject(data, *protometry.NewBoxOfSize(*protometry.NewVectorN(x, y, z), size))
}

func (o *Object) Equal(object Object) bool {
	return o.Data == object.Data && o.Bounds.Equal(object.Bounds)

}

func (o *Object) ToString() string {
	return fmt.Sprintf("Data:%v\nBounds:{\n%v\n}", o.Data, o.Bounds.ToString())
}