
package main

import (
//"fmt"
)

type LineChange struct {
  change_t ChangeType
  lnum     int
  cpos     int
  line     RLine
}

type ChangeList struct {
  data []*LineChange
}

func (m *ChangeList) Clear() {
  m.data = m.data[:0]
}

func (m *ChangeList) Len() int {
  return len( m.data )
}

func (m *ChangeList) Cap() int {
  return cap( m.data )
}

func (m *ChangeList) GetP( l_num int ) *LineChange {
  return m.data[ l_num ]
}

// Insert an empty line at position l_num
func (m *ChangeList) InsertP( l_num int, p_lc *LineChange ) {

  // First append an LineChange to make sure data is large enough:
  m.data = append( m.data, p_lc )

  // All the data from l_num on shift up so empty line is gone:
  copy( m.data[l_num+1:], m.data[l_num:] )
  // Put empty line at position l_num:
  m.data[ l_num ] = p_lc
}

func (m *ChangeList) RemoveP( l_num int ) *LineChange {

  var p_lc *LineChange = m.data[ l_num ]

  copy( m.data[l_num:], m.data[l_num+1:] )
  m.data = m.data[:len(m.data)-1]

  return p_lc
}

func (m *ChangeList) Push( p_lc *LineChange )  {
  m.data = append( m.data, p_lc )
}

func (m *ChangeList) Pop() *LineChange {

  var p_lc *LineChange = nil

  LEN := m.Len()

  if( 0 < LEN ) {
    p_lc = m.RemoveP( LEN-1 )
  }
  return p_lc
}

