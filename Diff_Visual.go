
package main

import (
//"fmt"
  "github.com/gdamore/tcell/v2"
  "regexp"
  "unicode"
)

// Returns true if something was changed, else false
//
func (m *Diff) Do_visualMode() bool {
  pV := m_vis.CV()

  changed := false
  m.MoveInBounds_Line()

  m.v_st_line = m.CrsLine();  m.v_fn_line = m.v_st_line
  m.v_st_char = m.CrsChar();  m.v_fn_char = m.v_st_char

  // Write current byte in visual:
  m.Replace_Crs_Char( &TS_VISUAL )

  for ; pV.inVisualMode || pV.inVisualBlock; {
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

func (m *Diff) Set_Visual_Mode( on bool ) {

  pV := m_vis.CV()

  if( on != pV.inVisualMode ) {
    pV.inVisualMode = on
    pV.PrintCmdLine()
    m_console.Show()
  }
}

func (m *Diff) Set_VisualB_Mode( on bool ) {

  pV := m_vis.CV()

  if( on != pV.inVisualBlock ) {
    pV.inVisualBlock = on
    pV.PrintCmdLine()
    m_console.Show()
  }
}

func (m *Diff) Replace_Crs_Char( p_S *tcell.Style ) {

  pV := m_vis.CV()
  pfb := pV.p_fb

  // Convert CL, which is diff line, to view line:
  CLv := m.ViewLine( pV, m.CrsLine() )

  LL := pfb.LineLen( CLv ) // Line length
  if( 0 < LL ) {
    R := pfb.GetR( CLv, m.CrsChar() )

    GL_ROW := pV.Row_Win_2_GL( m.crsRow )
    GL_COL := pV.Col_Win_2_GL( m.crsCol )

    m_console.SetR( GL_ROW, GL_COL, R, p_S )
  }
}

func (m *Diff) Do_v_Handle_g() {

  kr := m_key.In()

  if       ( kr.R == 'g' ) { m.GoToTopOfFile()
  } else if( kr.R == '0' ) { m.GoToStartOfRow()
  } else if( kr.R == '$' ) { m.GoToEndOfRow()
  } else if( kr.R == 'f' ) { m.Do_v_Handle_gf()
  } else if( kr.R == 'p' ) { m.Do_v_Handle_gp()
  }
}

func (m *Diff) PageDown_v() {

  NUM_LINES := m.NumLines()

  if( 0<NUM_LINES ) {
    OCLd := m.CrsLine() // Old cursor line diff

    NCLd := OCLd + m.WorkingRows( m_vis.CV() ) - 1 // New cursor line diff

    // Dont let cursor go past the end of the file:
    if( NUM_LINES-1 < NCLd ) { NCLd = NUM_LINES-1 }

    m.GoToCrsPos_Write( NCLd, 0 )
  }
}

// This one works better when IN visual mode:
func (m *Diff) PageUp_v() {

  NUM_LINES := m.NumLines()

  if( 0<NUM_LINES ) {
    OCLd := m.CrsLine() // Old cursor line diff

    NCLd := OCLd - m.WorkingRows( m_vis.CV() ) + 1 // New cursor line diff

    // Check for underflow:
    if( NCLd < 0 ) { NCLd = 0 }

    m.GoToCrsPos_Write( NCLd, 0 )
  }
}

func (m *Diff) Do_y_v() {

  pV := m_vis.CV()

  m_vis.reg.Clear()

  if( pV.inVisualBlock ) { m.Do_y_v_block()
  } else                 { m.Do_y_v_st_fn()
  }
}

func (m *Diff) Do_Y_v() {

  pV := m_vis.CV()

  m_vis.reg.Clear()

  if( pV.inVisualBlock ) { m.Do_y_v_block()
  } else                 { m.Do_Y_v_st_fn()
  }
}

func (m *Diff) Undo_v() {

  pV := m_vis.CV()

  pV.inVisualMode = false
  pV.inVisualBlock = false

  m.Update1V( pV )
}

func (m *Diff) Do_x_v() {

  pV := m_vis.CV()

  if( pV.inVisualBlock ) {
    m.Do_x_range_block( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
  } else {
    m.Do_x_range( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
  }
  // Visual mode is exited:
  pV.PrintCmdLine()
}

func (m *Diff) Do_D_v() {

  pV := m_vis.CV()

  if( pV.inVisualBlock ) {
    m.Do_x_range_block( m.v_st_line, m.v_st_char, m.v_fn_line, m.v_fn_char )
    // Visual block mode is exited:
    pV.PrintCmdLine()
  } else {
    m.Do_D_v_line()
  }
}

func (m *Diff) Do_s_v() {

  pV := m_vis.CV()
  // Need to know if cursor is at end of line before Do_x_v() is called:
  CURSOR_AT_END_OF_LINE := m.Do_s_v_cursor_at_end_of_line()

  m.Do_x_v()

  if( pV.inVisualBlock ) {
    if( CURSOR_AT_END_OF_LINE ) { m.Do_a_vb()
    } else                      { m.Do_i_vb()
    }
  } else {
    if( CURSOR_AT_END_OF_LINE ) { m.Do_a()
    } else                      { m.Do_i()
    }
  }
  pV.inVisualMode = false
}

func (m *Diff) Do_s_v_cursor_at_end_of_line() bool {

  cursor_at_end_of_line := false

  pV  := m_vis.CV()
  pfb := pV.p_fb

  DL := m.CrsLine()  // Diff line
  VL := m.ViewLine( pV, DL )
  LL := pfb.LineLen( VL )

  if( pV.inVisualBlock ) {
    if( 0 < LL ) { cursor_at_end_of_line = LL-1 <= m.CrsChar()
    } else       { cursor_at_end_of_line = 0    <  m.CrsChar()
    }
  } else {
    if( 0 < LL ) {
      cursor_at_end_of_line = LL-1 <= m.CrsChar()
    }
  }
  return cursor_at_end_of_line
}

func (m *Diff) Do_Tilda_v() {

  pV  := m_vis.CV()

  m.Swap_Visual_St_Fn_If_Needed()

  if( pV.inVisualBlock ) { m.Do_Tilda_v_block()
  } else                 { m.Do_Tilda_v_st_fn()
  }
  m.Set_Visual_Mode( false )
  m.Undo_v() //<- This will cause the tilda'ed characters to be redrawn
}

func (m *Diff) Do_v_Handle_gf() {

  if( m.v_st_line == m.v_fn_line ) {
    pV  := m_vis.CV()
    pfb := pV.p_fb

    m.Swap_Visual_St_Fn_If_Needed()

    VL := m.ViewLine( pV, m.v_st_line )

    fname := make( []rune, m.v_fn_char - m.v_st_char + 1 )

    for P := m.v_st_char; P<=m.v_fn_char; P++ {
      fname[P-m.v_st_char] = pfb.GetR( VL, P )
    }
    went_to_file := m_vis.GoToBuffer_Fname( string(fname) )

    if( went_to_file ) {
      // If we made it to buffer indicated by fname, no need to Undo_v() or
      // Remove_Banner() because the whole view pane will be redrawn
      m.Set_Visual_Mode( false )
    }
  }
}

func (m *Diff) Do_v_Handle_gp() {

  if( m.v_st_line == m.v_fn_line ) {
    pV  := m_vis.CV()
    pfb := pV.p_fb

    m.Swap_Visual_St_Fn_If_Needed()

    VL := m.ViewLine( pV, m.v_st_line )

    r_pattern := make( []rune, m.v_fn_char - m.v_st_char + 1 )

    for P := m.v_st_char; P<=m.v_fn_char; P++ {
      r_pattern[P-m.v_st_char] = pfb.GetR( VL, P  )
    }
    s_pattern := string(r_pattern)
    s_pattern_literal := regexp.QuoteMeta( s_pattern )

    m.Set_Visual_Mode( false )
    m.Undo_v()

    m_vis.Handle_Slash_GotPattern( s_pattern_literal, false )
  }
}

func (m *Diff) Do_y_v_block() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  old_v_st_line := m.v_st_line
  old_v_st_char := m.v_st_char

  m.Swap_Visual_St_Fn_If_Needed()

  for DL:=m.v_st_line; DL<=m.v_fn_line; DL++ {
    p_rl := new(RLine)

    VL := m.ViewLine( pV, DL )
    LL := pfb.LineLen( VL )

    for P := m.v_st_char; P<LL && P <= m.v_fn_char; P++ {
      p_rl.PushR( pfb.GetR( VL, P ) )
    }
    m_vis.reg.PushLP( p_rl )
  }
  m_vis.paste_mode = PM_BLOCK

  // Try to put cursor at (old_v_st_line, old_v_st_char), but
  // make sure the cursor is in bounds after the deletion:
  NUM_LINES := m.NumLines()
  ncl := old_v_st_line
  if( NUM_LINES <= ncl ) { ncl = NUM_LINES-1 }
  NLL := pfb.LineLen( m.ViewLine( pV, ncl ) )
  ncc := 0
  if( 0 < NLL ) {
    ncc = old_v_st_char
    if( NLL <=  old_v_st_char ) { NLL = NLL-1 }
  }
  m.GoToCrsPos_NoWrite( ncl, ncc )
}

func (m *Diff) Do_y_v_st_fn() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  m.Swap_Visual_St_Fn_If_Needed()

  for DL:=m.v_st_line; DL<=m.v_fn_line; DL++ {

    LINE_DIFF_TYPE := m.DiffType( pV, DL )
    if( LINE_DIFF_TYPE != DT_DELETED ) {
      p_rl := new(RLine)

      // Convert DL, which is diff line, to view line
      VL := m.ViewLine( pV, DL )
      LL := pfb.LineLen( VL )
      if( 0 < LL ) {
        P_st := 0
        if( DL == m.v_st_line ) { P_st = m.v_st_char }
        P_fn := LL-1
        if( DL == m.v_fn_line ) { P_fn = Min_i(LL-1,m.v_fn_char) }

        for P := P_st; P <= P_fn; P++ {
          p_rl.PushR( pfb.GetR( DL, P ) )
        }
      }
      m_vis.reg.PushLP( p_rl )
    }
  }
  m_vis.paste_mode = PM_ST_FN
}

func (m *Diff) Do_Y_v_st_fn() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  if( m.v_fn_line < m.v_st_line ) { Swap( &m.v_st_line, &m.v_fn_line ) }

  for DL:=m.v_st_line; DL<=m.v_fn_line; DL++ {

    LINE_DIFF_TYPE := m.DiffType( pV, DL )
    if( LINE_DIFF_TYPE != DT_DELETED ) {
      p_rl := new(RLine)

      // Convert DL, which is diff line, to view line
      VL := m.ViewLine( pV, DL )
      LL := pfb.LineLen( VL )

      if( 0 < LL ) {
        for P := 0; P <= LL-1; P++ {
          p_rl.PushR( pfb.GetR( VL, P ) )
        }
      }
      m_vis.reg.PushLP( p_rl )
    }
  }
  m_vis.paste_mode = PM_LINE
}

func (m *Diff) Do_x_range_block( st_line, st_char, fn_line, fn_char int ) {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  m.Do_x_range_pre( &st_line, &st_char, &fn_line, &fn_char )

  for L := st_line; L<=fn_line; L++ {
    p_rl := new(RLine)

    LL := pfb.LineLen( L )

    for P := st_char; P<LL && P <= fn_char; P++ {
      p_rl.PushR( pfb.RemoveR( L, st_char ) )
    }
    m_vis.reg.PushLP( p_rl )
  }
  m.Do_x_range_post( st_line, st_char )
}

func (m *Diff) Do_D_v_line() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  m.Swap_Visual_St_Fn_If_Needed()

  m_vis.reg.Clear()

  removed_line := false
  // 1. If m.v_st_line==0, fn_line will go negative in the loop below,
  //    so use int's instead of unsigned's
  // 2. Dont remove all lines in file to avoid crashing
  fn_line := m.v_fn_line
  for DL := m.v_st_line; 1 < pfb.NumLines() && DL<=fn_line; fn_line-- {
    VL := m.ViewLine( pV, DL )
    flp := pfb.RemoveLP( VL )
    m_vis.reg.PushLP( &flp.runes )

    // m.reg will delete nlp
    removed_line = true
  }
  m_vis.paste_mode = PM_LINE

  m.Set_Visual_Mode( false )
  // D'ed lines will be removed, so no need to Undo_v()

  if( removed_line ) {
    // Figure out and move to new cursor position:
    NUM_LINES := pfb.NumLines()

    ncl := m.v_st_line
    if( NUM_LINES-1 < ncl ) {
      ncl = 0
      if( 0 < m.v_st_line ) { ncl = m.v_st_line-1 }
    }
    NCLL := pfb.LineLen( ncl )
    ncc := 0
    if( 0 < NCLL ) {
      ncc = NCLL-1
      if( m.v_st_char < NCLL ) { ncc =  m.v_st_char }
    }
    m.GoToCrsPos_NoWrite( ncl, ncc )

    m.Update1V( pV )
  }
}

