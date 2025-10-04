
package main

import (
  "bytes"
  "fmt"
  "github.com/gdamore/tcell/v2"
  "regexp"
  "unicode"
)

type FileView struct {
  p_fb *FileBuf
  x,y int                // View position on screen
  nRows,nCols int        // Number of rows and columns in view, including border
  topLine, leftChar int  // View current top line and left character position
  crsRow, crsCol int     // View current cursor row and column in file

  inInsertMode, inReplaceMode bool
  inVisualMode, inVisualBlock bool
  inMap bool

  v_st_line, v_st_char int // Visual start line number, and char number on line
  v_fn_line, v_fn_char int // Visual ending line number, and char number on line

  tile_pos Tile_Pos
  in_diff_mode bool

  un_saved_change_sts bool
  external_change_sts bool

  p_diff *Diff
}

func (m *FileView) Init( file_buf *FileBuf ) {
  m.p_fb = file_buf
  m.nRows = m_console.Num_Rows()
  m.nCols = m_console.Num_Cols()
  m.tile_pos = TP_FULL
}

type IsWord_Func func(rune)bool

func (m *FileView) X() int { return m.x; }
func (m *FileView) Y() int { return m.y; }

func (m *FileView) WinRows() int { return m.nRows; }
func (m *FileView) WinCols() int { return m.nCols; }

func (m *FileView) CrsChar() int {
  return m.leftChar + m.crsCol
}

func (m *FileView) CrsLine() int {
  return m.topLine  + m.crsRow
}

// Converts from working view window row to global screen row
//
func (m *FileView) Row_Win_2_GL( win_row int ) int {
  return m.y + 1 + win_row
}

// Converts from working view window column to global screen column
//
func (m *FileView) Col_Win_2_GL( win_col int ) int {
  return m.x + 1 + win_col
}

// Translates zero based file line number to zero based global row
func (m *FileView) Line_2_GL( file_line int ) int {

  return m.y + 1 + file_line - m.topLine
}

// Translates zero based file line char position to zero based global column
//
func (m *FileView) Char_2_GL( line_char int ) int {

  return m.x + 1 + line_char - m.leftChar
}

// Returns working rows in view window
//
func (m *FileView) WorkingRows() int {
  return m.nRows - 5
}

// Returns working columns in view window
//
func (m *FileView) WorkingCols() int {
  return m.nCols - 2
}

func (m *FileView) BotLine() int {
  return m.topLine  + m.WorkingRows()-1
}

// (Index/Position) on current line of rune that is or would be displayed
// in right column
func (m *FileView) RightChar() int {
  return m.leftChar + m.WorkingCols()-1
}

func (m *FileView) Sts__Line_Row() int {
  return m.Row_Win_2_GL( m.WorkingRows() )
}

func (m *FileView) File_Line_Row() int {
  return m.Row_Win_2_GL( m.WorkingRows()+1 )
}

func (m *FileView) Cmd__Line_Row() int {
  return m.Row_Win_2_GL( m.WorkingRows()+2 )
}

func (m *FileView) Set_crsRowCol( row, col int ) {
  m.crsRow = row
  m.crsCol = col
}

func (m *FileView) GetTilePos() Tile_Pos {
  return m.tile_pos
}

func (m *FileView) SetTilePos( tp Tile_Pos ) {
  m.tile_pos = tp

  m.SetViewPos()
}

//func (m *FileView) PrintCursor() {
//
//  m_console.ShowCursor( m.Row_Win_2_GL( m.crsRow ), m.Col_Win_2_GL( m.crsCol ) )
//  m_console.Show()
//}

func (m *FileView) PrintCursor() {

  if( nil != m.p_diff ) {
    m.p_diff.PrintCursor( m )
  } else {
    m_console.ShowCursor( m.Row_Win_2_GL( m.crsRow ), m.Col_Win_2_GL( m.crsCol ) )
    m_console.Show()
  }
}

func (m *FileView) RepositionView() {
  // If a window re-size has taken place, and the window has gotten
  // smaller, change top line and left char if needed, so that the
  // cursor is in the FileView when it is re-drawn
  if( m.WorkingRows() <= m.crsRow ) {
    shift := m.crsRow - m.WorkingRows() + 1
    m.topLine += shift
    m.crsRow  -= shift
  }
  if( m.WorkingCols() <= m.crsCol ) {
    shift := m.crsCol - m.WorkingCols() + 1
    m.leftChar += shift
    m.crsCol   -= shift
  }
}

func (m *FileView) Update_not_PrintCursor() {

  if( !m_key.RunningDot() ) {

    if( nil != m.p_diff ) {
      m.p_diff.Update1V( m )
    } else {
      m.p_fb.Find_Styles( m.topLine + m.WorkingRows() )
      m.p_fb.Find_Regexs( m.topLine, m.WorkingRows() )

      m.RepositionView()
      m.PrintBorders()
      m.PrintWorkingView()
      m.PrintStsLine()
      m.PrintFileLine()
      m.PrintCmdLine()
    }
  }
}

func (m *FileView) Update_and_PrintCursor() {
  m.Update_not_PrintCursor()

  m.PrintCursor()
}

func (m *FileView) PrintWorkingView() {

  var NUM_LINES int = m.p_fb.NumLines()
  var WR        int = m.WorkingRows()
  var WC        int = m.WorkingCols()

  var row int = 0
  for k:=m.topLine; k<NUM_LINES && row<WR; k++ {
    // Dont allow line wrap:
    var LL    int = m.p_fb.LineLen( k )
    var G_ROW int = m.Row_Win_2_GL( row )
    var col int = 0
    for i:=m.leftChar; i<LL && col<WC; i++ {
      var p_TS *tcell.Style = m.Get_Style( k, i )
      var R rune = m.p_fb.GetR( k, i )

      var G_COL = m.Col_Win_2_GL( col )
      m.PrintWorkingView_Set( LL, G_ROW, G_COL, i, R, p_TS )
      col++
    }
    for ; col<WC; col++ {
      var G_COL = m.Col_Win_2_GL( col )
      m_console.SetR( G_ROW, G_COL, ' ', &TS_EMPTY )
    }
    row++
  }
  // Not enough lines to display, fill in with ~
  for ; row < WR; row++ {
    var G_ROW int = m.Row_Win_2_GL( row )

    m_console.SetR( G_ROW, m.Col_Win_2_GL( 0 ), '~', &TS_EOF )

    for col:=1; col<WC; col++ {
      m_console.SetR( G_ROW, m.Col_Win_2_GL( col ), ' ', &TS_EOF )
    }
  }
}

func (m *FileView) PrintWorkingView_Set( LL, G_ROW, G_COL, i int, R rune, p_TS *tcell.Style ) {

  if( '\r' == R && i==(LL-1) ) {
    // For readability, display carriage return at end of line as a space
    m_console.SetR( G_ROW, G_COL, ' ', &TS_NORMAL )
  } else {
    m_console.SetR( G_ROW, G_COL, R, p_TS )
  }
}

func (m *FileView) PrintBorders() {

  var HIGHLIGHT bool = ( 1 < m_vis.num_wins ) && ( m == m_vis.CV() )

  var p_S *tcell.Style = &TS_BORDER
  if( HIGHLIGHT ) {
    p_S = &TS_BORDER_HI
  }
  m.print_borders_top( p_S )
  m.print_borders_bottom( p_S )
  m.print_borders_left( p_S )
  m.print_borders_right( p_S )
}

func (m *FileView) PrintStsLine() {
  var CL int = m.CrsLine()
  var CC int = m.CrsChar()
  var LL int = m.p_fb.LineLen( CL )
  var WC int = m.WorkingCols()

  var fileSize int = m.p_fb.GetSize()
  var  crsByte int = m.p_fb.GetCursorByte( CL, CC )
  var percent int = int(100*float64(crsByte)/float64(fileSize) + 0.5)

  var buf bytes.Buffer

  fmt.Fprintf( &buf, "Pos=(%d,%d)  (%d%%, %d/%d)  Char=(",
                     CL+1, CC+1,
                     percent, crsByte, m.p_fb.GetSize() )
  if 0 < LL && CC < LL {
    var R rune = m.p_fb.GetR( CL, CC )
    fmt.Fprintf( &buf, "%d,%c", R, rune(R) )
  }
  fmt.Fprintf( &buf, ")" )

  for k:=buf.Len(); k<WC; k++ {
    fmt.Fprintf( &buf, " " )
  }
  m_console.SetBuffer( m.Sts__Line_Row(), m.Col_Win_2_GL( 0 ), &buf, &TS_BORDER )
}

func (m *FileView) PrintFileLine() {

  var WC int = m.WorkingCols()
  var buf bytes.Buffer

  fmt.Fprintf( &buf, "%s", m.p_fb.GetPath() )

  var s_b []byte = buf.Bytes()

  for k:=0; k<WC; k++ {
    if k < len( s_b ) {
      m_console.SetR( m.File_Line_Row(), m.Col_Win_2_GL( k ), rune(s_b[k]), &TS_BORDER )
    } else {
      m_console.SetR( m.File_Line_Row(), m.Col_Win_2_GL( k ), ' ', &TS_BORDER )
    }
  }
}

func (m *FileView) Set_Insert_Mode( on bool ) {

  if( on != m.inInsertMode ) {
    m.inInsertMode = on
    m.PrintCmdLine()
    m_console.Show()
  }
}

func (m *FileView) Set_Replace_Mode( on bool ) {

  if( on != m.inReplaceMode ) {
    m.inReplaceMode = on
    m.PrintCmdLine()
    m_console.Show()
  }
}

func (m *FileView) Set_Visual_Mode( on bool ) {

  if( on != m.inVisualMode ) {
    m.inVisualMode = on
    m.PrintCmdLine()
    m_console.Show()
  }
}

func (m *FileView) Set_VisualB_Mode( on bool ) {

  if( on != m.inVisualBlock ) {
    m.inVisualBlock = on
    m.PrintCmdLine()
    m_console.Show()
  }
}

func (m *FileView) PrintCmdLine() {
  // Command line row in window:
  var WIN_COL int = 0

  var G_ROW int = m.Cmd__Line_Row()
  var G_COL int = m.Col_Win_2_GL( WIN_COL )

  var col int = 0
  var buf bytes.Buffer

  if( m.inInsertMode ) {
    col,_ = fmt.Fprintf( &buf, "--INSERT--" )
    m_console.SetBuffer( G_ROW, G_COL, &buf, &TS_BANNER )

  } else if( m.inReplaceMode ) {
    col,_ = fmt.Fprintf( &buf, "--REPLACE--" )
    m_console.SetBuffer( G_ROW, G_COL, &buf, &TS_BANNER )

  } else if( m.inVisualMode ) {
    col,_ = fmt.Fprintf( &buf, "--VISUAL--" )
    m_console.SetBuffer( G_ROW, G_COL, &buf, &TS_BANNER )

  } else if( m.inVisualBlock ) {
    col,_ = fmt.Fprintf( &buf, "--VISUAL_B--" )
    m_console.SetBuffer( G_ROW, G_COL, &buf, &TS_BANNER )

  } else {
    col,_ = fmt.Fprintf( &buf, "            " )
    m_console.SetBuffer( G_ROW, G_COL, &buf, &TS_NORMAL )
  }

  if( m.inMap ) {
    mapping := []rune("--MAPPING--")
    mapping_len := len( mapping )
    map_col := m.nCols-2-mapping_len
    for k:=col; k<map_col; k++ {
      G_COL = m.Col_Win_2_GL( k )
      m_console.SetR( G_ROW, G_COL, ' ', &TS_NORMAL )
    }
    // --MAPPING-- banner is shown at the right side of the command line row:
    G_COL = m.Col_Win_2_GL( map_col )
    m_console.SetSR( G_ROW, G_COL, mapping, &TS_BANNER )
  } else {
    WC := m.WorkingCols()
    for k:=col; k<WC; k++ {
      G_COL = m.Col_Win_2_GL( k )
      m_console.SetR( G_ROW, G_COL, ' ', &TS_NORMAL )
    }
  }
}

func (m *FileView) Border_Rune_1( C_ex rune ) rune {

  border_rune := ' '

  if( m.p_fb.Changed() ) {
    border_rune = '+'
  } else if( m.p_fb.changed_externally ) {
    border_rune = C_ex
  }
  return border_rune
}

func (m *FileView) Border_Rune_2( C_ex rune ) rune {

  border_rune := ' '

  if( m.p_fb.changed_externally ) {
    border_rune = C_ex
  } else if( m.p_fb.Changed() ) {
    border_rune = '+'
  }
  return border_rune
}

func (m *FileView) print_borders_top( p_S *tcell.Style ) {

  BORDER_RUNE_1 := m.Border_Rune_1( DIR_DELIM )
  BORDER_RUNE_2 := m.Border_Rune_2( DIR_DELIM )

  var ROW_G int = m.y

  for k:=0; k<m.nCols; k++ {
    var COL_G int = m.x + k

    if( 0!=k%2 ) { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_2, p_S )
    } else       { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_1, p_S )
    }
  }
}

func (m *FileView) print_borders_bottom( p_S *tcell.Style ) {

  BORDER_RUNE_1 := m.Border_Rune_1( DIR_DELIM )
  BORDER_RUNE_2 := m.Border_Rune_2( DIR_DELIM )

  var ROW_G int = m.y + m.nRows - 1

  for k:=0; k<m.nCols; k++ {
    var COL_G int = m.x + k

    if( 0!=k%2 ) { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_2, p_S )
    } else       { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_1, p_S )
    }
  }
}

func (m *FileView) print_borders_left( p_S *tcell.Style ) {

  BORDER_RUNE_1 := m.Border_Rune_1( DIR_DELIM )
  BORDER_RUNE_2 := m.Border_Rune_2( DIR_DELIM )

  var COL_G int = m.x

  for k:=0; k<m.nRows; k++ {
    var ROW_G int = m.y + k

    if( 0!=k%2 ) { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_2, p_S )
    } else       { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_1, p_S )
    }
  }
}

func (m *FileView) print_borders_right( p_S *tcell.Style ) {

  BORDER_RUNE_1 := m.Border_Rune_1( DIR_DELIM )
  BORDER_RUNE_2 := m.Border_Rune_2( DIR_DELIM )

  var COL_G int = m.x + m.nCols - 1

  for k:=0; k<m.nRows; k++ {
    var ROW_G int = m.y + k
    // Do not print bottom right hand corner of console, because
    // on some terminals it scrolls the whole console screen up one line:
    if( 0!=k%2 ) { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_2, p_S )
    } else       { m_console.SetR( ROW_G, COL_G, BORDER_RUNE_1, p_S )
    }
  }
}

