
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

type JoinType int
const (
    STACK = iota
    AUGMENT
)

func NewJoin(how JoinType, mlist... *FloatMatrix) *FloatMatrix {
    var nrows, ncols, maxrow, maxcol int
    maxrow = 0
    maxcol = 0
    for _, m := range mlist {
        r, c := m.Size()
        nrows += r
        ncols += c
        if r > maxrow {
            maxrow = r
        }
        if c > maxcol {
            maxcol = c
        }
    }
    newr, newc := 0, 0
    switch how {
    case STACK:
        newr = nrows
        newc = maxcol
    case AUGMENT:
        fallthrough
    default:
        newc = ncols
        newr = maxrow
    }
    M := NewMatrix(newr, newc)
    
    crow := 0
    ccol := 0
    for _, m := range mlist {
        var T FloatMatrix
        r, c := m.Size()
        if how == STACK {
            T.SubMatrix(M, crow, 0, r, c)
        } else {
            T.SubMatrix(M, 0, ccol, r, c)
        }
        T.Copy(m)
        crow += r
        ccol += c
    }
    return M
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
