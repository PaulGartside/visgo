
package main

import (
//"bytes"
//"fmt"
//"github.com/gdamore/tcell/v2"
)

var top___border int = 1;
var bottomborder int = 1;
var left__border int = 1;
var right_border int = 1;
var NUM_ROWS     int = 1;
var WORKING_ROWS int = 1;


type LineView struct {
  p_fb *FileBuf
  x,y int                // View position on screen
  nCols int              // Number of columns in view, including border
  topLine, leftChar int  // View current top line and left character position
  crsCol int             // View current cursor column in file

  prefix_len int
  banner_delim rune

  inInsertMode bool
  inReplaceMode bool
  inVisualMode bool

  v_fn_line, v_fn_char int

  saved_line FLine
}

func (m *LineView) Init( p_file_buf *FileBuf, banner_delim rune ) {
  m.prefix_len = 2
  m.p_fb = p_file_buf
  m.banner_delim = banner_delim
}

func ( m *LineView ) WorkingCols() int {
  return m.nCols - left__border - right_border - m.prefix_len;
}

func ( m *LineView ) CrsChar() int {
  return m.leftChar + m.crsCol
}

func ( m *LineView ) CrsLine() int {
  return m.topLine
}

func ( m *LineView ) RightChar() int {
  return m.leftChar + m.WorkingCols()-1;
}

// Translates zero based working view window column to zero based global screen column
//
func ( m *LineView ) Col_Win_2_GL( win_col int ) int {

  return m.x + left__border + m.prefix_len + win_col;
}

// Translates zero based file line char position to zero based global column
//
func ( m *LineView ) Char_2_GL( line_char int ) int {

  return m.x + left__border + m.prefix_len - m.leftChar + line_char;
}

func ( m *LineView ) GetTopLine() int {
  return m.topLine;
}

func ( m *LineView ) SetTopLine( val int )  {
  m.topLine  = val;
}

func (m *LineView) SetContext( num_cols, x, y int ) {
  m.nCols    = num_cols;
  m.x        = x;
  m.y        = y;
  m.leftChar = 0;
}

//func (m *LineView) Save_Line() {
//  if( nil == m.p_saved_line || 0==m.p_saved_line.Len() ) {
//    if( nil == m.p_saved_line ) {
//      m.p_saved_line = new( Line )
//    }
//    var lp *Line = m.p_fb.GetLP( m.CrsLine() );
//    m.p_saved_line.Copy( lp )
//  }
//}

func (m *LineView) Save_Line() {

  if( 0==m.saved_line.Len() ) {
    var lp *FLine = m.p_fb.GetLP( m.CrsLine() );
    m.saved_line.CopyP( lp )
  }
}

func ( m *LineView ) Update() {

  m.RepositionView();
  m.DisplayBanner();
  m.PrintWorkingView();
  m.PrintCursor();
}

func ( m *LineView ) RepositionView() {
  // If a window re-size has taken place, and the window has gotten
  // smaller, change top line and left char if needed, so that the
  // cursor is in the FileView when it is re-drawn
  if( m.WorkingCols() <= m.crsCol ) {
    m.leftChar += ( m.crsCol - m.WorkingCols() + 1 );
    m.crsCol   -= ( m.crsCol - m.WorkingCols() + 1 );
  }
}

func ( m *LineView ) DisplayBanner() {

  var G_COL int = m.x + 1;

  if       ( m.inInsertMode  ) { m_console.SetR( m.y, G_COL, 'I', &TS_BANNER );
  } else if( m.inReplaceMode ) { m_console.SetR( m.y, G_COL, 'R', &TS_CONST );
  } else if( m.inVisualMode  ) { m_console.SetR( m.y, G_COL, 'V', &TS_DEFINE );
  } else                       { m_console.SetR( m.y, G_COL, 'E', &TS_CONTROL );
  }
  m_console.SetR( m.y, G_COL+1, m.banner_delim, &TS_NORMAL );
  m_console.Show();
}

func ( m *LineView ) PrintWorkingView() {

  var WC int = m.WorkingCols()

  // Dont allow line wrap:
  var k     int = m.topLine;
  var LL    int = m.p_fb.LineLen( k );
  var G_ROW int = m.y;

  var col int = 0;
  for i:=m.leftChar; i<LL && col<WC; i++ {
  //Style s    = Get_Style( m, k, i );
  //int   byte = m.fb.Get( k, i );
    var ru rune =  m.p_fb.GetR( k, i )

  //PrintWorkingView_Set( LL, G_ROW, col, i, byte, s );
    m_console.SetR( G_ROW, m.Col_Win_2_GL( col ), ru, &TS_NORMAL )
    col++
  }
  for ; col<WC; col++ {
    m_console.SetR( G_ROW, m.Col_Win_2_GL( col ), ' ', &TS_EMPTY );
  }
}

func ( m *LineView ) PrintCursor() {
  m_console.ShowCursor( m.y, m.Col_Win_2_GL( m.crsCol ) )
  m_console.Show()
}

