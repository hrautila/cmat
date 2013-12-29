
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/cmat package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package cmat

import (
    "encoding/gob"
    "bytes"
    "fmt"
    "errors"
    "strconv"
    "strings"
)

const encodeVersion = 1

// GobEncode matrix. If A is a submatrix elements outside submatrix are not included.
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


// Decode a matrix.
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

func (A *FloatMatrix) MarshalJSON() ([]byte, error) {
    s := fmt.Sprintf("{\"rows\":%d,\"cols\":%d,\"elems\":[", A.rows, A.cols)
    for i := 0; i < A.cols; i++ {
        if i > 0 {
            s += ","
        }
        for k, v := range A.elems[i*A.step:i*A.step+A.rows] {
            if k > 0 {
                s += ","
            }
            s += fmt.Sprintf("%.16e", v)
        }
    }
    s += "]}"
    return bytes.NewBufferString(s).Bytes(), nil
}

func (A *FloatMatrix) UnmarshalJSON(buf []byte) (err error) {
    var j int
    A.rows = 0
    A.cols = 0
    for _, part := range bytes.SplitN(buf, []byte(","), 3) {
        if bytes.Contains(part, []byte("rows")) {
            j = bytes.Index(part, []byte(":"))
            r, _ := strconv.ParseInt(string(part[j+1:]), 10, 0)
            A.rows = int(r)
            A.step = A.rows
        } else if bytes.Contains(part, []byte("cols")) {
            j = bytes.Index(part, []byte(":"))
            c, _ := strconv.ParseInt(string(part[j+1:]), 10, 0)
            A.cols = int(c)
        } else if A.rows*A.cols > 0 {
            // elements here
            j = bytes.Index(part, []byte("["))
            if j < 0 {
                return errors.New("matrix elements not found")
            }
            A.elems = make([]float64, A.rows*A.cols)
            input := bytes.NewBuffer(part[j+1:])
            for i := 0; i < A.rows*A.cols; i++ {
                s, err := input.ReadString(',')
                if err == nil {
                    A.elems[i], _ = strconv.ParseFloat(s[:len(s)-1], 64)
                } else {
                    //fmt.Printf("err: %v\ns:%v\n", err, s)
                    j = strings.Index(s, "]")
                    A.elems[i], _ = strconv.ParseFloat(s[:j], 64)
                }
            }
        }
    }

    return nil
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
