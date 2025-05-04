
package main

import (
//"fmt"
//"github.com/gdamore/tcell/v2"
)

func (m *Vis) InitViewFuncs() {

  m.view_funcs[ ':' ] = Handle_Colon
//m.view_funcs[ '+' ] = Handle_q
  m.view_funcs[ 'i' ] = Handle_i
//m.view_funcs[ 'j' ] = Handle_j
//m.view_funcs[ 'k' ] = Handle_k
  m.view_funcs[ 'l' ] = Handle_l
  m.view_funcs[ 'h' ] = Handle_h

  m.view_funcs[ 'v' ] = Handle_v
  m.view_funcs[ 'V' ] = Handle_V
  m.view_funcs[ 'a' ] = Handle_a
  m.view_funcs[ 'A' ] = Handle_A
  m.view_funcs[ 'o' ] = Handle_o
  m.view_funcs[ 'O' ] = Handle_O
  m.view_funcs[ 'x' ] = Handle_x
  m.view_funcs[ 's' ] = Handle_s
  m.view_funcs[ 'c' ] = Handle_c
  m.view_funcs[ 'C' ] = Handle_C
  m.view_funcs[ 'Q' ] = Handle_Q
  m.view_funcs[ 'H' ] = Handle_H
  m.view_funcs[ 'L' ] = Handle_L
  m.view_funcs[ 'M' ] = Handle_M
  m.view_funcs[ '0' ] = Handle_0
  m.view_funcs[ '$' ] = Handle_Dollar
  m.view_funcs[ '\n'] = Handle_Return
  m.view_funcs[ '\r'] = Handle_Return
  m.view_funcs[ 'G' ] = Handle_G
  m.view_funcs[ 'b' ] = Handle_b
  m.view_funcs[ 'w' ] = Handle_w
  m.view_funcs[ 'e' ] = Handle_e
  m.view_funcs[ 'f' ] = Handle_f
  m.view_funcs[ ';' ] = Handle_SemiColon
  m.view_funcs[ '%' ] = Handle_Percent
  m.view_funcs[ '{' ] = Handle_LeftSquigglyBracket
  m.view_funcs[ '}' ] = Handle_RightSquigglyBracket
  m.view_funcs[ 'F' ] = Handle_F
  m.view_funcs[ 'B' ] = Handle_B
  m.view_funcs[ '/' ] = Handle_Slash
  m.view_funcs[ '*' ] = Handle_Star
  m.view_funcs[ '&' ] = Handle_Ampersand
  m.view_funcs[ '.' ] = Handle_Dot
  m.view_funcs[ 'm' ] = Handle_m
  m.view_funcs[ 'g' ] = Handle_g
  m.view_funcs[ 'W' ] = Handle_W
  m.view_funcs[ 'd' ] = Handle_d
  m.view_funcs[ 'y' ] = Handle_y
  m.view_funcs[ 'D' ] = Handle_D
  m.view_funcs[ 'p' ] = Handle_p
  m.view_funcs[ 'P' ] = Handle_P
  m.view_funcs[ 'r' ] = Handle_r
  m.view_funcs[ 'R' ] = Handle_R
  m.view_funcs[ 'J' ] = Handle_J
  m.view_funcs[ '~' ] = Handle_Tilda
  m.view_funcs[ 'n' ] = Handle_n
  m.view_funcs[ 'N' ] = Handle_N
  m.view_funcs[ 'u' ] = Handle_u
  m.view_funcs[ 'U' ] = Handle_U
  m.view_funcs[ 'z' ] = Handle_z
}

func Handle_Colon( m *Vis ) {

  if( 0 == m.colon_file.NumLines() ) {
    m.colon_file.PushLE();
  }
  var p_cv *FileView = m.CV();
  var NUM_COLS int = p_cv.WinCols();
  var X        int = p_cv.X();
  var Y        int = p_cv.Cmd__Line_Row();

  m.colon_view.SetContext( NUM_COLS, X, Y );
  m.colon_mode = true;

  var CL int = m.colon_view.CrsLine();
  var LL int = m.colon_file.LineLen( CL );

  if( 0<LL ) {
    // Something on current line, so goto command line in escape mode
    m.colon_view.Update();
  } else {
    // Nothing on current line, so goto command line in insert mode
    L_Handle_i( m );
  }
}

func Handle_i( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('i') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_i()
  } else                  { p_cv.Do_i()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_j( m *Vis, num int ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.GoDown( num )
  } else                  { p_cv.GoDown( num )
  }
}

func Handle_k( m *Vis, num int ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.GoUp( num )
  } else                  { p_cv.GoUp( num )
  }
}

func Handle_l( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.GoRight()
  } else                  { p_cv.GoRight()
  }
}

