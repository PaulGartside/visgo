
package main

import (
  "bytes"
  "fmt"
//"io/ioutil"
//"log"
  "os"
  "strings"
//"time"

//"github.com/gdamore/tcell/v2"
)

var m_prog_name string
var m_wd string // vis Working directory
//var m_win_drive string
//var m_win_drive_and_path_sep string
var m_WIN_D        string // Windows drive
var m_WIN_D_len    int    // Windows drive length
var m_WIN_D_PS     string // Windows drive and path separator
var m_WIN_D_PS_len int    // Windows drive and path separator length

var m_BE_FILE    int = 0;    // Buffer editor file
var m_HELP_FILE  int = 1;    // Help          file
var m_MSG_FILE   int = 2;    // Message       file
var m_SHELL_FILE int = 3;    // Command Shell file
var m_COLON_FILE int = 4;    // Colon command file
var m_SLASH_FILE int = 5;    // Slash command file
var m_USER_FILE  int = 6;    // First user file

var  m_EDIT_BUF_NAME string = "BUFFER_EDITOR"
var  m_HELP_BUF_NAME string = "VIS_HELP"
var  m_MSG__BUF_NAME string = "MESSAGE_BUFFER"
var m_SHELL_BUF_NAME string = "SHELL_BUFFER"
var m_COLON_BUF_NAME string = "COLON_BUFFER"
var m_SLASH_BUF_NAME string = "SLASH_BUFFER"

var m_bb bytes.Buffer  // bytes Buffer

const MAX_WINS = 8
const DIR_DELIM = os.PathSeparator

var DIR_DELIM_S string = string(DIR_DELIM)

// This is how enums are created in go:
type Paste_Mode int

const (
  PM_LINE Paste_Mode = iota // Whole line from start line to finish line
  PM_ST_FN                  // Visual mode from start to finish
  PM_BLOCK                  // Visual mode rectangular block spanning one or more lines
)

type Paste_Pos int

const (
  PP_Before = iota
  PP_After
)

type File_Type int

const (
  FT_UNKNOWN File_Type = iota
  FT_BUFFER_EDITOR
  FT_DIR
  FT_BASH
  FT_CPP
  FT_GO
  FT_IDL
  FT_SQL
  FT_TEXT
  FT_XML
)

type ChangeType int

const (
  CT_INSERT_LINE = iota
  CT_REMOVE_LINE
  CT_INSERT_TEXT
  CT_REMOVE_TEXT
  CT_REPLACE_TEXT
)

type Tile_Pos int

const (
  TP_NONE Tile_Pos = iota
  // 1 x 1 tile:
  TP_FULL
  // 1 x 2 tiles:
  TP_LEFT_HALF
  TP_RITE_HALF
  // 2 x 1 tiles:
  TP_TOP__HALF
  TP_BOT__HALF
  // 2 x 2 tiles:
  TP_TOP__LEFT_QTR
  TP_TOP__RITE_QTR
  TP_BOT__LEFT_QTR
  TP_BOT__RITE_QTR
  // 1 x 4 tiles:
  TP_LEFT_QTR
  TP_RITE_QTR
  TP_LEFT_CTR__QTR
  TP_RITE_CTR__QTR
  // 2 x 4 tiles:
  TP_TOP__LEFT_8TH
  TP_TOP__RITE_8TH
  TP_TOP__LEFT_CTR_8TH
  TP_TOP__RITE_CTR_8TH
  TP_BOT__LEFT_8TH
  TP_BOT__RITE_8TH
  TP_BOT__LEFT_CTR_8TH
  TP_BOT__RITE_CTR_8TH
  // 1 x 3 tiles:
  TP_LEFT_THIRD
  TP_CTR__THIRD
  TP_RITE_THIRD
  TP_LEFT_TWO_THIRDS
  TP_RITE_TWO_THIRDS
)

var m_vis     Vis
var m_console Console
var m_key     Key

var m_rbuf RLine
var m_log []string

// Return filename from path.
// Ex: a/b/fname.txt -> fname.txt
// EX: fname.txt -> fname.txt
//
func cut_file_name( path string ) string {
  slash_idx := strings.LastIndex( path, string(DIR_DELIM) )
  return path[slash_idx+1:]
}

func die( msg string ) {
  fmt.Println( m_prog_name +" : "+ msg )
  os.Exit( 1 )
}

func print_log() {
  // Log()
  if( 0 < len(m_log) ) {
    fmt.Println()
    for _, msg := range m_log {
      fmt.Println( msg )
    }
    fmt.Println()
  }
}

func get_win_drive() {

  wd, err := os.Getwd()
  if( err != nil ) {
    die( fmt.Sprintf("os.Getwd() : %v", err) )
  }
  m_wd = wd
  for _,r := range m_wd {
    if( r == os.PathSeparator ) {
      break
    } else {
      m_WIN_D += string(r)
    }
  }
  m_WIN_D_PS = m_WIN_D + string(os.PathSeparator)

  m_WIN_D_len    = len( m_WIN_D )
  m_WIN_D_PS_len = len( m_WIN_D_PS )
}

func init_run() {
  m_vis.Init()
  m_vis.Run()
}

func main() {
  m_prog_name = cut_file_name( os.Args[0] ) // Trim path down to filename

  get_win_drive()
  m_console.Init()

  defer func() {
    panicked := false
    var r = recover()
    if r != nil {
      panicked = true
      fmt.Println(r)
    }
    m_console.Cleanup()
    print_log()

    if( panicked ) {
      panic( r )
      os.Exit(1)
    }
    os.Exit(0)
  }()
  init_run()
}

