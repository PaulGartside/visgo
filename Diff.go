
package main

import (
  "bytes"
  "fmt"
  "github.com/gdamore/tcell/v2"
  "strings"
  "time"
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

const max_files_added_per_diff = 10

type LineInfo = Vector[Diff_Type]

type Diff_Info struct {
  diff_type Diff_Type  // Diff type of line this Diff_Info refers to
  line_num  int        // Line number in file to which this Diff_Info applies (view line)
  pLineInfo *LineInfo  // Only non-nullptr if diff_type is DT_CHANGED
}

type SimLines struct { // Similar lines
  ln_s   int       // Line number in short comp area
  ln_l   int       // Line number in long  comp area
  nbytes int       // Number of bytes in common between lines
  li_s   *LineInfo // Line comparison info in short comp area
  li_l   *LineInfo // Line comparison info in long  comp area
}

type Diff struct {
  pvS *FileView
  pvL *FileView
  pfS *FileBuf
  pfL *FileBuf

  topLine  int  // top  of buffer view line number.
  leftChar int  // left of buffer view character number.
  crsRow   int  // cursor row    in buffer view. 0 <= crsRow < WorkingRows().
  crsCol   int  // cursor column in buffer view. 0 <= crsCol < WorkingCols().

  DI_List_S Vector[Diff_Info] // One Diff_Info per diff line
  DI_List_L Vector[Diff_Info] // One Diff_Info per diff line

  inVisualMode bool
  inVisualBlock bool

  v_st_line, v_st_char int
  v_fn_line, v_fn_char int

  printed_diff_ms bool
  diff_dur time.Duration

  sameList Vector[SameArea]
  diffList Vector[DiffArea]
//simiList Vector[SimLines]

  num_files_added_this_diff int
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

// Returns true if diff took place, else false
//
func (m *Diff) Run() bool {
  ran_diff := false
  if( m.pfS != m.pfL ) {
    m.pvS.p_diff = m
    m.pvL.p_diff = m
    // All lines in both files:
    DA := DiffArea{ 0, m.pfS.NumLines(), 0, m.pfL.NumLines() }
    m.RunDiff( DA )
    ran_diff = true

    m.Find_Context()
  }
  return ran_diff
}

func (m *Diff) NoDiff() {

  m.pvS.topLine  = m.ViewLineS( m.topLine )
  m.pvS.leftChar = m.leftChar
  m.pvS.crsRow   = m.crsRow
  m.pvS.crsCol   = m.crsCol

  m.pvL.topLine  = m.ViewLineL( m.topLine )
  m.pvL.leftChar = m.leftChar
  m.pvL.crsRow   = m.crsRow
  m.pvL.crsCol   = m.crsCol

  m.pvS.p_diff = nil
  m.pvL.p_diff = nil
}

func (m *Diff) NumLines() int {

  // DI_List_L and DI_List_S should be the same length
  return m.DI_List_L.Len()
}

func (m *Diff) CrsLine() int {

  return m.topLine  + m.crsRow;
}

func (m *Diff) CrsChar() int {

  return m.leftChar + m.crsCol;
}

//func (m *Diff) Row_Win_2_GL( pV *FileView, win_row int ) int {
//
//  return pV.Y() + 1 + win_row
//}

//func (m *Diff) Col_Win_2_GL( pV *FileView, win_col int ) int {
//
//  return pV.X() + 1 + win_col
//}

func (m *Diff) LineLen() int {

  pV := m_vis.CV()

  diff_line := m.CrsLine()

  var di Diff_Info = m.DiffInfo( pV, diff_line )

  if( DT_UNKN0WN == di.diff_type ||
      DT_DELETED == di.diff_type ) {
    return 0;
  }
  view_line := di.line_num

  return pV.p_fb.LineLen( view_line )
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

func (m *Diff) View_C( pV *FileView ) *FileView {
  var p_v *FileView
  if       ( pV == m.pvS ) { p_v = m.pvS
  } else if( pV == m.pvL ) { p_v = m.pvL
  }
  return p_v
}

func (m *Diff) View_O( pV *FileView ) *FileView {
  var p_v *FileView
  if       ( pV == m.pvS ) { p_v = m.pvL
  } else if( pV == m.pvL ) { p_v = m.pvS
  }
  return p_v
}

// View to Diff Info List of Current View
//
func (m *Diff) View_2_DI_List_C( pV *FileView ) *Vector[Diff_Info] {
  var p_DI_List *Vector[Diff_Info] = nil
  if       ( pV == m.pvS ) { p_DI_List = &m.DI_List_S
  } else if( pV == m.pvL ) { p_DI_List = &m.DI_List_L
  }
  return p_DI_List
}

// View to Diff Info List of Other View
//
func (m *Diff) View_2_DI_List_O( pV *FileView ) *Vector[Diff_Info] {
  var p_DI_List *Vector[Diff_Info] = nil
  if       ( pV == m.pvS ) { p_DI_List = &m.DI_List_L
  } else if( pV == m.pvL ) { p_DI_List = &m.DI_List_S
  }
  return p_DI_List
}

func (m *Diff) DiffType( pV *FileView, diff_line int ) Diff_Type {

  return m.DiffInfo( pV, diff_line ).diff_type
}

func (m *Diff) WorkingRows( pV *FileView ) int {

  return pV.WinRows() - 5
}

func (m *Diff) WorkingCols( pV *FileView ) int {

  return pV.WinCols() - 2
}

func (m *Diff) BotLine( pV *FileView ) int {

  return m.topLine + m.WorkingRows( pV )-1
}

func (m *Diff) RightChar( pV *FileView ) int {

  return m.leftChar + m.WorkingCols( pV )-1
}

func (m *Diff) Find_Context() {

  if( !m.Has_Context() ) {
    pV := m_vis.CV()

    if( pV.Has_Context() ) {
      m.Copy_ViewContext_2_DiffContext()
    } else {
      m.Do_n_Diff( false )
      m.MoveCurrLineCenter( false )
    }
  }
}

func (m *Diff) Has_Context() bool {

  return 0 != m.topLine ||
         0 != m.leftChar ||
         0 != m.crsRow ||
         0 != m.crsCol
}

func (m *Diff) Copy_ViewContext_2_DiffContext() {
  pV := m_vis.CV()

  // View context -> diff context
  diff_topLine := m.DiffLine( pV, pV.topLine )
  diff_crsLine := m.DiffLine( pV, pV.CrsLine() )
  diff_crsRow  := diff_crsLine - diff_topLine

  m.topLine  = diff_topLine
  m.leftChar = pV.leftChar
  m.crsRow   = diff_crsRow
  m.crsCol   = pV.crsCol
}

func (m *Diff) DiffLine( pV *FileView, view_line int ) int {
  dl := 0
  if       ( pV == m.pvS ) { dl = m.DiffLine_S( view_line )
  } else if( pV == m.pvL ) { dl = m.DiffLine_L( view_line )
  }
  return dl
}

// Return the diff line of the view line on the short side
func (m *Diff) DiffLine_S( view_line int ) int {

  diff_line := 0
  NUM_LINES_VS := m.pvS.p_fb.NumLines()

  if( 0 < NUM_LINES_VS ) {
    DI_LEN := m.DI_List_S.Len()

    if( NUM_LINES_VS <= view_line ) {
      diff_line = DI_LEN-1
    } else {
      // Diff line is greater or equal to view line,
      // so start at view line number and search forward
      k := view_line
      var di Diff_Info = m.DI_List_S.Get( view_line )
      k += view_line - di.line_num
      found := false

      for ; !found && k<DI_LEN; k += view_line - di.line_num {
        di = m.DI_List_S.Get( k )

        if( view_line == di.line_num ) {
          found = true
          diff_line = k
        }
      }
    }
  }
  return diff_line
}

// Return the diff line of the view line on the long side
func (m *Diff) DiffLine_L( view_line int ) int {

  diff_line := 0
  NUM_LINES_VL := m.pvL.p_fb.NumLines()

  if( 0 < NUM_LINES_VL ) {
    DI_LEN := m.DI_List_L.Len()

    if( NUM_LINES_VL <= view_line ) {
      diff_line = DI_LEN-1
    } else {
      // Diff line is greater or equal to view line,
      // so start at view line number and search forward
      k := view_line
      var di Diff_Info = m.DI_List_L.Get( view_line )
      k += view_line - di.line_num
      found := false

      for ; !found && k<DI_LEN; k += view_line - di.line_num {
        di = m.DI_List_L.Get( k )

        if( view_line == di.line_num ) {
          found = true
          diff_line = k
        }
      }
    }
  }
  return diff_line
}

func (m *Diff) RepositionViews() {
  // If a window re-size has taken place, and the window has gotten
  // smaller, change top line and left char if needed, so that the
  // cursor is in the buffer when it is re-drawn
  pV := m_vis.CV()

  if( m.WorkingRows( pV ) <= m.crsRow ) {
    m.topLine += ( m.crsRow - m.WorkingRows( pV ) + 1 )
    m.crsRow  -= ( m.crsRow - m.WorkingRows( pV ) + 1 )
  }
  if( m.WorkingCols( pV ) <= m.crsCol ) {
    m.leftChar += ( m.crsCol - m.WorkingCols( pV ) + 1 )
    m.crsCol   -= ( m.crsCol - m.WorkingCols( pV ) + 1 )
  }
}

// Update both views
func (m *Diff) UpdateBV() {

  m.RepositionViews()
  m.UpdateS()
  m.UpdateL()

  if( ! m.printed_diff_ms ) {
    msg := fmt.Sprintf("Diff took: %v ms", m.diff_dur.Milliseconds())
    m_vis.CmdLineMessage( msg )

    m.printed_diff_ms = true
  } else {
    m.PrintCursor( m_vis.CV() )  // Put cursor into position.
  }
}

// Update 1 view
func (m *Diff) Update1V( pV *FileView ) {

  m.RepositionViews()

  if       ( pV == m.pvS ) { m.UpdateS()
  } else if( pV == m.pvL ) { m.UpdateL()
  }
}

func (m *Diff) UpdateS() {

  // Update short view:
  m.pfS.Find_Styles( m.ViewLineS( m.topLine ) + m.WorkingRows( m.pvS ) );
  m.pfS.Find_Regexs( m.ViewLineS( m.topLine ), m.WorkingRows( m.pvS ) );

  m.pvS.PrintBorders()
  m.PrintWorkingView( m.pvS )
  m.PrintStsLine( m.pvS )
  m.pvS.PrintFileLine()

  m.PrintCmdLine( m.pvS )
}

func (m *Diff) UpdateL() {

  // Update long view:
  m.pfL.Find_Styles( m.ViewLineL( m.topLine ) + m.WorkingRows( m.pvL ) );
  m.pfL.Find_Regexs( m.ViewLineL( m.topLine ), m.WorkingRows( m.pvL ) );

  m.pvL.PrintBorders()
  m.PrintWorkingView( m.pvL )
  m.PrintStsLine( m.pvL )
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

    G_ROW := pV.Row_Win_2_GL( row )
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
    m_console.SetR( G_ROW, pV.Col_Win_2_GL( col ), '~', &TS_DIFF_DEL )
  }
}

func (m *Diff) PrintWorkingView_DT_DELETED( pV *FileView, WC, G_ROW int ) {

  for col:=0; col<WC; col++ {
    m_console.SetR( G_ROW, pV.Col_Win_2_GL( col ), '-', &TS_DIFF_DEL )
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
      G_COL := pV.Col_Win_2_GL( col )
      var dt Diff_Type = di.pLineInfo.Get(i)

      if( DT_SAME == dt ) {
        var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, cp )
        R := pV.p_fb.GetR( vl, cp )
        pV.PrintWorkingView_Set( LL, G_ROW, G_COL, cp, R, p_TS )
        cp++

      } else if( DT_CHANGED == dt || DT_INSERTED == dt ) {
        var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, cp )
        p_TS = DiffStyle( p_TS )
        R := pV.p_fb.GetR( vl, cp )
        pV.PrintWorkingView_Set( LL, G_ROW, G_COL, cp, R, p_TS )
        cp++

      } else if( DT_DELETED == dt ) {
        m_console.SetR( G_ROW, G_COL, '-', &TS_DIFF_DEL )

      } else { //( DT_UNKN0WN  == dt )
        m_console.SetR( G_ROW, G_COL, '~', &TS_DIFF_DEL )
      }
      col++
    }
    for ; col<WC; col++ {
      G_COL := pV.Col_Win_2_GL( col )
      m_console.SetR( G_ROW, G_COL, ' ', &TS_EMPTY );
    }
  } else {
    for i:=m.leftChar; i<LL && col<WC; i++ {
      var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, i ); p_TS = DiffStyle( p_TS )
      G_COL := pV.Col_Win_2_GL( col )
      R := pV.p_fb.GetR( vl, i )
      pV.PrintWorkingView_Set( LL, G_ROW, G_COL, i, R, p_TS )
      col++
    }
    for ; col<WC; col++ {
      G_COL := pV.Col_Win_2_GL( col )
      m_console.SetR( G_ROW, G_COL, ' ', &TS_DIFF_NORMAL )
    }
  }
}