func (m *Diff) Do_a_vb() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  DL := m.CrsLine()
  VL := m.ViewLine( pV, DL ) // View line number
  LL := pfb.LineLen( VL )

  if( 0==LL ) { m.Do_i_vb(); return }

  CURSOR_AT_EOL := ( m.CrsChar() == LL-1 )
  if( CURSOR_AT_EOL ) {
    m.GoToCrsPos_NoWrite( DL, LL )
  }
  CURSOR_AT_RIGHT_COL := ( m.crsCol == m.WorkingCols( pV )-1 )

  if( CURSOR_AT_RIGHT_COL ) {
    // Only need to scroll window right, and then enter insert i:
    m.leftChar++ //< This increments m.CrsChar()
  } else if( !CURSOR_AT_EOL ) { // If cursor was at EOL, already moved cursor forward
    // Only need to move cursor right, and then enter insert i:
    m.crsCol += 1 //< This increments m.CrsChar()
  }
  m.Update1V( pV )

  m.Do_i_vb()
}

func (m *Diff) Do_i_vb() {

  pV  := m_vis.CV()

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
        m.Update1V( pV )
      }
    } else {
      m.InsertAddChar_vb( kr.R )
      count++
      m.Update1V( pV )
    }
  }
  m.Set_Insert_Mode( false )
}

