
package main

import (
//"fmt"
)

type Highlight_Code struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_Code) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_Code) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if       ( c1=='/' && c0 == '/' ) { p--; m.state = m.Hi_CPP_Comment
      } else if( c1=='/' && c0 == '*' ) { p--; m.state = m.Hi_BegC_Comment
      } else if(            c0 == '#' ) { m.state = m.Hi_In_Define
      } else if( Quote_Start('\'',c2,c1,c0) ) { m.state = m.Hi_In_SingleQuote
      } else if( Quote_Start('"' ,c2,c1,c0) ) { m.state = m.Hi_In_DoubleQuote

      } else if( !IsIdent( c1 ) && IsDigit(c0) ) { m.state = m.Hi_NumberBeg
      } else if( (c1==':' && c0==':') || (c1=='-' && c0=='>') ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_VARTYPE )
        m.p_fb.SetSyntaxStyle( l, p  , HI_VARTYPE )
      } else if( TwoControl( c1, c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_CONTROL )
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONTROL )
      } else if( OneVarType( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      } else if( OneControl( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
      } else if( c0 < 32 || 126 < c0 ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Code) Hi_CPP_Comment( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  p=0; l++
  m.state = m.Hi_In_None

  return l,p
}

func (m *Highlight_Code) Hi_BegC_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p++
  m.state = m.Hi_In_C_Comment
  return l,p
}

func (m *Highlight_Code) Hi_In_C_Comment( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c1=='*' && c0=='/' ) {
        m.state = m.Hi_EndC_Comment
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Code) Hi_EndC_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p++
  m.state = m.Hi_In_None
  return l,p
}

func (m *Highlight_Code) Hi_In_Define( l, p int ) (int,int) {
  m.state = nil
  LL := m.p_fb.LineLen( l )

  var ce rune = 0; // character at end of line
  for ; p<LL; p++ {
    // c0 is ahead of c1: (c1,c0)
    var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
    var c0 rune =                     m.p_fb.GetR( l, p )

    if( c1=='/' && c0=='/' ) {
      m.state = m.Hi_CPP_Comment
      p--
    } else if( c1=='/' && c0=='*' ) {
      m.state = m.Hi_BegC_Comment
      p--
    } else {
      m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
    }
    if( nil != m.state ) { return l,p }
    ce = c0
  }
  p=0; l++

  if( ce == '\\' ) {
    m.state = m.Hi_In_Define
  } else {
    m.state = m.Hi_In_None
  }
  return l,p
}

//func (m *Highlight_Code) Hi_In_SingleQuote( l, p int ) (int,int) {
//  m.state = nil
//  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//  p++
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLen( l )
//
//    slash_escaped := false
//    for ; p<LL; p++ {
//      // c0 is ahead of c1: (c1,c0)
//      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
//      var c0 rune =                     m.p_fb.GetR( l, p )
//
//      if( (c1==0    && c0=='\'') ||
//          (c1!='\\' && c0=='\'') ||
//          (c1=='\\' && c0=='\'' && slash_escaped) ) {
//        // End of single quote:
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//        p++
//        m.state = m.Hi_In_None
//      } else {
//        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
//        } else                     { slash_escaped = false
//        }
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//    p = 0
//  }
//  return l,p
//}

// This shows one way to re-use class methods in Go:
//
func (m *Highlight_Code) Hi_In_SingleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_SingleQuote_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

//func (m *Highlight_Code) Hi_In_DoubleQuote( l, p int ) (int,int) {
//  m.state = nil
//  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//  p++
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLen( l )
//
//    slash_escaped := false
//    for ; p<LL; p++ {
//      // c0 is ahead of c1: (c1,c0)
//      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
//      var c0 rune =                     m.p_fb.GetR( l, p )
//
//      if( (c1==0    && c0=='"') ||
//          (c1!='\\' && c0=='"') ||
//          (c1=='\\' && c0=='"' && slash_escaped) ) {
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//        p++
//        m.state = m.Hi_In_None
//      } else {
//        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
//        } else                     { slash_escaped = false
//        }
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//    p = 0
//  }
//  return l,p
//}

// This shows one way to re-use class methods in Go:
//
func (m *Highlight_Code) Hi_In_DoubleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_DoubleQuote_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

func (m *Highlight_Code) Hi_NumberBeg( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

  var c1 rune = m.p_fb.GetR( l, p )
  p++
  m.state = m.Hi_NumberIn

  LL := m.p_fb.LineLen( l )
  if( '0' == c1 && (p+1)<LL ) {
    var c0 rune = m.p_fb.GetR( l, p )
    if( 'x' == c0 || 'X' == c0 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberHex
      p++
    }
  }
  return l,p
}

// Need to add highlighting for:
//   L = long
//   U = unsigned
//   UL = unsigned long
//   ULL = unsigned long long
//   F = float
// at the end of numbers
func (m *Highlight_Code) Hi_NumberIn( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = m.p_fb.GetR( l, p )

    if( '.'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberFraction
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberExponent
      p++
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p )
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
          p++
        }
      }
    } else if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else if( c1=='L' || c1=='F' || c1=='U' ) {
      m.state = m.Hi_NumberTypeSpec
    } else if( c1=='\'' && (p+1)<LL ) {
      // ' is followed by another digit on line
      var c0 rune = m.p_fb.GetR( l, p+1 )

      if( IsDigit( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONST )
        m.p_fb.SetSyntaxStyle( l, p+1, HI_CONST )
        p += 2
      } else {
        m.state = m.Hi_In_None
      }
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Code) Hi_NumberHex( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsXDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Code) Hi_NumberFraction( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberExponent
      p++
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p )
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
          p++
        }
      }
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Code) Hi_NumberExponent( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Code) Hi_NumberTypeSpec( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l )

  if( p < LL ) {
    var c0 rune = m.p_fb.GetR( l, p )

    if( c0=='L' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      p++
      m.state = m.Hi_In_None
    } else if( c0=='F' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      p++
      m.state = m.Hi_In_None
    } else if( c0=='U' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
      if( p<LL ) {
        var c1 rune = m.p_fb.GetR( l, p )
        if( c1=='L' ) { // UL
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
          if( p<LL ) {
            var c2 rune = m.p_fb.GetR( l, p )
            if( c2=='L' ) { // ULL
              m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
            }
          }
        }
      }
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