func (m *FileView) Get_Style( line, pos int ) *tcell.Style {

  var p_TS *tcell.Style = &TS_NORMAL

  if( m.InVisualArea( line, pos ) ) {
    p_TS = &TS_RV_NORMAL

    if       ( m.InStar    ( line, pos ) ) { p_TS = &TS_RV_STAR
    } else if( m.InStarInF ( line, pos ) ) { p_TS = &TS_RV_STAR_IN_F
    } else if( m.InDefine  ( line, pos ) ) { p_TS = &TS_RV_DEFINE
    } else if( m.InComment ( line, pos ) ) { p_TS = &TS_RV_COMMENT
    } else if( m.InConst   ( line, pos ) ) { p_TS = &TS_RV_CONST
    } else if( m.InControl ( line, pos ) ) { p_TS = &TS_RV_CONTROL
    } else if( m.InVarType ( line, pos ) ) { p_TS = &TS_RV_VARTYPE
    } else if( m.InNonAscii( line, pos ) ) { p_TS = &TS_RV_NONASCII
    }
  } else if( m.InStar    ( line, pos ) ) { p_TS = &TS_STAR
  } else if( m.InStarInF ( line, pos ) ) { p_TS = &TS_STAR_IN_F
  } else if( m.InDefine  ( line, pos ) ) { p_TS = &TS_DEFINE
  } else if( m.InComment ( line, pos ) ) { p_TS = &TS_COMMENT
  } else if( m.InConst   ( line, pos ) ) { p_TS = &TS_CONST
  } else if( m.InControl ( line, pos ) ) { p_TS = &TS_CONTROL
  } else if( m.InVarType ( line, pos ) ) { p_TS = &TS_VARTYPE
  } else if( m.InNonAscii( line, pos ) ) { p_TS = &TS_NONASCII
  }
  return p_TS
}

func (m *FileView) InVisualArea( line, pos int ) bool {

  if       ( m.inVisualMode )  { return m.InVisualStFn ( line, pos )
  } else if( m.inVisualBlock ) { return m.InVisualBlock( line, pos )
  }
  return false
}

func (m *FileView) InVisualStFn( line, pos int ) bool {

  if( m.v_st_line == line && line == m.v_fn_line ) {
    return (m.v_st_char <= pos && pos <= m.v_fn_char) ||
           (m.v_fn_char <= pos && pos <= m.v_st_char)
  } else if( (m.v_st_line < line && line < m.v_fn_line) ||
           (m.v_fn_line < line && line < m.v_st_line) ) {
    return true
  } else if( m.v_st_line == line && line < m.v_fn_line ) {
    return m.v_st_char <= pos
  } else if( m.v_fn_line == line && line < m.v_st_line ) {
    return m.v_fn_char <= pos
  } else if( m.v_st_line < line && line == m.v_fn_line ) {
    return pos <= m.v_fn_char
  } else if( m.v_fn_line < line && line == m.v_st_line ) {
    return pos <= m.v_st_char
  }
  return false
}

func (m *FileView) InVisualBlock( line, pos int ) bool {

  return ( m.v_st_line <= line && line <= m.v_fn_line &&
           m.v_st_char <= pos  && pos  <= m.v_fn_char ) || // bot rite
         ( m.v_st_line <= line && line <= m.v_fn_line &&
           m.v_fn_char <= pos  && pos  <= m.v_st_char ) || // bot left
         ( m.v_fn_line <= line && line <= m.v_st_line &&
           m.v_st_char <= pos  && pos  <= m.v_fn_char ) || // top rite
         ( m.v_fn_line <= line && line <= m.v_st_line &&
           m.v_fn_char <= pos  && pos  <= m.v_st_char )    // top left
}

func (m *FileView) InStar( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_STAR )
}

func (m *FileView) InStarInF( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_STAR_IN_F )
}

func (m *FileView) InStarOrStarInF( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_STAR | HI_STAR_IN_F )
}

func (m *FileView) InDefine( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_DEFINE )
}

func (m *FileView) InComment( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_COMMENT )
}

func (m *FileView) InConst( line, pos int ) bool {
  return m.p_fb.HasStyle( line, pos, HI_CONST )
}

func (m *FileView) InControl( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_CONTROL )
}

func (m *FileView) InVarType( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_VARTYPE )
}

func (m *FileView) InNonAscii( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_NONASCII )
}

func (m *FileView) Set_Cmd_Line_Msg( msg string ) {
  // FIXME
}

func (m *FileView) Has_Context() bool {

  return 0 != m.topLine  ||
         0 != m.leftChar ||
         0 != m.crsRow   ||
         0 != m.crsCol
}

func (m *FileView) Clear_Context() {
  m.topLine  = 0
  m.leftChar = 0
  m.crsRow   = 0
  m.crsCol   = 0
}

func (m *FileView) Check_Context() {

  var NUM_LINES int = m.p_fb.NumLines()

  if( 0 == NUM_LINES ) {
    m.Clear_Context()

  } else {
    var changed bool = false
    var CL int = m.CrsLine()

    if( NUM_LINES <= CL ) {
      CL = NUM_LINES-1
      changed = true
    }
    var LL int = m.p_fb.LineLen( CL )
    var CP int = m.CrsChar()
    if( LL <= CP ) {
      CP = LLM1(LL)
      changed = true
    }
    if( changed ) {
      m.GoToCrsPos_NoWrite( CL, CP )
    }
  }
}

func (m *FileView) Set_Context_pFV( p_fv *FileView ) {

  m.topLine  = p_fv.topLine
  m.leftChar = p_fv.leftChar
  m.crsRow   = p_fv.crsRow
  m.crsCol   = p_fv.crsCol
}

func (m *FileView) Set_Context_4Is( topLine, leftChar, crsRow, crsCol int ) {

  m.topLine  = topLine
  m.leftChar = leftChar
  m.crsRow   = crsRow
  m.crsCol   = crsCol
}

func (m *FileView) SetViewPos() {

  m.TilePos_2_x()
  m.TilePos_2_y()
  m.TilePos_2_nRows()
  m.TilePos_2_nCols()
}

func (m *FileView) TilePos_2_x() {
  // TP_FULL     , TP_BOT__HALF    , TP_LEFT_QTR
  // TP_LEFT_HALF, TP_TOP__LEFT_QTR, TP_TOP__LEFT_8TH
  // TP_TOP__HALF, TP_BOT__LEFT_QTR, TP_BOT__LEFT_8TH
  // TP_LEFT_THIRD, TP_LEFT_TWO_THIRDS
  m.x = 0

  if( TP_RITE_HALF         == m.tile_pos ||
      TP_TOP__RITE_QTR     == m.tile_pos ||
      TP_BOT__RITE_QTR     == m.tile_pos ||
      TP_RITE_CTR__QTR     == m.tile_pos ||
      TP_TOP__RITE_CTR_8TH == m.tile_pos ||
      TP_BOT__RITE_CTR_8TH == m.tile_pos ) {

    m.x = m.Cols_Left_Half()

  } else if( TP_LEFT_CTR__QTR     == m.tile_pos ||
             TP_TOP__LEFT_CTR_8TH == m.tile_pos ||
             TP_BOT__LEFT_CTR_8TH == m.tile_pos ) {

    m.x = m.Cols_Left_Far_Qtr()

  } else if( TP_RITE_QTR      == m.tile_pos ||
             TP_TOP__RITE_8TH == m.tile_pos ||
             TP_BOT__RITE_8TH == m.tile_pos ) {

    m.x = m.Cols_Left_Half() + m.Cols_Rite_Ctr_Qtr()

  } else if( TP_CTR__THIRD      == m.tile_pos ||
             TP_RITE_TWO_THIRDS == m.tile_pos ) {

    m.x = m.Cols_Left_Third()

  } else if( TP_RITE_THIRD == m.tile_pos ) {

    m.x = m.Cols_Left_Third() + m.Cols_Ctr__Third()
  }
}

func (m *FileView) TilePos_2_y() {

  var CON_ROWS int = m_console.Num_Rows()

  // TP_FULL         , TP_LEFT_CTR__QTR
  // TP_LEFT_HALF    , TP_RITE_CTR__QTR
  // TP_RITE_HALF    , TP_RITE_QTR
  // TP_TOP__HALF    , TP_TOP__LEFT_8TH
  // TP_TOP__LEFT_QTR, TP_TOP__LEFT_CTR_8TH
  // TP_TOP__RITE_QTR, TP_TOP__RITE_CTR_8TH
  // TP_LEFT_QTR     , TP_TOP__RITE_8TH
  // TP_LEFT_THIRD   , TP_CTR__THIRD, TP_RITE_THIRD
  // TP_LEFT_TWO_THIRDS, TP_RITE_TWO_THIRDS
  m.y = 0

  if( TP_BOT__HALF         == m.tile_pos ||
      TP_BOT__LEFT_QTR     == m.tile_pos ||
      TP_BOT__RITE_QTR     == m.tile_pos ||
      TP_BOT__LEFT_8TH     == m.tile_pos ||
      TP_BOT__LEFT_CTR_8TH == m.tile_pos ||
      TP_BOT__RITE_CTR_8TH == m.tile_pos ||
      TP_BOT__RITE_8TH     == m.tile_pos ) {

    m.y = CON_ROWS/2
  }
}

func (m *FileView) TilePos_2_nRows() {

  var CON_ROWS int = m_console.Num_Rows()

  var ODD_ROWS bool = 0 != CON_ROWS%2

  // TP_TOP__HALF        , TP_BOT__HALF        ,
  // TP_TOP__LEFT_QTR    , TP_BOT__LEFT_QTR    ,
  // TP_TOP__RITE_QTR    , TP_BOT__RITE_QTR    ,
  // TP_TOP__LEFT_8TH    , TP_BOT__LEFT_8TH    ,
  // TP_TOP__LEFT_CTR_8TH, TP_BOT__LEFT_CTR_8TH,
  // TP_TOP__RITE_CTR_8TH, TP_BOT__RITE_CTR_8TH,
  // TP_TOP__RITE_8TH    , TP_BOT__RITE_8TH    ,
  m.nRows = CON_ROWS/2

  if( TP_FULL            == m.tile_pos ||
      TP_LEFT_HALF       == m.tile_pos ||
      TP_RITE_HALF       == m.tile_pos ||
      TP_LEFT_QTR        == m.tile_pos ||
      TP_LEFT_CTR__QTR   == m.tile_pos ||
      TP_RITE_CTR__QTR   == m.tile_pos ||
      TP_RITE_QTR        == m.tile_pos ||
      TP_LEFT_THIRD      == m.tile_pos ||
      TP_CTR__THIRD      == m.tile_pos ||
      TP_RITE_THIRD      == m.tile_pos ||
      TP_LEFT_TWO_THIRDS == m.tile_pos ||
      TP_RITE_TWO_THIRDS == m.tile_pos ) {

    m.nRows = CON_ROWS
  }
  if( ODD_ROWS && ( TP_BOT__HALF         == m.tile_pos ||
                    TP_BOT__LEFT_QTR     == m.tile_pos ||
                    TP_BOT__RITE_QTR     == m.tile_pos ||
                    TP_BOT__LEFT_8TH     == m.tile_pos ||
                    TP_BOT__LEFT_CTR_8TH == m.tile_pos ||
                    TP_BOT__RITE_CTR_8TH == m.tile_pos ||
                    TP_BOT__RITE_8TH     == m.tile_pos ) ) {
    m.nRows++
  }
}

func (m *FileView) TilePos_2_nCols() {

  if( TP_FULL      == m.tile_pos ||
      TP_TOP__HALF == m.tile_pos ||
      TP_BOT__HALF == m.tile_pos ) {

    m.nCols = m_console.Num_Cols()

  } else if( TP_LEFT_HALF     == m.tile_pos ||
             TP_TOP__LEFT_QTR == m.tile_pos ||
             TP_BOT__LEFT_QTR == m.tile_pos ) {

    m.nCols = m.Cols_Left_Half()

  } else if( TP_RITE_HALF     == m.tile_pos ||
             TP_TOP__RITE_QTR == m.tile_pos ||
             TP_BOT__RITE_QTR == m.tile_pos ) {

    m.nCols = m.Cols_Rite_Half()

  } else if( TP_LEFT_QTR      == m.tile_pos ||
             TP_TOP__LEFT_8TH == m.tile_pos ||
             TP_BOT__LEFT_8TH == m.tile_pos ) {

    m.nCols = m.Cols_Left_Far_Qtr()

  } else if( TP_LEFT_CTR__QTR     == m.tile_pos ||
             TP_TOP__LEFT_CTR_8TH == m.tile_pos ||
             TP_BOT__LEFT_CTR_8TH == m.tile_pos ) {

    m.nCols = m.Cols_Left_Ctr_Qtr()

  } else if( TP_RITE_CTR__QTR     == m.tile_pos ||
             TP_TOP__RITE_CTR_8TH == m.tile_pos ||
             TP_BOT__RITE_CTR_8TH == m.tile_pos ) {

    m.nCols = m.Cols_Rite_Ctr_Qtr()

  } else if( TP_RITE_QTR      == m.tile_pos ||
             TP_TOP__RITE_8TH == m.tile_pos ||
             TP_BOT__RITE_8TH == m.tile_pos ) {

    m.nCols = m.Cols_Rite_Far_Qtr()

  } else if( TP_LEFT_THIRD == m.tile_pos ) {

    m.nCols = m.Cols_Left_Third()

  } else if( TP_CTR__THIRD == m.tile_pos ) {

    m.nCols = m.Cols_Ctr__Third()

  } else if( TP_RITE_THIRD == m.tile_pos ) {

    m.nCols = m.Cols_Rite_Third()

  } else if( TP_LEFT_TWO_THIRDS == m.tile_pos ) {

    m.nCols = m.Cols_Left_Third() + m.Cols_Ctr__Third()

  } else if( TP_RITE_TWO_THIRDS == m.tile_pos ) {

    m.nCols = m.Cols_Ctr__Third() + m.Cols_Rite_Third()
  }
}

func (m *FileView) Cols_Left_Half() int {

  var CON_COLS int = m_console.Num_Cols()

  if( 0 != CON_COLS%2 ) {
    return CON_COLS/2+1 //< Left side gets extra column
  }
  return CON_COLS/2;  //< Both sides get equal
}

func (m *FileView) Cols_Rite_Half() int {

  var CON_COLS int = m_console.Num_Cols()

  return CON_COLS - m.Cols_Left_Half()
}

func (m *FileView) Cols_Left_Far_Qtr() int {

  var COLS_LEFT_HALF int = m.Cols_Left_Half()

  if( 0 != COLS_LEFT_HALF%2 ) {
    return COLS_LEFT_HALF/2+1 //< Left ctr qtr gets extra column
  }
  return COLS_LEFT_HALF/2; //< Both qtrs get equal
}

func (m *FileView) Cols_Left_Ctr_Qtr() int {

  return m.Cols_Left_Half() - m.Cols_Left_Far_Qtr()
}

func (m *FileView) Cols_Rite_Far_Qtr() int {

  var COLS_RITE_HALF int = m.Cols_Rite_Half()

  if( 0 != COLS_RITE_HALF%2 ) {
    return COLS_RITE_HALF/2+1 //< Rite ctr qtr gets extra column
  }
  return COLS_RITE_HALF/2 //< Both sides get equal
}

func (m *FileView) Cols_Rite_Ctr_Qtr() int {

  return m.Cols_Rite_Half() - m.Cols_Rite_Far_Qtr()
}

func (m *FileView) Cols_Left_Third() int {

  var CON_COLS int = m_console.Num_Cols()

  if( CON_COLS%3==2 ) {
    return CON_COLS/3+1 // Ctr and left get extra column
  }
  return CON_COLS/3
}

func (m *FileView) Cols_Ctr__Third() int {

  var CON_COLS int = m_console.Num_Cols()

  if( CON_COLS%3==1 ) {
    return CON_COLS/3+1 // Ctr third gets extra column
  }
  return CON_COLS/3
}

func (m *FileView) Cols_Rite_Third() int {

  var CON_COLS int = m_console.Num_Cols()

  return CON_COLS - m.Cols_Left_Third() - m.Cols_Ctr__Third()
}

