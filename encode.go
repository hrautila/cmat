
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

import (
    "encoding/gob"
    "bytes"
)

const encodeVersion = 1

func (A *FloatMatrix) GobEncode() ([]byte, error) {
    var prefix uint8 = encodeVersion
    var b bytes.Buffer
    enc := gob.NewEncoder(&b)
    enc.Encode(prefix)
    enc.Encode(A.rows)
    enc.Encode(A.cols)
    for i := 0; i < A.cols; i++ {
        col := A.elems[i*A.step:i*A.step+A.rows]
        enc.Encode(col)
    }
    return b.Bytes(), nil
}


func (A *FloatMatrix) GobDecode(buf []byte) (err error) {
    var prefix uint8

    b := bytes.NewBuffer(buf)
    dec := gob.NewDecoder(b)
    err = dec.Decode(&prefix)
    if err != nil { return }

    err = dec.Decode(&A.rows)
    if err != nil { return }

    err = dec.Decode(&A.cols)
    if err != nil { return }

    A.step = A.rows
    A.elems = make([]float64, A.rows*A.cols, A.rows*A.cols)
    for i := 0; i < A.cols; i++ {
        var ebuf []float64
        err = dec.Decode(&ebuf)
        if err != nil { return }
        copy(A.elems[i*A.step:], ebuf)
    }
    return
}



// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