func (m *Diff) PrintWorkingView_DT_DIFF_FILES( pV *FileView, WC, G_ROW, dl int ) {

  vl := m.ViewLine( pV, dl ) //(vl=view line)
  LL := pV.p_fb.LineLen( vl )
  col := 0

  for i:=m.leftChar; i<LL && col<WC; i++ {
    G_COL := pV.Col_Win_2_GL( col )
    R := pV.p_fb.GetR( vl, i )
    var p_TS *tcell.Style = m.Get_Style( pV, dl, vl, i )

    pV.PrintWorkingView_Set( LL, G_ROW, G_COL, i, R, p_TS )
    col++
  }
  for ; col<WC; col++ {
    G_COL := pV.Col_Win_2_GL( col )
    if( col%2==0 ) {
      m_console.SetR( G_ROW, G_COL, ' ', &TS_NORMAL )
    } else {
      m_console.SetR( G_ROW, G_COL, ' ', &TS_DIFF_NORMAL )
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
    G_COL := pV.Col_Win_2_GL( col )
    pV.PrintWorkingView_Set( LL, G_ROW, G_COL, i, R, p_TS );
    col++
  }
  for ; col<WC; col++ {
    if( DT==DT_SAME ) {
      m_console.SetR( G_ROW, pV.Col_Win_2_GL( col ), ' ', &TS_EMPTY )
    } else {
      m_console.SetR( G_ROW, pV.Col_Win_2_GL( col ), ' ', &TS_DIFF_NORMAL )
    }
  }
}

func (m *Diff) PrintWorkingView_EOF( pV *FileView, WR, WC, row int ) {

  // Not enough lines to display, fill in with ~
  for ; row < WR; row++ {
    G_ROW := pV.Row_Win_2_GL( row )

    m_console.SetR( G_ROW, pV.Col_Win_2_GL( 0 ), '~', &TS_EOF )

    for col:=1; col<WC; col++ {
      m_console.SetR( G_ROW, pV.Col_Win_2_GL( col ), ' ', &TS_EOF )
    }
  }
}

func (m *Diff) PrintCursor( pV *FileView ) {

  m_console.ShowCursor( pV.Row_Win_2_GL( m.crsRow ), pV.Col_Win_2_GL( m.crsCol ) )
  m_console.Show()
}

func (m *Diff) PrintStsLine( pV *FileView ) {

  var p_DI_List *Vector[Diff_Info] = m.View_2_DI_List_C( pV )
  pfb := pV.p_fb
  CLd := m.CrsLine()                   // Line position diff
  CLv := p_DI_List.Get( CLd ).line_num // Line position view
  CC  := m.CrsChar()                   // Char position
  LL := 0
  if( 0 < m.NumLines() &&  0 < pfb.NumLines() ) {
    LL = pfb.LineLen( CLv )
  }
  WC := m.WorkingCols( pV )

  fileSize := pfb.GetSize()
  crsByte  := pfb.GetCursorByte( CLv, CC )
  percent := int(100*float64(crsByte)/float64(fileSize) + 0.5)

 var buf bytes.Buffer

  fmt.Fprintf( &buf, "Pos=(%d,%d)  (%d%%, %d/%d)  Char=(",
                     CLv+1, CC+1,
                     percent, crsByte, fileSize )
  if 0 < LL && CC < LL {
    var R rune = pfb.GetR( CLv, CC )
    fmt.Fprintf( &buf, "%d,%c", R, R )
  }
  fmt.Fprintf( &buf, ")" )

  for k:=buf.Len(); k<WC; k++ {
    fmt.Fprintf( &buf, " " )
  }
  if( WC < buf.Len() ) {
    buf.Truncate( WC )
  }
  m_console.SetBuffer( pV.Sts__Line_Row(), pV.Col_Win_2_GL( 0 ), &buf, &TS_BORDER )
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
  // FIXME
}

func ( m *Diff ) GoDown( num int ) {

  NUM_LINES := m.NumLines()
  OCL       := m.CrsLine() // Old cursor line

  if( 0 < NUM_LINES && OCL < NUM_LINES-1 ) {
    NCL := OCL+num // New cursor line

    if( NUM_LINES-1 < NCL ) { NCL = NUM_LINES-1 }

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func ( m *Diff ) GoUp( num int ) {

  NUM_LINES := m.NumLines()
  OCL       := m.CrsLine() // Old cursor line

  if( 0 < NUM_LINES && 0 < OCL ) {
    NCL := OCL-num // New cursor line

    if( NCL < 0 ) { NCL = 0 }

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func ( m *Diff ) GoRight( num int ) {

  if( 0<m.NumLines() ) {
    LL  := m.LineLen()
    OCP := m.CrsChar() // Old cursor position

    if( 0<LL && OCP < LL-1 ) {
      NCP := OCP+num // New cursor position

      if( LL-1 < NCP ) { NCP = LL-1 }

      m.GoToCrsPos_Write( m.CrsLine(), NCP )
    }
  }
}

func ( m *Diff ) GoLeft( num int ) {

  OCP := m.CrsChar() // Old cursor position

  if( 0 < m.NumLines() && 0 < OCP ) {
    NCP := OCP-num // New cursor position

    if( NCP < 0 ) { NCP = 0 }

    m.GoToCrsPos_Write( m.CrsLine(), NCP )
  }
}

func (m *Diff) Do_n() {

  if( 0<len(m_vis.regex_str) ) { m.Do_n_Pattern()
  } else                       { m.Do_n_Diff(true)
  }
}

func (m *Diff) Do_N() {

  if( 0<len(m_vis.regex_str) ) { m.Do_N_Pattern()
  } else                       { m.Do_N_Diff()
  }
}

func (m *Diff) Do_v() bool {
  // FIXME
  return false
}

func (m *Diff) Do_V() bool {
  // FIXME
  return false
}

func (m *Diff) Do_a() {
  // FIXME
}

func (m *Diff) Do_A() {
  // FIXME
}

func (m *Diff) Do_o() {
  // FIXME
}

func (m *Diff) Do_O() {
  // FIXME
}

func (m *Diff) Do_x() {
  // FIXME
}

func (m *Diff) Do_s() {
  // FIXME
}

func (m *Diff) Do_cw() {
  // FIXME
}

func (m *Diff) Do_D() {
  // FIXME
}

func (m *Diff) GoToTopLineInView() {

  m.GoToCrsPos_Write( m.topLine, m.CrsChar() )
}

func (m *Diff) GoToBotLineInView() {

  pV := m_vis.CV()

  NUM_LINES := m.NumLines()

  bottom_line_in_view := m.topLine + m.WorkingRows( pV )-1

  bottom_line_in_view = Min_i( NUM_LINES-1, bottom_line_in_view )

  m.GoToCrsPos_Write( bottom_line_in_view, m.CrsChar() )
}

func (m *Diff) GoToMidLineInView() {

  pV := m_vis.CV()

  NUM_LINES := m.NumLines()

  // Default: Last line in file is not in view
  crsLine := m.topLine + m.WorkingRows( pV )/2

  if( NUM_LINES-1 < m.BotLine( pV ) ) {
    // Last line in file above bottom of view
    crsLine = m.topLine + (NUM_LINES-1 - m.topLine)/2
  }
  m.GoToCrsPos_Write( crsLine, 0 )
}

func (m *Diff) GoToBegOfLine() {

  if( 0 < m.NumLines() ) {
    CL := m.CrsLine() // Cursor line

    m.GoToCrsPos_Write( CL, 0 )
  }
}

func (m *Diff) GoToEndOfLine() {

  if( 0 < m.NumLines() ) {
    LL := m.LineLen()

    OCL := m.CrsLine() // Old cursor line
    NCP := LLM1(LL)

    m.GoToCrsPos_Write( OCL, NCP );
  }
}

func (m *Diff) GoToEndOfNextLine() {

  NUM_LINES := m.NumLines() // Diff

  if( 0<NUM_LINES ) {
    OCL := m.CrsLine() // Old cursor diff line

    if( OCL < (NUM_LINES-1) ) {
      // Before last line, so can go down
      pV  := m_vis.CV()
      pfb := pV.p_fb
      VL := m.ViewLine( pV, OCL+1 ) // View line
      LL := pfb.LineLen( VL )

      m.GoToCrsPos_Write( OCL+1, LLM1(LL) )
    }
  }
}

func (m *Diff) GoToEndOfFile() {

  NUM_LINES := m.NumLines()

  if( 0<NUM_LINES ) {
    m.GoToCrsPos_Write( NUM_LINES-1, 0 )
  }
}

func (m *Diff) GoToNextWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToNextWord_GetPosition( &ncp ) ) {
    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *Diff) GoToPrevWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToPrevWord_GetPosition( &ncp ) ) {
    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *Diff) GoToEndOfWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToEndOfWord_GetPosition( &ncp ) ) {
    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *Diff) Do_f( FAST_RUNE rune ) {

  NUM_LINES := m.NumLines()

  if( 0 < NUM_LINES ) {
    pV := m_vis.CV()
    pfb := pV.p_fb

    DL  := m.CrsLine()          // Diff line
    VL  := m.ViewLine( pV, DL ) // View line
    LL  := pfb.LineLen( VL )    // Line length
    OCP := m.CrsChar()          // Old cursor position

    if( OCP < LL-1 ) {
      NCP := 0
      found_rune := false
      for p:=OCP+1; !found_rune && p<LL; p++ {
        R := pfb.GetR( VL, p )

        if( R == FAST_RUNE ) {
          NCP = p
          found_rune = true
        }
      }
      if( found_rune ) {
        m.GoToCrsPos_Write( DL, NCP )
      }
    }
  }
}

func (m *Diff) GoToOppositeBracket() {

  pV := m_vis.CV()

  m.MoveInBounds_Line()

  NUM_LINES := pV.p_fb.NumLines()
  CL        := m.ViewLine( pV, m.CrsLine() ) //< View line
  CC        := m.CrsChar()
  LL        := m.LineLen()

  if( 0 < NUM_LINES && 0 < LL ) {

    R := pV.p_fb.GetR( CL, CC )

    if( R=='{' || R=='[' || R=='(' ) {
      var finish_rune rune = 0
      if       ( R=='{' ) { finish_rune = '}'
      } else if( R=='[' ) { finish_rune = ']'
      } else if( R=='(' ) { finish_rune = ')'
      } else              { ; // Un-handled case
      }
      m.GoToOppositeBracket_Forward( R, finish_rune )

    } else if( R=='}' || R==']' || R==')' ) {
      var finish_rune rune = 0
      if       ( R=='}' ) { finish_rune = '{'
      } else if( R==']' ) { finish_rune = '['
      } else if( R==')' ) { finish_rune = '('
      } else              { ; // Un-handled case
      }
      m.GoToOppositeBracket_Backward( R, finish_rune )
    }
  }
}

func (m *Diff) GoToLeftSquigglyBracket() {

  m.MoveInBounds_Line()

  var  start_char rune = '}'
  var finish_char rune = '{'

  m.GoToOppositeBracket_Backward( start_char, finish_char )
}

func (m *Diff) GoToRightSquigglyBracket() {

  m.MoveInBounds_Line()

  var  start_char rune = '{'
  var finish_char rune = '}'

  m.GoToOppositeBracket_Forward( start_char, finish_char )
}

func (m *Diff) PageDown() {

  NUM_LINES := m.NumLines()
  if( 0 < NUM_LINES ) {

    pV := m_vis.CV()

    // new diff top line:
    newTopLine := m.topLine + m.WorkingRows( pV ) - 1
    // Subtracting 1 above leaves one line in common between the 2 pages.

    if( newTopLine < NUM_LINES ) {
      m.crsCol = 0
      m.topLine = newTopLine

      // Dont let cursor go past the end of the file:
      if( NUM_LINES <= m.CrsLine() ) {
        // This line places the cursor at the top of the screen, which I prefer:
        m.crsRow = 0
      }
      m.UpdateBV();
    }
  }
}

func (m *Diff) PageUp() {

  // Dont scroll if we are at the top of the file:
  if( 0 < m.topLine ) {
    //Leave crsRow unchanged.
    m.crsCol = 0

    pV := m_vis.CV()

    // Dont scroll past the top of the file:
    if( m.topLine < m.WorkingRows( pV ) - 1 ) {
      m.topLine = 0;
    } else {
      m.topLine -= m.WorkingRows( pV ) - 1
    }
    m.UpdateBV()
  }
}

func (m *Diff) Do_Star_GetNewPattern() string {
  return ""
}

func (m *Diff) GoToTopOfFile() {

  m.GoToCrsPos_Write( 0, 0 )
}

func (m *Diff) GoToStartOfRow() {

  if( 0 < m.NumLines() ) {
    OCL := m.CrsLine() // Old cursor line

    m.GoToCrsPos_Write( OCL, m.leftChar )
  }
}

func (m *Diff) GoToEndOfRow() {

  if( 0 < m.NumLines() ) {
    pV  := m_vis.CV()
    pfb := pV.p_fb

    DL := m.CrsLine()          // Diff line
    VL := m.ViewLine( pV, DL ) // View line

    LL := pfb.LineLen( VL )
    if( 0 < LL ) {
      NCP := Min_i( LL-1, m.leftChar + m.WorkingCols( pV ) - 1 )

      m.GoToCrsPos_Write( DL, NCP )
    }
  }
}

//func (m *Diff) GoToFile() {
//  // FIXME
//}

func (m *Diff) Do_dd() {
  // FIXME
}

func (m *Diff) Do_dw() {
  // FIXME
}

func (m *Diff) Do_yy() {
  // FIXME
}

func (m *Diff) Do_yw() {
  // FIXME
}

func (m *Diff) Do_p() {
  // FIXME
}

func (m *Diff) Do_P() {
  // FIXME
}

func (m *Diff) Do_r() {
  // FIXME
}

func (m *Diff) Do_R() {
  // FIXME
}

func (m *Diff) Do_J() {
  // FIXME
}

func (m *Diff) Do_Tilda() {
  // FIXME
}

func (m *Diff) Do_u() {
  // FIXME
}

func (m *Diff) Do_U() {
  // FIXME
}

func (m *Diff) MoveCurrLineToTop() {

  if( 0<m.crsRow ) {
    m.topLine += m.crsRow
    m.crsRow = 0
    m.UpdateBV()
  }
}

func (m *Diff) MoveCurrLineCenter( write bool ) {

  pV := m_vis.CV()

  center := int(0.5*float64(m.WorkingRows(pV)) + 0.5)

  OCL := m.CrsLine() // Old cursor line

  if( 0 < OCL && OCL < center && 0 < m.topLine ) {
    // Cursor line cannot be moved to center, but can be moved closer to center
    // CrsLine(m) does not change:
    m.crsRow += m.topLine
    m.topLine = 0

    if( write ) { m.UpdateBV() }

  } else if( center < OCL && center != m.crsRow ) {
    m.topLine += m.crsRow - center
    m.crsRow = center

    if( write ) { m.UpdateBV() }
  }
}

func (m *Diff) MoveCurrLineToBottom() {

  if( 0 < m.topLine ) {
    pV := m_vis.CV()

    WR  := m.WorkingRows( pV )
    OCL := m.CrsLine() // Old cursor line

    if( WR-1 <= OCL ) {
      m.topLine -= WR - m.crsRow - 1
      m.crsRow = WR-1
      m.UpdateBV()
    } else {
      // Cursor line cannot be moved to bottom, but can be moved closer to bottom
      // CrsLine(m) does not change:
      m.crsRow += m.topLine
      m.topLine = 0;
      m.UpdateBV()
    }
  }
}

func (m *Diff) RunDiff( DA DiffArea ) {

  var t1 time.Time = time.Now()

  m.Popu_SameList( DA )
  m.Sort_SameList()
//PrintSameList();
  m.Popu_DiffList( DA )
//PrintDiffList();
  m.Popu_DI_List( DA )
//PrintDI_List( DA );

  var t2 time.Time = time.Now()
  m.diff_dur = t2.Sub( t1 )
  m.printed_diff_ms = false
}

// Find the largest SameArea in the DiffArea, da.
// Largest SameArea is determined by the most matching bytes.
//
func (m *Diff) Find_Max_Same( da DiffArea ) SameArea {

  var max_same SameArea

  for _ln_s := da.ln_s; _ln_s<da.fnl_s()-max_same.nlines; _ln_s++ {
    var ln_s int = _ln_s;
    var cur_same SameArea
    for ln_l := da.ln_l; ln_s<da.fnl_s() && ln_l<da.fnl_l(); ln_l++ {
      var ls *FLine = m.pfS.GetLP( ln_s )
      var ll *FLine = m.pfL.GetLP( ln_l )

      if( ls.Chksum() != ll.Chksum() ) { cur_same.Clear(); ln_s = _ln_s;
      } else {
        if( 0 == max_same.nlines || // First line match
            0 == cur_same.nlines ) {// First line match this outer loop
          cur_same.Init( ln_s, ln_l, ls.Len()+1 ); // Add one to account for line delimiter

        } else { // Continuation of cur_same
          cur_same.Inc( Min_i( ls.Len()+1, ll.Len()+1 ) ); // Add one to account for line delimiter
        }
        if( max_same.nbytes < cur_same.nbytes ) { max_same.Set( cur_same ); }
        ln_s++;
      }
    }
    // This line makes the diff run faster:
    if( 0 < max_same.nlines ) { _ln_s = Max_i( _ln_s, max_same.ln_s+max_same.nlines-1 ); }
  }
  return max_same;
}

func (m *Diff) Popu_SameList( DA DiffArea ) {

  m.sameList.Clear()
  var compList Vector[DiffArea]
      compList.Push( DA )

  var da DiffArea
  for compList.Pop( &da ) {
    var same SameArea = m.Find_Max_Same( da )

    if( 0<same.nlines && 0<same.nbytes ) { //< Dont count a single empty line as a same area
      m.sameList.Push( same )

      SAME_FNL_S := same.ln_s+same.nlines // Same finish line short
      SAME_FNL_L := same.ln_l+same.nlines // Same finish line long

      if( ( same.ln_s == da.ln_s || same.ln_l == da.ln_l ) &&
          SAME_FNL_S < da.fnl_s() &&
          SAME_FNL_L < da.fnl_l() ) {
        // Only one new DiffArea after same:
        ca1 := DiffArea{ SAME_FNL_S, da.fnl_s()-SAME_FNL_S,
                         SAME_FNL_L, da.fnl_l()-SAME_FNL_L }
        compList.Push( ca1 )

      } else if( ( SAME_FNL_S == da.fnl_s() || SAME_FNL_L == da.fnl_l() ) &&
                 da.ln_s < same.ln_s &&
                 da.ln_l < same.ln_l ) {
        // Only one new DiffArea before same:
        ca1 := DiffArea{ da.ln_s, same.ln_s-da.ln_s,
                         da.ln_l, same.ln_l-da.ln_l }
        compList.Push( ca1 );

      } else if( da.ln_s < same.ln_s && SAME_FNL_S < da.fnl_s() &&
                 da.ln_l < same.ln_l && SAME_FNL_L < da.fnl_l() ) {
        // Two new DiffArea's, one before same, and one after same:
        ca1 := DiffArea{ da.ln_s, same.ln_s-da.ln_s,
                         da.ln_l, same.ln_l-da.ln_l }
        ca2 := DiffArea{ SAME_FNL_S, da.fnl_s()-SAME_FNL_S,
                         SAME_FNL_L, da.fnl_l()-SAME_FNL_L }
        compList.Push( ca1 )
        compList.Push( ca2 )
      }
    }
  }
}

// Sort m.sameList from least ln_l to greatest ln_l.
// DiffArea.ln_l is beginning line number in long file.
//
func (m *Diff) Sort_SameList() {

  SLL := m.sameList.Len()

  for k:=0; k<SLL; k++ {
    for j:=SLL-1; k<j; j-- {
      var sa0 SameArea = m.sameList.Get( j-1 )
      var sa1 SameArea = m.sameList.Get( j   )

      if( sa1.ln_l < sa0.ln_l ) {
        m.sameList.Set( j-1, sa1 )
        m.sameList.Set( j  , sa0 )
      }
    }
  }
}

func (m *Diff) Popu_DiffList( DA DiffArea ) {

  m.diffList.Clear()

  m.Popu_DiffList_Begin( DA )

  SLL := m.sameList.Len()

  for k:=1; k<SLL; k++ {
    var sa0 SameArea = m.sameList.Get( k-1 )
    var sa1 SameArea = m.sameList.Get( k   )

    da_ln_s := sa0.ln_s+sa0.nlines;
    da_ln_l := sa0.ln_l+sa0.nlines;

    da := DiffArea{ da_ln_s, sa1.ln_s - da_ln_s,
                    da_ln_l, sa1.ln_l - da_ln_l }

    m.diffList.Push( da )
  }
  m.Popu_DiffList_End( DA )
}

func (m *Diff) Popu_DiffList_Begin( DA DiffArea ) {
  // Add DiffArea before first SameArea if needed:
  if( 0 < m.sameList.Len() ) {
    var sa SameArea = m.sameList.Get( 0 )

    nlines_s_da := sa.ln_s - DA.ln_s // Num lines in short diff area
    nlines_l_da := sa.ln_l - DA.ln_l // Num lines in long  diff area

    if( 0 < nlines_s_da || 0 < nlines_l_da ) {
      // DiffArea at beginning of DiffArea:
      da := DiffArea{ DA.ln_s, nlines_s_da, DA.ln_l, nlines_l_da }
      m.diffList.Push( da )
    }
  }
}

func (m *Diff) Popu_DiffList_End( DA DiffArea ) {

  SLL := m.sameList.Len()

  if( 0 < SLL ) { // Add DiffArea after last SameArea if needed:
    var sa SameArea = m.sameList.Get( SLL-1 )
    sa_s_end := sa.ln_s + sa.nlines
    sa_l_end := sa.ln_l + sa.nlines

    if( sa_s_end < DA.fnl_s() ||
        sa_l_end < DA.fnl_l() ) { // DiffArea at end of file:
      // Number of lines of short and long equal to
      // start of SameArea short and long
      da := DiffArea{ sa_s_end, DA.fnl_s() - sa_s_end,
                      sa_l_end, DA.fnl_l() - sa_l_end }
      m.diffList.Push( da )
    }
  } else { // No SameArea, so whole DiffArea is a DiffArea:
    m.diffList.Push( DA )
  }
}

func (m *Diff) Popu_DI_List( CA DiffArea ) {

  SLL := m.sameList.Len()
  DLL := m.diffList.Len()

  if       ( SLL == 0 ) { m.Popu_DI_List_NoSameArea()
  } else if( DLL == 0 ) { m.Popu_DI_List_NoDiffArea()
  } else                { m.Popu_DI_List_DiffAndSame( CA )
  }
}

func (m *Diff) Popu_DI_List_NoSameArea() {

  // Should only be one DiffArea, which is the whole DiffArea:
  // DLL := m.diffList.Len()
  // ASSERT( DLL==1 )

  m.Popu_DI_List_AddDiff( m.diffList.Get(0) )
}

func (m *Diff) Popu_DI_List_AddDiff( da DiffArea ) {

  if( da.nlines_s < da.nlines_l ) {
    m.Popu_DI_List_AddDiff_Common( da,
                                   &m.DI_List_S,
                                   &m.DI_List_L,
                                   m.pfS, m.pfL )

  } else if( da.nlines_l < da.nlines_s ) {
    // Since the long file has a shorter DiffArea than the short file,
    // pass in a reversed DiffArea:
    da_r := DiffArea{ da.ln_l, da.nlines_l, da.ln_s, da.nlines_s }

    m.Popu_DI_List_AddDiff_Common( da_r,
                                   &m.DI_List_L,
                                   &m.DI_List_S,
                                   m.pfL, m.pfS )
  } else { // da.nlines_s == da.nlines_l
    for k:=0; k<da.nlines_l; k++ {
      var ls *FLine = m.pfS.GetLP( da.ln_s+k )
      var ll *FLine = m.pfL.GetLP( da.ln_l+k )

      var li_s LineInfo
      var li_l LineInfo

      m.Compare_Lines( ls, &li_s, ll, &li_l )

      dis := Diff_Info{ DT_CHANGED, da.ln_s+k, &li_s }
      dil := Diff_Info{ DT_CHANGED, da.ln_l+k, &li_l }

      m.DI_List_S.Push( dis )
      m.DI_List_L.Push( dil )
    }
  }
}

func (m *Diff) Popu_DI_List_NoDiffArea() {
  // Should only be one SameArea, which is the whole DiffArea:
  // SLL := m.sameList.Len()
  // ASSERT( 1 == SLL )
  m.Popu_DI_List_AddSame( m.sameList.Get(0) )
}

func (m *Diff) Popu_DI_List_DiffAndSame( DA DiffArea ) {

  SLL := m.sameList.Len()
  DLL := m.diffList.Len()

  var da DiffArea = m.diffList.Get( 0 )

  if( DA.ln_s==da.ln_s && DA.ln_l==da.ln_l ) {
    // Start with DiffArea, and then alternate between SameArea and DiffArea.
    // There should be at least as many DiffArea's as SameArea's.
    // ASSERT( SLL<=DLL )
    for k:=0; k<SLL; k++ {
      var da DiffArea = m.diffList.Get( k ); m.Popu_DI_List_AddDiff( da )
      var sa SameArea = m.sameList.Get( k ); m.Popu_DI_List_AddSame( sa )
    }
    if( SLL < DLL ) {
      // ASSERT( SLL+1==DLL )
      var da DiffArea = m.diffList.Get( DLL-1 ); m.Popu_DI_List_AddDiff( da )
    }
  } else {
    // Start with SameArea, and then alternate between DiffArea and SameArea.
    // There should be at least as many SameArea's as DiffArea's.
    // ASSERT( DLL<=SLL )
    for k:=0; k<DLL; k++ {
      var sa SameArea = m.sameList.Get( k ); m.Popu_DI_List_AddSame( sa )
      var da DiffArea = m.diffList.Get( k ); m.Popu_DI_List_AddDiff( da )
    }
    if( DLL < SLL ) {
      // ASSERT( DLL+1==SLL )
      var sa SameArea = m.sameList.Get( SLL-1 ); m.Popu_DI_List_AddSame( sa )
    }
  }
}

func (m *Diff) Popu_DI_List_AddDiff_Common( da DiffArea,
                                            p_DI_List_s, p_DI_List_l *Vector[Diff_Info],
                                            pfs, pfl *FileBuf ) {
  var simiList Vector[SimLines]

  m.Popu_SimiList( da, pfs, pfl, &simiList )

  m.Sort_SimiList( &simiList )
//PrintSimiList();

  m.SimiList_2_DI_Lists( da, p_DI_List_s, p_DI_List_l, &simiList )
}

// Returns number of bytes that are the same between the two lines
// and fills in li_s and li_l
func (m *Diff) Compare_Lines( ls *FLine, li_s *LineInfo,
                              ll *FLine, li_l *LineInfo ) int {
  if( 0==ls.Len() && 0==ll.Len() ) { return 1 }

  li_s.Clear(); li_l.Clear()
  var pls *FLine = ls; var pli_s *LineInfo = li_s;
  var pll *FLine = ll; var pli_l *LineInfo = li_l;
  if( ll.Len() < ls.Len() ) { pls = ll; pli_s = li_l;
                              pll = ls; pli_l = li_s; }
  SLL := pls.Len()
  LLL := pll.Len()

  pli_l.SetLen( LLL )
  pli_s.SetLen( LLL )

  num_same := 0
  i_s := 0
  i_l := 0

  for i_s < SLL && i_l < LLL {
    cs := pls.GetR( i_s )
    cl := pll.GetR( i_l )

    if( cs == cl ) {
      num_same++
      pli_s.Set( i_s, DT_SAME ); i_s++
      pli_l.Set( i_l, DT_SAME ); i_l++
    } else {
      remaining_s := SLL - i_s
      remaining_l := LLL - i_l

      if( 0<remaining_s &&
          0<remaining_l &&
          remaining_s == remaining_l ) {

        pli_s.Set( i_s, DT_CHANGED ); i_s++
        pli_l.Set( i_l, DT_CHANGED ); i_l++

      } else if( remaining_s < remaining_l ) { pli_l.Set( i_l, DT_INSERTED ); i_l++
      } else if( remaining_l < remaining_s ) { pli_s.Set( i_s, DT_INSERTED ); i_s++
      }
    }
  }
  for k:=SLL; k<LLL; k++ { pli_s.Set( k, DT_DELETED )  }
  for k:=i_l; k<LLL; k++ { pli_l.Set( k, DT_INSERTED ) }

  return num_same
}

// Returns true if the two lines, line_s and line_l, in the two files
// being compared, are the names of files that differ
func (m *Diff) Popu_DI_List_Have_Diff_Files( line_s, line_l int ) bool {
  files_differ := false

  if( m.pfS.is_dir && m.pfL.is_dir ) {
    // fname_s and fname_l are head names
    var fname_s string = m.pfS.GetLP( line_s ).to_str()
    var fname_l string = m.pfL.GetLP( line_l ).to_str()

    if( (fname_s != "..") && !strings.HasSuffix( fname_s, DIR_DELIM_S ) &&
        (fname_l != "..") && !strings.HasSuffix( fname_l, DIR_DELIM_S ) ) {
      // fname_s and fname_l should now be full path names,
      // tail and head, of regular files
      fname_s =  m.pfS.dir_name + fname_s
      fname_l =  m.pfL.dir_name + fname_l

      var pfb_s *FileBuf = m_vis.GetFileBuf_s( fname_s )
      var pfb_l *FileBuf = m_vis.GetFileBuf_s( fname_l )

      // If one side is in ram, read in the other side:
      if       ( (nil == pfb_s) && (nil != pfb_l) ) { m_vis.NotHaveFileAddFile( fname_s );
      } else if( (nil != pfb_s) && (nil == pfb_l) ) { m_vis.NotHaveFileAddFile( fname_l );
      } else if( (nil == pfb_s) && (nil == pfb_l) ) {
        // Adding files is slow because of all the new'ing, so limit
        // the number of files that can be added per diff:
        if( m.num_files_added_this_diff < max_files_added_per_diff ) {
          var added_s bool = m_vis.NotHaveFileAddFile( fname_s )
          var added_l bool = m_vis.NotHaveFileAddFile( fname_l )

          if( added_s ) { m.num_files_added_this_diff++ }
          if( added_l ) { m.num_files_added_this_diff++ }
        }
      }
      pfb_s = m_vis.GetFileBuf_s( fname_s );
      pfb_l = m_vis.GetFileBuf_s( fname_l );

      if( (nil == pfb_s) || (nil == pfb_l) ) {
        // Slow: Compare the files in NVM:
        same,_ := Files_Are_Same_p( fname_s, fname_l )
        files_differ = !same
      } else {
        // Fast: Compare files already cached in memory:
        same := Files_Are_Same_o( pfb_s, pfb_l )
        files_differ = !same
      }
    }
  }
  return files_differ;
}

func (m *Diff) Popu_DI_List_AddSame( sa SameArea ) {

  for k:=0; k<sa.nlines; k++ {
    var DT Diff_Type = DT_SAME
    if( m.Popu_DI_List_Have_Diff_Files( sa.ln_s+k, sa.ln_l+k ) ) {
      DT = DT_DIFF_FILES
    }
    dis := Diff_Info{ DT, sa.ln_s+k, nil }
    dil := Diff_Info{ DT, sa.ln_l+k, nil }

    m.DI_List_S.Push( dis )
    m.DI_List_L.Push( dil )
  }
}

func (m *Diff) Popu_SimiList( da DiffArea,
                              pfs, pfl *FileBuf,
                              p_simiList *Vector[SimLines] ) {
//m.simiList.Clear()

  if( 0<da.nlines_s && 0<da.nlines_l ) {
    ca := da

    var compList Vector[DiffArea]
    compList.Push( ca )

    for (p_simiList.Len() < da.nlines_s) && compList.Pop( &ca ) {

      var siml SimLines = m.Find_Lines_Most_Same( ca, pfs, pfl )
      p_simiList.Push( siml )

      if( ( siml.ln_s == ca.ln_s || siml.ln_l == ca.ln_l ) &&
          siml.ln_s+1 < ca.fnl_s() &&
          siml.ln_l+1 < ca.fnl_l() ) {
        // Only one new DiffArea after siml:
        ca1 := DiffArea{ siml.ln_s+1, ca.fnl_s()-siml.ln_s-1,
                         siml.ln_l+1, ca.fnl_l()-siml.ln_l-1 }
        compList.Push( ca1 )

      } else if( ( siml.ln_s+1 == ca.fnl_s() || siml.ln_l+1 == ca.fnl_l() ) &&
                 ca.ln_s < siml.ln_s &&
                 ca.ln_l < siml.ln_l ) {
        // Only one new DiffArea before siml:
        ca1 := DiffArea{ ca.ln_s, siml.ln_s-ca.ln_s,
                         ca.ln_l, siml.ln_l-ca.ln_l }
        compList.Push( ca1 )

      } else if( ca.ln_s < siml.ln_s && siml.ln_s+1 < ca.fnl_s() &&
                 ca.ln_l < siml.ln_l && siml.ln_l+1 < ca.fnl_l() ) {
        // Two new DiffArea's, one before siml, and one after siml:
        ca1 := DiffArea{ ca.ln_s, siml.ln_s-ca.ln_s,
                         ca.ln_l, siml.ln_l-ca.ln_l }
        ca2 := DiffArea{ siml.ln_s+1, ca.fnl_s()-siml.ln_s-1,
                         siml.ln_l+1, ca.fnl_l()-siml.ln_l-1 }
        compList.Push( ca1 )
        compList.Push( ca2 )
      }
    }
  }
}

func (m *Diff) Sort_SimiList( p_simiList *Vector[SimLines] ) {

  SLL := p_simiList.Len()

  for k:=0; k<SLL; k++ {
    for j:=SLL-1; k<j; j-- {
      var sl0 SimLines = p_simiList.Get( j-1 )
      var sl1 SimLines = p_simiList.Get( j   )

      if sl1.ln_l < sl0.ln_l {
        p_simiList.Set( j-1, sl1 )
        p_simiList.Set( j  , sl0 )
      }
    }
  }
}

func (m *Diff) SimiList_2_DI_Lists( da DiffArea,
                                    p_DI_List_s, p_DI_List_l *Vector[Diff_Info],
                                    p_simiList *Vector[SimLines] ) {
  // dis_ln = Diff info short line number.
  // Diff_Info.diff_type on the short side defaults to DT_DELETED.
  // A deleted line does not have a line number,
  //   so use the line number of the previous line:
  dis_ln := 0; if( 0<da.ln_s ) { dis_ln = da.ln_s-1 }

  for ln_l:=da.ln_l; ln_l<da.fnl_l(); ln_l++ {
    dis := Diff_Info{ DT_DELETED , dis_ln, nil }
    dil := Diff_Info{ DT_INSERTED, ln_l  , nil }

    // j start. Used for loop optimization
    j_st := 0
    for j:=j_st; j<p_simiList.Len(); j++ {
      var p_siml *SimLines = p_simiList.GetP( j )

      if( p_siml.ln_l == ln_l ) {
        // The Diff_Info.diff_type on the short side is being set to DT_CHANGED,
        // so a line numbe can be assigned:
        dis.line_num  = p_siml.ln_s
        dis.diff_type = DT_CHANGED
        dis.pLineInfo = p_siml.li_s; p_siml.li_s = nil; // Transfer ownership of LineInfo from p_siml to dis

        dil.diff_type = DT_CHANGED
        dil.pLineInfo = p_siml.li_l; p_siml.li_l = nil; // Transfer ownership of LineInfo from p_siml to dil

        dis_ln = dis.line_num
        j_st = j + 1
        break
      }
    }
    // DI_List_s and DI_List_l now own LineInfo objects:
    p_DI_List_s.Push( dis )
    p_DI_List_l.Push( dil )
  }
}

//func (m *Diff) Find_Lines_Most_Same( ca DiffArea, pfs, pfl *FileBuf ) SimLines {
//
//  // LD = Length Difference between long area and short area
//  var LD int = ca.nlines_l - ca.nlines_s
//
//  most_same := SimLines{ 0, 0, 0, nil, nil }
//  for ln_s := ca.ln_s; ln_s<ca.fnl_s(); ln_s++ {
//    var ST_L int = ca.ln_l+(ln_s-ca.ln_s)
//
////  for ln_l := ST_L; ln_l<ca.fnl_l() && ln_l<ST_L+LD+1 && ln_l<pfl.NumLines(); ln_l++ {}
//    for ln_l := ST_L; ln_l<ca.fnl_l() && ln_l<ST_L+LD+1; ln_l++ {
//
//      var ls *FLine = pfs.GetLP( ln_s ) // Line from short area
////Log( fmt.Sprintf("Find_Lines_Most_Same: ln_l=%v, pfl.NumLines()=%v, ca.fnl_l()=%v, ST_L+LD+1=%v",
////                                        ln_l,    pfl.NumLines(),    ca.fnl_l(),    ST_L+LD+1) )
//      var ll *FLine = pfl.GetLP( ln_l ) // Line from long  area
//
//      var li_s LineInfo
//      var li_l LineInfo
//      var bytes_same int = m.Compare_Lines( ls, &li_s, ll, &li_l )
//
//      if( most_same.nbytes < bytes_same ) {
//        most_same.ln_s   = ln_s
//        most_same.ln_l   = ln_l
//        most_same.nbytes = bytes_same
//        most_same.li_s   = &li_s; // Hand off li_s
//        most_same.li_l   = &li_l; // and      li_l
//      }
//    }
//  }
//  if( 0==most_same.nbytes ) {
//    // This if() block ensures that each line in the short DiffArea is
//    // matched to a line in the long DiffArea.  Each line in the short
//    // DiffArea must be matched to a line in the long DiffArea or else
//    // SimiList_2_DI_Lists wont work right.
//    most_same.ln_s   = ca.ln_s
//    most_same.ln_l   = ca.ln_l
//    most_same.nbytes = 1
//  }
//  return most_same
//}

func (m *Diff) Find_Lines_Most_Same( ca DiffArea, pfs, pfl *FileBuf ) SimLines {

  most_same := SimLines{ 0, 0, 0, nil, nil }
  for ln_s := ca.ln_s; ln_s<ca.fnl_s(); ln_s++ {

    for ln_l := ca.ln_l; ln_l<ca.fnl_l(); ln_l++ {

      var ls *FLine = pfs.GetLP( ln_s ) // Line from short area
//Log( fmt.Sprintf("Find_Lines_Most_Same: ln_l=%v, pfl.NumLines()=%v, ca.fnl_l()=%v, ST_L+LD+1=%v",
//                                        ln_l,    pfl.NumLines(),    ca.fnl_l(),    ST_L+LD+1) )
      var ll *FLine = pfl.GetLP( ln_l ) // Line from long  area

      var li_s LineInfo
      var li_l LineInfo
      var bytes_same int = m.Compare_Lines( ls, &li_s, ll, &li_l )

      if( most_same.nbytes < bytes_same ) {
        most_same.ln_s   = ln_s
        most_same.ln_l   = ln_l
        most_same.nbytes = bytes_same
        most_same.li_s   = &li_s; // Hand off li_s
        most_same.li_l   = &li_l; // and      li_l
      }
    }
  }
  if( 0==most_same.nbytes ) {
    // This if() block ensures that each line in the short DiffArea is
    // matched to a line in the long DiffArea.  Each line in the short
    // DiffArea must be matched to a line in the long DiffArea or else
    // SimiList_2_DI_Lists wont work right.
    most_same.ln_s   = ca.ln_s
    most_same.ln_l   = ca.ln_l
    most_same.nbytes = 1
  }
  return most_same
}

func (m *Diff) GoToCrsPos_Write( ncp_crsLine, ncp_crsChar int ) {

  pV := m_vis.CV()

  OCL := m.CrsLine()
  OCP := m.CrsChar()
  NCL := ncp_crsLine
  NCP := ncp_crsChar

  if( OCL == NCL && OCP == NCP ) {
    // Not moving to new cursor line so just put cursor back where is was
    m.PrintCursor( pV )
  } else {
    if( m.inVisualMode ) {
      m.v_fn_line = NCL
      m.v_fn_char = NCP
    }
    // These moves refer to View of buffer:
    MOVE_DOWN  := m.BotLine( pV )   < NCL
    MOVE_RIGHT := m.RightChar( pV ) < NCP
    MOVE_UP    := NCL < m.topLine
    MOVE_LEFT  := NCP < m.leftChar

    redraw := MOVE_DOWN || MOVE_RIGHT || MOVE_UP || MOVE_LEFT

    if( redraw ) {
      if       ( MOVE_DOWN ) { m.topLine = NCL - m.WorkingRows( pV ) + 1
      } else if( MOVE_UP )   { m.topLine = NCL
      }
      if       ( MOVE_RIGHT ) { m.leftChar = NCP - m.WorkingCols( pV ) + 1
      } else if( MOVE_LEFT  ) { m.leftChar = NCP
      }
      // crsRow and crsCol must be set to new values before calling CalcNewCrsByte
      m.crsRow = NCL - m.topLine
      m.crsCol = NCP - m.leftChar

      m.UpdateBV()

    } else if( m.inVisualMode ) {
      if( m.inVisualBlock ) { m.GoToCrsPos_Write_VisualBlock( OCL, OCP, NCL, NCP )
      } else                { m.GoToCrsPos_Write_Visual     ( OCL, OCP, NCL, NCP )
      }
    } else {
      // crsRow and crsCol must be set to new values before calling CalcNewCrsByte and PrintCursor
      m.crsRow = NCL - m.topLine
      m.crsCol = NCP - m.leftChar

      m.PrintStsLine( m.pvS )
      m.PrintStsLine( m.pvL )
      m.PrintCursor( pV )  // Put cursor into position.
    }
  }
}

func (m *Diff) GoToCrsPos_NoWrite( ncp_crsLine, ncp_crsChar int ) {

  pV := m_vis.CV()

  // These moves refer to View of buffer:
  MOVE_DOWN  := m.BotLine( pV )   < ncp_crsLine
  MOVE_RIGHT := m.RightChar( pV ) < ncp_crsChar
  MOVE_UP    := ncp_crsLine < m.topLine
  MOVE_LEFT  := ncp_crsChar < m.leftChar

  if     ( MOVE_DOWN ) { m.topLine = ncp_crsLine - m.WorkingRows( pV ) + 1
  } else if( MOVE_UP ) { m.topLine = ncp_crsLine
  }
  m.crsRow = ncp_crsLine - m.topLine

  if       ( MOVE_RIGHT ) { m.leftChar = ncp_crsChar - m.WorkingCols( pV ) + 1
  } else if( MOVE_LEFT  ) { m.leftChar = ncp_crsChar
  }
  m.crsCol = ncp_crsChar - m.leftChar
}

func (m *Diff) GoToCrsPos_Write_VisualBlock( OCL, OCP, NCL, NCP int ) {
  // FIXME:
}

func (m *Diff) GoToCrsPos_Write_Visual( OCL, OCP, NCL, NCP int ) {
  // FIXME:
}

func (m *Diff) Do_n_Diff( write bool ) {

  if( 0 < m.NumLines() ) {
  //m.Set_Cmd_Line_Msg("Searching down for diff");

    dl := m.CrsLine() // Diff line, changed by search methods below

    pV := m_vis.CV()

    p_DI_List := m.View_2_DI_List_C( pV )

    var DT Diff_Type = p_DI_List.Get(dl).diff_type // Current diff type

    found_same := true

    if( DT == DT_CHANGED ||
        DT == DT_INSERTED ||
        DT == DT_DELETED ||
        DT == DT_DIFF_FILES ) {
      // If currently on a diff, search for same before searching for diff
      found_same = m.Do_n_Search_for_Same( &dl, p_DI_List )
    }
    if( found_same ) {
      found_diff := m.Do_n_Search_for_Diff( &dl, p_DI_List )

      var NCL, NCP int
      if( found_diff ) {
        NCL = dl;
        NCP = m.Do_n_Find_Crs_Pos( NCL, p_DI_List )
      } else { // Could not find a difference.
               // Check if one file ends in LF and the other does not:
        if( m.pfS.lines.LF_at_EOF != m.pfL.lines.LF_at_EOF ) {
          found_diff = true
          NCL = p_DI_List.Len() - 1
          NCP = pV.p_fb.LineLen( p_DI_List.Get( NCL ).line_num )
        }
      }
      if( found_diff ) {
        if( write ) { m.GoToCrsPos_Write( NCL, NCP )
        } else      { m.GoToCrsPos_NoWrite( NCL, NCP )
        }
      }
    }
  }
}

func (m *Diff) Do_N_Diff() {

  if( 0 < m.NumLines() ) {
  //m.diff.Set_Cmd_Line_Msg("Searching up for diff");

    dl := m.CrsLine() // Diff line, changed by search methods below

    pV := m_vis.CV()

    p_DI_List := m.View_2_DI_List_C( pV )

    var DT Diff_Type = p_DI_List.Get(dl).diff_type // Current diff type

    found_same := true

    if( DT == DT_CHANGED ||
        DT == DT_INSERTED ||
        DT == DT_DELETED ||
        DT == DT_DIFF_FILES ) {
      // If currently on a diff, search for same before searching for diff
      found_same = m.Do_N_Search_for_Same( &dl, p_DI_List )
    }
    if( found_same ) {
      found_diff := m.Do_N_Search_for_Diff( &dl, p_DI_List )

      if( found_diff ) {
        NCL := dl
        NCP := m.Do_n_Find_Crs_Pos( NCL, p_DI_List )

        m.GoToCrsPos_Write( NCL, NCP )
      }
    }
  }
}

func (m *Diff) Do_n_Pattern() {

  pV := m_vis.CV();

  NUM_LINES := pV.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
  //String msg("/");
  //m.diff.Set_Cmd_Line_Msg( msg += m_vis.GetRegex() );

    ncp := CrsPos{ 0, 0 } // Next cursor position

    if( m.Do_n_FindNextPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    }
  }
}

func (m *Diff) Do_N_Pattern() {

  pV := m_vis.CV();

  NUM_LINES := pV.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
  //String msg("/");
  //m.diff.Set_Cmd_Line_Msg( msg += m_vis.GetRegex() );

    ncp := CrsPos{ 0, 0 } // Next cursor position

    if( m.Do_N_FindPrevPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    }
  }
}

