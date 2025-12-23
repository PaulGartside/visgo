
package main

import (
  "fmt"
)

// Diff or Comparison area
//
type DiffArea struct {
  ln_s     int // Beginning line number in short file
  nlines_s int // Number of consecutive lines different in short file
  ln_l     int // Beginning line number in long  file
  nlines_l int // Number of consecutive lines different in long  file
}

//func (m *DiffArea) Init( ln_s, nlines_s, ln_l, nlines_l int ) {
//  m.ln_s     = ln_s    
//  m.ln_l     = ln_l    
//  m.nlines_s = nlines_s
//  m.nlines_l = nlines_l
//}

// fnl_s() is one past last line of DiffArea on the short side
//
func (m *DiffArea) fnl_s() int {
  return m.ln_s + m.nlines_s
}

// fnl_l() is one past last line of DiffArea on the long side
//
func (m *DiffArea) fnl_l() int {
  return m.ln_l + m.nlines_l
}

func (m *DiffArea) Print() {
  msg := fmt.Sprintf("DiffArea: lines_s=(%u,%u) lines_l=(%u,%u)\n",
                     m.ln_s+1, m.fnl_s(), m.ln_l+1, m.fnl_l())
  Log( msg )
}

