
package main

import (
  "bufio"
  "fmt"
  "io"
  "io/fs"
  "log"
  "os"
//"sort"
  "regexp"
  "strings"
  "time"
)

type FileBuf struct {
//p_vis     *Vis
  p_lv      *LineView

  path_name, dir_name string
//file_name string

  is_dir, is_regular bool
  lines     FLineList
  views     FileViewList

  changed_externally bool

  file_type       File_Type
  Hi              Highlight_Base
//hi_touched_line int

  regex_str string
  p_regex_obj *regexp.Regexp

  history ChangeHist
  save_history bool

  mod_time time.Time
}

func (m *FileBuf) Init_FB_Common( path_name string, FT File_Type ) {
  m.path_name = path_name

  m.is_dir     = IsDir( m.path_name )
  m.is_regular = IsReg( m.path_name )

  if( m.is_dir ) {
    m.dir_name = m.path_name
  } else {
    m.dir_name = GetFnameTail( m.path_name )
  }
  if( FT == FT_BUFFER_EDITOR ) {
    m.file_type = FT_BUFFER_EDITOR;
    m.Hi = new( Highlight_BufferEditor )
    m.Hi.Init( m )

  } else if( FT_UNKNOWN == m.file_type ) {
    m.Find_File_Type_Suffix();
  }
  m_vis.Add_FileBuf_2_Lists_Create_Views( m )
  m.history.Init( m )
}

func (m *FileBuf) Init_FB( path_name string, FT File_Type ) {

  m.Init_FB_Common( path_name, FT )

  if( m_USER_FILE <= m_vis.Buf2FileNum( m ) ) {
    m.ReadFile()
    m.save_history = true
  }
}

func (m *FileBuf) Init_FB_2( path_name string, FT File_Type, p_other *FileBuf ) {

  m.Init_FB_Common( path_name, FT )

  m.lines.CopyP( &p_other.lines )
//m.LF_at_EOF = p_other.LF_at_EOF
}

func (m *FileBuf) Find_File_Type_CPP() bool {

  if( strings.HasSuffix(m.path_name, ".h") ||
      strings.HasSuffix(m.path_name, ".h.new") ||
      strings.HasSuffix(m.path_name, ".h.old") ||
      strings.HasSuffix(m.path_name, ".c") ||
      strings.HasSuffix(m.path_name, ".c.new") ||
      strings.HasSuffix(m.path_name, ".c.old") ||
      strings.HasSuffix(m.path_name, ".hh") ||
      strings.HasSuffix(m.path_name, ".hh.new") ||
      strings.HasSuffix(m.path_name, ".hh.old") ||
      strings.HasSuffix(m.path_name, ".cc") ||
      strings.HasSuffix(m.path_name, ".cc.new") ||
      strings.HasSuffix(m.path_name, ".cc.old") ||
      strings.HasSuffix(m.path_name, ".hpp") ||
      strings.HasSuffix(m.path_name, ".hpp.new") ||
      strings.HasSuffix(m.path_name, ".hpp.old") ||
      strings.HasSuffix(m.path_name, ".cpp") ||
      strings.HasSuffix(m.path_name, ".cpp.new") ||
      strings.HasSuffix(m.path_name, ".cpp.old") ||
      strings.HasSuffix(m.path_name, ".cxx") ||
      strings.HasSuffix(m.path_name, ".cxx.new") ||
      strings.HasSuffix(m.path_name, ".cxx.old") ) {

    m.file_type = FT_CPP;
    m.Hi = new( Highlight_CPP )
    m.Hi.Init( m )
    return true;
  }
  return false;
}

func (m *FileBuf) Find_File_Type_Go() bool {

  if( strings.HasSuffix(m.path_name, ".go") ||
      strings.HasSuffix(m.path_name, ".go.new") ||
      strings.HasSuffix(m.path_name, ".go.old") ) {

    m.file_type = FT_GO;
    m.Hi = new( Highlight_Go )
    m.Hi.Init( m )
    return true;
  }
  return false;
}

func (m *FileBuf) Find_File_Type_Suffix() {

  if( m.is_dir ) {
    m.file_type = FT_DIR;
    m.Hi = new( Highlight_Dir )
    m.Hi.Init( m )

  } else if( m.Find_File_Type_CPP() ||
             m.Find_File_Type_Go() ) {
    // File type found
  } else {
    // File type NOT found based on suffix.
    // File type will be found in Find_File_Type_FirstLine()
  }
}

func (m *FileBuf) AddView( p_fv *FileView ) {

  m.views.PushPFv( p_fv )
}

func (m *FileBuf) ReadFile() {
  if       ( m.is_dir     ) { m.ReadExistingDir( m.path_name )
  } else if( m.is_regular ) { m.ReadExistingFile( m.path_name )
  } else {
    // File does not exist, so add an empty line:
    m.PushLE()
  }
  m.mod_time = ModificationTime( m.path_name )
}

