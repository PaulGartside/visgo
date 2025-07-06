
package main

import (
//"fmt"
)

type Highlight_XML struct {
  p_fb  *FileBuf
  state HiStateFunc
  qtXSt HiStateFunc // qtXSt = Quote Exit State
}

func (m *Highlight_XML) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_XML) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_XML) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_XML[:], st, fn )
}

var HiPairs_XML = [...]HiKeyVal {
  // HTML tags:
  { "xml"     , HI_CONTROL },
  { "version" , HI_CONTROL },
  { "encoding", HI_CONTROL },

  { ""        , 0 },
}

func (m *Highlight_XML) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c3,c2,c1,c0)
      var c3 rune = 0; if( 2<p ) { c3 = m.p_fb.GetR( l, p-3 ) }
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c1=='<' && c0!='!' && c0!='/') {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE )
        m.state = m.Hi_OpenTag_ElemName

      } else if( c1=='<' && c0=='/') {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE ); //< '<'
        m.p_fb.SetSyntaxStyle( l, p  , HI_DEFINE ); //< '/'
        p++; // Move past '/'
        m.state = m.Hi_OpenTag_ElemName

      } else if( c3=='<' && c2=='!' && c1=='-' && c0=='-' ) {
        m.p_fb.SetSyntaxStyle( l, p-3, HI_COMMENT ); //< '<'
        m.p_fb.SetSyntaxStyle( l, p-2, HI_COMMENT ); //< '!'
        m.p_fb.SetSyntaxStyle( l, p-1, HI_COMMENT ); //< '-'
        m.p_fb.SetSyntaxStyle( l, p  , HI_COMMENT ); //< '-'
        p++; // Move past '-'
        m.state = m.Hi_Comment

      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
        m.state = m.Hi_NumberBeg

      } else {
        ; //< No syntax highlighting on content outside of <>tags
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_OpenTag_ElemName( l, p int ) (int,int) {
  m.state = nil
  found_elem_name := false

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++; // Move past '>'
        m.state = m.Hi_In_None

      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( !found_elem_name ) {
        if( IsXML_Ident( c0 ) ) {
          found_elem_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else if( c0==' ' || c0=='\t' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else if( found_elem_name ) {
        if( IsXML_Ident( c0 ) ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else if( c0==' ' || c0=='\t' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
          p++; //< Move past white space
          m.state = m.Hi_OpenTag_AttrName
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_OpenTag_AttrName( l, p int ) (int,int) {
  m.state = nil
  found_attr_name := false
  past__attr_name := false

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++; // Move past '>'
        m.state = m.Hi_In_None

      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( !found_attr_name ) {
        if( IsXML_Ident( c0 ) ) {
          found_attr_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else if( c0==' ' || c0=='\t' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else if( found_attr_name && !past__attr_name ) {
        if( IsXML_Ident( c0 ) ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else if( c0==' ' || c0=='\t' ) {
          past__attr_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else if( c0=='=' ) {
          past__attr_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
          p++; //< Move past '='
          m.state = m.Hi_OpenTag_AttrVal
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else if( found_attr_name && past__attr_name ) {
        if( c0=='=' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
          p++; //< Move past '='
          m.state = m.Hi_OpenTag_AttrVal
        } else if( c0==' ' || c0=='\t' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_OpenTag_AttrVal( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++; // Move past '>'
        m.state = m.Hi_In_None
      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
      } else if( c0=='\'' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++; // Move past '\''
        m.state = m.Hi_In_SingleQuote
        m.qtXSt = m.Hi_OpenTag_AttrName
      } else if( c0=='"' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++; //< Move past '\"'
        m.state = m.Hi_In_DoubleQuote
        m.qtXSt = m.Hi_OpenTag_AttrName
      } else if( c0==' ' || c0=='\t' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_CloseTag( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++; // Move past '>'
        m.state = m.Hi_In_None

      } else if( c0=='/' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( IsXML_Ident( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )

      } else if( c0==' ' || c0=='\t' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_Comment( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c2=='-' && c1=='-' && c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p-2, HI_COMMENT ); //< '-'
        m.p_fb.SetSyntaxStyle( l, p-1, HI_COMMENT ); //< '-'
        m.p_fb.SetSyntaxStyle( l, p  , HI_COMMENT ); //< '>'
        p++; // Move past '>'
        m.state = m.Hi_In_None
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_In_SingleQuote( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                    m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='\'') ||
          (c1!='\\' && c0=='\'') ||
          (c1=='\\' && c0=='\'' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++; // Move past '\''
        m.state = m.qtXSt
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else {                     slash_escaped = false
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_XML) Hi_In_DoubleQuote( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                    m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='"') ||
          (c1!='\\' && c0=='"') ||
          (c1=='\\' && c0=='"' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++; //< Move past '\"'
        m.state = m.qtXSt
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else {                     slash_escaped = false
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

//func (m *Highlight_XML) Hi_NumberBeg( l, p int ) (int,int) {
//  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//
//  var c1 rune = m.p_fb.GetR( l, p )
//  p++
//  m.state = m.Hi_NumberIn
//
//  LL := m.p_fb.LineLen( l )
//  if( '0' == c1 && (p+1)<LL ) {
//    var c0 rune = m.p_fb.GetR( l, p )
//    if( 'x' == c0 || 'X' == c0 ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      m.state = m.Hi_NumberHex
//      p++
//    }
//  }
//  return l,p
//}

func (m *Highlight_XML) Hi_NumberBeg( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberBeg_Base( l,p, m.p_fb, m.Hi_NumberIn, m.Hi_NumberHex )

  return l,p
}

//func (m *Highlight_XML) Hi_NumberIn( l, p int ) (int,int) {
//  LL := m.p_fb.LineLen( l )
//  if( LL <= p ) { m.state = m.Hi_In_None
//  } else {
//    var c1 rune = m.p_fb.GetR( l, p )
//
//    if( '.'==c1 ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      m.state = m.Hi_NumberFraction
//      p++
//    } else if( 'e'==c1 || 'E'==c1 ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      m.state = m.Hi_NumberExponent
//      p++
//      if( p<LL ) {
//        var c0 rune = m.p_fb.GetR( l, p )
//        if( '+' == c0 || '-' == c0 ) {
//          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//          p++
//        }
//      }
//    } else if( IsDigit(c1) ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      p++
//    } else {
//      m.state = m.Hi_In_None
//    }
//  }
//  return l,p
//}

func (m *Highlight_XML) Hi_NumberIn( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberIn_Base( l,p, m.p_fb, m.Hi_In_None, m.Hi_NumberFraction, m.Hi_NumberExponent )

  return l,p
}

//func (m *Highlight_XML) Hi_NumberHex( l, p int ) (int,int) {
//  LL := m.p_fb.LineLen( l )
//  if( LL <= p ) { m.state = m.Hi_In_None
//  } else {
//    var c1 rune = m.p_fb.GetR( l, p )
//    if( IsXDigit(c1) ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      p++
//    } else {
//      m.state = m.Hi_In_None
//    }
//  }
//  return l,p
//}

func (m *Highlight_XML) Hi_NumberHex( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberHex_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

//func (m *Highlight_XML) Hi_NumberFraction( l, p int ) (int,int) {
//  LL := m.p_fb.LineLen( l )
//  if( LL <= p ) { m.state = m.Hi_In_None
//  } else {
//    var c1 rune = m.p_fb.GetR( l, p )
//    if( IsDigit(c1) ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      p++
//    } else if( 'e'==c1 || 'E'==c1 ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      m.state = m.Hi_NumberExponent
//      p++
//      if( p<LL ) {
//        var c0 rune = m.p_fb.GetR( l, p )
//        if( '+' == c0 || '-' == c0 ) {
//          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//          p++
//        }
//      }
//    } else {
//      m.state = m.Hi_In_None
//    }
//  }
//  return l,p
//}

func (m *Highlight_XML) Hi_NumberFraction( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberFraction_Base( l,p, m.p_fb, m.Hi_In_None, m.Hi_NumberExponent )

  return l,p
}

//func (m *Highlight_XML) Hi_NumberExponent( l, p int ) (int,int) {
//  LL := m.p_fb.LineLen( l )
//  if( LL <= p ) { m.state = m.Hi_In_None
//  } else {
//    var c1 rune = m.p_fb.GetR( l, p )
//    if( IsDigit(c1) ) {
//      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      p++
//    } else {
//      m.state = m.Hi_In_None
//    }
//  }
//  return l,p
//}

func (m *Highlight_XML) Hi_NumberExponent( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberExponent_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