func (m *Diff) Do_Tilda_v_block() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  for DL := m.v_st_line; DL<=m.v_fn_line; DL++ {
    VL := m.ViewLine( pV, DL )
    LL := pfb.LineLen( VL )

    for P := m.v_st_char; P<LL && P <= m.v_fn_char; P++ {
      R := pfb.GetR( VL, P )
      changed := false
      if       ( unicode.IsUpper( R ) ) { R = unicode.ToLower( R ); changed = true
      } else if( unicode.IsLower( R ) ) { R = unicode.ToUpper( R ); changed = true
      }
      if( changed ) { pfb.SetR( VL, P, R, true ) }
    }
  }
}

func (m *Diff) Do_Tilda_v_st_fn() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  for DL := m.v_st_line; DL<=m.v_fn_line; DL++ {
    VL := m.ViewLine( pV, DL )
    LL := pfb.LineLen( VL )
    P_st := 0
    if( DL==m.v_st_line ) { P_st = m.v_st_char }
    P_fn := LLM1( LL )
    if( DL==m.v_fn_line ) { P_fn = m.v_fn_char }

    for P := P_st; P <= P_fn; P++ {
      R := pfb.GetR( VL, P )
      changed := false
      if       ( unicode.IsUpper( R ) ) { R = unicode.ToLower( R ); changed = true
      } else if( unicode.IsLower( R ) ) { R = unicode.ToUpper( R ); changed = true
      }
      if( changed ) { pfb.SetR( VL, P, R, true ) }
    }
  }
}