func (m *FileBuf) ReReadFile() {

  // Can only re-read user files
  if( m_USER_FILE <= m_vis.Buf2FileNum( m ) ) {
    m.ClearChanged();
    m.ClearLines();

    m.save_history = false; //< Gets turned back on in ReadFile()

    m.ReadFile();

    // Reposition cursor in each FileView of this file if needed:
    for w:=0; w<MAX_WINS; w++ {
      var p_V *FileView = m.views.GetPFv( w );

      p_V.Check_Context();
    }
    m.save_history    = true;
  }
}

func (m *FileBuf) ReadExistingDir( dir_path string ) {

  var dir_fs fs.FS = os.DirFS( dir_path )

  var de_s []fs.DirEntry
  var err  error

  de_s, err = fs.ReadDir( dir_fs, "." )

  if err != nil {
    log.Fatal( err )
  } else {
    m.PushLSR( []rune("..") )
    for _, de := range de_s {
      var de_name string = de.Name()
      if de.IsDir() {
        de_name = AppendDirDelim( de_name )
      }
      m.PushLSR( []rune(de_name) )
    }
  }
  m.ReadExistingDir_Sort()
}

func (m *FileBuf) ReadExistingDir_Sort() {

  var NUM_LINES int = m.NumLines()

  // Sort lines (file names), least to greatest:
  for i:=NUM_LINES-1; 0<i; i-- {
    for k:=0; k<i; k++ {

      var p_l_0 *FLine = m.lines.GetLP( k   )
      var p_l_1 *FLine = m.lines.GetLP( k+1 )

      // *p_l_0 is greater then *p_l_1:
      if( 1 < p_l_0.Compare( p_l_1 ) ) { m.lines.Swap( k, k+1 )
      }
    }
  }

  // Move non-directory files to end:
  for i:=NUM_LINES-1; 0<i; i-- {
    for k:=0; k<i; k++ {

      var p_l_0 *FLine = m.lines.GetLP( k   )

      if( ! p_l_0.EqualStr("..") ) {
        var p_l_1 *FLine = m.lines.GetLP( k+1 )

        S_DIR_DELIM := string(DIR_DELIM)

        if( !p_l_0.ends_with( S_DIR_DELIM ) &&
             p_l_1.ends_with( S_DIR_DELIM ) ) { m.lines.Swap( k, k+1 )
        }
      }
    }
  }
}

//func (m *FileBuf) ReadExistingFile( file_path string ) {
//
//  if( !FileExists( file_path ) ) {
//    // File does not exist, so just add an empty line:
//    m.PushLE()
//  } else {
//    infile_bytes, err := os.ReadFile( file_path )
//    if err != nil {
//      // File exists but could not be read.  Drop out.
//      m_vis.CmdLineMessage( fmt.Sprintf("\"%s\" exists but could not be read", file_path) )
//    } else {
//      var p_fl *FLine
//      p_fl = m.add_infile_bytes( p_fl, infile_bytes )
//      if nil != p_fl {
//        m.lines.PushLP( p_fl )
//      }
//    }
//  }
//}

func (m *FileBuf) ReadExistingFile( file_path string ) {
  if( !FileExists( file_path ) ) {
    // File does not exist, so just add an empty line:
    m.PushLE()
  } else {
    p_f, err := os.Open( file_path )
    if( err != nil ) {
      m_vis.CmdLineMessage( fmt.Sprintf("\"%s\" exists but could not be read", file_path) )
    } else {
      defer p_f.Close()
      m.read_open_file_in_chunks( p_f )
    }
  }
}

func (m *FileBuf) read_open_file_in_chunks( p_f *os.File ) {
  var p_fl *FLine = nil
  infile_bytes := make([]byte, 128)
  var offset int64 = 0
  done := false
  for !done {
    n, err := p_f.ReadAt( infile_bytes, offset )
    if( n < 128 ) {
      done = true
      if( err != io.EOF ) {
        m_vis.CmdLineMessage( fmt.Sprintf("Error reading from: \"%s\" at %v: %v", p_f.Name(), offset, err ) )
      }
    }
    if( 0 < n ) {
      p_fl = m.add_infile_bytes( p_fl, infile_bytes[:n] )
      offset += int64(n)
    }
  }
  if nil != p_fl { m.lines.PushLP( p_fl )
  }
}

func (m *FileBuf) add_infile_bytes( p_fl *FLine, infile_bytes []byte ) *FLine {
  for _, B := range infile_bytes {
    if nil == p_fl { p_fl = new(FLine) }

    if '\n' == B {
      m.PushLP( p_fl )
      p_fl = nil
      m.lines.LF_at_EOF = true;
    } else {
      p_fl.PushR( rune(B) )
      m.lines.LF_at_EOF = false;
    }
  }
  return p_fl
}

func (m *FileBuf) ClearLines() {
  m.lines.Clear()

//m.lineRegexsValid.clear();
}

