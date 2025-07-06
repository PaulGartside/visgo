
package main

import (
//"fmt"
)

type Highlight_Bash struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_Bash) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_Bash) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_Bash) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_Bash[:], st, fn )
}

func (m *Highlight_Bash) Hi_In_None( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_In_Comment( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  p=0; l++
  m.state = m.Hi_In_None

  return l,p
}

func (m *Highlight_Bash) Hi_SingleQuote( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_DoubleQuote( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_NumberBeg( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_NumberIn( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_NumberHex( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_NumberFraction( l, p int ) (int,int) {
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

func (m *Highlight_Bash) Hi_NumberExponent( l, p int ) (int,int) {
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

var HiPairs_Bash = [...]HiKeyVal {
  { "if"                 , HI_CONTROL },
  { "fi"                 , HI_CONTROL },
  { "else"               , HI_CONTROL },
  { "elif"               , HI_CONTROL },
  { "for"                , HI_CONTROL },
  { "done"               , HI_CONTROL },
  { "while"              , HI_CONTROL },
  { "do"                 , HI_CONTROL },
  { "return"             , HI_CONTROL },
  { "switch"             , HI_CONTROL },
  { "case"               , HI_CONTROL },
  { "esac"               , HI_CONTROL },
  { "break"              , HI_CONTROL },
  { "then"               , HI_CONTROL },
  { "bg"                 , HI_CONTROL },
  { "bind"               , HI_CONTROL },
  { "builtin"            , HI_CONTROL },
  { "caller"             , HI_CONTROL },
  { "cd"                 , HI_CONTROL },
  { "command"            , HI_CONTROL },
  { "compgen"            , HI_CONTROL },
  { "complete"           , HI_CONTROL },
  { "compopt"            , HI_CONTROL },
  { "continue"           , HI_CONTROL },
  { "echo"               , HI_CONTROL },
  { "enable"             , HI_CONTROL },
  { "eval"               , HI_CONTROL },
  { "exec"               , HI_CONTROL },
  { "exit"               , HI_CONTROL },
  { "export"             , HI_CONTROL },
  { "fc"                 , HI_CONTROL },
  { "fg"                 , HI_CONTROL },
  { "function"           , HI_CONTROL },
  { "hash"               , HI_CONTROL },
  { "help"               , HI_CONTROL },
  { "history"            , HI_CONTROL },
  { "jobs"               , HI_CONTROL },
  { "kill"               , HI_CONTROL },
  { "logout"             , HI_CONTROL },
  { "popd"               , HI_CONTROL },
  { "printf"             , HI_CONTROL },
  { "pushd"              , HI_CONTROL },
  { "pwd"                , HI_CONTROL },
  { "return"             , HI_CONTROL },
  { "set"                , HI_CONTROL },
  { "shift"              , HI_CONTROL },
  { "shopt"              , HI_CONTROL },
  { "source"             , HI_CONTROL },
  { "suspend"            , HI_CONTROL },
  { "test"               , HI_CONTROL },
  { "times"              , HI_CONTROL },
  { "trap"               , HI_CONTROL },
  { "ulimit"             , HI_CONTROL },
  { "umask"              , HI_CONTROL },
  { "unalias"            , HI_CONTROL },
  { "unset"              , HI_CONTROL },
  { "wait"               , HI_CONTROL },

  { "declare"            , HI_VARTYPE },
  { "dirs"               , HI_VARTYPE },
  { "disown"             , HI_VARTYPE },
  { "getopts"            , HI_VARTYPE },
  { "let"                , HI_VARTYPE },
  { "local"              , HI_VARTYPE },
  { "mapfile"            , HI_VARTYPE },
  { "read"               , HI_VARTYPE },
  { "readonly"           , HI_VARTYPE },
  { "type"               , HI_VARTYPE },
  { "typeset"            , HI_VARTYPE },
  { "@"                  , HI_VARTYPE },
  { "#"                  , HI_VARTYPE },

  { "false"              , HI_CONST   },
  { "true"               , HI_CONST   },

  { "alias"              , HI_DEFINE  },
  { ""                   , 0 },
}

