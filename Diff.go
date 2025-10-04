
package main

import (
  "github.com/gdamore/tcell/v2"
)

type Diff_Type int

const (
  DT_UNKN0WN Diff_Type = iota
  DT_SAME
  DT_CHANGED
  DT_INSERTED
  DT_DELETED
  DT_DIFF_FILES
)

type LineInfo = Vector[Diff_Type]

type Diff_Info struct {
  diff_type Diff_Type  // Diff type of line this Diff_Info refers to
  line_num  int        // Line number in file to which this Diff_Info applies (view line)
  pLineInfo *LineInfo  // Only non-nullptr if diff_type is DT_CHANGED
};

type Diff struct {
  pvS *FileView
  pvL *FileView
  pfS *FileBuf
  pfL *FileBuf

  topLine  int  // top  of buffer view line number.
  leftChar int  // left of buffer view character number.
  crsRow   int  // cursor row    in buffer view. 0 <= crsRow < WorkingRows().
  crsCol   int  // cursor column in buffer view. 0 <= crsCol < WorkingCols().

  DI_List_S Vector[Diff_Info]
  DI_List_L Vector[Diff_Info]

  inVisualMode bool
  inVisualBlock bool

  v_st_line, v_st_char int
  v_fn_line, v_fn_char int
}

func (m *Diff) Init( pv0 *FileView, pv1 *FileView ) {

  if( pv0.p_fb.NumLines() < pv1.p_fb.NumLines() ) {
    m.pvS = pv0
    m.pvL = pv1
  } else {
    m.pvS = pv1
    m.pvL = pv0
  }
  m.pfS = m.pvS.p_fb
  m.pfL = m.pvL.p_fb
}

func (m *Diff) Run() bool {
  m.pvS.p_diff = m
  m.pvL.p_diff = m
  return true
}

//func (m *Diff) NoDiff() {
//  m.pvS.p_diff = nil
//  m.pvL.p_diff = nil
//}

func (m *Diff) NumLines() int {

  // DI_List_L and DI_List_S should be the same length
  return m.DI_List_L.Len()
}

func (m *Diff) Row_Win_2_GL( pV *FileView, win_row int ) int {

  return pV.Y() + 1 + win_row
}

func (m *Diff) Col_Win_2_GL( pV *FileView, win_col int ) int {

  return pV.X() + 1 + win_col
}

func (m *Diff) ViewLine( pV *FileView, diff_line int ) int {
  var ln int
  if       ( pV == m.pvS ) { ln = m.DI_List_S.Get( diff_line ).line_num
  } else if( pV == m.pvL ) { ln = m.DI_List_L.Get( diff_line ).line_num
  }
  return ln
}

func (m *Diff) DiffInfo( pV *FileView, diff_line int ) Diff_Info {
  var di Diff_Info
  if       ( pV == m.pvS ) { di = m.DI_List_S.Get( diff_line )
  } else if( pV == m.pvL ) { di = m.DI_List_L.Get( diff_line )
  }
  return di
}

func (m *Diff) DiffType( pV *FileView, diff_line int ) Diff_Type {

  return m.DiffInfo( pV, diff_line ).diff_type
}

func (m *Diff) Update( pV *FileView ) {

  if       ( pV == m.pvS ) { m.UpdateS()
  } else if( pV == m.pvL ) { m.UpdateL()
  }
}

//func (m *Diff) UpdateS() {
//
//  // Update short view:
//  m.pfS.Find_Styles( m.ViewLineS( m.topLine ) + m.pvS.WorkingRows() )
//  m.pfS.Find_Regexs( m.ViewLineS( m.topLine ), m.pvS.WorkingRows() )
//
//  m.pvS.PrintBorders()
//  m.PrintWorkingView( m.pvS )
//  m.PrintStsLine( m.pvS )
//  m.pvS.PrintFileLine()
//  m.PrintCmdLine( m.pvS )
//}

