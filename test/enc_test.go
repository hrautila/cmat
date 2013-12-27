
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package test

import (
    "github.com/hrautila/cmat"
    "testing"
    "encoding/gob"
    "bytes"
)

func TestGob(t *testing.T) {
    var B cmat.FloatMatrix
    var network bytes.Buffer
    N := 16
    A := cmat.NewMatrix(N, N)
    zeromean := cmat.NewFloatNormSource()
    A.SetFrom(zeromean)

    enc := gob.NewEncoder(&network)
    dec := gob.NewDecoder(&network)
    
    // encode to network
    err := enc.Encode(A)
    if err != nil {
        t.Logf("encode error: %v\n", err)
        t.FailNow()
    }

    // decode from network
    err = dec.Decode(&B)
    if err != nil {
        t.Logf("decode error: %v\n", err)
        t.FailNow()
    }

    t.Logf("A == B: %v\n", B.AllClose(A))
}

func TestSubMatrixGob(t *testing.T) {
    var B, As cmat.FloatMatrix
    var network bytes.Buffer
    N := 32
    A := cmat.NewMatrix(N, N)
    zeromean := cmat.NewFloatNormSource()
    A.SetFrom(zeromean)
    As.SubMatrix(A, 3, 3, N-6, N-6)

    enc := gob.NewEncoder(&network)
    dec := gob.NewDecoder(&network)
    
    // encode to network
    err := enc.Encode(&As)
    if err != nil {
        t.Logf("encode error: %v\n", err)
        t.FailNow()
    }

    // decode from network
    err = dec.Decode(&B)
    if err != nil {
        t.Logf("decode error: %v\n", err)
        t.FailNow()
    }

    ar,ac  := As.Size()
    br, bc := B.Size()
    t.Logf("As[%d,%d] == B[%d,%d]: %v\n", ar, ac, br, bc, B.AllClose(&As))
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
