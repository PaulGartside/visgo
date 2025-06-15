
package main

import (
  "bytes"
)

type BLine struct {
  data []byte
}

// Create slice of length filled with zeros
func (m *BLine) Init( length int ) {
  m.data = make( []byte, length )
}

func (m *BLine) Len() int {
  return len( m.data )
}

func (m *BLine) Cap() int {
  return cap( m.data )
}

// Set length to 0
func (m *BLine) Clear() {
  m.data = m.data[:0]
}

// Set all elements to zero
func (m *BLine) Zeroize() {
  // Sets all values in m.data.to 0 but does not change its length
  clear( m.data )
  // Manual way of zeroizing m.data:
//for k := range m.data {
//  m.data[k] = 0
//}
}

// Copy src_ln.data into m.data
func (m *BLine) Copy( src_ln BLine ) {
  m.SetLen( src_ln.Len() )
  copy( m.data[:], src_ln.data[:] )
}

// Set length without guaranteeing existing contents remain the same:
func (m *BLine) SetLen( length int ) {

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

func (m *BLine) GetB( b_num int ) byte {
  return m.data[ b_num ]
}

func (m *BLine) SetB( b_num int, B byte ) {
  m.data[ b_num ] = B
}

func (m *BLine) RemoveB( b_num int ) byte {

  var B byte = m.data[ b_num ]
  copy( m.data[b_num:], m.data[b_num+1:] )
  m.data = m.data[:len(m.data)-1]
  return B
}

func (m *BLine) PushB( B byte ) {
  m.data = append( m.data, B )
} 

func (m *BLine) PushSB( s_b []byte ) {
  m.data = append( m.data, s_b... )
} 

func (m *BLine) PushL( ln BLine ) {
  m.data = append( m.data, ln.data... )
}

func (m *BLine) InsertB( b_num int, B byte ) {
  // First append B to make sure data is large enough:
  m.PushB( B ) 
  copy( m.data[b_num+1:], m.data[b_num:] )
  m.data[ b_num ] = B
}

func (m *BLine) EqualL( ln BLine ) bool {

  return bytes.Equal( m.data, ln.data )
}

func (m *BLine) EqualStr( S string ) bool {

  return bytes.Equal( m.data, []byte(S) )
}

func (m *BLine) to_str() string {
  return string(m.data)
}

func (m *BLine) RemoveSpaces() {

  for k:=0; k<len(m.data); k++ {

    if( IsSpace( rune(m.data[k]) ) ) {
      copy( m.data[k:], m.data[k+1:] )
      m.data = m.data[:len(m.data)-1]
      k--
    }
  }
}

