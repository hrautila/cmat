
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package test

import "testing"
import "github.com/hrautila/cmat"

func TestCopy(t *testing.T) {
    M := 9
    N := 9
    A := cmat.NewMatrix(M, N)
    B := cmat.NewMatrix(M, N)
    twos := cmat.NewFloatConstSource(2.0)
    B.SetFrom(twos)
    A.Copy(B)
    ok := A.AllClose(B)  
    if ! ok {
        t.Logf("copy status: %v\n", ok)
        if (N < 9) {
            t.Logf("A\n%v\n", A)
        }
    }
    if (N < 10) {
        t.Logf("A\n%v\n", A)
    }
}


func TestSubMatrixCopy(t *testing.T) {
    var subA, subB cmat.FloatMatrix
    M := 9
    N := 9
    A := cmat.NewMatrix(M, N)
    B := cmat.NewMatrix(M, N)
    twos := cmat.NewFloatConstSource(2.0)
    B.SetFrom(twos)
    subA.SubMatrix(A, 1, 1, M-2, N-2)
    subB.SubMatrix(B, 1, 1, M-2, N-2)
    subA.Copy(&subB)
    ok := subA.AllClose(&subB)  
    if ! ok {
        t.Logf("copy status: %v\n", ok)
        if (N < 9) {
            t.Logf("subA\n%v\n", subA)
        }
    }
    if (N < 10) {
        t.Logf("A\n%v\n", A)
    }
}

// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
