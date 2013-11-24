
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

import (
    "math"
    "fmt"
)

func indexMin(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Column majoe double precision matrix.
type FloatMatrix struct {
    elems []float64
    step int
    rows int
    cols int
}

type FlagBits int

const (
    UPPER = 1 << iota
    LOWER
    SYMM
    UNIT
    NONE = 0
)

// Make new matrix of size r rows, c cols.
func NewMatrix(r, s int) *FloatMatrix {
    ebuf := make([]float64, r*s, r*s)
    return &FloatMatrix{ebuf, r, r, s}
}

// Make a new copy of matrix
func NewCopy(A *FloatMatrix) *FloatMatrix {
    B := NewMatrix(A.Size())
    B.Copy(A)
    return B
}

// Make a new matrix and use ebuf as element storage. cap(ebuf) must not be less than
// rows*cols.
func MakeMatrix(rows, cols int, ebuf []float64) *FloatMatrix {
    if int(cap(ebuf)) < rows*cols {
        return nil;
    }
    return &FloatMatrix{ebuf, rows, rows, cols}
}

// Set matrix size and storage. Minimum size for ebuf is stride*cols.
// If stride zero or negative then rows is used as row stride.
// Returns nil if buffer capasity too small. Otherwise returns A. 
func (A *FloatMatrix) SetBuf(rows, cols, stride int, ebuf []float64) *FloatMatrix {
    if stride <= 0 {
        stride = rows
    }
    if int(cap(ebuf)) < stride*cols {
        return nil;
    }
    A.elems = ebuf
    A.rows = rows
    A.cols = cols
    A.step = stride
    return A
}

// Get size of the matrix as tuple (rows, cols).
func (A *FloatMatrix) Size() (int, int) {
    return A.rows, A.cols
}

// Get row stride of the matrix.
func (A *FloatMatrix) Stride() int {
    return A.step
}

// Get number of elements in matrix.
func (A *FloatMatrix) Len() int {
    return A.rows*A.cols;
}

func (A *FloatMatrix) IsVector() bool {
    return A != nil && (A.rows == 1 || A.cols == 1)
}

// Return raw element array.
func (A *FloatMatrix) Data() []float64 {
    return A.elems
}

// Make A submatrix of B.  Returns A.
func (A *FloatMatrix) SubMatrix(B *FloatMatrix, row, col int, sizes ...int) *FloatMatrix {
    var nr, nc, step int
    if row < 0 {
        row += B.rows
    }
    if col < 0 {
        col += B.cols
    }
    nr = B.rows - row
    nc = B.cols - col
    step = B.step
    switch len(sizes) {
    case 2:
        nr = sizes[0]
        nc = sizes[1]
        step = B.step
    case 3:
        nr = sizes[0]
        nc = sizes[1]
        step = sizes[2]
    }
    A.step = step
    A.rows = nr
    A.cols = nc
    if row >= 0 && row < B.rows && col >= 0 && col < B.cols {
        A.elems = B.elems[row+col*B.step:]
    } else {
        A.elems = nil
        A.rows = 0
        A.cols = 0
    }
    return A
}

// Make R a row vector of A i.e. R = A[row,:]
func (R *FloatMatrix) Row(A *FloatMatrix, row int, sizes ...int) *FloatMatrix {
    if row >= A.rows {
        return nil
    }
    if row < 0 {
        row += A.rows
    }
    var col int = 0
    var nc int = A.cols
    if len(sizes) == 1 {
        col = sizes[0]
        nc = A.cols - col
    } else if len(sizes) == 2 {
        col = sizes[0]
        nc = sizes[1]
    }
    if col + nc > A.cols {
        return nil
    }
    R.step = A.step
    R.rows = 1
    R.cols = nc
    if row >= 0 && row < A.rows && col < A.cols {
        R.elems = A.elems[row+col*A.step:]
    } else {
        R.elems = nil
        R.rows = 0
        R.cols = 0
    }
    return R
}

// Make C column of A. C = A[:,col]. Parameter sizes is singleton (row) and column
// vector starts at `row` and extends to the last element of the column. Alternatively
// sizes can be tuple of (row, numelems) and column vector starts at `row` and extends
// `numelems` elements. Function returns C.
func (C *FloatMatrix) Column(A *FloatMatrix, col int, sizes ...int) *FloatMatrix {
    if A == nil || C == nil {
        return nil
    }
    if col >= A.cols {
        return nil
    }
    var row int = 0
    var nr int = A.rows
    if len(sizes) == 1 {
        row = sizes[0]
        nr = A.rows - row
    } else if len(sizes) == 2 {
        row = sizes[0]
        nr = sizes[1]
    }
    if row + nr > A.rows {
        return nil
    }
    C.step = A.step
    C.rows = nr
    C.cols = 1
    if row < A.rows && col < A.cols {
        C.elems = A.elems[row+col*A.step:]
    } else {
        C.elems = nil
        C.rows = 0
        C.cols = 0
    }       
    return C
}

// Return matrix diagonal as row vector.
func (D *FloatMatrix) Diag(A *FloatMatrix) *FloatMatrix {
    return D.SubMatrix(A, 0, 0, 1, A.cols, A.step+1)
}



// Get element at [i, j]. Returns NaN if indexes are invalid.
func (A *FloatMatrix) Get(i, j int) float64 {
    if A.rows == 0 || A.cols == 0 {
        return 0.0
    }
    if i < 0 {
        i += A.rows
    }
    if j < 0 {
        j += A.cols
    }
    if i < 0 || i >= A.rows || j < 0 || j >= A.cols {
        return math.NaN()
    }
    return A.elems[i+j*A.step]
}

// Get element at index i. Returns NaN if index is invalid.
func (A *FloatMatrix) GetAt(i int) float64 {
    if i < 0 {
        i += A.rows*A.cols
    }
    if i < 0 || i >= A.rows*A.cols {
        return math.NaN()
    }
    c := i / A.rows
    r := i % A.rows
    return A.elems[r+c*A.step]
}

// Set element at [i, j]
func (A *FloatMatrix) Set(i, j int, v float64) {
    if A.rows == 0 || A.cols == 0 {
        return
    }
    if i < 0 {
        i += A.rows
    }
    if j < 0 {
        j += A.cols
    }
    if i < 0 || i >= A.rows || j < 0 || j >= A.cols {
        return
    }
    A.elems[i+j*A.step] = v
}

// Set element at index i. 
func (A *FloatMatrix) SetAt(i int, v float64) {
    if i < 0 {
        i += A.rows*A.cols
    }
    if i < 0 || i >= A.rows*A.cols {
        return
    }
    c := i / A.rows
    r := i % A.rows
    A.elems[r+c*A.step] = v;
}

// Make A copy of B.
func (A *FloatMatrix) Copy(B *FloatMatrix) *FloatMatrix {
    if B == nil || A == nil {
        return nil
    }
    if B.rows != A.rows || B.cols != A.cols {
        return nil;
    }
    if B.rows == 1 {
        // row vector
        for j := 0; j < B.cols; j++ {
            A.elems[j*A.step] = B.elems[j*B.step]
        }
        return B
    }
    // copy by column
    for j := 0; j < B.cols; j++ {
        copy(A.elems[j*A.step:], B.elems[j*B.step:B.rows+j*B.step])
    }
    return B
}

// Transpose matrix, A = B.T
func (A *FloatMatrix) Transpose(B *FloatMatrix) *FloatMatrix {
    if B == nil || A == nil {
        return nil
    }
    if A.rows != B.cols || A.cols != B.rows {
        return nil
    }
    for j := 0; j < B.cols; j++ {
        for i := 0; i < B.rows; i++ {
            A.elems[j+i*A.step] = B.elems[i+j*B.step]
        }
    }
    return B
}

// Absolute tolerance. Values v1, v2 are equal within tolerance if ABS(v1-v2) < ABSTOL + RELTOL*ABS(v2)
const ABSTOL = 1e-8
// Relative tolerance
const RELTOL = 1.0000000000000001e-05

func inTolerance(a, b, atol, rtol float64) bool {
    df := math.Abs(a - b)
    ref := atol + rtol * math.Abs(b)
    if df > ref {
        return false
    }
    return true
}

// Test if matrix A is equal to B within given tolenrances. Tolerances are given
// as tuple (abstol, reltol). If no tolerances are given default constants ABSTOL and
// RELTOL are used.
func (A *FloatMatrix) AllClose(B *FloatMatrix, tols ...float64) bool {
    var atol, rtol float64 = ABSTOL, RELTOL
    if A.rows != B.rows || A.cols != B.cols {
        return false
    }
    if len(tols) == 2 {
        atol = tols[0]
        rtol = tols[1]
    }
    if A.rows == 1 {
        // row vector
        for j := 0; j < A.cols; j++ {
            if ! inTolerance(A.elems[j*A.step], B.elems[j*B.step], atol, rtol) {
                return false
            }
        }
        return true
    }
    for j := 0; j < A.cols; j++ {
        for i := 0; i < A.rows; i++ {
            if ! inTolerance(A.elems[i+j*A.step], B.elems[i+j*B.step], atol, rtol) {
                return false
            }
        }
    }
    return true
}



// Convert matrix to string with spesific element format.
func (A *FloatMatrix) ToString(format string) string {
    s := ""
    if A == nil {
        return "<nil>"
    }
    for i := 0; i < A.rows; i++ {
        if i > 0 {
            s += "\n"
        }
        s += "["
        for j := 0; j < A.cols; j++ {
            if j > 0 {
                s += ", "
            }
            s += fmt.Sprintf(format, A.elems[i+j*A.step])
        }
        s += "]"
    }
    return s
}

func (A *FloatMatrix) String() string {
    return A.ToString("%9.2e")
}

// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
