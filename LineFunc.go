
package main

import (
//"fmt"
  "github.com/gdamore/tcell/v2"
)

func (m *Vis) InitLineFuncs() {

  m.line_funcs[':' ]         = L_Handle_Escape;
  m.line_funcs[tcell.KeyESC] = L_Handle_Escape;
  m.line_funcs[tcell.KeyLF ] = L_Handle_Return; // 10, Line feed
  m.line_funcs[tcell.KeyCR ] = L_Handle_Return; // 13, Carriage return, used on Windows

  m.line_funcs[ 'a' ] = L_Handle_a;
  m.line_funcs[ 'A' ] = L_Handle_A;
  m.line_funcs[ 'b' ] = L_Handle_b;
  m.line_funcs[ 'c' ] = L_Handle_c;
  m.line_funcs[ 'd' ] = L_Handle_d;
  m.line_funcs[ 'D' ] = L_Handle_D;
  m.line_funcs[ 'e' ] = L_Handle_e;
  m.line_funcs[ 'f' ] = L_Handle_f;
  m.line_funcs[ 'g' ] = L_Handle_g;
  m.line_funcs[ 'G' ] = L_Handle_G;
  m.line_funcs[ 'h' ] = L_Handle_h;
  m.line_funcs[ 'i' ] = L_Handle_i;
  m.line_funcs[ 'j' ] = L_Handle_j;
  m.line_funcs[ 'J' ] = L_Handle_J;
  m.line_funcs[ 'k' ] = L_Handle_k;
  m.line_funcs[ 'l' ] = L_Handle_l;
  m.line_funcs[ 'o' ] = L_Handle_o;
  m.line_funcs[ 'n' ] = L_Handle_n;
  m.line_funcs[ 'N' ] = L_Handle_N;
  m.line_funcs[ 'p' ] = L_Handle_p;
  m.line_funcs[ 'P' ] = L_Handle_P;
  m.line_funcs[ 'R' ] = L_Handle_R;
  m.line_funcs[ 's' ] = L_Handle_s;
  m.line_funcs[ 'v' ] = L_Handle_v;
  m.line_funcs[ 'w' ] = L_Handle_w;
  m.line_funcs[ 'x' ] = L_Handle_x;
  m.line_funcs[ 'y' ] = L_Handle_y;
  m.line_funcs[ '~' ] = L_Handle_Tilda;
  m.line_funcs[ '$' ] = L_Handle_Dollar;
  m.line_funcs[ '%' ] = L_Handle_Percent;
  m.line_funcs[ '0' ] = L_Handle_0;
  m.line_funcs[ ';' ] = L_Handle_SemiColon;
  m.line_funcs[ ':' ] = L_Handle_Colon;
  m.line_funcs[ '/' ] = L_Handle_Slash;
  m.line_funcs[ '.' ] = L_Handle_Dot;
}

func L_Handle_Escape( m *Vis ) {

  if       ( m.colon_mode ) { L_Handle_Colon( m )
  } else if( m.slash_mode ) { L_Handle_Slash( m )
  }
}

func L_Handle_Return( m *Vis ) {
//Log("Top: L_Handle_Return()")

  if( m.colon_mode ) {
    m.colon_mode = false;
    m.colon_view.HandleReturn();
    m.Handle_Colon_Cmd();

  } else if( m.slash_mode ) {
    m.slash_mode = false;
    m.slash_view.HandleReturn();
    m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
  }
//Log("Bot: L_Handle_Return()")
}

func L_Handle_a( m *Vis ) {

  if( m.colon_mode ) {
    end_of_line_delim := m.colon_view.Do_a();

    if( end_of_line_delim ) {
      m.colon_mode = false;
      m.Handle_Colon_Cmd();
    }
  } else if( m.slash_mode ) {
    end_of_line_delim := m.slash_view.Do_a();

    if( end_of_line_delim ) {
      m.slash_mode = false;
      m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
    }
  }
}

func L_Handle_A( m *Vis ) {

  if( m.colon_mode ) {
    end_of_line_delim := m.colon_view.Do_A();

    if( end_of_line_delim ) {
      m.colon_mode = false;
      m.Handle_Colon_Cmd();
    }
  } else if( m.slash_mode ) {
    end_of_line_delim := m.slash_view.Do_A();

    if( end_of_line_delim ) {
      m.slash_mode = false;
      m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
    }
  }
}

