
package main

import (
//"fmt"

  "github.com/gdamore/tcell/v2"
)

type Key_rune struct {
  K tcell.Key
  R rune
}

func make_Key_rune( R rune ) Key_rune {

  return Key_rune{ tcell.KeyRune, R }
}

// Returns true if this Key_rune represents a printable character
func (m *Key_rune) IsKeyRune() bool {

  return m.K == tcell.KeyRune
}

// Is Escape
func (m *Key_rune) IsESC() bool {

  return m.K == tcell.KeyESC
}

// Is Backspace
func (m *Key_rune) IsBS() bool {

  return m.K == tcell.KeyBS
}

// Is Delete
func (m *Key_rune) IsDEL() bool {

  return m.K == tcell.KeyDEL
}

func (m *Key_rune) IsEndOfLineDelim() bool {

  is_EOL_delim := false

  if( m.K == tcell.KeyLF || m.K == tcell.KeyCR ) {
    is_EOL_delim = true
  }
  return is_EOL_delim
}