func (m *FileView) GoToLine( user_line_num int ) {

  var NL int = m.p_fb.NumLines()

  if( 0 == NL ) { m.PrintCursor()
  } else {
    // Internal line number is 1 less than user line number:
    var NCL int = user_line_num - 1; // New cursor line number

    if( NCL  < 0   ) { NCL = 0 }
    if( NL-1 < NCL ) { NCL = NL-1 }

    m.GoToCrsPos_Write( NCL, 0 )
  }
}

// GoDown Internal
//
func (m *FileView) GoDown_i( num int ) {
  var NUM_LINES int = m.p_fb.NumLines()
  var OCL       int = m.CrsLine()

  if 0<NUM_LINES && OCL < NUM_LINES-1 {
    var NCL int = OCL+num; // New cursor line

    if( NUM_LINES-1 < NCL ) { NCL = NUM_LINES-1; }

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func (m *FileView) GoDown( num int ) {

  if( nil != m.p_diff ) {
    m.p_diff.GoDown( num )
  } else {
    m.GoDown_i( num )
  }
}

// GoUp Internal
//
func (m *FileView) GoUp_i( num int ) {
  var NUM_LINES int = m.p_fb.NumLines()
  var OCL       int = m.CrsLine()

  if 0<NUM_LINES && 0 < OCL {
    var NCL int = OCL-num; // New cursor line

    if( NCL < 0 ) { NCL = 0; }

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func (m *FileView) GoUp( num int ) {

  if( nil != m.p_diff ) {
    m.p_diff.GoUp( num )
  } else {
    m.GoUp_i( num )
  }
}

func (m *FileView) GoRight_i( num int ) {
  var NUM_LINES int = m.p_fb.NumLines()

  if 0<NUM_LINES {
    var CL  int = m.CrsLine() // Cursor line
    var LL  int = m.p_fb.LineLen( CL )
    var OCP int = m.CrsChar() // Old cursor position

    if( 0<LL && OCP < LL-1 ) {
      var num int = 1
      var NCP int = OCP+num // New cursor position

      if( LL-1 < NCP ) { NCP = LL-1 }

      m.GoToCrsPos_Write( CL, NCP )
    }
  }
}

func (m *FileView) GoRight( num int ) {

  if( nil != m.p_diff ) {
    m.p_diff.GoRight( num )
  } else {
    m.GoRight_i( num )
  }
}

func (m *FileView) GoLeft_i( num int ) {
  var NUM_LINES int = m.p_fb.NumLines()

  var OCP int = m.CrsChar(); // Old cursor position

  if( 0<NUM_LINES && 0 < OCP ) {
    var num int = 1
    var NCP int = OCP-num; // New cursor position

    if( NCP < 0 ) { NCP = 0; }

    m.GoToCrsPos_Write( m.CrsLine(), NCP )
  }
}

func (m *FileView) GoLeft( num int ) {

  if( nil != m.p_diff ) {
    m.p_diff.GoLeft( num )
  } else {
    m.GoLeft_i( num )
  }
}

func (m *FileView) GoToCrsPos_NoWrite( ncp_crsLine, ncp_crsChar int ) {

  var NCL int = ncp_crsLine
  var NCP int = ncp_crsChar

  // These moves refer to View of buffer:
  var MOVE_DOWN  bool = m.BotLine()   < NCL
  var MOVE_RIGHT bool = m.RightChar() < NCP
  var MOVE_UP    bool = NCL < m.topLine
  var MOVE_LEFT  bool = NCP < m.leftChar

  if       ( MOVE_DOWN ) { m.topLine = NCL - m.WorkingRows() + 1
  } else if( MOVE_UP )   { m.topLine = NCL
  }
  m.crsRow = ncp_crsLine - m.topLine

  if       ( MOVE_RIGHT ) { m.leftChar = NCP - m.WorkingCols() + 1
  } else if( MOVE_LEFT  ) { m.leftChar = NCP
  }
  m.crsCol = ncp_crsChar - m.leftChar
}

func (m *FileView) GoToCrsPos_Write( ncp_crsLine, ncp_crsChar int ) {

  var OCL int = m.CrsLine()
  var OCP int = m.CrsChar()
  var NCL int = ncp_crsLine
  var NCP int = ncp_crsChar

  if( OCL == NCL && OCP == NCP ) {
    m.PrintCursor();  // Put cursor back into position.
    return
  }
  if( m.inVisualMode || m.inVisualBlock ) {
    m.v_fn_line = NCL
    m.v_fn_char = NCP
  }
  // These moves refer to View of buffer:
  var MOVE_DOWN  bool = m.BotLine()   < NCL
  var MOVE_RIGHT bool = m.RightChar() < NCP
  var MOVE_UP    bool = NCL < m.topLine
  var MOVE_LEFT  bool = NCP < m.leftChar

  var redraw bool = MOVE_DOWN || MOVE_RIGHT || MOVE_UP || MOVE_LEFT

  if( redraw ) {
    if       ( MOVE_DOWN ) { m.topLine = NCL - m.WorkingRows() + 1
    } else if( MOVE_UP )   { m.topLine = NCL
    }
    if       ( MOVE_RIGHT ) { m.leftChar = NCP - m.WorkingCols() + 1
    } else if( MOVE_LEFT  ) { m.leftChar = NCP
    }
    // m.crsRow and m.crsCol must be set to new values before calling CalcNewCrsByte
    m.Set_crsRowCol( NCL - m.topLine, NCP - m.leftChar )

    m.Update_and_PrintCursor()
  } else {
    if       ( m.inVisualMode  ) { m.GoToCrsPos_Write_Visual     ( OCL, OCP, NCL, NCP )
    } else if( m.inVisualBlock ) { m.GoToCrsPos_Write_VisualBlock( OCL, OCP, NCL, NCP )
    } else {
      // m.crsRow and m.crsCol must be set to new values before calling CalcNewCrsByte and PrintCursor
      m.Set_crsRowCol( NCL - m.topLine, NCP - m.leftChar )

      m.PrintStsLine()
      m.PrintCursor()  // Put cursor into position.
    }
  }
}

func (m *FileView) GoToCrsPos_Write_Visual( OCL, OCP, NCL, NCP int ) {
  // (old cursor pos) < (new cursor pos)
  OCP_LT_NCP := OCL < NCL || (OCL == NCL && OCP < NCP)

  if( OCP_LT_NCP ) { // Cursor moved forward
    m.GoToCrsPos_WV_Forward( OCL, OCP, NCL, NCP )
  } else { // NCP_LT_OCP // Cursor moved backward
    m.GoToCrsPos_WV_Backward( OCL, OCP, NCL, NCP )
  }
  m.crsRow = NCL - m.topLine
  m.crsCol = NCP - m.leftChar
  m.PrintCursor()
}

func (m *FileView) GoToCrsPos_Write_VisualBlock( OCL, OCP, NCL, NCP int ) {
  // m.v_fn_line == NCL && v_fn_char == NCP, so dont need to include
  // m.v_fn_line       and v_fn_char in Min and Max calls below:
  vis_box_left := Min_i( m.v_st_char, Min_i( OCP, NCP ) )
  vis_box_rite := Max_i( m.v_st_char, Max_i( OCP, NCP ) )
  vis_box_top  := Min_i( m.v_st_line, Min_i( OCL, NCL ) )
  vis_box_bot  := Max_i( m.v_st_line, Max_i( OCL, NCL ) )

  draw_box_left := Max_i( m.leftChar   , vis_box_left )
  draw_box_rite := Min_i( m.RightChar(), vis_box_rite )
  draw_box_top  := Max_i( m.topLine    , vis_box_top  )
  draw_box_bot  := Min_i( m.BotLine()  , vis_box_bot  )

  for l:=draw_box_top; l<=draw_box_bot; l++ {
    LL := m.p_fb.LineLen( l )

    for k:=draw_box_left; k<LL && k<=draw_box_rite; k++ {
      // On some terminals, the cursor on reverse video on white space does not
      // show up, so to prevent that, do not reverse video the cursor position:
      R  := m.p_fb.GetR( l, k )
      style := m.Get_Style( l, k )

      if( NCL==l && NCP==k ) {
        if( RV_Style( style ) ) {
          NonRV_style := RV_Style_2_NonRV( style )

          m_console.SetR( m.Line_2_GL( l ), m.Char_2_GL( k ), R, NonRV_style )
        }
      } else {
        m_console.SetR( m.Line_2_GL( l ), m.Char_2_GL( k ), R, style )
      }
    }
  }
  m.crsRow = NCL - m.topLine
  m.crsCol = NCP - m.leftChar
//Console::Update()
  m.PrintCursor()
//m.sts_line_needs_update = true
}

func (m *FileView) GoToTopLineInView() {

  m.GoToCrsPos_Write( m.topLine, m.CrsChar() )
}

func (m *FileView) GoToBotLineInView() {

  var NUM_LINES int = m.p_fb.NumLines()

  var bottom_line_in_view int = m.topLine + m.WorkingRows()-1

  bottom_line_in_view = Min_i( NUM_LINES-1, bottom_line_in_view )

  m.GoToCrsPos_Write( bottom_line_in_view, m.CrsChar() )
}

func (m *FileView) GoToMidLineInView() {

  var NUM_LINES int = m.p_fb.NumLines()

  // Default: Last line in file is not in view
  var NCL int = m.topLine + m.WorkingRows()/2; // New cursor line

  if( NUM_LINES-1 < m.BotLine() ) {
    // Last line in file above bottom of view
    NCL = m.topLine + (NUM_LINES-1 - m.topLine)/2
  }
  m.GoToCrsPos_Write( NCL, 0 )
}

func (m *FileView) GoToEndOfLine() {

  if( 0<m.p_fb.NumLines() ) {

    var LL  int = m.p_fb.LineLen( m.CrsLine() )
    var OCL int = m.CrsLine(); // Old cursor line

    if( m.inVisualBlock ) {
      // In Visual Block, $ puts cursor at the position
      // of the end of the longest line in the block
      var max_LL int = LL

      for L:=m.v_st_line; L<=m.v_fn_line; L++ {
        max_LL = Max_i( max_LL, m.p_fb.LineLen( L ) )
      }
      m.GoToCrsPos_Write( OCL, LLM1( max_LL ) )
    } else {
      m.GoToCrsPos_Write( OCL, LLM1( LL ) )
    }
  }
}

func (m *FileView) GoToBegOfLine() {

  if( 0<m.p_fb.NumLines() ) {
    var OCL int = m.CrsLine(); // Old cursor line

    m.GoToCrsPos_Write( OCL, 0 )
  }
}

func (m *FileView) GoToEndOfNextLine() {

  var NUM_LINES int = m.p_fb.NumLines()

  if( 0<NUM_LINES ) {
    var OCL int = m.CrsLine(); // Old cursor line

    if( OCL < (NUM_LINES-1) ) {
      // Before last line, so can go down
      var LL int = m.p_fb.LineLen( OCL+1 )

      m.GoToCrsPos_Write( OCL+1, LLM1( LL ) )
    }
  }
}

func (m *FileView) GoToEndOfFile() {

  var NUM_LINES int = m.p_fb.NumLines()

  if( 0<NUM_LINES ) {
    m.GoToCrsPos_Write( NUM_LINES-1, 0 )
  }
}

func (m *FileView) GoToTopOfFile() {

  m.GoToCrsPos_Write( 0, 0 )
}

func (m *FileView) GoToStartOfRow() {

  if( 0<m.p_fb.NumLines() ) {
    var OCL int = m.CrsLine(); // Old cursor line

    m.GoToCrsPos_Write( OCL, m.leftChar )
  }
}

func (m *FileView) GoToEndOfRow() {

  if( 0 < m.p_fb.NumLines() ) {
    var OCL int = m.CrsLine(); // Old cursor line

    var LL int = m.p_fb.LineLen( OCL )
    if( 0 < LL ) {
      var NCP int = Min_i( LL-1, m.leftChar + m.WorkingCols() - 1 )

      m.GoToCrsPos_Write( OCL, NCP )
    }
  }
}

func (m *FileView) GoToFile() {

  var fname string
  var ok bool

  if( m_vis.Is_BE_FILE( m.p_fb ) ) {

    fname, ok = m.GetFileName_WholeLine()
  } else {
    fname, ok = m.GetFileName_PartialLine()
  }

  if( ok ) { m_vis.GoToBuffer_Fname( fname ) }
}

// Returns true if found next word, else false
//
func (m *FileView) GoToNextWord_GetPosition( ncp *CrsPos ) bool {

  var NUM_LINES int = m.p_fb.NumLines()
  if( 0==NUM_LINES ) { return false; }

  var found_space bool = false
  var found_word  bool = false
  var OCL int = m.CrsLine(); // Old cursor line
  var OCP int = m.CrsChar(); // Old cursor position

  var isWord IsWord_Func = IsWord_Ident

  // Find white space, and then find non-white space
  for l:=OCL; (!found_space || !found_word) && l<NUM_LINES; l++ {

    var LL int = m.p_fb.LineLen( l )
    if( LL==0 || OCL<l ) {
      found_space = true
      // Once we have encountered a space, word is anything non-space.
      // An empty line is considered to be a space.
      isWord = NotSpace
    }
    var START_C int = True_1_else_2( OCL==l, OCP, 0 )

    for p:=START_C; (!found_space || !found_word) && p<LL; p++ {

      ncp.crsLine = l
      ncp.crsChar = p

      var R rune = m.p_fb.GetR( l, p )

      if( found_space  ) {
        if( isWord( R ) ) { found_word = true; }
      } else {
        if( !isWord( R ) ) { found_space = true; }
      }
      // Once we have encountered a space, word is anything non-space
      if( IsSpace( R ) ) { isWord = NotSpace; }
    }
  }
  return found_space && found_word
}

// Return true if new cursor position found, else false
//
func (m *FileView) GoToPrevWord_GetPosition( ncp *CrsPos ) bool {

  var NUM_LINES int = m.p_fb.NumLines()
  if( 0==NUM_LINES ) { return false; }

  var OCL int = m.CrsLine(); // Old cursor line
  var LL int = m.p_fb.LineLen( OCL )

  if( LL < m.CrsChar() ) { // Since cursor is now allowed past EOL,
                           // it may need to be moved back:
    if( 0<LL && !IsSpace( m.p_fb.GetR( OCL, LL-1 ) ) ) {
      // Backed up to non-white space, which is previous word, so return true
      ncp.crsLine = OCL
      ncp.crsChar = LL-1
      return true
    } else {
      m.GoToCrsPos_NoWrite( OCL, LLM1( LL ) )
    }
  }
  var found_space bool = false
  var found_word  bool = false
  var OCP int = m.CrsChar(); // Old cursor position

  var isWord IsWord_Func = NotSpace

  // Find word to non-word transition
  for l:=OCL; (!found_space || !found_word) && -1<l; l-- {

    var LL int = m.p_fb.LineLen( l )
    if( LL==0 || l<OCL ) {
      // Once we have encountered a space, word is anything non-space.
      // An empty line is considered to be a space.
      isWord = NotSpace
    }
    var START_C int = True_1_else_2( OCL==l, OCP-1, LL-1 )

    for p:=START_C; (!found_space || !found_word) && -1<p; p-- {

      ncp.crsLine = l
      ncp.crsChar = p

      var R rune = m.p_fb.GetR( l, p )

      if( found_word  ) {
        if( !isWord( R ) || p==0 ) { found_space = true; }
      } else {
        if( isWord( R ) ) {
          found_word = true
          if( 0==p ) { found_space = true; }
        }
      }
      // Once we have encountered a space, word is anything non-space
      if( IsWord_Ident( R ) ) { isWord = IsWord_Ident; }
    }
    if( found_space && found_word ) {
      if( 0 < ncp.crsChar && ncp.crsChar < LL-1 ) { ncp.crsChar++; }
    }
  }
  return found_space && found_word
}

// Returns true if found end of word, else false
// 1. If at end of word, or end of non-word, move to next char
// 2. If on white space, skip past white space
// 3. If on word, go to end of word
// 4. If on non-white-non-word, go to end of non-white-non-word
func (m *FileView) GoToEndOfWord_GetPosition( ncp *CrsPos ) bool {

  var NUM_LINES int = m.p_fb.NumLines()
  if( 0==NUM_LINES ) { return false; }

  var CL int = m.CrsLine(); // Cursor line
  var LL int = m.p_fb.LineLen( CL )
  var CP int = m.CrsChar(); // Cursor position

  // At end of line, or line too short:
  if( (LL-1) <= CP || LL < 2 ) { return false; }

  var CR rune = m.p_fb.GetR( CL, CP );   // Current byte
  var NR rune = m.p_fb.GetR( CL, CP+1 ); // Next byte

  // 1. If at end of word, or end of non-word, move to next byte
  if( (IsWord_Ident   ( CR ) && !IsWord_Ident   ( NR )) ||
      (IsWord_NonIdent( CR ) && !IsWord_NonIdent( NR )) ) {
    CP++
  }
  // 2. If on white space, skip past white space
  if( IsSpace( m.p_fb.GetR(CL, CP) ) ) {
    for ; CP<LL && IsSpace( m.p_fb.GetR(CL, CP) ); CP++ { ; }
    if( LL <= CP ) { return false; } // Did not find non-white space
  }
  // At this point (CL,CP) should be non-white space
  CR = m.p_fb.GetR( CL, CP );  // Current char

  ncp.crsLine = CL

  if( IsWord_Ident( CR ) ) { // On identity
    // 3. If on word space, go to end of word space
    for ; CP<LL && IsWord_Ident( m.p_fb.GetR(CL, CP) ); CP++ {
      ncp.crsChar = CP
    }
  } else if( IsWord_NonIdent( CR ) ) { // On Non-identity, non-white space
    // 4. If on non-white-non-word, go to end of non-white-non-word
    for ; CP<LL && IsWord_NonIdent( m.p_fb.GetR(CL, CP) ); CP++ {
      ncp.crsChar = CP
    }
  } else { // Should never get here:
    return false
  }
  return true
}

func (m *FileView) GoToNextWord() {

  var ncp = CrsPos{ 0, 0 }

  if( m.GoToNextWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *FileView) GoToPrevWord() {

  var ncp = CrsPos{ 0, 0 }

  if( m.GoToPrevWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *FileView) GoToEndOfWord() {

  var ncp = CrsPos{ 0, 0 }

  if( m.GoToEndOfWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *FileView) GoToOppositeBracket() {

  m.MoveInBounds_Line()

  var NUM_LINES int = m.p_fb.NumLines()
  var CL int = m.CrsLine()
  var CC int = m.CrsChar()
  var LL int = m.p_fb.LineLen( CL )

  if( 0<NUM_LINES && 0<LL ) {

    var R rune = m.p_fb.GetR( CL, CC )

    if( R=='{' || R=='[' || R=='(' ) {
      var finish_rune rune = 0
      if       ( R=='{' ) { finish_rune = '}'
      } else if( R=='[' ) { finish_rune = ']'
      } else if( R=='(' ) { finish_rune = ')'
      }
      m.GoToOppositeBracket_Forward( R, finish_rune )

    } else if( R=='}' || R==']' || R==')' ) {
      var finish_rune rune = 0
      if       ( R=='}' ) { finish_rune = '{'
      } else if( R==']' ) { finish_rune = '['
      } else if( R==')' ) { finish_rune = '('
      }
      m.GoToOppositeBracket_Backward( R, finish_rune )
    }
  }
}

func (m *FileView) GoToOppositeBracket_Forward( ST_R, FN_R rune ) {

  var NUM_LINES int = m.p_fb.NumLines()
  var CL int = m.CrsLine()
  var CC int = m.CrsChar()

  // Search forward
  var level int = 0
  var found bool = false

  for l:=CL; !found && l<NUM_LINES; l++ {

    var LL int = m.p_fb.LineLen( l )

    for p:=True_1_else_2(CL==l,CC+1,0); !found && p<LL; p++ {

      var R rune = m.p_fb.GetR( l, p )

      if       ( R==ST_R ) { level++
      } else if( R==FN_R ) {
        if( 0 < level ) { level--
        } else {
          found = true

          m.GoToCrsPos_Write( l, p )
        }
      }
    }
  }
}

func (m *FileView) GoToOppositeBracket_Backward( ST_R, FN_R rune ) {

  var CL int = m.CrsLine()
  var CC int = m.CrsChar()

  // Search forward
  var level int = 0
  var found bool = false

  for l:=CL; !found && 0<=l; l-- {

    var LL int = m.p_fb.LineLen( l )

    for p:=True_1_else_2( CL==l, CC-1, LL-1); !found && 0<=p; p-- {

      var R rune = m.p_fb.GetR( l, p )

      if       ( R==ST_R ) { level++
      } else if( R==FN_R ) {
        if( 0 < level ) { level--
        } else {
          found = true

          m.GoToCrsPos_Write( l, p )
        }
      }
    }
  }
}

func (m *FileView) GoToLeftSquigglyBracket() {

  m.MoveInBounds_Line()

  var  start_rune rune = '}'
  var finish_rune rune = '{'
  m.GoToOppositeBracket_Backward( start_rune, finish_rune )
}

func (m *FileView) GoToRightSquigglyBracket() {

  m.MoveInBounds_Line()

  var  start_rune rune = '{'
  var finish_rune rune = '}'
  m.GoToOppositeBracket_Forward( start_rune, finish_rune )
}

// Cursor is moving forward
// Write out from (OCL,OCP) up to but not including (NCL,NCP)
func (m *FileView) GoToCrsPos_WV_Forward( OCL, OCP, NCL, NCP int ) {
  if( OCL == NCL ) { // Only one line:
    for k:=OCP; k<NCP; k++ {
      R := m.p_fb.GetR( OCL, k )
      m_console.SetR( m.Line_2_GL( OCL ), m.Char_2_GL( k ), R, m.Get_Style(OCL,k) )
    }
  } else { // Multiple lines
    // Write out first line:
    OCLL := m.p_fb.LineLen( OCL ) // Old cursor line length
    END_FIRST_LINE := Min_i( m.RightChar()+1, OCLL )
    for k:=OCP; k<END_FIRST_LINE; k++ {
      R := m.p_fb.GetR( OCL, k )
      m_console.SetR( m.Line_2_GL( OCL ), m.Char_2_GL( k ), R, m.Get_Style(OCL,k) )
    }
    // Write out intermediate lines:
    for l:=OCL+1; l<NCL; l++ {
      LL := m.p_fb.LineLen( l ) // Line length
      END_OF_LINE := Min_i( m.RightChar()+1, LL )
      for k:=m.leftChar; k<END_OF_LINE; k++ {
        R := m.p_fb.GetR( l, k )
        m_console.SetR( m.Line_2_GL( l ), m.Char_2_GL( k ), R, m.Get_Style(l,k) )
      }
    }
    // Write out last line:
    // Print from beginning of next line to new cursor position:
    NCLL := m.p_fb.LineLen( NCL ) // Line length
    END := Min_i( NCLL, NCP )
    for k:=m.leftChar; k<END; k++ {
      R := m.p_fb.GetR( NCL, k )
      m_console.SetR( m.Line_2_GL( NCL ), m.Char_2_GL( k ), R, m.Get_Style(NCL,k)  )
    }
  }
}

// Cursor is moving backwards
// Write out from (OCL,OCP) back to but not including (NCL,NCP)
func (m *FileView) GoToCrsPos_WV_Backward( OCL, OCP, NCL, NCP int ) {
  if( OCL == NCL ) { // Only one line:
    LL := m.p_fb.LineLen( OCL ) // Line length
    if( 0 < LL ) {
      START := Min_i( OCP, LL-1 )
      for k:=START; NCP<k; k-- {
        R := m.p_fb.GetR( OCL, k )
        m_console.SetR( m.Line_2_GL( OCL ) , m.Char_2_GL( k ), R, m.Get_Style(OCL,k) )
      }
    }
  } else { // Multiple lines
    // Write out first line:
    OCLL := m.p_fb.LineLen( OCL ) // Old cursor line length
    if( 0 < OCLL ) {
      for k:=Min_i(OCP,OCLL-1); m.leftChar<=k; k-- {
        R := m.p_fb.GetR( OCL, k )
        m_console.SetR( m.Line_2_GL( OCL ), m.Char_2_GL( k ), R, m.Get_Style(OCL,k) )
      }
    }
    // Write out intermediate lines:
    for l:=OCL-1; NCL<l; l-- {
      LL := m.p_fb.LineLen( l ) // Line length
      if( 0 < LL ) {
        END_OF_LINE := Min_i( m.RightChar(), LL-1 )
        for k:=END_OF_LINE; m.leftChar<=k; k-- {
          R := m.p_fb.GetR( l, k )
          m_console.SetR( m.Line_2_GL( l ), m.Char_2_GL( k ), R, m.Get_Style(l,k) )
        }
      }
    }
    // Write out last line:
    // Go down to beginning of last line:
    NCLL := m.p_fb.LineLen( NCL ) // New cursor line length
    if( 0 < NCLL ) {
      END_LAST_LINE := Min_i( m.RightChar(), NCLL-1 )

      // Print from beginning of next line to new cursor position:
      for k:=END_LAST_LINE; NCP<=k; k-- {
        R := m.p_fb.GetR( NCL, k )
        m_console.SetR( m.Line_2_GL( NCL ), m.Char_2_GL( k ), R, m.Get_Style(NCL,k) )
      }
    }
  }
}

func (m *FileView) PageDown() {

  var NUM_LINES int = m.p_fb.NumLines()

  if( 0<NUM_LINES ) {
    var newTopLine int = m.topLine + m.WorkingRows() - 1
    // Subtracting 1 above leaves one line in common between the 2 pages.

    if( newTopLine < NUM_LINES ) {
      m.crsCol = 0
      m.topLine = newTopLine

      // Dont let cursor go past the end of the file:
      if( NUM_LINES <= m.CrsLine() ) {
        // This line places the cursor at the top of the screen, which I prefer:
        m.crsRow = 0
      }
      m.Update_and_PrintCursor()
    }
  }
}

func (m *FileView) PageDown_v() {

  NUM_LINES := m.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
    OCL := m.CrsLine()  // Old cursor line

    NCL := OCL + m.WorkingRows() - 1  // New cursor line

    // Dont let cursor go past the end of the file:
    if( NUM_LINES-1 < NCL ) { NCL = NUM_LINES-1 }

    m.GoToCrsPos_Write( NCL, 0 )
  }
}

func (m *FileView) PageUp() {

  // Dont scroll if we are at the top of the file:
  if( 0 < m.topLine ) {
    //Leave m.crsRow unchanged.
    m.crsCol = 0

    // Dont scroll past the top of the file:
    if( m.topLine < m.WorkingRows() - 1 ) {
      m.topLine = 0
    } else {
      m.topLine -= m.WorkingRows() - 1
    }
    m.Update_and_PrintCursor()
  }
}

func (m *FileView) PageUp_v() {

  NUM_LINES := m.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
    OCL := m.CrsLine()  // Old cursor line

    NCL := OCL - m.WorkingRows() + 1  // New cursor line

    // Check for underflow:
    if( NCL < 0 ) { NCL = 0 }

    m.GoToCrsPos_Write( NCL, 0 )
  }
}

func (m *FileView) MoveCurrLineToTop() {

  if( 0 < m.crsRow ) {
    m.topLine += m.crsRow
    m.crsRow = 0
    m.Update_and_PrintCursor()
  }
}

func (m *FileView) MoveCurrLineCenter() {

  var center int = int( 0.5*float32(m.WorkingRows()) + 0.5 )

  var OCL int = m.CrsLine(); // Old cursor line

  if( 0 < OCL && OCL < center && 0 < m.topLine ) {
    // Cursor line cannot be moved to center, but can be moved closer to center
    // CrsLine() does not change:
    m.crsRow += m.topLine
    m.topLine = 0
    m.Update_and_PrintCursor()

  } else if( center <= OCL &&
             center != m.crsRow ) {

    m.topLine += m.crsRow - center
    m.crsRow = center
    m.Update_and_PrintCursor()
  }
}

func (m *FileView) MoveCurrLineToBottom() {

  if( 0 < m.topLine ) {
    var WR  int = m.WorkingRows()
    var OCL int = m.CrsLine(); // Old cursor line

    if( WR-1 <= OCL ) {
      m.topLine -= WR - m.crsRow - 1
      m.crsRow = WR-1
      m.Update_and_PrintCursor()

    } else {
      // Cursor line cannot be moved to bottom, but can be moved closer to bottom
      // CrsLine() does not change:
      m.crsRow += m.topLine
      m.topLine = 0
      m.Update_and_PrintCursor()
    }
  }
}

// If past end of line, move back to end of line.
//
func (m *FileView) MoveInBounds_Line() {

  var CL  int = m.CrsLine()
  var LL  int = m.p_fb.LineLen( CL )
  var EOL int = LLM1( LL )

  if( EOL < m.CrsChar() ) {
    m.GoToCrsPos_NoWrite( CL, EOL )
  }
}

func (m *FileView) Do_i() {
  m.Set_Insert_Mode( true )

  if( 0 == m.p_fb.NumLines() ) { m.p_fb.PushLE(); }

  var LL int = m.p_fb.LineLen( m.CrsLine() );  // Line length

  // Since cursor is now allowed past EOL, it may need to be moved back:
  if LL < m.CrsChar() {
    // For user friendlyness, move cursor to new position immediately:
    m.GoToCrsPos_Write( m.CrsLine(), LL )
  }
  var count int
  for kr := m_key.In(); ! kr.IsESC(); kr = m_key.In() {
    if kr.IsEndOfLineDelim() {
      m.InsertAddReturn()
    } else if( kr.IsBS() || kr.IsDEL() ) {
      if( 0 < count ) { m.InsertBackspace(); }
    } else {
      m.InsertAddR( kr.R )
    }
    if( kr.IsBS() || kr.IsDEL() ) {
      if( 0 < count ) { count--; }
    } else { count++;
    }
  }
  m.Set_Insert_Mode( false )

  // Move cursor back one space:
  if( 0 < m.crsCol ) {
    m.crsCol--
    m.p_fb.Update()
  }
}

func (m *FileView) Do_a() {

  if( 0<m.p_fb.NumLines() ) {
    var CL int = m.CrsLine()
    var CC int = m.CrsChar()
    var LL int = m.p_fb.LineLen( CL )

    if( LL < CC ) {
      m.GoToCrsPos_NoWrite( CL, LL )
      m.p_fb.Update()
    } else if( CC < LL ) {
      m.GoToCrsPos_NoWrite( CL, CC+1 )
      m.p_fb.Update()
    }
  }
  m.Do_i()
}

func (m *FileView) Do_A() {

  m.GoToEndOfLine()

  m.Do_a()
}

func (m *FileView) Do_o() {

  var ONL int = m.p_fb.NumLines() //< Old number of lines
  var OCL int = m.CrsLine()       //< Old cursor line

  // New cursor line
  var NCL = 0
  if( 0 < ONL ) { NCL = OCL+1 }

  // Add the new line:
  m.p_fb.InsertLE( NCL )

  m.GoToCrsPos_NoWrite( NCL, 0 )

  m.p_fb.Update()

  m.Do_i()
}

func (m *FileView) Do_O() {
  // Add the new line:
  var new_line_num int = m.CrsLine()
  m.p_fb.InsertLE( new_line_num )
  m.crsCol = 0
  m.leftChar = 0

  m.p_fb.Update()

  m.Do_i()
}

func (m *FileView) Do_x() {

  // If there is nothing to 'x', just return:
  if( 0 < m.p_fb.NumLines() ) {

    var CL int = m.CrsLine()
    var LL int = m.p_fb.LineLen( CL )

    // If nothing on line, just return:
    if( 0 < LL )  {

      // If past end of line, move to end of line:
      if( LL-1 < m.CrsChar() ) {
        m.GoToCrsPos_Write( CL, LL-1 )
      }
      var R rune = m.p_fb.RemoveR( CL, m.CrsChar() )

      // Put char x'ed into register:
      var nlp *RLine = new( RLine )
      nlp.PushR( R )
      m_vis.reg.Clear()
      m_vis.reg.PushLP( nlp )
      m_vis.paste_mode = PM_ST_FN

      var NLL int = m.p_fb.LineLen( CL ); // New line length

      // Reposition the cursor:
      if( NLL <= m.leftChar+m.crsCol ) {
        // The char x'ed is the last char on the line, so move the cursor
        //   back one space.  Above, a char was removed from the line,
        //   but m.crsCol has not changed, so the last char is now NLL.
        // If cursor is not at beginning of line, move it back one more space.
        if( 0 < m.crsCol ) { m.crsCol-- }
      }
      m.p_fb.Update()
    }
  }
}

func (m *FileView) Do_s() {

  var CL  int = m.CrsLine()
  var LL  int = m.p_fb.LineLen( CL )
  var EOL int = LLM1( LL )
  var CP  int = m.CrsChar()

  if( CP < EOL ) {
    m.Do_x()
    m.Do_i()
  } else { // EOL <= CP
    m.Do_x()
    m.Do_a()
  }
}

func (m *FileView) Do_dw_get_fn( st_line, st_char int,
                                 fn_line, fn_char *int ) bool {

  var LL int  = m.p_fb.LineLen( st_line )
  var R  rune = m.p_fb.GetR( st_line, st_char )

  if( IsSpace( R ) ||         // On white space
      ( st_char < LLM1(LL) && // On non-white space before white space
        IsSpace( m.p_fb.GetR( st_line, st_char+1 ) ) ) ) {
    // w:
    ncp_w := CrsPos{ 0, 0 }
    var ok bool = m.GoToNextWord_GetPosition( &ncp_w )
    if( ok && 0 < ncp_w.crsChar ) { ncp_w.crsChar-- }
    if( ok && st_line == ncp_w.crsLine &&
              st_char <= ncp_w.crsChar ) {
      *fn_line = ncp_w.crsLine
      *fn_char = ncp_w.crsChar
      return true
    }
  }
  // if not on white space, and
  // not on non-white space before white space,
  // or fell through, try e:
  ncp_e := CrsPos{ 0, 0 }
  var ok bool = m.GoToEndOfWord_GetPosition( &ncp_e )

  if( ok && st_line == ncp_e.crsLine &&
            st_char <= ncp_e.crsChar ) {
    *fn_line = ncp_e.crsLine
    *fn_char = ncp_e.crsChar
    return true
  }
  return false
}

// If nothing was deleted, return 0.
// If last char on line was deleted, return 2,
// Else return 1.
func (m *FileView) Do_dw() int {

  NUM_LINES := m.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
    var st_line int = m.CrsLine()
    var st_char int = m.CrsChar()

    var LL int = m.p_fb.LineLen( st_line )

    // If past end of line, nothing to do
    if( st_char < LL ) {
      // Determine fn_line, fn_char:
      var fn_line int = 0
      var fn_char int = 0

      if( m.Do_dw_get_fn( st_line, st_char, &fn_line, &fn_char ) ) {

        m.Do_x_range( st_line, st_char, fn_line, fn_char )

        var deleted_last_char bool = (fn_char == LL-1)

        return True_1_else_2( deleted_last_char, 2, 1 )
      }
    }
  }
  return 0
}

func (m *FileView) Do_cw() {

  var result int = m.Do_dw()

  if       ( result==1 ) { m.Do_i()
  } else if( result==2 ) { m.Do_a()
  }
}

func (m *FileView) Do_D() {

  var NUM_LINES int = m.p_fb.NumLines()
  var OCL int = m.CrsLine();  // Old cursor line
  var OCP int = m.CrsChar();  // Old cursor position
  var OLL int = m.p_fb.LineLen( OCL );  // Old line length

  // If there is nothing to 'D', just return:
  if( 0<NUM_LINES && 0<OLL && OCP<OLL ) {

    var nlp *RLine = new( RLine )

    for k:=OCP; k<OLL; k++ {
      var R rune = m.p_fb.RemoveR( OCL, OCP )
      nlp.PushR( R )
    }
    m_vis.reg.Clear()
    m_vis.reg.PushLP( nlp )
    m_vis.paste_mode = PM_ST_FN

    // If cursor is not at beginning of line, move it back one space.
    if( 0 < m.crsCol ) { m.crsCol--; }

    m.p_fb.Update()
  }
}

func (m *FileView) Do_x_range( st_line, st_char, fn_line, fn_char int ) {

  m.Do_x_range_pre( &st_line, &st_char, &fn_line, &fn_char )

  if( st_line == fn_line ) {
    m.Do_x_range_single( st_line, st_char, fn_char )
  } else {
    m.Do_x_range_multiple( st_line, st_char, fn_line, fn_char )
  }
  m.Do_x_range_post( st_line, st_char )
}

func (m *FileView) Do_x_range_pre( p_st_line, p_st_char, p_fn_line, p_fn_char *int ) {

  if( m.inVisualBlock ) {
    if( *p_fn_line < *p_st_line ) { Swap( p_st_line, p_fn_line ); }
    if( *p_fn_char < *p_st_char ) { Swap( p_st_char, p_fn_char ); }
  } else {
    if( *p_fn_line < *p_st_line ||
        (*p_fn_line == *p_st_line && *p_fn_char < *p_st_char) ) {
      Swap( p_st_line, p_fn_line )
      Swap( p_st_char, p_fn_char )
    }
  }
  m_vis.reg.Clear()
}

func (m *FileView) Do_x_range_single( L, st_char, fn_char int ) {

  var OLL int = m.p_fb.LineLen( L ); // Original line length

  if( 0<OLL ) {
    var nlp *RLine = new( RLine )

    var P_st int = Min_i( st_char, OLL-1 )
    var P_fn int = Min_i( fn_char, OLL-1 )

    var LL int = OLL

    // Dont remove a single line, or else Q wont work right
    for P := P_st; P_st < LL && P <= P_fn; P++ {

      nlp.PushR( m.p_fb.RemoveR( L, P_st ) )

      LL = m.p_fb.LineLen( L ); // Removed a char, so re-calculate LL
    }
    m_vis.reg.PushLP( nlp )
  }
}

func (m *FileView) Do_x_range_multiple( st_line, st_char, fn_line, fn_char int ) {

  var started_in_middle bool = false
  var ended___in_middle bool = false

  var n_fn_line int = fn_line; // New finish line

  for L := st_line; L<=n_fn_line; {
    var nlp *RLine = new( RLine )

    var OLL int = m.p_fb.LineLen( L ); // Original line length

    var P_st int = True_1_else_2( (L==  st_line), Min_i(st_char, OLL-1), 0 )
    var P_fn int = True_1_else_2( (L==n_fn_line), Min_i(fn_char, OLL-1), OLL-1 )

    if(   st_line == L && 0    < P_st  ) { started_in_middle = true; }
    if( n_fn_line == L && P_fn < OLL-1 ) { ended___in_middle = true; }

    var LL int = OLL

    for P := P_st; P_st < LL && P <= P_fn; P++ {

      nlp.PushR( m.p_fb.RemoveR( L, P_st ) )

      LL = m.p_fb.LineLen( L ); // Removed a char, so re-calculate LL
    }
    if( 0 == P_st && OLL-1 == P_fn ) { // Removed entire line
      m.p_fb.RemoveLP( L )
      n_fn_line--
    } else {
      L++
    }
    m_vis.reg.PushLP( nlp )
  }
  if( started_in_middle && ended___in_middle ) {

    var p_fl *FLine = m.p_fb.RemoveLP( st_line+1 )
    m.p_fb.AppendLineToLine( st_line, p_fl )
  }
}

func (m *FileView) Do_x_range_post( st_line, st_char int ) {

  if( m.inVisualBlock ) { m_vis.paste_mode = PM_BLOCK
  } else                { m_vis.paste_mode = PM_ST_FN
  }
  // Try to put cursor at (st_line, st_char), but
  // make sure the cursor is in bounds after the deletion:
  var NUM_LINES int = m.p_fb.NumLines()
  var ncl int = st_line
  if( NUM_LINES <= ncl ) { ncl = NUM_LINES-1; }
  var NLL int = m.p_fb.LineLen( ncl )
  var ncc int = 0
  if( 0 < NLL ) { ncc = Min_i( NLL-1, st_char ) }

  m.GoToCrsPos_NoWrite( ncl, ncc )

//m.Set_Visual_Mode( false )

//m.p_fb.Update(); //<- No need to Undo_v() or Remove_Banner() because of this
  m.Undo_v()
}

func (m *FileView) Do_f( FAST_RUNE rune ) {

  var NUM_LINES int = m.p_fb.NumLines()

  if( 0< NUM_LINES ) {

    var OCL int = m.CrsLine();           // Old cursor line
    var LL  int = m.p_fb.LineLen( OCL ); // Line length
    var OCP int = m.CrsChar();           // Old cursor position

    if( OCP < LLM1(LL) ) {

      var NCP int = 0
      var found_rune bool = false

      for p:=OCP+1; !found_rune && p<LL; p++ {

        var R rune = m.p_fb.GetR( OCL, p )

        if( R == FAST_RUNE ) {
          NCP = p
          found_rune = true
        }
      }
      if( found_rune ) {
        m.GoToCrsPos_Write( OCL, NCP )
      }
    }
  }
}

func (m *FileView) Do_dd() {

  var ONL int = m.p_fb.NumLines(); // Old number of lines

  // If there is nothing to 'dd', just return:
  if( 1 < ONL ) {
    if( m.p_fb == m_vis.GetFileBuf( m_BE_FILE ) ) {
      m.Do_dd_BufferEditor( ONL )
    } else {
      m.Do_dd_Normal( ONL )
    }
  }
}

func (m *FileView) Do_dd_BufferEditor( ONL int ) {

  var OCL int = m.CrsLine(); // Old cursor line

  // Can only delete one of the user files out of buffer editor
  if( m_USER_FILE <= OCL ) {

    var lp *FLine = m.p_fb.GetLP( OCL )

    fname := lp.to_str()

    if( !m_vis.FileName_Is_Displayed( fname ) ) {
      m_vis.ReleaseFileName( fname )

      m.Do_dd_Normal( ONL )
    }
  }
}

func (m *FileView) Do_dd_Normal( ONL int ) {

  var OCL int = m.CrsLine();           // Old cursor line
  var OCP int = m.CrsChar();           // Old cursor position

  var DELETING_LAST_LINE bool = OCL == ONL-1

  var NCL int = True_1_else_2( DELETING_LAST_LINE, OCL-1, OCL ); // New cursor line
  var NLL int = True_1_else_2( DELETING_LAST_LINE, m.p_fb.LineLen( NCL ),
                                                   m.p_fb.LineLen( NCL + 1 ) )
  var NCP int = Min_i( OCP, LLM1( NLL ) )

  // Remove line from FileBuf and save in paste register:
  var p_fl *FLine = m.p_fb.RemoveLP( OCL )

  m_vis.reg.Clear()
  m_vis.reg.PushLP( &p_fl.runes )
  m_vis.paste_mode = PM_LINE

  m.GoToCrsPos_NoWrite( NCL, NCP )

  m.p_fb.Update()
}

func (m *FileView) Do_Star_GetNewPattern() string {

  var pattern string

  if( m.p_fb.NumLines() == 0 ) { return pattern }

  var CL int = m.CrsLine()
  var LL int = m.p_fb.LineLen( CL )

  if( 0<LL ) {
    m.MoveInBounds_Line()
    var CC int = m.CrsChar()

    var R rune = m.p_fb.GetR( CL,  CC )

    if( IsAlnum( R ) || R=='_' ) {
      pattern += string( R )

      // Search forward:
      for k:=CC+1; k<LL; k++ {
        R = m.p_fb.GetR( CL, k )
        if( IsAlnum( R ) || R=='_' ) { pattern += string( R )
        } else                       { break
        }
      }
      // Search backward:
      for k:=CC-1; 0<=k; k-- {
        R = m.p_fb.GetR( CL, k )
        if( IsAlnum( R ) || R=='_' ) { pattern = string(R) + pattern
        } else                       {  break
        }
      }
    } else {
      if( (R != ' ') &&  unicode.IsGraphic( R ) ) { pattern += string( R ) }
    }
    if( 0 < len(pattern) ) {
      pattern = string("\\b") + pattern + string("\\b")
    }
  }
  return pattern
}

// Goto next pattern or goto next dir in buffer editor or do nothing
func (m *FileView) Do_n_i() {

  if( 0 < len(m_vis.regex_str) ) {
    m.Do_n_Pattern()

  } else if( m.p_fb == m_vis.GetFileBuf( m_BE_FILE ) ) {
    m.Do_n_NextDir()
  }
}

func (m *FileView) Do_n() {

  if( nil != m.p_diff ) {
    m.p_diff.Do_n()
  } else {
    m.Do_n_i()
  }
}

// Goto previous pattern or goto previous dir in buffer editor or do nothing
func (m *FileView) Do_N_i() {

  if( 0 < len(m_vis.regex_str) ) {
    m.Do_N_Pattern()

  } else if( m.p_fb == m_vis.GetFileBuf( m_BE_FILE ) ) {
    m.Do_N_PrevDir()
  }
}

func (m *FileView) Do_N() {

  if( nil != m.p_diff ) {
    m.p_diff.Do_N()
  } else {
    m.Do_N_i()
  }
}

// Go to next pattern
func (m *FileView) Do_n_Pattern() {

  if( 0 < m.p_fb.NumLines() ) {
    var msg string = "/" + m_vis.regex_str
    m.Set_Cmd_Line_Msg( msg )

    var ncp CrsPos // Next cursor position

    if( m.Do_n_FindNextPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    } else {
      // Pattern not found, so put cursor back in view:
      m.PrintCursor()
    }
  }
}

// Go to previous pattern
func (m *FileView) Do_N_Pattern() {

  if( 0 < m.p_fb.NumLines() ) {
    var msg string = "/" + m_vis.regex_str
    m.Set_Cmd_Line_Msg( msg )

    var ncp CrsPos // Prev cursor position

    if( m.Do_N_FindPrevPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    } else {
      // Pattern not found, so put cursor back in view:
      m.PrintCursor()
    }
  }
}

func (m *FileView) Do_n_NextDir() {

  if( 1 < m.p_fb.NumLines() ) {
    m.Set_Cmd_Line_Msg("Searching down for dir")
    var dl int = m.CrsLine(); // Dir line, changed by search methods below

    var found_line bool = true

    if( m.Line_is_dir_name( dl ) ) {
      // If currently on a dir, go to next line before searching for dir
      found_line = m.Do_n_NextDir_Next_Line( &dl )
    }
    if( found_line ) {
      var found_dir bool = m.Do_n_NextDir_Search_for_Dir( &dl )

      if( found_dir ) {
        var NCL int = dl
        var NCP int = LLM1( m.p_fb.LineLen( NCL ) )

        m.GoToCrsPos_Write( NCL, NCP )
      }
    }
  }
}

func (m *FileView) Do_N_PrevDir() {

  if( 1 < m.p_fb.NumLines() ) {
    m.Set_Cmd_Line_Msg("Searching up for dir")
    var dl int = m.CrsLine(); // Dir line, changed by search methods below

    var found_line bool = true

    if( m.Line_is_dir_name( dl ) ) {
      // If currently on a dir, go to prev line before searching for dir
      found_line = m.Do_N_PrevDir_Prev_Line( &dl )
    }
    if( found_line ) {
      var found_dir bool = m.Do_N_PrevDir_Search_for_Dir( &dl )

      if( found_dir ) {
        var NCL int = dl
        var NCP int = LLM1( m.p_fb.LineLen( NCL ) )

        m.GoToCrsPos_Write( NCL, NCP )
      }
    }
  }
}

func (m *FileView) Do_n_FindNextPattern( ncp *CrsPos ) bool {

  var NUM_LINES int = m.p_fb.NumLines()

  var OCL int = m.CrsLine(); var st_l int = OCL
  var OCC int = m.CrsChar(); var st_c int = OCC

  var found_next_star bool = false

  // Move past current pattern:
  var LL int = m.p_fb.LineLen( OCL )

  m.p_fb.Check_4_New_Regex()
  m.p_fb.Find_Regexs_4_Line( OCL )
  for ; st_c<LL && m.InStarOrStarInF(OCL,st_c); st_c++ {
  }
  // If at end of current line, go down to next line:
  if( LL <= st_c ) { st_c=0; st_l++; }

  // Search for first pattern position past current position
  for l:=st_l; !found_next_star && l<NUM_LINES; l++ {

    m.p_fb.Find_Regexs_4_Line( l )

    var LL int = m.p_fb.LineLen( l )

    for p:=st_c; !found_next_star && p<LL; p++ {

      if( m.InStarOrStarInF(l,p) ) {

        found_next_star = true
        ncp.crsLine = l
        ncp.crsChar = p
      }
    }
    // After first line, always start at beginning of line
    st_c = 0
  }
  // Reached end of file and did not find any patterns,
  // so go to first pattern in file
  if( !found_next_star ) {
    for l:=0; !found_next_star && l<=OCL; l++ {
      m.p_fb.Find_Regexs_4_Line( l )

      var LL int = m.p_fb.LineLen( l )
      var END_C int = True_1_else_2( (OCL==l), Min_i( OCC, LL ), LL )

      for p:=0; !found_next_star && p<END_C; p++ {

        if( m.InStarOrStarInF(l,p) ) {
          found_next_star = true
          ncp.crsLine = l
          ncp.crsChar = p
        }
      }
    }
  }
  return found_next_star
}

func (m *FileView) Do_N_FindPrevPattern(  ncp *CrsPos ) bool {

  m.MoveInBounds_Line()

  var NUM_LINES int = m.p_fb.NumLines()

  var OCL int = m.CrsLine()
  var OCC int = m.CrsChar()

  m.p_fb.Check_4_New_Regex()

  var found_prev_star bool = false

  // Search for first star position before current position
  for l:=OCL; !found_prev_star && 0<=l; l-- {

    m.p_fb.Find_Regexs_4_Line( l )

    var LL int = m.p_fb.LineLen( l )

    var p int =LL-1
    if( OCL==l ) { p = True_1_else_2( (0<OCC), OCC-1, 0 ) }

    for ; 0<p && !found_prev_star; p-- {
      for ; 0<=p && m.InStarOrStarInF(l,p); p-- {
        found_prev_star = true
        ncp.crsLine = l
        ncp.crsChar = p
      }
    }
  }
  // Reached beginning of file and did not find any patterns,
  // so go to last pattern in file
  if( !found_prev_star ) {
    for l:=NUM_LINES-1; !found_prev_star && OCL<l; l-- {
      m.p_fb.Find_Regexs_4_Line( l )

      var LL int = m.p_fb.LineLen( l )

      var p int =LL-1
      if( OCL==l ) { p = True_1_else_2( (0<OCC), OCC-1, 0 ) }

      for ; 0<p && !found_prev_star; p-- {
        for ; 0<=p && m.InStarOrStarInF(l,p); p-- {
          found_prev_star = true
          ncp.crsLine = l
          ncp.crsChar = p
        }
      }
    }
  }
  return found_prev_star
}

func (m *FileView) Do_n_NextDir_Next_Line( dl *int ) bool {

  var NUM_LINES int = m.p_fb.NumLines()

  // Search forward for next line
  var found bool = false

  if( 1 < NUM_LINES ) {
    *dl = True_1_else_2( ( NUM_LINES-1 <= *dl ), 0, *dl+1 )

    found = true
  }
  return found
}

func (m *FileView) Do_N_PrevDir_Prev_Line( dl *int ) bool {

  var NUM_LINES int = m.p_fb.NumLines()

  // Search backward for prev line
//var found bool = false
  found := false

  if( 1 < NUM_LINES ) {
    *dl = True_1_else_2( (0 == *dl), NUM_LINES-1, *dl-1 )

    found = true
  }
  return found
}

// This function will only be called if 1<NUM_LINES
func (m *FileView) Do_n_NextDir_Search_for_Dir( dl *int ) bool {

  var found_dir bool = false

  var NUM_LINES int = m.p_fb.NumLines()
  var dl_st int = *dl

  // Search forward from dl_st:
  for ; !found_dir && *dl<NUM_LINES; {
    if( m.Line_is_dir_name( *dl ) ) { found_dir = true
    } else                          { (*dl)++
    }
  }
  if( !found_dir ) {
    // Wrap around back to top and search down to dl_st:
    *dl = 0
    for ; !found_dir && *dl<dl_st ; {

      if( m.Line_is_dir_name( *dl ) ) { found_dir = true
      } else                          { (*dl)++
      }
    }
  }
  return found_dir
}

// This function will only be called if 1<NUM_LINES
func (m *FileView) Do_N_PrevDir_Search_for_Dir( dl *int ) bool {

  found_dir := false

  NUM_LINES := m.p_fb.NumLines()
  dl_st := *dl

  // Search backward from dl_st:
  for ; !found_dir && 0<*dl; {
    if( m.Line_is_dir_name( *dl ) ) { found_dir = true
    } else                         { (*dl)--
    }
  }
  if( !found_dir && 0==*dl && m.Line_is_dir_name( *dl ) ) {
    found_dir = true
  }
  if( !found_dir ) {
    // Wrap around back to bottom and search up to dl_st:
    *dl = NUM_LINES-1
    for ; !found_dir && dl_st<*dl; {

      if( m.Line_is_dir_name( *dl ) ) { found_dir = true
      } else                          { (*dl)--
      }
    }
  }
  return found_dir
}

func (m *FileView) Do_v() bool {

  m.Set_Visual_Mode( true )
//m.Set_VisualB_Mode( false )

  return m.Do_visualMode()
}

func (m *FileView) Do_V() bool {

//m.Set_Visual_Mode( false )
  m.Set_VisualB_Mode( true )

  return m.Do_visualMode()
}

func (m *FileView) Do_yy() {
  // If there is nothing to 'yy', just return:
  if( 0<m.p_fb.NumLines() ) {
    var lp *FLine = m.p_fb.GetLP( m.CrsLine() )

    m_vis.reg.Clear()
    m_vis.reg.PushLP( &lp.runes )

    m_vis.paste_mode = PM_LINE
  }
}

func (m *FileView) Do_yw() {
  // If there is nothing to 'yw', just return:
  if( 0 < m.p_fb.NumLines() ) {
    var st_line int = m.CrsLine()
    var st_char int = m.CrsChar()

    // Determine fn_line, fn_char:
    var fn_line int = 0
    var fn_char int = 0

    if( m.Do_dw_get_fn( st_line, st_char, &fn_line, &fn_char ) ) {
      var nlp *RLine = new( RLine )
      // st_line and fn_line should be the same
      for k:=st_char; k<=fn_char; k++ {
        nlp.PushR( m.p_fb.GetR( st_line, k ) )
      }
      m_vis.reg.Clear()
      m_vis.reg.PushLP( nlp )
      m_vis.paste_mode = PM_ST_FN
    }
  }
}

func (m *FileView) Do_p() {

  PM := m_vis.paste_mode

  if       ( PM_ST_FN == PM ) { m.Do_p_or_P_st_fn( PP_After )
  } else if( PM_BLOCK == PM ) { m.Do_p_block()
  } else /*( PM_LINE  == PM*/ { m.Do_p_line()
  }
}

func (m *FileView) Do_p_or_P_st_fn( paste_pos Paste_Pos ) {

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    NLL := m_vis.reg.GetLP(k).Len()  // New line length
    OCL := m.CrsLine()               // Old cursor line

    if( 0 == k ) { // Add to current line
      m.MoveInBounds_Line()
      OLL := m.p_fb.LineLen( OCL )
      OCP := m.CrsChar()  // Old cursor position

      // If line we are pasting to is zero length, dont paste a space forward
      forward := 0
      if( (0 < OLL) && (paste_pos==PP_After) ) { forward = 1 }

      for i:=0; i<NLL; i++ {
        R := m_vis.reg.GetLP(k).GetR(i)
        m.p_fb.InsertR( OCL, OCP+i+forward, R )
      }
      if( 1 < N_REG_LINES && OCP+forward < OLL ) { // Move rest of first line onto new line below
        m.p_fb.InsertLE( OCL+1 )
        for i:=0; i<(OLL-OCP-forward); i++ {
          R := m.p_fb.RemoveR( OCL, OCP + NLL+forward )
          m.p_fb.PushR( OCL+1, R )
        }
      }
    } else if( N_REG_LINES-1 == k ) {
      // Insert a new line if at end of file:
      if( m.p_fb.NumLines() == OCL+k ) { m.p_fb.InsertLE( OCL+k ) }

      for i:=0; i<NLL; i++ {
        R := m_vis.reg.GetLP(k).GetR(i)
        m.p_fb.InsertR( OCL+k, i, R )
      }
    } else {
      // Put reg on line below:
      m.p_fb.InsertRLP( OCL+k, m_vis.reg.GetLP(k) )
    }
  }
  // Update current view after other views, so that the cursor will be put back in place
  m.p_fb.Update()
}

func (m *FileView) Do_p_block() {

  OCL := m.CrsLine()    // Old cursor line
  OCP := m.CrsChar()    // Old cursor position
  OLL := m.p_fb.LineLen( OCL ) // Old line length
  ISP := OCP+1          // Insert position
  if( 0 == OCP ) {
    if( 0 < OLL ) { ISP = 1 } else { ISP = 0 }
  }
  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    if( m.p_fb.NumLines()<=OCL+k ) { m.p_fb.InsertLE( OCL+k ) }
    LL := m.p_fb.LineLen( OCL+k )
    if( LL < ISP ) {
      // Fill in line with white space up to ISP:
      for i:=0; i<(ISP-LL); i++ {
        // Insert at end of line so undo will be atomic:
        NLL := m.p_fb.LineLen( OCL+k ) // New line length
        m.p_fb.InsertR( OCL+k, NLL, ' ' )
      }
    }
    var p_reg_line *RLine = m_vis.reg.GetLP(k)
    RLL := p_reg_line.Len()

    for i:=0; i<RLL; i++ {
      R := p_reg_line.GetR(i)
      m.p_fb.InsertR( OCL+k, ISP+i, R )
    }
  }
  m.p_fb.Update()
}

func (m *FileView) Do_p_line() {

  OCL := m.CrsLine()  // Old cursor line
  NUM_LINES := m_vis.reg.Len()

  for k:=0; k<NUM_LINES; k++ {
    // Put reg on line below:
    p_rl := new( RLine )
    p_rl.Copy( *m_vis.reg.GetLP(k) )
    m.p_fb.InsertRLP( OCL+k+1, p_rl )
  }
  // Update current view after other views,
  // so that the cursor will be put back in place
  m.p_fb.Update()
}

func (m *FileView) Do_P() {

  var PM Paste_Mode = m_vis.paste_mode

  if       ( PM_ST_FN == PM ) { m.Do_p_or_P_st_fn( PP_Before )
  } else if( PM_BLOCK == PM ) { m.Do_P_block()
  } else /*( PM_LINE  == PM*/ { m.Do_P_line()
  }
}

func (m *FileView) Do_P_block() {
  OCL := m.CrsLine()  // Old cursor line
  OCP := m.CrsChar()  // Old cursor position

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    if( m.p_fb.NumLines()<=OCL+k ) { m.p_fb.InsertLE( OCL+k ) }

    LL := m.p_fb.LineLen( OCL+k )
    if( LL < OCP ) {
      // Fill in line with white space up to OCP:
      for i:=0; i<(OCP-LL); i++ {
        m.p_fb.InsertR( OCL+k, LL, ' ' )
      }
    }
    var p_reg_line *RLine = m_vis.reg.GetLP(k)
    RLL := p_reg_line.Len()

    for i:=0; i<RLL; i++ {
      R := p_reg_line.GetR(i)
      m.p_fb.InsertR( OCL+k, OCP+i, R )
    }
  }
  m.p_fb.Update()
}

func (m *FileView) Do_P_line() {

  OCL := m.CrsLine()  // Old cursor line
  NUM_LINES := m_vis.reg.Len()

  for k:=0; k<NUM_LINES; k++ {
    // Put reg on line above:
    p_rl := new( RLine )
    p_rl.Copy( *m_vis.reg.GetLP(k) )
    m.p_fb.InsertRLP( OCL+k, p_rl )
  }
  m.p_fb.Update()
}

func (m *FileView) Do_r() {

  OCL := m.CrsLine()           // Old cursor line
  OCP := m.CrsChar()           // Old cursor position
  OLL := m.p_fb.LineLen( OCL ) // Old line length
  ISP := 0                     // Insert position
  if( 0<OLL ) { ISP = OCP+1 }

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    // Make sure file has a line where register line will be inserted:
    if( m.p_fb.NumLines() <= OCL+k ) {
      m.p_fb.InsertLE( OCL+k )
    }
    LL := m.p_fb.LineLen( OCL+k )

    // Make sure file line is as long as ISP before inserting register line:
    if( LL < ISP ) {
      // Fill in line with white space up to ISP:
      for i:=0; i<(ISP-LL); i++ {
        // Insert at end of line so undo will be atomic:
        NLL := m.p_fb.LineLen( OCL+k )  // New line length
        m.p_fb.InsertR( OCL+k, NLL, ' ' )
      }
    }
    m.Do_r_replace_white_space_with_register_line( k, OCL, ISP )
  }
  m.p_fb.Update()
}

