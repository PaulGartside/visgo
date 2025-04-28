
package main

import (
//"fmt"
//"strings"
)

type Highlight_BufferEditor struct {
  p_fb *FileBuf
  state HiStateFunc
}

func (m *Highlight_BufferEditor) Init(p_fb *FileBuf) {
  m.p_fb = p_fb
}

func (m *Highlight_BufferEditor) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_BufferEditor) Run_Range( st CrsPos, fn int )")

  m.state = m.Hi_In_None

  l := st.crsLine;
  p := st.crsChar;

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
}

func (m *Highlight_BufferEditor) Hi_In_None( l, p int ) (int,int) {
  m.state = nil;
  for ; l<m.p_fb.NumLines(); l++ {
    var p_fl *FLine = m.p_fb.GetLP( l );
    var LL int   = m.p_fb.LineLen( l );

    if( 0<LL ) {
      var c_end rune = m.p_fb.GetR( l, LL-1 );

      if( p_fl.EqualStr( m_EDIT_BUF_NAME ) ||
          p_fl.EqualStr( m_HELP_BUF_NAME ) ||
          p_fl.EqualStr( m_MSG__BUF_NAME ) ||
          p_fl.EqualStr( m_SHELL_BUF_NAME ) ||
          p_fl.EqualStr( m_COLON_BUF_NAME ) ||
          p_fl.EqualStr( m_SLASH_BUF_NAME ) ) {
        for k:=0; k<LL; k++ {
          m.p_fb.SetSyntaxStyle( l, k, HI_DEFINE );
        }
      } else if( c_end == DIR_DELIM ) {
        for k:=0; k<LL; k++ {
          var R rune = m.p_fb.GetR( l, k );
          if( R == DIR_DELIM ) {
            m.p_fb.SetSyntaxStyle( l, k, HI_CONST );
          } else {
            m.p_fb.SetSyntaxStyle( l, k, HI_CONTROL );
          }
        }
      } else {
        for k:=0; k<LL; k++ {
          var R rune = m.p_fb.GetR( l, k );
          if( R == DIR_DELIM ) {
            m.p_fb.SetSyntaxStyle( l, k, HI_CONST );
          }
        }
      }
    }
    p = 0;
  }
  return l,p
}