func (m *FileBuf) Write() bool {
  var ok bool = false;

//if( ENC_BYTE == m.encoding )
//{
    ok = m.Write_p( &m.lines );
//}
//else if( ENC_HEX == m.encoding )
//{
//  Array_t<BLine*> n_lines;
//  bool LF_at_EOF = true;
//
//  if( HEX_to_BYTE_get_lines( m, n_lines, LF_at_EOF ) )
//  {
//    ok = Write_p( m, n_lines, LF_at_EOF );
//  }
//}
//else {
//  m.p_vis.Window_Message("\nUnhandled Encoding: %s\n\n"
//                        , Encoding_Str( m.encoding ) );
//}
  return ok;
}

func (m *FileBuf) Write_p( p_lines *FLineList ) bool {

  var ok bool = true

  if( 0==len( m.path_name ) ) {
    // No file name message:
    ok = false;
    m_vis.CmdLineMessage("No file name to write to");
  } else {
    p_f, err := os.OpenFile( m.path_name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644 )

    if( err != nil ) {
      // Could not open file for writing message:
      ok = false;
      m_vis.Window_Message( fmt.Sprintf("\nCould not open:\n\n%s\n\nfor writing\n\n",
                            m.path_name) )
    } else {
      defer p_f.Close()
      var NUM_LINES int = p_lines.Len()

      for k:=0; err==nil && k<NUM_LINES; k++ {
        _,err = p_f.Write( p_lines.GetLP(k).to_SB(0) )

        if( k<NUM_LINES-1 || p_lines.LF_at_EOF ) {
          _,err = p_f.Write( []byte("\n") );
        }
      }
      if( err != nil ) {
        m_vis.CmdLineMessage( fmt.Sprintf("Failed to write: \"%s\"", m.path_name) );
      } else {
        m.mod_time = ModificationTime( m.path_name )
        m.changed_externally = false;

        m.history.Clear();
        // Wrote to file message:
        m_vis.CmdLineMessage( fmt.Sprintf("\"%s\" written", m.path_name) );
      }
    }
  }
  return ok;
}

func (m *FileBuf) GetPath() string {
  return m.path_name
}

func (m *FileBuf) NumLines() int {
  return m.lines.Len()
}

func (m *FileBuf) LineLen( k int ) int {
  if 0 <= k && k < m.lines.Len() {
    return m.lines.LineLen( k )
  }
  return 0
}

func (m *FileBuf) GetR( l_num, r_num int ) rune {

  return m.lines.GetR( l_num, r_num )
}

func (m *FileBuf) SetR( l_num, r_num int, R rune, continue_last_update bool ) {

  old_R := m.lines.GetR( l_num, r_num )

  if( old_R != R ) {
    m.lines.SetR( l_num, r_num, R )
    if( m.save_history ) {
      m.history.Save_Set( l_num, r_num, old_R, continue_last_update )
    }
  }
}

//func (m *FileBuf) PushL( ln FLine ) {
//
//  m.lines.PushLP( &ln )
//}

// Push pointer to FLine
func (m *FileBuf) PushLP( p_fl *FLine ) {

  m.lines.PushLP( p_fl )

  if( m.save_history ) {
    m.history.Save_InsertLine( m.lines.Len()-1 )
  }
}

// Push Line Empty
func (m *FileBuf) PushLE() {
  p_fl := new(FLine)
  m.PushLP( p_fl )
}

// Append Line made from a Slice of runes
func (m *FileBuf) PushLSR( s_r []rune ) {
  p_fl := new(FLine)
  p_fl.PushSR( s_r )
  m.PushLP( p_fl )
}

// Append line ln to end of line l_num
//
func (m *FileBuf) AppendLineToLine( l_num int, p_fl *FLine ) {

  m.lines.AppendLineToLine( l_num, p_fl )

  if( m.save_history ) {
    NEW_LL := p_fl.Len()
    first_insert := m.lines.LineLen( l_num ) - NEW_LL
    for k:=0; k<NEW_LL; k++ {
      m.history.Save_InsertRune( l_num, first_insert + k )
    }
  }
}

func (m *FileBuf) InsertR( l_num, r_num int, R rune ) {

  m.lines.InsertR( l_num, r_num, R )

  if( m.save_history ) {
    m.history.Save_InsertRune( l_num, r_num )
  }
}

func (m *FileBuf) PushR( l_num int, R rune ) {

  m.lines.PushR( l_num, R )

  if( m.save_history ) {
    pushed_pos := m.lines.LineLen( l_num )-1
    m.history.Save_InsertRune( l_num, pushed_pos )
  }
}

func (m *FileBuf) GetLP( l_num int ) *FLine {

  var p_fl *FLine = m.lines.GetLP( l_num )
  return p_fl
}

//func (m *FileBuf) GetLP( l_num int ) *FLine {
//
//  var p_ln *FLine = m.lines.GetLP( l_num )
//  return p_ln
//}

func (m *FileBuf) RemoveLP( l_num int ) *FLine {

  p_fl := m.lines.RemoveLP( l_num )

  if( m.save_history ) {
    m.history.Save_RemoveLine( l_num, p_fl )
  }
  return p_fl
}

