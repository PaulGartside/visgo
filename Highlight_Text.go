
package main

type Highlight_Text struct {
  p_fb *FileBuf
}

func (m *Highlight_Text) Init(p_fb *FileBuf) {
  m.p_fb = p_fb
}

func (m *Highlight_Text) Run_Range( st CrsPos, fn int ) {
//Log("In func (m *Highlight_Text) Run_Range( st CrsPos, fn int )")
}

