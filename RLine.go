
package main

import (
//"bytes"
//"fmt"
  "slices"
//"unicode/utf8"
)

type RLine struct {
  data []rune
  enc_utf8 bool
}

// Create slice of length filled with zeros
func (m *RLine) Init( length int ) {
  m.data = make( []rune, length )
}

func (m *RLine) Len() int {
  return len( m.data )
}

func (m *RLine) Cap() int {
  return cap( m.data )
}

// Set length to zero
func (m *RLine) Clear() {
  m.data = m.data[:0]
}

// Set all elements to zero
func (m *RLine) Zeroize() {
  // Sets all values in m.data.to 0 but does not change its length
  clear( m.data )
  // Manual way of zeroizing m.data:
//for k := range m.data {
//  m.data[k] = 0
//}
}

// Copy src_ln.data into m.data
//func (m *RLine) Copy( src_ln RLine ) {
//
//  src_ll := src_ln.Len()
//
//  if( m.Cap() < src_ll ) {
//    m.Init( src_ll )
//  }
//  copy( m.data[0:src_ll], src_ln.data[0:src_ll] )
////copy( m.data[:], src_ln.data[:] )
//}

// Copy src_ln.data into m.data
func (m *RLine) Copy( src_ln RLine ) {

  m.SetLen( src_ln.Len() )
  copy( m.data[:], src_ln.data[:] )
}

// Set length without guaranteeing existing contents remain the same:
func (m *RLine) SetLen( length int ) {

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

func (m *RLine) GetR( r_num int ) rune {
  return m.data[ r_num ]
}

func (m *RLine) to_SB( st int ) []byte {
  m_bb.Reset()
  for k:=st; k<len(m.data); k++ {
    if( m.enc_utf8 ) { m_bb.WriteRune( m.data[ k ] )
    } else           { m_bb.WriteByte( byte(m.data[ k ]) )
    }
  }
  return m_bb.Bytes()
}

func (m *RLine) SetR( r_num int, R rune ) {
  m.data[ r_num ] = R
}

//func (m *RLine) GetB( b_num int ) byte {
//  // In insert mode, at the end of a line, b_num == len( m.data ),
//  if b_num < len( m.data ) {
//    return m.data[ b_num ]
//  }
//  return ' '
//}

func (m *RLine) RemoveR( r_num int ) rune {

  var R rune = m.data[ r_num ]
  copy( m.data[r_num:], m.data[r_num+1:] )
  m.data = m.data[:len(m.data)-1]
  return R
}

func (m *RLine) PushR( R rune ) {
  m.data = append( m.data, R )
} 

func (m *RLine) PushSR( s_r []rune ) {
  m.data = append( m.data, s_r... )
} 

func (m *RLine) PushL( ln RLine ) {
  m.data = append( m.data, ln.data... )
}

func (m *RLine) InsertR( r_num int, R rune ) {
  // First append R to make sure data is large enough:
  m.PushR( R ) 
  copy( m.data[r_num+1:], m.data[r_num:] )
  m.data[ r_num ] = R
}

func (m *RLine) EqualL( ln RLine ) bool {

  return slices.Equal( m.data, ln.data )
}

func (m *RLine) EqualStr( S string ) bool {

  return slices.Equal( m.data, []rune(S) )
}

func (m *RLine) to_str() string {
  return string(m.data)
}

func (m *RLine) from_str( S string ) {
  m.data = []rune(S)
}

// Not sure if this method is need.  Just use to_str().
//
//func (m *RLine) to_bytes() []byte {
//  return []byte( m.to_str() )
//}

// Longer version of to_bytes().
// Not sure if this method is need.  Just use to_str().
//
//func (m *RLine) to_bytes() []byte {
//  num_bytes := 0
//  for _, R := range m.data {
//    num_bytes += utf8.RuneLen( R )
//  }
//  s_b := make( []byte, num_bytes )
//
//  byte_offset := 0
//  for _, R := range m.data {
//    byte_offset += utf8.EncodeRune( s_b[byte_offset:], R )
//  }
//  return s_b
//}

func (m *RLine) RemoveSpaces() {

  for k:=0; k<len(m.data); k++ {

    if( IsSpace( m.data[k] ) ) {
      copy( m.data[k:], m.data[k+1:] )
      m.data = m.data[:len(m.data)-1]
      k--;
    }
  }
}

