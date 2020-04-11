// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
)

type View struct {
	x  Operand
	sx int
	sy int
	lx int // x length
	ly int // y length
}

func NewView(x Operand, sx, sy, lx, ly int) *View {
	return &View{x: x, sx: sx, sy: sy, lx: lx, ly: ly}
}

// Forward computes the output of the function.
func (r *View) Forward() mat.Matrix {
	y := mat.NewEmptyDense(r.lx, r.ly)
	for i := 0; i < r.lx; i++ {
		for j := 0; j < r.ly; j++ {
			y.Set(i, j, r.x.Value().At(i+r.sx, j+r.sy))
		}
	}
	return y
}

func (r *View) Backward(gy mat.Matrix) {
	if r.x.RequiresGrad() {
		gx := mat.NewEmptyDense(r.x.Value().Dims())
		for i := 0; i < r.lx; i++ {
			for j := 0; j < r.ly; j++ {
				gx.Set(i+r.sx, j+r.sy, gy.At(i, j))
			}
		}
		r.x.PropagateGrad(gx)
	}
}
