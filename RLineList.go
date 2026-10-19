
package main

type RLineList struct {
  lines []*RLine
}

// Create slice of length filled with empty RLine's
//
func (m *RLineList) Init( length int ) {
  m.lines = make( []*RLine, length )
  for k:=0; k<length; k++ {
    m.lines[k] = new(RLine)
  }
}

func (m *RLineList) Clear() {
  m.lines = m.lines[:0]
}

func (m *RLineList) Len() int {
  return len( m.lines )
}

func (m *RLineList) Cap() int {
  return cap( m.lines )
}

// Copy *p_src into self
func (m *RLineList) CopyP( p_src *RLineList ) {
  LEN := p_src.Len()
  m.SetLen( LEN )

  for k:=0; k<LEN; k++ {
    m.lines[k].CopyP( p_src.lines[k] )
  }
}

// Set length while guaranteeing existing contents remain the same:
//
func (m *RLineList) SetLen( length int ) {
  LEN_OLD := m.Len()

  if( length < m.Len() ) {
    m.lines = m.lines[:length]

  } else if( LEN_OLD < length ) {
    if( length <= m.Cap() ) {
      for ; m.Len() < length; {
        p_rl := new(RLine)
        m.lines = append( m.lines, p_rl )
      }
    } else {
      var lines_old []*RLine = m.lines
      m.lines = make( []*RLine, length )
      for k:=0; k<LEN_OLD; k++ {
        m.lines[k] = lines_old[k]
      }
      for k:=LEN_OLD; k<length; k++ {
        m.lines[k] = new(RLine)
      }
    }
  }
}

//func (m *RLineList) LineLen( l_num int ) int {
//
//  return m.lines[ l_num ].Len()
//}

func (m *RLineList) GetLP( l_num int ) *RLine {

  return m.lines[ l_num ]
}

//func (m *RLineList) GetR( l_num, r_num int ) rune {
//
//  return m.lines[ l_num ].GetR( r_num )
//}

func (m *RLineList) RemoveR( l_num, r_num int ) rune {

  return m.lines[ l_num ].RemoveR( r_num )
}

func (m *RLineList) InsertR( l_num, r_num int, R rune ) {

  m.lines[ l_num ].InsertR( r_num, R )
}

func (m *RLineList) PushR( l_num int, R rune ) {

  m.lines[ l_num ].PushR( R )
}

// Insert an empty line at position l_num
func (m *RLineList) InsertLE( l_num int ) {

  // First append an empty line to make sure lines is large enough:
  var ln RLine
  m.lines = append( m.lines, &ln )

  // All the lines from l_num on shift up so empty line is gone:
  copy( m.lines[l_num+1:], m.lines[l_num:] )
  // Put empty line at position l_num:
  m.lines[ l_num ] = &ln
}

func (m *RLineList) RemoveLP( l_num int ) *RLine {

  var p_ln *RLine = m.lines[ l_num ]

  copy( m.lines[l_num:], m.lines[l_num+1:] )
  m.lines = m.lines[:len(m.lines)-1]

  return p_ln
}

func (m *RLineList) PushLP( p_ln *RLine ) {

  m.lines = append( m.lines, p_ln )
}

func (m *RLineList) PushLE() {
  var ln RLine
  m.lines = append( m.lines, &ln )
}

func (m *RLineList) AppendLineToLine( l_num int, p_ln *RLine ) {

//m.lines[ l_num ].PushLP( p_ln )
  m_p_ln := m.lines[ l_num ]

  for k:=0; k<p_ln.LenB(); k++ {
    m_p_ln.PushB( p_ln.GetB( k ) )
  }
}