// Remove from FileBuf and return the byte at line l_num and position c_num
//
func (m *FileBuf) RemoveR( l_num, c_num int ) rune {

//Line* lp =  m.lines[ l_num ];
//Line* sp = m.styles[ l_num ];

  var R rune = m.lines.RemoveR( l_num, c_num )
//byte C = 0;
//bool ok = lp->remove( c_num, C )
       //&& sp->remove( c_num )
       //&& m.lineRegexsValid.set( l_num, false );

//if( SavingHist( m ) ) m.history.Save_RemoveChar( l_num, c_num, C );

  if( m.save_history ) {
    m.history.Save_RemoveRune( l_num, c_num, R )
  }
  return R;
}

// Insert a new empty line on line l_num.
// l_num can be len( m.lines ).
//
func (m *FileBuf) InsertLE( l_num int ) {

  m.lines.InsertLE( l_num )

  if( m.save_history ) {
    m.history.Save_InsertLine( l_num )
  }
  m.InsertLine_Adjust_Views_topLines( l_num );
}

// Insert RLine
func (m *FileBuf) InsertRLP( l_num int, p_rl *RLine ) {

  m.lines.InsertRLP( l_num, p_rl )

  if( m.save_history ) {
    m.history.Save_InsertLine( l_num )
  }
  m.InsertLine_Adjust_Views_topLines( l_num );
}

func (m *FileBuf) GetSize() int {

  return m.lines.GetSize()
}

func (m *FileBuf) GetCursorByte( CL, CC int ) int {

  return m.lines.GetCursorByte( CL, CC )
}

//func (m *FileBuf) UpdateWinViews( PRINT_CMD_LINE bool ) {
//
//  for w:=0; w<MAX_WINS; w++ {
//
//    pV *View = m.views[w];
//
//    for w2:=0; w2<m_vis.GetNumWins(); w2++ {
//
//      if( pV == m_vis.WinView( w2 ) )
//      {
//      //m.self.Find_Styles( pV->GetTopLine() + pV->WorkingRows() );
//      //m.self.Find_Regexs( pV->GetTopLine(), pV->WorkingRows() );
//
//        pV.RepositionView();
//        pV.Print_Borders();
//        pV.PrintWorkingView();
//        pV.PrintStsLine();
//        pV.PrintFileLine();
//
//        if( PRINT_CMD_LINE ) pV.PrintCmdLine();
//
//        pV->SetStsLineNeedsUpdate( true );
//      }
//    }
//  }
//if m.p_fv == m_vis.CV() {
//  // m.p_fv is currently displayed, so update it:
//  m.p_fv.Update_not_PrintCursor()
//}
//}

func (m *FileBuf) Update() {

  m_vis.Update_Change_Statuses()

  m_vis.UpdateViewsOfFile( m )

  // Put cursor back into current window
  m_vis.CV().PrintCursor()
}

func (m *FileBuf) UpdateCmd() {

//m_vis.Update_Change_Statuses();

//m_vis.UpdateViewsOfFile( m );

  if( nil != m.p_lv ) {

    m.p_lv.RepositionView();
    m.p_lv.PrintWorkingView();

    // Put cursor back into current window
    m.p_lv.PrintCursor();
  }
}

func (m *FileBuf) InsertLine_Adjust_Views_topLines( l_num int ) {

  for w:=0; w<m.views.Len(); w++ {

    var p_fv *FileView = m.views.GetPFv( w );

    p_fv.InsertedLine_Adjust_TopLine( l_num );
  }
}

func (m *FileBuf) Update_Styles_Find_St( first_line int ) CrsPos {
  // Default start position is beginning of file
  st := CrsPos{ 0, 0 }

  // Find first position without a style before first_line, because the
  // style code assumes it is finding styles starting on an empty style
  // CrsPos:
  var done bool = false
  for l:=first_line-1; !done && 0<=l; l-- {
    var LL int = m.LineLen( l );
    for p:=LL-1; !done && 0<=p; p-- {
      var S byte = m.lines.GetLP(l).GetStyle(p)
      if( 0==S ) {
        st.crsLine = l;
        st.crsChar = p;
        done = true;
      }
    }
  }
  return st;
}

// Find m.styles up to but not including up_to_line number
func (m *FileBuf) Find_Styles( up_to_line int ) {

  var NUM_LINES int = m.NumLines()

  if( 0<NUM_LINES ) {
    m.lines.hi_touched_line = Min_i( m.lines.hi_touched_line, NUM_LINES-1 );

    if( m.lines.hi_touched_line < up_to_line ) {
      // Find m.styles for some EXTRA_LINES beyond where we need to find
      // m.styles for the moment, so that when the user is scrolling down
      // through an area of a file that has not yet been syntax highlighed,
      // Find_Styles_In_Range() does not need to be called every time the
      // user scrolls down another line.  Find_Styles_In_Range() will only
      // be called once for every EXTRA_LINES scrolled down.
      var EXTRA_LINES int = 10;

      var st CrsPos = m.Update_Styles_Find_St( m.lines.hi_touched_line );
      var fn int    = Min_i( up_to_line+EXTRA_LINES, NUM_LINES );

      m.Find_Styles_In_Range( st, fn );

      m.lines.hi_touched_line = fn;
    }
  }
}

