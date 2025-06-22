
package main

import (
//"bytes"
  "fmt"
//"io/ioutil"
//"log"
  "os"
  "strconv"
//"strings"
  "time"
)

type CmdFunc func(*Vis)
type CmdFunc2 func(*Vis,int)

type Vis struct {
  running bool

  files   FileBufList
  views   [MAX_WINS]FileViewList

  view_funcs [128]CmdFunc
  line_funcs [128]CmdFunc

  colon_file FileBuf
  slash_file FileBuf
  colon_view LineView
  slash_view LineView

  colon_mode bool
  slash_mode bool

  reg RLineList
  paste_mode Paste_Mode

  fast_rune rune // Char on line to goto when ';' is entered

  win      int // Sub-window index
  num_wins int // Number of sub-windows currently on screen
  file_hist [MAX_WINS]IntList

  regex_str string

  diff Diff
}

func (m *Vis) HaveFile( path string, p_file_index *int ) bool {

  var already_have_file bool = false

  for k:=0; k<m.files.Len(); k++ {

    var fpath string = m.files.GetPFb(k).GetPath()
    if path == fpath {
      already_have_file = true
      if( p_file_index != nil ) {
        *p_file_index = k
      }
      break
    }
  }
  return already_have_file
}

func (m *Vis) NotHaveFileAddFile( path string ) bool {

  var added_file bool = false

  if( !m.HaveFile( path, nil ) ) {
    m.CreateFile( path, FT_UNKNOWN )

    added_file = true
  }
  return added_file
}

func (m *Vis) Add_FileBuf_2_Lists_Create_Views( p_fb *FileBuf ) {

  m.files.PushPFb( p_fb )

  for w:=0; w<MAX_WINS; w++ {
    p_v := new( FileView )
    p_v.Init( p_fb )
    m.views[w].PushPFv( p_v )
    p_fb.AddView( p_v )
  }
  // Push file name onto buffer editor buffer
  m.AddToBufferEditor( p_fb.path_name )
}

func (m *Vis) AddToBufferEditor( fname string ) {
  p_fl := new(FLine)
  p_fl.PushSR( []rune(fname) )

  var p_fb *FileBuf = m.views[0].GetPFv( m_BE_FILE ).p_fb
  p_fb.PushLP( p_fl )
  p_fb.BufferEditor_SortName()
  p_fb.ClearChanged()

  // Since buffer editor file has been re-arranged, make sure none of its
  // views have the cursor position past the end of the line
  for k:=0; k<MAX_WINS; k++ {
    var p_V *FileView = m.views[k].GetPFv( m_BE_FILE )

    var CL int = p_V.CrsLine()
    var CP int = p_V.CrsChar()
    var LL int = p_fb.LineLen( CL )

    if( LL <= CP ) {
      p_V.GoToCrsPos_NoWrite( CL, LL-1 )
    }
  }
}

func (m *Vis) GetFileBuf( index int ) *FileBuf {

  return m.files.GetPFb( index )
}

func (m *Vis) GetFileBuf_s( fname string ) *FileBuf {

  for k:=0; k<m.files.Len(); k++ {
    var pfb_k *FileBuf = m.files.GetPFb( k )

    if( fname == pfb_k.path_name ) {
      return pfb_k
    }
  }
  return nil
}

// Return true if went to buffer indicated by fname, else false
//func (m *Vis) GoToBuffer_Fname( fname string ) bool {
//
//  var went_to_buffer bool = false
//
//  // 1. Search for fname in buffer list, and if found, go to that buffer:
//  var file_index int = 0
//  if( m.HaveFile( fname, &file_index ) ) {
//    m.GoToBuffer( file_index )
//
//    went_to_buffer = true
//
//  // 2. Get full file name of fname relative to dir of current file
//  } else if( m.GetFullFileNameRelative2CurrFile( fname ) ) {
//    // 3. Search for fname in buffer list, and if found, go to that buffer:
//    if( m.HaveFile( fname, &file_index ) ) {
//      m.GoToBuffer( file_index )
//
//      went_to_buffer = true
//    // 4. See if file exists, and if so, add a file buffer, and go to that buffer
//    } else if( FileExists( fname ) ) {
//      p_fb := new( FileBuf )
//      p_fb.Init( fname )
//
//      if( m.HaveFile( fname, &file_index ) ) {
//        m.GoToBuffer( file_index 
//
//        went_to_buffer = true
//      }
//    }
//  }
//
//  if( ! went_to_buffer ) {
//    m.CmdLineMessage( fmt.Sprintf("Could not find file: %s", fname) )
//  }
//  return went_to_buffer
//}

