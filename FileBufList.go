
package main

type FileBufList struct {
  files []*FileBuf
}

func (m *FileBufList) Clear() {
  // Set all pointers to nil so the objects they point to
  // can be garbage collected:
  for k:=0; k<len(m.files); k++ {
    m.files[k] = nil
  }
  m.files = m.files[:0]
}

func (m *FileBufList) Len() int {
  return len( m.files )
}

func (m *FileBufList) GetPFb( f_num int ) *FileBuf {
  return m.files[ f_num ]
}

func (m *FileBufList) RemovePFb( f_num int ) *FileBuf {

  var p_fb *FileBuf = m.files[ f_num ]
  copy( m.files[f_num:], m.files[f_num+1:] )

  // Help with garbage collection:
  m.files[len(m.files)-1] = nil

  m.files = m.files[:len(m.files)-1]
  return p_fb
}

func (m *FileBufList) PushPFb( p_fb *FileBuf ) {
  m.files = append( m.files, p_fb )
} 

func (m *FileBufList) InsertPFb( f_num int, p_fb *FileBuf ) {
  // First append p_fb to make sure files is large enough:
  m.PushPFb( p_fb ) 
  copy( m.files[f_num+1:], m.files[f_num:] )
  m.files[ f_num ] = p_fb
}

