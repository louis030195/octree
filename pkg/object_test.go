package octree

import (
    "github.com/louis030195/protometry/api/volume"
    "testing"
)

func TestObject_Equal(t *testing.T) {
	type fields struct {
		id     uint64
		Data   interface{}
		Bounds volume.Box
	}
	type args struct {
		object Object
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				id:     1, // Because NewObject increments the id automatically
				Data:   nil,
				Bounds: volume.Box{},
			},
			args: args{object: *NewObject(nil, volume.Box{})},
			want: true,
		},
		{
			fields: fields{
				id:     2, // Because NewObject increments the id automatically
				Data:   2728624,
				Bounds: volume.Box{},
			},
			args: args{object: *NewObject(nil, volume.Box{})},
			want: true,
		},
		{ // Equality is only checked on id, not data or bounds
			fields: fields{
				id:     2424224,
				Data:   1,
				Bounds: volume.Box{},
			},
			args: args{object: *NewObject(1, volume.Box{})},
			want: false,
		},
		{ // Equality is only checked on id, not data or bounds
			fields: fields{
				id:     4,
				Data:   nil,
				Bounds: *volume.NewBoxOfSize(0, 27.332, 0, 1),
			},
			args: args{object: *NewObject(nil, *volume.NewBoxOfSize(8726.1, 0, 0, 1))},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Object{
				id:     tt.fields.id,
				Data:   tt.fields.Data,
				Bounds: tt.fields.Bounds,
			}
			if got := o.Equal(tt.args.object); got != tt.want {
				t.Errorf("Equal() = %v, want %v, objects: %v, %v", got, tt.want, o, tt.args.object)
			}
		})
	}
}