func ( m *LineView ) GoToCrsPos_NoWrite( ncp_crsLine, ncp_crsChar int ) {

  m.topLine = ncp_crsLine;

  // These moves refer to View of buffer:
  var MOVE_RIGHT bool = m.RightChar() < ncp_crsChar;
  var MOVE_LEFT  bool = ncp_crsChar < m.leftChar;

  if       ( MOVE_RIGHT ) { m.leftChar = ncp_crsChar - m.WorkingCols()+1;
  } else if( MOVE_LEFT  ) { m.leftChar = ncp_crsChar;
  }
  m.crsCol = ncp_crsChar - m.leftChar;
}

func ( m *LineView ) GoToCrsPos_Write( ncp_crsLine , ncp_crsChar int ) {

  var OCP int = m.CrsChar();
  var NCL int = ncp_crsLine;
  var NCP int = ncp_crsChar;

  if( m.topLine == NCL && OCP == NCP ) {
    m.PrintCursor();  // Put cursor back into position.
    return;
  }
  if( m.inVisualMode ) {
    m.v_fn_line = NCL;
    m.v_fn_char = NCP;
  }
  // These moves refer to View of buffer:
  var MOVE_UP_DN bool = NCL != m.topLine;
  var MOVE_RIGHT bool = m.RightChar() < NCP;
  var MOVE_LEFT  bool = NCP < m.leftChar;

  var redraw bool = MOVE_UP_DN || MOVE_RIGHT || MOVE_LEFT;

  if( redraw ) {
    m.topLine = NCL;

    if       ( MOVE_RIGHT ) { m.leftChar = NCP - m.WorkingCols() + 1;
    } else if( MOVE_LEFT  ) { m.leftChar = NCP;
    }
    // m.crsRow and m.crsCol must be set to new values before calling CalcNewCrsByte
    m.crsCol = NCP - m.leftChar;

    m.Update();
  } else if( m.inVisualMode ) {
    m.GoToCrsPos_Write_Visual( OCP, NCP );
  } else {
    m.crsCol = NCP - m.leftChar;

    m.PrintCursor();  // Put cursor into position.
  }
}

func ( m *LineView ) GoToCrsPos_Write_Visual( OCP, NCP int ) {
}

// Returns true if end of line delimiter was entered, else false
func ( m *LineView ) Do_i() bool {
//Log("Top: Do_i()")

  if( 0 == m.p_fb.NumLines() ) { m.p_fb.PushLE(); }

  m.Save_Line();

  m.inInsertMode = true;
  m.Update(); //< Clear any possible message left on command line

  var LL int = m.p_fb.LineLen( m.CrsLine() );  // Line length

  // For user friendlyness, move cursor to new position immediately:
  // Since cursor is not allowed past EOL, it may need to be moved back:
  var CC int = m.CrsChar()
  if( LL < CC ) { CC = LL; }
  m.GoToCrsPos_Write( m.CrsLine(), CC );

  var EOL_DELIM_ENTERED bool = m.Do_i_tabs_or_normal(LL);

  if( !EOL_DELIM_ENTERED ) {
    // Move cursor back one space:
    if( 0 < m.crsCol ) { m.crsCol--; }

    m.inInsertMode = false;
    m.DisplayBanner();
    m.PrintCursor();
  }
//Log( fmt.Sprintf("Bot: Do_i(): %v", EOL_DELIM_ENTERED) )
  return EOL_DELIM_ENTERED;
}

func ( m *LineView ) Do_i_tabs_or_normal( LL int ) bool {

  var EOL_DELIM_ENTERED bool = false;

  var CURSOR_AT_EOL bool = m.CrsChar() == LL;

  if( CURSOR_AT_EOL ) { EOL_DELIM_ENTERED = m.Do_i_tabs();
  } else              { EOL_DELIM_ENTERED = m.Do_i_normal();
  }
  return EOL_DELIM_ENTERED;
}

func ( m *LineView ) Do_i_tabs() bool {
  return m.Do_i_normal();
}

func ( m *LineView ) Do_i_normal() bool {

  var count int = 0;

  for kr := m_key.In(); !kr.IsESC(); kr = m_key.In() {

    if( kr.IsEndOfLineDelim() ) {
      m.HandleReturn();
      return true;
    } else if( kr.IsBS() || kr.IsDEL() ) {
      if( 0<count ) {
        m.InsertBackspace();
        count--;
      }
    } else if( kr.IsKeyRune() ) {
      m.InsertAddR( kr.R );
      count++;
    }
  }
  return false;
}