//func (m *Diff) UpdateL() {
//
//  // Update long view:
//  m.pfL.Find_Styles( m.ViewLineL( m.topLine ) + m.pvL.WorkingRows() )
//  m.pfL.Find_Regexs( m.ViewLineL( m.topLine ), m.pvL.WorkingRows() )
//
//  m.pvL.PrintBorders()
//  m.PrintWorkingView( m.pvL )
//  m.PrintStsLine( m.pvL )
//  m.pvL.PrintFileLine()
//  m.PrintCmdLine( m.pvL )
//}

func (m *Diff) UpdateS() {

  // Update short view:
  m.pvS.p_fb.Find_Styles( m.pvS.topLine + m.pvS.WorkingRows() )
  m.pvS.p_fb.Find_Regexs( m.pvS.topLine, m.pvS.WorkingRows() )

  m.pvS.RepositionView()
  m.pvS.PrintBorders()
  m.pvS.PrintWorkingView()
  m.pvS.PrintStsLine()
  m.pvS.PrintFileLine()

  m.PrintCmdLine( m.pvS )
}

func (m *Diff) UpdateL() {

  // Update long view:
  m.pvL.p_fb.Find_Styles( m.pvS.topLine + m.pvS.WorkingRows() )
  m.pvL.p_fb.Find_Regexs( m.pvS.topLine, m.pvS.WorkingRows() )

  m.pvL.RepositionView()
  m.pvL.PrintBorders()
  m.pvL.PrintWorkingView()
  m.pvL.PrintStsLine()
  m.pvL.PrintFileLine()

  m.PrintCmdLine( m.pvL )
}

//func (m *Diff)  ViewLine( pV *FileView, diff_line int ) int {
//
//  return ( pV == m.pvS ) ? m.DI_List_S[ diff_line ].line_num
//                         : m.DI_List_L[ diff_line ].line_num;
//}

func (m *Diff)  ViewLineS( diff_line int ) int {

  return m.DI_List_S.Get( diff_line ).line_num
}

func (m *Diff)  ViewLineL( diff_line int ) int {

  return m.DI_List_L.Get( diff_line ).line_num
}

//Style DiffStyle( const Style S )
func DiffStyle( p_TS *tcell.Style ) *tcell.Style {
  // If S is already a DIFF style, just return it
  var p_TS_diff *tcell.Style = p_TS

  if       ( p_TS == &TS_NORMAL   ) { p_TS_diff = &TS_DIFF_NORMAL   
  } else if( p_TS == &TS_STAR     ) { p_TS_diff = &TS_DIFF_STAR     
  } else if( p_TS == &TS_STAR_IN_F) { p_TS_diff = &TS_DIFF_STAR_IN_F
  } else if( p_TS == &TS_COMMENT  ) { p_TS_diff = &TS_DIFF_COMMENT  
  } else if( p_TS == &TS_DEFINE   ) { p_TS_diff = &TS_DIFF_DEFINE   
  } else if( p_TS == &TS_CONST    ) { p_TS_diff = &TS_DIFF_CONST    
  } else if( p_TS == &TS_CONTROL  ) { p_TS_diff = &TS_DIFF_CONTROL  
  } else if( p_TS == &TS_VARTYPE  ) { p_TS_diff = &TS_DIFF_VARTYPE  
  } else if( p_TS == &TS_VISUAL   ) { p_TS_diff = &TS_DIFF_VISUAL   
  }
  return p_TS_diff;
}

