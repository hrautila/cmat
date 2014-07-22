
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

import (
	"math/rand"
	"time"
)

// Interface for providing float values.
type FloatSource interface {
    Get(i, j int) float64
}

// Interface for pushing out elements.
type FloatSink interface {
    Put(i, j int, val float64)
}

// Source that produces const values.
type FloatConstSource struct {
	Const float64
}

func NewFloatConstSource(val float64) *FloatConstSource {
    return &FloatConstSource{val}
}

// Get 
func (s *FloatConstSource) Get(i, j int) float64 {
	return s.Const;
}

// Float value source for normally distributed values with mean `Mean` and
// standard deviation `StdDev`.
type FloatNormSource struct {
    // mean of the distribution
    Mean float64
    // required standard deviation
    StdDev float64
    Rnd  *rand.Rand
}

// Create a new source of Normally distributed float values. Optional parameters
// is either StdDev or tuple of (StdDev, Mean). Initial random number is multiplied with
// StdDev and Mean is added to the result. Default values for (StdDev,Mean) is (1.0, 0.0)
func NewFloatNormSource(params...float64) *FloatNormSource {
    var mean, stddev float64 = 0.0, 1.0
    switch len(params) {
    case 1:
        stddev = params[0]
        mean = 0.0
    case 2:
        stddev = params[0]
        mean = params[1]
    }
    return &FloatNormSource{mean, stddev, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Fetch a new value from distribution.
func (s *FloatNormSource) Get(i, j int) float64 {
    if s.Rnd == nil {
        s.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
    }
    return s.Rnd.NormFloat64()*s.StdDev + s.Mean
}

// Float value source for uniformly distributed in range [0.0,1.0) shifted
// and scaled and returning value Scale*(Uniform()-Shift)
type FloatUniformSource struct {
    // shifting factor
    Shift float64
    // scaling factor
    Scale float64
    Rnd  *rand.Rand
}

// Create a new source of uniformly distributed float values. Optional parameters
// is either shift or tuple of (shift, scale). Value of shift is added to 
// initial value which is then multiplied with the scaling factor. Default values
// for shift and scale are (0.0, 1.0) that creates source for random numbers from [0.0, 1.0)
func NewFloatUniformSource(params... float64) *FloatUniformSource {
    var shift, scale float64 = 0.0, 1.0
    switch len(params) {
    case 1:
        scale = params[0]
        shift = 0.0
    case 2:
        scale = params[0]
        shift = params[1]
    }
    return &FloatUniformSource{shift, scale, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Fetch a new value from distribution.
func (s *FloatUniformSource) Get(i, j int) float64 {
    if s.Rnd == nil {
        s.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
    }
    return s.Scale * (s.Rnd.Float64() + s.Shift)
}


// Source for retrieving matrix elements from a table. Default is returned
// if element index outside table.
type FloatTableSource struct {
    Data [][]float64
    Default  float64
}

func (ts *FloatTableSource) Get(i, j int) float64 {
    if i >= len(ts.Data) {
        return ts.Default
    }
    if j >= len(ts.Data[i]) {
        return ts.Default
    }
    return ts.Data[i][j]
}

// Return source table dimensions.
func (ts *FloatTableSource) Size() (int, int) {
    cols := 0
    rows := len(ts.Data)
    for i := 0; i < rows; i++ {
        if cols < len(ts.Data[i]) {
            cols = len(ts.Data[i])
        }
    }
    return rows, cols
}

// Create a new table source.
func NewFloatTableSource(data [][]float64, defval float64) *FloatTableSource {
    return &FloatTableSource{data, defval}
}


// Set matrix elements from source. Optional bits define which part of
// the matrix is accessed. Default is to set all entries. 
// 
// Bits
//   UPPER        set upper triangular/trapezoidal part
//   UPPER|UNIT   set strictly upper triangular/trapezoidal part 
//   LOWER        set lower triangular/trapezoidal part
//   LOWER|UNIT   set strictly lower triangular/trapezoidal part 
//   SYMM         set lower part symmetrically to upper part
//
// To set strictly lower part of a matrix: A.SetFrom(src, LOWER|UNIT)
//
func (m *FloatMatrix) SetFrom(source FloatSource, bits ...int) {
    var flags int = NONE
    if len(bits) > 0 {
        flags = bits[0]
    }
    unit := 0
    if flags & UNIT != 0 {
        unit = 1
    }
    switch {
    case flags & UPPER != 0:
        // upper triangular/trapezoidial, by rows
        for i := 0; i < m.rows; i++ {
            for j := i+unit; j < m.cols; j++ {
                m.elems[i+j*m.step] = source.Get(i, j)
            }
        }
        return
    case flags & LOWER != 0:
        // lower triangular/trapezoidial, by columns
        for j := 0; j < m.cols; j++ {
            for i := j+unit; i < m.rows; i++ {
                m.elems[i+j*m.step] = source.Get(i, j)
            }
        }
        return
    case flags & SYMM != 0:
        if m.rows != m.cols {
            return
        }
        for j := 0; j < m.cols; j++ {
            for i := 0; i < j; i++ {
                m.elems[i+j*m.step] = source.Get(i, j)
                m.elems[j+i*m.step] = m.elems[i+j*m.step]
            }
            m.elems[j+j*m.step] = source.Get(j, j)
        }
        return
    }
    // normal matrix here
    for j := 0; j < m.cols; j++ {
        for i := 0; i < m.rows; i++ {
            m.elems[i+j*m.step] = source.Get(i, j)
        }
    }
}

// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
