
package main

import (
  "bytes"
  "fmt"
  "log"
//"os"
//"strings"
//"time"

  "github.com/gdamore/tcell/v2"
)

// TS = tcell.Style
var TS_NORMAL    tcell.Style
var TS_BORDER    tcell.Style
var TS_BORDER_HI tcell.Style
var TS_EOF       tcell.Style
var TS_BANNER    tcell.Style
var TS_CONST     tcell.Style
var TS_DEFINE    tcell.Style
var TS_CONTROL   tcell.Style
var TS_EMPTY     tcell.Style
var TS_VISUAL    tcell.Style
var TS_STAR      tcell.Style
var TS_STAR_IN_F tcell.Style
var TS_COMMENT   tcell.Style
var TS_VARTYPE   tcell.Style
var TS_NONASCII  tcell.Style

var TS_RV_NORMAL    tcell.Style
var TS_RV_STAR      tcell.Style
var TS_RV_STAR_IN_F tcell.Style
var TS_RV_DEFINE    tcell.Style
var TS_RV_COMMENT   tcell.Style
var TS_RV_CONST     tcell.Style
var TS_RV_CONTROL   tcell.Style
var TS_RV_VARTYPE   tcell.Style
var TS_RV_NONASCII  tcell.Style

// Windows is limited to the following colors:
//
// ColorBlack,
// ColorMaroon,
// ColorGreen,
// ColorNavy,
// ColorOlive,
// ColorPurple,
// ColorTeal,
// ColorSilver,
// ColorGray,
// ColorRed,
// ColorLime,
// ColorBlue,
// ColorYellow,
// ColorFuchsia,
// ColorAqua,
// ColorWhite,

func init() {
  m_console.Set_Color_Scheme_1()
}

type Console struct {
  screen   tcell.Screen
  running  bool
}

func (m *Console) Init() {
  var e error
  if m.screen, e = tcell.NewScreen(); e != nil { log.Fatal( e ) }

  if e = m.screen.Init(); e != nil { log.Fatal( e ) }

  m.screen.SetCursorStyle( tcell.CursorStyleSteadyBlock )

//if e = clipboard.Init(); e != nil { log.Fatal( e ) }

  m.running = true
}

func (m *Console) Cleanup() {
  m.running = false

  m.screen.Fini()
}

//func (m *Console) Size() (int, int) {
//func (m *Console) Get_Cols_Rows() (int, int) {
//  return m.screen.Size()
//}

func (m *Console) Num_Cols() int {

  n_cols, _ := m.screen.Size()

  return n_cols
}

func (m *Console) Num_Rows() int {

  _, n_rows := m.screen.Size()

  return n_rows
}

func (m *Console) ShowCursor( row, col int ) {
  m.screen.ShowCursor( col, row )
}

func (m *Console) Show() {
  m.screen.Show()
}

// Set a rune on screen
func (m *Console) SetR( row,col int, ru rune, p_S *tcell.Style ) {
  m.screen.SetContent( col, row, ru, nil, *p_S )
}

// Set a slice of bytes on screen
//func (m *Console) SetSB( row,col int, s_b []byte, p_S *tcell.Style ) {
//  s_b_len := len( s_b )
//  for k:=0; k<s_b_len; k++ {
//    m.screen.SetContent( col+k, row, rune(s_b[k]), nil, *p_S )
//  }
//}

// Set a slice of runes on screen
func (m *Console) SetSR( row,col int, s_r []rune, p_S *tcell.Style ) {
  s_r_len := len( s_r )
  for k:=0; k<s_r_len; k++ {
    m.screen.SetContent( col+k, row, s_r[k], nil, *p_S )
  }
}

func (m *Console) SetBuffer( row,col int, buf *bytes.Buffer, p_S *tcell.Style ) {
  done := false
  for k:=0; !done; k++ {
    R, _, err := buf.ReadRune()
    if( err == nil ) {
      m.SetR( row, col+k, R, p_S )
    } else {
      done = true
    }
  }
}

func (m *Console) SetString( row,col int, msg string, p_S *tcell.Style ) {
  var buf bytes.Buffer
  _, err := buf.WriteString( msg )
  if( err != nil ) {
    Log( fmt.Sprintf("Failed to write %s to bytes.Buffer", msg) )
  } else {
    m.SetBuffer( row, col, &buf, p_S )
  }
}

// EventError
// EventTime
// EventInterrupt
// EventKey
// EventMouse
// EventPaste
// EventResize
//
func (m *Console) Key_In() Key_rune {
  var K tcell.Key
  var R rune

  var rcvd_key bool = false

  for m.running && !rcvd_key {
    switch ev := m.screen.PollEvent().(type) {

    case nil:
      m.running = false

    case *tcell.EventResize:
      m_vis.Handle_Resize()
      m.screen.Sync()

    case *tcell.EventKey:
      rcvd_key = true
      K = ev.Key()
      if K == tcell.KeyRune {
        R = ev.Rune()
      }

    default:
    }
  }
  return Key_rune{ K, R }
}

func (m *Console) HasPendingEvent() bool {

  return m.screen.HasPendingEvent()
}

