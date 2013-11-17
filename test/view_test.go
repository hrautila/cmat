
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package test

import "testing"
import "github.com/hrautila/cmat"

func TestViews(t *testing.T) {
    var As cmat.FloatMatrix
    N := 7
    A := cmat.NewMatrix(N, N)
    zeromean := cmat.NewFloatNormSource()
    A.SetFrom(zeromean)
    dlast := A.Get(-1, -1)
    t.Logf("A[-1,-1] = %10.3e\n", dlast)
    for i := 0; i < N; i++ {
        As.SubMatrix(A, i, i)
        As.Set(0, 0, As.Get(0, 0)+float64(i+1))
        if N < 10 {
            t.Logf("A[%d,%d]   = %10.3e\n", i, i, As.Get(0,0))
        }
    }
    ok := float64(N) + dlast == A.Get(-1, -1)
    nC := int(A.Get(-1, -1) - dlast) 
    t.Logf("add 1.0 %d times [%10.3e to %10.3e]: %v\n", nC, dlast, A.Get(-1,-1), ok)
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
