
package main

import (
//"fmt"
)

type Highlight_Rust struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_Rust) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_Rust) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_Rust) Run_Range( st CrsPos, fn int )")

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_Rust) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_Rust[:], st, fn )
}

func (m *Highlight_Rust) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLenB( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = rune(m.p_fb.GetB( l, p-2 )) }
      var c1 rune = 0; if( 0<p ) { c1 = rune(m.p_fb.GetB( l, p-1 )) }
      var c0 rune =                     rune(m.p_fb.GetB( l, p ))

      if       ( c1=='/' && c0 == '/' ) { p--; m.state = m.Hi_BegCPP_Comment
    //} else if( c1=='/' && c0 == '*' ) { p--; m.state = m.Hi_BegC_Comment

      } else if( Quote_Start('\'',c2,c1,c0) ) { m.state = m.Hi_In_SingleQuote
      } else if( Quote_Start('"' ,c2,c1,c0) ) { m.state = m.Hi_In_DoubleQuote
      } else if( Quote_Start('`' ,c2,c1,c0) ) { m.state = m.Hi_In_Back__Quote

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

//func (m *Highlight_Rust) Hi_BegC_Comment( l, p int ) (int,int) {
//  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
//  p++
//  m.state = m.Hi_In_C_Comment
//  return l,p
//}

//func (m *Highlight_Rust) Hi_In_C_Comment( l, p int ) (int,int) {
//  m.state = nil
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLenB( l )
//
//    for ; p<LL; p++ {
//      // c0 is ahead of c1: (c1,c0)
//      var c1 rune = 0; if( 0<p ) { c1 = rune(m.p_fb.GetB( l, p-1 )) }
//      var c0 rune =                     rune(m.p_fb.GetB( l, p ))
//
//      if( c1=='*' && c0=='/' ) {
//        m.state = m.Hi_EndC_Comment
//      } else {
//        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//    p = 0
//  }
//  return l,p
//}

//func (m *Highlight_Rust) Hi_EndC_Comment( l, p int ) (int,int) {
//  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
//  p++
//  m.state = m.Hi_In_None
//  return l,p
//}

func (m *Highlight_Rust) Hi_BegCPP_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p++
  m.state = m.Hi_In_CPP_Comment
  return l,p
}

func (m *Highlight_Rust) Hi_In_CPP_Comment( l, p int ) (int,int) {
  var LL int = m.p_fb.LineLenB( l )
  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  p--
  m.state = m.Hi_EndCPP_Comment
  return l,p
}

func (m *Highlight_Rust) Hi_EndCPP_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p=0; l++
  m.state = m.Hi_In_None
  return l,p
}

//func (m *Highlight_Rust) Hi_In_SingleQuote( l, p int ) (int,int) {
//  m.state = nil
//  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//  p++
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLen( l )
//
//    var slash_escaped bool = false
//    for ; p<LL; p++ {
//      // c0 is ahead of c1: (c1,c0)
//      var c1 rune = 0; if( 0<p) { c1 = m.p_fb.GetR( l, p-1 ) }
//      var c0 rune =                    m.p_fb.GetR( l, p )
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
//        } else {                     slash_escaped = false
//        }
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//    p = 0
//  }
//  return l,p
//}

// This shows one way to re-use class methods in Rust:
//
func (m *Highlight_Rust) Hi_In_SingleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_SingleQuote_Base( l,p, m.p_fb, m.Hi_In_None, HI_CONST )

  return l,p
}

//func (m *Highlight_Rust) Hi_In_DoubleQuote( l, p int ) (int,int) {
//  m.state = nil
//  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//  p++
//  for ; l<m.p_fb.NumLines(); l++ {
//    LL := m.p_fb.LineLen( l )
//
//    var slash_escaped bool = false
//    for ; p<LL; p++ {
//      // c0 is ahead of c1: (c1,c0)
//      var c1 rune = 0; if( 0<p) { c1 = m.p_fb.GetR( l, p-1 ) }
//      var c0 rune =                    m.p_fb.GetR( l, p )
//
//      if( (c1==0    && c0=='"') ||
//          (c1!='\\' && c0=='"') ||
//          (c1=='\\' && c0=='"' && slash_escaped) ) {
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//        p++
//        m.state = m.Hi_In_None
//      } else {
//        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
//        } else {                     slash_escaped = false
//        }
//        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
//      }
//      if( nil != m.state ) { return l,p }
//    }
//    p = 0
//  }
//  return l,p
//}

// This shows one way to re-use class methods in Rust:
//
func (m *Highlight_Rust) Hi_In_DoubleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_DoubleQuote_Base( l,p, m.p_fb, m.Hi_In_None, HI_CONST )

  return l,p
}

