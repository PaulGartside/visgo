
package main

import (
//"fmt"
)

type ChangeHist struct {
  p_fb *FileBuf
  changes ChangeList
}

func (m *ChangeHist) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *ChangeHist) Has_Changes() bool {
  return 0 < m.changes.Len()
}

func (m *ChangeHist) Clear() {
  m.changes.Clear()
}

func (m *ChangeHist) Undo( p_V *FileView ) {

  var plc *LineChange = m.changes.Pop()

  if( nil != plc ) {
    ct := plc.change_t

    if( p_V.in_diff_mode ) {
      if       ( ct ==  CT_INSERT_LINE ) { m.Undo_InsertLine_Diff( plc, p_V )
      } else if( ct ==  CT_REMOVE_LINE ) { m.Undo_RemoveLine_Diff( plc, p_V )
      } else if( ct ==  CT_INSERT_TEXT ) { m.Undo_InsertChar_Diff( plc, p_V )
      } else if( ct ==  CT_REMOVE_TEXT ) { m.Undo_RemoveChar_Diff( plc, p_V )
      } else if( ct == CT_REPLACE_TEXT ) { m.Undo_Set_Diff       ( plc, p_V )
      }
    } else {
      if       ( ct ==  CT_INSERT_LINE ) { m.Undo_InsertLine( plc, p_V )
      } else if( ct ==  CT_REMOVE_LINE ) { m.Undo_RemoveLine( plc, p_V )
      } else if( ct ==  CT_INSERT_TEXT ) { m.Undo_InsertChar( plc, p_V )
      } else if( ct ==  CT_REMOVE_TEXT ) { m.Undo_RemoveChar( plc, p_V )
      } else if( ct == CT_REPLACE_TEXT ) { m.Undo_Set       ( plc, p_V )
      }
    }
  }
}

func (m *ChangeHist) UndoAll( p_V *FileView ) {

  for 0 < m.changes.Len() {
    m.Undo( p_V );
  }
}

func (m *ChangeHist) Save_Set( l_num, c_pos int,
                               old_R rune,
                               continue_last_update bool ) {

  continuation_of_previous_replacement := false

  if( continue_last_update ) {

    NUM_CHANGES := m.changes.Len();
    if( 0<NUM_CHANGES && 0<c_pos ) {

      plc := m.changes.GetP( NUM_CHANGES-1 )
      if( CT_REPLACE_TEXT == plc.change_t &&
          l_num           == plc.lnum &&
          c_pos           == plc.cpos + plc.line.Len() ) {

        plc.line.PushR( old_R )
        continuation_of_previous_replacement = true
      }
    }
  }
  if( !continuation_of_previous_replacement ) {
    // Start of new replacement:
    plc := new( LineChange )
    plc.change_t = CT_REPLACE_TEXT
    plc.lnum = l_num
    plc.cpos = c_pos
    plc.line.PushR( old_R )

    m.changes.Push( plc )
  }
}

func (m *ChangeHist) Save_InsertLine( l_num int ) {

  plc := new( LineChange )
  plc.change_t = CT_INSERT_LINE
  plc.lnum = l_num

  m.changes.Push( plc )
}

func (m *ChangeHist) Save_InsertRune( l_num, c_pos int ) {

  continuation_of_previous_insertion := false

  NUM_CHANGES := m.changes.Len()

  if( 0<NUM_CHANGES  && 0<c_pos ) {

    plc := m.changes.GetP( NUM_CHANGES-1 )
    if( CT_INSERT_TEXT == plc.change_t &&
        l_num          == plc.lnum &&
        c_pos          == ( plc.cpos + plc.line.Len() ) ) {

      plc.line.PushR( 0 )
      continuation_of_previous_insertion = true
    }
  }
  if( !continuation_of_previous_insertion ) {
    // Start of new insertion:
    plc := new( LineChange )
    plc.change_t = CT_INSERT_TEXT
    plc.lnum = l_num
    plc.cpos = c_pos
    plc.line.PushR( 0 )

    m.changes.Push( plc )
  }
}

func (m *ChangeHist) Save_RemoveLine( l_num int, p_fl *FLine ) {

  plc := new( LineChange )
  plc.change_t = CT_REMOVE_LINE
  plc.lnum     = l_num

  // Copy p_fl into plc.line:
  plc.line.Copy( p_fl.runes )

  m.changes.Push( plc )
}

func (m *ChangeHist) Save_RemoveRune( l_num, c_pos int, old_R rune ) {

  continuation_of_previous_removal := false

  NUM_CHANGES := m.changes.Len()

  if( 0<NUM_CHANGES ) {

    plc := m.changes.GetP( NUM_CHANGES-1 )
    if( CT_REMOVE_TEXT == plc.change_t &&
        l_num          == plc.lnum &&
        c_pos          == plc.cpos ) {

      plc.line.PushR( old_R )
      continuation_of_previous_removal = true
    }
  }
  if( !continuation_of_previous_removal ) {
    // Start of new removal:
    plc := new( LineChange )
    plc.change_t = CT_REMOVE_TEXT
    plc.lnum     = l_num
    plc.cpos     = c_pos
    plc.line.PushR( old_R )

    m.changes.Push( plc )
  }
}

func (m *ChangeHist) Undo_InsertLine( plc *LineChange, p_V *FileView ) {
  // Undo an inserted line by removing the inserted line
  m.p_fb.RemoveLP( plc.lnum );

  // If last line of file was just removed, plc.lnum is out of range,
  // so go to NUM_LINES-1 instead:
  NUM_LINES := m.p_fb.NumLines();
  LINE_NUM  := Min_i( plc.lnum, NUM_LINES-1 )

  p_V.GoToCrsPos_NoWrite( LINE_NUM, plc.cpos );

  m.p_fb.Update();
}