func (m *Diff) Do_n_Search_for_Same( p_dl *int,
                                     p_DI_List *Vector[Diff_Info] ) bool {

  NUM_LINES := m.NumLines()
  dl_st := *p_dl

  // Search forward for DT_SAME
  found := false

  if( 1 < NUM_LINES ) {
    for !found && *p_dl<NUM_LINES {
      var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

      if( DT == DT_SAME ) { found = true
      } else              { *p_dl++
      }
    }
    if( !found ) {
      // Wrap around back to top and search again:
      *p_dl = 0
      for( !found && *p_dl<dl_st ) {
        var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

        if( DT == DT_SAME ) { found = true
        } else              { *p_dl++
        }
      }
    }
  }
  return found
}

func (m *Diff) Do_n_Search_for_Diff( p_dl *int,
                                     p_DI_List *Vector[Diff_Info] ) bool {

  dl_st := *p_dl

  // Search forward for non-DT_SAME
  found_diff := false

  if( 1 < m.NumLines() ) {
    found_diff = m.Do_n_Search_for_Diff_DT( p_dl, p_DI_List )

    if( !found_diff ) {
      *p_dl = dl_st
      found_diff = m.Do_n_Search_for_Diff_WhiteSpace( p_dl, p_DI_List )
    }
  }
  return found_diff
}