func (m *FileView) Do_r_replace_white_space_with_register_line( k, OCL, ISP int ) {
  // Replace white space with register line, insert after white space used:
  var p_reg_line *RLine = m_vis.reg.GetLP(k)
  RLL := p_reg_line.Len()
  OLL := m.p_fb.LineLen( OCL+k )

  continue_last_update := false

  for i:=0; i<RLL; i++ {
    R_reg := p_reg_line.GetR(i)

    replaced_space := false

    if( ISP+i < OLL ) {
      R_old := m.p_fb.GetR( OCL+k, ISP+i )

      if( R_old == ' ' ) {
        // Replace ' ' with R_reg:
        m.p_fb.SetR( OCL+k, ISP+i, R_reg, continue_last_update )
        replaced_space = true
        continue_last_update = true
      }
    }
    // No more spaces or end of line, so insert:
    if( !replaced_space ) { m.p_fb.InsertR( OCL+k, ISP+i, R_reg ) }
  }
}

func (m *FileView) Do_R() {

  m.Set_Replace_Mode( true )

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE() }

  for kr := m_key.In(); !kr.IsESC(); kr = m_key.In() {

    if( kr.IsBS() || kr.IsDEL() ) {
      m.p_fb.Undo( m )
    } else {
      if( '\n' == kr.R ) { m.ReplaceAddReturn()
      } else             { m.ReplaceAddChars( kr.R )
      }
    }
  }
  m.Set_Replace_Mode( false )

  // Move cursor back one space:
  if( 0 < m.crsCol ) {
    m.crsCol--;  // Move cursor back one space.
  }
  m.p_fb.Update()
}