// Return true if went to buffer indicated by fname, else false
func (m *Vis) GoToBuffer_Fname( fname string ) bool {

  var went_to_buffer bool = false

  // 1. Search for fname in buffer list, and if found, go to that buffer:
  var file_index int = 0
  if( m.HaveFile( fname, &file_index ) ) {
    m.GoToBuffer( file_index )
    went_to_buffer = true

  // 2. Get full file name of fname relative to dir of current file
  } else {
    fpath, ok := m.GetFullFileNameRelative2CurrFile( fname )
    if( ok ) {
      // 3. Search for fpath in buffer list, and if found, go to that buffer:
      if( m.HaveFile( fpath, &file_index ) ) {
        m.GoToBuffer( file_index )

        went_to_buffer = true
      // 4. See if file exists, and if so, add a file buffer, and go to that buffer
      } else if( FileExists( fpath ) ) {
        m.CreateFile( fpath, FT_UNKNOWN )

        if( m.HaveFile( fpath, &file_index ) ) {
          m.GoToBuffer( file_index )

          went_to_buffer = true
        } else {
          // Should never get here:
        }
      }
    }
  }

  if( ! went_to_buffer ) {
    m.CmdLineMessage( fmt.Sprintf("Could not find file: %s", fname) )
  }
  return went_to_buffer
}

func (m *Vis) GoToPrevBuffer() {

  var FILE_HIST_LEN int = m.file_hist[m.win].Len()

  if( FILE_HIST_LEN <= 1 ) {
    // Nothing to do, so just put cursor back
    m.CV().PrintCursor()

  } else {
    var went_back_to_prev_dir_diff bool = false

  //if( m.CV().in_diff_mode )
  //{
  //  went_back_to_prev_dir_diff = WentBackToPrevDirDiff(m)
  //
  //  if( !went_back_to_prev_dir_diff ) NoDiff_CV(m)
  //}
    if( !went_back_to_prev_dir_diff ) {

    //var pV_old *FileView = m.CV()
    //Tile_Pos const tp_old = pV_old->GetTilePos()

      // Move view index at front to back of m.file_hist
      var view_index_old int = m.file_hist[m.win].Remove( 0 )
      m.file_hist[m.win].Push( view_index_old )

      // For DIR and BUFFER_EDITOR, invalidate regex's so that files that
      // no longer contain the current regex are no longer highlighted
    //if( m.CV().GetFB().GetFileType() == FT_DIR
    // || m.CV().GetFB().GetFileType() == FT_BUFFER_EDITOR )
    //{
    //  m.CV().GetFB().Invalidate_Regexs()
    //}
      // Redisplay current window with new view:
    //CV(m)->SetTilePos( tp_old )
      m.CV().Update_and_PrintCursor()
    }
  }
}

func (m *Vis) GoToPoundBuffer() {

  m.NoDiff_CV()

  if( m_BE_FILE == m.file_hist[m.win].Get(1) ) {

    m.GoToBuffer( m.file_hist[m.win].Get(2) )
  } else {
    m.GoToBuffer( m.file_hist[m.win].Get(1) )
  }
}

func (m *Vis) GoToCurrBuffer() {

  // CVI = Current View Index
  var CVI int = m.file_hist[m.win].Get(0)

  if( CVI == m_BE_FILE    ||
      CVI == m_HELP_FILE  ||
      CVI == m_COLON_FILE ) {
    // FIXME:
    //CVI == m_SLASH_FILE )

    m.NoDiff_CV()
    m.GoToBuffer( m.file_hist[m.win].Get(1) )

  } else {
    // User asked for view that is currently displayed.
    // Dont do anything, just put cursor back in place.
    m.CV().PrintCursor()
  }
}

func (m *Vis) GoToBufferEditor() {

  m.GoToBuffer( m_BE_FILE )
}

func (m *Vis) GoToMsgBuffer() {

  m.GoToBuffer( m_MSG_FILE )
}

func (m *Vis) NoDiff_CV() {
  // FIXME
}

func (m *Vis) GetFullFileNameRelative2CurrFile( fname string ) (string, bool) {

  var pname string
  var got_full_file_name bool = false

  // 2. Get full file name
  if( fname[0] == DIR_DELIM && FileExists( fname ) ) {
    pname = fname
    got_full_file_name = true; // fname is already a full file name

  } else {
    pname = FindFullFileNameRel2( m.CV().p_fb.dir_name, fname )

    got_full_file_name = true; // fname now contains full file name
  }
  return pname, got_full_file_name
}

func (m *Vis) GoToBuffer( buf_idx int ) {

  if( m.views[m.win].Len() <= buf_idx ) {
    m.CmdLineMessage( fmt.Sprintf("Buffer %lu does not exist", buf_idx) )

  } else {
    m.NoDiff_CV()

    if( buf_idx == m.file_hist[m.win].Get(0) ) {
      // User asked for view that is currently displayed.
      // Dont do anything, just put cursor back in place.
      m.CV().PrintCursor()

    } else {
      m.file_hist[m.win].Insert( 0, buf_idx )

      // Remove subsequent buf_idx's from m.file_hist[m.win]:
      for k:=1; k<m.file_hist[m.win].Len(); k++  {
        if( buf_idx == m.file_hist[m.win].Get(k) ) {
          m.file_hist[m.win].Remove( k )
        }
      }
      // FIXME:
      var p_nv *FileView = m.CV(); // New FileView to display
      if( ! p_nv.Has_Context() ) {
        // Look for context for the new view:
        var found_context bool = false
        for w:=0; !found_context && w<m.num_wins; w++ {
          var p_v *FileView = m.views[w].GetPFv( buf_idx )
          if( p_v.Has_Context() ) {
            found_context = true
            p_nv.Set_Context( p_v )
          }
        }
      }
      // For DIR and BUFFER_EDITOR, invalidate regex's so that files that
      // no longer contain the current regex are no longer highlighted
    // FIXME:
    //if( p_nv.GetFB()->GetFileType() == FT_DIR
    // || p_nv.GetFB()->GetFileType() == FT_BUFFER_EDITOR )
    //{
    //  p_nv.GetFB()->Invalidate_Regexs()
    //}
      p_nv.SetTilePos( m.PV().GetTilePos() )
      p_nv.Update_and_PrintCursor()
    }
  }
}

