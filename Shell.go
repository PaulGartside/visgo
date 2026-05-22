
package main

import (
  "fmt"
  "os/exec"
  "unicode/utf8"
  "strings"
)

type Shell struct {
  p_fv *FileView
  pfb *FileBuf

  concatinated_cmd string
  cmd_name string
  cmd_args []string

  hash_divider [40]rune
  line_divider [40]rune

  // Variables used when starting a command in the background:
  running bool

  sb_stdout strings.Builder
  sb_stderr strings.Builder

  p_cmd *exec.Cmd
  p_wait_done_chan chan error
}

func (m *Shell) Init( p_fb *FileBuf ) {

  m.pfb = p_fb

  // Fill in divideer
  for k:=0; k<40; k++ {
    m.hash_divider[k] = '#'
    m.line_divider[k] = '-'
    m.pfb.PushR( 0, '#' )
  }
  m.pfb.PushLE()

  NUM_LINES := m.pfb.NumLines()

  for w:=0; w<m.pfb.views.Len(); w++ {
    var p_fv *FileView = m.pfb.views.Get( w )

    p_fv.GoToCrsPos_NoWrite( NUM_LINES-1, 0 )
  }
}

func (m *Shell) Clean_Up() {

  m.running = false

  m.p_cmd = nil
  m.p_wait_done_chan = nil

  m.sb_stdout.Reset()
  m.sb_stderr.Reset()
}

// Run a command and block until the command completes:
//
func (m *Shell) Run() {

  m.p_fv = m_vis.CV()

  m.cmd_name = ""
  m.cmd_args = m.cmd_args[:0]

  ok := m.Get_Shell_Cmd()

  if( ok ) {
    m.Run_Cmd()
  } else {
    m.p_fv.PrintCursor()
  }
}

// Start a command and do not block while the command runs:
//
func (m *Shell) Start() {

  m.p_fv = m_vis.CV()

  m.cmd_name = ""
  m.cmd_args = m.cmd_args[:0]

  ok := m.Get_Shell_Cmd()

  if( ok ) {
    m.Start_Cmd()
  } else {
    m.p_fv.PrintCursor()
  }
}

func (m *Shell) Get_Shell_Cmd() bool {
  got_cmd := false

  if( 0 < m.pfb.NumLines() ) {
    // Find first line, which is line below last line matching '^[ ]*#':
    first_line := m.Get_first_line()

    LAST_LINE := m.pfb.NumLines()-1

    if( first_line <= LAST_LINE ) {

      m.concatinated_cmd = m.Concatenate_cmd_lines( first_line )

      if(0 < len(m.concatinated_cmd)) {
        got_cmd = m.Get_cmd_name_args( m.concatinated_cmd )
      }
    }
  }
  return got_cmd
}

// Find first line, which is line below last line matching '^[ ]*#':
//
func (m *Shell) Get_first_line() int {
  first_line := 0
  found_first_line := false
  LAST_LINE := m.pfb.NumLines()-1
  for l:=LAST_LINE; !found_first_line && 0<=l; l-- {
    LL := m.pfb.LineLen( l )
    first_non_white := 0
    for first_non_white<LL && IsSpace( m.pfb.GetR( l, first_non_white ) ) {
      first_non_white++;
    }
    if( first_non_white<LL && '#' == m.pfb.GetR( l, first_non_white ) ) {
      found_first_line = true;
      first_line = l+1;
    }
  }
  return first_line
}

func (m *Shell) Concatenate_cmd_lines( first_line int ) string {

  var sb strings.Builder
  // Concatenate all command lines into String cmd:
  LAST_LINE := m.pfb.NumLines()-1

  for k:=first_line; k<=LAST_LINE; k++ {
    LL := m.pfb.LineLen( k )
    for p:=0; p<LL; p++ {
      R := m.pfb.GetR( k, p )
      if( R == '#' ) { break } //< Ignore # to end of line
      sb.WriteRune( R )
    }
    // In the SHELL buffer, commands broken up onto multiple lines are
    // concatinated together with a space separating the lines:
    if( 0<LL && k<LAST_LINE ) { sb.WriteRune(' ') }
  }
  // Remove leading and ending white space
  concatinated_str := strings.TrimSpace( sb.String() )

  return concatinated_str
}

func (m *Shell) Get_cmd_name_args( concatinated_str string ) bool {
  got_cmd := false
  concatinated_str_len := len(concatinated_str)
  var sb strings.Builder
  var R rune
  var k, R_sz int

  for k=0; k<concatinated_str_len; k+=R_sz {
    R,R_sz = utf8.DecodeRuneInString(concatinated_str[k:])
    if( IsSpace( R ) ) {
      break
    } else {
      sb.WriteRune( R )
    }
  }
  m.cmd_name = sb.String()
  got_cmd = true

  sb.Reset()
  done := false
  for( !done ) {
    // Skip past white space:
    for ; k<concatinated_str_len; k+=R_sz {
      R,R_sz = utf8.DecodeRuneInString(concatinated_str[k:])
      if( !IsSpace( R ) ) {
        break
      }
    }
    for ; k<concatinated_str_len; k+=R_sz {
      R,R_sz = utf8.DecodeRuneInString(concatinated_str[k:])
      if( IsSpace( R ) ) {
        cmd_arg := sb.String()
        sb.Reset()
        m.cmd_args = append( m.cmd_args, cmd_arg )
        k+=R_sz
        break
      } else {
        sb.WriteRune( R )
      }
    }
    if( concatinated_str_len <= k ) {
      done = true
      cmd_arg := sb.String()
      if( 0<len(cmd_arg) ) {
        m.cmd_args = append( m.cmd_args, cmd_arg )
      }
    }
  }
  return got_cmd
}

