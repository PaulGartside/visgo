
package main

import (
//"bytes"
//"fmt"
  "hash/crc32"
  "slices"
//"strings"
  "unicode"
  "unicode/utf8"
)

type RLine struct {
  data []byte

  last_R_idx int           // Index of rune last accessed
  last_R_idx_B_offset int  // Byte offset of last_R_idx

  chksum_diff uint32
  chksum_diff_valid bool
}

// Create slice of size bytes filled with zeros
func (m *RLine) Init( size int ) {
  m.data = make( []byte, size )
}

func (m *RLine) Size() int {
  return len( m.data )
}

// Length in bytes, or number of contained bytes
//
func (m *RLine) LenB() int {
  return len( m.data )
}

// Length in runes, or number of contained runes
//
func (m *RLine) LenR() int {
  return utf8.RuneCount( m.data )
}

// Capacity in Bytes
//
func (m *RLine) Cap() int {
  return cap( m.data )
}

// Set length to zero
func (m *RLine) Clear() {
  m.data = m.data[:0]

  m.last_R_idx = 0
  m.last_R_idx_B_offset = 0

  m.chksum_diff_valid = false
}

// Set all elements to zero
func (m *RLine) Zeroize() {
  // Sets all values in m.data.to 0 but does not change its length
  clear( m.data )
  // Manual way of zeroizing m.data:
//for k := range m.data {
//  m.data[k] = 0
//}
  m.last_R_idx = 0
  m.last_R_idx_B_offset = 0

  m.chksum_diff_valid = false
}

// Increase capacity by N bytes and guarantees existing contents remain the same.
// Length is not changed.
func (m *RLine) Inc_Cap( N int ) {
  OLD_LEN := m.LenB()
  NEW_CAP := m.LenB() + N

  var old_data []byte = m.data
  m.data = make( []byte, OLD_LEN, NEW_CAP )
  copy( m.data, old_data )
}

// Set length while guaranteeing existing contents remain the same:
//
func (m *RLine) SetLen( length int ) {

  if( length < m.LenB() ) {
    // Contents up to length-1 preserved:
    m.data = m.data[:length]

    if( length <= m.last_R_idx_B_offset ) {
      m.last_R_idx = 0
      m.last_R_idx_B_offset = 0
    }
  } else if( m.LenB() < length ) {
    if( length <= m.Cap() ) {
      // Contents preserved, zero values appended to end:
      for ; m.LenB() < length; { m.data = append( m.data, 0 ) }
    } else {
      // Capacity increased. Contents preserved, zero values appended to end:
      var old []byte = m.data
      len_old := len(old)
      m.Init( length )
      for k:=0; k < len_old; k++ { m.data[k] = old[k] }
    }
  }
}

// Copy src_ln.data into m.data
func (m *RLine) Copy( src_ln RLine ) {

  m.SetLen( src_ln.LenB() )
  copy( m.data[:], src_ln.data[:] )

  m.last_R_idx = 0
  m.last_R_idx_B_offset = 0

  m.chksum_diff_valid = false
}

func (m *RLine) GetB( B_num int ) byte {
  return m.data[ B_num ]
}

// Gets R_num rune in m.data.
// Returns rune and index in m.data of that rune
//
//func (m *RLine) GetR( R_num int ) rune, int, int {
//  var R rune = 0
//  B_offset_data := 0 // Byte offset in m.data
//  LEN := len(m.data)
//  for R_offset_data:=0; B_offset_data<LEN; R_offset_data++ {
//    R_t, R_size := utf8.DecodeRune( m.data[B_offset_data:] )
//    if( R_num == R_offset_data ) {
//      R = R_t
//      break
//    }
//    B_offset_data += R_size
//  }
//  return R, R_size, B_offset_data
//}

