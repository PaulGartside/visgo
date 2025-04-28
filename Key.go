
package main

import (
//"fmt"

//"github.com/gdamore/tcell/v2"
)

type Key struct {
  // Public:
  save_2_dot_buf_n bool // Normal view
  save_2_dot_buf_l bool // Line view
  save_2_vis_buf bool
  save_2_map_buf bool
  get_from_dot_buf_n bool // Normal view
  get_from_dot_buf_l bool // Line view
  get_from_map_buf bool

  dot_buf_n Vector[Key_rune] // Dot buf for normal view
  dot_buf_l Vector[Key_rune] // Dot buf for line view
  vis_buf   Vector[Key_rune]
  map_buf   Vector[Key_rune]

  // Private:
  dot_buf_index_n int // Normal view
  dot_buf_index_l int // Line view
  map_buf_index   int

  console Console
}

func (m *Key) In() Key_rune {
  var kr Key_rune

  if       ( m.get_from_map_buf   ) { kr = m.In_MapBuf()
  } else if( m.get_from_dot_buf_n ) { kr = m.In_DotBuf_n()
  } else if( m.get_from_dot_buf_l ) { kr = m.In_DotBuf_l()
  } else                            { kr = m_console.Key_In()
  }
  if( m.save_2_map_buf   ) { m.map_buf.Push( kr ) }
  if( m.save_2_dot_buf_n ) { m.dot_buf_n.Push( kr ) }
  if( m.save_2_dot_buf_l ) { m.dot_buf_l.Push( kr ) }
  if( m.save_2_vis_buf   ) { m.vis_buf.Push( kr ) }

  return kr
}

func (m *Key) In_DotBuf_n() Key_rune {

  kr := m.dot_buf_n.Get( m.dot_buf_index_n )
  m.dot_buf_index_n++

  if( m.dot_buf_n.Len() <= m.dot_buf_index_n ) {
    m.get_from_dot_buf_n = false
    m.dot_buf_index_n    = 0
  }
  return kr
}

func (m *Key) In_DotBuf_l() Key_rune {

  kr := m.dot_buf_l.Get( m.dot_buf_index_l )
  m.dot_buf_index_l++

  if( m.dot_buf_l.Len() <= m.dot_buf_index_l ) {
    m.get_from_dot_buf_l = false
    m.dot_buf_index_l    = 0
  }
  return kr
}

func (m *Key) In_MapBuf() Key_rune {

  kr := m.map_buf.Get( m.map_buf_index )
  m.map_buf_index++

  if( m.map_buf.Len() <= m.map_buf_index ) {
    m.get_from_map_buf = false
    m.map_buf_index    = 0
  }
  return kr
}

//func (m *Key) MapBuf_pop() {
//  LEN_MAP_BUF := m.map_buf.Len()
//  if( 0 < LEN_MAP_BUF ) {
//    m.map_buf = m.map_buf[:LEN_MAP_BUF-1]
//  }
//}