func (m *FileView) Do_J() {

  NUM_LINES := m.p_fb.NumLines() // Number of lines
  CL        := m.CrsLine()       // Cursor line

  ON_LAST_LINE := ( CL == NUM_LINES-1 )

  if( !ON_LAST_LINE && 2 <= NUM_LINES ) {

    m.GoToEndOfLine()

    p_fl := m.p_fb.RemoveLP( CL+1 )
    m.p_fb.AppendLineToLine( CL, p_fl )

    // Update() is less efficient than only updating part of the screen,
    //   but it makes the code simpler.
    m.p_fb.Update()
  }
}

func (m *FileView) Do_Tilda() {

  if( 0 < m.p_fb.NumLines() ) {
    OCL := m.CrsLine() // Old cursor line
    OCP := m.CrsChar() // Old cursor position
    LL  := m.p_fb.LineLen( OCL )

    if( 0 < LL && OCP < LL ) {
      R := m.p_fb.GetR( m.CrsLine(), m.CrsChar() )
      changed := false
      if       ( unicode.IsUpper( R ) ) { R = unicode.ToLower( R ); changed = true
      } else if( unicode.IsLower( R ) ) { R = unicode.ToUpper( R ); changed = true
      }
      if( m.crsCol < Min_i( LL-1, m.WorkingCols()-1 ) ) {
        if( changed ) { m.p_fb.SetR( m.CrsLine(), m.CrsChar(), R, true ) }
        // Need to move cursor right:
        m.crsCol++
      } else if( m.RightChar() < LL-1 ) {
        // Need to scroll window right:
        if( changed ) { m.p_fb.SetR( m.CrsLine(), m.CrsChar(), R, true ) }
        m.leftChar++
      } else { // m.RightChar() == LL-1
        // At end of line so cant move or scroll right:
        if( changed ) { m.p_fb.SetR( m.CrsLine(), m.CrsChar(), R, true ) }
      }
      m.p_fb.Update()
    }
  }
}