// Gets R_num rune in m.data.
// Returns (Rune) and
//         (size in bytes of Rune) and
//         (Byte position in line) of that Rune
//
//func (m *RLine) GetR( R_num int ) (rune, int, int) {
//  var R rune = 0
//  R_size := 0
//  B_offset_data := 0 // Byte offset in m.data
//  LEN := len(m.data)
//  for R_index:=0; B_offset_data<LEN; R_index++ {
//    var R_t rune
//    R_t, R_size = utf8.DecodeRune( m.data[B_offset_data:] )
//    if( R_num == R_index ) {
//      R = R_t
//      break
//    }
//    B_offset_data += R_size
//  }
//  return R, R_size, B_offset_data
//}

// Gets R_num rune in m.data.
// Returns (Rune) and
//         (size in bytes of Rune) and
//         (Byte position in line) of that Rune
//
func (m *RLine) GetR( R_num int ) (rune, int, int) {
  var R rune = 0
  R_size_data := 0
  B_offset_data := 0 // Byte offset in m.data
  LEN := len(m.data)

  if( 0 < LEN ) {
    if( R_num <= 0 ) {
      R, R_size_data = utf8.DecodeRune( m.data[0:] )
      m.last_R_idx = 0
      m.last_R_idx_B_offset = 0

    } else if( R_num == m.last_R_idx ) {
      R, R_size_data = utf8.DecodeRune( m.data[m.last_R_idx_B_offset:] )
      B_offset_data = m.last_R_idx_B_offset
      // m.last_R_idx and m.last_R_idx_B_offset are unchanged

    } else if( m.last_R_idx < R_num ) {
      // Search forward for R_num
      B_offset_data = m.last_R_idx_B_offset
      for R_index:=m.last_R_idx; B_offset_data<LEN; R_index++ {
        var R_t rune
        R_t, R_size_data = utf8.DecodeRune( m.data[B_offset_data:] )
        if( R_index == R_num ) {
          R = R_t
          m.last_R_idx = R_index
          m.last_R_idx_B_offset = B_offset_data
          break
        }
        B_offset_data += R_size_data
      }
    } else { // R_num < m.last_R_idx
      if( m.last_R_idx == m.last_R_idx_B_offset ) {
        R = rune(m.data[ R_num ])
        R_size_data = 1
        B_offset_data = R_num
        // No multi-byte utf8 before m.last_R_idx, each rune is only one byte
      } else {
        // Search backwards for R_num
        B_offset_data = m.last_R_idx_B_offset
        for R_index:=m.last_R_idx-1; 0<=B_offset_data && B_offset_data<LEN; R_index-- {
          B_offset_data = find_prior_rune_start_or_invalid_utf8( m.data, B_offset_data )
          if( R_index == R_num ) {
            R, R_size_data = utf8.DecodeRune( m.data[B_offset_data:] )
            m.last_R_idx = R_index
            m.last_R_idx_B_offset = B_offset_data
            break
          }
        }
      }
    }
  }
  return R, R_size_data, B_offset_data
}

// Gets rune in m.data at byte offset in m.data
// Returns rune, rune size
//
func (m *RLine) GetRatB( B_offset_data int ) (rune, int) {
  var R rune = 0
  R_size := 0
  if( B_offset_data<len(m.data) ) {
    R, R_size = utf8.DecodeRune( m.data[B_offset_data:] )
  }
  return R, R_size
}

// Returns a copy of m.data
//
func (m *RLine) to_SB( st int ) []byte {
  m_bb.Reset()
  for k:=st; k<len(m.data); k++ {
    m_bb.WriteByte( m.data[ k ] )
  }
  return m_bb.Bytes()
}

func (m *RLine) SetB( B_num int, B byte ) {
  m.data[ B_num ] = B

  if( B_num < m.last_R_idx_B_offset ) {
    m.last_R_idx = 0
    m.last_R_idx_B_offset = 0
  }
  m.chksum_diff_valid = false
}