// Find m.styles starting at st up to but not including fn line number
func (m *FileBuf) Find_Styles_In_Range( st CrsPos, fn int ) {
  // Hi should have already have been allocated, but just in case
  // check and allocate here if needed.
  if( nil == m.Hi ) {
    m.file_type = FT_TEXT;
    m.Hi = new( Highlight_Text )
    m.Hi.Init( m )
  }
  m.Hi.Run_Range( st, fn );
}

// Leave star style unchanged, and clear syntax m.styles
func (m *FileBuf) ClearSyntaxStyles( l_num, c_num int ) {

//var sr BLine = m.lines.GetLP( l_num ).styles
  var lp *FLine = m.lines.GetLP( l_num )

  // Clear everything except star and in-file
//sr.SetB( c_num, sr.GetB( c_num ) & ( HI_STAR | HI_STAR_IN_F ) );
  lp.SetStyle( c_num, lp.GetStyle( c_num ) & ( HI_STAR | HI_STAR_IN_F ) );
}

// Leave star and in-file styles unchanged, and set syntax style
func (m *FileBuf) SetSyntaxStyle( l_num, c_num int, style byte ) {

//var sr BLine = m.lines.GetLP( l_num ).styles
  var lp *FLine = m.lines.GetLP( l_num )
//var S  byte = sr.GetB( c_num )
  var S  byte = lp.GetStyle( c_num )

  //< Clear everything except star and in-file
  S &= ( HI_STAR | HI_STAR_IN_F )
  S |= style;   //< Set style

//sr.SetB( c_num, S )
  lp.SetStyle( c_num, S )
}

func (m *FileBuf) HasStyle( l_num, c_num int, style byte ) bool {

  var lp *FLine = m.lines.GetLP( l_num )

  var S byte = lp.GetStyle( c_num )

  return 0 != (S & style)
}

//func (m *FileBuf) BufferEditor_SortName() {
//  sort.Sort( SLineSlice(m.lines.lines) )
//}

// Move largest file name to bottom.
// Files are grouped under directories.
// Returns true if any lines were swapped.
func (m *FileBuf) BufferEditor_SortName() bool {

  var changed bool = false
  var NUM_LINES int = m.NumLines()

  var NUM_BUILT_IN_FILES int = m_USER_FILE
//var FNAME_START_CHAR   int = 0

  // Sort lines (file names), least to greatest:
  for i:=NUM_LINES-1; NUM_BUILT_IN_FILES<i; i-- {
    for k:=NUM_BUILT_IN_FILES; k<i; k++ {

      var p_l_0 *FLine = m.lines.GetLP( k   )
      var p_l_1 *FLine = m.lines.GetLP( k+1 )

      if( m.BufferEditor_SortName_Swap( p_l_0, p_l_1 ) ) {
      //SwapLines( m, k, k+1 );
        m.lines.Swap( k, k+1 )
        changed = true;
      }
    }
  }
  return changed;
}

// Return true if p_l_0 (dname,fname)
// is greater then p_l_1 (dname,fname)
func (m *FileBuf) BufferEditor_SortName_Swap( p_l_0, p_l_1 *FLine ) bool {
  var swap bool = false

  // Tail is the directory name:
  var l_0_dn string = GetFnameTail( p_l_0.to_str() )
  var l_1_dn string = GetFnameTail( p_l_1.to_str() )

  var dn_compare int = strings.Compare( l_0_dn, l_1_dn )

  if( 0<dn_compare ) {
    // l_0 dname is greater than l_1 dname
    swap = true;

  } else if( 0==dn_compare ) {
    // l_0 dname == l_1 dname
    // Head is the file name:
    var l_0_fn string = GetFnameHead( p_l_0.to_str() )
    var l_1_fn string = GetFnameHead( p_l_1.to_str() )

    if( 0<strings.Compare( l_0_fn, l_1_fn ) ) {
      // l_0 fname is greater than l_1 fname
      swap = true;
    }
  }
  return swap;
}

func (m *FileBuf) Changed() bool {
  return m.history.Has_Changes()
}

func (m *FileBuf) ClearChanged() {
  m.history.Clear();
}

func (m *FileBuf) Find_Regexs( start_line, num_lines int ) {

  m.Check_4_New_Regex()

//if( m.p_regex_obj != nil ) {
    var up_to_line int = Min_i( start_line+num_lines, m.NumLines() );

    for k:=start_line; k<up_to_line; k++ {
      m.Find_Regexs_4_Line( k );
    }
//}
}