//func (m *Vis) Create_FileBuf_and_LineViews( path string ) *LineView {
//  p_fb := new( FileBuf )
//  m.files = append( m.files, p_fb )
//
//  p_v := new( FileView )
//  p_v.Init( p_fb )
//  m.views = append( m.views, p_v )
//
//  p_fb.Init( path, p_v )
//
//  return p_v
//}

func (m *Vis) CreateFile( pname string, file_type File_Type ) *FileBuf {

  p_fb := new( FileBuf )
  p_fb.Init_FB( pname, file_type )

//m.Add_FileBuf_2_Lists_Create_Views( p_fb )

  return p_fb
}

func (m *Vis) InitBufferEditor() {

  m.CreateFile( m_EDIT_BUF_NAME, FT_BUFFER_EDITOR )
}

func (m *Vis) InitHelpBuffer() {

  m.CreateFile( m_HELP_BUF_NAME, FT_TEXT )
}

func (m *Vis) InitMsgBuffer() {

  m.CreateFile( m_MSG__BUF_NAME, FT_TEXT )
}

func (m *Vis) InitShellBuffer() {

  m.CreateFile( m_SHELL_BUF_NAME, FT_TEXT )
}

func (m *Vis) InitColonBuffer() {

  m.colon_file.Init_FB( m_COLON_BUF_NAME, FT_TEXT )
  m.colon_file.p_lv = &m.colon_view

  m.colon_view.Init( &m.colon_file, ':' )

//m.Add_FileBuf_2_Lists_Create_Views( &m.colon_file )
}

func (m *Vis) InitSlashBuffer() {

  m.slash_file.Init_FB( m_SLASH_BUF_NAME, FT_TEXT )
  m.slash_file.p_lv = &m.slash_view

  m.slash_view.Init( &m.slash_file, '/' )

//m.Add_FileBuf_2_Lists_Create_Views( &m.slash_file )
}

func (m *Vis) InitUserFiles_AddFile( relative_name string ) {

  path := FindFullFileNameRel2CWD( relative_name )

  if !m.HaveFile( path, nil ) {
    m.CreateFile( path, FT_UNKNOWN )
  }
}

func (m *Vis) InitUserFiles() {

  ARGC := len( os.Args )
  for k:=1; k<ARGC; k++ {
    m.InitUserFiles_AddFile( os.Args[k] )
  }
  m.InitUserFiles_AddFile(".")
}

func (m *Vis) InitFileHistory() {

  for w:=0; w<MAX_WINS; w++ {

    m.file_hist[w].Push( m_BE_FILE )
    m.file_hist[w].Push( m_HELP_FILE )

    if( m_USER_FILE<m.views[w].Len() ) {

      m.file_hist[w].Insert( 0, m_USER_FILE )

      for f_num:=m.views[w].Len()-1; (m_USER_FILE+1)<=f_num; f_num-- {
        m.file_hist[w].Push( f_num )
      }
    }
  }
}

func (m *Vis) FileName_Is_Displayed( full_fname string ) bool {

  var file_num int = 0

  if( m.FName_2_FNum( full_fname, &file_num ) ) {

    return m.FileNum_Is_Displayed( file_num )
  }
  return false
}

func (m *Vis) FileNum_Is_Displayed( file_num int ) bool {

  for w:=0; w<m.num_wins; w++ {

    if( file_num == m.file_hist[ w ].Get( 0 ) ) {
      return true
    }
  }
  return false
}

func (m *Vis) ReleaseFileName( full_fname string ) {

  var file_num int = 0

  if( m.FName_2_FNum( full_fname, &file_num ) ) {

    m.ReleaseFileNum( file_num )
  }
}

func (m *Vis) ReleaseFileNum( file_num int ) {

  m.files.RemovePFb( file_num )

  for k:=0; k<MAX_WINS; k++  {
    // Remove m.views[file_num]
    m.views[k].RemovePFv( file_num )

    var p_file_hist_k *IntList = &m.file_hist[k]

    // Remove all file_num's from m_file_hist
    for i:=0; i<p_file_hist_k.Len(); i++ {

      if( file_num == p_file_hist_k.Get( i ) ) {
        p_file_hist_k.Remove( i )
      }
    }
    // Decrement all file_hist numbers greater than file_num
    for i:=0; i<p_file_hist_k.Len(); i++ {

      var val int = p_file_hist_k.Get( i )

      if( file_num < val ) {
        p_file_hist_k.Set( i, val-1 )
      }
    }
  }
}