//func (m *RLine) SetRatB( B_num int, R rune) int {
//
//  m.chksum_diff_valid = false
//}

// Sets R_num rune in m.data to R.
// Returns rune size, and byte offset in m.data of that rune
//
//func (m *RLine) SetR( R_num int, R_in rune ) int {
//  R_size_in := utf8.RuneLen(R_in)
//  B_offset_data := 0 // Byte offset in m.data
//
//  if( 0 < R_size_in ) {
//    LEN := len(m.data)
//    // FIXME: use m.last_R_idx and m.last_R_idx_B_offset to speed up search:
//    for R_index:=0; B_offset_data<LEN; R_index++ {
//      _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//
//      if( R_num == R_index ) {
//        m.last_R_idx = R_num
//        m.last_R_idx_B_offset = B_offset_data
//
//        m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
//        break
//      }
//      B_offset_data += R_size_data
//    }
//    m.chksum_diff_valid = false
//  }
//  return B_offset_data
//}

// Sets R_num rune in m.data to R_in.
// Returns byte offset in m.data of set rune
//
//func (m *RLine) SetR( R_num int, R_in rune ) int {
//  B_offset_data := 0 // Byte offset in m.data
//  R_size_in := utf8.RuneLen(R_in)
//  LEN := len(m.data)
//
//  if( 0 < LEN && 0 < R_size_in ) {
//    if( R_num <= 0 ) {
//      _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//      m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
//      m.last_R_idx = 0
//      m.last_R_idx_B_offset = 0
//
//    } else if( R_num == m.last_R_idx ) {
//      _, R_size_data := utf8.DecodeRune( m.data[m.last_R_idx_B_offset:] )
//      B_offset_data = m.last_R_idx_B_offset
//      m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
//      // m.last_R_idx and m.last_R_idx_B_offset are unchanged
//
//    } else if( m.last_R_idx < R_num ) {
//      // Search forward for R_num
//      B_offset_data = m.last_R_idx_B_offset
//      for R_index:=m.last_R_idx; B_offset_data<LEN; R_index++ {
//        _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//
//        if( R_index == R_num ) {
//          m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
//          m.last_R_idx = R_index
//          m.last_R_idx_B_offset = B_offset_data
//          break
//        }
//        B_offset_data += R_size_data
//      }
//    } else { // R_num < m.last_R_idx
//      // Search backwards for R_num
//      B_offset_data = m.last_R_idx_B_offset
//      for R_index:=m.last_R_idx-1; 0<=B_offset_data && B_offset_data<LEN; R_index-- {
//        B_offset_data = find_prior_rune_start_or_invalid_utf8( m.data, B_offset_data )
//        if( R_index == R_num ) {
//          _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//          m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
//          m.last_R_idx = R_index
//          m.last_R_idx_B_offset = B_offset_data
//          break
//        }
//      }
//    }
//    m.chksum_diff_valid = false
//  }
//  return B_offset_data
//}

// Sets R_num rune in m.data to R_in.
// Returns byte offset in m.data of set rune
//
func (m *RLine) SetR( R_num int, R_in rune ) int {

  _, R_size_data, B_offset_data := m.GetR( R_num ) 

  if( 0 < R_size_data ) {
    // A rune was found, so R_num is not past end of m.data.
    R_size_in := utf8.RuneLen(R_in)
    m.SetR_ReplaceRatB( R_in,B_offset_data, R_size_in,R_size_data )
    // m.last_R_idx and m.last_R_idx_B_offset were set in m.GetR()
  } else {
    // A rune was NOT found, so R_num is past end of m.data.
    // Append R_in:
    LEN := len(m.data)
    if( LEN <= B_offset_data ) {
      // (LEN == 0) or (m.LenR() <= R_num), so append R_in:
      m.PushR( R_in ) //< This will increase m.data cap if needed

      // When R_in is appended, m.last_R_idx and m.last_R_idx_B_offset
      // dont need to change.
    }
  }
  m.chksum_diff_valid = false

  return B_offset_data
}