//func (m *FileBuf) Find_Regexs( start_line, num_lines int ) {
//
//  m.Check_4_New_Regex()
//
//  var up_to_line int = Min_i( start_line+num_lines, m.NumLines() );
//
//  for k:=start_line; k<up_to_line; k++ {
//
//    var lp *FLine = m.lines.GetLP( k )
//
//    if( !lp.star_styles_valid ) {
//      m.Find_Regexs_4_Line( lp );
//      lp.star_styles_valid = true
//    }
//  }
//}

func (m *FileBuf) Check_4_New_Regex() {

  if( m.regex_str != m_vis.regex_str ) {

    m.Invalidate_Regexs()

    m.regex_str = m_vis.regex_str

    if( 0 == len(m.regex_str) ) {
      m.p_regex_obj = nil
    } else {
      var err error
      m.p_regex_obj, err = regexp.Compile( m.regex_str )
      if( err != nil ) {
        m.p_regex_obj = nil
      }
    }
  }
}

func (m *FileBuf) Invalidate_Regexs() {

  for k:=0; k<m.lines.Len(); k++ {
    var lp *FLine = m.lines.GetLP(k)
    lp.star_styles_valid = false
  }
}

func (m *FileBuf) Find_Regexs_4_Line( line_num int ) {

  var lp *FLine = m.lines.GetLP(line_num)

  if( !lp.star_styles_valid ) {

    lp.ClearStarAndInFileStyles()

    if( nil != m.p_regex_obj ) {

      if( m.file_type == FT_BUFFER_EDITOR ||
          m.file_type == FT_DIR ) {

        if( m.Other_File_Has_My_Regex( lp.to_str() ) ) {
          LL := lp.Len()
          for k:=0; k<LL; k++ {
            lp.Set__StarInFStyle( k )
          }
        }
      }
      m.Find_Regexs_4_Line_Plain( lp )
    }
    lp.star_styles_valid = true
  }
}

//func (s string) ends_with( suffix string ) bool {
//  len_s := len(s)
//  len_suffix := len(suffix)
//  ends_w := len_suffix <= len_s && s[ len_s - len_suffix: ] == suffix
//  return ends_w
//}

func (m *FileBuf) Other_File_Has_My_Regex( file_name string ) bool {

  if( m.file_type == FT_DIR ) {
    // For FT_DIR, file_name does not contain the directory
    var path_name string = m.dir_name + file_name

    if( m.Filename_Is_Relevant( path_name ) ) {
      return m.Have_Regex_In_File( path_name )
    }
  } else if( m.file_type == FT_BUFFER_EDITOR ) {
    // For FT_BUFFER_EDITOR, file_name is the line in the BUFFER_EDITOR
    if( file_name !=  m_EDIT_BUF_NAME &&
        file_name !=  m_HELP_BUF_NAME &&
        file_name !=  m_MSG__BUF_NAME &&
        file_name != m_SHELL_BUF_NAME &&
        file_name != m_COLON_BUF_NAME &&
        file_name != m_SLASH_BUF_NAME &&
        !strings.HasSuffix( file_name, string(DIR_DELIM) ) ) {

      var pfb *FileBuf = m_vis.GetFileBuf_s( file_name );
      if( nil != pfb ) {
        return pfb.Has_Regex( m.p_regex_obj );
      }
    }
  }
  return false;
}

