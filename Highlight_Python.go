
package main

import (
//"fmt"
)

type Highlight_Python struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_Python) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_Python) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_Python) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_Python[:], st, fn )
}

func (m *Highlight_Python) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      comment := c0=='#' && (0==p || c1!='$')

      if       ( comment )                    { m.state = m.Hi_In_Comment
      } else if( Quote_Start('\'',c2,c1,c0) ) { m.state = m.Hi_SingleQuote
      } else if( Quote_Start('"' ,c2,c1,c0) ) { m.state = m.Hi_DoubleQuote
      } else if( !IsIdent(c1) && IsDigit(c0)) { m.state = m.Hi_NumberBeg
      } else if( (c1==':' && c0==':') || (c1=='-' && c0=='>') ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_VARTYPE )
        m.p_fb.SetSyntaxStyle( l, p  , HI_VARTYPE )
      } else if( TwoControl( c1, c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_CONTROL )
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONTROL )
      } else if( c0=='$' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
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

func (m *Highlight_Python) Hi_In_Comment( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  p=0; l++
  m.state = m.Hi_In_None

  return l,p
}

func (m *Highlight_Python) Hi_SingleQuote( l, p int ) (int,int) {
  m.state = nil
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
  p++
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='\'') ||
          (c1!='\\' && c0=='\'') ||
          (c1=='\\' && c0=='\'' && slash_escaped) ) {
        // End of single quote:
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++
        m.state = m.Hi_In_None
      } else {
        if( (c1!='\\' && c0=='$') ||
            (c1=='\\' && c0=='$' && slash_escaped) ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        }
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else                     { slash_escaped = false
        }
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Python) Hi_DoubleQuote( l, p int ) (int,int) {
  m.state = nil
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST ); p++

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: c1,c0
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='"') ||
          (c1!='\\' && c0=='"') ||
          (c1=='\\' && c0=='"' && slash_escaped) ) {
        // End of double quote:
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST ); p++
        m.state = m.Hi_In_None
      } else {
        if( (c1!='\\' && c0=='$') ||
            (c1=='\\' && c0=='$' && slash_escaped) ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        } else {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        }
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else                     { slash_escaped = false
        }
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Python) Hi_NumberBeg( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

  var c1 rune = m.p_fb.GetR( l, p )
  p++
  m.state = m.Hi_NumberIn

  LL := m.p_fb.LineLen( l )
  if( '0' == c1 && (p+1)<LL ) {
    var c0 rune = m.p_fb.GetR( l, p )
    if( 'x' == c0 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      m.state = m.Hi_NumberHex
      p++
    }
  }
  return l,p
}

func (m *Highlight_Python) Hi_NumberIn( l, p int ) (int,int) {
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
    } else {
      m.state = m.Hi_In_None
    }
  }
  return l,p
}