func (m *Diff) Set_Insert_Mode( on bool ) {

  pV := m_vis.CV()

  if( on != pV.inInsertMode ) {
    pV.inInsertMode = on
    m.PrintCmdLine( pV )
    m_console.Show()
  }
}

func (m *Diff) InsertBackspace_vb() {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  DL := m.CrsLine()          // Old cursor line
  VL := m.ViewLine( pV, DL ) // View line number
  CP := m.CrsChar()          // Old cursor position

  if( 0<CP ) {
    N_REG_LINES := m_vis.reg.Len()

    for k:=0; k<N_REG_LINES; k++ {
      pfb.RemoveR( VL+k, CP-1 )

      m.Patch_Diff_Info_Changed( pV, DL+k )
    }
    m.GoToCrsPos_NoWrite( DL, CP-1 )
  }
}

func (m *Diff) InsertAddChar_vb( R rune ) {

  pV  := m_vis.CV()
  pfb := pV.p_fb

  DL := m.CrsLine()          // Old cursor line
  VL := m.ViewLine( pV, DL ) // View line number
  CP := m.CrsChar()         // Old cursor position

  N_REG_LINES := m_vis.reg.Len()

  for k:=0; k<N_REG_LINES; k++ {
    LL := pfb.LineLen( VL+k )

    if( LL < CP ) {
      // Fill in line with white space up to OCP:
      for i:=0; i<(CP-LL); i++ {
        // Insert at end of line so undo will be atomic:
        NLL := pfb.LineLen( VL+k ) // New line length
        pfb.InsertR( VL+k, NLL, ' ' )
      }
    }
    pfb.InsertR( VL+k, CP, R )

    m.Patch_Diff_Info_Changed( pV, DL+k )
  }
  m.GoToCrsPos_NoWrite( DL, CP+1 )
}

