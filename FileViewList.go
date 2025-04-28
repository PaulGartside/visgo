
package main

type FileViewList struct {
  views []*FileView
}

// Create slice of length filled with nil *FileView's
func (m *FileViewList) Init( length int ) {
  m.views = make( []*FileView, length )
}

func (m *FileViewList) Len() int {
  return len( m.views )
}

func (m *FileViewList) Cap() int {
  return cap( m.views )
}

// Set length to 0
func (m *FileViewList) Clear() {
  // Set all pointers to nil so the objects they point to
  // can be garbage collected:
  for k:=0; k<len(m.views); k++ {
    m.views[k] = nil
  }
  m.views = m.views[:0]
}

// Set length without guaranteeing existing contents remain the same:
func (m *FileViewList) SetLen( length int ) {

  if( length < m.Len() ) {
    m.views = m.views[:length]
  } else if( m.Len() < length ) {
    if( length <= m.Cap() ) {
      for ; m.Len() < length; { m.views = append( m.views, nil ) }
    } else {
      m.Init( length )
    }
  }
}

func (m *FileViewList) CopyP( p_src *FileViewList ) {
  m.SetLen( p_src.Len() )
  copy( m.views, p_src.views )
}

func (m *FileViewList) GetPFv( f_num int ) *FileView {
  return m.views[ f_num ]
}

func (m *FileViewList) RemovePFv( f_num int ) *FileView {

  var p_fb *FileView = m.views[ f_num ]
  copy( m.views[f_num:], m.views[f_num+1:] )

  // Help with garbage collection:
  m.views[len(m.views)-1] = nil

  m.views = m.views[:len(m.views)-1]
  return p_fb
}

func (m *FileViewList) PushPFv( p_fb *FileView ) {
  m.views = append( m.views, p_fb )
} 

func (m *FileViewList) InsertPFv( f_num int, p_fb *FileView ) {
  // First append p_fb to make sure views is large enough:
  m.PushPFv( p_fb ) 
  copy( m.views[f_num+1:], m.views[f_num:] )
  m.views[ f_num ] = p_fb
}

