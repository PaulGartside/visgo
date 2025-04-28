
package main

import (
//"fmt"
)

type Vector[T any] struct {
  vals []T
}

// Create slice of length filled with zeros
func (m *Vector[T]) Init( length int ) {
  m.vals = make( []T, length )
}

func (m *Vector[T]) Len() int {
  return len(m.vals)
}

func (m *Vector[T]) Cap() int {
  return cap( m.vals )
}

// Set length to 0
func (m *Vector[T]) Clear() {
  m.vals = m.vals[:0]
}

func (m *Vector[T]) Get( idx int ) T {
  return m.vals[ idx ]
}

func (m *Vector[T]) Push( val T ) {
  m.vals = append( m.vals, val )
}

func (m *Vector[T]) Pop() (T, bool) {
  var rval T
  ok := false
  LEN := len(m.vals)

  if( 0 < LEN ) {
    rval = m.vals[ LEN-1 ]
    ok = true
    m.vals = m.vals[:LEN-1]
  }
  return rval, ok
}

// Copy src_vec.vals into m.vals
func (m *Vector[T]) Copy( src_vec Vector[T] ) {
  m.SetLen( src_vec.Len() )
  copy( m.vals[:], src_vec.vals[:] )
}

// Set length without guaranteeing existing contents remain the same:
func (m *Vector[T]) SetLen( length int ) {

  if( length < m.Len() ) {
    m.vals = m.vals[:length]
  } else if( m.Len() < length ) {
    if( length <= m.Cap() ) {
      var t T
      for ; m.Len() < length; { m.vals = append( m.vals, t ) }
    } else {
      m.Init( length )
    }
  }
}

