
package main

import (
//"fmt"
)

type Highlight_CPP struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_CPP) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_CPP) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_Go) Run_Range( st CrsPos, fn int )")

  m.state = m.Hi_In_None

  l := st.crsLine;
  p := st.crsChar;

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_CPP) Hi_In_None( l, p int ) (int,int) {
  m.state = nil;
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l );

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p );

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if       ( c1=='/' && c0 == '/' ) { p--; m.state = m.Hi_CPP_Comment;
      } else if( c1=='/' && c0 == '*' ) { p--; m.state = m.Hi_BegC_Comment;
      } else if(            c0 == '#' ) { m.state = m.Hi_In_Define;
      } else if( Quote_Start('\'',c2,c1,c0) ) { m.state = m.Hi_In_SingleQuote;
      } else if( Quote_Start('"' ,c2,c1,c0) ) { m.state = m.Hi_In_DoubleQuote;

      } else if( !IsIdent( c1 ) && IsDigit(c0) ) { m.state = m.Hi_NumberBeg;
      } else if( (c1==':' && c0==':') || (c1=='-' && c0=='>') ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_VARTYPE );
        m.p_fb.SetSyntaxStyle( l, p  , HI_VARTYPE );
      } else if( TwoControl( c1, c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_CONTROL );
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONTROL );
      } else if( OneVarType( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE );
      } else if( OneControl( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL );
      } else if( c0 < 32 || 126 < c0 ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII );
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0;
  }
  return l,p
}

func (m *Highlight_CPP) Hi_CPP_Comment( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT );
  }
  p=0; l++;
  m.state = m.Hi_In_None;

  return l,p
}

func (m *Highlight_CPP) Hi_BegC_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT );
  p++;
  m.state = m.Hi_In_C_Comment;
  return l,p
}

func (m *Highlight_CPP) Hi_In_C_Comment( l, p int ) (int,int) {
  m.state = nil;
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l );

    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c1=='*' && c0=='/' ) {
        m.state = m.Hi_EndC_Comment;
      } else {
        m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT );
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0;
  }
  return l,p
}

func (m *Highlight_CPP) Hi_EndC_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT );
  p++;
  m.state = m.Hi_In_None;
  return l,p
}

func (m *Highlight_CPP) Hi_In_Define( l, p int ) (int,int) {
  m.state = nil;
  LL := m.p_fb.LineLen( l );

  var ce rune = 0; // character at end of line
  for ; p<LL; p++ {
    // c0 is ahead of c1: (c1,c0)
    var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
    var c0 rune =                     m.p_fb.GetR( l, p )

    if( c1=='/' && c0=='/' ) {
      m.state = m.Hi_CPP_Comment;
      p--;
    } else if( c1=='/' && c0=='*' ) {
      m.state = m.Hi_BegC_Comment;
      p--;
    } else {
      m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE );
    }
    if( nil != m.state ) { return l,p }
    ce = c0;
  }
  p=0; l++;

  if( ce == '\\' ) {
    m.state = m.Hi_In_Define;
  } else {
    m.state = m.Hi_In_None;
  }
  return l,p
}

func (m *Highlight_CPP) Hi_In_SingleQuote( l, p int ) (int,int) {
  m.state = nil;
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
  p++;
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l );

    slash_escaped := false;
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='\'') ||
          (c1!='\\' && c0=='\'') ||
          (c1=='\\' && c0=='\'' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
        p++;
        m.state = m.Hi_In_None;
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped;
        } else                     { slash_escaped = false;
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0;
  }
  return l,p
}

func (m *Highlight_CPP) Hi_In_DoubleQuote( l, p int ) (int,int) {
  m.state = nil;
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
  p++;
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l );

    slash_escaped := false;
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='"') ||
          (c1!='\\' && c0=='"') ||
          (c1=='\\' && c0=='"' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
        p++;
        m.state = m.Hi_In_None;
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped;
        } else                     { slash_escaped = false;
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0;
  }
  return l,p
}

func (m *Highlight_CPP) Hi_NumberBeg( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST );

  var c1 rune = m.p_fb.GetR( l, p );
  p++;
  m.state = m.Hi_NumberIn;

  LL := m.p_fb.LineLen( l );

  if( '0' == c1 && (p+1)<LL ) {
    var c0 rune = m.p_fb.GetR( l, p );

    if( 'x' == c0 || 'X' == c0 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      m.state = m.Hi_NumberHex;
      p++;
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
func (m *Highlight_CPP) Hi_NumberIn( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );
  if( LL <= p ) { m.state = m.Hi_In_None;
  } else {
    var c1 rune = m.p_fb.GetR( l, p );

    if( '.'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      m.state = m.Hi_NumberFraction;
      p++;
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      m.state = m.Hi_NumberExponent;
      p++;
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p );
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
          p++;
        }
      }
    } else if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      p++;
    } else if( c1=='L' || c1=='F' || c1=='U' ) {
      m.state = m.Hi_NumberTypeSpec;
    } else if( c1=='\'' && (p+1)<LL ) {
      // ' is followed by another digit on line
      var c0 rune = m.p_fb.GetR( l, p+1 );

      if( IsDigit( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONST );
        m.p_fb.SetSyntaxStyle( l, p+1, HI_CONST );
        p += 2;
      } else {
        m.state = m.Hi_In_None;
      }
    } else {
      m.state = m.Hi_In_None;
    }
  }
  return l,p
}

func (m *Highlight_CPP) Hi_NumberHex( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );
  if( LL <= p ) { m.state = m.Hi_In_None;
  } else {
    var c1 rune = m.p_fb.GetR( l, p );
    if( IsXDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      p++;
    } else {
      m.state = m.Hi_In_None;
    }
  }
  return l,p
}