func (m *FileView) Do_u() {
  m.p_fb.Undo( m )
}

func (m *FileView) Do_U() {
  m.p_fb.UndoAll( m )
}

// Returns true if something was changed, else false
//
func (m *FileView) Do_visualMode() bool {
  changed := false
  m.MoveInBounds_Line()

  m.v_st_line = m.CrsLine();  m.v_fn_line = m.v_st_line
  m.v_st_char = m.CrsChar();  m.v_fn_char = m.v_st_char

  // Write current byte in visual:
  m.Replace_Crs_Char( &TS_VISUAL )

  for ; m.inVisualMode || m.inVisualBlock; {
    kr := m_key.In()

    if       ( kr.R == 'l' ) { m.GoRight(1)
    } else if( kr.R == 'h' ) { m.GoLeft(1)
    } else if( kr.R == 'j' ) { m.GoDown(1)
    } else if( kr.R == 'k' ) { m.GoUp(1)
    } else if( kr.R == 'H' ) { m.GoToTopLineInView()
    } else if( kr.R == 'L' ) { m.GoToBotLineInView()
    } else if( kr.R == 'M' ) { m.GoToMidLineInView()
    } else if( kr.R == 'n' ) { m.Do_n()
    } else if( kr.R == 'N' ) { m.Do_N()
    } else if( kr.R == '0' ) { m.GoToBegOfLine()
    } else if( kr.R == '$' ) { m.GoToEndOfLine()
    } else if( kr.R == 'g' ) { m.Do_v_Handle_g()
    } else if( kr.R == 'G' ) { m.GoToEndOfFile()
    } else if( kr.R == 'F' ) { m.PageDown_v()
    } else if( kr.R == 'B' ) { m.PageUp_v()
    } else if( kr.R == 'b' ) { m.GoToPrevWord()
    } else if( kr.R == 'w' ) { m.GoToNextWord()
    } else if( kr.R == 'e' ) { m.GoToEndOfWord()
    } else if( kr.R == '%' ) { m.GoToOppositeBracket()
    } else if( kr.R == 'z' ) { m_vis.Handle_z()
    } else if( kr.R == 'f' ) { m_vis.Handle_f()
    } else if( kr.R == ';' ) { m_vis.Handle_SemiColon()
    } else if( kr.R == 'y' ) { m.Do_y_v(); m.Undo_v()
    } else if( kr.R == 'Y' ) { m.Do_Y_v(); m.Undo_v()
    } else if( kr.R == 'r' ) { m.Do_r_v(); m.Undo_v()
    } else if( kr.R == 'R' ) { m.Do_R_v(); m.Undo_v()
    } else if( kr.R == 'x' ||
               kr.R == 'd' ) { m.Do_x_v();     changed = true; m.Undo_v()
    } else if( kr.R == 'D' ) { m.Do_D_v();     changed = true; m.Undo_v()
    } else if( kr.R == 's' ) { m.Do_s_v();     changed = true; m.Undo_v()
    } else if( kr.R == '~' ) { m.Do_Tilda_v(); changed = true; m.Undo_v()
    } else if( kr.IsESC() ) { m.Undo_v()
    }
  }
  return changed
}