func (m *Vis) FName_2_FNum( full_fname string, file_num *int ) bool {

  var found bool = false

  for k:=0; !found && k<m.files.Len(); k++ {

    if( full_fname == m.files.GetPFb( k ).path_name ) {
      found = true
      *file_num = k
    }
  }
  return found
}

func (m *Vis) Exe_Colon_detab() {

  if( 6 < m_rbuf.Len() ) {
    S := m_rbuf.to_str()
    if tab_sz, err := strconv.Atoi( S[6:] ); err != nil {
      m.CmdLineMessage( fmt.Sprintf("Could not convert to int: %s", S[6:]) )
    } else {
      if( 0 < tab_sz && tab_sz <= 32 ) {
        m.CV().p_fb.RemoveTabs_SpacesAtEOLs( tab_sz )
      }
    }
  }
}

func (m *Vis) Exe_Colon_dos2unix() {

  m.CV().p_fb.dos2unix()
}

func (m *Vis) Exe_Colon_unix2dos() {

  m.CV().p_fb.unix2dos()
}

func (m *Vis) Exe_Colon_b() {

  if( 1 == m_rbuf.Len() ) { // :b
    m.GoToPrevBuffer()

  } else if( 2 <= m_rbuf.Len() ) {
    var r1 rune = m_rbuf.GetR(1)
    if       ( '#' == r1 ) { m.GoToPoundBuffer(); // :b#
    } else if( 'c' == r1 ) { m.GoToCurrBuffer();  // :bc
    } else if( 'e' == r1 ) { m.GoToBufferEditor();// :be
    } else if( 'm' == r1 ) { m.GoToMsgBuffer();   // :bm
    } else {                                      // :b<number>
      if buffer_num,err := strconv.Atoi( string(m_rbuf.data[1:]) ); err == nil {
        m.GoToBuffer( buffer_num )
      }
    }
  }
}

func (m *Vis) Set_Syntax() {

  if( 4 < m_rbuf.Len() ) {
    S := m_rbuf.to_str()
    m.CV().p_fb.Set_File_Type( S[4:] )
  }
}

func (m *Vis) Quit() {

  if( m.num_wins <= 1 ) {
    m.QuitAll()
  } else {
    m.Quit_One()
  }
}

func (m *Vis) QuitAll() {

  m.running = false
}

func (m *Vis) Quit_One() {

  var p_cv *FileView = m.CV()

  var TP Tile_Pos = p_cv.GetTilePos()

  if( p_cv.in_diff_mode ) {
  // FIXME:
  //NoDiff_4_FileBuf( m, m.diff.GetViewShort()->GetFB() )
  //NoDiff_4_FileBuf( m, m.diff.GetViewLong() ->GetFB() )

  //m.diff.Copy_DiffContext_2_Remaining_ViewContext()
  }
  if( m.win < m.num_wins-1 ) {
    m.Quit_ShiftDown()
  }
  if( 0 < m.win ) { m.win--; }
  m.num_wins--

  m.Quit_JoinTiles( TP )

  m.UpdateViews( false )

  m.CV().PrintCursor()
}

func (m *Vis) Quit_ShiftDown() {

  // Make copy of win's list of views and view history:
//ViewList win_views    ( m.views    [m.win] )
// unsList win_file_hist( m.file_hist[m.win] )

  // Make copy of win's list of views and view history:
//var win_views FileViewList; win_views.Copy( &m.views[m.win] )
//var win_file_hist  IntList; win_file_hist.Copy( &m.file_hist[m.win] )

  // Make copy of win's list of views and view history:
  var win_views FileViewList; win_views     = m.views[m.win]
  var win_file_hist  IntList; win_file_hist = m.file_hist[m.win]

  // Shift everything down
  for w:=m.win+1; w<m.num_wins; w++ {
    m.views    [w-1] = m.views    [w]
    m.file_hist[w-1] = m.file_hist[w]
  }
  // Put win's list of views at end of views:
  // Put win's view history at end of view historys:
  m.views    [m.num_wins-1] = win_views
  m.file_hist[m.num_wins-1] = win_file_hist
}