func Handle_h( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.GoLeft()
  } else                  { p_cv.GoLeft()
  }
}

func Handle_v( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.vis_buf.Clear()
    m_key.vis_buf.Push( make_Key_rune('v') )
    m_key.save_2_vis_buf = true
  }
  var p_cv *FileView = m.CV()

  copy_vis_buf_2_dot_buf_n := false
  if( p_cv.in_diff_mode ) { copy_vis_buf_2_dot_buf_n = m.diff.Do_v()
  } else                  { copy_vis_buf_2_dot_buf_n = p_cv.Do_v()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_vis_buf = false;

    if( copy_vis_buf_2_dot_buf_n ) {
      m_key.dot_buf_n.Copy( m_key.vis_buf );
    }
  }
}

func Handle_V( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.vis_buf.Clear()
    m_key.vis_buf.Push( make_Key_rune('V') )
    m_key.save_2_vis_buf = true
  }
  var p_cv *FileView = m.CV()

  copy_vis_buf_2_dot_buf_n := false
  if( p_cv.in_diff_mode ) { copy_vis_buf_2_dot_buf_n = m.diff.Do_V()
  } else                  { copy_vis_buf_2_dot_buf_n = p_cv.Do_V()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_vis_buf = false;

    if( copy_vis_buf_2_dot_buf_n ) {
      m_key.dot_buf_n.Copy( m_key.vis_buf )
    }
  }
}

func Handle_a( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('a') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_a()
  } else                  { p_cv.Do_a()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_A( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('A') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_A()
  } else                  { p_cv.Do_A()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_o( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('o') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_o()
  } else                  { p_cv.Do_o()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_O( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('O') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_O()
  } else                  { p_cv.Do_O()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_x( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('x') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_x()
  } else                  { p_cv.Do_x()
  }
}

func Handle_s( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('s') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_s()
  } else                  { p_cv.Do_s()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_c( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {

    if( kr.R == 'w' ) {
      if( !m_key.get_from_dot_buf_n ) {
        m_key.dot_buf_n.Clear()
        m_key.dot_buf_n.Push( make_Key_rune('c') )
        m_key.dot_buf_n.Push( make_Key_rune('w') )
        m_key.save_2_dot_buf_n = true
      }
      var p_cv *FileView = m.CV();

      if( p_cv.in_diff_mode ) { m.diff.Do_cw()
      } else                  { p_cv.Do_cw()
      }
      if( !m_key.get_from_dot_buf_n ) {
        m_key.save_2_dot_buf_n = false
      }
    } else if( kr.R == '$' ) {
      if( !m_key.get_from_dot_buf_n ) {
        m_key.dot_buf_n.Clear()
        m_key.dot_buf_n.Push( make_Key_rune('c') )
        m_key.dot_buf_n.Push( make_Key_rune('$') )
        m_key.save_2_dot_buf_n = true
      }
      var p_cv *FileView = m.CV();

      if( p_cv.in_diff_mode ) { m.diff.Do_D(); m.diff.Do_a()
      } else                  { p_cv.Do_D();     p_cv.Do_a()
      }
      if( !m_key.get_from_dot_buf_n ) {
        m_key.save_2_dot_buf_n = false
      }
    }
  }
}

func Handle_C( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if       ( kr.R == 'C' ) { m_console.copy_paste_buf_2_system_clipboard()
    } else if( kr.R == 'P' ) { m_console.copy_system_clipboard_2_paste_buf()
    }
  }
}

func Handle_Q( m *Vis ) {

  Handle_Dot( m )
  Handle_j( m, 1 )
  Handle_0( m )
}

func Handle_H( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToTopLineInView()
  } else                  {   p_cv.GoToTopLineInView()
  }
}

func Handle_L( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToBotLineInView()
  } else                  {   p_cv.GoToBotLineInView()
  }
}

func Handle_M( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToMidLineInView()
  } else                  {   p_cv.GoToMidLineInView()
  }
}

func Handle_0( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToBegOfLine()
  } else                  {   p_cv.GoToBegOfLine()
  }
}

func Handle_Dollar( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToEndOfLine()
  } else                  {   p_cv.GoToEndOfLine()
  }
}

func Handle_Return( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToEndOfNextLine()
  } else {
    if( m_SLASH_FILE == m.Curr_FileNum() ) {
      // In search buffer, search for pattern on current line:
      var lp *FLine = p_cv.p_fb.GetLP( p_cv.CrsLine() );

      m.Handle_Slash_GotPattern( lp.to_str() );
    } else {
      // Normal case:
      p_cv.GoToEndOfNextLine();
    }
  }
}