func (m *Diff) GoToCrsPos_Write_VisualBlock( OCL, OCP, NCL, NCP int ) {

  pV  := m_vis.CV()
  pfb := pV.p_fb
  // m.v_fn_line == NCL && m.v_fn_char == NCP, so dont need to include
  // m.v_fn_line       and m.v_fn_char in Min and Max calls below:
  vis_box_left := Min_i( m.v_st_char, Min_i( OCP, NCP ) )
  vis_box_rite := Max_i( m.v_st_char, Max_i( OCP, NCP ) )
  vis_box_top  := Min_i( m.v_st_line, Min_i( OCL, NCL ) )
  vis_box_bot  := Max_i( m.v_st_line, Max_i( OCL, NCL ) )

  draw_box_left := Max_i( m.leftChar     , vis_box_left )
  draw_box_rite := Min_i( m.RightChar(pV), vis_box_rite )
  draw_box_top  := Max_i( m.topLine      , vis_box_top  )
  draw_box_bot  := Min_i( m.BotLine(pV)  , vis_box_bot  )

  for DL:=draw_box_top; DL<=draw_box_bot; DL++ {
    VL := m.ViewLine( pV, DL ); // View line number

    LL := pfb.LineLen( VL )

    for k:=draw_box_left; k<LL && k<=draw_box_rite; k++ {
      R := pfb.GetR( VL, k )
      var p_TS *tcell.Style = m.Get_Style( pV, DL, VL, k )

      m_console.SetR( pV.Line_2_GL( DL ), pV.Char_2_GL( k ), R, p_TS )
    }
  }
  m.crsRow = NCL - m.topLine
  m.crsCol = NCP - m.leftChar

  m.PrintCursor( pV ) // Does Console::Update()
}