func (m *Vis) Quit_JoinTiles( TP Tile_Pos ) {

  // win is disappearing, so move its screen space to another view:
  if       ( TP == TP_LEFT_HALF )         { m.Quit_JoinTiles_LEFT_HALF()
  } else if( TP == TP_RITE_HALF )         { m.Quit_JoinTiles_RITE_HALF()
  } else if( TP == TP_TOP__HALF )         { m.Quit_JoinTiles_TOP__HALF()
  } else if( TP == TP_BOT__HALF )         { m.Quit_JoinTiles_BOT__HALF()
  } else if( TP == TP_TOP__LEFT_QTR )     { m.Quit_JoinTiles_TOP__LEFT_QTR()
  } else if( TP == TP_TOP__RITE_QTR )     { m.Quit_JoinTiles_TOP__RITE_QTR()
  } else if( TP == TP_BOT__LEFT_QTR )     { m.Quit_JoinTiles_BOT__LEFT_QTR()
  } else if( TP == TP_BOT__RITE_QTR )     { m.Quit_JoinTiles_BOT__RITE_QTR()
  } else if( TP == TP_LEFT_QTR )          { m.Quit_JoinTiles_LEFT_QTR()
  } else if( TP == TP_RITE_QTR )          { m.Quit_JoinTiles_RITE_QTR()
  } else if( TP == TP_LEFT_CTR__QTR )     { m.Quit_JoinTiles_LEFT_CTR__QTR()
  } else if( TP == TP_RITE_CTR__QTR )     { m.Quit_JoinTiles_RITE_CTR__QTR()
  } else if( TP == TP_TOP__LEFT_8TH )     { m.Quit_JoinTiles_TOP__LEFT_8TH()
  } else if( TP == TP_TOP__RITE_8TH )     { m.Quit_JoinTiles_TOP__RITE_8TH()
  } else if( TP == TP_TOP__LEFT_CTR_8TH ) { m.Quit_JoinTiles_TOP__LEFT_CTR_8TH()
  } else if( TP == TP_TOP__RITE_CTR_8TH ) { m.Quit_JoinTiles_TOP__RITE_CTR_8TH()
  } else if( TP == TP_BOT__LEFT_8TH )     { m.Quit_JoinTiles_BOT__LEFT_8TH()
  } else if( TP == TP_BOT__RITE_8TH )     { m.Quit_JoinTiles_BOT__RITE_8TH()
  } else if( TP == TP_BOT__LEFT_CTR_8TH ) { m.Quit_JoinTiles_BOT__LEFT_CTR_8TH()
  } else if( TP == TP_BOT__RITE_CTR_8TH ) { m.Quit_JoinTiles_BOT__RITE_CTR_8TH()
  } else if( TP == TP_LEFT_THIRD )        { m.Quit_JoinTiles_LEFT_THIRD()
  } else if( TP == TP_CTR__THIRD )        { m.Quit_JoinTiles_CTR__THIRD()
  } else if( TP == TP_RITE_THIRD )        { m.Quit_JoinTiles_RITE_THIRD()
  } else if( TP == TP_LEFT_TWO_THIRDS )   { m.Quit_JoinTiles_LEFT_TWO_THIRDS()
  } else if( TP == TP_RITE_TWO_THIRDS )   { m.Quit_JoinTiles_RITE_TWO_THIRDS()
  }
}

func ( m *Vis ) Exe_Colon_e() {

  var p_cv *FileView = m.CV()

  if( 1 == m_rbuf.Len() ) { // :e
    var p_fb *FileBuf = p_cv.p_fb
    p_fb.ReReadFile()

    for w:=0; w<m.num_wins; w++ {
      if( p_fb == m.GetView_Win( w ).p_fb ) {
        // View is currently displayed, perform needed update:
        m.GetView_Win( w ).Update_and_PrintCursor()
      }
    }
  } else { // :e file_name
    // Edit file of supplied file name:
    var fname string = m_rbuf.to_str()[1:]

    pname := FindFullFileNameRel2( p_cv.p_fb.dir_name, fname )

    m.NotHaveFileAddFile( pname )

    var file_index int = 0

    if( m.HaveFile( pname, &file_index ) ) {
      m.GoToBuffer( file_index )
    }
  }
}

func ( m *Vis ) Exe_Colon_w() {
//Log("Exe_Colon_w()")

  var p_cv *FileView = m.CV()

  var file_written bool = false

  if( m_rbuf.EqualStr("w") || m_rbuf.EqualStr("wq") ) {

    if( p_cv == m.views[ m.win ].GetPFv( m_SHELL_FILE ) ) {
      // Dont allow SHELL_BUFFER to be saved with :w.
      // Require :w filename.
      p_cv.PrintCursor()

    } else {
      // If the file gets written, CmdLineMessage will be called,
      // which will put the cursor back in position,
      // else Window_Message will be called
      // which will put the cursor back in the message window
      file_written = p_cv.p_fb.Write()
    }
    if( m_rbuf.EqualStr("wq") ) {
      m.Quit()
    }
  } else { // :w file_name
    // Write file of supplied file name:
    var fname string = m_rbuf.to_str()[1:]

    pname := FindFullFileNameRel2( p_cv.p_fb.dir_name, fname )

    var file_index int = -1
    if( m.HaveFile( pname, &file_index ) ) {
      m.files.GetPFb( file_index ).Write()

    } else if( DIR_DELIM != pname[ len( pname )-1 ] ) {
    //p_fb := m.CreateFile( pname, FT_UNKNOWN )
      p_fb := new( FileBuf )
      p_fb.Init_FB_2( pname, FT_UNKNOWN, p_cv.p_fb )

      file_written = p_fb.Write()
    }
  }
  if( file_written ) {
    if( m.Update_Change_Statuses() ) {
    //if( p_cv.m_in_diff ) {
    //  m_diff.PrintCursor(); // Does m_console.Update()
    //}
    }
  }
}