// Helper function of SetR.
// Replace rune at byte B_offset_data with R_in.
// R_size_in = size of R_in
// R_size_data = size of rune at byte B_offset_data
//
func (m *RLine) SetR_ReplaceRatB( R_in rune, B_offset_data int, R_size_in,R_size_data int ) {

  if( R_size_in == R_size_data ) {
    utf8.EncodeRune( m.data[B_offset_data:], R_in )

  } else if( R_size_in < R_size_data ) {
    utf8.EncodeRune( m.data[B_offset_data:], R_in )
    copy( m.data[(B_offset_data+R_size_in):], m.data[(B_offset_data+R_size_data):] )
    size_diff := R_size_data - R_size_in
    m.data = m.data[:len(m.data)-size_diff]

  } else { // ( R_size_data < R_size_in )
    size_diff := R_size_in - R_size_data
    NEW_LEN := m.LenB() + size_diff
    if( m.Cap() < NEW_LEN ) {
      m.Inc_Cap( size_diff + 16 )
    }
    m.data = m.data[:NEW_LEN]
    copy( m.data[(B_offset_data+R_size_in):], m.data[(B_offset_data+R_size_data):] )
    utf8.EncodeRune( m.data[B_offset_data:], R_in )
  }
}

func (m *RLine) RemoveB( B_num int ) byte {

  var B byte = m.data[ B_num ]
  copy( m.data[B_num:], m.data[B_num+1:] )
  m.data = m.data[:len(m.data)-1]

  if( B_num < m.last_R_idx_B_offset ) {
    m.last_R_idx = 0
    m.last_R_idx_B_offset = 0
  }
  m.chksum_diff_valid = false

  return B
}

//func (m *RLine) RemoveR( R_num int ) rune {
//  var R rune = 0
//
//  B_offset_data := 0 // Byte offset in m.data
//  LEN := len(m.data)
//  // FIXME: use m.last_R_idx and m.last_R_idx_B_offset to speed up search:
//  for R_index:=0; B_offset_data<LEN; R_index++ {
//    R_t, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//    if( R_num == R_index ) {
//      R = R_t
//      copy( m.data[B_offset_data:], m.data[B_offset_data+R_size_data:] )
//      m.data = m.data[:LEN-R_size_data]
//
//      m.last_R_idx = R_index
//      m.last_R_idx_B_offset = B_offset_data
//      m.chksum_diff_valid = false
//      break
//    }
//    B_offset_data += R_size_data
//  }
//  return R
//}

func (m *RLine) RemoveR( R_num int ) rune {

  R, R_size_data, B_offset_data := m.GetR( R_num ) 

  if( 0 < R_size_data ) {
    // m.last_R_idx and m.last_R_idx_B_offset were set in m.GetR().
    LEN := len(m.data)
    copy( m.data[B_offset_data:], m.data[B_offset_data+R_size_data:] )
    m.data = m.data[:LEN-R_size_data]

    if( len(m.data) <= B_offset_data ) {
      // R_num was last rune, and it was just removed,
      // so clear m.last_R_idx and m.last_R_idx_B_offset:
      m.last_R_idx = 0
      m.last_R_idx_B_offset = 0
    }
    m.chksum_diff_valid = false
  }
  return R
}

func (m *RLine) PushB( B byte ) {
  m.data = append( m.data, B )

  m.chksum_diff_valid = false
}

func (m *RLine) PushR( R rune ) {
  R_size_in := utf8.RuneLen(R)
  if( 0 < R_size_in ) {
    OLD_LEN := m.LenB()
    NEW_LEN := m.LenB() + R_size_in
    if( m.Cap() < NEW_LEN ) {
      m.Inc_Cap( R_size_in + 16 )
    }
    // Increase m.data length by R_size_in:
    m.data = m.data[:OLD_LEN+R_size_in]
    utf8.EncodeRune( m.data[OLD_LEN:], R )

    m.chksum_diff_valid = false
  }
}