func (m *Diff) PrintWorkingView( pV *FileView ) {

  NUM_LINES := m.NumLines()
  WR        := pV.WorkingRows()
  WC        := pV.WorkingCols()

  row := 0; // (dl=diff line)
  for dl:=m.topLine; dl<NUM_LINES && row<WR; dl++ {

    G_ROW := m.Row_Win_2_GL( pV, row )
    DT := m.DiffType( pV, dl )

    if( DT == DT_UNKN0WN ) {
      m.PrintWorkingView_DT_UNKN0WN( pV, WC, G_ROW )
    } else if( DT == DT_DELETED ) {
      m.PrintWorkingView_DT_DELETED( pV, WC, G_ROW )
    } else if( DT == DT_CHANGED ) {
      m.PrintWorkingView_DT_CHANGED( pV, WC, G_ROW, dl )
    } else if( DT == DT_DIFF_FILES ) {
      m.PrintWorkingView_DT_DIFF_FILES( pV, WC, G_ROW, dl )
    } else { // DT == DT_INSERTED || DT == DT_SAME
      m.PrintWorkingView_DT_INSERTED_SAME( pV, WC, G_ROW, dl, DT )
    }
    row++
  }
  m.PrintWorkingView_EOF( pV, WR, WC, row )
}

func (m *Diff) PrintWorkingView_DT_UNKN0WN( pV *FileView, WC, G_ROW int ) {

  for col:=0; col<WC; col++ {
    m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), '~', &TS_DIFF_DEL )
  }
}

func (m *Diff) PrintWorkingView_DT_DELETED( pV *FileView, WC, G_ROW int ) {

  for col:=0; col<WC; col++ {
    m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), '-', &TS_DIFF_DEL )
  }
}

func (m *Diff) PrintWorkingView_DT_CHANGED( pV *FileView, WC, G_ROW, dl int ) {

  vl := m.ViewLine( pV, dl ) //(vl=view line)
  LL := pV.p_fb.LineLen( vl )

  var di Diff_Info = m.DiffInfo( pV, dl )

  col := 0

  if( nil != di.pLineInfo ) {
    LIL := di.pLineInfo.Len()
    cp := m.leftChar // char position
    for i:=m.leftChar; cp<LL && i<LIL && col<WC; i++ {
      var dt Diff_Type = di.pLineInfo.Get(i)

      if( DT_SAME == dt ) {
        var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, cp )
        R := pV.p_fb.GetR( vl, cp )
        pV.PrintWorkingView_Set( LL, G_ROW, col, cp, R, p_TS )
        cp++

      } else if( DT_CHANGED == dt || DT_INSERTED == dt ) {
        var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, cp )
        p_TS = DiffStyle( p_TS )
        R := pV.p_fb.GetR( vl, cp )
        pV.PrintWorkingView_Set( LL, G_ROW, col, cp, R, p_TS )
        cp++

      } else if( DT_DELETED == dt ) {
        m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), '-', &TS_DIFF_DEL )

      } else { //( DT_UNKN0WN  == dt )
        m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), '~', &TS_DIFF_DEL )
      }
      col++
    }
    for ; col<WC; col++ {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_EMPTY );
    }
  } else {
    for i:=m.leftChar; i<LL && col<WC; i++ {
      var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, i ); p_TS = DiffStyle( p_TS )
      R := pV.p_fb.GetR( vl, i )
      pV.PrintWorkingView_Set( LL, G_ROW, col, i, R, p_TS )
      col++
    }
    for ; col<WC; col++ {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_DIFF_NORMAL )
    }
  }
}

func (m *Diff) PrintWorkingView_DT_DIFF_FILES( pV *FileView, WC, G_ROW, dl int ) {

  vl := m.ViewLine( pV, dl ) //(vl=view line)
  LL := pV.p_fb.LineLen( vl )
  col := 0

  for i:=m.leftChar; i<LL && col<WC; i++ {
    R := pV.p_fb.GetR( vl, i )
    var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, i )

    pV.PrintWorkingView_Set( LL, G_ROW, col, i, R, p_TS )
    col++
  }
  for ; col<WC; col++ {
    if( col%2==0 ) {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_NORMAL )
    } else {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_DIFF_NORMAL )
    }
  }
}

