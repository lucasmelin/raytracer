package geometry

import (
	"fmt"
	"reflect"
	"testing"
)

const epsilon = 0.00000001

type RndMock struct {
	floats []float64
	idx    int
}

func (rnd *RndMock) Float64() float64 {
	idx := rnd.idx
	rnd.idx++
	return rnd.floats[idx]
}

func unitEquals(a Unit, b Unit) bool {
	if (a.X-b.X) >= epsilon || (b.X-a.X) >= epsilon {
		return false
	}
	if (a.Y-b.Y) >= epsilon || (b.Y-a.Y) >= epsilon {
		return false
	}
	if (a.Z-b.Z) >= epsilon || (b.Z-a.Z) >= epsilon {
		return false
	}
	return true
}

func TestRandUnit(t *testing.T) {
	tests := []struct {
		name string
		rnd  RndMock
		want Unit
	}{
		{
			name: "basic",
			rnd: RndMock{
				floats: []float64{0.8, 0.7, 0.6},
			},
			want: NewVec(2.0*0.8-1.0, 2.0*0.7-1.0, 2.0*0.6-1.0).ToUnit(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandUnit(&tt.rnd); !unitEquals(got, tt.want) {
				t.Errorf("RandUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Scale(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name     string
		original fields
		factor   float64
		want     Vec
	}{
		{
			name: "no scaling",
			original: fields{
				X: 1.0,
				Y: 2.0,
				Z: 3.0,
			},
			factor: 1.0,
			want:   NewVec(1.0, 2.0, 3.0),
		},
		{
			name: "double",
			original: fields{
				X: 1.0,
				Y: 2.0,
				Z: 3.0,
			},
			factor: 2.0,
			want:   NewVec(2.0, 4.0, 6.0),
		},
		{
			name: "half",
			original: fields{
				X: 1.0,
				Y: 2.0,
				Z: 3.0,
			},
			factor: 0.5,
			want:   NewVec(0.5, 1.0, 1.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.original.X,
				Y: tt.original.Y,
				Z: tt.original.Z,
			}
			if got := v.Scale(tt.factor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Zero(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}

	type testCase struct {
		name   string
		fields fields
		want   bool
	}

	var tests []testCase
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				want := false
				name := fmt.Sprintf("non-zero vector {%.2f %.2f %.2f}", i, j, k)
				if i == 0 && j == 0 && k == 0 {
					want = true
					name = "zero vector"
				}
				tests = append(tests, testCase{
					name: name,
					fields: fields{
						X: i,
						Y: j,
						Z: k,
					},
					want: want,
				})
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Zero(); got != tt.want {
				t.Errorf("Zero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Min(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}

	tests := []struct {
		name   string
		fields fields
		v2     Vec
		want   Vec
	}{
		{
			name: "first vector smaller",
			fields: fields{
				X: 1,
				Y: 2,
				Z: 3,
			},
			v2:   NewVec(10, 10, 10),
			want: NewVec(1, 2, 3),
		},
		{
			name: "second vector smaller",
			fields: fields{
				X: 10,
				Y: 10,
				Z: 10,
			},
			v2:   NewVec(1, 2, 3),
			want: NewVec(1, 2, 3),
		},
		{
			name: "second vector X smaller",
			fields: fields{
				X: 10,
				Y: 2,
				Z: 3,
			},
			v2:   NewVec(1, 10, 10),
			want: NewVec(1, 2, 3),
		},
		{
			name: "second vector Y smaller",
			fields: fields{
				X: 1,
				Y: 10,
				Z: 3,
			},
			v2:   NewVec(10, 2, 10),
			want: NewVec(1, 2, 3),
		},
		{
			name: "second vector Z smaller",
			fields: fields{
				X: 1,
				Y: 2,
				Z: 10,
			},
			v2:   NewVec(10, 10, 3),
			want: NewVec(1, 2, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Min(tt.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Inv(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		want   Vec
	}{
		{
			name: "positive vector",
			fields: fields{
				X: 1,
				Y: 2,
				Z: 3,
			},
			want: NewVec(-1, -2, -3),
		},
		{
			name: "negative vector",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			want: NewVec(1, 2, 3),
		},
		{
			name: "mixed sign vector",
			fields: fields{
				X: -1,
				Y: 2,
				Z: -3,
			},
			want: NewVec(1, -2, 3),
		},
		{
			name: "zero vector",
			fields: fields{
				X: 0,
				Y: 0,
				Z: 0,
			},
			want: NewVec(0, 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Inv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Add(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		v2     Vec
		want   Vec
	}{
		{
			name: "two positive vectors",
			fields: fields{
				X: 1,
				Y: 2,
				Z: 3,
			},
			v2:   NewVec(4, 5, 6),
			want: NewVec(5, 7, 9),
		},
		{
			name: "two negative vectors",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(-4, -5, -6),
			want: NewVec(-5, -7, -9),
		},
		{
			name: "positive and negative vector",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(4, 5, 6),
			want: NewVec(3, 3, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Add(tt.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Sub(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		v2     Vec
		want   Vec
	}{
		{
			name: "two positive vectors",
			fields: fields{
				X: 4,
				Y: 5,
				Z: 6,
			},
			v2:   NewVec(1, 2, 3),
			want: NewVec(3, 3, 3),
		},
		{
			name: "two negative vectors",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(-4, -5, -6),
			want: NewVec(3, 3, 3),
		},
		{
			name: "positive and negative vector",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(4, 5, 6),
			want: NewVec(-5, -7, -9),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Sub(tt.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_ToUnit(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		want   Unit
	}{
		{
			name: "3^2 + 4^2 + 12^2 = 13^2",
			fields: fields{
				X: 3,
				Y: 4,
				Z: 12,
			},
			want: Unit{NewVec(float64(3.0)/13, float64(4.0)/13, float64(12.0)/13)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.ToUnit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Max(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		v2     Vec
		want   Vec
	}{
		{
			name: "first vector larger",
			fields: fields{
				X: 10,
				Y: 20,
				Z: 30,
			},
			v2:   NewVec(1, 1, 1),
			want: NewVec(10, 20, 30),
		},
		{
			name: "second vector larger",
			fields: fields{
				X: 1,
				Y: 1,
				Z: 1,
			},
			v2:   NewVec(10, 20, 30),
			want: NewVec(10, 20, 30),
		},
		{
			name: "second vector X larger",
			fields: fields{
				X: 1,
				Y: 20,
				Z: 30,
			},
			v2:   NewVec(10, 2, 3),
			want: NewVec(10, 20, 30),
		},
		{
			name: "second vector Y larger",
			fields: fields{
				X: 10,
				Y: 2,
				Z: 30,
			},
			v2:   NewVec(1, 20, 3),
			want: NewVec(10, 20, 30),
		},
		{
			name: "second vector Z larger",
			fields: fields{
				X: 10,
				Y: 20,
				Z: 3,
			},
			v2:   NewVec(1, 2, 30),
			want: NewVec(10, 20, 30),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Max(tt.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec_Mul(t *testing.T) {
	type fields struct {
		X float64
		Y float64
		Z float64
	}
	tests := []struct {
		name   string
		fields fields
		v2     Vec
		want   Vec
	}{
		{
			name: "two unit vectors",
			fields: fields{
				X: 1,
				Y: 1,
				Z: 1,
			},
			v2:   NewVec(1, 1, 1),
			want: NewVec(1, 1, 1),
		},
		{
			name: "two positive vectors",
			fields: fields{
				X: 1,
				Y: 2,
				Z: 3,
			},
			v2:   NewVec(4, 5, 6),
			want: NewVec(4, 10, 18),
		},
		{
			name: "two negative vectors",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(-4, -5, -6),
			want: NewVec(4, 10, 18),
		},
		{
			name: "one negative one positive vector",
			fields: fields{
				X: -1,
				Y: -2,
				Z: -3,
			},
			v2:   NewVec(4, 5, 6),
			want: NewVec(-4, -10, -18),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vec{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Mul(tt.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}
