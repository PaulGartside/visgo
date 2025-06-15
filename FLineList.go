
package main

//import "fmt"

type FLineList struct {
  lines []*FLine
  LF_at_EOF bool
  offsets []int // Line Offests: Byte offset of the first rune of each line in file
  hi_touched_line int // Line before which highlighting is valid
}

// Create slice of length filled with empty FLine's
func (m *FLineList) Init( length int ) {
  m.lines = make( []*FLine, length )
}

func (m *FLineList) Clear() {
  m.lines = m.lines[:0]

  m.ChangedLine( 0 )
}

func (m *FLineList) Len() int {
  return len( m.lines )
}

func (m *FLineList) Cap() int {
  return cap( m.lines )
}

// Return size in bytes of file represented by this FLineList
//
func (m *FLineList) GetSize() int {
  size := 0
  NUM_LINES := m.Len()
  if( 0 < NUM_LINES ) {
    // Absolute byte offset of beginning of first line in file is always zero:
    if( 0 == len( m.offsets ) ) {
      m.offsets = append( m.offsets, 0 )
    }
    // Old line offsets length:
    OLOL := len( m.offsets )

    // New line offsets length:
    m.SetOffsetsLen( NUM_LINES )

    for k:=OLOL; k<NUM_LINES; k++ {
      m.offsets[ k ] = m.offsets[ k-1 ] + m.lines[ k-1 ].Len() + 1 //< Add 1 for '\n'
    }
    size = m.offsets[ NUM_LINES-1 ] + m.lines[ NUM_LINES-1 ].Len()
    if( m.LF_at_EOF ) { size++ }
  }
  return size
}

func (m *FLineList) GetCursorByte( CL, CC int ) int {
  crs_byte := 0
  NUM_LINES := m.Len()
  if 0 < NUM_LINES {
    // Make sure CL are is range:
    CL = Min_i( CL, NUM_LINES-1 )
    // Make sure CC is in range:
    CLL := m.GetLP(CL).Len()
    if( CLL <= CC ) {
      CC = 0
      if( 0 < CLL ) { CC = CLL-1 }
    }
    // Old line offsets length:
    OLOL := len( m.offsets )

    if( OLOL != NUM_LINES ) {
      m.SetOffsetsLen( NUM_LINES )
      for k:=OLOL; k<NUM_LINES; k++ {
        m.offsets[ k ] = m.offsets[ k-1 ] + m.LineLen( k-1 ) + 1 //< Add 1 for '\n'
      }
    }
    crs_byte = m.offsets[ CL ] + CC
  }
  return crs_byte
}

// Copy *p_src into self
func (m *FLineList) CopyP( p_src *FLineList ) {
  m.SetLen( p_src.Len() )
  copy( m.lines[:], p_src.lines[:] )
  m.LF_at_EOF = p_src.LF_at_EOF

  m.ChangedLine( 0 )
}

// Set m.lines len to length.
// If capacity is increased, existing contents of m.lines is lost.
//
func (m *FLineList) SetLen( length int ) {

  if( length < m.Len() ) {
    m.lines = m.lines[:length]
  } else if( m.Len() < length ) {
    if( length <= m.Cap() ) {
      for ; m.Len() < length; {
        p_fl := new(FLine)
        m.lines = append( m.lines, p_fl )
      }
    } else {
      m.Init( length )
    }
  }
}

// Set m.offsets len to length.
// Existing contents of m.offsets are preserved.
//
func (m *FLineList) SetOffsetsLen( length int ) {

  // Old line offsets length:
  OLOL := len( m.offsets )

  if( length < OLOL ) {
    m.offsets = m.offsets[:length]
  } else if( OLOL < length ) {
    if( length <= cap(m.offsets) ) {
      for ; len(m.offsets) < length; {
        m.offsets = append( m.offsets, 0 )
      }
    } else {
      old_offsets := m.offsets
      m.offsets = make( []int, length )
      copy( m.offsets[:OLOL], old_offsets[:OLOL] )
    }
  }
}

func (m *FLineList) LineLen( l_num int ) int {

  return m.lines[ l_num ].Len()
}

func (m *FLineList) GetLP( l_num int ) *FLine {

  return m.lines[ l_num ]
}

func (m *FLineList) GetR( l_num, r_num int ) rune {

  return m.lines[ l_num ].GetR( r_num )
}

func (m *FLineList) SetR( l_num, r_num int, R rune ) {

  m.lines[ l_num ].SetR( r_num, R )

  m.ChangedLine( l_num )
}

func (m *FLineList) RemoveR( l_num, r_num int ) rune {

  m.ChangedLine( l_num )

  R := m.lines[ l_num ].RemoveR( r_num )
  return R
}

func (m *FLineList) InsertR( l_num, r_num int, R rune ) {

  m.lines[ l_num ].InsertR( r_num, R )

  m.ChangedLine( l_num )
}

func (m *FLineList) PushR( l_num int, R rune ) {

  m.lines[ l_num ].PushR( R )

  m.ChangedLine( l_num )
}

// Insert an empty line at position l_num
func (m *FLineList) InsertLE( l_num int ) {

  // First append an empty line to make sure lines is large enough:
  var ln FLine
  m.lines = append( m.lines, &ln )

  // All the lines from l_num on shift up so empty line is gone:
  copy( m.lines[l_num+1:], m.lines[l_num:] )
  // Put empty line at position l_num:
  m.lines[ l_num ] = &ln

  m.ChangedLine( l_num )
}

func (m *FLineList) InsertRLP( l_num int, p_rln *RLine ) {

  // First append an empty line to make sure lines is large enough:
  p_fln := new( FLine )
  m.lines = append( m.lines, p_fln )

  // All the lines from l_num on shift up so empty line is gone:
  copy( m.lines[l_num+1:], m.lines[l_num:] )
  // Turn empty line into non-empty line:
  p_fln.CopyPRL( p_rln )
  // Put non-emtpy line at position l_num:
  m.lines[ l_num ] = p_fln

  m.ChangedLine( l_num )
}

func (m *FLineList) RemoveLP( l_num int ) *FLine {

  var p_ln *FLine = m.lines[ l_num ]

  copy( m.lines[l_num:], m.lines[l_num+1:] )
  m.lines = m.lines[:len(m.lines)-1]

  m.ChangedLine( l_num )

  return p_ln
}

func (m *FLineList) PushLP( p_ln *FLine ) {

  m.lines = append( m.lines, p_ln )
}

func (m *FLineList) AppendLineToLine( l_num int, p_fl *FLine ) {

  m.lines[ l_num ].PushLP( p_fl )

  m.ChangedLine( l_num )
}

func (m *FLineList) Swap( l_num_1, l_num_2 int ) {

  m.lines[l_num_1], m.lines[l_num_2] = m.lines[l_num_2], m.lines[l_num_1]

  m.ChangedLine( Min_i(l_num_1, l_num_2) )
}

func (m *FLineList) ChangedLine( line_num int ) {

  if( 0<=line_num && line_num<len(m.offsets) ) {
    m.SetOffsetsLen( line_num )
  }
  m.hi_touched_line = Min_i( m.hi_touched_line, line_num )
}