func (m *Diff) PrintWorkingView_DT_INSERTED_SAME( pV *FileView, WC, G_ROW, dl int, DT Diff_Type ) {

  vl := m.ViewLine( pV, dl ) //(vl=view line)
  LL := pV.p_fb.LineLen( vl )
  col := 0

  for i:=m.leftChar; i<LL && col<WC; i++ {
    R := pV.p_fb.GetR( vl, i )
    var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, i )

    if( DT == DT_INSERTED ) {
      p_TS = DiffStyle( p_TS )
    }
    pV.PrintWorkingView_Set( LL, G_ROW, col, i, R, p_TS );
    col++
  }
  for ; col<WC; col++ {
    if( DT==DT_SAME ) {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_EMPTY )
    } else {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_DIFF_NORMAL )
    }
  }
}

func (m *Diff) PrintWorkingView_EOF( pV *FileView, WR, WC, row int ) {

  // Not enough lines to display, fill in with ~
  for ; row < WR; row++ {
    G_ROW := m.Row_Win_2_GL( pV, row )

    m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, 0 ), '~', &TS_EOF )

    for col:=1; col<WC; col++ {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( pV, col ), ' ', &TS_EOF )
    }
  }
}

func (m *Diff) PrintCursor( pV *FileView ) {
  // FIXME:
}

func (m *Diff) PrintStsLine( pV *FileView ) {
}

func (m *Diff) PrintCmdLine( pV *FileView ) {
  // Prints "--INSERT--" banner, and/or clears command line
  i:=0
  // Draw insert banner if needed
  if( pV.inInsertMode ) {
    i=10 // Strlen of "--INSERT--"
    m_console.SetString( pV.Cmd__Line_Row(), pV.Col_Win_2_GL( 0 ), "--INSERT--", &TS_BANNER )
  }
  WC := pV.WorkingCols()

  for ; i<WC-7; i++ {
    m_console.SetR( pV.Cmd__Line_Row(), pV.Col_Win_2_GL( i ), ' ', &TS_NORMAL )
  }
  m_console.SetString( pV.Cmd__Line_Row(), pV.Col_Win_2_GL( WC-8 ), "--DIFF--", &TS_BANNER )
}

func (m *Diff) InVisualArea( pV *FileView, DL, pos int ) bool {

  // Only one diff view, current view, can be in visual mode.
  if( m_vis.CV() == pV && m.inVisualMode ) {
    if( m.inVisualBlock ) { return m.InVisualBlock( DL, pos )
    } else                { return m.InVisualStFn ( DL, pos )
    }
  }
  return false
}

func (m *Diff) InVisualBlock( DL, pos int ) bool {

  return ( m.v_st_line <= DL  && DL  <= m.v_fn_line &&
           m.v_st_char <= pos && pos <= m.v_fn_char ) || // bot rite
         ( m.v_st_line <= DL  && DL  <= m.v_fn_line &&
           m.v_fn_char <= pos && pos <= m.v_st_char ) || // bot left
         ( m.v_fn_line <= DL  && DL  <= m.v_st_line &&
           m.v_st_char <= pos && pos <= m.v_fn_char ) || // top rite
         ( m.v_fn_line <= DL  && DL  <= m.v_st_line &&
           m.v_fn_char <= pos && pos <= m.v_st_char );// top left
}

func (m *Diff) InVisualStFn( DL, pos int ) bool {

  if( !m.inVisualMode ) { return false }

  if( m.v_st_line == DL && DL == m.v_fn_line ) {
    return (m.v_st_char <= pos && pos <= m.v_fn_char) ||
           (m.v_fn_char <= pos && pos <= m.v_st_char)

  } else if( (m.v_st_line < DL && DL < m.v_fn_line) ||
           (m.v_fn_line < DL && DL < m.v_st_line) ) {
    return true

  } else if( m.v_st_line == DL && DL < m.v_fn_line ) {
    return m.v_st_char <= pos

  } else if( m.v_fn_line == DL && DL < m.v_st_line ) {
    return m.v_fn_char <= pos

  } else if( m.v_st_line < DL && DL == m.v_fn_line ) {
    return pos <= m.v_fn_char

  } else if( m.v_fn_line < DL && DL == m.v_st_line ) {
    return pos <= m.v_st_char
  }
  return false
}