//func (m *RLine) PushSR( s_r []rune ) {
//  S := string( s_r )
//  s_b := []byte( S )
//  m.data = append( m.data, s_b... )
//
////m.chksum_valid = false
//  m.chksum_diff_valid = false
//}

func (m *RLine) PushStr( S string ) {

  s_b := []byte( S )
  m.data = append( m.data, s_b... )

  m.chksum_diff_valid = false
}

func (m *RLine) PushL( ln RLine ) {
  m.data = append( m.data, ln.data... )

  m.chksum_diff_valid = false
}

func (m *RLine) InsertB( B_num int, B byte ) {
  // First append B to make sure data is large enough:
  m.PushB( B )
  copy( m.data[B_num+1:], m.data[B_num:] )
  m.data[ B_num ] = B

  if( B_num < m.last_R_idx_B_offset ) {
    m.last_R_idx = 0
    m.last_R_idx_B_offset = 0
  }
  m.chksum_diff_valid = false
}

//func (m *RLine) InsertR( R_num int, R rune ) {
//
//  R_size_in := utf8.RuneLen(R)
//  if( 0 < R_size_in ) {
//    B_offset_data := 0 // Byte offset in m.data
//    LEN := len(m.data)
//    for R_index:=0; B_offset_data<LEN; R_index++ {
//      _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//      if( R_num == R_index ) {
//        // Insert R into m.data at B_offset_data
//        m.PushR( R ) //< This will increase m.data cap if needed
//        copy( m.data[B_offset_data+R_size_in:], m.data[B_offset_data:] )
//        utf8.EncodeRune( m.data[B_offset_data:], R )
//
//        m.chksum_diff_valid = false
//        break
//      }
//      B_offset_data += R_size_data
//    }
//  }
//}

//func (m *RLine) InsertR( R_pos int, R rune ) {
//
//  R_size_in := utf8.RuneLen(R)
//  if( 0 < R_size_in ) {
//    LEN_B := m.LenB()
//    LEN_R := m.LenR()
//    if( LEN_R <= R_pos ) {
//      m.PushR( R ) //< This will increase m.data cap if needed
//      m.chksum_diff_valid = false
//    } else {
//      B_offset_data := 0 // Byte offset in m.data
//      for R_index:=0; B_offset_data<LEN_B; R_index++ {
//        _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//        if( R_pos == R_index ) {
//          // Insert R into m.data at B_offset_data
//          m.PushR( R ) //< This will increase m.data cap if needed
//          copy( m.data[B_offset_data+R_size_in:], m.data[B_offset_data:] )
//          utf8.EncodeRune( m.data[B_offset_data:], R )
//
//          m.chksum_diff_valid = false
//          break
//        }
//        B_offset_data += R_size_data
//      }
//    }
//  }
//}

// If R_pos is greater then
//func (m *RLine) InsertR( R_pos int, R rune ) {
//
//  R_size_in := utf8.RuneLen(R)
//  if( 0 < R_size_in ) {
//    inserted_R := false
//    B_offset_data := 0 // Byte offset in m.data
//    LEN := len(m.data)
//    // FIXME: use m.last_R_idx and m.last_R_idx_B_offset to speed up search:
//    for R_index:=0; B_offset_data<LEN; R_index++ {
//      _, R_size_data := utf8.DecodeRune( m.data[B_offset_data:] )
//      if( R_index == R_pos ) {
//        // Insert R into m.data at B_offset_data
//        m.PushR( R ) //< This will increase m.data cap if needed
//        copy( m.data[B_offset_data+R_size_in:], m.data[B_offset_data:] )
//        utf8.EncodeRune( m.data[B_offset_data:], R )
//
//        m.last_R_idx = R_index
//        m.last_R_idx_B_offset = B_offset_data
//        m.chksum_diff_valid = false
//        inserted_R = true
//        break
//      }
//      B_offset_data += R_size_data
//    }
//    if( !inserted_R ) {
//      // (LEN == 0) or (m.LenR() <= R_pos), so append R:
//      m.PushR( R ) //< This will increase m.data cap if needed
//      m.chksum_diff_valid = false
//    }
//  }
//}

