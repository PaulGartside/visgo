// go:build windows

package main

import (
  "fmt"
  "strings"
  "golang.design/x/clipboard"
)

var initialized_clipboard bool

func Init_clipboard() bool {

  if( !initialized_clipboard ) {

    var e error = clipboard.Init()
    if( e != nil ) {
      msg := fmt.Sprintf("clipboard.Init() Error: %v", e )
      m_vis.CmdLineMessage( msg )
    } else {
      initialized_clipboard = true
    }
  }
  return initialized_clipboard
}

func Copy_paste_buf_2_system_clipboard() {

  if( Init_clipboard() ) {
    var S string

    reg_len := m_vis.reg.Len()

    for k:=0; k<reg_len; k++ {
      p_rl := m_vis.reg.GetLP( k )
      S += p_rl.to_str()
      if( k < reg_len-1 ) { S += "\n"
      }
    }
    clipboard.Write( clipboard.FmtText, []byte( S ) )

    S_len := len( S )

    if( 0 < S_len ) {
      msg := fmt.Sprintf("Copied %v lines, %v bytes to system clipboard",
                         reg_len, S_len )
      m_vis.CmdLineMessage( msg )
    } else {
      m_vis.CmdLineMessage("Cleared system clipboard");
    }
  }
}

func Copy_system_clipboard_2_paste_buf() {

  if( Init_clipboard() ) {
    m_vis.reg.Clear()

    cb_str := string( clipboard.Read( clipboard.FmtText ) )
    cb_str = strings.ReplaceAll( cb_str, "\r\n", "\n" )
    cb_str_len := len( cb_str )

    if( cb_str_len == 0 ) {
      m_vis.CmdLineMessage("Cleared paste buffer");
    } else {
      var cb_lines []string = strings.Split( cb_str, "\n" )
      cb_lines_len := len( cb_lines )
      for _, cb_line := range cb_lines {
        p_rl := new( RLine )
        p_rl.from_str( cb_line )
        m_vis.reg.PushLP( p_rl )
      }
      msg := fmt.Sprintf("Copied %v lines, %v bytes to paste buffer",
                         cb_lines_len, cb_str_len )
      m_vis.CmdLineMessage( msg )
    }
  }
}