func L_Handle_b( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToPrevWord();
  } else if( m.slash_mode ) { m.slash_view.GoToPrevWord();
  }
}

func L_Handle_c( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if( kr.R == 'w' ) {
      if       ( m.colon_mode ) { m.colon_view.Do_cw();
      } else if( m.slash_mode ) { m.slash_view.Do_cw();
      }
    } else if( kr.R == '$' ) {
      if( m.colon_mode ) {
        m.colon_view.Do_D();
        m.colon_view.Do_a();
      } else if( m.slash_mode ) {
        m.slash_view.Do_D();
        m.slash_view.Do_a();
      }
    }
  }
}

func L_Handle_d( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if( kr.R == 'd' ) {
      if       ( m.colon_mode ) { m.colon_view.Do_dd();
      } else if( m.slash_mode ) { m.slash_view.Do_dd();
      }
    } else if( kr.R == 'w' ) {
      if       ( m.colon_mode ) { m.colon_view.Do_dw();
      } else if( m.slash_mode ) { m.slash_view.Do_dw();
      }
    }
  }
}

func L_Handle_D( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_D();
  } else if( m.slash_mode ) { m.slash_view.Do_D();
  }
}

func L_Handle_e( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToEndOfWord();
  } else if( m.slash_mode ) { m.slash_view.GoToEndOfWord();
  }
}

func L_Handle_f( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    m.fast_rune = kr.R
    if       ( m.colon_mode ) { m.colon_view.Do_f( m.fast_rune );
    } else if( m.slash_mode ) { m.slash_view.Do_f( m.fast_rune );
    }
  }
}

func L_Handle_g( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if( kr.R == 'g' ) {
      if       ( m.colon_mode ) { m.colon_view.GoToTopOfFile();
      } else if( m.slash_mode ) { m.slash_view.GoToTopOfFile();
      }
    } else if( kr.R == '0' ) {
      if       ( m.colon_mode ) { m.colon_view.GoToStartOfRow();
      } else if( m.slash_mode ) { m.slash_view.GoToStartOfRow();
      }
    } else if( kr.R == '$' ) {
      if       ( m.colon_mode ) { m.colon_view.GoToEndOfRow();
      } else if( m.slash_mode ) { m.slash_view.GoToEndOfRow();
      }
    }
  }
}

func L_Handle_G( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToEndOfFile();
  } else if( m.slash_mode ) { m.slash_view.GoToEndOfFile();
  }
}

func L_Handle_h( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoLeft();
  } else if( m.slash_mode ) { m.slash_view.GoLeft();
  }
}

func L_Handle_i( m *Vis ) {
//Log("Top: L_Handle_i()")

  if( m.colon_mode ) {
    var end_of_line_delim bool = m.colon_view.Do_i();

    if( end_of_line_delim ) {
      m.colon_mode = false;
      m.Handle_Colon_Cmd();
    }
  } else if( m.slash_mode ) {
    var end_of_line_delim bool = m.slash_view.Do_i();

    if( end_of_line_delim ) {
      m.slash_mode = false;

//Log( fmt.Sprintf("m_rbuf.to_str()=%v", m_rbuf.to_str()) )
      m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
    }
  }
//Log("Bot: L_Handle_i()")
}

func L_Handle_j( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoDown();
  } else if( m.slash_mode ) { m.slash_view.GoDown();
  }
}

func L_Handle_J( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_J();
  } else if( m.slash_mode ) { m.slash_view.Do_J();
  }
}

func L_Handle_k( m *Vis ) {
//Log("Top: L_Handle_k()")
  if       ( m.colon_mode ) { m.colon_view.GoUp();
  } else if( m.slash_mode ) {
//Log("m.slash_view.GoUp()")
    m.slash_view.GoUp();
  }
//Log("Bot: L_Handle_k()")
}

func L_Handle_l( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoRight();
  } else if( m.slash_mode ) { m.slash_view.GoRight();
  }
}

