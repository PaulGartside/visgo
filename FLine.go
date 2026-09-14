
package main

import (
//"bytes"
//"fmt"
//"slices"
//"hash/alder32"
//"hash/crc32"
  "unicode/utf8"
)

// File line
type FLine struct {
  // runes and styles should be the same length
  runes  RLine // rune Line
  styles BLine // byte Line
  star_styles_valid bool // star styles is short for regex search styles
}

func (m *FLine) Init( length int ) {
  m.runes.Init( length )
  m.styles.Init( length )
}

// Set length to zero
func (m *FLine) Clear() {
  m.runes.Clear()
  m.styles.Clear()
  m.star_styles_valid = false
}

// Copy *p_src_ln into self
func (m *FLine) CopyP( p_src_ln *FLine ) {
  m.runes.Copy( p_src_ln.runes )
  m.styles.Copy( p_src_ln.styles )
  m.star_styles_valid = false
}

// Copy rune line pointer
// Copy *p_src_ln into self
//
func (m *FLine) CopyPRL( p_src_ln *RLine ) {
  m.runes.Copy( *p_src_ln )
  m.styles.SetLen( m.runes.LenB() )
  m.styles.Zeroize()
  m.star_styles_valid = false
}

//func (m *FLine) Size() int {
//  return m.runes.Size()
//}

func (m *FLine) LenB() int {
  return m.runes.LenB()
}

func (m *FLine) LenR() int {
  return m.runes.LenR()
}

// Get byte
func (m *FLine) GetB( idx int ) byte {
  return m.runes.GetB( idx )
}

// Get rune
func (m *FLine) GetR( R_num int ) (rune, int, int) {
  return m.runes.GetR( R_num )
}

func (m *FLine) GetRatB( B_num int ) (rune, int) {
  return m.runes.GetRatB( B_num )
}

// Get style
func (m *FLine) GetStyle( idx int ) byte {
  return m.styles.GetB( idx )
}

func (m *FLine) SetB( idx int, B byte ) {
  m.runes.SetB( idx, B )
  m.star_styles_valid = false
}

// Sets idx rune in m.runes to R.
// Returns byte offset in m.runes of that rune
//
func (m *FLine) SetR( idx int, R rune ) int {
  B_pos := m.runes.SetR( idx, R )
  m.star_styles_valid = false
  return B_pos
}

//func (m *FLine) SetRatB( ) {
//  m.runes.SetRatB()
//  m.star_styles_valid = false
//}

// Set style
func (m *FLine) SetStyle( idx int, S byte ) {
  m.styles.SetB( idx, S )
}

func (m *FLine) RemoveB( idx int ) byte {

  var B byte = m.runes.RemoveB( idx )
               m.styles.RemoveB( idx )
  m.star_styles_valid = false
  return B
}

// Remove rune
func (m *FLine) RemoveR( idx int ) rune {

  var R rune = m.runes.RemoveR( idx )

  R_size := utf8.RuneLen(R)
  for k:=0; k<R_size; k++ {
    m.styles.RemoveB( idx )
  }
  m.star_styles_valid = false
  return R
}

// Push byte
func (m *FLine) PushB( B byte ) {

  m.runes.PushB( B )
  m.styles.PushB( 0 )
  m.star_styles_valid = false
} 

// Push rune
func (m *FLine) PushR( R rune ) {

  m.runes.PushR( R )

  R_size := utf8.RuneLen(R)
  for k:=0; k<R_size; k++ {
    m.styles.PushB( 0 )
  }
  m.star_styles_valid = false
} 

// Push a slice of runes(s_r)
//func (m *FLine) PushSR( s_r []rune ) {
//
//  s_b := make( []byte, len(s_r) )
//
//  m.runes.PushSR( s_r )
//  m.styles.PushSB( s_b )
//  m.star_styles_valid = false
//} 

func (m *FLine) PushStr( S string ) {

  m.runes.PushStr( S )

  len_S := len(S)
  for k:=0; k<len_S; k++ {
    m.styles.PushB( 0 )
  }
  m.star_styles_valid = false
}

// Push line
func (m *FLine) PushL( ln FLine ) {
  m.runes.PushL( ln.runes )
  m.styles.PushL( ln.styles ) // Should we push the styles ?
  m.star_styles_valid = false
}

func (m *FLine) PushLP( p_fl *FLine ) {
  m.runes.PushL( p_fl.runes )
  m.styles.PushL( p_fl.styles ) // Should we push the styles ?
  m.star_styles_valid = false
}

func (m *FLine) InsertB( idx int, B byte ) {

  m.runes.InsertB( idx, B )
  m.styles.InsertB( idx, 0 )
  m.star_styles_valid = false
}

// Insert rune
func (m *FLine) InsertR( idx int, R rune ) {

  m.runes.InsertR( idx, R )

  R_size := utf8.RuneLen(R)
  for k:=0; k<R_size; k++ {
    m.styles.InsertB( idx, 0 )
  }
  m.star_styles_valid = false
}

// Equal line pointer
func (m *FLine) EqualLP( pln *FLine ) bool {

  return m.runes.EqualL( pln.runes )
}

// Equal line pointer
func (m *FLine) EqualDiffLP( pln *FLine ) bool {

  return m.runes.EqualDiffL( pln.runes )
}

// Equal string
func (m *FLine) EqualStr( S string ) bool {

  return m.runes.EqualStr( S )
}

func (m *FLine) from_str( S string ) {
  m.runes.from_str( S )
  m.styles.Clear()
  m.styles.SetLen( m.runes.LenB() )
}

// Convert to string
func (m *FLine) to_str() string {

  return m.runes.to_str()
}

func (m *FLine) Compare( pln *FLine ) int {

  return m.runes.Compare( pln.runes )
}

// Convert to slice of bytes
func (m *FLine) to_SB( st int ) []byte {
  return m.runes.to_SB( st )
}

// Returns true if m.runes has tag at byte_pos
//
func (m *FLine) has_at( tag string, byte_pos int ) bool {

  return m.runes.has_at( tag, byte_pos )
}

// Returns true if m.runes has tag at byte_pos case insensitively
//
func (m *FLine) has_at_ci( tag string, byte_pos int ) bool {

  return m.runes.has_at_ci( tag, byte_pos )
}

//func (m *FLine) RemoveSpaces() {
//
//  for k:=0; k<len(m.data); k++ {
//
//    if( IsSpace( m.data[k] ) ) {
//      copy( m.data[k:], m.data[k+1:] )
//      m.data = m.data[:len(m.data)-1]
//      k--
//    }
//  }
//}

func (m *FLine) ClearStarAndInFileStyles() {

  for k:=0; k<m.styles.Len(); k++ {

    var S byte = m.styles.GetB( k )

    S &^= HI_STAR      // Clear HI_STAR(0x01) in S
    S &^= HI_STAR_IN_F // Clear HI_STAR_IN_F(0x02) in S

    m.styles.SetB( k, S )
  }
  m.star_styles_valid = false
}

// Leave syntax m.styles unchanged, and set star-in-file style
func (m *FLine) Set__StarInFStyle( idx int ) {

  m.styles.SetB( idx, m.styles.GetB( idx ) | HI_STAR_IN_F )
}

func (m *FLine) EndsWith( suffix string ) bool {

  return m.runes.EndsWith( suffix )
}

//func (m *FLine) Chksum() uint32 {
//
//  return m.runes.Chksum()
//}

func (m *FLine) Chksum_diff() uint32 {

  return m.runes.Chksum_diff()
}