func Handle_G( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToEndOfFile()
  } else                  {   p_cv.GoToEndOfFile()
  }
}

func Handle_b( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToPrevWord()
  } else                  {   p_cv.GoToPrevWord()
  }
}

func Handle_w( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToNextWord()
  } else                  {   p_cv.GoToNextWord()
  }
}

func Handle_e( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToEndOfWord()
  } else                  {   p_cv.GoToEndOfWord()
  }
}

func Handle_f( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    m.fast_rune = kr.R;

    var p_cv *FileView = m.CV();

    if( p_cv.in_diff_mode ) { m.diff.Do_f( m.fast_rune )
    } else                  {   p_cv.Do_f( m.fast_rune )
    }
  }
}

func Handle_SemiColon( m *Vis ) {

  if( 0 <= m.fast_rune ) {

    var p_cv *FileView = m.CV();

    if( p_cv.in_diff_mode ) { m.diff.Do_f( m.fast_rune )
    } else                  {   p_cv.Do_f( m.fast_rune )
    }
  }
}

func Handle_Percent( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToOppositeBracket()
  } else                  {   p_cv.GoToOppositeBracket()
  }
}

func Handle_LeftSquigglyBracket( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToLeftSquigglyBracket()
  } else                  {   p_cv.GoToLeftSquigglyBracket()
  }
}

func Handle_RightSquigglyBracket( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.GoToRightSquigglyBracket()
  } else                  {   p_cv.GoToRightSquigglyBracket()
  }
}

func Handle_F( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.PageDown()
  } else                  {   p_cv.PageDown()
  }
}

func Handle_B( m *Vis ) {

  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.PageUp()
  } else                  {   p_cv.PageUp()
  }
}

func Handle_Slash( m *Vis ) {

  if( 0 == m.slash_file.NumLines() ) {
    m.slash_file.PushLE();
  }
  var p_cv *FileView = m.CV();
  var NUM_COLS int = p_cv.WinCols();
  var X        int = p_cv.X();
  var Y        int = p_cv.Cmd__Line_Row();

  m.slash_view.SetContext( NUM_COLS, X, Y );
  m.slash_mode = true;

  var CL int = m.slash_view.CrsLine();
  var LL int = m.slash_file.LineLen( CL );

  if( 0<LL ) {
    // Something on current line, so goto command line in escape mode
    m.slash_view.Update();
  } else {
    // Nothing on current line, so goto command line in insert mode
    L_Handle_i( m );
  }
}

func Handle_Star( m *Vis ) {

  var p_cv *FileView = m.CV()

  var pattern string
  if( p_cv.in_diff_mode ) { pattern = m.diff.Do_Star_GetNewPattern()
  } else                  { pattern =   p_cv.Do_Star_GetNewPattern()
  }
  if( pattern != m.regex_str ) {

    m.regex_str = pattern;

    if( 0 < len(m.regex_str) ) {
      m.Do_Star_Update_Search_Editor()
    }
    // Show new star pattern for all windows currently displayed:
    m.UpdateViews( true );
  }
}

func Handle_Ampersand( m *Vis ) {
  // FIXME:
}

func Handle_Dot( m *Vis ) {

  if( 0<m_key.dot_buf_n.Len() ) {
    if( m_key.save_2_map_buf ) {
      // Pop '.' off map_buf, because the contents of m_key.map_buf
      // will be saved to m_key.map_buf.
      m_key.map_buf.Pop()
    }
    m_key.get_from_dot_buf_n = true;

    for m_key.get_from_dot_buf_n {
      kr := m_key.In()

      var cf CmdFunc = m.GetCmdFunc( kr )
      if( nil != cf ) { cf(m) }
    }
    var p_cv *FileView = m.CV();

    if( p_cv.in_diff_mode ) {
      // Diff does its own update every time a command is run
    } else {
      // Dont update until after all the commands have been executed:
      p_cv.p_fb.Update()
    }
  }
}

func Handle_m( m *Vis ) {
  // FIXME:
}

func Handle_g( m *Vis ) {

  if( m.running ) {
    kr := m_key.In()

    if( kr.IsKeyRune() ) {
      p_cv := m.CV();

      if( p_cv.in_diff_mode ) {
        if       ( kr.R == 'g' ) { m.diff.GoToTopOfFile()
        } else if( kr.R == '0' ) { m.diff.GoToStartOfRow()
        } else if( kr.R == '$' ) { m.diff.GoToEndOfRow()
        } else if( kr.R == 'f' ) { m.diff.GoToFile()
        }
      } else {
        if       ( kr.R == 'g' ) { p_cv.GoToTopOfFile()
        } else if( kr.R == '0' ) { p_cv.GoToStartOfRow()
        } else if( kr.R == '$' ) { p_cv.GoToEndOfRow()
        } else if( kr.R == 'f' ) { p_cv.GoToFile()
        }
      }
    }
  }
}