func ( m *Vis ) MoveToLine() {
  // Move cursor to line:
  line_num,err := strconv.Atoi( m_rbuf.to_str() )

  if( nil == err ) {
    var p_cv *FileView = m.CV()

    p_cv.GoToLine( line_num )
  }
}

func ( m *Vis ) Handle_Colon_Cmd() {

  m_rbuf.RemoveSpaces()

  if( 0 == m_rbuf.Len() ) {
    m.CV().PrintCursor()
  } else {
    if       ( m_rbuf.EqualStr("q") )        { m.Quit()
    } else if( m_rbuf.EqualStr("qa") )       { m.QuitAll()
    } else if( m_rbuf.EqualStr("vsp") )      { m.VSplitWindow()
    } else if( m_rbuf.EqualStr("sp") )       { m.HSplitWindow()
    } else if( IsDigit(m_rbuf.GetR(0)) )     { m.MoveToLine()
    } else if( m_rbuf.GetR(0)=='b' )         { m.Exe_Colon_b()
    } else if( m_rbuf.GetR(0)=='e' )         { m.Exe_Colon_e()
    } else if( m_rbuf.GetR(0)=='w' )         { m.Exe_Colon_w()
    } else if( m_rbuf.EqualStr("dos2unix") ) { m.Exe_Colon_dos2unix()
    } else if( m_rbuf.EqualStr("unix2dos") ) { m.Exe_Colon_unix2dos()
    } else if( m_rbuf.StartsWith("syn=") )   { m.Set_Syntax()
    } else if( m_rbuf.StartsWith("detab=") ) { m.Exe_Colon_detab()
    } else {
      m.CV().PrintCursor()
    }
  }
}

//func ( m *Vis ) Handle_Slash_GotPattern( pattern string ) {
//
//  m.regex_str = pattern
//
//  if( 0<len(m.regex_str) ) {
//    m.Do_Star_Update_Search_Editor()
//
//    var p_cv *FileView = m.CV()
//
//    if( p_cv.in_diff_mode ) { m.diff.Do_n()
//    } else                  {   p_cv.Do_n()
//    }
//  }
//  // Show new slash pattern for all windows currently displayed:
//  m.UpdateViews( true )
//}

func ( m *Vis ) Handle_Slash_GotPattern( pattern string, goto_pattern bool ) {

  m.regex_str = pattern

  if( 0<len(m.regex_str) ) {
    m.Do_Star_Update_Search_Editor()

    if( goto_pattern ) {
      var p_cv *FileView = m.CV()
      if( p_cv.in_diff_mode ) { m.diff.Do_n()
      } else                  {   p_cv.Do_n()
      }
    }
  }
  // Show new slash pattern for all windows currently displayed:
  m.UpdateViews( true )
}

// 1. Search for regex pattern in search editor.
// 2. If regex pattern is found in search editor,
//         move pattern to end of search editor
//    else add regex pattern to end of search editor
// 3. If search editor is displayed, update search editor window
//
func ( m *Vis ) Do_Star_Update_Search_Editor() {
  var pfb *FileBuf = &m.slash_file

  // If last line in SLASH_BUFFER is blank, remove it:
  NUM_SE_LINES := pfb.NumLines(); // Number of search editor lines
  if( 0<NUM_SE_LINES && 0 == pfb.LineLen( NUM_SE_LINES-1 ) ) {
    pfb.RemoveLP( NUM_SE_LINES-1 )
    NUM_SE_LINES = pfb.NumLines()
  }
  // 1. Search for regex pattern in search editor.
  found_pattern_in_search_editor := false
  line_in_search_editor := 0

  for ln:=0; !found_pattern_in_search_editor && ln<NUM_SE_LINES; ln++ {
    lp := pfb.GetLP( ln )
    if( lp.EqualStr( m.regex_str ) ) {
      found_pattern_in_search_editor = true
      line_in_search_editor = ln
    }
  }
  // 2. If regex pattern is found in search editor,
  //         move pattern to end of search editor
  //    else add regex pattern to end of search editor
  if( found_pattern_in_search_editor ) {
    // Move pattern to end of search editor, so newest searches are at bottom of file
    if( line_in_search_editor < NUM_SE_LINES-1 ) {
      p_fl := pfb.RemoveLP( line_in_search_editor )
      pfb.PushLP( p_fl )
    }
  } else {
    // Push regex onto search editor buffer
    p_fl := new(FLine)
  //for( const char* p=m.regex.c_str(); *p; p++ ) line.push( *p )
    for _, R := range m.regex_str {
      p_fl.PushR( R )
    }
    pfb.PushLP( p_fl )
  }
  // Push an emtpy line onto slash buffer to leave empty / prompt:
  pfb.PushLE()

  // 3. If search editor is displayed, update search editor window
  p_cv := m.CV()
  m.slash_view.SetContext( p_cv.WinCols(), p_cv.X(), p_cv.Cmd__Line_Row() )
  m.slash_view.GoToCrsPos_NoWrite( pfb.NumLines()-1, 0 )
  pfb.Update()
}

