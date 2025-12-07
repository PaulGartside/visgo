
package main

import (
//"fmt"
)

type Vector[T any] struct {
  data []T
}

// Create slice of length filled with zeros
func (m *Vector[T]) Init( length int ) {
  m.data = make( []T, length )
}

func (m *Vector[T]) Len() int {
  return len(m.data)
}

func (m *Vector[T]) Cap() int {
  return cap( m.data )
}

// Set length to 0
func (m *Vector[T]) Clear() {
  m.data = m.data[:0]
}

func (m *Vector[T]) Get( idx int ) T {
  return m.data[ idx ]
}

func (m *Vector[T]) GetP( idx int ) *T {
  return &m.data[ idx ]
}

func (m *Vector[T]) Set( idx int, val T ) {
  m.data[ idx ] = val
}

func (m *Vector[T]) Push( val T ) {
  m.data = append( m.data, val )
}

//func (m *Vector[T]) Pop() (T, bool) {
//  var rval T
//  ok := false
//  LEN := len(m.data)
//
//  if( 0 < LEN ) {
//    rval = m.data[ LEN-1 ]
//    ok = true
//    m.data = m.data[:LEN-1]
//  }
//  return rval, ok
//}

//func (m *Vector[T]) Pop(p_rval *T) bool {
//  ok := false
//  if( nil != p_rval ) {
//    LEN := len(m.data)
//
//    if( 0 < LEN ) {
//      *p_rval = m.data[ LEN-1 ]
//      ok = true
//      m.data = m.data[:LEN-1]
//    }
//  }
//  return ok
//}

func (m *Vector[T]) Pop(p_rval *T) bool {
  ok := false
  LEN := len(m.data)

  if( 0 < LEN ) {
    ok = true
    if( nil != p_rval ) {
      *p_rval = m.data[ LEN-1 ]
    }
    m.data = m.data[:LEN-1]
  }
  return ok
}

// Copy src_vec.data into m.data
//
func (m *Vector[T]) Copy( src_vec Vector[T] ) {
  m.SetLen( src_vec.Len() )
  copy( m.data[:], src_vec.data[:] )
}

// Set length while guaranteeing existing contents remain the same:
//
func (m *Vector[T]) SetLen( length int ) {

  if( length < m.Len() ) {
    // Contents up to length-1 preserved:
    m.data = m.data[:length]

  } else if( m.Len() < length ) {
    if( length <= m.Cap() ) {
      // Contents preserved, zero values appended to end:
      var t T
      for ; m.Len() < length; { m.data = append( m.data, t ) }
    } else {
      // Capacity increased. Contents preserved, zero values appended to end:
      var old []T = m.data
      len_old := len(old)
      m.Init( length )
      for k:=0; k < len_old; k++ { m.data[k] = old[k] }
    }
  }
}

func (m *Vector[T]) Insert( idx int, val T ) {
  // First append val to make sure data is large enough:
  m.Push( val ) 
  copy( m.data[idx+1:], m.data[idx:] )
  m.data[ idx ] = val
}

func (m *Vector[T]) Remove( idx int ) T {

  var val T = m.data[ idx ]
  copy( m.data[idx:], m.data[idx+1:] )

  m.data = m.data[:len(m.data)-1]
  return val
}

type IntList = Vector[int]