// Insert rune R_in at R_pos.
// If (m.LenR() <= R_pos), append R_in.
//
func (m *RLine) InsertR( R_num int, R_in rune ) {

  _, R_size_data, B_offset_data := m.GetR( R_num ) 

  if( 0 < R_size_data ) {
    // A rune was found, so R_num is not past end of m.data.
    // Insert R_in into m.data at B_offset_data
    m.PushR( R_in ) //< This will increase m.data cap if needed
    R_size_in := utf8.RuneLen(R_in)
    copy( m.data[B_offset_data+R_size_in:], m.data[B_offset_data:] )
    utf8.EncodeRune( m.data[B_offset_data:], R_in )

    // m.last_R_idx and m.last_R_idx_B_offset were set in m.GetR()
    m.chksum_diff_valid = false

  } else {
    // A rune was NOT found, so R_num is past end of m.data.
    // Append R_in:
    LEN := len(m.data)
    if( LEN <= B_offset_data ) {
      // (LEN == 0) or (m.LenR() <= R_num), so append R_in:
      m.PushR( R_in ) //< This will increase m.data cap if needed

      // When R_in is appended, m.last_R_idx and m.last_R_idx_B_offset
      // dont need to change.
      m.chksum_diff_valid = false
    }
  }
}

func (m *RLine) EqualL( ln RLine ) bool {

//if( m.chksum_valid && ln.chksum_valid ) {
//  return m.chksum == ln.chksum
//}
  return slices.Equal( m.data, ln.data )
}

func (m *RLine) EqualDiffL( ln RLine ) bool {

  if( m.chksum_diff_valid && ln.chksum_diff_valid ) {
    return m.chksum_diff == ln.chksum_diff
  }
  return slices.Equal( m.data, ln.data )
}

func (m *RLine) EqualStr( S string ) bool {

  return slices.Equal( m.data, []byte(S) )
}

//func (m *RLine) StartsWith( S string ) bool {
//
//  return strings.HasPrefix( m.to_str(), S )
//}

//func (m *RLine) StartsWith( S string ) bool {
//
//  starts_with_S := true
//  len_self := len( m.data )
//  len_S := len( S )
//
//  if( len_self < len_S ) {
//    starts_with_S = false
//  } else {
//    for k:=0; (starts_with_S && k<len_S); k++ {
//      if( S[k] != m.data[k] ) {
//        starts_with_S = false
//      }
//    }
//  }
//  return starts_with_S
//}

func (m *RLine) StartsWith( S string ) bool {

  return m.has_at( S, 0 )
}

func (m *RLine) EndsWith( S string ) bool {

  ends_with_S := true
  len_self := len( m.data )
  len_S    := len( S )

  if( len_self < len_S ) {
    ends_with_S = false
  } else {
    size_diff := len_self - len_S
    for k:=0; (ends_with_S && k<len_S); k++ {
      if( S[k] != m.data[k+size_diff] ) {
        ends_with_S = false
      }
    }
  }
  return ends_with_S
}

func (m *RLine) to_str() string {
  return string(m.data)
}

//func (m *RLine) to_str( pos int ) string {
//  if( pos < len(m.data) ) {
//    return string(m.data[pos:])
//  }
//  return ""
//}

func (m *RLine) from_str( S string ) {
  m.data = []byte(S)

  m.last_R_idx = 0
  m.last_R_idx_B_offset = 0
  m.chksum_diff_valid = false
}