func (m *Console) Set_Color_Scheme_1() {

  TS_NORMAL    = tcell.StyleDefault
  TS_BORDER    = TS_BORDER   .Background( tcell.ColorBlue  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_BORDER_HI = TS_BORDER_HI.Background( tcell.ColorAqua  ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_EOF       = TS_EOF      .Background( tcell.ColorGray  ).Foreground( tcell.ColorRed ).Bold(true)
  TS_BANNER    = TS_BANNER   .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_CONST     = TS_CONST    .Background( tcell.ColorBlack ).Foreground( tcell.ColorAqua ).Bold(true)
  TS_DEFINE    = TS_DEFINE   .Background( tcell.ColorBlack ).Foreground( tcell.ColorMaroon ).Bold(true)
  TS_CONTROL   = TS_CONTROL  .Background( tcell.ColorBlack ).Foreground( tcell.ColorYellow ).Bold(true)
  TS_EMPTY     = tcell.StyleDefault
  TS_VISUAL    = TS_VISUAL   .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)

  TS_STAR      = TS_STAR     .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_STAR_IN_F = TS_STAR_IN_F.Background( tcell.ColorBlue  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_COMMENT   = TS_COMMENT  .Background( tcell.ColorBlack ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_VARTYPE   = TS_VARTYPE  .Background( tcell.ColorBlack ).Foreground( tcell.ColorLime ).Bold(true)
  TS_NONASCII  = TS_NONASCII .Background( tcell.ColorPurple).Foreground( tcell.ColorAqua ).Bold(true)

  TS_RV_NORMAL    = TS_RV_NORMAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_RV_STAR      = TS_RV_STAR     .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_RV_STAR_IN_F = TS_RV_STAR_IN_F.Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_RV_DEFINE    = TS_RV_DEFINE   .Background( tcell.ColorWhite ).Foreground( tcell.ColorPurple ).Bold(true)
  TS_RV_COMMENT   = TS_RV_COMMENT  .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_RV_CONST     = TS_RV_CONST    .Background( tcell.ColorWhite ).Foreground( tcell.ColorAqua ).Bold(true)
  TS_RV_CONTROL   = TS_RV_CONTROL  .Background( tcell.ColorWhite ).Foreground( tcell.ColorFuchsia ).Bold(true)
  TS_RV_VARTYPE   = TS_RV_VARTYPE  .Background( tcell.ColorWhite ).Foreground( tcell.ColorLime  ).Bold(true)
  TS_RV_NONASCII  = TS_RV_NONASCII .Background( tcell.ColorRed  ).Foreground( tcell.ColorBlue ).Bold(true)
}

func (m *Console) Set_Color_Scheme_2() {

  TS_NORMAL    = tcell.StyleDefault
  TS_BORDER    = TS_BORDER   .Background( tcell.ColorBlue  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_BORDER_HI = TS_BORDER_HI.Background( tcell.ColorAqua  ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_EOF       = TS_EOF      .Background( tcell.ColorGray  ).Foreground( tcell.ColorRed ).Bold(true)
  TS_BANNER    = TS_BANNER   .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_CONST     = TS_CONST    .Background( tcell.ColorBlack ).Foreground( tcell.ColorAqua ).Bold(true)
  TS_DEFINE    = TS_DEFINE   .Background( tcell.ColorBlack ).Foreground( tcell.ColorMaroon ).Bold(true)
  TS_CONTROL   = TS_CONTROL  .Background( tcell.ColorBlack ).Foreground( tcell.ColorYellow ).Bold(true)
  TS_EMPTY     = TS_EOF
  TS_VISUAL    = TS_VISUAL   .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)

  TS_STAR      = TS_STAR     .Background( tcell.ColorRed   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_STAR_IN_F = TS_STAR_IN_F.Background( tcell.ColorBlue  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_COMMENT   = TS_COMMENT  .Background( tcell.ColorBlack ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_VARTYPE   = TS_VARTYPE  .Background( tcell.ColorBlack ).Foreground( tcell.ColorLime ).Bold(true)
  TS_NONASCII  = TS_NONASCII .Background( tcell.ColorPurple).Foreground( tcell.ColorAqua ).Bold(true)

  TS_RV_NORMAL    = TS_RV_NORMAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_RV_STAR      = TS_RV_STAR     .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_RV_STAR_IN_F = TS_RV_STAR_IN_F.Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_RV_DEFINE    = TS_RV_DEFINE   .Background( tcell.ColorWhite ).Foreground( tcell.ColorPurple ).Bold(true)
  TS_RV_COMMENT   = TS_RV_COMMENT  .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue ).Bold(true)
  TS_RV_CONST     = TS_RV_CONST    .Background( tcell.ColorWhite ).Foreground( tcell.ColorAqua ).Bold(true)
  TS_RV_CONTROL   = TS_RV_CONTROL  .Background( tcell.ColorWhite ).Foreground( tcell.ColorFuchsia ).Bold(true)
  TS_RV_VARTYPE   = TS_RV_VARTYPE  .Background( tcell.ColorWhite ).Foreground( tcell.ColorLime  ).Bold(true)
  TS_RV_NONASCII  = TS_RV_NONASCII .Background( tcell.ColorRed  ).Foreground( tcell.ColorBlue ).Bold(true)
}

func (m *Console) Set_Color_Scheme_3() {

  TS_NORMAL    = TS_NORMAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlack ).Bold(true)

  TS_BORDER    = TS_BORDER   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_BORDER_HI = TS_BORDER_HI.Background( tcell.ColorBlue  ).Foreground( tcell.ColorAqua  ).Bold(true)
  TS_EOF       = TS_EOF      .Background( tcell.ColorRed   ).Foreground( tcell.ColorGray  ).Bold(true)
  TS_BANNER    = TS_BANNER   .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_CONST     = TS_CONST    .Background( tcell.ColorAqua  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_DEFINE    = TS_DEFINE   .Background( tcell.ColorMaroon).Foreground( tcell.ColorBlack ).Bold(true)
  TS_CONTROL   = TS_CONTROL  .Background( tcell.ColorYellow).Foreground( tcell.ColorBlack ).Bold(true)
  TS_EMPTY     = TS_NORMAL
  TS_VISUAL    = TS_VISUAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed ).Bold(true)

  TS_STAR      = TS_STAR     .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_STAR_IN_F = TS_STAR_IN_F.Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_COMMENT   = TS_COMMENT  .Background( tcell.ColorBlue  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_VARTYPE   = TS_VARTYPE  .Background( tcell.ColorLime  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_NONASCII  = TS_NONASCII .Background( tcell.ColorAqua  ).Foreground( tcell.ColorPurple).Bold(true)

  TS_RV_NORMAL    = TS_RV_NORMAL   .Background( tcell.ColorBlack  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_STAR      = TS_RV_STAR     .Background( tcell.ColorRed    ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_STAR_IN_F = TS_RV_STAR_IN_F.Background( tcell.ColorBlue   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_DEFINE    = TS_RV_DEFINE   .Background( tcell.ColorPurple ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_COMMENT   = TS_RV_COMMENT  .Background( tcell.ColorBlue   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_CONST     = TS_RV_CONST    .Background( tcell.ColorAqua   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_CONTROL   = TS_RV_CONTROL  .Background( tcell.ColorFuchsia).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_VARTYPE   = TS_RV_VARTYPE  .Background( tcell.ColorLime   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_NONASCII  = TS_RV_NONASCII .Background( tcell.ColorBlue   ).Foreground( tcell.ColorRed   ).Bold(true)
}

func (m *Console) Set_Color_Scheme_4() {

  TS_NORMAL    = TS_NORMAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlack ).Bold(true)

  TS_BORDER    = TS_BORDER   .Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_BORDER_HI = TS_BORDER_HI.Background( tcell.ColorBlue  ).Foreground( tcell.ColorAqua  ).Bold(true)
  TS_EOF       = TS_EOF      .Background( tcell.ColorRed   ).Foreground( tcell.ColorGray  ).Bold(true)
  TS_BANNER    = TS_BANNER   .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_CONST     = TS_CONST    .Background( tcell.ColorAqua  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_DEFINE    = TS_DEFINE   .Background( tcell.ColorMaroon).Foreground( tcell.ColorBlack ).Bold(true)
  TS_CONTROL   = TS_CONTROL  .Background( tcell.ColorYellow).Foreground( tcell.ColorBlack ).Bold(true)
  TS_EMPTY     = TS_EOF
  TS_VISUAL    = TS_VISUAL   .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed ).Bold(true)

  TS_STAR      = TS_STAR     .Background( tcell.ColorWhite ).Foreground( tcell.ColorRed   ).Bold(true)
  TS_STAR_IN_F = TS_STAR_IN_F.Background( tcell.ColorWhite ).Foreground( tcell.ColorBlue  ).Bold(true)
  TS_COMMENT   = TS_COMMENT  .Background( tcell.ColorBlue  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_VARTYPE   = TS_VARTYPE  .Background( tcell.ColorLime  ).Foreground( tcell.ColorBlack ).Bold(true)
  TS_NONASCII  = TS_NONASCII .Background( tcell.ColorAqua  ).Foreground( tcell.ColorPurple).Bold(true)

  TS_RV_NORMAL    = TS_RV_NORMAL   .Background( tcell.ColorBlack  ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_STAR      = TS_RV_STAR     .Background( tcell.ColorRed    ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_STAR_IN_F = TS_RV_STAR_IN_F.Background( tcell.ColorBlue   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_DEFINE    = TS_RV_DEFINE   .Background( tcell.ColorPurple ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_COMMENT   = TS_RV_COMMENT  .Background( tcell.ColorBlue   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_CONST     = TS_RV_CONST    .Background( tcell.ColorAqua   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_CONTROL   = TS_RV_CONTROL  .Background( tcell.ColorFuchsia).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_VARTYPE   = TS_RV_VARTYPE  .Background( tcell.ColorLime   ).Foreground( tcell.ColorWhite ).Bold(true)
  TS_RV_NONASCII  = TS_RV_NONASCII .Background( tcell.ColorBlue   ).Foreground( tcell.ColorRed   ).Bold(true)
}

