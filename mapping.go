
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

import (
    "math"
)

// Interface for mapping element values to new values.
type FloatMapping interface {
    Eval(i, j int, val float64) float64
}

// Change matrix elements with  a Mapping. If flag bit UPPER is set then strictly lower
// part of the matrix is not touched. If flag bit LOWER is set then strictly upper
// part of the matrix is not touched. If bit SYMM is set then transformer is evaluated
// only for upper part indexes and strictly lower part is set symmetrically.
func (A *FloatMatrix) Map(t FloatMapping, bits ...int) {
    var flags int = NONE
    if len(bits) > 0 {
        flags = bits[0]
    }
    switch {
    case flags & UPPER != 0:
        // upper triangular/trapezoidial 
        for i := 0; i < A.rows; i++ {
            for j := i; j < A.cols; j++ {
                A.elems[i+j*A.step] = t.Eval(i, j, A.elems[i+j*A.step])
            }
        }
        return
    case flags & LOWER != 0:
        // lower triangular/trapezoidial
        for j := 0; j < A.cols; j++ {
            for i := j; i < A.rows; i++ {
                A.elems[i+j*A.step] = t.Eval(i, j, A.elems[i+j*A.step])
            }
        }
        return
    case flags & SYMM != 0:
        if A.rows != A.cols {
            return
        }
        for j := 0; j < A.cols; j++ {
            for i := 0; i < j; i++ {
                A.elems[i+j*A.step] = t.Eval(i, j, A.elems[i+j*A.step])
                A.elems[j+i*A.step] = A.elems[i+j*A.step]
            }
            A.elems[j+j*A.step] = t.Eval(j, j, A.elems[j+j*A.step])
        }
        return
    }
    // normal matrix here; access in memory order
    for j := 0; j < A.cols; j++ {
        for i := 0; i < A.rows; i++ {
            A.elems[i+j*A.step] = t.Eval(i, j, A.elems[i+j*A.step])
        }
    }
}

// Simple wrapper for element location dependent mapping.
type FloatEvaluator struct {
    Callable func(int, int, float64) float64
}

func (t *FloatEvaluator) Eval(i, j int, val float64) float64 {
    return t.Callable(i, j, val)
}

// Single parameter mapping from float value to another float value.
type FloatFunction struct {
    // float valued callable
    Callable func(float64) float64
}

// Evaluate callable with argument.
func (t *FloatFunction) Eval(i, j int, val float64) float64 {
    return t.Callable(val)
}

// Two parameter mapping from float value to new ie. newval = fn(oldval, const)
type Float2Function struct {
    // float valued callable function
    Callable func(float64, float64) float64
    // value of second argument for callable
    Constant float64
}

// Evaluate two parameter callable as Callable(val, Constant)
func (t *Float2Function) Eval(i, j int, val float64) float64 {
    return t.Callable(val, t.Constant)
}

// Element-wise logarithm.
func (A *FloatMatrix) Log() {
    A.Map(&FloatFunction{math.Log}, NONE)
}

// Element-wise scaling.
func (A *FloatMatrix) Scale(val float64) {
    fnc := func(a float64) float64 {
        return a*val
    }
    A.Map(&FloatFunction{fnc}, NONE)
}

// Element-wise adding
func (A *FloatMatrix) Add(val float64) {
    fnc := func(a float64) float64 {
        return a+val
    }
    A.Map(&FloatFunction{fnc}, NONE)
}

// Make matrix UPPER triangular matrix ie. set strictly lower part to zero.
func TriU(A *FloatMatrix, bits int) *FloatMatrix {
    fnc := func(i, j int, val float64) float64 {
        if i > j {
            return 0.0
        }
        if i == j && bits & UNIT != 0 {
            return 1.0
        }
        return val
    }
    // traverse lower part
    A.Map(&FloatEvaluator{fnc}, LOWER)
    return A
}

// Make matrix LOWER triangular matrix ie. set strictly upper part to zero.
func TriL(A *FloatMatrix, bits int) *FloatMatrix {
    fnc := func(i, j int, val float64) float64 {
        if j > i {
            return 0.0
        }
        if i == j && bits & UNIT != 0 {
            return 1.0
        }
        return val
    }
    // traverse upper part
    A.Map(&FloatEvaluator{fnc}, UPPER)
    return A
}

// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
