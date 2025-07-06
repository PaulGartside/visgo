
package main

import (
//"bytes"
//"fmt"
  "github.com/gdamore/tcell/v2"
  "regexp"
  "unicode"
)

type LineView struct {
  p_fb *FileBuf
  x,y int                // View position on screen
  nCols int              // Number of columns in view, including border
  topLine, leftChar int  // View current top line and left character position
  crsCol int             // View current cursor column in file

  inInsertMode, inReplaceMode bool
  inVisualMode bool

  v_st_line, v_st_char int // Visual start line number, and char number on line
  v_fn_line, v_fn_char int // Visual ending line number, and char number on line

  prefix_len int
  banner_delim rune

  saved_line FLine
}

func (m *LineView) Init( p_file_buf *FileBuf, banner_delim rune ) {
  m.prefix_len = 2
  m.p_fb = p_file_buf
  m.banner_delim = banner_delim
}

func (m *LineView) CrsChar() int {
  return m.leftChar + m.crsCol
}

func (m *LineView) CrsLine() int {
  return m.topLine
}

// Translates zero based working view window column to zero based global screen column
//
func (m *LineView) Col_Win_2_GL( win_col int ) int {

  return m.x + 1 + m.prefix_len + win_col
}

// Translates zero based file line char position to zero based global column
//
func (m *LineView) Char_2_GL( line_char int ) int {

  return m.x + 1 + m.prefix_len + line_char - m.leftChar
}

// Returns working columns in view window
//
func (m *LineView) WorkingCols() int {
  return m.nCols - 2 - m.prefix_len
}

// (Index/Position) on current line of rune that is or would be displayed
// in right column
func (m *LineView) RightChar() int {
  return m.leftChar + m.WorkingCols()-1
}

func (m *LineView) PrintCursor() {

  m_console.ShowCursor( m.y, m.Col_Win_2_GL( m.crsCol ) )
  m_console.Show()
}

func (m *LineView) RepositionView() {
  // If a window re-size has taken place, and the window has gotten
  // smaller, change left char if needed, so that the
  // cursor is in the LineView when it is re-drawn
  if( m.WorkingCols() <= m.crsCol ) {
    shift := m.crsCol - m.WorkingCols() + 1
    m.leftChar += shift
    m.crsCol   -= shift
  }
}

func (m *LineView) Update() {

  m.RepositionView()
  m.DisplayBanner()
  m.PrintWorkingView()
  m.PrintCursor()
}

func (m *LineView) PrintWorkingView() {

  var WC int = m.WorkingCols()

  // Dont allow line wrap:
  var k     int = m.topLine
  var LL    int = m.p_fb.LineLen( k )
  var G_ROW int = m.y

  var col int = 0
  for i:=m.leftChar; i<LL && col<WC; i++ {
    var ru rune =  m.p_fb.GetR( k, i )

    m_console.SetR( G_ROW, m.Col_Win_2_GL( col ), ru, &TS_NORMAL )
    col++
  }
  for ; col<WC; col++ {
    m_console.SetR( G_ROW, m.Col_Win_2_GL( col ), ' ', &TS_EMPTY )
  }
}

func (m *LineView) Get_Style( line, pos int ) *tcell.Style {

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

func (m *LineView) InVisualArea( line, pos int ) bool {

  if( m.inVisualMode ) {
    return m.InVisualStFn( line, pos )
  }
  return false
}

func (m *LineView) InVisualStFn( line, pos int ) bool {

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

func (m *LineView) InStar( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_STAR )
}

func (m *LineView) InStarInF( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_STAR_IN_F )
}

//func (m *LineView) InStarOrStarInF( line, pos int ) bool {
//
//  return m.p_fb.HasStyle( line, pos, HI_STAR | HI_STAR_IN_F )
//}

func (m *LineView) InDefine( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_DEFINE )
}

func (m *LineView) InComment( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_COMMENT )
}

func (m *LineView) InConst( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_CONST )
}

func (m *LineView) InControl( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_CONTROL )
}

func (m *LineView) InVarType( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_VARTYPE )
}

func (m *LineView) InNonAscii( line, pos int ) bool {

  return m.p_fb.HasStyle( line, pos, HI_NONASCII )
}

func (m *LineView) DisplayBanner() {

  var G_COL int = m.x + 1

  if       ( m.inInsertMode  ) { m_console.SetR( m.y, G_COL, 'I', &TS_BANNER )
  } else if( m.inReplaceMode ) { m_console.SetR( m.y, G_COL, 'R', &TS_CONST )
  } else if( m.inVisualMode  ) { m_console.SetR( m.y, G_COL, 'V', &TS_DEFINE )
  } else                       { m_console.SetR( m.y, G_COL, 'E', &TS_CONTROL )
  }
  m_console.SetR( m.y, G_COL+1, m.banner_delim, &TS_NORMAL )
  m_console.Show()
}

func (m *LineView) DisplayBanner_PrintCursor() {

  m.DisplayBanner()

  m.PrintCursor()
}