func L_Handle_o( m *Vis ) {

  if( m.colon_mode ) {
    var end_of_line_delim bool = m.colon_view.Do_o();

    if( end_of_line_delim ) {
      m.colon_mode = false;

      m.Handle_Colon_Cmd();
    }
  } else if( m.slash_mode ) {
    var end_of_line_delim bool = m.slash_view.Do_o();

    if( end_of_line_delim ) {
      m.slash_mode = false;

      m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
    }
  }
}

func L_Handle_n( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_n();
  } else if( m.slash_mode ) { m.slash_view.Do_n();
  }
}

func L_Handle_N( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_N();
  } else if( m.slash_mode ) { m.slash_view.Do_N();
  }
}

func L_Handle_p( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_p();
  } else if( m.slash_mode ) { m.slash_view.Do_p();
  }
}

func L_Handle_P( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_P();
  } else if( m.slash_mode ) { m.slash_view.Do_P();
  }
}

func L_Handle_R( m *Vis ) {

  if( m.colon_mode ) {
    var end_of_line_delim bool = m.colon_view.Do_R();

    if( end_of_line_delim ) {
      m.colon_mode = false;

      m.Handle_Colon_Cmd();
    }
  } else if( m.slash_mode ) {
    var end_of_line_delim bool = m.slash_view.Do_R();

    if( end_of_line_delim ) {
      m.slash_mode = false;

      m.Handle_Slash_GotPattern( m_rbuf.to_str(), true );
    }
  }
}

func L_Handle_s( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_s();
  } else if( m.slash_mode ) { m.slash_view.Do_s();
  }
}

func L_Handle_v( m *Vis ) {

  copy_vis_buf_2_dot_buf_l := false;

  if       ( m.colon_mode ) { copy_vis_buf_2_dot_buf_l = m.colon_view.Do_v();
  } else if( m.slash_mode ) { copy_vis_buf_2_dot_buf_l = m.slash_view.Do_v();
  }

  if( copy_vis_buf_2_dot_buf_l ) {}
}

func L_Handle_w( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToNextWord();
  } else if( m.slash_mode ) { m.slash_view.GoToNextWord();
  }
}

func L_Handle_x( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_x();
  } else if( m.slash_mode ) { m.slash_view.Do_x();
  }
}

func L_Handle_y( m *Vis ) {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    if( kr.R == 'y' ) {
      if       ( m.colon_mode ) { m.colon_view.Do_yy();
      } else if( m.slash_mode ) { m.slash_view.Do_yy();
      }
    } else if( kr.R == 'w' ) {
      if       ( m.colon_mode ) { m.colon_view.Do_yw();
      } else if( m.slash_mode ) { m.slash_view.Do_yw();
      }
    }
  }
}

func L_Handle_Tilda( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.Do_Tilda();
  } else if( m.slash_mode ) { m.slash_view.Do_Tilda();
  }
}

func L_Handle_Dollar( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToEndOfLine();
  } else if( m.slash_mode ) { m.slash_view.GoToEndOfLine();
  }
}

func L_Handle_Percent( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToOppositeBracket();
  } else if( m.slash_mode ) { m.slash_view.GoToOppositeBracket();
  }
}

func L_Handle_0( m *Vis ) {

  if       ( m.colon_mode ) { m.colon_view.GoToBegOfLine();
  } else if( m.slash_mode ) { m.slash_view.GoToBegOfLine();
  }
}

func L_Handle_SemiColon( m *Vis ) {

  if( 0 <= m.fast_rune ) {
    if       ( m.colon_mode ) { m.colon_view.Do_f( m.fast_rune );
    } else if( m.slash_mode ) { m.slash_view.Do_f( m.fast_rune );
    }
  }
}

func L_Handle_Colon( m *Vis ) {

  m.colon_mode = false;

  p_cv := m.CV()

  if( p_cv.in_diff_mode ) { m.diff.PrintCursor( p_cv );
  } else                  { p_cv.PrintCursor();
  }
}

func L_Handle_Slash( m *Vis ) {

  m.slash_mode = false;

  p_cv := m.CV()

  if( p_cv.in_diff_mode ) { m.diff.PrintCursor( p_cv );
  } else                  { p_cv.PrintCursor();
  }
}

func L_Handle_Dot( m *Vis ) {
}