func (m *Vis) Init() {

  m.num_wins = 1

  m.InitBufferEditor()
  m.InitHelpBuffer()
  m.InitMsgBuffer()
  m.InitShellBuffer()
  m.InitColonBuffer()
  m.InitSlashBuffer()

  m.InitUserFiles()
  m.InitFileHistory()
  m.InitViewFuncs()
  m.InitLineFuncs()
}

func (m *Vis) CV() *FileView {
  return m.views[m.win].GetPFv( m.file_hist[m.win].Get(0) )
}

func (m *Vis) PV() *FileView {
  return m.views[m.win].GetPFv( m.file_hist[m.win].Get(1) )
}

func (m *Vis) GetView_Win( w int ) *FileView {

  return m.views[w].GetPFv( m.file_hist[w].Get( 0 ) )
}

func (m *Vis) Buf2FileNum( p_fb *FileBuf ) int {

  for k:=0; k<m.views[0].Len(); k++ {
    if( m.views[0].GetPFv( k ).p_fb == p_fb ) {
      return k
    }
  }
  return 0
}

func (m *Vis) Curr_FileNum() int {

  return m.file_hist[ m.win ].Get( 0 )
}

func (m *Vis) UpdateViews( show_search bool ) {

  for w:=0; w<m.num_wins; w++ {

    var pv *FileView = m.GetView_Win( w )

    if( ! pv.in_diff_mode ) {
      if( show_search ) {
        var msg string = "/" + m.regex_str
        pv.Set_Cmd_Line_Msg( msg )
      }
      pv.Update_not_PrintCursor()
    }
  }
// FIXME:
//if( m.InDiffMode() ) {
//  if( show_search ) {
//    string msg("/")
//    m.diff.Set_Cmd_Line_Msg( msg += m.regex )
//  }
//  m.diff.Update()
//}
  var p_cv *FileView = m.CV()

  if( !p_cv.in_diff_mode ) { p_cv.PrintCursor(); }
}

func (m *Vis) UpdateViewsOfFile( p_fb *FileBuf ) {

  // Update displayed views of file referred to by fb:
  for w:=0; w<m.num_wins; w++ {
    // V is currently displayed view in pane w:
    var V *FileView = m.GetView_Win( w )

    if( V.p_fb == p_fb ) {
      // View V is of fb, so update:
      V.Update_not_PrintCursor()
    }
  }
}

func (m *Vis) Is_BE_FILE( p_fb *FileBuf ) bool {

  return p_fb == m.views[0].GetPFv( m_BE_FILE ).p_fb
}

func (m *Vis) Handle_f() {

  kr := m_key.In()

  if( kr.IsKeyRune() ) {
    m.fast_rune = kr.R

    cv := m.CV()

    if( cv.in_diff_mode ) { m.diff.Do_f( m.fast_rune )
    } else                {     cv.Do_f( m.fast_rune )
    }
  }
}

func (m *Vis) Handle_SemiColon() {

  if( 0 <= m.fast_rune ) {
    cv := m.CV()

    if( cv.in_diff_mode ) { m.diff.Do_f( m.fast_rune )
    } else                {     cv.Do_f( m.fast_rune )
    }
  }
}

func (m *Vis) Handle_z() {

  cv := m.CV()
  kr := m_key.In()

  if( kr.R == 't' || kr.IsEndOfLineDelim() ) {
    if( cv.in_diff_mode ) { m.diff.MoveCurrLineToTop()
    } else                {    cv.MoveCurrLineToTop()
    }
  } else if( kr.R == 'z' ) {
    if( cv.in_diff_mode ) { m.diff.MoveCurrLineCenter()
    } else                {    cv.MoveCurrLineCenter()
    }
  } else if( kr.R == 'b' ) {
    if( cv.in_diff_mode ) { m.diff.MoveCurrLineToBottom()
    } else               {     cv.MoveCurrLineToBottom()
    }
  }
}

func (m *Vis) Handle_Resize() {

  m.AdjustViews()

//m.CV().Update_and_PrintCursor()
  m.UpdateViews( false )
}

func (m *Vis) AdjustViews() {

  for w:=0; w<m.num_wins; w++ {

    m.GetView_Win( w ).SetViewPos()
  }
}

func (m *Vis) CmdLineMessage( msg string ) {

  var pV *FileView = m.CV()

  WC  := pV.WorkingCols()
  ROW := pV.Cmd__Line_Row()
  COL := pV.Col_Win_2_GL( 0 )
  MSG_LEN := len( msg )

  if( WC < MSG_LEN ) {
    // messaged does not fit, so truncate beginning
    m_console.SetString( ROW, COL, msg[MSG_LEN-WC:], &TS_NORMAL )
  } else {
    // messaged fits, add spaces at end
    m_console.SetString( ROW, COL, msg, &TS_NORMAL )
    for k:=0; k<(WC-MSG_LEN); k++ {
      m_console.SetR( ROW, COL+MSG_LEN+k, ' ', &TS_NORMAL )
    }
  }

  if( pV.in_diff_mode ) { m.diff.PrintCursor( pV )
  } else                {     pV.PrintCursor()
  }
}

