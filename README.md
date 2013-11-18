CMAT - Column major MATrix
--------------------------


Another implementation of a column major matrix package. Similiar to the MATRIX package but with
more restricted interface. Linear algebra operations available in GOMAS or NETLIB package.


### FloatMatrix 


  Creating instances
  
    NewMatrix(r, c)                Create a new FloatMatrix of size r rows, c columns
    NewCopy(A)                     Create a new FloatMatrix as copy of A
    NewJoin(join, mlist...)        Create a new compound FloatMatrix from argument matrices
    MakeMatrix(buf, r, c)          Create a new FloatMatrix, use buf as element store

  Basic attributes
  
    A.Size() (int,int)             Get size of A as (rows, cols)
    A.Len() int                    Get number of element in A
    A.Stride() int                 Get row stride of A
    A.Data() []float64             Get raw data elements

  Getting/setting elements
  
    A.Get(i, j) float64            Get value at [i,j]
    A.Set(i, j, v)                 Set element at [i,j] to value v
    A.GetAt(i) float64             Get i'th element 
    A.SetAt(i, v)                  Set i'th element to value v

  Strings
  
    A.ToString(format) string      Convert to string, format is element format
    A.String() string              Convert to string

  Testing
  
    A.AllClose(B, tols...) bool    Test if A is close to B, within tolerance

  Copying
  
    A.Copy(B)                      Copy B to A.
    A.Transpose(B)                 Copy B.T to A

  Matrix views
  
    A.SubMatrix(B, r, c, nr, nc)   Make A submatrix of B starting at [r,c]
    A.Column(B, col)               Make A column vector of B
    A.Row(B, row)                  Make A row vector of B
    A.Diag(B)                      Make A diagonal row vector of B

  Transforming and setting
  
    A.SetFrom(src, bits)           Get values for A from source
    A.Map(mapping)                 Apply mapping to all elements of A


### Data sources


    interface FloatSource         
       Get(i, j) float64

    NewFloatConstSource(val)           Create source that provides const value
    NewFloatNormSource(stddev,mean)    Create source of normally distributed floats
    NewFloatUniformSource(scale,shift) Create source of uniformly distributed floats


### Mapping types


    interface FloatMapping
       Eval(i, j, v) float64

    Following types implement FloatMapping interface for Map() method.


    type FloatEvaluator {       Generic from value v at [i,j] to new value
         Callable func(i, j, v) float64  
    }

    type FloatFunction {        Mapping to new value Callable(v)
         Callable func(v) float64  
    }

    type Float2Function {       Mapping to new value Callable(v, Const)
         Callable func(v, c) float64  
	 Const float64
    }

    See file mapping.go for some examples.