func (m *Diff) Get_Style( pV *FileView, DL, VL, pos int ) *tcell.Style {

  var p_TS *tcell.Style = &TS_EMPTY

  if( VL < pV.p_fb.NumLines() && pos < pV.p_fb.LineLen( VL ) ) {
    p_TS = &TS_NORMAL;

    if       ( m.InVisualArea( pV, DL, pos ) ) { p_TS = &TS_RV_VISUAL
    } else if( pV.InStar   ( VL, pos ) ) { p_TS = &TS_STAR
    } else if( pV.InStarInF( VL, pos ) ) { p_TS = &TS_STAR_IN_F
    } else if( pV.InDefine ( VL, pos ) ) { p_TS = &TS_DEFINE
    } else if( pV.InComment( VL, pos ) ) { p_TS = &TS_COMMENT
    } else if( pV.InConst  ( VL, pos ) ) { p_TS = &TS_CONST
    } else if( pV.InControl( VL, pos ) ) { p_TS = &TS_CONTROL
    } else if( pV.InVarType( VL, pos ) ) { p_TS = &TS_VARTYPE
    }
  }
  return p_TS
}

func (m *Diff) Do_i() {
}

func ( m *Diff ) GoDown( num int ) {
}

func ( m *Diff ) GoUp( num int ) {
}

func ( m *Diff ) GoRight() {
}

func ( m *Diff ) GoLeft() {
}

func (m *Diff) Do_n() {
}

func (m *Diff) Do_N() {
}

func (m *Diff) Do_v() bool {
  return false
}

func (m *Diff) Do_V() bool {
  return false
}

func (m *Diff) Do_a() {
}

func (m *Diff) Do_A() {
}

func (m *Diff) Do_o() {
}

func (m *Diff) Do_O() {
}

func (m *Diff) Do_x() {
}

func (m *Diff) Do_s() {
}

func (m *Diff) Do_cw() {
}

func (m *Diff) Do_D() {
}

func (m *Diff) GoToTopLineInView() {
}

func (m *Diff) GoToBotLineInView() {
}

func (m *Diff) GoToMidLineInView() {
}

func (m *Diff) GoToBegOfLine() {
}

func (m *Diff) GoToEndOfLine() {
}

func (m *Diff) GoToEndOfNextLine() {
}

func (m *Diff) GoToEndOfFile() {
}

func (m *Diff) GoToPrevWord() {
}

func (m *Diff) GoToNextWord() {
}

func (m *Diff) GoToEndOfWord() {
}

func (m *Diff) Do_f( FAST_RUNE rune ) {
}

func (m *Diff) GoToOppositeBracket() {
}

func (m *Diff) GoToLeftSquigglyBracket() {
}

func (m *Diff) GoToRightSquigglyBracket() {
}

func (m *Diff) PageDown() {
}

func (m *Diff) PageUp() {
}

func (m *Diff) Do_Star_GetNewPattern() string {
  return ""
}

func (m *Diff) GoToTopOfFile() {
}

func (m *Diff) GoToStartOfRow() {
}

func (m *Diff) GoToEndOfRow() {
}

func (m *Diff) GoToFile() {
}

func (m *Diff) Do_dd() {
}

func (m *Diff) Do_dw() {
}

func (m *Diff) Do_yy() {
}

func (m *Diff) Do_yw() {
}

func (m *Diff) Do_p() {
}

func (m *Diff) Do_P() {
}

func (m *Diff) Do_r() {
}

func (m *Diff) Do_R() {
}

func (m *Diff) Do_J() {
}

func (m *Diff) Do_Tilda() {
}

func (m *Diff) Do_u() {
}

func (m *Diff) Do_U() {
}

func (m *Diff) MoveCurrLineToTop() {
}

func (m *Diff) MoveCurrLineCenter() {
}

func (m *Diff) MoveCurrLineToBottom() {
}