func ( m *LineView ) HandleReturn() bool {

  m.inInsertMode = false;

  var CL int = m.topLine;
  var LL int = m.p_fb.LineLen( CL ); // Current line length

  // 1. Remove current colon command line and copy it into m_rbuf:
  var p_fl *FLine = m.p_fb.RemoveLP( CL );
  m_rbuf.Clear() //< Set m_rbuf length to zero
  for k:=0; k<LL; k++ {
    m_rbuf.PushR( p_fl.GetR(k) )
  }

  // 2. If last line is blank, remove it:
  var NL int = m.p_fb.NumLines(); // Number of colon file lines
  if( 0<NL && 0 == m.p_fb.LineLen( NL-1 ) ) {
    m.p_fb.RemoveLP( NL-1 ); NL--;
  }
  // 3. Remove any other lines in colon file that match current colon command:
  for ln:=0; ln<NL; ln++ {

    var lp *FLine = m.p_fb.GetLP( ln );
    if( m_rbuf.EqualL( lp.runes ) ) {
      m.p_fb.RemoveLP( ln ); NL--; ln--;
    }
  }
  // 4. Add current colon command to end of colon file:
  if( 0 < m.saved_line.Len() ) {
    p_fl := new(FLine)
    p_fl.CopyP( &m.saved_line )
    m.p_fb.PushLP( p_fl ); NL++;
    m.saved_line.Clear()
  }
  m.p_fb.PushLP( p_fl ); NL++;

  if( 0 < LL ) { m.p_fb.PushLE(); NL++; }
  m.GoToCrsPos_NoWrite( NL-1, 0 );

  m.p_fb.UpdateCmd();

  return true;
}

func ( m *LineView ) InsertBackspace() {

  // If no lines in buffer, no backspacing to be done
  if( 0 < m.p_fb.NumLines() ) {

    var CL int = m.CrsLine(); // Cursor line
    var CP int = m.CrsChar(); // Cursor position

    if( 0<CP ) {
      m.p_fb.RemoveR( CL, CP-1 );

      if( 0 < m.crsCol ) { m.crsCol -= 1;
      } else             { m.leftChar -= 1;
      }
      m.p_fb.UpdateCmd();
    }
  }
}

func ( m *LineView ) InsertAddR( R rune ) {

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE(); }

  m.p_fb.InsertR( m.CrsLine(), m.CrsChar(), R );

  if( m.WorkingCols() <= m.crsCol+1 ) {
    // On last working column, need to scroll right:
    m.leftChar++;
  } else {
    m.crsCol += 1;
  }
  m.p_fb.UpdateCmd();
}

func ( m *LineView ) GoUp() {

  var NUM_LINES int = m.p_fb.NumLines();
  var OCL       int = m.CrsLine(); // Old cursor line

  if( 0<NUM_LINES && 0<OCL ) {
    var NCL int = OCL-1; // New cursor line

    m.GoToCrsPos_Write( NCL, m.CrsChar() );
  }
}

func ( m *LineView ) GoDown() {

  var NUM_LINES int = m.p_fb.NumLines();
  var NCL       int = m.CrsLine()+1; // New cursor line

  if( 0<NUM_LINES && NCL<NUM_LINES ) {

    m.GoToCrsPos_Write( NCL, m.CrsChar() );
  }
}

func ( m *LineView ) GoLeft() {
  // FIXME:
}

func ( m *LineView ) GoRight() {
  // FIXME:
}

func ( m *LineView ) Do_a() bool {
  // FIXME:
  return false
}

func ( m *LineView ) Do_A() bool {
  // FIXME:
  return false
}

func ( m *LineView ) GoToPrevWord() {
  // FIXME:
}

func ( m *LineView ) Do_cw() {
  // FIXME:
}

func ( m *LineView ) Do_D() {
  // FIXME:
}

func ( m *LineView ) Do_dd() {
  // FIXME:
}

// If nothing was deleted, return 0.
// If last char on line was deleted, return 2,
// Else return 1.
func ( m *LineView ) Do_dw() int {
  // FIXME:
  return 0
}

func ( m *LineView ) GoToEndOfWord() {
  // FIXME:
}

func ( m *LineView ) Do_f( FAST_CHAR rune ) {
  // FIXME:
}

func ( m *LineView ) GoToTopOfFile() {
  // FIXME:
}

func ( m *LineView ) GoToEndOfFile() {
  // FIXME:
}

func ( m *LineView ) GoToStartOfRow() {
  // FIXME:
}

func ( m *LineView ) GoToEndOfRow() {
  // FIXME:
}

func ( m *LineView ) Do_J() {
  // FIXME:
}

func ( m *LineView ) Do_o() bool {
  // FIXME:
  return false
}

func ( m *LineView ) Do_n() {
  // FIXME:
}

func ( m *LineView ) Do_N() {
  // FIXME:
}

func ( m *LineView ) Do_p() {
  // FIXME:
}

func ( m *LineView ) Do_P() {
  // FIXME:
}

func ( m *LineView ) Do_R() bool {
  // FIXME:
  return false
}

func ( m *LineView ) Do_s() {
  // FIXME:
}

func ( m *LineView ) Do_v() bool {
  // FIXME:
  return false
}

func ( m *LineView ) GoToNextWord() {
  // FIXME:
}

func ( m *LineView ) Do_x() {
  // FIXME:
}

func ( m *LineView ) Do_yy() {
  // FIXME:
}

func ( m *LineView ) Do_yw() {
  // FIXME:
}

func ( m *LineView ) Do_Tilda() {
  // FIXME:
}

func ( m *LineView ) GoToEndOfLine() {
  // FIXME:
}

func ( m *LineView ) GoToOppositeBracket() {
  // FIXME:
}

func ( m *LineView ) GoToBegOfLine() {
  // FIXME:
}