func (m *Diff) Do_n_Find_Crs_Pos( NCL int,
                                  p_DI_List *Vector[Diff_Info] ) int {
  NCP := 0

  var DT_new Diff_Type = p_DI_List.Get( NCL ).diff_type

  if( DT_new == DT_CHANGED ) {
    var pLI_s *LineInfo = m.DI_List_S.Get( NCL ).pLineInfo
    var pLI_l *LineInfo = m.DI_List_L.Get( NCL ).pLineInfo

    for k:=0; nil != pLI_s && k<pLI_s.Len() &&
              nil != pLI_l && k<pLI_l.Len(); k++ {

      var dt_s Diff_Type = pLI_s.Get( k )
      var dt_l Diff_Type = pLI_l.Get( k )

      if( dt_s != DT_SAME || dt_l != DT_SAME ) {
        NCP = k
        break
      }
    }
  }
  return NCP
}

func (m *Diff) Do_N_Search_for_Same( p_dl *int,
                                     p_DI_List *Vector[Diff_Info] ) bool {

  NUM_LINES := m.NumLines()
  dl_st := *p_dl

  // Search backwards for DT_SAME
  found := false

  if( 1 < NUM_LINES ) {
    for !found && 0<=*p_dl {
      var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

      if( DT == DT_SAME ) { found = true
      } else              { *p_dl--
      }
    }
    if( !found ) {
      // Wrap around back to bottom and search again:
      *p_dl = NUM_LINES-1
      for !found && dl_st<*p_dl  {
        var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

        if( DT == DT_SAME ) { found = true
        } else              { *p_dl--
        }
      }
    }
  }
  return found
}

