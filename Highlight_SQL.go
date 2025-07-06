
package main

import (
//"fmt"
)

type Highlight_SQL struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_SQL) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_SQL) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_SQL) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_SQL[:], st, fn )
}

var HiPairs_SQL = [...]HiKeyVal {

  { "PRAGMA", HI_DEFINE },

  { "NULL", HI_CONST },

  { "AUTOINCREMENT", HI_CONTROL },
  { "BEGIN"        , HI_CONTROL },
  { "CASCADE"      , HI_CONTROL },
  { "CHECK"        , HI_CONTROL },
  { "COMMIT"       , HI_CONTROL },
  { "CREATE"       , HI_CONTROL },
  { "DEFAULT"      , HI_CONTROL },
  { "DELETE"       , HI_CONTROL },
  { "DROP"         , HI_CONTROL },
  { "EXISTS"       , HI_CONTROL },
  { "FROM"         , HI_CONTROL },
  { "IF"           , HI_CONTROL },
  { "INSERT"       , HI_CONTROL },
  { "INTO"         , HI_CONTROL },
  { "NOT"          , HI_CONTROL },
  { "ON"           , HI_CONTROL },
  { "TRANSACTION"  , HI_CONTROL },
  { "UPDATE"       , HI_CONTROL },
  { "VALUES"       , HI_CONTROL },

  { "FOREIGN"   , HI_VARTYPE },
  { "KEY"       , HI_VARTYPE },
  { "BOOL"      , HI_VARTYPE },
  { "INTEGER"   , HI_VARTYPE },
  { "NUMERIC"   , HI_VARTYPE },
  { "PRIMARY"   , HI_VARTYPE },
  { "REFERENCES", HI_VARTYPE },
  { "TABLE"     , HI_VARTYPE },
  { "TEXT"      , HI_VARTYPE },

  { "", 0 },
}

func (m *Highlight_SQL) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if       ( c1=='-' && c0=='-'  ) { m.state = m.Hi_Beg_Comment
      } else if(            c0=='\'' ) { m.state = m.Hi_In_SingleQuote
      } else if(            c0=='"'  ) { m.state = m.Hi_In_DoubleQuote
      } else if( !IsIdent( c1 ) && IsDigit(c0) ) { m.state = m.Hi_NumberBeg
      } else if( (c1==':' && c0==':') || (c1=='-' && c0=='>') ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_VARTYPE )
        m.p_fb.SetSyntaxStyle( l, p  , HI_VARTYPE )

      } else if( (c1=='=' && c0=='=') ||
                 (c1=='&' && c0=='&') ||
                 (c1=='|' && c0=='|') ||
                 (c1=='|' && c0=='=') ||
                 (c1=='&' && c0=='=') ||
                 (c1=='!' && c0=='=') ||
                 (c1=='+' && c0=='=') ||
                 (c1=='-' && c0=='=') ) {
         m.p_fb.SetSyntaxStyle( l, p-1, HI_CONTROL )
         m.p_fb.SetSyntaxStyle( l, p  , HI_CONTROL )

      } else if( c0=='&' ||
                 c0=='.' ||
                 c0=='*' ||
                 c0=='[' ||
                 c0==']' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE );

      } else if( c0=='~' ||
                 c0=='=' || c0=='^' ||
                 c0==':' || c0=='%' ||
                 c0=='+' || c0=='-' ||
                 c0=='<' || c0=='>' ||
                 c0=='!' || c0=='?' ||
                 c0=='(' || c0==')' ||
                 c0=='{' || c0=='}' ||
                 c0==',' || c0==';' ||
                 c0=='/' || c0=='|' ) {
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

func (m *Highlight_SQL) Hi_Beg_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p++
  m.state = m.Hi_In__Comment
  return l,p
}

func (m *Highlight_SQL) Hi_In__Comment( l, p int ) (int,int) {
  var LL int = m.p_fb.LineLen( l )
  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  p--
  m.state = m.Hi_End_Comment
  return l,p
}

func (m *Highlight_SQL) Hi_End_Comment( l, p int ) (int,int) {
  m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  p=0; l++
  m.state = m.Hi_In_None
  return l,p
}

// This shows one way to re-use class methods in Go:
//
func (m *Highlight_SQL) Hi_In_SingleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_SingleQuote_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

// This shows one way to re-use class methods in Go:
//
func (m *Highlight_SQL) Hi_In_DoubleQuote( l, p int ) (int,int) {

  l,p, m.state = Hi_In_DoubleQuote_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

func (m *Highlight_SQL) Hi_NumberBeg( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberBeg_Base( l,p, m.p_fb, m.Hi_NumberIn, m.Hi_NumberHex )

  return l,p
}

func (m *Highlight_SQL) Hi_NumberIn( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberIn_Base( l,p, m.p_fb, m.Hi_In_None, m.Hi_NumberFraction, m.Hi_NumberExponent )

  return l,p
}

func (m *Highlight_SQL) Hi_NumberHex( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberHex_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

func (m *Highlight_SQL) Hi_NumberFraction( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberFraction_Base( l,p, m.p_fb, m.Hi_In_None, m.Hi_NumberExponent )

  return l,p
}

func (m *Highlight_SQL) Hi_NumberExponent( l, p int ) (int,int) {

  l,p, m.state = Hi_NumberExponent_Base( l,p, m.p_fb, m.Hi_In_None )

  return l,p
}