func (m *Highlight_Rust) Hi_In_Back__Quote( l, p int ) (int,int) {
  m.state = nil
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
  p++
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLenB( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c0 rune = rune(m.p_fb.GetB( l, p ))

      if( c0=='`' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++
        m.state = m.Hi_In_None
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Rust) Hi_NumberBeg( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

  var c1 rune = rune(m.p_fb.GetB( l, p ))
  p++
  m.state = m.Hi_NumberIn

  LL := m.p_fb.LineLenB( l )

  if( '0' == c1 && (p+1)<LL ) {
    var c0 rune = rune(m.p_fb.GetB( l, p ))

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
func (m *Highlight_Rust) Hi_NumberIn( l, p int ) (int,int) {
  LL := m.p_fb.LineLenB( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = rune(m.p_fb.GetB( l, p ))

    if( '.'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberFraction
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberExponent
      p++
      if( p<LL ) {
        var c0 rune = rune(m.p_fb.GetB( l, p ))
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
      var c0 rune = rune(m.p_fb.GetB( l, p+1 ))

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

func (m *Highlight_Rust) Hi_NumberHex( l, p int ) (int,int) {
  LL := m.p_fb.LineLenB( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = rune(m.p_fb.GetB( l, p ))
    if( IsXDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Rust) Hi_NumberFraction( l, p int ) (int,int) {
  LL := m.p_fb.LineLenB( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = rune(m.p_fb.GetB( l, p ))
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberExponent
      p++
      if( p<LL ) {
        var c0 rune = rune(m.p_fb.GetB( l, p ))
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

func (m *Highlight_Rust) Hi_NumberExponent( l, p int ) (int,int) {
  LL := m.p_fb.LineLenB( l )
  if( LL <= p ) { m.state = m.Hi_In_None
  } else {
    var c1 rune = rune(m.p_fb.GetB( l, p ))
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Rust) Hi_NumberTypeSpec( l, p int ) (int,int) {
  LL := m.p_fb.LineLenB( l )

  if( p < LL ) {
    var c0 rune = rune(m.p_fb.GetB( l, p ))

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
        var c1 rune = rune(m.p_fb.GetB( l, p ))
        if( c1=='L' ) { // UL
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
          if( p<LL ) {
            var c2 rune = rune(m.p_fb.GetB( l, p ))
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

var HiPairs_Rust = [...]HiKeyVal {
  // Keywords:
  { "as"                 , HI_CONTROL },
  { "async"              , HI_CONTROL },
  { "await"              , HI_CONTROL },
  { "break"              , HI_CONTROL },
  { "const"              , HI_CONTROL },
  { "continue"           , HI_CONTROL },
  { "crate"              , HI_CONTROL },
  { "dyn"                , HI_CONTROL },
  { "else"               , HI_CONTROL },
  { "enum"               , HI_CONTROL },
  { "extern"             , HI_CONTROL },
  { "false"              , HI_VARTYPE },
  { "fn"                 , HI_CONTROL },
  { "for"                , HI_CONTROL },
  { "if"                 , HI_CONTROL },
  { "impl"               , HI_CONTROL },
  { "in"                 , HI_CONTROL },
  { "let"                , HI_CONTROL },
  { "loop"               , HI_CONTROL },
  { "match"              , HI_CONTROL },
  { "mod"                , HI_CONTROL },
  { "move"               , HI_CONTROL },
  { "mut"                , HI_CONTROL },
  { "pub"                , HI_CONTROL },
  { "ref"                , HI_CONTROL },
  { "return"             , HI_CONTROL },
  { "Self"               , HI_CONTROL },
  { "self"               , HI_CONTROL },
  { "static"             , HI_CONTROL },
  { "struct"             , HI_CONTROL },
  { "super"              , HI_CONTROL },
  { "trait"              , HI_CONTROL },
  { "true"               , HI_VARTYPE },
  { "type"               , HI_CONTROL },
  { "union"              , HI_CONTROL },
  { "unsafe"             , HI_CONTROL },
  { "use"                , HI_CONTROL },
  { "where"              , HI_CONTROL },
  { "while"              , HI_CONTROL },

  // Keywords reserved for future use:
  { "abstract"           , HI_CONTROL },
  { "become"             , HI_CONTROL },
  { "box"                , HI_CONTROL },
  { "do"                 , HI_CONTROL },
  { "final"              , HI_CONTROL },
  { "gen"                , HI_CONTROL },
  { "macro"              , HI_CONTROL },
  { "override"           , HI_CONTROL },
  { "priv"               , HI_CONTROL },
  { "try"                , HI_CONTROL },
  { "typeof"             , HI_CONTROL },
  { "unsized"            , HI_CONTROL },
  { "virtual"            , HI_CONTROL },
  { "yield"              , HI_CONTROL },

  // Types
  { "i8"                 , HI_VARTYPE },
  { "i16"                , HI_VARTYPE },
  { "i32"                , HI_VARTYPE },
  { "i64"                , HI_VARTYPE },
  { "i128"               , HI_VARTYPE },
  { "isize"              , HI_VARTYPE },
  { "u8"                 , HI_VARTYPE },
  { "u16"                , HI_VARTYPE },
  { "u32"                , HI_VARTYPE },
  { "u64"                , HI_VARTYPE },
  { "u128"               , HI_VARTYPE },
  { "usize"              , HI_VARTYPE },
  { "f32"                , HI_VARTYPE },
  { "f64"                , HI_VARTYPE },
  { "bool"               , HI_VARTYPE },
  { "char"               , HI_VARTYPE },
  { "tup"                , HI_VARTYPE },
  { "String"             , HI_VARTYPE }, // String is a standard struct, not a built in type

  // Constants
//{ ""              , HI_CONST   },

//{ ""             , HI_DEFINE  },
  { ""                   , 0 },
}