func (m *Diff) Do_N_Search_for_Diff( p_dl *int,
                                     p_DI_List *Vector[Diff_Info] ) bool {

  dl_st := *p_dl

  // Search backwards for non-DT_SAME
  found_diff := false

  if( 1 < m.NumLines() ) {
    found_diff = m.Do_N_Search_for_Diff_DT( p_dl, p_DI_List )

    if( !found_diff ) {
      *p_dl = dl_st
      found_diff = m.Do_N_Search_for_Diff_WhiteSpace( p_dl, p_DI_List )
    }
  }
  return found_diff;
}

func (m *Diff) Do_n_FindNextPattern( p_ncp *CrsPos ) bool {

  pV := m_vis.CV()
  pfb := pV.p_fb

  NUM_LINES := pfb.NumLines()

  OCL := m.CrsLine() // Diff line
  OCC := m.CrsChar()

  OCLv := m.ViewLine( pV, OCL ) // View line

  st_l := OCLv
  st_c := OCC

  found_next_star := false

  // Move past current star:
  LL := pfb.LineLen( OCLv )

  pfb.Check_4_New_Regex()
  pfb.Find_Regexs_4_Line( OCL )

  // Move past current pattern:
  for ; st_c<LL && pV.InStarOrStarInF(OCLv,st_c); st_c++ {
  }
  // If at end of current line, go down to next line
  if( LL <= st_c ) { st_c=0; st_l++ }

  // Search for first pattern position past current position
  for l:=st_l; !found_next_star && l<NUM_LINES; l++ {
    pfb.Find_Regexs_4_Line( l )

    LL := pfb.LineLen( l )

    for p:=st_c; !found_next_star && p<LL; p++ {
      if( pV.InStarOrStarInF(l,p) ) {
        found_next_star = true;
        // Convert from view line back to diff line:
        dl := m.DiffLine( pV, l )
        p_ncp.crsLine = dl
        p_ncp.crsChar = p
      }
    }
    // After first line, always start at beginning of line
    st_c = 0
  }
  // Near end of file and did not find any patterns, so go to first pattern in file
  if( !found_next_star ) {
    for l:=0; !found_next_star && l<=OCLv; l++ {
      pfb.Find_Regexs_4_Line( l )

      LL := pfb.LineLen( l )
      END_C := LL
      if( OCLv==l ) { END_C = Min_i( OCC, LL ) }

      for p:=0; !found_next_star && p<END_C; p++ {
        if( pV.InStarOrStarInF(l,p) ) {
          found_next_star = true;
          // Convert from view line back to diff line:
          dl := m.DiffLine( pV, l )
          p_ncp.crsLine = dl
          p_ncp.crsChar = p
        }
      }
    }
  }
  return found_next_star
}

