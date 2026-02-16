
package main

import (
  "bytes"
  "fmt"
  "io"
  "io/fs"
  "os"
  "path/filepath"
  "strings"
  "regexp"
  "time"
  "unicode"
)

// Log( fmt.Sprintf("", ) )
//
func Log( msg string ) {

  m_log = append( m_log, msg )
}

func HasPrefix( s, prefix string ) bool {
  prefix_len := len(prefix)
  return prefix_len <= len(s) && s[:prefix_len] == prefix
}

//func FindFullFileNameRel2CWD( in_fname string ) string {
//
//  var full_fname string
//
//  if( m_WIN_D_PS_len<len(in_fname) &&
//      HasPrefix(in_fname, m_WIN_D_PS) ) {
//    // in_fname is already a full path, so just return it:
//    full_fname = in_fname
//
//  } else if( 0<len(in_fname) &&
//             os.PathSeparator == in_fname[0] ) {
//    // in_fname is already a full path, but just needs windows drive:
//    full_fname = m_WIN_D + in_fname
//
//  } else {
//    full_fname = fmt.Sprintf("%s%c%s", m_wd, os.PathSeparator, in_fname )
//  }
//  full_fname = filepath.Clean( full_fname )
//
//  if( IsDir( full_fname ) ) {
//    full_fname = AppendDirDelim( full_fname )
//  }
//  return full_fname
//}

func FindFullFileNameRel2CWD( in_fname string ) string {

  var full_fname string

  if( 0 < len(m_WIN_D) ) {
    full_fname = FindFullFileNameRel2CWD_win( in_fname )
  } else {
    full_fname = FindFullFileNameRel2CWD_unix( in_fname )
  }
  full_fname = filepath.Clean( full_fname )

  if( IsDir( full_fname ) ) {
    full_fname = AppendDirDelim( full_fname )
  }
  return full_fname
}

func FindFullFileNameRel2CWD_win( in_fname string ) string {

  var full_fname string

  if( m_WIN_D_PS_len<len(in_fname) &&
      HasPrefix(in_fname, m_WIN_D_PS) ) {
    // in_fname is already a full path, so just return it:
    full_fname = in_fname

  } else if( 0<len(in_fname) &&
             os.PathSeparator == in_fname[0] ) {
    // in_fname is already a full path, but just needs windows drive:
    full_fname = m_WIN_D + in_fname

  } else {
    full_fname = fmt.Sprintf("%s%c%s", m_wd, os.PathSeparator, in_fname )
  }
  return full_fname
}

func FindFullFileNameRel2CWD_unix( in_fname string ) string {

  var full_fname string

  if( 0<len(in_fname) && os.PathSeparator == in_fname[0] ) {
    // in_fname is already a full path, so just return it:
    full_fname = in_fname

  } else {
    full_fname = fmt.Sprintf("%s%c%s", m_wd, os.PathSeparator, in_fname )
  }
  return full_fname
}

// Finds full file name of file or directory of in_fname relative to rel_2_dir.
// The full file name found does not need to exist to return success.
// Returns the full file name.
//
func FindFullFileNameRel2( rel_2_dir, in_fname string ) string {

  var full_fname string

  if( 0==len( rel_2_dir ) ||
      "." == rel_2_dir ) {
    full_fname = FindFullFileNameRel2CWD( in_fname )

  } else if( 0<len( in_fname ) && DIR_DELIM == in_fname[0] ) {
    // in_fname is already a full path, but just needs windows drive:
    full_fname = m_WIN_D + in_fname

  } else {
    full_fname = fmt.Sprintf("%s%c%s", rel_2_dir, os.PathSeparator, in_fname )
  }
  full_fname = filepath.Clean( full_fname )

  if( IsDir( full_fname ) ) {
    full_fname = AppendDirDelim( full_fname )
  }
  return full_fname
}

func FileExists( file_name string ) bool {
  // var info FileInfo
  _, err := os.Stat( file_name )

  if os.IsNotExist( err ) {
    return false
  }
  return true
}

//func IsDir( file_name string ) bool {
//  // var info FileInfo
//  info, err := os.Stat( file_name )
//  if os.IsNotExist( err ) {
//    return false
//  }
//  var file_mode fs.FileMode = info.Mode()
//
//  return file_mode.IsDir()
//}

//func IsDir( file_name string ) bool {
//  // var info FileInfo
//  info, err := os.Stat( file_name )
//  if( err != nil ) {
//    return false
//  }
//  var file_mode fs.FileMode = info.Mode()
//
//  return file_mode.IsDir()
//}