func (m *Diff) GoToCrsPos_Write_Visual( OCL, OCP, NCL, NCP int ) {

  pV  := m_vis.CV()

  // (old cursor pos) < (new cursor pos)
  OCP_LT_NCP := OCL < NCL || (OCL == NCL && OCP < NCP)

  if( OCP_LT_NCP ) { // Cursor moved forward
    m.GoToCrsPos_WV_Forward( OCL, OCP, NCL, NCP )
  } else { // NCP_LT_OCP // Cursor moved backward
    m.GoToCrsPos_WV_Backward( OCL, OCP, NCL, NCP )
  }
  m.crsRow = NCL - m.topLine
  m.crsCol = NCP - m.leftChar
//Console::Update()
  m.PrintCursor( pV )
//m.sts_line_needs_update = true;
}

// Cursor is moving forward
// Write out from (OCL,OCP) up to but not including (NCL,NCP)
func (m *Diff) GoToCrsPos_WV_Forward( OCL, OCP, NCL, NCP int ) {

  pV  := m_vis.CV()
  pfb := pV.p_fb
  // Convert OCL and NCL, which are diff lines, to view lines:
  OCLv := m.ViewLine( pV, OCL )
  NCLv := m.ViewLine( pV, NCL )

  if( OCL == NCL ) { // Only one line:
    for k:=OCP; k<NCP; k++ {
      R := pfb.GetR( OCLv, k )
      var p_TS *tcell.Style = m.Get_Style( pV, OCL, OCLv, k )
      m_console.SetR( pV.Line_2_GL( OCL ), pV.Char_2_GL( k ), R, p_TS )
    }
  } else { // Multiple lines
    // Write out first line:
    var FIRST_LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, OCL )
    if( FIRST_LINE_DIFF_TYPE != DT_DELETED ) {
      OCLL := pfb.LineLen( OCLv ) // Old cursor line length
      END_FIRST_LINE := Min_i( m.RightChar( pV )+1, OCLL )
      for k:=OCP; k<END_FIRST_LINE; k++ {
        R := pfb.GetR( OCLv, k )
        var p_TS *tcell.Style = m.Get_Style( pV, OCL, OCLv, k )
        m_console.SetR( pV.Line_2_GL( OCL ), pV.Char_2_GL( k ), R, p_TS )
      }
    }
    // Write out intermediate lines:
    for l:=OCL+1; l<NCL; l++ {
      var LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, l )
      if( LINE_DIFF_TYPE != DT_DELETED ) {
        // Convert OCL, which is diff line, to view line
        Vl := m.ViewLine( pV, l )
        LL := pfb.LineLen( Vl ) // Line length
        END_OF_LINE := Min_i( m.RightChar( pV )+1, LL )
        for k:=m.leftChar; k<END_OF_LINE; k++ {
          R := pfb.GetR( Vl, k )
          var p_TS *tcell.Style = m.Get_Style( pV, l, Vl, k )
          m_console.SetR( pV.Line_2_GL( l ), pV.Char_2_GL( k ), R, p_TS )
        }
      }
    }
    // Write out last line:
    var LAST_LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, NCL )
    if( LAST_LINE_DIFF_TYPE != DT_DELETED ) {
      // Print from beginning of next line to new cursor position:
      NCLL := pfb.LineLen( NCLv ) // Line length
      END_LAST_LINE := Min_i( NCLL, NCP )
      for k:=m.leftChar; k<END_LAST_LINE; k++ {
        R := pfb.GetR( NCLv, k )
        var p_TS *tcell.Style = m.Get_Style( pV, NCL, NCLv, k )
        m_console.SetR( pV.Line_2_GL( NCL ), pV.Char_2_GL( k ), R, p_TS )
      }
    }
  }
}