func (m *ChangeHist) Undo_RemoveLine( plc *LineChange, p_V *FileView ) {
  // Undo a removed line by inserting the removed line
  m.p_fb.InsertRLP( plc.lnum, &plc.line );

  p_V.GoToCrsPos_NoWrite( plc.lnum, plc.cpos );

  m.p_fb.Update();
}

func (m *ChangeHist) Undo_InsertChar( plc *LineChange, p_V *FileView ) {
  LINE_LEN := plc.line.Len();

  // Undo inserted chars by removing the inserted chars
  for k:=0; k<LINE_LEN; k++ {
    m.p_fb.RemoveR( plc.lnum, plc.cpos );
  }
  p_V.GoToCrsPos_NoWrite( plc.lnum, plc.cpos );

  m.p_fb.Update();
}

func (m *ChangeHist) Undo_RemoveChar( plc *LineChange, p_V *FileView ) {
  LINE_LEN := plc.line.Len();

  // Undo removed chars by inserting the removed chars
  for k:=0; k<LINE_LEN; k++ {
    R := plc.line.GetR(k);

    m.p_fb.InsertR( plc.lnum, plc.cpos+k, R );
  }
  p_V.GoToCrsPos_NoWrite( plc.lnum, plc.cpos );

  m.p_fb.Update();
}

func (m *ChangeHist) Undo_Set( plc *LineChange, p_V *FileView ) {
  LINE_LEN := plc.line.Len();

  for k:=0; k<LINE_LEN; k++ {
    R := plc.line.GetR(k);

    m.p_fb.SetR( plc.lnum, plc.cpos+k, R, false );
  }
  p_V.GoToCrsPos_NoWrite( plc.lnum, plc.cpos );

  m.p_fb.Update();
}

func (m *ChangeHist) Undo_Set_Diff( plc *LineChange, p_V *FileView ) {
// FIXME:
//  LINE_LEN := plc.line.Len();
//
//  for k:=0; k<LINE_LEN; k++ {
//    R := plc.line.GetR(k);
//
//    m.p_fb.SetR( plc.lnum, plc.cpos+k, R, false );
//  }
//  p_Diff := &m_vis.diff
//
//  DL := p_Diff.DiffLine( p_V, plc.lnum );
//  p_Diff.Patch_Diff_Info_Changed( p_V, DL );
//
//  p_Diff.GoToCrsPos_NoWrite( DL, plc.cpos );
//
//  if( !p_Diff.ReDiff() ) { p_Diff.Update() }
}

func (m *ChangeHist) Undo_InsertLine_Diff( plc *LineChange, p_V *FileView ) {
// FIXME:
//  // Undo an inserted line by removing the inserted line
//  m.p_fb.RemoveLP( plc.lnum );
//
//  // If last line of file was just removed, plc.lnum is out of range,
//  // so go to NUM_LINES-1 instead:
//  NUM_LINES := m.p_fb.NumLines();
//  LINE_NUM  := Min_i( plc.lnum, NUM_LINES-1 )
//
//  p_Diff = &m_vis.diff
//
//  DL := p_Diff.DiffLine( p_V, LINE_NUM );
//
//  p_Diff.Patch_Diff_Info_Deleted( p_V, DL );
//
//  p_Diff.GoToCrsPos_NoWrite( DL, plc.cpos );
//
//  if( !p_Diff.ReDiff() ) { p_Diff.Update() }
}

func (m *ChangeHist) Undo_RemoveLine_Diff( plc *LineChange, p_V *FileView ) {
// FIXME:
//  // Undo a removed line by inserting the removed line
//  m.p_fb.InsertRLP( plc.lnum, &plc.line );
//
//  p_Diff = &m_vis.diff
//
//  DL    := p_Diff.DiffLine( p_V, plc.lnum );
//  ODVL0 := p_Diff.On_Deleted_View_Line_Zero( DL );
//
//  p_Diff.Patch_Diff_Info_Inserted( p_V, DL, ODVL0 );
//
//  p_Diff.GoToCrsPos_NoWrite( plc.lnum, plc.cpos );
//
//  if( !p_Diff.ReDiff() ) { p_Diff.Update() }
}

func (m *ChangeHist) Undo_InsertChar_Diff( plc *LineChange, p_V *FileView ) {
// FIXME:
//  LINE_LEN := plc.line.Len();
//
//  // Undo inserted chars by removing the inserted chars
//  for k:=0; k<LINE_LEN; k++ {
//    m.p_fb.RemoveR( plc.lnum, plc.cpos );
//  }
//  p_Diff = &m_vis.diff
//
//  DL := p_Diff.DiffLine( p_V, plc.lnum );
//  p_Diff.Patch_Diff_Info_Changed( p_V, DL );
//
//  p_Diff.GoToCrsPos_NoWrite( DL, plc.cpos );
//
//  if( !p_Diff.ReDiff() ) { p_Diff.Update() }
}

func (m *ChangeHist) Undo_RemoveChar_Diff( plc *LineChange, p_V *FileView ) {
// FIXME:
//  LINE_LEN := plc.line.len();
//
//  // Undo removed chars by inserting the removed chars
//  for k:=0; k<LINE_LEN; k++ {
//    R := plc.line.GetR(k);
//
//    m.p_fb.InsertR( plc.lnum, plc.cpos+k, R );
//  }
//  p_Diff = &m_vis.diff
//
//  DL := p_Diff.DiffLine( p_V, plc.lnum );
//  p_Diff.Patch_Diff_Info_Changed( p_V, DL );
//
//  p_Diff.GoToCrsPos_NoWrite( DL, plc.cpos );
//
//  if( !p_Diff.ReDiff() ) { p_Diff.Update() }
}