func (m *Highlight_Python) Hi_NumberHex( l, p int ) (int,int) {
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

func (m *Highlight_Python) Hi_NumberFraction( l, p int ) (int,int) {
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

func (m *Highlight_Python) Hi_NumberExponent( l, p int ) (int,int) {
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

var HiPairs_Python = [...]HiKeyVal {
  { "and"         , HI_CONTROL },
  { "break"       , HI_CONTROL },
  { "continue"    , HI_CONTROL },
  { "elif"        , HI_CONTROL },
  { "else"        , HI_CONTROL },
  { "except"      , HI_CONTROL },
  { "finally"     , HI_CONTROL },
  { "for"         , HI_CONTROL },
  { "if"          , HI_CONTROL },
  { "in"          , HI_CONTROL },
  { "is"          , HI_CONTROL },
  { "not"         , HI_CONTROL },
  { "or"          , HI_CONTROL },
  { "pass"        , HI_CONTROL },
  { "raise"       , HI_CONTROL },
  { "return"      , HI_CONTROL },
  { "try"         , HI_CONTROL },
  { "while"       , HI_CONTROL },

  { "as"          , HI_VARTYPE },
  { "class"       , HI_VARTYPE },
  { "def"         , HI_VARTYPE },
  { "del"         , HI_VARTYPE },
  { "global"      , HI_VARTYPE },
  { "int"         , HI_VARTYPE },
  { "with"        , HI_VARTYPE },

  // Built in functions:
  { "abs"         , HI_VARTYPE },
  { "all"         , HI_VARTYPE },
  { "any"         , HI_VARTYPE },
  { "ascii"       , HI_VARTYPE },
  { "bin"         , HI_VARTYPE },
  { "bool"        , HI_VARTYPE },
  { "bytearrary"  , HI_VARTYPE },
  { "bytes"       , HI_VARTYPE },
  { "callable"    , HI_VARTYPE },
  { "chr"         , HI_VARTYPE },
  { "classmethod" , HI_VARTYPE },
  { "compile"     , HI_VARTYPE },
  { "complex"     , HI_VARTYPE },
  { "delattr"     , HI_VARTYPE },
  { "dict"        , HI_VARTYPE },
  { "dir"         , HI_VARTYPE },
  { "divmod"      , HI_VARTYPE },
  { "enumerate"   , HI_VARTYPE },
  { "eval"        , HI_VARTYPE },
  { "exec"        , HI_VARTYPE },
  { "filter"      , HI_VARTYPE },
  { "float"       , HI_VARTYPE },
  { "format"      , HI_VARTYPE },
  { "frozenset"   , HI_VARTYPE },
  { "getattr"     , HI_VARTYPE },
  { "globals"     , HI_VARTYPE },
  { "hasattr"     , HI_VARTYPE },
  { "hash"        , HI_VARTYPE },
  { "help"        , HI_VARTYPE },
  { "hex"         , HI_VARTYPE },
  { "id"          , HI_VARTYPE },
  { "input"       , HI_VARTYPE },
  { "int"         , HI_VARTYPE },
  { "isinstance"  , HI_VARTYPE },
  { "issubclass"  , HI_VARTYPE },
  { "iter"        , HI_VARTYPE },
  { "len"         , HI_VARTYPE },
  { "list"        , HI_VARTYPE },
  { "locals"      , HI_VARTYPE },
  { "map"         , HI_VARTYPE },
  { "max"         , HI_VARTYPE },
  { "memoryview"  , HI_VARTYPE },
  { "min"         , HI_VARTYPE },
  { "next"        , HI_VARTYPE },
  { "object"      , HI_VARTYPE },
  { "oct"         , HI_VARTYPE },
  { "open"        , HI_VARTYPE },
  { "ord"         , HI_VARTYPE },
  { "pow"         , HI_VARTYPE },
  { "print"       , HI_VARTYPE },
  { "property"    , HI_VARTYPE },
  { "range"       , HI_VARTYPE },
  { "repr"        , HI_VARTYPE },
  { "reversed"    , HI_VARTYPE },
  { "round"       , HI_VARTYPE },
  { "set"         , HI_VARTYPE },
  { "setattr"     , HI_VARTYPE },
  { "slice"       , HI_VARTYPE },
  { "sorted"      , HI_VARTYPE },
  { "staticmethod", HI_VARTYPE },
  { "str"         , HI_VARTYPE },
  { "sum"         , HI_VARTYPE },
  { "super"       , HI_VARTYPE },
  { "tuple"       , HI_VARTYPE },
  { "type"        , HI_VARTYPE },
  { "vars"        , HI_VARTYPE },
  { "zip"         , HI_VARTYPE },

  { "assert"      , HI_DEFINE  },
  { "from"        , HI_DEFINE  },
  { "import"      , HI_DEFINE  },
  { "__import__"  , HI_DEFINE  },
  { "__name__"    , HI_DEFINE  },

  { "False"       , HI_CONST   },
  { "None"        , HI_CONST   },
  { "True"        , HI_CONST   },
  { ""            , 0 },
}