// This implementation avoids RLine.to_str(), which allocates a new string.
// Return -1 if m is less than    ln
// Return  0 if m is the same as  ln
// Return +1 if m is greater than ln
//
func (m *RLine) Compare( ln RLine ) int {
  rval := 0
  len_m  := len( m.data )
  len_ln := len( ln.data )
  min_len := Min_i( len_m, len_ln )

  for k:=0; (rval == 0 && k<min_len); k++ {
    if       ( m.data[k] < ln.data[k] ) { rval = -1
    } else if( m.data[k] > ln.data[k] ) { rval =  1
    }
  }
  if( 0 == rval ) {
    if       ( len_m < len_ln ) { rval = -1
    } else if( len_m > len_ln ) { rval =  1
    }
  }
  return rval
}

// Returns true if m.data has S at byte_pos
//
//func (m *RLine) has_at( S string, byte_pos int ) bool {
//  has := false
//  len_S := len(S)
//
//  if( byte_pos+len_S < len(m.data) ) {
//    has = true
//    for k:=0; has && k<len_S; k++ {
//      if( S[k] != m.data[byte_pos+k] ) {
//        has = false
//      }
//    }
//  }
//  return has
//}

// Returns true if m.data has S at byte_pos
//
func (m *RLine) has_at( S string, byte_pos int ) bool {

  has := true
  len_self := len( m.data )
  len_S := len( S )

  if( len_self < len_S+byte_pos ) {
    has = false
  } else {
    for k:=0; (has && k<len_S); k++ {
      if( S[k] != m.data[k+byte_pos] ) {
        has = false
      }
    }
  }
  return has
}

// Returns true if m.data has tag at byte_pos case insensitively
//
//func (m *RLine) has_at_ci( tag string, byte_pos int ) bool {
//  has := false
//  tag_len := len(tag)
//  if( byte_pos+tag_len < len(m.data) ) {
//    has = true
//    for k:=0; has && k<tag_len; k++ {
//      if( unicode.ToLower(rune(tag[k])) != unicode.ToLower(rune(m.data[byte_pos+k])) ) {
//        has = false
//      }
//    }
//  }
//  return has
//}

// Returns true if m.data has tag at byte_pos case insensitively
//
func (m *RLine) has_at_ci( S string, byte_pos int ) bool {

  has := true
  len_self := len( m.data )
  len_S := len( S )

  if( len_self < len_S+byte_pos ) {
    has = false
  } else {
    for k:=0; (has && k<len_S); k++ {
      if( unicode.ToLower(rune(S[k])) != unicode.ToLower(rune(m.data[k+byte_pos])) ) {
        has = false
      }
    }
  }
  return has
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

  LEN := len(m.data)
  for k:=0; k<LEN; k++ {

    if( IsSpace( rune(m.data[k]) ) ) {
      copy( m.data[k:], m.data[k+1:] )
      m.data = m.data[:len(m.data)-1]
      k--
      m.chksum_diff_valid = false
      LEN = len(m.data)

      m.last_R_idx = 0
      m.last_R_idx_B_offset = 0
    }
  }
}

//func (m *RLine) Chksum() uint32 {
//
//  if( !m.chksum_valid ) {
//    m.chksum = crc32.ChecksumIEEE( m.data )
//    m.chksum_valid = true
//  }
//  return m.chksum
//}

// Chksum_diff ignores white space at beginning and end of line:
//
func (m *RLine) Chksum_diff() uint32 {

  if( !m.chksum_diff_valid ) {
    st := 0
    fn := m.LenB()

    for( (0 < fn) && IsSpace( rune(m.data[fn-1]) ) ) { fn-- }
    for( (st < fn) && IsSpace( rune(m.data[st]) ) ) { st++ }

    m.chksum_diff = crc32.ChecksumIEEE( m.data[st:fn] )
    m.chksum_diff_valid = true
  }
  return m.chksum_diff
}

