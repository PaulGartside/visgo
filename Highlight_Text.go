
package main

import (
//"fmt"
)

type Highlight_Text struct {
  p_fb  *FileBuf
  state HiStateFunc
}

func (m *Highlight_Text) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_Text) Run_Range( st CrsPos, fn int ) {

  m.state = m.Hi_In_None

  l := st.crsLine;
  p := st.crsChar;

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
}

func (m *Highlight_Text) Hi_In_None( l, p int ) (int,int) {
  m.state = nil;
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      var R rune = m.p_fb.GetR( l, p )

      if( R == '#' ) {
        m.state = m.Hi_In_Define
      } else if( R != 27 && (R < 32 || 126 < R) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( nil != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Text) Hi_In_Define( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
  }
  p=0; l++;
  m.state = m.Hi_In_None

  return l,p
}