func (m *FileBuf) Filename_Is_Relevant( fname string ) bool {

  return strings.HasSuffix( fname, ".txt") ||
         strings.HasSuffix( fname, ".txt.new") ||
         strings.HasSuffix( fname, ".txt.old") ||
         strings.HasSuffix( fname, ".sh") ||
         strings.HasSuffix( fname, ".sh.new"  ) ||
         strings.HasSuffix( fname, ".sh.old"  ) ||
         strings.HasSuffix( fname, ".bash"    ) ||
         strings.HasSuffix( fname, ".bash.new") ||
         strings.HasSuffix( fname, ".bash.old") ||
         strings.HasSuffix( fname, ".alias"   ) ||
         strings.HasSuffix( fname, ".bash_profile") ||
         strings.HasSuffix( fname, ".bash_logout") ||
         strings.HasSuffix( fname, ".bashrc" ) ||
         strings.HasSuffix( fname, ".profile") ||
         strings.HasSuffix( fname, ".h"      ) ||
         strings.HasSuffix( fname, ".h.new"  ) ||
         strings.HasSuffix( fname, ".h.old"  ) ||
         strings.HasSuffix( fname, ".c"      ) ||
         strings.HasSuffix( fname, ".c.new"  ) ||
         strings.HasSuffix( fname, ".c.old"  ) ||
         strings.HasSuffix( fname, ".hh"     ) ||
         strings.HasSuffix( fname, ".hh.new" ) ||
         strings.HasSuffix( fname, ".hh.old" ) ||
         strings.HasSuffix( fname, ".cc"     ) ||
         strings.HasSuffix( fname, ".cc.new" ) ||
         strings.HasSuffix( fname, ".cc.old" ) ||
         strings.HasSuffix( fname, ".hpp"    ) ||
         strings.HasSuffix( fname, ".hpp.new") ||
         strings.HasSuffix( fname, ".hpp.old") ||
         strings.HasSuffix( fname, ".cpp"    ) ||
         strings.HasSuffix( fname, ".cpp.new") ||
         strings.HasSuffix( fname, ".cpp.old") ||
         strings.HasSuffix( fname, ".cxx"    ) ||
         strings.HasSuffix( fname, ".cxx.new") ||
         strings.HasSuffix( fname, ".cxx.old") ||
         strings.HasSuffix( fname, ".idl"    ) ||
         strings.HasSuffix( fname, ".idl.new") ||
         strings.HasSuffix( fname, ".idl.old") ||
         strings.HasSuffix( fname, ".idl.in"    ) ||
         strings.HasSuffix( fname, ".idl.in.new") ||
         strings.HasSuffix( fname, ".idl.in.old") ||
         strings.HasSuffix( fname, ".html"    ) ||
         strings.HasSuffix( fname, ".html.new") ||
         strings.HasSuffix( fname, ".html.old") ||
         strings.HasSuffix( fname, ".htm"     ) ||
         strings.HasSuffix( fname, ".htm.new" ) ||
         strings.HasSuffix( fname, ".htm.old" ) ||
         strings.HasSuffix( fname, ".java"    ) ||
         strings.HasSuffix( fname, ".java.new") ||
         strings.HasSuffix( fname, ".java.old") ||
         strings.HasSuffix( fname, ".js"    ) ||
         strings.HasSuffix( fname, ".js.new") ||
         strings.HasSuffix( fname, ".js.old") ||
         strings.HasSuffix( fname, ".Make"    ) ||
         strings.HasSuffix( fname, ".make"    ) ||
         strings.HasSuffix( fname, ".Make.new") ||
         strings.HasSuffix( fname, ".make.new") ||
         strings.HasSuffix( fname, ".Make.old") ||
         strings.HasSuffix( fname, ".make.old") ||
         strings.HasSuffix( fname, "Makefile" ) ||
         strings.HasSuffix( fname, "makefile" ) ||
         strings.HasSuffix( fname, "Makefile.new") ||
         strings.HasSuffix( fname, "makefile.new") ||
         strings.HasSuffix( fname, "Makefile.old") ||
         strings.HasSuffix( fname, "makefile.old") ||
         strings.HasSuffix( fname, ".stl"    ) ||
         strings.HasSuffix( fname, ".stl.new") ||
         strings.HasSuffix( fname, ".stl.old") ||
         strings.HasSuffix( fname, ".ste"    ) ||
         strings.HasSuffix( fname, ".ste.new") ||
         strings.HasSuffix( fname, ".ste.old") ||
         strings.HasSuffix( fname, ".py"    ) ||
         strings.HasSuffix( fname, ".py.new") ||
         strings.HasSuffix( fname, ".py.old") ||
         strings.HasSuffix( fname, ".sql"    ) ||
         strings.HasSuffix( fname, ".sql.new") ||
         strings.HasSuffix( fname, ".sql.old") ||
         strings.HasSuffix( fname, ".xml"     ) ||
         strings.HasSuffix( fname, ".xml.new" ) ||
         strings.HasSuffix( fname, ".xml.old" ) ||
         strings.HasSuffix( fname, ".xml.in"    ) ||
         strings.HasSuffix( fname, ".xml.in.new") ||
         strings.HasSuffix( fname, ".xml.in.old") ||
         strings.HasSuffix( fname, ".cmake"     ) ||
         strings.HasSuffix( fname, ".cmake.new" ) ||
         strings.HasSuffix( fname, ".cmake.old" ) ||
         strings.HasSuffix( fname, ".cmake"     ) ||
         strings.HasSuffix( fname, ".cmake.new" ) ||
         strings.HasSuffix( fname, ".cmake.old" ) ||
         strings.HasSuffix( fname, "CMakeLists.txt") ||
         strings.HasSuffix( fname, "CMakeLists.txt.old") ||
         strings.HasSuffix( fname, "CMakeLists.txt.new") ||
         strings.HasSuffix( fname, "go") ||
         strings.HasSuffix( fname, "go.old") ||
         strings.HasSuffix( fname, "go.new")
}

func (m *FileBuf) Have_Regex_In_File( pname string ) bool {
  found := false
  done := false
  if( IsReg( pname ) ) {
    // p_f *os.File
    p_f, err := os.Open( pname )

    if( err == nil ) {
      defer p_f.Close()

      var p_reader *bufio.Reader = bufio.NewReader( p_f )
      m_bb.Reset()
      for( !found && !done ) {
        B, err := p_reader.ReadByte()
        if( err != nil ) {
          done = true
        } else {
          if( '\n' == B ) {
            found = Bytes_Has_Regex( m_bb.Bytes(), m.p_regex_obj );
            m_bb.Reset();
          } else {
            m_bb.WriteByte( B );
          }
        }
      }
    }
  }
  return found;
}