func (m *LineView) SetContext( num_cols, x, y int ) {
  m.nCols    = num_cols
  m.x        = x
  m.y        = y
  m.leftChar = 0
}

func (m *LineView) GoDown() {

  var NUM_LINES int = m.p_fb.NumLines()
  var NCL       int = m.CrsLine()+1; // New cursor line

  if( 0<NUM_LINES && NCL<NUM_LINES ) {

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func (m *LineView) GoUp() {

  var NUM_LINES int = m.p_fb.NumLines()
  var OCL       int = m.CrsLine(); // Old cursor line

  if( 0<NUM_LINES && 0<OCL ) {
    var NCL int = OCL-1; // New cursor line

    m.GoToCrsPos_Write( NCL, m.CrsChar() )
  }
}

func (m *LineView) GoRight() {

  if( 0<m.p_fb.NumLines() ) {
    CL := m.CrsLine() // Cursor line
    LL := m.p_fb.LineLen( CL )
    CP := m.CrsChar() // Cursor position

    if( 0<LL && CP < LL-1 ) {
      m.GoToCrsPos_Write( CL, CP+1 )
    }
  }
}

func (m *LineView) GoLeft() {

  CP := m.CrsChar() // Cursor position

  if( 0<m.p_fb.NumLines() && 0<CP ) {
    m.GoToCrsPos_Write( m.CrsLine(), CP-1 )
  }
}

func (m *LineView) GoToCrsPos_NoWrite( ncp_crsLine, ncp_crsChar int ) {

  m.topLine = ncp_crsLine

  // These moves refer to View of buffer:
  var MOVE_RIGHT bool = m.RightChar() < ncp_crsChar
  var MOVE_LEFT  bool = ncp_crsChar < m.leftChar

  if       ( MOVE_RIGHT ) { m.leftChar = ncp_crsChar - m.WorkingCols()+1
  } else if( MOVE_LEFT  ) { m.leftChar = ncp_crsChar
  }
  m.crsCol = ncp_crsChar - m.leftChar
}

func (m *LineView) GoToCrsPos_Write( ncp_crsLine , ncp_crsChar int ) {

  var OCP int = m.CrsChar()
  var NCL int = ncp_crsLine
  var NCP int = ncp_crsChar

  if( m.topLine == NCL && OCP == NCP ) {
    m.PrintCursor();  // Put cursor back into position.
    return
  }
  if( m.inVisualMode ) {
    m.v_fn_line = NCL
    m.v_fn_char = NCP
  }
  // These moves refer to View of buffer:
  var MOVE_UP_DN bool = NCL != m.topLine
  var MOVE_RIGHT bool = m.RightChar() < NCP
  var MOVE_LEFT  bool = NCP < m.leftChar

  var redraw bool = MOVE_UP_DN || MOVE_RIGHT || MOVE_LEFT

  if( redraw ) {
    m.topLine = NCL

    if       ( MOVE_RIGHT ) { m.leftChar = NCP - m.WorkingCols() + 1
    } else if( MOVE_LEFT  ) { m.leftChar = NCP
    }
    // m.crsRow and m.crsCol must be set to new values before calling CalcNewCrsByte
    m.crsCol = NCP - m.leftChar

    m.Update()
  } else if( m.inVisualMode ) {
    m.GoToCrsPos_Write_Visual( OCP, NCP )
  } else {
    m.crsCol = NCP - m.leftChar

    m.PrintCursor();  // Put cursor into position.
  }
}

func (m *LineView) GoToCrsPos_Write_Visual( OCP, NCP int ) {
  // (old cursor pos) < (new cursor pos)
  var OCP_LT_NCP bool = OCP < NCP

  if( OCP_LT_NCP ) { // Cursor moved forward
    m.GoToCrsPos_WV_Forward( OCP, NCP )
  } else { // NCP_LT_OCP // Cursor moved backward
    m.GoToCrsPos_WV_Backward( OCP, NCP )
  }
  m.crsCol = NCP - m.leftChar
  m.PrintCursor()
}

func (m *LineView) GoToBegOfLine() {

  if( 0<m.p_fb.NumLines() ) {
    OCL := m.CrsLine() // Old cursor line

    m.GoToCrsPos_Write( OCL, 0 )
  }
}

func (m *LineView) GoToEndOfLine() {

  if( 0 < m.p_fb.NumLines() ) {

    LL  := m.p_fb.LineLen( m.CrsLine() )
    OCL := m.CrsLine() // Old cursor line

    m.GoToCrsPos_Write( OCL, LLM1( LL ) )
  }
}

func (m *LineView) GoToTopOfFile() {

  m.GoToCrsPos_Write( 0, 0 )
}

func (m *LineView) GoToEndOfFile() {

  NUM_LINES := m.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
    m.GoToCrsPos_Write( NUM_LINES-1, 0 )
  }
}

func (m *LineView) GoToStartOfRow() {

  if( 0<m.p_fb.NumLines() ) {
    OCL := m.CrsLine(); // Old cursor line

    m.GoToCrsPos_Write( OCL, m.leftChar )
  }
}

func (m *LineView) GoToEndOfRow() {

  if( 0 < m.p_fb.NumLines() ) {
    OCL := m.CrsLine(); // Old cursor line

    LL := m.p_fb.LineLen( OCL )
    if( 0 < LL ) {
      NCP := Min_i( LL-1, m.leftChar + m.WorkingCols()-1 )

      m.GoToCrsPos_Write( OCL, NCP )
    }
  }
}

// Returns true if found next word, else false
//
func (m *LineView) GoToNextWord_GetPosition( ncp *CrsPos ) bool {

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
func (m *LineView) GoToPrevWord_GetPosition( ncp *CrsPos ) bool {

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
func (m *LineView) GoToEndOfWord_GetPosition( ncp *CrsPos ) bool {

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

func (m *LineView) GoToNextWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToNextWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *LineView) GoToPrevWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToPrevWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *LineView) GoToEndOfWord() {

  ncp := CrsPos{ 0, 0 }

  if( m.GoToEndOfWord_GetPosition( &ncp ) ) {

    m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
  }
}

func (m *LineView) GoToOppositeBracket() {

  m.MoveInBounds_Line()

  NUM_LINES := m.p_fb.NumLines()
  CL := m.CrsLine()
  CC := m.CrsChar()
  LL := m.p_fb.LineLen( CL )

  if( 0<NUM_LINES && 0<LL ) {

    R := m.p_fb.GetR( CL, CC )

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

func (m *LineView) GoToOppositeBracket_Forward( ST_R, FN_R rune ) {

  NUM_LINES := m.p_fb.NumLines()
  CL := m.CrsLine()
  CC := m.CrsChar()

  // Search forward
  level := 0
  found := false

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

func (m *LineView) GoToOppositeBracket_Backward( ST_R, FN_R rune ) {

  CL := m.CrsLine()
  CC := m.CrsChar()

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

// Cursor is moving forward
// Write out from (OCL,OCP) up to but not including (NCL,NCP)
func (m *LineView) GoToCrsPos_WV_Forward( OCP, NCP int ) {

  CL := m.CrsLine()

  for k:=OCP; k<NCP; k++  {
    R := m.p_fb.GetR( CL, k )

    m_console.SetR( m.y, m.Char_2_GL( k ), R, m.Get_Style(CL,k) )
  }
}

// Cursor is moving backwards
// Write out from (OCL,OCP) back to but not including (NCL,NCP)
func (m *LineView) GoToCrsPos_WV_Backward( OCP, NCP int ) {

  CL := m.CrsLine()
  LL := m.p_fb.LineLen( CL ) // Line length

  if( 0 < LL ) {
    START := Min_i( OCP, LL-1 )

    for k:=START; NCP<k; k-- {
      R := m.p_fb.GetR( CL, k )

      m_console.SetR( m.y, m.Char_2_GL( k ), R, m.Get_Style(CL,k) )
    }
  }
}

// If past end of line, move back to end of line.
//
func (m *LineView) MoveInBounds_Line() {

  var CL  int = m.CrsLine()
  var LL  int = m.p_fb.LineLen( CL )
  var EOL int = LLM1( LL )

  if( EOL < m.CrsChar() ) {
    m.GoToCrsPos_NoWrite( CL, EOL )
  }
}

// Returns true if end of line delimiter was entered, else false
func (m *LineView) Do_i() bool {

  if( 0 == m.p_fb.NumLines() ) { m.p_fb.PushLE(); }

  m.Save_Line()

  m.inInsertMode = true
  m.Update(); //< Clear any possible message left on command line

  var LL int = m.p_fb.LineLen( m.CrsLine() );  // Line length

  // For user friendlyness, move cursor to new position immediately:
  // Since cursor is not allowed past EOL, it may need to be moved back:
  var CC int = m.CrsChar()
  if( LL < CC ) { CC = LL; }
  m.GoToCrsPos_Write( m.CrsLine(), CC )

  var EOL_DELIM_ENTERED bool = m.Do_i_tabs_or_normal(LL)

  if( !EOL_DELIM_ENTERED ) {
    // Move cursor back one space:
    if( 0 < m.crsCol ) { m.crsCol--; }

    m.inInsertMode = false
    m.DisplayBanner()
    m.PrintCursor()
  }
  return EOL_DELIM_ENTERED
}

func (m *LineView) Do_i_tabs_or_normal( LL int ) bool {

  var EOL_DELIM_ENTERED bool = false

  var CURSOR_AT_EOL bool = m.CrsChar() == LL

  if( CURSOR_AT_EOL ) { EOL_DELIM_ENTERED = m.Do_i_tabs()
  } else              { EOL_DELIM_ENTERED = m.Do_i_normal()
  }
  return EOL_DELIM_ENTERED
}

func (m *LineView) Do_i_tabs() bool {
  return m.Do_i_normal()
}

func (m *LineView) Do_i_normal() bool {
  var count int = 0
  for kr := m_key.In(); !kr.IsESC(); kr = m_key.In() {
    if( kr.IsEndOfLineDelim() ) {
      m.HandleReturn()
      return true
    } else if( kr.IsBS() || kr.IsDEL() ) {
      if( 0<count ) {
        m.InsertBackspace()
        count--
      }
    } else if( kr.IsKeyRune() ) {
      m.InsertAddR( kr.R )
      count++
    }
  }
  return false
}

func (m *LineView) Do_a() bool {

  if( 0<m.p_fb.NumLines() ) {
    CL := m.CrsLine()
    CC := m.CrsChar()
    LL := m.p_fb.LineLen( CL )

    if( LL < CC ) {
      m.GoToCrsPos_NoWrite( CL, LL )
      m.p_fb.UpdateCmd()
    } else if( CC < LL ) {
      m.GoToCrsPos_NoWrite( CL, CC+1 )
      m.p_fb.UpdateCmd()
    }
  }
  return m.Do_i()
}

func (m *LineView) Do_A() bool {
  // FIXME:
  m.GoToEndOfLine()

  return m.Do_a()
}

func (m *LineView) Do_o() bool {

  ONL := m.p_fb.NumLines() //< Old number of lines
  OCL := m.CrsLine()         //< Old cursor line

  // Add the new line:
  var NCL = 0
  if( 0 < ONL ) { NCL = OCL+1 }

  m.p_fb.InsertLE( NCL )

  m.GoToCrsPos_NoWrite( NCL, 0 )

  m.p_fb.UpdateCmd()

  return m.Do_i()
}

func (m *LineView) Do_x() {

  // If there is nothing to 'x', just return:
  if( 0 < m.p_fb.NumLines() ) {

    m.Save_Line()

    CL := m.CrsLine()
    LL := m.p_fb.LineLen( CL )

    // If nothing on line, just return:
    if( 0 < LL ) {
      // If past end of line, move to end of line:
      if( LL-1 < m.CrsChar() ) {
        m.GoToCrsPos_Write( CL, LL-1 )
      }
      R := m.p_fb.RemoveR( CL, m.CrsChar() )

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
        if( 0 < m.crsCol ) { m.crsCol--; }
      }
      m.p_fb.UpdateCmd()
    }
  }
}

func (m *LineView) Do_s() {

  m.Save_Line()
  CL  := m.CrsLine()
  LL  := m.p_fb.LineLen( CL )
  EOL := LLM1( LL )
  CP  := m.CrsChar()

  if( CP < EOL ) {
    m.Do_x()
    m.Do_i()
  } else { // EOL <= CP
    m.Do_x()
    m.Do_a()
  }
}

func (m *LineView) Do_dw_get_fn( st_line, st_char int,
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
  ok := m.GoToEndOfWord_GetPosition( &ncp_e )

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
func (m *LineView) Do_dw() int {

  NUM_LINES := m.p_fb.NumLines()

  if( 0 < NUM_LINES ) {
    m.Save_Line()
    st_line := m.CrsLine()
    st_char := m.CrsChar()

    LL := m.p_fb.LineLen( st_line )

    // If past end of line, nothing to do
    if( st_char < LL ) {
      // Determine fn_line, fn_char:
      fn_line := 0
      fn_char := 0

      if( m.Do_dw_get_fn( st_line, st_char, &fn_line, &fn_char ) ) {

        m.Do_x_range( st_line, st_char, fn_line, fn_char )

        deleted_last_char := (fn_char == LL-1)

        return True_1_else_2( deleted_last_char, 2, 1 )
      }
    }
  }
  return 0
}

func (m *LineView) Do_cw() {
  // FIXME:
  m.Save_Line()

  result := m.Do_dw()

  if       ( result==1 ) { m.Do_i()
  } else if( result==2 ) { m.Do_a()
  }
}

func (m *LineView) Do_D() {

  NUM_LINES := m.p_fb.NumLines()
  OCL := m.CrsLine()  // Old cursor line
  OCP := m.CrsChar()  // Old cursor position
  OLL := m.p_fb.LineLen( OCL )  // Old line length

  // If there is nothing to 'D', just return:
  if( 0<NUM_LINES && 0<OLL && OCP<OLL ) {
    m.Save_Line()

    nlp := new( RLine )

    for k:=OCP; k<OLL; k++ {
      R := m.p_fb.RemoveR( OCL, OCP )
      nlp.PushR( R )
    }
    m_vis.reg.Clear()
    m_vis.reg.PushLP( nlp )
    m_vis.paste_mode = PM_ST_FN

    // If cursor is not at beginning of line, move it back one space.
    if( 0 < m.crsCol ) { m.crsCol-- }

    m.p_fb.UpdateCmd()
  }
}

func (m *LineView) Do_x_range( st_line, st_char, fn_line, fn_char int ) {

  m.Do_x_range_pre( &st_line, &st_char, &fn_line, &fn_char )

  if( st_line == fn_line ) {
    m.Do_x_range_single( st_line, st_char, fn_char )
  } else {
    m.Do_x_range_multiple( st_line, st_char, fn_line, fn_char )
  }
  m.Do_x_range_post( st_line, st_char )
}

func (m *LineView) Do_x_range_pre( p_st_line, p_st_char, p_fn_line, p_fn_char *int ) {

  if( *p_fn_line < *p_st_line ||
      (*p_fn_line == *p_st_line && *p_fn_char < *p_st_char) ) {
    Swap( p_st_line, p_fn_line )
    Swap( p_st_char, p_fn_char )
  }
  m_vis.reg.Clear()
}

func (m *LineView) Do_x_range_single( L, st_char, fn_char int ) {

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

func (m *LineView) Do_x_range_multiple( st_line, st_char, fn_line, fn_char int ) {

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

func (m *LineView) Do_x_range_post( st_line, st_char int ) {

  m_vis.paste_mode = PM_ST_FN

  // Try to put cursor at (st_line, st_char), but
  // make sure the cursor is in bounds after the deletion:
  var NUM_LINES int = m.p_fb.NumLines()
  var ncl int = st_line
  if( NUM_LINES <= ncl ) { ncl = NUM_LINES-1; }
  var NLL int = m.p_fb.LineLen( ncl )
  var ncc int = 0
  if( 0 < NLL ) { ncc = Min_i( NLL-1, st_char ) }

  m.GoToCrsPos_NoWrite( ncl, ncc )

  m.inVisualMode = false

  m.p_fb.UpdateCmd(); //<- No need to Undo_v() or Remove_Banner() because of this
}

func (m *LineView) Do_f( FAST_RUNE rune ) {

  if( 0 < m.p_fb.NumLines() ) {
    OCL := m.CrsLine()           // Old cursor line
    LL  := m.p_fb.LineLen( OCL ) // Line length
    OCP := m.CrsChar()           // Old cursor position

    if( OCP < LLM1(LL) ) {
      NCP := 0
      found_char := false
      for p:=OCP+1; !found_char && p<LL; p++ {
        R := m.p_fb.GetR( OCL, p )

        if( R == FAST_RUNE ) {
          NCP = p
          found_char = true
        }
      }
      if( found_char ) {
        m.GoToCrsPos_Write( OCL, NCP )
      }
    }
  }
}

func (m *LineView) Do_dd() {

  ONL := m.p_fb.NumLines(); // Old number of lines

  // If there is nothing to 'dd', just return:
  if( 1 < ONL ) {
    m.Do_dd_Normal( ONL )
  }
}

func (m *LineView) Do_dd_Normal( ONL int ) {

  OCL := m.CrsLine();           // Old cursor line
  OCP := m.CrsChar();           // Old cursor position
//OLL := m.p_fb.LineLen( OCL ); // Old line length

  var DELETING_LAST_LINE bool = OCL == ONL-1

  var NCL int = True_1_else_2( DELETING_LAST_LINE, OCL-1, OCL ); // New cursor line
  var NLL int = True_1_else_2( DELETING_LAST_LINE, m.p_fb.LineLen( NCL ),
                                                   m.p_fb.LineLen( NCL + 1 ) )
  var NCP int = Min_i( OCP, LLM1( NLL ) )

  // Remove line from FileBuf and save in paste register:
  var p_fl *FLine = m.p_fb.RemoveLP( OCL )

  // m_vis.reg will own nlp
  m_vis.reg.Clear()
  m_vis.reg.PushLP( &p_fl.runes )
  m_vis.paste_mode = PM_LINE

  m.GoToCrsPos_NoWrite( NCL, NCP )

  m.p_fb.UpdateCmd()
}

// Go to next pattern
func (m *LineView) Do_n() {

  if( 0 < m.p_fb.NumLines() ) {
    ncp := CrsPos{ 0, 0 } // Next cursor position

    if( m.Do_n_FindNextPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    } else {
      // Pattern not found, so put cursor back in view:
      m.PrintCursor()
    }
  }
}

// Go to previous pattern
func (m *LineView) Do_N() {

  if( 0 < m.p_fb.NumLines() ) {
    ncp := CrsPos{ 0, 0 } // Next cursor position

    if( m.Do_N_FindPrevPattern( &ncp ) ) {
      m.GoToCrsPos_Write( ncp.crsLine, ncp.crsChar )
    } else {
      // Pattern not found, so put cursor back in view:
      m.PrintCursor()
    }
  }
}

func (m *LineView) Do_n_FindNextPattern( ncp *CrsPos ) bool {

  var NUM_LINES int = m.p_fb.NumLines()

  var OCL int = m.CrsLine(); var st_l int = OCL
  var OCC int = m.CrsChar(); var st_c int = OCC

  var found_next_star bool = false

  // Move past current pattern:
  var LL int = m.p_fb.LineLen( OCL )

  m.p_fb.Check_4_New_Regex()
  m.p_fb.Find_Regexs_4_Line( OCL )
  for ; st_c<LL && m.InStar(OCL,st_c); st_c++ { 
  }
  // If at end of current line, go down to next line:
  if( LL <= st_c ) { st_c=0; st_l++; }

  // Search for first pattern position past current position
  for l:=st_l; !found_next_star && l<NUM_LINES; l++ {

    m.p_fb.Find_Regexs_4_Line( l )

    var LL int = m.p_fb.LineLen( l )

    for p:=st_c; !found_next_star && p<LL; p++ {

      if( m.InStar(l,p) ) {

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

        if( m.InStar(l,p) ) {
          found_next_star = true
          ncp.crsLine = l
          ncp.crsChar = p
        }
      }
    }
  }
  return found_next_star
}

func (m *LineView) Do_N_FindPrevPattern(  ncp *CrsPos ) bool {

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
      for ; 0<=p && m.InStar(l,p); p-- {
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
        for ; 0<=p && m.InStar(l,p); p-- {
          found_prev_star = true
          ncp.crsLine = l
          ncp.crsChar = p
        }
      }
    }
  }
  return found_prev_star
}

// Returns true if something was changed, else false
//
func (m *LineView) Do_v() bool {

  if( 0 == m.p_fb.NumLines() ) { m.p_fb.PushLE() }
  m.Save_Line()

  m.inVisualMode = true
  m.DisplayBanner()

  LL := m.p_fb.LineLen( m.CrsLine() );  // Line length

  // For user friendlyness, move cursor to new position immediately:
  // Since cursor is now allowed past EOL, it may need to be moved back:
  if( LL < m.CrsChar() ) {
    m.GoToCrsPos_Write( m.CrsLine(), LL )
  }
  m.v_st_line = m.CrsLine();  m.v_fn_line = m.v_st_line
  m.v_st_char = m.CrsChar();  m.v_fn_char = m.v_st_char

  // Write current byte in visual:
  m.Replace_Crs_Char( &TS_VISUAL )

  for ; m.inVisualMode ; {
    kr := m_key.In()

    if       ( kr.R == 'l' ) { m.GoRight()
    } else if( kr.R == 'h' ) { m.GoLeft()
    } else if( kr.R == '0' ) { m.GoToBegOfLine()
    } else if( kr.R == '$' ) { m.GoToEndOfLine()
    } else if( kr.R == 'g' ) { m.Do_v_Handle_g()
    } else if( kr.R == 'b' ) { m.GoToPrevWord()
    } else if( kr.R == 'w' ) { m.GoToNextWord()
    } else if( kr.R == 'e' ) { m.GoToEndOfWord()
    } else if( kr.R == 'f' ) { L_Handle_f( &m_vis )
    } else if( kr.R == ';' ) { L_Handle_SemiColon( &m_vis )
    } else if( kr.R == 'y' ) { m.Do_y_v(); goto EXIT_VISUAL
    } else if( kr.R == 'Y' ) { m.Do_Y_v(); goto EXIT_VISUAL
    } else if( kr.R == 'x' ||
               kr.R == 'd' ) { m.Do_x_v(); return true
    } else if( kr.R == 'D' ) { m.Do_D_v(); return true
    } else if( kr.R == 's' ) { m.Do_s_v(); return true
    } else if( kr.R == '~' ) { m.Do_Tilda_v(); return true
    } else if( kr.IsESC() ) { goto EXIT_VISUAL
    }
  }
  return false

EXIT_VISUAL:
  m.inVisualMode = false
  m.DisplayBanner()
  m.Undo_v()
  return false
}

func (m *LineView) Do_yy() {
  // FIXME:
  // If there is nothing to 'yy', just return:
  if( 0 < m.p_fb.NumLines() ) {
    var lp *FLine = m.p_fb.GetLP( m.CrsLine() )

    m_vis.reg.Clear()
    m_vis.reg.PushLP( &lp.runes )

    m_vis.paste_mode = PM_LINE
  }
}

func (m *LineView) Do_yw() {
  // If there is nothing to 'yw', just return:
  if( 0 < m.p_fb.NumLines() ) {
    st_line := m.CrsLine()
    st_char := m.CrsChar()

    // Determine fn_line, fn_char:
    fn_line := 0
    fn_char := 0

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

func (m *LineView) Do_p() {

  PM := m_vis.paste_mode

  if( PM_ST_FN == PM ) {
    m.Do_p_or_P_st_fn( PP_After )
  } else { // PM_LINE == PM
    LL := m.p_fb.LineLen( m.CrsLine() )

    // If there is nothing on the current line, paste onto current line,
    // else insert into line below:
    if( 0 == LL ) { m.Do_p_or_P_st_fn( PP_After )
    } else        { m.Do_p_line()
    }
  }
}

func (m *LineView) Do_p_or_P_st_fn( paste_pos Paste_Pos ) {

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    NLL := m_vis.reg.GetLP(k).Len()  // New line length
    OCL := m.CrsLine()           // Old cursor line

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
  m.p_fb.UpdateCmd()
}

func (m *LineView) Do_p_line() {

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
  m.p_fb.UpdateCmd()
}

func (m *LineView) Do_P() {

  var PM Paste_Mode = m_vis.paste_mode

  if       ( PM_ST_FN == PM ) { m.Do_p_or_P_st_fn( PP_Before )
  } else /*( PM_LINE  == PM*/ { m.Do_P_line()
  }
}

func (m *LineView) Do_P_line() {

  OCL := m.CrsLine()  // Old cursor line
  NUM_LINES := m_vis.reg.Len()

  for k:=0; k<NUM_LINES; k++ {
    // Put reg on line above:
    p_rl := new( RLine )
    p_rl.Copy( *m_vis.reg.GetLP(k) )
    m.p_fb.InsertRLP( OCL+k, p_rl )
  }
  m.p_fb.UpdateCmd()
}

// Returns true if end of line delimiter was entered, else false
func (m *LineView) Do_R() bool {

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE() }

  m.Save_Line()

  m.inReplaceMode = true
  m.DisplayBanner_PrintCursor()

  for kr := m_key.In(); !kr.IsESC(); kr = m_key.In() {

    if( kr.IsEndOfLineDelim() ) {
      m.HandleReturn()
      return true
    } else { m.ReplaceAddChars( kr.R )
    }
  }
  // Move cursor back one space:
  if( 0 < m.crsCol ) { m.crsCol-- }

  m.inReplaceMode = false
  m.DisplayBanner_PrintCursor()

  return false
}

func (m *LineView) Do_J() {

  m.Save_Line()
  NUM_LINES := m.p_fb.NumLines() // Number of lines
  CL        := m.CrsLine()         // Cursor line

  ON_LAST_LINE := ( CL == NUM_LINES-1 )

  if( !ON_LAST_LINE && 2 <= NUM_LINES ) {

    m.GoToEndOfLine()

    p_fl := m.p_fb.RemoveLP( CL+1 )
    m.p_fb.AppendLineToLine( CL, p_fl )

    // Update() is less efficient than only updating part of the screen,
    //   but it makes the code simpler.
    m.p_fb.UpdateCmd()
  }
}

func (m *LineView) Do_Tilda() {

  if( 0 < m.p_fb.NumLines() ) {
    m.Save_Line()
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
      } else { // RightChar() == LL-1
        // At end of line so cant move or scroll right:
        if( changed ) { m.p_fb.SetR( m.CrsLine(), m.CrsChar(), R, true ) }
      }
      m.p_fb.UpdateCmd()
    }
  }
}

// Returns true if still in visual mode, else false
//
func (m *LineView) Do_v_Handle_g() {

  kr := m_key.In()

  if       ( kr.R == 'g' ) { m.GoToTopOfFile()
  } else if( kr.R == '0' ) { m.GoToStartOfRow()
  } else if( kr.R == '$' ) { m.GoToEndOfRow()
  } else if( kr.R == 'f' ) { m.Do_v_Handle_gf()
  } else if( kr.R == 'p' ) { m.Do_v_Handle_gp()
  }
}

func (m *LineView) Do_y_v() {

  m_vis.reg.Clear()

  m.Do_y_v_st_fn()
}

func (m *LineView) Do_Y_v() {

  m_vis.reg.Clear()

  m.Do_Y_v_st_fn()
}

func (m *LineView) Do_x_v() {

  m.Do_x_range( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )

  m.DisplayBanner_PrintCursor()
}

func (m *LineView) Do_D_v() {

  m.Do_D_v_line()
}

func (m *LineView) Do_s_v() {

  // Need to know if cursor is at end of line before Do_x_v() is called:
  CURSOR_AT_END_OF_LINE := m.Do_s_v_cursor_at_end_of_line()

  m.Do_x_v()

  if( CURSOR_AT_END_OF_LINE ) { m.Do_a()
  } else                      { m.Do_i()
  }
}

func (m *LineView) Do_s_v_cursor_at_end_of_line() bool {

  CURSOR_AT_END_OF_LINE := false

  LL := m.p_fb.LineLen( m.CrsLine() )

  if( 0 < LL ) { CURSOR_AT_END_OF_LINE = m.CrsChar() == LL-1
  }
  return CURSOR_AT_END_OF_LINE
}

func (m *LineView) Do_Tilda_v() {

  m.Swap_Visual_St_Fn_If_Needed()

  m.Do_Tilda_v_st_fn()

  m.inVisualMode = false
  m.DisplayBanner()
  m.Undo_v(); //<- This will cause the tilda'ed characters to be redrawn
}

func (m *LineView) Do_v_Handle_gf() {

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
      m.inVisualMode = false
    }
  }
}

func (m *LineView) Do_v_Handle_gp() {

  if( m.v_st_line == m.v_fn_line ) {
    m.Swap_Visual_St_Fn_If_Needed()

    r_pattern := make( []rune, m.v_fn_char - m.v_st_char + 1 )

    for P := m.v_st_char; P<=m.v_fn_char; P++ {
      r_pattern[P-m.v_st_char] = m.p_fb.GetR( m.v_st_line, P  )
    }
    s_pattern := string(r_pattern)
    s_pattern_literal := regexp.QuoteMeta( s_pattern )

    m.inVisualMode = false
    m.Undo_v()
    m.DisplayBanner_PrintCursor()

    m_vis.Handle_Slash_GotPattern( s_pattern_literal, false )
  }
}

func (m *LineView) Do_y_v_st_fn() {

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

func (m *LineView) Do_Y_v_st_fn() {

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

func (m *LineView) Do_D_v_line() {

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

    // m_vis.reg will delete nlp
    removed_line = true
  }
  m_vis.paste_mode = PM_LINE

  m.inVisualMode = false
  m.DisplayBanner()
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

    m.p_fb.UpdateCmd()
  }
}

func (m *LineView) Do_Tilda_v_st_fn() {

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

func (m *LineView) Undo_v() {

  m.p_fb.UpdateCmd()
}

func (m *LineView) InsertAddR( R rune ) {

  if( m.p_fb.NumLines()==0 ) { m.p_fb.PushLE(); }

  m.p_fb.InsertR( m.CrsLine(), m.CrsChar(), R )

  if( m.WorkingCols() <= m.crsCol+1 ) {
    // On last working column, need to scroll right:
    m.leftChar++
  } else {
    m.crsCol += 1
  }
  m.p_fb.UpdateCmd()
}

func (m *LineView) InsertBackspace() {

  // If no lines in buffer, no backspacing to be done
  if( 0 < m.p_fb.NumLines() ) {

    var CL int = m.CrsLine(); // Cursor line
    var CP int = m.CrsChar(); // Cursor position

    if( 0<CP ) {
      m.p_fb.RemoveR( CL, CP-1 )

      if( 0 < m.crsCol ) { m.crsCol -= 1
      } else             { m.leftChar -= 1
      }
      m.p_fb.UpdateCmd()
    }
  }
}

//func (m *LineView) Save_Line() {
//  if( nil == m.p_saved_line || 0==m.p_saved_line.Len() ) {
//    if( nil == m.p_saved_line ) {
//      m.p_saved_line = new( Line )
//    }
//    var lp *Line = m.p_fb.GetLP( m.CrsLine() )
//    m.p_saved_line.Copy( lp )
//  }
//}

func (m *LineView) Save_Line() {

  if( 0 == m.saved_line.Len() ) {
    p_fl := m.p_fb.GetLP( m.CrsLine() )

    m.saved_line.CopyP( p_fl )
  }
}

func (m *LineView) HandleReturn() bool {

  m.inInsertMode = false

  var CL int = m.topLine
  var LL int = m.p_fb.LineLen( CL ); // Current line length

  // 1. Remove current colon command line and copy it into m_rbuf:
  var p_fl *FLine = m.p_fb.RemoveLP( CL )
  m_rbuf.Clear() //< Set m_rbuf length to zero
  for k:=0; k<LL; k++ {
    m_rbuf.PushR( p_fl.GetR(k) )
  }

  // 2. If last line is blank, remove it:
  var NL int = m.p_fb.NumLines(); // Number of colon file lines
  if( 0<NL && 0 == m.p_fb.LineLen( NL-1 ) ) {
    m.p_fb.RemoveLP( NL-1 ); NL--
  }
  // 3. Remove any other lines in colon file that match current colon command:
  for ln:=0; ln<NL; ln++ {

    var lp *FLine = m.p_fb.GetLP( ln )
    if( m_rbuf.EqualL( lp.runes ) ) {
      m.p_fb.RemoveLP( ln ); NL--; ln--
    }
  }
  // 4. Add current colon command to end of colon file:
  if( 0 < m.saved_line.Len() ) {
    p_fl := new(FLine)
    p_fl.CopyP( &m.saved_line )
    m.p_fb.PushLP( p_fl ); NL++
    m.saved_line.Clear()
  }
  m.p_fb.PushLP( p_fl ); NL++

  if( 0 < LL ) { m.p_fb.PushLE(); NL++; }
  m.GoToCrsPos_NoWrite( NL-1, 0 )

  m.p_fb.UpdateCmd()

  return true
}

func (m *LineView) ReplaceAddChars( R rune ) {

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
  m.p_fb.UpdateCmd()
}

func (m *LineView) Replace_Crs_Char( p_S *tcell.Style ) {

  LL := m.p_fb.LineLen( m.CrsLine() ); // Line length

  if( 0 < LL ) {
    R := m.p_fb.GetR( m.CrsLine(), m.CrsChar() )

    GL_COL := m.Col_Win_2_GL( m.crsCol )

    m_console.SetR( m.y, GL_COL, R, p_S )
  }
}


func (m *LineView) Swap_Visual_St_Fn_If_Needed() {

  if( m.v_fn_line < m.v_st_line ||
      (m.v_fn_line == m.v_st_line && m.v_fn_char < m.v_st_char) ) {
    // Visual mode went backwards over multiple lines, or
    // Visual mode went backwards over one line
    Swap( &m.v_st_line, &m.v_fn_line )
    Swap( &m.v_st_char, &m.v_fn_char )
  }
}

