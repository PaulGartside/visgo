
package main

import (
//"bytes"
//"fmt"
  "slices"
  "strings"
  "unicode/utf8"
)

type RLine struct {
  data []byte
  enc_utf8 bool //< If this is false then data is byte encoded
}

// Create slice of length filled with zeros
func (m *RLine) Init( length int ) {
  m.data = make( []byte, length )
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

// Increase capacity by N and guarantees existing contents remain the same.
// Length is not changed.
func (m *RLine) Inc_Cap( N int ) {
  OLD_LEN := m.Len()
  NEW_CAP := m.Len() + N

  var old_data []byte = m.data
  m.data = make( []byte, OLD_LEN, NEW_CAP )
  copy( m.data, old_data )
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

func (m *RLine) GetB( B_num int ) byte {
  return m.data[ B_num ]
}

func (m *RLine) GetR( R_num int ) rune {
  var R rune = 0
  if( !m.enc_utf8 ) {
    R = rune(m.data[ R_num ])
  } else {
    B_offset_data := 0 // Byte offset in m.data
    for R_offset_data:=0; B_offset_data<len(m.data); R_offset_data++ {
      R_t, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
      if( R_num == R_offset_data ) {
        R = R_t
        break;
      }
      B_offset_data += R_size_data
    }
  }
  return R
}

func (m *RLine) to_SB( st int ) []byte {
  m_bb.Reset()
  for k:=st; k<len(m.data); k++ {
    m_bb.WriteByte( m.data[ k ] )
  }
  return m_bb.Bytes()
}

func (m *RLine) SetB( b_num int, B byte ) {
    m.data[ b_num ] = B
}

func (m *RLine) SetR( R_num int, R rune ) {

  if( !m.enc_utf8 ) {
    m.data[ R_num ] = byte(R)
  } else {
    R_size_in := utf8.RuneLen(R)
    if( 0 < R_size_in ) {
      B_offset_data := 0 // Byte offset in m.data
      for R_offset_data:=0; B_offset_data<len(m.data); R_offset_data++ {
        _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
        if( R_num == R_offset_data ) {
          if( R_size_in == R_size_data ) {
            utf8.EncodeRune( m.data[R_offset_data:], R )
          } else if( R_size_in < R_size_data ) {
            utf8.EncodeRune( m.data[R_offset_data:], R )
            copy( m.data[(R_offset_data+R_size_in):], m.data[(R_offset_data+R_size_data):] )
            size_diff := R_size_data - R_size_in
            m.data = m.data[:len(m.data)-size_diff]
          } else { // ( R_size_data < R_size_in )
            size_diff := R_size_in - R_size_data
            NEW_LEN := m.Len() + size_diff
            if( m.Cap() < NEW_LEN ) {
              m.Inc_Cap( size_diff + 16 )
            }
            m.data = m.data[:NEW_LEN]
            copy( m.data[(R_offset_data+R_size_in):], m.data[(R_offset_data+R_size_data):] )
            utf8.EncodeRune( m.data[R_offset_data:], R )
          }
          break;
        }
        B_offset_data += R_size_data
      }
    }
  }
}

func (m *RLine) RemoveR( R_num int ) rune {
  var R rune = 0

  if( !m.enc_utf8 ) {
    R = rune(m.data[ R_num ])
    copy( m.data[R_num:], m.data[R_num+1:] )
    m.data = m.data[:len(m.data)-1]
  } else {
    B_offset_data := 0 // Byte offset in m.data
    for R_offset_data:=0; B_offset_data<len(m.data); R_offset_data++ {
      R_t, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
      if( R_num == R_offset_data ) {
        R = R_t
        copy( m.data[R_offset_data:], m.data[R_offset_data+R_size_data:] )
        m.data = m.data[:len(m.data)-R_size_data]
        break;
      }
      B_offset_data += R_size_data
    }
  }
  return R
}

func (m *RLine) PushB( B byte ) {
  m.data = append( m.data, B )
}

func (m *RLine) PushR( R rune ) {
  if( !m.enc_utf8 ) {
    m.data = append( m.data, byte(R) )
  } else {
    R_size_in := utf8.RuneLen(R)
    if( 0 < R_size_in ) {
      OLD_LEN := m.Len()
      NEW_LEN := m.Len() + R_size_in
      if( m.Cap() < NEW_LEN ) {
        m.Inc_Cap( R_size_in + 16 )
      }
      // Increase m.data length by R_size_in:
      m.data = m.data[:OLD_LEN+R_size_in]
      utf8.EncodeRune( m.data[OLD_LEN:], R )
    }
  }
}

func (m *RLine) PushSR( s_r []rune ) {
  S := string( s_r )
  s_b := []byte( S )
  m.data = append( m.data, s_b... )
}

//func (m *RLine) PushStr( S string ) {
//  m.data = append( m.data, []byte(S) )
//}

func (m *RLine) PushL( ln RLine ) {
  m.data = append( m.data, ln.data... )
}

func (m *RLine) InsertR( R_num int, R rune ) {
  if( !m.enc_utf8 ) {
    B := byte(R)
    // First append B to make sure data is large enough:
    m.PushB( B )
    copy( m.data[R_num+1:], m.data[R_num:] )
    m.data[ R_num ] = B
  } else {
    R_size_in := utf8.RuneLen(R)
    if( 0 < R_size_in ) {
      B_offset_data := 0 // Byte offset in m.data
      for R_offset_data:=0; B_offset_data<len(m.data); R_offset_data++ {
        _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
        if( R_num == R_offset_data ) {
          // Insert R into m.data at B_offset_data
          m.PushR( R ) //< This will increase m.data cap if needed
          copy( m.data[R_offset_data+R_size_in:], m.data[R_offset_data:] )
          utf8.EncodeRune( m.data[R_offset_data:], R )
          break;
        }
        B_offset_data += R_size_data
      }
    }
  }
}

func (m *RLine) EqualL( ln RLine ) bool {

  return slices.Equal( m.data, ln.data )
}

func (m *RLine) EqualStr( S string ) bool {

  return slices.Equal( m.data, []byte(S) )
}

func (m *RLine) to_str() string {
  return string(m.data)
}

func (m *RLine) from_str( S string ) {
  m.data = []byte(S)
}

func (m *RLine) Compare( ln RLine ) int {

  // FIXME: Do this in a way that does not require memory allocation:
  return strings.Compare( m.to_str(), ln.to_str() )
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

    if( IsSpace( rune(m.data[k]) ) ) {
      copy( m.data[k:], m.data[k+1:] )
      m.data = m.data[:len(m.data)-1]
      k--;
    }
  }
}

func (m *RLine) ends_with( suffix string ) bool {
  S := m.to_str()
  len_S := len(S)
  len_suffix := len(suffix)
  ends_w := len_suffix <= len_S && S[ len_S - len_suffix: ] == suffix
  return ends_w
}