//func (m *FileBuf) Has_Regex( p_regex_obj *regexp.Regexp ) bool {
//
//  NUM_LINES := m.lines.Len()
//  for k:=0; k<NUM_LINES; k++ {
//    var lp *FLine = m.lines.GetLP( k )
//    LL := lp.Len()
//
//    if( 0 < LL ) {
//      // Dont use FLine.to_SB because that creates a temporary string and []byte.
//      // Use m_bb because the re-uses memory already allocated.
//      m_bb.Reset()
//      for ln:=0; ln<LL; ln++ {
//        R := lp.GetR( ln )
//        m_bb.WriteRune( R )
//      }
//      if( Bytes_Has_Regex( m_bb.Bytes(), p_regex_obj ) ) {
//        return true;
//      }
//    }
//  }
//  return false;
//}

func (m *FileBuf) Has_Regex( p_regex_obj *regexp.Regexp ) bool {

  NUM_LINES := m.lines.Len()

  for k:=0; k<NUM_LINES; k++ {

    var lp *FLine = m.lines.GetLP( k )
    if( 0 < lp.Len() ) {
      if( Bytes_Has_Regex( lp.to_SB(0), p_regex_obj ) ) {
        return true;
      }
    }
  }
  return false;
}

func (m *FileBuf) Find_Regexs_4_Line_Plain( lp *FLine ) {
  var LL int = lp.Len();
  // Find the patterns for the line:
  var found bool = true;
  for p:=0; found && p<LL; {
    var match_pos, match_len int

  //found = m.Regex_Search( lp.runes.data[p:], &match_pos, &match_len )
    found = m.Regex_Search( lp.to_SB(p), &match_pos, &match_len )
    if( found ) {
      var match_st int = p + match_pos;
      var match_fn int = p + match_pos + match_len;

      for pos:=match_st; pos<LL && pos<match_fn; pos++ {
        lp.styles.data[pos] |= HI_STAR
      }
      p = match_fn;
    }
  }
}

//func (m *FileBuf) Find_Regexs_4_Line( lp *FLine ) {
//
//  lp.ClearStarAndInFileStyles()
//
//  if( nil != m.p_regex_obj ) {
//    var LL int = lp.Len();
//    // Find the patterns for the line:
//    var found bool = true;
//    for p:=0; found && p<LL; {
//      var match_pos, match_len int
//
//      found = m.Regex_Search( lp.runes.data[p:],
//                              &match_pos,
//                              &match_len ) && 0 < match_len
//      if( found ) {
//        var match_st int = p + match_pos;
//        var match_fn int = p + match_pos + match_len;
//
//        for pos:=match_st; pos<LL && pos<match_fn; pos++ {
//          lp.styles.data[pos] |= HI_STAR
//        }
//        p = match_fn;
//      }
//    }
//  }
//}

//func (m *FileBuf) Regex_Search( runes_to_search []rune,
//                                match_pos, match_len *int ) bool {
//  var found bool
//
//  str_to_search := string( runes_to_search )
////bytes_to_search := []byte( str_to_search )
//
//  var loc []int = m.p_regex_obj.FindStringIndex( str_to_search )
//
//  found = (loc != nil) && (loc[0] < loc[1])
//  if( found ) {
//    *match_pos = loc[0]
//    *match_len = loc[1] - loc[0]
//  }
//  return found
//}

//func (m *FileBuf) Regex_Search( runes_to_search []rune,
//                                match_pos, match_len *int ) bool {
//  var found bool
//
//  m_bb.Reset()
//  for _, R := range runes_to_search {
//    n_r_bytes,_ := m_bb.WriteRune( R )
//    if( 1 < n_r_bytes ) {
//      Log( fmt.Sprintf("R=(%v,%c) is %v bytes", R, R, n_r_bytes) )
//    }
//  }
//  var loc []int = m.p_regex_obj.FindIndex( m_bb.Bytes() )
//
//  found = (loc != nil) && (loc[0] < loc[1])
//  if( found ) {
//    *match_pos = loc[0]
//    *match_len = loc[1] - loc[0]
//  }
//  return found
//}

func (m *FileBuf) Regex_Search( bytes_to_search []byte,
                                match_pos, match_len *int ) bool {
  var found bool

//m_bb.Reset()
//for _, R := range runes_to_search {
//  m_bb.WriteByte( byte(R) )
//}
  var loc []int = m.p_regex_obj.FindIndex( bytes_to_search )

  found = (loc != nil) && (loc[0] < loc[1])
  if( found ) {
    *match_pos = loc[0]
    *match_len = loc[1] - loc[0]
  }
  return found
}

func (m *FileBuf) Undo( p_fv *FileView ) {

  if( m.save_history ) {
    m.save_history = false

    m.history.Undo( p_fv )

    m.save_history = true
  }
}

func (m *FileBuf) UndoAll( p_fv *FileView ) {

  if( m.save_history ) {
    m.save_history = false

    m.history.UndoAll( p_fv )

    m.save_history = true
  }
}