// Returns true if still in visual mode, else false
//
func (m *FileView) Do_v_Handle_g() {

  kr := m_key.In()

  if       ( kr.R == 'g' ) { m.GoToTopOfFile()
  } else if( kr.R == '0' ) { m.GoToStartOfRow()
  } else if( kr.R == '$' ) { m.GoToEndOfRow()
  } else if( kr.R == 'f' ) { m.Do_v_Handle_gf()
  } else if( kr.R == 'p' ) { m.Do_v_Handle_gp()
  }
}

func (m *FileView) Do_y_v() {

  m_vis.reg.Clear()

  if( m.inVisualBlock ) { m.Do_y_v_block()
  } else                { m.Do_y_v_st_fn()
  }
}

func (m *FileView) Do_Y_v() {

  m_vis.reg.Clear()

  if( m.inVisualBlock ) { m.Do_y_v_block()
  } else                { m.Do_Y_v_st_fn()
  }
}

func (m *FileView) Do_r_v() {

  m_vis.reg.Clear()

  if( m.inVisualBlock ) { m.Do_r_v_block()
  } else                { m.Do_r_v_st_fn()
  }
}

func (m *FileView) Do_R_v() {

  m_vis.reg.Clear()

  if( m.inVisualBlock ) { m.Do_r_v_block()
  } else                { m.Do_R_v_st_fn()
  }
}