func (m *Highlight_CPP) Hi_NumberFraction( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );
  if( LL <= p ) { m.state = m.Hi_In_None;
  } else {
    var c1 rune = m.p_fb.GetR( l, p );
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      p++;
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      m.state = m.Hi_NumberExponent;
      p++;
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p );
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
          p++;
        }
      }
    } else {
      m.state = m.Hi_In_None;
    }
  }
  return l,p
}

func (m *Highlight_CPP) Hi_NumberExponent( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );
  if( LL <= p ) { m.state = m.Hi_In_None;
  } else {
    var c1 rune = m.p_fb.GetR( l, p );
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST );
      p++;
    } else {
      m.state = m.Hi_In_None;
    }
  }
  return l,p
}

func (m *Highlight_CPP) Hi_NumberTypeSpec( l, p int ) (int,int) {
  LL := m.p_fb.LineLen( l );

  if( p < LL ) {
    var c0 rune = m.p_fb.GetR( l, p );

    if( c0=='L' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE );
      p++;
      m.state = m.Hi_In_None;
    } else if( c0=='F' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE );
      p++;
      m.state = m.Hi_In_None;
    } else if( c0=='U' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++;
      if( p<LL ) {
        var c1 rune = m.p_fb.GetR( l, p );
        if( c1=='L' ) { // UL
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++;
          if( p<LL ) {
            var c2 rune = m.p_fb.GetR( l, p );
            if( c2=='L' ) { // ULL
              m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++;
            }
          }
        }
      }
      m.state = m.Hi_In_None;
    }
  }
  return l,p
}

func (m *Highlight_CPP) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  HiPairs := [...]HiKeyVal {
    { "if"                 , HI_CONTROL },
    { "else"               , HI_CONTROL },
    { "for"                , HI_CONTROL },
    { "while"              , HI_CONTROL },
    { "do"                 , HI_CONTROL },
    { "return"             , HI_CONTROL },
    { "switch"             , HI_CONTROL },
    { "case"               , HI_CONTROL },
    { "break"              , HI_CONTROL },
    { "default"            , HI_CONTROL },
    { "continue"           , HI_CONTROL },
    { "template"           , HI_CONTROL },
    { "public"             , HI_CONTROL },
    { "protected"          , HI_CONTROL },
    { "private"            , HI_CONTROL },
    { "typedef"            , HI_CONTROL },
    { "delete"             , HI_CONTROL },
    { "operator"           , HI_CONTROL },
    { "sizeof"             , HI_CONTROL },
    { "using"              , HI_CONTROL },
    { "namespace"          , HI_CONTROL },
    { "goto"               , HI_CONTROL },
    { "friend"             , HI_CONTROL },
    { "try"                , HI_CONTROL },
    { "catch"              , HI_CONTROL },
    { "throw"              , HI_CONTROL },
    { "and"                , HI_CONTROL },
    { "or"                 , HI_CONTROL },
    { "not"                , HI_CONTROL },
    { "new"                , HI_CONTROL },
    { "const_cast"         , HI_CONTROL },
    { "static_cast"        , HI_CONTROL },
    { "dynamic_cast"       , HI_CONTROL },
    { "reinterpret_cast"   , HI_CONTROL },
    { "override"           , HI_CONTROL },

    // Types
    { "auto"               , HI_VARTYPE },
    { "int"                , HI_VARTYPE },
    { "long"               , HI_VARTYPE },
    { "void"               , HI_VARTYPE },
    { "this"               , HI_VARTYPE },
    { "bool"               , HI_VARTYPE },
    { "char"               , HI_VARTYPE },
    { "const"              , HI_VARTYPE },
    { "constexpr"          , HI_VARTYPE },
    { "short"              , HI_VARTYPE },
    { "float"              , HI_VARTYPE },
    { "double"             , HI_VARTYPE },
    { "signed"             , HI_VARTYPE },
    { "unsigned"           , HI_VARTYPE },
    { "extern"             , HI_VARTYPE },
    { "static"             , HI_VARTYPE },
    { "enum"               , HI_VARTYPE },
    { "uint8_t"            , HI_VARTYPE },
    { "uint16_t"           , HI_VARTYPE },
    { "uint32_t"           , HI_VARTYPE },
    { "uint64_t"           , HI_VARTYPE },
    { "size_t"             , HI_VARTYPE },
    { "int8_t"             , HI_VARTYPE },
    { "int16_t"            , HI_VARTYPE },
    { "int32_t"            , HI_VARTYPE },
    { "int64_t"            , HI_VARTYPE },
    { "float32_t"          , HI_VARTYPE },
    { "float64_t"          , HI_VARTYPE },
    { "FILE"               , HI_VARTYPE },
    { "DIR"                , HI_VARTYPE },
    { "class"              , HI_VARTYPE },
    { "struct"             , HI_VARTYPE },
    { "union"              , HI_VARTYPE },
    { "typename"           , HI_VARTYPE },
    { "virtual"            , HI_VARTYPE },
    { "inline"             , HI_VARTYPE },
    { "explicit"           , HI_VARTYPE },

    // Constants
    { "true"               , HI_CONST   },
    { "false"              , HI_CONST   },
    { "NULL"               , HI_CONST   },
    { "nullptr"            , HI_CONST   },

    { "__FUNCTION__"       , HI_DEFINE  },
    { "__PRETTY_FUNCTION__", HI_DEFINE  },
    { "__TIMESTAMP__"      , HI_DEFINE  },
    { "__FILE__"           , HI_DEFINE  },
    { "__func__"           , HI_DEFINE  },
    { "__LINE__"           , HI_DEFINE  },
    { ""                   , 0 },
  }
  Hi_FindKey_In_Range( m.p_fb, HiPairs[:], st, fn );
}