func (m *Shell) Run_Cmd() {
  // Add ######################################
  m.pfb.PushLSR( m.hash_divider[:] )

  var p_cmd *exec.Cmd = exec.Command( m.cmd_name, m.cmd_args... )
  var sb_stdout strings.Builder
  var sb_stderr strings.Builder
  p_cmd.Stdout = &sb_stdout
  p_cmd.Stderr = &sb_stderr

  err := p_cmd.Run()

  stdout_str := sb_stdout.String()
  stderr_str := sb_stderr.String()

  m.Print_output_str( stdout_str )

  if( 0<len(stdout_str) && 0<len(stderr_str) ) {
    m.pfb.PushLSR( m.line_divider[:] )
  }
  m.Print_output_str( stderr_str )

  if( err != nil ) {
    m.pfb.PushLSR( m.line_divider[:] )
    err_msg := fmt.Sprintf("%v: %s", err, m.concatinated_cmd)
    p_fl := new( FLine )
    p_fl.from_str( err_msg )
    m.pfb.PushLP( p_fl )
  }
  m.Print_Divider_Move_Cursor_2_Bottom()
}

func (m *Shell) Handle_Err( err error ) {
  m.pfb.PushLSR( m.line_divider[:] )
  err_msg := fmt.Sprintf("%v: %s", err, m.concatinated_cmd)
  p_fl := new( FLine )
  p_fl.from_str( err_msg )
  m.pfb.PushLP( p_fl )

  m.Print_Divider_Move_Cursor_2_Bottom()

  m.Clean_Up();
}

func (m *Shell) Start_Cmd() {
  // Add ######################################
  m.pfb.PushLSR( m.hash_divider[:] )

  m.p_cmd = exec.Command( m.cmd_name, m.cmd_args... )

  m.sb_stdout.Reset()
  m.sb_stderr.Reset()
  m.p_cmd.Stdout = &m.sb_stdout
  m.p_cmd.Stderr = &m.sb_stderr

  err := m.p_cmd.Start()

  if( err != nil ) {
    m.Handle_Err( err )
  } else {
    m.p_wait_done_chan = make(chan error, 1)

    go m.Wait_Done()
    m.running = true

    // Move cursor to bottom of SHELL file view:
    NUM_LINES := m.pfb.NumLines()
    m.p_fv.GoToCrsPos_NoWrite( NUM_LINES-1, 0 )
    m.pfb.Update()
  }
}

func (m *Shell) Wait_Done() {

  m.p_wait_done_chan <- m.p_cmd.Wait()
}

func (m *Shell) Print_output_str( out_str string ) {

  p_fl := new( FLine )
  for _,R := range out_str {
    if( R == '\n' ) {
      m.pfb.PushLP( p_fl )
      p_fl = new( FLine )
    } else {
      p_fl.PushR( R )
    }
  }
  if( 0 < p_fl.Len() ) { m.pfb.PushLP( p_fl ) }
}

func (m *Shell) Add_Divider() {
  // Add ###################################### followed by empty line
  m.pfb.PushLSR( m.hash_divider[:] )
  m.pfb.PushLE()
}

func (m *Shell) Print_Divider_Move_Cursor_2_Bottom() {

  m.Add_Divider()

  // Move cursor to bottom of file
  NUM_LINES := m.pfb.NumLines()
  m.p_fv.GoToCrsPos_NoWrite( NUM_LINES-1, 0 )
  m.pfb.Update()
}

func (m *Shell) Update() {

  m.Run_Non_Blocking();
}

func (m *Shell) Run_Non_Blocking() {

  select {
    case err := <-m.p_wait_done_chan:
      m.Handle_Cmd_Done( err )

    default:
      m.Run_Non_Blocking_Read_Display_Stdout()
  }
}

func (m *Shell) Handle_Cmd_Done( err error ) {

  stdout_str := m.Run_Non_Blocking_Read_Display_Stdout()
  stderr_str := m.sb_stderr.String()

  m.sb_stderr.Reset()

  if( 0<len(stdout_str) && 0<len(stderr_str) ) {
    m.pfb.PushLSR( m.line_divider[:] )
  }
  m.Print_output_str( stderr_str )

  if( err != nil ) {
    m.pfb.PushLSR( m.line_divider[:] )
    err_msg := fmt.Sprintf("%v: %s", err, m.concatinated_cmd)
    p_fl := new( FLine )
    p_fl.from_str( err_msg )
    m.pfb.PushLP( p_fl )
  }
  m.Print_Divider_Move_Cursor_2_Bottom()
  m.Clean_Up();
}

func (m *Shell) Run_Non_Blocking_Read_Display_Stdout() string {

  // Note that the following two lines is a potential race condition.
  // The started command will be putting stdout into m.sb_stdout.
  // If after m.sb_stdout.String() is called, and before m.sb_stdout.Reset()
  // is called, the started command puts new stdout into m.sb_stdout,
  // it will be lost.
  stdout_str := m.sb_stdout.String()
  m.sb_stdout.Reset()
  m.Print_output_str( stdout_str )

  return stdout_str
}