func (m *Diff) Do_N_FindPrevPattern( p_ncp *CrsPos ) bool {

  m.MoveInBounds_Line()

  pV := m_vis.CV()
  pfb := pV.p_fb

  NUM_LINES := pfb.NumLines()

  OCL := m.CrsLine()
  OCC := m.CrsChar()

  OCLv := m.ViewLine( pV, OCL ) // View line

  pfb.Check_4_New_Regex()

  found_prev_star := false

  // Search for first star position before current position
  for l:=OCLv; !found_prev_star && 0<=l; l-- {
    pfb.Find_Regexs_4_Line( l )

    LL := pfb.LineLen( l )

    p := LL-1
    if( OCLv==l ) {
      p = 0
      if( 0<OCC ) { p = OCC-1 }
    }
    for ; 0<p && !found_prev_star; p-- {
      for ; 0<=p && pV.InStarOrStarInF(l,p); p-- {
        found_prev_star = true
        // Convert from view line back to diff line:
        dl := m.DiffLine( pV, l )
        p_ncp.crsLine = dl
        p_ncp.crsChar = p
      }
    }
  }
  // Near beginning of file and did not find any patterns, so go to last pattern in file
  if( !found_prev_star ) {
    for l:=NUM_LINES-1; !found_prev_star && OCLv<l; l-- {
      pfb.Find_Regexs_4_Line( l )

      LL := pfb.LineLen( l )

      p := LL-1
      if( OCLv==l ) {
        p = 0
        if( 0<OCC ) { p = OCC-1 }
      }
      for ; 0<p && !found_prev_star; p-- {
        for ; 0<=p && pV.InStarOrStarInF(l,p); p-- {
          found_prev_star = true;
          // Convert from view line back to diff line:
          dl := m.DiffLine( pV, l )
          p_ncp.crsLine = dl
          p_ncp.crsChar = p
        }
      }
    }
  }
  return found_prev_star
}