// Cursor is moving backwards
// Write out from (OCL,OCP) back to but not including (NCL,NCP)
func (m *Diff) GoToCrsPos_WV_Backward( OCL, OCP, NCL, NCP int ) {

  pV  := m_vis.CV()
  pfb := pV.p_fb
  // Convert OCL and NCL, which are diff lines, to view lines:
  OCLv := m.ViewLine( pV, OCL )
  NCLv := m.ViewLine( pV, NCL )

  if( OCL == NCL ) { // Only one line:
    for k:=OCP; NCP<k; k-- {
      R := pfb.GetR( OCLv, k )
      var p_TS *tcell.Style = m.Get_Style( pV, OCL, OCLv, k )
      m_console.SetR( pV.Line_2_GL( OCL ), pV.Char_2_GL( k ), R, p_TS )
    }
  } else { // Multiple lines
    // Write out first line:
    var FIRST_LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, OCL )
    if( FIRST_LINE_DIFF_TYPE != DT_DELETED ) {
      OCLL := pfb.LineLen( OCLv ) // Old cursor line length
      RIGHT_MOST_POS := Min_i( OCP, LLM1(OCLL) )
      for k:=RIGHT_MOST_POS; m.leftChar<k; k-- {
        R := pfb.GetR( OCLv, k )
        var p_TS *tcell.Style = m.Get_Style( pV, OCL, OCLv, k )
        m_console.SetR( pV.Line_2_GL( OCL ), pV.Char_2_GL( k ), R, p_TS )
      }
      if( m.leftChar < OCLL ) {
        R := pfb.GetR( OCLv, m.leftChar )
        var p_TS *tcell.Style = m.Get_Style( pV, OCL, OCLv, m.leftChar )
        m_console.SetR( pV.Line_2_GL( OCL ), pV.Char_2_GL( m.leftChar ), R, p_TS )
      }
    }
    // Write out intermediate lines:
    for l:=OCL-1; NCL<l; l-- {
      var LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, l )
      if( LINE_DIFF_TYPE != DT_DELETED ) {
        // Convert l, which is diff line, to view line:
        Vl := m.ViewLine( pV, l )
        LL := pfb.LineLen( Vl ) // Line length
        END_OF_LINE := Min_i( m.RightChar( pV ), LLM1(LL) )
        for k:=END_OF_LINE; m.leftChar<k; k-- {
          R := pfb.GetR( Vl, k )
          var p_TS *tcell.Style = m.Get_Style( pV, l, Vl, k )
          m_console.SetR( pV.Line_2_GL( l ), pV.Char_2_GL( k ), R, p_TS )
        }
      }
    }
    // Write out last line:
    var LAST_LINE_DIFF_TYPE Diff_Type = m.DiffType( pV, NCL )
    if( LAST_LINE_DIFF_TYPE != DT_DELETED ) {
      // Print from end of last line to new cursor position:
      NCLL := pfb.LineLen( NCLv ) // New cursor line length
      END_LAST_LINE := Min_i( m.RightChar( pV ), LLM1(NCLL) )
      for k:=END_LAST_LINE; NCP<k; k-- {
        R := pfb.GetR( NCLv, k )
        var p_TS *tcell.Style = m.Get_Style( pV, NCL, NCLv, k )
        m_console.SetR( pV.Line_2_GL( NCL ), pV.Char_2_GL( k ), R, p_TS )
      }
      if( NCP < NCLL ) {
        R := pfb.GetR( NCLv, NCP )
        var p_TS *tcell.Style = m.Get_Style( pV, NCL, NCLv, NCP )
        m_console.SetR( pV.Line_2_GL( NCL ), pV.Char_2_GL( NCP ), R, p_TS )
      }
    }
  }
}