func (m *Vis) Window_Message( msg string ) {
  // FIXME:
}

func (m *Vis) Speed_up_scrolling( _ru rune, _handler CmdFunc2 ) (bool, Key_rune) {
  var l_have_saved_key_ru bool = false
  var l_saved_kr Key_rune

  var num int = 1
  var done bool = false
  for !done {
    if( !m_console.HasPendingEvent() ) {
      done = true
    } else {
      l_kr := m_key.In()
      if( l_kr.IsKeyRune() && l_kr.R == _ru ) {
        num++
      } else {
        done = true
        l_have_saved_key_ru = true
        l_saved_kr = l_kr
      }
    }
  }
  // Fast scrolling:
  if( 1<num ) { _handler( m, num*2 )
  } else      { _handler( m, num )
  }
  return l_have_saved_key_ru, l_saved_kr
}

func (m *Vis) CheckFileModTime() {
  // m.file_hist[m.win].Get(0) is the current file number of the current window
  if( m_USER_FILE <= m.file_hist[m.win].Get(0) ) {
    pfb := m.CV().p_fb

    var curr_mod_time time.Time = ModificationTime( pfb.path_name )

    if(  curr_mod_time.After(pfb.mod_time) ) {

      if( pfb.is_regular ) {
        // Update file modification time so that the message window
        // will not keep popping up:
        pfb.changed_externally = true

      } else if( pfb.is_dir ) {
        // Dont ask the user, just read in the directory.
        // pfb->GetModTime() will get updated in pfb->ReReadFile()
        pfb.ReReadFile()

        for w:=0; w<m.num_wins; w++ {
          pV := m.GetView_Win( w )
          if( pfb == pV.p_fb ) {
            // View is currently displayed, perform needed update:
            pV.Update_not_PrintCursor()
          }
        }
      }
    }
  }
}

// This ensures that proper change status is displayed around each window:
// '    ' for file in vis same as file on file system,
// '++++' for changes in vis not written to file system,
// '////' for file on file system changed externally to vis,
// '+/+/' for changes in vis and on file system
func (m *Vis) Update_Change_Statuses() bool {

  // Update buffer changed status around windows:
  updated_change_sts := false

  for w:=0; w<m.num_wins; w++ {
    // pV points to currently displayed view in window w:
    var pV *FileView = m.GetView_Win( w )

    if( (pV.un_saved_change_sts != pV.p_fb.Changed()) ||
        (pV.external_change_sts != pV.p_fb.changed_externally) ) {

      pV.PrintBorders()

      pV.un_saved_change_sts = pV.p_fb.Changed()
      pV.external_change_sts = pV.p_fb.changed_externally

      updated_change_sts = true
    }
  }
  if( updated_change_sts ) { m_console.Show()
  }
  return updated_change_sts
}

func (m *Vis) GetCmdFunc( kr Key_rune ) CmdFunc {
  var cf CmdFunc
  if( kr.IsKeyRune() ) { cf = m.view_funcs[ kr.R ]
  } else               { cf = m.view_funcs[ kr.K ]
  }
  return cf
}

func (m *Vis) GetLineFunc( kr Key_rune ) CmdFunc {
  var cf CmdFunc
  if( kr.IsKeyRune() ) { cf = m.line_funcs[ kr.R ]
  } else               { cf = m.line_funcs[ kr.K ]
  }
  return cf
}

func (m *Vis) Run() {
  var have_saved_key_ru bool = false
  var saved_kr Key_rune

  var kr Key_rune

  m.running = true
  for m.running {
    if( have_saved_key_ru ) {
      kr = saved_kr
      have_saved_key_ru = false
    } else {
      kr = m_key.In()
    }

    if( kr.IsKeyRune() ) {
      // kr.R is valid
      var cf CmdFunc
      if( m.colon_mode || m.slash_mode ) {
        cf = m.line_funcs[ kr.R ]
      } else {
        if( kr.R == 'j' ) {
          have_saved_key_ru, saved_kr = m.Speed_up_scrolling( kr.R, Handle_j )
        } else if( kr.R == 'k' ) {
          have_saved_key_ru, saved_kr = m.Speed_up_scrolling( kr.R, Handle_k )
        } else {
          cf = m.view_funcs[ kr.R ]
        }
      }
      if nil != cf {
        cf(m)
      }
    } else {
      // kr.R is NOT valid. kr.K is tcell.KeyESC, tcell.KeyLF or tcell.KeyCR
      var cf CmdFunc
      if( m.colon_mode || m.slash_mode ) { cf = m.line_funcs[ kr.K ]
      } else                             { cf = m.view_funcs[ kr.K ]
      }
      if nil != cf {
        cf(m)
      }
    }
    m.CheckFileModTime()
  //m.Update_Change_Statuses()
  //updated_chg_sts := m.Update_Change_Statuses()
  }
}