// Look for difference based on Diff_Info:
func (m *Diff) Do_n_Search_for_Diff_DT( p_dl *int,
                                        p_DI_List *Vector[Diff_Info] ) bool {
  found_diff := false

  NUM_LINES := m.NumLines()
  dl_st := *p_dl

  for !found_diff && *p_dl<NUM_LINES {
    var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

    if( DT == DT_CHANGED ||
        DT == DT_INSERTED ||
        DT == DT_DELETED ||
        DT == DT_DIFF_FILES ) { found_diff = true
    } else                    { *p_dl++
    }
  }
  if( !found_diff ) {
    // Wrap around back to top and search again:
    *p_dl = 0
    for !found_diff && *p_dl<dl_st {
      var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

      if( DT == DT_CHANGED ||
          DT == DT_INSERTED ||
          DT == DT_DELETED ||
          DT == DT_DIFF_FILES ) { found_diff = true
      } else                    { *p_dl++
      }
    }
  }
  return found_diff
}

// Look for difference in white space at beginning or ending of lines:
func (m *Diff) Do_n_Search_for_Diff_WhiteSpace( p_dl *int,
                                                p_DI_List *Vector[Diff_Info] ) bool {

  found_diff := false

  NUM_LINES := m.NumLines()

  p_DI_List_o := &m.DI_List_S
  pF_o := m.pfS
  pF_m := m.pfL

  if( p_DI_List == &m.DI_List_S ) {
    p_DI_List_o = &m.DI_List_L
    pF_o = m.pfL
    pF_m = m.pfS
  }
  // If the current line has a difference in white space at beginning or end, start
  // searching on next line so the current line number is not automatically returned.
  var curr_line_has_LT_WS_diff bool =
    m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, *p_dl, p_DI_List, p_DI_List_o, pF_m, pF_o )
  dl_st := *p_dl
  if( curr_line_has_LT_WS_diff ) { dl_st = (*p_dl + 1) % NUM_LINES }

  // Search from dl_st to end for lines of different length:
  for k:=dl_st; !found_diff && k<NUM_LINES; k++ {
    found_diff = m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, k, p_DI_List, p_DI_List_o, pF_m, pF_o )
  }
  if( !found_diff ) {
    // Search from top to dl_st for lines of different length:
    for k:=0; !found_diff && k<dl_st; k++ {
      found_diff = m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, k, p_DI_List, p_DI_List_o, pF_m, pF_o )
    }
  }
  return found_diff;
}

func (m *Diff) Line_Has_Leading_or_Trailing_WS_Diff( p_dl *int,
                                                     k int,
                                                     p_DI_List, p_DI_List_o *Vector[Diff_Info],
                                                     pF_m, pF_o *FileBuf ) bool {
  L_T_WS_diff := false

  var Di_m Diff_Info = p_DI_List.Get( k )
  var Di_o Diff_Info = p_DI_List_o.Get( k )

  if( Di_m.diff_type == DT_SAME && Di_o.diff_type == DT_SAME ) {
    var p_lm *FLine = pF_m.GetLP( Di_m.line_num ) // Line from my    view
    var p_lo *FLine = pF_o.GetLP( Di_o.line_num ) // Line from other view

    if( p_lm.Len() != p_lo.Len() ) {
      L_T_WS_diff = true
      *p_dl = k
    }
  }
  return L_T_WS_diff
}

// If past end of line, move back to end of line.
// Returns true if moved, false otherwise.
//
func (m *Diff) MoveInBounds_Line() {

  pV := m_vis.CV()

  DL  := m.CrsLine()  // Diff line
  VL  := m.ViewLine( pV, DL )      // View line
  LL  := pV.p_fb.LineLen( VL )
  EOL := 0; if( 0<LL ) { EOL = LL-1 }

  if( EOL < m.CrsChar() ) { // Since cursor is now allowed past EOL,
                            // it may need to be moved back:
    m.GoToCrsPos_NoWrite( DL, EOL )
  }
}

// Look for difference based on Diff_Info:
func (m *Diff) Do_N_Search_for_Diff_DT( p_dl *int,
                                        p_DI_List *Vector[Diff_Info] ) bool {
  found_diff := false

  NUM_LINES := m.NumLines()
  dl_st := *p_dl

  for( !found_diff && 0<=*p_dl ) {
    var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

    if( DT == DT_CHANGED ||
        DT == DT_INSERTED ||
        DT == DT_DELETED ||
        DT == DT_DIFF_FILES ) { found_diff = true
    } else                    { *p_dl--
    }
  }
  if( !found_diff ) {
    // Wrap around back to bottom and search again:
    *p_dl = NUM_LINES-1
    for !found_diff && dl_st<*p_dl {
      var DT Diff_Type = p_DI_List.Get(*p_dl).diff_type

      if( DT == DT_CHANGED ||
          DT == DT_INSERTED ||
          DT == DT_DELETED ||
          DT == DT_DIFF_FILES ) { found_diff = true
      } else                    { *p_dl--
      }
    }
  }
  return found_diff;
}

// Look for difference in white space at beginning or ending of lines:
func (m *Diff) Do_N_Search_for_Diff_WhiteSpace( p_dl *int,
                                                p_DI_List *Vector[Diff_Info] ) bool {
  found_diff := false

  NUM_LINES := m.NumLines()

  p_DI_List_o := &m.DI_List_S
  pF_o := m.pfS
  pF_m := m.pfL

  if( p_DI_List == &m.DI_List_S ) {
    p_DI_List_o = &m.DI_List_L
    pF_o = m.pfL
    pF_m = m.pfS
  }
  // If the current line has a difference in white space at beginning or end, start
  // searching on next line so the current line number is not automatically returned.
  var curr_line_has_LT_WS_diff bool =
    m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, *p_dl, p_DI_List, p_DI_List_o, pF_m, pF_o )

  dl_st := *p_dl
  if( curr_line_has_LT_WS_diff ) {
    dl_st = NUM_LINES-1
    if( 0 < *p_dl ) { dl_st = (*p_dl - 1) % NUM_LINES}
  }
  // Search from dl_st to end for lines of different length:
  for k:=dl_st; !found_diff && 0<=k; k-- {
    found_diff = m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, k, p_DI_List, p_DI_List_o, pF_m, pF_o )
  }
  if( !found_diff ) {
    // Search from top to dl_st for lines of different length:
    for k:=NUM_LINES-1; !found_diff && dl_st<k; k-- {
      found_diff = m.Line_Has_Leading_or_Trailing_WS_Diff( p_dl, k, p_DI_List, p_DI_List_o, pF_m, pF_o )
    }
  }
  return found_diff;
}

// Returns true if found next word, else false
//
func (m *Diff) GoToNextWord_GetPosition( ncp *CrsPos ) bool {

  pV := m_vis.CV()

  NUM_LINES := pV.p_fb.NumLines()
  if( 0==NUM_LINES ) { return false }

  found_space := false
  found_word  := false

  // Convert from diff line (CrsLine(m)), to view line:
  OCL := m.ViewLine( pV, m.CrsLine() ) //< Old cursor view line
  OCP := m.CrsChar()                   //< Old cursor position

  var isWord IsWord_Func = IsWord_Ident

  // Find white space, and then find non-white space
  for vl:=OCL; (!found_space || !found_word) && vl<NUM_LINES; vl++ {
    LL := pV.p_fb.LineLen( vl )
    if( LL == 0 || OCL < vl ) {
      found_space = true
      // Once we have encountered a space, word is anything non-space.
      // An empty line is considered to be a space.
      isWord = NotSpace
    }
    START_C := 0; if( OCL==vl ) { START_C = OCP }

    for p:=START_C; (!found_space || !found_word) && p<LL; p++ {
      R := pV.p_fb.GetR( vl, p )

      if( found_space  ) {
        if( isWord( R ) ) { found_word = true }
      } else {
        if( !isWord( R ) ) { found_space = true }
      }
      // Once we have encountered a space, word is anything non-space
      if( IsSpace( R ) ) { isWord = NotSpace }

      if( found_space && found_word ) {
        // Convert from view line back to diff line:
        dl := m.DiffLine( pV, vl )

        ncp.crsLine = dl
        ncp.crsChar = p
      }
    }
  }
  return found_space && found_word
}

