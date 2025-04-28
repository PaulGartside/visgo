
package main

type Highlight_Dir struct {
  p_fb *FileBuf
}

func (m *Highlight_Dir) Init(p_fb *FileBuf) {
  m.p_fb = p_fb
}

func (m *Highlight_Dir) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_Dir) Run_Range( st CrsPos, fn int )")
}

