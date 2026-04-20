
package main

import (
//"fmt"
)

type Highlight_JSON struct {
  p_fb  *FileBuf
  state HiStateFunc
  const_style byte
}

func (m *Highlight_JSON) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_JSON) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_JSON) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_JSON[:], st, fn )
}

var HiPairs_JSON = [...]HiKeyVal {

  { "true" , HI_DEFINE },
  { "false", HI_DEFINE },
  { "null" , HI_DEFINE },
  { ""     , 0 },
}

func (m *Highlight_JSON) Hi_In_None( l, p int ) (int,int) {
  m.const_style = HI_CONST
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c0=='{' || c0=='}' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )

      } else if( c0=='[' || c0==']' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )

      } else if( c0==':' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        m.state = m.Hi_Value
        p++
      } else if( c0==',' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )

      } else if( Quote_Start('"' ,c2,c1,c0) ) {
        m.state = m.Hi_In_DoubleQuote

      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
        m.state = m.Hi_NumberBeg

      } else if( c0 < 32 || 126 < c0 ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_JSON) Hi_In_DoubleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_DoubleQuote_Base( l,p, m.p_fb, m.Hi_In_None, m.const_style )

  return l,p
}

func (m *Highlight_JSON) Hi_NumberBeg( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberBeg_Base( l,p, m.p_fb, m.Hi_NumberIn,
                                                 m.Hi_NumberHex )
  return l,p
}

func (m *Highlight_JSON) Hi_NumberIn( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberIn_Base( l,p, m.p_fb, m.Hi_NumberIn,
                                                m.Hi_NumberFraction,
                                                m.Hi_NumberExponent,
                                                m.Hi_In_None )
  return l,p
}

func (m *Highlight_JSON) Hi_NumberHex( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberHex_Base( l,p, m.p_fb, m.Hi_NumberHex,
                                                 m.Hi_In_None )
  return l,p
}

func (m *Highlight_JSON) Hi_NumberFraction( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberFraction_Base( l,p, m.p_fb, m.Hi_NumberFraction,
                                                      m.Hi_NumberExponent,
                                                      m.Hi_In_None )
  return l,p
}

func (m *Highlight_JSON) Hi_NumberExponent( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberExponent_Base( l,p, m.p_fb, m.Hi_NumberExponent,
                                                      m.Hi_In_None )
  return l,p
}

//func (m *Highlight_JSON) Hi_Value( l, p int ) (int,int) {
//  m.const_style = HI_DEFINE
//  m.state = nil
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLen( l )
//
//    for ; p<LL; p++ {
//      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
//      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
//      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
//      var c0 rune =                     m.p_fb.GetR( l, p )
//
//      if( c0 == '{' || c0 == '}' || c0 == ',') {
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
//        m.state = m.Hi_In_None
//        p++
//      } else if( c0 == '[' || c0==']' ) {
//        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
//      //m.state = m.Hi_In_None
//      //m.const_style = HI_CONST
//      //p++
//      } else if( Quote_Start('"' ,c2,c1,c0) ) {
//        m.state = m.Hi_In_DoubleQuote
//
//      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
//        m.state = m.Hi_NumberBeg
//
//      } else if( c0 < 32 || 126 < c0 ) {
//        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//  }
//  return l,p
//}

func (m *Highlight_JSON) Hi_Value( l, p int ) (int,int) {
  m.const_style = HI_DEFINE
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c0 == '{' || c0 == '}' || c0 == ',') {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        m.state = m.Hi_In_None
        p++
      } else if( c0 == '[' || c0==']' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        m.state = m.Hi_In_None
        p++
      } else if( Quote_Start('"' ,c2,c1,c0) ) {
        m.state = m.Hi_In_DoubleQuote

      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
        m.state = m.Hi_NumberBeg

      } else if( c0 < 32 || 126 < c0 ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
  }
  return l,p
}