// Return true if new cursor position found, else false
func (m *Diff) GoToPrevWord_GetPosition( ncp *CrsPos ) bool {

  pV := m_vis.CV()

  NUM_LINES := pV.p_fb.NumLines()
  if( 0==NUM_LINES ) { return false }

  // Convert from diff line (CrsLine(m)), to view line:
  OCL := m.ViewLine( pV, m.CrsLine() )
  LL  := pV.p_fb.LineLen( OCL )

  if( LL < m.CrsChar() ) { // Since cursor is now allowed past EOL,
                           // it may need to be moved back:
    if( 0 < LL && !IsSpace( pV.p_fb.GetR( OCL, LL-1 ) ) ) {
      // Backed up to non-white space, which is previous word, so return true
      // Convert from view line back to diff line:
      ncp.crsLine = m.CrsLine() //< diff line
      ncp.crsChar = LL-1
      return true
    } else {
      m.GoToCrsPos_NoWrite( m.CrsLine(), LLM1( LL ) )
    }
  }
  found_space := false
  found_word  := false
  OCP := m.CrsChar() // Old cursor position

  var isWord IsWord_Func = NotSpace

  // Find word to non-word transition
  for vl:=OCL; (!found_space || !found_word) && -1<vl; vl-- {
    LL := pV.p_fb.LineLen( vl )
    if( LL==0 || vl<OCL ) {
      // Once we have encountered a space, word is anything non-space.
      // An empty line is considered to be a space.
      isWord = NotSpace
    }
    START_C := LL-1; if( OCL==vl ) { START_C = OCP-1 }

    for p:=START_C; (!found_space || !found_word) && -1<p; p-- {
      R := pV.p_fb.GetR( vl, p)

      if( found_word  ) {
        if( !isWord( R ) || p==0 ) { found_space = true }
      } else {
        if( isWord( R ) ) {
          found_word = true
          if( 0==p ) { found_space = true }
        }
      }
      // Once we have encountered a space, word is anything non-space
      if( IsIdent( R ) ) { isWord = IsWord_Ident }

      if( found_space && found_word ) {
        // Convert from view line back to diff line:
        dl := m.DiffLine( pV, vl )

        ncp.crsLine = dl
        ncp.crsChar = p
      }
    }
    if( found_space && found_word ) {
      if( 0 < ncp.crsChar && ncp.crsChar < LL-1 ) { ncp.crsChar++ }
    }
  }
  return found_space && found_word;
}

// Returns true if found end of word, else false
// 1. If at end of word, or end of non-word, move to next char
// 2. If on white space, skip past white space
// 3. If on word, go to end of word
// 4. If on non-white-non-word, go to end of non-white-non-word
func (m *Diff) GoToEndOfWord_GetPosition( ncp *CrsPos ) bool {

  pV := m_vis.CV()
  pfb := pV.p_fb

  NUM_LINES := pfb.NumLines()
  if( 0==NUM_LINES ) { return false }

  // Convert from diff line (CrsLine(m)), to view line:
  CL := m.ViewLine( pV, m.CrsLine() )
  LL := pfb.LineLen( CL )
  CP := m.CrsChar() // Cursor position

  // At end of line, or line too short:
  if( (LL-1) <= CP || LL < 2 ) { return false }

  CR := pfb.GetR( CL, CP )   // Current rune
  NR := pfb.GetR( CL, CP+1 ) // Next rune

  // 1. If at end of word, or end of non-word, move to next char
  if( (IsWord_Ident   ( CR ) && !IsWord_Ident   ( NR )) ||
      (IsWord_NonIdent( CR ) && !IsWord_NonIdent( NR )) ) {
    CP++
  }
  // 2. If on white space, skip past white space
  if( IsSpace( pfb.GetR(CL, CP) ) ) {
    for ; CP<LL && IsSpace( pfb.GetR(CL, CP) ); CP++ {;}
    if( LL <= CP ) { return false } // Did not find non-white space
  }
  // At this point (CL,CP) should be non-white space
  CR = pfb.GetR( CL, CP )  // Current char

  ncp.crsLine = m.CrsLine() // Diff line

  if( IsWord_Ident( CR ) ) { // On identity
    // 3. If on word space, go to end of word space
    for ; CP<LL && IsWord_Ident( pfb.GetR(CL, CP) ); CP++ {
      ncp.crsChar = CP
    }
  } else if( IsWord_NonIdent( CR ) ) { // On Non-identity, non-white space
    // 4. If on non-white-non-word, go to end of non-white-non-word
    for ; CP<LL && IsWord_NonIdent( pfb.GetR(CL, CP) ); CP++ {
      ncp.crsChar = CP
    }
  } else { // Should never get here:
    return false
  }
  return true
}

func (m *Diff)  GoToOppositeBracket_Forward( ST_R, FN_R rune ) {

  pV := m_vis.CV()

  NUM_LINES := pV.p_fb.NumLines()

  // Convert from diff line (CrsLine(m)), to view line:
  CL := m.ViewLine( pV, m.CrsLine() )
  CC := m.CrsChar()

  // Search forward
  level := 0
  found := false

  for vl:=CL; !found && vl<NUM_LINES; vl++ {
    LL := pV.p_fb.LineLen( vl )

    p := 0; if( CL==vl ) { p = CC+1 }

    for ; !found && p<LL; p++ {
      R := pV.p_fb.GetR( vl, p )
      if       ( R==ST_R ) { level++;
      } else if( R==FN_R ) {
        if( 0 < level ) { level--
        } else {
          found = true
          // Convert from view line back to diff line:
          dl := m.DiffLine(pV, vl)

          m.GoToCrsPos_Write( dl, p )
        }
      }
    }
  }
}

func (m *Diff) GoToOppositeBracket_Backward( ST_R, FN_R rune ) {

  pV := m_vis.CV()

  // Convert from diff line (CrsLine(m)), to view line:
  CL := m.ViewLine( pV, m.CrsLine() )
  CC := m.CrsChar()

  // Search backward
  level := 0
  found := false

  for vl:=CL; !found && 0<=vl; vl-- {
    LL := pV.p_fb.LineLen( vl )

    p := LL-1; if( CL==vl ) { p = CC-1 }
    for ; !found && 0<=p; p-- {
      R := pV.p_fb.GetR( vl, p )

      if       ( R==ST_R ) { level++
      } else if( R==FN_R ) {
        if( 0 < level ) { level--
        } else {
          found = true
          // Convert from view line back to dif line:
          dl := m.DiffLine(pV, vl)

          m.GoToCrsPos_Write( dl, p )
        }
      }
    }
  }
}

func (m *Diff) GoToFile() {

  pV := m_vis.CV()

  cDI_List := m.View_2_DI_List_C( pV )
  oDI_List := m.View_2_DI_List_O( pV )

  var cDT Diff_Type = cDI_List.Get( m.CrsLine() ).diff_type // Current diff type
  var oDT Diff_Type = oDI_List.Get( m.CrsLine() ).diff_type // Other   diff type

  var fname_vec Vector[rune]
  if( m.GetFileName_UnderCursor( &fname_vec ) ) {
    fname := string( fname_vec.data )
    did_diff := false
    // Special case, look at two files in diff mode:
    cV := m.View_C( pV ) // Current view
    oV := m.View_O( pV ) // Other   view

    var cPath string = FindFullFileNameRel2( cV.p_fb.dir_name, fname ) // Current side file to diff full fname
    var oPath string = FindFullFileNameRel2( oV.p_fb.dir_name, fname ) // Other   side file to diff full fname

    c_file_idx := 0 // Current side index of file to diff
    o_file_idx := 0 // Other   side index of file to diff
    if( m.GetBufferIndex( cPath, &c_file_idx ) &&
        m.GetBufferIndex( oPath, &o_file_idx ) ) {

      var c_file_buf *FileBuf = m_vis.GetFileBuf( c_file_idx )
      var o_file_buf *FileBuf = m_vis.GetFileBuf( o_file_idx )
      // Files with same name and different contents
      // or directories with same name but different paths
      if( (cDT == DT_DIFF_FILES && oDT == DT_DIFF_FILES) ||
          (cV.p_fb.is_dir && oV.p_fb.is_dir &&
           (c_file_buf.file_name == o_file_buf.file_name) &&
           (c_file_buf.dir_name  != o_file_buf.dir_name) ) ) {
        // Save current view context for when we come back
        cV_vl_cl := m.ViewLine( cV, m.CrsLine() )
        cV_vl_tl := m.ViewLine( cV, m.topLine )
        cV.topLine  = cV_vl_tl
        cV.crsRow   = cV_vl_cl - cV_vl_tl
        cV.leftChar = m.leftChar
        cV.crsCol   = m.crsCol

        // Save other view context for when we come back
        oV_vl_cl := m.ViewLine( oV, m.CrsLine() )
        oV_vl_tl := m.ViewLine( oV, m.topLine )
        oV.topLine  = oV_vl_tl
        oV.crsRow   = oV_vl_cl - oV_vl_tl
        oV.leftChar = m.leftChar
        oV.crsCol   = m.crsCol

        did_diff = m_vis.Diff_By_File_Indexes( cV, c_file_idx, oV, o_file_idx )
      }
    }
    if( !did_diff ) {
      // Normal case, dropping out of diff mode to look at file:
      m_vis.GoToBuffer_Fname( fname )
    }
  }
}

func (m *Diff) GetFileName_UnderCursor( fname *Vector[rune] ) bool {

  pV  := m_vis.CV()
  pfb := pV.p_fb
  got_filename := false

  DL := m.CrsLine()            // Diff line number
  VL := m.ViewLine( pV, DL ) // View line number
  LL := pfb.LineLen( VL )

  if( 0 < LL ) {
    m.MoveInBounds_Line()
    CP := m.CrsChar()
    R := pfb.GetR( VL, CP )

    if( IsFileNameChar( R ) ) {
      // Get the file name:
      got_filename = true

      fname.Push( R )

      // Search backwards, until white space is found:
      for k:=CP-1; -1<k; k--  {
        R = pfb.GetR( VL, k )

        if( !IsFileNameChar( R ) ) { break
        } else { fname.Insert( 0, R )
        }
      }
      // Search forwards, until white space is found:
      for k:=CP+1; k<LL; k++ {
        R = pfb.GetR( VL, k )

        if( !IsFileNameChar( R ) ) { break
        } else { fname.Push( R )
        }
      }
    //EnvKeys2Vals( fname );
    }
  }
  return got_filename
}

func (m *Diff) GetBufferIndex( file_path string, file_index *int ) bool {

  got_buffer_index := false

  // 1. Search for file_path in buffer list
  if( m_vis.HaveFile( file_path, file_index ) ) {
    got_buffer_index = true

  // 2. See if file exists, and if so, add a file buffer
  } else if( FileExists( file_path ) ) {
    // pfb gets added to m_vis.files in Add_FileBuf_2_Lists_Create_Views()
    p_fb := new( FileBuf )
    p_fb.Init_FB( file_path, FT_UNKNOWN )

    if( m_vis.HaveFile( file_path, file_index ) ) {
      got_buffer_index = true
    }
  }
  return got_buffer_index
}