func IsDir( file_name string ) bool {
  is_dir := false
  // var info FileInfo
  info, err := os.Stat( file_name )
  if( err == nil ) {
    var file_mode fs.FileMode = info.Mode()

    is_dir = file_mode.IsDir()
  }
  return is_dir
}

//func IsReg( file_name string ) bool {
//  // var info FileInfo
//  info, err := os.Stat( file_name )
//  if os.IsNotExist( err ) {
//    return false
//  }
//  var file_mode fs.FileMode = info.Mode()
//
//  return file_mode.IsRegular()
//}

//func IsReg( file_name string ) bool {
//  // var info FileInfo
//  info, err := os.Stat( file_name )
//  if( err != nil ) {
//    return false
//  }
//  var file_mode fs.FileMode = info.Mode()
//
//  return file_mode.IsRegular()
//}

func IsReg( file_name string ) bool {
  is_reg := false
  // var info FileInfo
  info, err := os.Stat( file_name )
  if( err == nil ) {
    var file_mode fs.FileMode = info.Mode()

    is_reg = file_mode.IsRegular()
  }
  return is_reg
}

func AppendDirDelim( s string ) string {
  var s_len int = len( s )

  if 0<s_len && DIR_DELIM != s[s_len-1] {
    s = s + DIR_DELIM_S
  }
  return s
}

func IsSpace( R rune ) bool {

  return R == ' ' || R == '\t' || R == '\n' || R == '\r'
}

func NotSpace( R rune ) bool {

  return !IsSpace( R )
}

func IsDigit( R rune ) bool {

  return '0' <= R && R <= '9'
}

func IsXDigit( R rune ) bool {

  return ('0' <= R && R <= '9') ||
         ('a' <= R && R <= 'f') ||
         ('A' <= R && R <= 'F')
}

//func RemoveSpaces( s_b []byte ) []byte {
//
//  for k:=0; k<len(s_b); k++ {
//
//    if( IsSpace( s_b[k] ) ) {
//      copy( s_b[k:], s_b[k+1:] )
//      s_b = s_b[:len(s_b)-1]
//      k--
//    }
//  }
//  return s_b
//}

func Min_i( a,b int ) int {

  if( a < b ) { return a }
  return b
}

func Max_i( a,b int ) int {

  if( a < b ) { return b }
  return a
}

// Line length minus 1
//
func LLM1( k int ) int {

  if( 0 < k ) { return k-1 }
  return 0
}

func True_1_else_2( condition bool, v1, v2 int ) int {

  if( condition ) { return v1; }
  return v2
}

func True_1_else_2_b( condition bool, v1, v2 byte ) byte {

  if( condition ) { return v1; }
  return v2
}

func True_1_else_2_r( condition bool, v1, v2 rune ) rune {

  if( condition ) { return v1; }
  return v2
}

func IsAlnum( R rune ) bool {

  return unicode.IsLetter( R ) || unicode.IsNumber( R )
}

func IsWord_Ident( R rune ) bool {

  return IsAlnum( R ) || R == '_'
}

func IsWord_NonIdent( R rune ) bool {

  return !IsSpace( R ) && !IsWord_Ident( R )
}

func IsXML_Ident( R rune ) bool {

  return IsAlnum( R ) ||
         R == '_' ||
         R == '-' ||
         R == '.' ||
         R == ':'
}

func Swap( A, B *int ) {

  var T int = *B
  *B = *A
  *A = T
}

func IsFileNameChar( R rune ) bool {
  return '$' == R               ||
         '+' == R               ||
         '-' == R               ||
         '.' == R               ||
    DIR_DELIM== R               ||
       ( '0' <= R && R <= '9' ) ||
       ( 'A' <= R && R <= 'Z' ) ||
         '_' == R               ||
       ( 'a' <= R && R <= 'z' ) ||
         '~' == R               ||
         ' ' == R
}

// Remove leading and trailing white space
func Trim( ln RLine ) {
  Trim_Beg( ln )
  Trim_End( ln )
}

// Remove leading white space
func Trim_Beg( ln RLine ) {

  var done bool = false
  for k:=0; !done && k<ln.Len(); k++ {

    if( IsSpace( ln.GetR( k ) ) ) {
      ln.RemoveR( k )
      // Since we just shifted down over current char, re-check current char
      k--
    } else {
      done = true
    }
  }
}

// Remove trailing white space
func Trim_End( ln RLine ) {

  var LEN int = ln.Len()
  if( 0 < LEN ) {
    var done bool = false
    for k:=LEN-1; !done && -1<k; k-- {

      if( IsSpace( ln.GetR( k ) ) ) {
        ln.RemoveR( k )
      } else {
        done = true
      }
    }
  }
}

//func IsAlnum( R rune ) bool {
//
//  return ('0' <= R && R <= '9') ||
//         ('a' <= R && R <= 'z') ||
//         ('A' <= R && R <= 'Z')
//}

