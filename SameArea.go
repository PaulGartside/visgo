
package main

type SameArea struct {
  ln_s   int // Beginning line number in short file
  ln_l   int // Beginning line number in long  file
  nlines int // Number of consecutive lines the same
  nbytes int // Number of bytes in consecutive lines the same
}

func (m *SameArea) Clear() {
  m.ln_s   = 0
  m.ln_l   = 0
  m.nlines = 0
  m.nbytes = 0
}

func (m *SameArea) Init( _ln_s, _ln_l, _nbytes int ) {
  m.ln_s   = _ln_s
  m.ln_l   = _ln_l
  m.nlines = 1
  m.nbytes = _nbytes
}

func (m *SameArea) Inc( _nbytes int ) {
  m.nlines += 1
  m.nbytes += _nbytes
}

func (m *SameArea) Set( sa SameArea ) {
  m.ln_s   = sa.ln_s
  m.ln_l   = sa.ln_l
  m.nlines = sa.nlines
  m.nbytes = sa.nbytes
}

