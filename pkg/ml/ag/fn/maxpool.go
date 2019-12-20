// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"brillion.io/spago/pkg/mat"
	"brillion.io/spago/pkg/utils"
	"math"
)

type MaxPool struct {
	x    Operand
	rows int
	cols int
	// initialized during the forward pass
	y       mat.Matrix
	argmaxi [][]int
	argmaxj [][]int
}

func NewMaxPool(x Operand, r, c int) *MaxPool {
	return &MaxPool{
		x:       x,
		rows:    r,
		cols:    c,
		y:       nil,
		argmaxi: nil,
		argmaxj: nil,
	}
}

// Forward computes the output of the function.
func (r *MaxPool) Forward() mat.Matrix {
	if !(r.x.Value().Rows()%r.rows == 0 && r.x.Value().Columns()%r.cols == 0) {
		panic("fn: size mismatch")
	}

	r.y = mat.NewEmptyDense(r.x.Value().Rows()/r.rows, r.x.Value().Columns()/r.cols)
	r.argmaxi = utils.MakeIntMatrix(r.y.Dims()) // output argmax row index
	r.argmaxj = utils.MakeIntMatrix(r.y.Dims()) // output argmax column index

	for row := 0; row < r.y.Rows(); row++ {
		for col := 0; col < r.y.Columns(); col++ {
			max := math.SmallestNonzeroFloat64
			for i := row * r.rows; i < (row*r.rows)+r.rows; i++ {
				for j := col * r.cols; j < (col*r.cols)+r.rows; j++ {
					val := r.x.Value().At(i, j)
					if val > max {
						max = val
						r.argmaxi[row][col] = i
						r.argmaxj[row][col] = j
					}
				}
			}
			r.y.Set(max, row, col)
		}
	}

	return r.y
}

func (r *MaxPool) Backward(gy mat.Matrix) {
	if r.x.RequiresGrad() {
		gx := r.x.Value().ZerosLike()
		for row := 0; row < r.y.Rows(); row++ {
			for col := 0; col < r.y.Columns(); col++ {
				gx.Set(gy.At(row, col), r.argmaxi[row][col], r.argmaxj[row][col])
			}
		}
		r.x.PropagateGrad(gx)
	}
}