func IsIdent( R rune ) bool {

  return IsAlnum( R ) || R == '_'
}

func line_start_or_prev_C_non_ident( line RLine, p int ) bool {

  if( 0==p ) {
    return true // p is on line start
  }
  // At this point 0 < p
  var C rune = line.GetR( p-1 )
  if( !IsAlnum( C ) && C!='_' ) {
    // C is not an identifier
    return true
  }
  // On identifier
  return false
}

func line_end_or_non_ident( line RLine, LL, p int ) bool {

  if( p == LL-1 ) {
    return true // p is on line end
  }
  if( p < LL-1 ) {
    // At this point p should always be less than LL-1,
    // but put the check in above just to be safe.
    // The check above could also be implemented as an ASSERT.
    var C rune = line.GetR(p+1)
    if( !IsAlnum( C ) && C!='_' ) {
      // C is not an identifier
      return true
    }
  }
  // C is an identifier
  return false
}

func Quote_Start( qt, c2, c1, c0 rune ) bool {

  return (c1==0    && c0==qt ) || //< Quote at beginning of line
         (c1!='\\' && c0==qt ) || //< Non-escaped quote
         (c2=='\\' && c1=='\\' && c0==qt ) //< Escaped escape before quote
}

func OneVarType( c0 rune ) bool {

  return (c0=='&') ||
         (c0=='.'  || c0=='*') ||
         (c0=='['  || c0==']')
}

func OneControl( c0 rune ) bool {

  return c0=='=' || c0=='^' || c0=='~' ||
         c0==':' || c0=='%' ||
         c0=='+' || c0=='-' ||
         c0=='<' || c0=='>' ||
         c0=='!' || c0=='?' ||
         c0=='(' || c0==')' ||
         c0=='{' || c0=='}' ||
         c0==',' || c0==';' ||
         c0=='/' || c0=='|'
}
func TwoControl( c1, c0 rune ) bool {

  return (c1=='=' && c0=='=') ||
         (c1=='&' && c0=='&') ||
         (c1=='|' && c0=='|') ||
         (c1=='|' && c0=='=') ||
         (c1=='&' && c0=='=') ||
         (c1=='!' && c0=='=') ||
         (c1=='+' && c0=='=') ||
         (c1=='-' && c0=='=')
}

//func GetFnameHead( string in_full_fname ) string {
//
//  string head
//
//  // This const_cast is okay because we are not changing in_fname_cp:
//  char* in_fname_cp = CCast<char*>(in_full_fname)
//  char* const last_slash = strrchr( in_fname_cp, DIR_DELIM )
//  if( last_slash )
//  {
//    for( const char* cp = last_slash + 1; *cp; cp++ ) head.push( *cp )
//  }
//  else {
//    // No tail, all head:
//    for( const char* cp = in_full_fname; *cp; cp++ ) head.push( *cp )
//  }
//  return head
//}

// Here Head is file name.
// Return portion of in_full_fname after DIR_DELIM.
// If DIR_DELIM is not in in_full_fname, return in_full_fname.
//
func GetFnameHead( in_full_fname string ) string {

  var head string

  var index_of_last_slash int = strings.LastIndexByte( in_full_fname, DIR_DELIM )

  if( -1 < index_of_last_slash ) {
    head = in_full_fname[index_of_last_slash+1:]
  } else {
    // No tail, all head:
    head = in_full_fname
  }
  return head
}

// Here Tail is directory.
// Return portion of in_full_fname before last DIR_DELIM.
// If DIR_DELIM is not in in_full_fname, return empty string.
//
func GetFnameTail( in_full_fname string ) string {

  var tail string

  var index_of_last_slash int = strings.LastIndexByte( in_full_fname, DIR_DELIM )

  // This const_cast is okay because we are not changing in_fname_cp:
  if( -1 < index_of_last_slash ) {
    tail = in_full_fname[0:index_of_last_slash]
  }
  return tail
}

func Bytes_Has_Regex( s_b []byte, p_regex_obj *regexp.Regexp ) bool {

  var loc []int = p_regex_obj.FindIndex( s_b )

  has := (loc != nil) && (loc[0] < loc[1])
  return has
}

// type FileInfo interface {
//   Name() string
//   Size() int64
//   Mode() FileMode
//   ModTime() time.Time
//   IsDir() bool
//   Sys() any
// }
func ModificationTime( fname string ) time.Time {

  var mod_time time.Time

  file_info, err := os.Stat( fname )
  if( err == nil ) {
    mod_time = file_info.ModTime()
  }
  return mod_time
}