func (m *FileView) Do_x_v() {

  if( m.inVisualBlock ) {
    m.Do_x_range_block( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
  } else {
    m.Do_x_range( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
  }
  m.PrintCmdLine()
}

func (m *FileView) Do_D_v() {

  if( m.inVisualBlock ) {
    m.Do_x_range_block( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
    m.PrintCmdLine()
  } else {
    m.Do_D_v_line()
  }
}

func (m *FileView) Do_s_v() {

  LL := m.p_fb.LineLen( m.CrsLine() )

  CURSOR_AT_END_OF_LINE := false
  if( 0 < m.v_st_char && 0 < m.v_fn_char && 0 < LL ) {

    CURSOR_AT_END_OF_LINE = LL-1 <= m.v_st_char || LL-1 <= m.v_fn_char
  }
  was_in_visual_block :=  m.inVisualBlock

  m.Do_x_v()

  if( was_in_visual_block ) {
    if( CURSOR_AT_END_OF_LINE ) { m.Do_a_vb()
    } else                      { m.Do_i_vb()
    }
  } else {
    if( CURSOR_AT_END_OF_LINE ) { m.Do_a()
    } else                      { m.Do_i()
    }
  }
}

func (m *FileView) Do_Tilda_v() {

  m.Swap_Visual_St_Fn_If_Needed()

  if( m.inVisualBlock ) { m.Do_Tilda_v_block()
  } else                { m.Do_Tilda_v_st_fn()
  }
  m.Set_Visual_Mode( false )
  m.Undo_v() //<- This will cause the tilda'ed characters to be redrawn
}

func (m *FileView) Do_v_Handle_gf() {

  if( m.v_st_line == m.v_fn_line ) {
    m.Swap_Visual_St_Fn_If_Needed()

    fname := make( []rune, m.v_fn_char - m.v_st_char + 1 )

    for P := m.v_st_char; P<=m.v_fn_char; P++ {
      fname[P-m.v_st_char] = m.p_fb.GetR( m.v_st_line, P  )
    }
    went_to_file := m_vis.GoToBuffer_Fname( string(fname) )

    if( went_to_file ) {
      // If we made it to buffer indicated by fname, no need to Undo_v() or
      // Remove_Banner() because the whole view pane will be redrawn
      m.Set_Visual_Mode( false )
    }
  }
}

func (m *FileView) Do_v_Handle_gp() {

  if( m.v_st_line == m.v_fn_line ) {
    m.Swap_Visual_St_Fn_If_Needed()

    r_pattern := make( []rune, m.v_fn_char - m.v_st_char + 1 )

    for P := m.v_st_char; P<=m.v_fn_char; P++ {
      r_pattern[P-m.v_st_char] = m.p_fb.GetR( m.v_st_line, P  )
    }
    s_pattern := string(r_pattern)
    s_pattern_literal := regexp.QuoteMeta( s_pattern )

    m.Set_Visual_Mode( false )
    m.Undo_v()

    m_vis.Handle_Slash_GotPattern( s_pattern_literal, false )
  }
}

func (m *FileView) Do_y_v_block() {

  old_v_st_line := m.v_st_line
  old_v_st_char := m.v_st_char

  m.Swap_Visual_St_Fn_If_Needed()

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen( L )

    for P := m.v_st_char; P<LL && P <= m.v_fn_char; P++ {
      p_rl.PushR( m.p_fb.GetR( L, P ) )
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_BLOCK

  // Try to put cursor at (old_v_st_line, old_v_st_char), but
  // make sure the cursor is in bounds after the deletion:
  NUM_LINES := m.p_fb.NumLines()
  ncl := old_v_st_line
  if( NUM_LINES <= ncl ) { ncl = NUM_LINES-1 }
  NLL := m.p_fb.LineLen( ncl )
  ncc := 0
  if( 0 < NLL ) {
    ncc = old_v_st_char
    if( NLL <=  old_v_st_char ) { NLL = NLL-1 }
  }
  m.GoToCrsPos_NoWrite( ncl, ncc )
}

func (m *FileView) Do_y_v_st_fn() {

  m.Swap_Visual_St_Fn_If_Needed()

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen( L )
    if( 0 < LL ) {
      P_st := 0
      if( L == m.v_st_line ) { P_st = m.v_st_char }
      P_fn := LL-1
      if( L == m.v_fn_line ) { P_fn = Min_i(LL-1,m.v_fn_char) }

      for P := P_st; P <= P_fn; P++ {
        p_rl.PushR( m.p_fb.GetR( L, P ) )
      }
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_ST_FN
}

func (m *FileView) Do_Y_v_st_fn() {

  if( m.v_fn_line < m.v_st_line ) { Swap( &m.v_st_line, &m.v_fn_line ) }

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen(L)

    if( 0 < LL ) {
      for P := 0; P <= LL-1; P++ {
        p_rl.PushR( m.p_fb.GetR( L, P ) )
      }
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_LINE
}

func (m *FileView) Do_r_v_block() {

  m.Swap_Visual_St_Fn_If_Needed()

  old_v_st_line := m.v_st_line
  old_v_st_char := m.v_st_char

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen( L )

    continue_last_update := false

    for P := m.v_st_char; P<LL && P <= m.v_fn_char; P++ {
      p_rl.PushR( m.p_fb.GetR( L, P ) )
                  m.p_fb.SetR( L, P, ' ', continue_last_update )

      continue_last_update = true
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_BLOCK

  // Try to put cursor at (old_v_st_line, old_v_st_char), but
  // make sure the cursor is in bounds after the deletion:
  NUM_LINES := m.p_fb.NumLines()
  ncl := old_v_st_line
  if( NUM_LINES <= old_v_st_line ) { ncl = NUM_LINES-1 }

  NLL := m.p_fb.LineLen( ncl )
  ncc := old_v_st_char-1
  if       ( NLL <= 0 || old_v_st_char <= 0 ) { ncc = 0
  } else if( NLL <= old_v_st_char-1 )         { ncc = NLL-1
  }
  m.GoToCrsPos_NoWrite( ncl, ncc )
}

func (m *FileView) Do_r_v_st_fn() {

  m.Swap_Visual_St_Fn_If_Needed()

  old_v_st_line := m.v_st_line
  old_v_st_char := m.v_st_char

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen( L )
    if( 0 < LL ) {
      P_st := 0
      if( L==m.v_st_line ) { P_st = m.v_st_char }
      P_fn := LL-1
      if( L==m.v_fn_line ) { P_fn = Min_i( LL-1, m.v_fn_char ) }

      continue_last_update := false

      for P := P_st; P <= P_fn; P++ {
        p_rl.PushR( m.p_fb.GetR( L, P ) )
                    m.p_fb.SetR( L, P, ' ', continue_last_update )

        continue_last_update = true
      }
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_ST_FN

  // Try to put cursor at (old_v_st_line, old_v_st_char-1), but
  // make sure the cursor is in bounds after the deletion:
  NUM_LINES := m.p_fb.NumLines()
  ncl := old_v_st_line
  if( NUM_LINES <= old_v_st_line ) { ncl = NUM_LINES-1 }

  NLL := m.p_fb.LineLen( ncl )
  ncc := old_v_st_char-1
  if       ( NLL <= 0 || old_v_st_char <= 0 ) { ncc = 0
  } else if( NLL <= old_v_st_char-1 )         { ncc = NLL-1
  }
  m.GoToCrsPos_NoWrite( ncl, ncc )
}

func (m *FileView) Do_R_v_st_fn() {

  if( m.v_fn_line < m.v_st_line ) { Swap( &m.v_st_line, &m.v_fn_line ) }

  for L:=m.v_st_line; L<=m.v_fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen(L)

    if( 0 < LL ) {
      continue_last_update := false

      for P := 0; P <= LL-1; P++ {
        p_rl.PushR( m.p_fb.GetR( L, P ) )
                    m.p_fb.SetR( L, P, ' ', continue_last_update )

        continue_last_update = true
      }
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_LINE
}

func (m *FileView) Do_x_range_block( st_line, st_char, fn_line, fn_char int ) {

  m.Do_x_range_pre( &st_line, &st_char, &fn_line, &fn_char )

  for L := st_line; L<=fn_line; L++ {
    p_rl := new(RLine)

    LL := m.p_fb.LineLen( L )

    for P := st_char; P<LL && P <= fn_char; P++ {
      p_rl.PushR( m.p_fb.RemoveR( L, st_char ) )
    }
    m_vis.reg.PushLP( p_rl )
  }
  m.Do_x_range_post( st_line, st_char )
}

func (m *FileView) Do_D_v_line() {

  m.Swap_Visual_St_Fn_If_Needed()

  m_vis.reg.Clear()

  removed_line := false
  // 1. If m.v_st_line==0, fn_line will go negative in the loop below,
  //    so use int's instead of unsigned's
  // 2. Dont remove all lines in file to avoid crashing
  fn_line := m.v_fn_line
  for L := m.v_st_line; 1 < m.p_fb.NumLines() && L<=fn_line; fn_line-- {
    flp := m.p_fb.RemoveLP( L )
    m_vis.reg.PushLP( &flp.runes )

    // m.reg will delete nlp
    removed_line = true
  }
  m_vis.paste_mode = PM_LINE

  m.Set_Visual_Mode( false )
  // D'ed lines will be removed, so no need to Undo_v()

  if( removed_line ) {
    // Figure out and move to new cursor position:
    NUM_LINES := m.p_fb.NumLines()

    ncl := m.v_st_line
    if( NUM_LINES-1 < ncl ) {
      ncl = 0
      if( 0 < m.v_st_line ) { ncl = m.v_st_line-1 }
    }
    NCLL := m.p_fb.LineLen( ncl )
    ncc := 0
    if( 0 < NCLL ) {
      ncc = NCLL-1
      if( m.v_st_char < NCLL ) { ncc =  m.v_st_char }
    }
    m.GoToCrsPos_NoWrite( ncl, ncc )

    m.p_fb.Update()
  }
}

func (m *FileView) Do_a_vb() {

  CL := m.CrsLine()
  LL := m.p_fb.LineLen( CL )
  if( 0==LL ) { m.Do_i_vb(); return }

  CURSOR_AT_EOL := ( m.CrsChar() == LL-1 )
  if( CURSOR_AT_EOL ) {
    m.GoToCrsPos_NoWrite( CL, LL )
  }
  CURSOR_AT_RIGHT_COL := ( m.crsCol == m.WorkingCols()-1 )

  if( CURSOR_AT_RIGHT_COL ) {
    // Only need to scroll window right, and then enter insert i:
    m.leftChar++ //< This increments m.CrsChar()
  } else if( !CURSOR_AT_EOL ) { // If cursor was at EOL, already moved cursor forward
    // Only need to move cursor right, and then enter insert i:
    m.crsCol += 1 //< This increments m.CrsChar()
  }
  m.p_fb.Update()

  m.Do_i_vb()
}

func (m *FileView) Do_i_vb() {

  m.Set_Insert_Mode( true )

  count := 0
//for( char c=m.key.In(); c != ESC; c=m.key.In() )
//key,R := m_key.In()
  for kr := m_key.In(); !kr.IsESC(); kr = m_key.In() {
    if( kr.IsEndOfLineDelim() ) {
      ; // Ignore end of line delimiters
    } else if( kr.IsBS() || kr.IsDEL() ) {
      if( 0 < count ) {
        m.InsertBackspace_vb()
        count--
        m.p_fb.Update()
      }
    } else {
      m.InsertAddChar_vb( kr.R )
      count++
      m.p_fb.Update()
    }
  }
  m.Set_Insert_Mode( false )
}

func (m *FileView) Do_Tilda_v_block() {

  for L := m.v_st_line; L<=m.v_fn_line; L++ {
    LL := m.p_fb.LineLen( L )

    for P := m.v_st_char; P<LL && P <= m.v_fn_char; P++ {
      R := m.p_fb.GetR( L, P )
      changed := false
      if       ( unicode.IsUpper( R ) ) { R = unicode.ToLower( R ); changed = true
      } else if( unicode.IsLower( R ) ) { R = unicode.ToUpper( R ); changed = true
      }
      if( changed ) { m.p_fb.SetR( L, P, R, true ) }
    }
  }
}

func (m *FileView) Do_Tilda_v_st_fn() {

  for L := m.v_st_line; L<=m.v_fn_line; L++ {
    LL := m.p_fb.LineLen( L )
    P_st := 0
    if( L==m.v_st_line ) { P_st = m.v_st_char }
    P_fn := LL-1
    if( L==m.v_fn_line ) { P_fn = m.v_fn_char }

    for P := P_st; P <= P_fn; P++ {
      R := m.p_fb.GetR( L, P )
      changed := false
      if       ( unicode.IsUpper( R ) ) { R = unicode.ToLower( R ); changed = true
      } else if( unicode.IsLower( R ) ) { R = unicode.ToUpper( R ); changed = true
      }
      if( changed ) { m.p_fb.SetR( L, P, R, true ) }
    }
  }
}

func (m *FileView) Undo_v() {

  m.inVisualMode = false
  m.inVisualBlock = false

  m.p_fb.Update()
}

func (m *FileView) InsertAddR( R rune ) {

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE(); }

  m.p_fb.InsertR( m.CrsLine(), m.CrsChar(), R )

  if( m.WorkingCols() <= m.crsCol+1 ) {
    // On last working column, need to scroll right:
    m.leftChar++
  } else {
    m.crsCol += 1
  }
  m.Update_and_PrintCursor()
}

func (m *FileView) InsertAddReturn() {
  // The lines in p_fb do not end with '\n's.
  // When the file is written, '\n's are added to the ends of the lines.
  p_nl := new( RLine )
  var OLL int = m.p_fb.LineLen( m.CrsLine() );  // Old line length
  var OCP int = m.CrsChar();                    // Old cursor position

  for k:=OCP; k<OLL; k++ {
    var C rune = m.p_fb.RemoveR( m.CrsLine(), OCP )
    p_nl.PushR( C )
  }
  // Truncate the rest of the old line:
  // Add the new line:
  new_line_num := m.CrsLine()+1
  m.p_fb.InsertRLP( new_line_num, p_nl )
  m.crsCol = 0
  m.leftChar = 0
  if( m.CrsLine() < m.BotLine() ) { m.crsRow++
  } else {
    // If we were on the bottom working line, scroll screen down
    // one line so that the cursor line is not below the screen.
    m.topLine++
  }
  m.p_fb.Update()
}

func (m *FileView) InsertBackspace_RmC( OCL, OCP int ) {

  m.p_fb.RemoveR( OCL, OCP-1 )

  if( 0 < m.crsCol ) { m.crsCol -= 1
  } else             { m.leftChar -= 1; }

  m.p_fb.Update()
}

func (m *FileView) InsertBackspace_RmNL( OCL int ) {
  // Cursor Line Position is zero, so:
  // 1. Save previous line, end of line + 1 position
  ncp := CrsPos{ OCL-1, m.p_fb.LineLen( OCL-1 ) }

  // 2. Remove the line
  var p_fl *FLine = m.p_fb.RemoveLP( OCL )

  // 3. Append rest of line to previous line
  m.p_fb.AppendLineToLine( OCL-1, p_fl )

  // 4. Put cursor at the old previous line end of line + 1 position
  var MOVE_UP    bool = ncp.crsLine < m.topLine
  var MOVE_RIGHT bool = m.RightChar() < ncp.crsChar

  if( MOVE_UP ) { m.topLine = ncp.crsLine; }
                  m.crsRow = ncp.crsLine - m.topLine

  if( MOVE_RIGHT ) { m.leftChar = ncp.crsChar - m.WorkingCols() + 1; }
                     m.crsCol   = ncp.crsChar - m.leftChar

  // 5. Removed a line, so update to re-draw window view
  m.p_fb.Update()
}

func (m *FileView) InsertBackspace() {
  // If no lines in buffer, no backspacing to be done
  if 0 < m.p_fb.NumLines() {
    var OCL int = m.CrsLine(); // Old cursor line
    var OCP int = m.CrsChar(); // Old cursor position

    if( 0 < OCP ) { m.InsertBackspace_RmC ( OCL, OCP )
    } else        { m.InsertBackspace_RmNL( OCL ); }
  }
}

func (m *FileView) InsertedLine_Adjust_TopLine( l_num int ) {

  if( l_num < m.topLine ) { m.topLine++ }
}

func (m *FileView) InsertBackspace_vb() {

  OCL := m.CrsLine()  // Old cursor line
  OCP := m.CrsChar()  // Old cursor position

  if( 0<OCP ) {
    N_REG_LINES := m_vis.reg.Len()

    for k:=0; k<N_REG_LINES; k++ {
      m.p_fb.RemoveR( OCL+k, OCP-1 )
    }
    m.GoToCrsPos_NoWrite( OCL, OCP-1 )
  }
}

func (m *FileView) InsertAddChar_vb( R rune ) {

  OCL := m.CrsLine()  // Old cursor line
  OCP := m.CrsChar()  // Old cursor position

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    LL := m.p_fb.LineLen( OCL+k )

    if( LL < OCP ) {
      // Fill in line with white space up to OCP:
      for i:=0; i<(OCP-LL); i++ {
        // Insert at end of line so undo will be atomic:
        NLL := m.p_fb.LineLen( OCL+k ) // New line length
        m.p_fb.InsertR( OCL+k, NLL, ' ' )
      }
    }
    m.p_fb.InsertR( OCL+k, OCP, R )
  }
  m.GoToCrsPos_NoWrite( OCL, OCP+1 )
}

func (m *FileView) GetFileName_PartialLine() (string, bool) {

  var fname RLine
  var got_filename bool = false

  var CL int = m.CrsLine()
  var LL int = m.p_fb.LineLen( CL )

  if( 0 < LL ) {
    m.MoveInBounds_Line()
    var CP int = m.CrsChar()
    var R rune = m.p_fb.GetR( CL, CP )

    if( IsFileNameChar( R ) ) {
      // Get the file name:
      got_filename = true

      fname.PushR( R )

      // Search backwards, until non-filename char found:
      for k:=CP-1; -1<k; k-- {
        R = m.p_fb.GetR( CL, k )

        if( !IsFileNameChar( R ) ) { break
        } else { fname.InsertR( 0, R )
        }
      }
      // Search forwards, until non-filename char found:
      for k:=CP+1; k<LL; k++ {
        R = m.p_fb.GetR( CL, k )

        if( !IsFileNameChar( R ) ) { break
        } else { fname.PushR( R )
        }
      }
      // Trim white space off beginning and ending of fname:
      Trim( fname )
    }
  }
  return fname.to_str(), got_filename
}

func (m *FileView) GetFileName_WholeLine() (string, bool) {

  var fname *FLine
  var got_filename bool = false

  var CL int = m.CrsLine()
  var LL int = m.p_fb.LineLen( CL )

  if( 0 < LL ) {
    fname = m.p_fb.GetLP( CL )
    got_filename = true
  }
  return fname.to_str(), got_filename
}

func (m *FileView) Line_is_dir_name( line_num int ) bool {
  var line_is_dir bool = false

  var fname string = m.p_fb.GetLP( line_num ).to_str()

  var pfb *FileBuf = m_vis.GetFileBuf_s( fname )

  if( nil != pfb && pfb.is_dir ) {
    line_is_dir = true
  }
  return line_is_dir
}

func (m *FileView) ReplaceAddChars( R rune ) {

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE() }

  CL := m.CrsLine()
  CP := m.CrsChar()
  LL := m.p_fb.LineLen( CL )
  EOL := 0
  if( 0 < LL ) { EOL = LL-1 }

  if( EOL < CP ) {
    // Extend line out to where cursor is:
    for k:=LL; k<CP; k++ {  m.p_fb.PushR( CL, ' ' ) }
  }
  // Put char back in file buffer
  continue_last_update := false
  if( CP < LL ) { m.p_fb.SetR( CL, CP, R, continue_last_update )
  } else {
    m.p_fb.PushR( CL, R )
  }
  if( m.crsCol < m.WorkingCols()-1 ) {
    m.crsCol++
  } else {
    m.leftChar++
  }
  m.p_fb.Update()
}

func (m *FileView) ReplaceAddReturn() {

  // The lines in fb do not end with '\n's.
  // When the file is written, '\n's are added to the ends of the lines.
  OCL := m.CrsLine()
  OCP := m.CrsChar()
  OLL := m.p_fb.LineLen( OCL )
  p_new_line := new( RLine )

  for k:=OCP; k<OLL; k++ {
    R := m.p_fb.RemoveR( OCL, OCP )
    p_new_line.PushR( R )
  }
  // Truncate the rest of the old line:
  // Add the new line:
  new_line_num := OCL+1
  m.p_fb.InsertRLP( new_line_num, p_new_line )
  m.crsCol = 0
  m.leftChar = 0
  if( OCL < m.BotLine() ) { m.crsRow++
  } else {
    // If we were on the bottom working line, scroll screen down
    // one line so that the cursor line is not below the screen.
    m.topLine++
  }
  m.p_fb.Update()
}

func (m *FileView) Replace_Crs_Char( p_S *tcell.Style ) {

  LL := m.p_fb.LineLen( m.CrsLine() ) // Line length

  if( 0 < LL ) {
    R := m.p_fb.GetR( m.CrsLine(), m.CrsChar() )

    GL_ROW := m.Row_Win_2_GL( m.crsRow )
    GL_COL := m.Col_Win_2_GL( m.crsCol )

    m_console.SetR( GL_ROW, GL_COL, R, p_S )
  }
}

func (m *FileView) Swap_Visual_St_Fn_If_Needed() {

  if( m.inVisualBlock ) {
    if( m.v_fn_line < m.v_st_line ) { Swap( &m.v_st_line, &m.v_fn_line ) }
    if( m.v_fn_char < m.v_st_char ) { Swap( &m.v_st_char, &m.v_fn_char ) }
  } else {
    if( m.v_fn_line < m.v_st_line ||
        (m.v_fn_line == m.v_st_line && m.v_fn_char < m.v_st_char) ) {
      // Visual mode went backwards over multiple lines, or
      // Visual mode went backwards over one line
      Swap( &m.v_st_line, &m.v_fn_line )
      Swap( &m.v_st_char, &m.v_fn_char )
    }
  }
}

func RV_Style( S *tcell.Style ) bool {

  return S == &TS_RV_NORMAL    ||
         S == &TS_RV_STAR      ||
         S == &TS_RV_STAR_IN_F ||
         S == &TS_RV_DEFINE    ||
         S == &TS_RV_COMMENT   ||
         S == &TS_RV_CONST     ||
         S == &TS_RV_CONTROL   ||
         S == &TS_RV_VARTYPE
}

func RV_Style_2_NonRV( RVS *tcell.Style ) *tcell.Style {

  S := &TS_NORMAL

  if       ( RVS == &TS_RV_STAR     ) { S = &TS_STAR
  } else if( RVS == &TS_RV_STAR_IN_F) { S = &TS_STAR_IN_F
  } else if( RVS == &TS_RV_DEFINE   ) { S = &TS_DEFINE
  } else if( RVS == &TS_RV_COMMENT  ) { S = &TS_COMMENT
  } else if( RVS == &TS_RV_CONST    ) { S = &TS_CONST
  } else if( RVS == &TS_RV_CONTROL  ) { S = &TS_CONTROL
  } else if( RVS == &TS_RV_VARTYPE  ) { S = &TS_VARTYPE
  }
  return S
}

