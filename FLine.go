
package main

import (
//"bytes"
//"fmt"
//"slices"
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
}

// Copy *p_src_ln into self
func (m *FLine) CopyPRL( p_src_ln *RLine ) {
  m.runes.Copy( *p_src_ln )
  m.styles.SetLen( m.runes.Len() )
  m.styles.Zeroize()
  m.star_styles_valid = false
}

func (m *FLine) Len() int {
  return m.runes.Len()
}

// Get rune
func (m *FLine) GetR( idx int ) rune {
  return m.runes.GetR( idx )
}

// Get style
func (m *FLine) GetStyle( idx int ) byte {
  return m.styles.GetB( idx )
}

// Set rune
func (m *FLine) SetR( idx int, R rune ) {
  m.runes.SetR( idx, R )
  m.star_styles_valid = false
}

// Set style
func (m *FLine) SetStyle( idx int, S byte ) {
  m.styles.SetB( idx, S )
}

// Remove rune
func (m *FLine) RemoveR( idx int ) rune {

  var R rune = m.runes.RemoveR( idx )
               m.styles.RemoveB( idx )
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
  m.styles.PushB( 0 )
  m.star_styles_valid = false
} 

// Push a slice of runes(s_r)
func (m *FLine) PushSR( s_r []rune ) {

  s_b := make( []byte, len(s_r) )

  m.runes.PushSR( s_r )
  m.styles.PushSB( s_b )
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

// Insert rune
func (m *FLine) InsertR( idx int, R rune ) {

  m.runes.InsertR( idx, R )
  m.styles.InsertB( idx, 0 )
  m.star_styles_valid = false
}

// Equal line pointer
func (m *FLine) EqualLP( pln *FLine ) bool {

  return m.runes.EqualL( pln.runes )
}

// Equal string
func (m *FLine) EqualStr( S string ) bool {

  return m.runes.EqualStr( S )
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

//func (m *FLine) ClearStarAndInFileStyles() {
//
//  for k:=0; k<styles.Len(); k++ {
//    var S byte = styles.data[k]
//
//    S &= ~HI_STAR
//    S &= ~HI_STAR_IN_F
//
//    styles.data[k] = S
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

func (m *FLine) ends_with( suffix string ) bool {

  return m.runes.ends_with( suffix )
}

