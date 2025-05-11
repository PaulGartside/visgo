
package main

import (
//"fmt"
)

type Highlight_Dir struct {
  p_fb *FileBuf
  state HiStateFunc
}

func (m *Highlight_Dir) Init(p_fb *FileBuf) {
  m.p_fb = p_fb
}

func (m *Highlight_Dir) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_Dir) Run_Range( st CrsPos, fn int )")

  m.state = m.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.state && l<fn {
    l,p = m.state( l, p )
  }
}

func (m *Highlight_Dir) Hi_In_None( l, p int ) (int,int) {
  m.state = nil
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    if( 0<LL ) {
      R_end := m.p_fb.GetR( l, LL-1 )

      if       ( R_end == DIR_DELIM ) { m.Hi_In_None_Dir( l, LL )
      } else if( 1<LL )               { m.Hi_In_None_File( l, LL )
      }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_Dir) Hi_In_None_Dir( l, LL int ) {

  for k:=0; k<LL-1; k++ {
    // R0 is ahead of R1: (R1,R0)
    var R1 rune = 0; if( 0<k ) { R1 = m.p_fb.GetR( l, k-1 ) }
    var R0 rune =                     m.p_fb.GetR( l, k )

    if( R0 == '.' ) {
      m.p_fb.SetSyntaxStyle( l, k, HI_VARTYPE )
    } else if( R1 == '-' && R0 == '>' ) {
      // -> means symbolic link
      m.p_fb.SetSyntaxStyle( l, k-1, HI_DEFINE )
      m.p_fb.SetSyntaxStyle( l, k  , HI_DEFINE )
    } else {
      m.p_fb.SetSyntaxStyle( l, k, HI_CONTROL )
    }
  }
  m.p_fb.SetSyntaxStyle( l, LL-1, HI_CONST )
}

func (m *Highlight_Dir) Hi_In_None_File( l, LL int ) {

  // P1 is ahead of P0: (P0,P1)
  P0 := m.p_fb.GetR( l, 0 )
  P1 := m.p_fb.GetR( l, 1 )

  if( P0=='.' && P1=='.' ) {
    m.p_fb.SetSyntaxStyle( l, 0, HI_DEFINE )
    m.p_fb.SetSyntaxStyle( l, 1, HI_DEFINE )
  } else {
    found_sym_link := false
    for k:=0; k<LL; k++ {
      // R0 is ahead of R1: (R1,R0)
      var R1 rune = 0; if( 0<k ) { R1 = m.p_fb.GetR( l, k-1 ) }
      var R0 rune =                     m.p_fb.GetR( l, k )

      if( R0 == '.' ) {
        m.p_fb.SetSyntaxStyle( l, k, HI_VARTYPE )
      } else if( R1 == '-' && R0 == '>' ) {
        found_sym_link = true
        // -> means symbolic link
        m.p_fb.SetSyntaxStyle( l, k-1, HI_DEFINE )
        m.p_fb.SetSyntaxStyle( l, k  , HI_DEFINE )
      } else if( found_sym_link && R0 == DIR_DELIM ) {
        m.p_fb.SetSyntaxStyle( l, k, HI_CONST )
      }
    }
  }
}