// Returns true,nil if path_1 and path_2 are the same size.
// Returns false,nil if path_1 and path_2 are NOT the same size,
//   and no error occured.
// Returns false,err if an error occured.
//
func Files_Are_Same_Size( path_1, path_2 string ) (bool, error) {

  fi_1, err := os.Stat(path_1)
  if err != nil {
    return false, err
  }
  fi_2, err := os.Stat(path_2)
  if err != nil {
    return false, err
  }
  return (fi_1.Size() == fi_2.Size()), nil
}

// See of Files are Same using File Paths.
// Returns true,nil if contents of path_1 and path_2 are the same.
// Returns false,nil if contents of path_1 and path_2 are NOT the same,
//   and no error occured.
// Returns false,err if an error occured.
//
func Files_Are_Same_p( path_1, path_2 string ) (bool, error) {
  // Assume files are the same until proven otherwise:
  same,err := Files_Are_Same_Size( path_1, path_2 )
  if( same ) {
    f1, err := os.Open(path_1)
    same = (err == nil)
    if( same ) {
      defer f1.Close()

      f2, err := os.Open(path_2)
      same = (err == nil)
      if( same ) {
        defer f2.Close()

        chunkSize := 1024
        buf_1 := make([]byte, chunkSize)
        buf_2 := make([]byte, chunkSize)

        for same {
          n1, err1 := f1.Read(buf_1)
          n2, err2 := f2.Read(buf_2)

          if( err1 == io.EOF && err2 == io.EOF ) {
            // Reached end of both files, and all chunks matched
            break
          } else if( (err1 == io.EOF && err2 != io.EOF) ||
                     (err1 != io.EOF && err2 == io.EOF) ) {
            // Content differs. Should never get here file sizes were verified to be the same above.
            same = false
          } else { // err1 != io.EOF && err2 != io.EOF
            if( err1 != nil ) {
              same = false
              err = err1
            } else if( err2 != nil ) {
              same = false
              err = err2
            } else { // No errors, compare file contents:
              if( n1 != n2 || !bytes.Equal(buf_1[:n1], buf_2[:n2]) ) {
                same = false
              }
            }
          }
        }
      }
    }
  }
  return same, err
}

// See of Files are Same using FileBuf objects.
// Returns true,nil if contents of path_1 and path_2 are the same.
// Returns false,nil if contents of path_1 and path_2 are NOT the same,
//   and no error occured.
// Returns false,err if an error occured.
//
func Files_Are_Same_o( pfb_1, pfb_2 *FileBuf ) bool {
  files_are_same := false

  if( pfb_1.is_regular && pfb_2.is_regular ) {
    num_lines_s := pfb_1.NumLines()
    num_lines_l := pfb_2.NumLines()

    if( (num_lines_s == num_lines_l) &&
        (pfb_1.lines.LF_at_EOF == pfb_2.lines.LF_at_EOF) ) {
      files_are_same = true

      for k:=0; files_are_same && k<num_lines_s; k++ {
        var p_fln_1 *FLine = pfb_1.GetLP( k )
        var p_fln_2 *FLine = pfb_2.GetLP( k )

        if( !p_fln_1.EqualLP( p_fln_2 ) ) {
          files_are_same = false
        }
      }
    }
  }
  return files_are_same
}

// dir1 is direct parent of dir2
// Example:
// dir1 = /a/b/c/
// dir2 = /a/b/c/d/
//
func dir1_is_parent_dir_of_dir2( dir1, dir2 string ) bool {
  if( len(dir1)+1 < len(dir2) ) {
    if( strings.HasPrefix( dir2, dir1) ) {
      start_idx  := len(dir1)
      finish_idx := len(dir2)

      for dir2[start_idx]   ==DIR_DELIM { start_idx++ }
      for dir2[finish_idx-1]==DIR_DELIM { finish_idx-- }

      for k:=start_idx; k<finish_idx; k++ {
        if( dir2[k] == DIR_DELIM ) { return false }
      }
      return true
    }
  }
  return false
}

// Example: If dir_name = /a/b/c/d/
// returns: d
//
func get_last_dir_of( dir_name string ) string {
  var last_dir strings.Builder

  if( 1 < len(dir_name) ) {
    finish_idx := len(dir_name) - 1

    // Backup finish_idx to before one or more DIR_DELIM's at end of dir_name:
    for 0<finish_idx && dir_name[finish_idx]==DIR_DELIM { finish_idx-- }

    start_idx := finish_idx;
    for( 0<=start_idx && dir_name[start_idx]!=DIR_DELIM ) { start_idx-- }
    start_idx++

    for k:=start_idx; k<=finish_idx; k++ { last_dir.WriteByte( dir_name[k] ) }
  }
  return last_dir.String()
}