func Handle_W( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {

    if       ( kr.R == 'W' ) { m.GoToNextWindow()
    } else if( kr.R == 'l' ) { m.GoToNextWindow_l()
    } else if( kr.R == 'h' ) { m.GoToNextWindow_h()
    } else if( kr.R == 'j' ||
               kr.R == 'k' ) { m.GoToNextWindow_jk()
    } else if( kr.R == 'R' ) { m.FlipWindows()
    }
  }
}

func Handle_d( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if( kr.R == 'd' ) {
      if( !m_key.get_from_dot_buf_n ) {
        m_key.dot_buf_n.Clear()
        m_key.dot_buf_n.Push( make_Key_rune('d') )
        m_key.dot_buf_n.Push( make_Key_rune('d') )
      }
      var p_cv *FileView = m.CV();
      if( p_cv.in_diff_mode ) { m.diff.Do_dd()
      } else                  {   p_cv.Do_dd()
      }
    } else if( kr.R == 'w' ) {
      if( !m_key.get_from_dot_buf_n ) {
        m_key.dot_buf_n.Clear()
        m_key.dot_buf_n.Push( make_Key_rune('d') )
        m_key.dot_buf_n.Push( make_Key_rune('w') )
      }
      var p_cv *FileView = m.CV();
      if( p_cv.in_diff_mode ) { m.diff.Do_dw()
      } else                  {   p_cv.Do_dw()
      }
    }
  }
}

func Handle_y( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {

    var p_cv *FileView = m.CV();
    if( p_cv.in_diff_mode ) {
      if       ( kr.R == 'y' ) { m.diff.Do_yy()
      } else if( kr.R == 'w' ) { m.diff.Do_yw()
      }
    } else {
      if       ( kr.R == 'y' ) { p_cv.Do_yy()
      } else if( kr.R == 'w' ) { p_cv.Do_yw()
      }
    }
  }
}

func Handle_D( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('D') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_D()
  } else                  {   p_cv.Do_D()
  }
}

func Handle_p( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('p') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_p()
  } else                  {   p_cv.Do_p()
  }
}

func Handle_P( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('P') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_P()
  } else                  {   p_cv.Do_P()
  }
}

func Handle_r( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('r') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_r()
  } else                  {   p_cv.Do_r()
  }
}

func Handle_R( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('R') )
    m_key.save_2_dot_buf_n = true
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_R()
  } else                  {   p_cv.Do_R()
  }
  if( !m_key.get_from_dot_buf_n ) {
    m_key.save_2_dot_buf_n = false
  }
}

func Handle_J( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('J') )
  }
  var p_cv *FileView = m.CV();

  if( p_cv.in_diff_mode ) { m.diff.Do_J()
  } else                  {   p_cv.Do_J()
  }
}

func Handle_Tilda( m *Vis ) {

  if( !m_key.get_from_dot_buf_n ) {
    m_key.dot_buf_n.Clear()
    m_key.dot_buf_n.Push( make_Key_rune('~') )
  }
  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_Tilda()
  } else                  {   p_cv.Do_Tilda()
  }
}

func Handle_n( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_n()
  } else                  {   p_cv.Do_n()
  }
}

func Handle_N( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_N()
  } else                  {   p_cv.Do_N()
  }
}

func Handle_u( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_u()
  } else                  {   p_cv.Do_u()
  }
}

func Handle_U( m *Vis ) {

  var p_cv *FileView = m.CV()

  if( p_cv.in_diff_mode ) { m.diff.Do_U()
  } else                  {   p_cv.Do_U()
  }
}

func Handle_z( m *Vis ) {

  if( m.running ) {
    kr := m_key.In()

    is_eol  := kr.IsEndOfLineDelim()
    is_rune := kr.IsKeyRune()

    var p_cv *FileView = m.CV()

    if( p_cv.in_diff_mode ) {
      if( is_eol ||
                 (is_rune && kr.R == 't') ) { m.diff.MoveCurrLineToTop()
      } else if( (is_rune && kr.R == 'z') ) { m.diff.MoveCurrLineCenter()
      } else if( (is_rune && kr.R == 'b') ) { m.diff.MoveCurrLineToBottom()
      }
    } else {
      if( is_eol ||
                 (is_rune && kr.R == 't') ) { p_cv.MoveCurrLineToTop()
      } else if( (is_rune && kr.R == 'z') ) { p_cv.MoveCurrLineCenter()
      } else if( (is_rune && kr.R == 'b') ) { p_cv.MoveCurrLineToBottom()
      }
    }
  }
}

