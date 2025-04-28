
package main

type IntList struct {
  data []int
}

// Create slice of length filled with zeros
func (m *IntList) Init( length int ) {
  m.data = make( []int, length )
}

func (m *IntList) Len() int {
  return len( m.data )
}

func (m *IntList) Cap() int {
  return cap( m.data )
}

// Set length to 0
func (m *IntList) Clear() {
  m.data = m.data[:0]
}

func (m *IntList) CopyP( p_src *IntList ) {
  m.SetLen( p_src.Len() )
  copy( m.data, p_src.data )
}

// Set length without guaranteeing existing contents remain the same:
func (m *IntList) SetLen( length int ) {

  if( length < m.Len() ) {
    m.data = m.data[:length]
  } else if( m.Len() < length ) {
    if( length <= m.Cap() ) {
      for ; m.Len() < length; { m.data = append( m.data, 0 ) }
    } else {
      m.Init( length )
    }
  }
}

func (m *IntList) Get( idx int ) int {
  return m.data[ idx ]
}

func (m *IntList) Set( idx,val int ) {
  m.data[ idx ] = val
}

func (m *IntList) Remove( idx int ) int {

  var val int = m.data[ idx ]
  copy( m.data[idx:], m.data[idx+1:] )

  m.data = m.data[:len(m.data)-1]
  return val
}

func (m *IntList) Push( val int ) {
  m.data = append( m.data, val )
} 

func (m *IntList) Insert( idx int, val int ) {
  // First append val to make sure data is large enough:
  m.Push( val ) 
  copy( m.data[idx+1:], m.data[idx:] )
  m.data[ idx ] = val
}

