
package main

import (
//"fmt"
//"strings"
)

type Hi_State int

const (
  St_In_None Hi_State = iota
  St_OpenTag_ElemName
  St_OpenTag_AttrName
  St_OpenTag_AttrVal
  St_CloseTag
  St_XML_Comment
  St_SingleQuote
  St_DoubleQuote
  St_NumberBeg
  St_NumberDec
  St_NumberHex
  St_NumberExponent
  St_NumberFraction
  St_NumberTypeSpec

  St_JS_None
  St_JS_Define
  St_JS_C_Comment
  St_JS_CPP_Comment
  St_JS_SingleQuote
  St_JS_DoubleQuote
  St_JS_NumberBeg
  St_JS_NumberDec
  St_JS_NumberHex
  St_JS_NumberExponent
  St_JS_NumberFraction
  St_JS_NumberTypeSpec

  St_CS_None
  St_CS_C_Comment
  St_CS_SingleQuote
  St_CS_DoubleQuote

  St_Done
)

// JavaScript or CSS edges
type Edges struct {
  beg CrsPos
  end CrsPos
}

// cp is between beg and end
func (m *Edges) contains( cp CrsPos ) bool {

  cp_past_beg_edge := m.beg.crsLine < cp.crsLine ||
                      ( m.beg.crsLine == cp.crsLine && m.beg.crsChar < cp.crsChar )

  // end == 0 means end edge is somewhere ahead of cp
  cp_before_end_edge := ( 0 == m.end.crsLine && 0 == m.end.crsChar ) ||
                        ( cp.crsLine < m.end.crsLine ||
                          ( cp.crsLine == m.end.crsLine && cp.crsChar < m.end.crsChar ) )

  return cp_past_beg_edge && cp_before_end_edge
}

// cp is less than beg, or
// beg is greater then or equal to cp
func (m *Edges) ge( cp CrsPos ) bool {

  // beg == null means beg edge is somewhere ahead of cp
  cp_before_beg_edge := cp.crsLine < m.beg.crsLine ||
                        ( cp.crsLine == m.beg.crsLine && cp.crsChar < m.beg.crsChar )

  return cp_before_beg_edge
}

type Highlight_HTML struct {
  p_fb   *FileBuf
//l      int // Line
//p      int // Position on line
  state  Hi_State // Current state
  qtXSt  Hi_State // Quote exit state
  ccXSt  Hi_State // C comment exit state
  numXSt Hi_State // Number exit state

  // Variables to go in and out of JS:
  OpenTag_was_script bool
  OpenTag_was_style  bool
  JS_edges           Vector[Edges]
  CS_edges           Vector[Edges]
}

func (m *Highlight_HTML) Init( p_fb *FileBuf ) {
  m.p_fb = p_fb
}

func (m *Highlight_HTML) Run_Range( st CrsPos, fn int ) {

  m.state = m.Run_Range_Get_Initial_State( st )

  l := st.crsLine
  p := st.crsChar

  for St_Done != m.state && l<fn {
    state_was_JS := m.JS_State( m.state )
    st_l := l
    st_p := p

    l,p = m.Run_State( l, p )

    if( state_was_JS ) {
      st := CrsPos{ st_l, st_p }

      m.Find_Styles_Keys_In_Range( st, l+1 )
    }
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_HTML) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.p_fb, HiPairs_JS[:], st, fn )
}

var HTML_Tags = [...]string {
  "DOCTYPE",
  "abbr"    , "address"   , "area"      , "article" ,
  "aside"   , "audio"     , "a"         , "base"    ,
  "bdi"     , "bdo"       , "blockquote", "body"    ,
  "br"      , "button"    , "b"         , "canvas"  ,
  "caption" , "cite"      , "code"      , "col"     ,
  "colgroup", "datalist"  , "dd"        , "del"     ,
  "details" , "dfn"       , "dialog"    , "div"     ,
  "dl"      , "dt"        , "em"        , "embed"   ,
  "fieldset", "figcaption", "figure"    , "footer"  ,
  "form"    , "h1"        , "h2"        , "h3"      ,
  "h4"      , "h5"        , "h6"        , "head"    ,
  "header"  , "hr"        , "html"      , "ifname"  ,
  "img"     , "input"     , "ins"       , "i"       ,
  "kbd"     , "keygen"    , "label"     , "legend"  ,
  "link"    , "li"        , "main"      , "map"     ,
  "mark"    , "menu"      , "menuitem"  , "meta"    ,
  "meter"   , "nav"       , "noscript"  , "object"  ,
  "ol"      , "optgroup"  , "option"    , "p"       ,
  "param"   , "picture"   , "pre"       , "progress",
  "q"       , "rp"        , "rt"        , "ruby"    ,
  "samp"    , "script"    , "section"   , "select"  ,
  "small"   , "source"    , "span"      , "strong"  ,
  "style"   , "sub"       , "summary"   , "sup"     ,
  "s"       , "table"     , "tbody"     , "td"      ,
  "textarea", "tfoot"     , "thread"    , "th"      ,
  "time"    , "title"     , "tr"        , "track"   ,
  "ul"      , "u"         , "var"       , "video"   ,
  "wbr"     ,
}

var HiPairs_JS = [...]HiKeyVal {
  // Keywords:
  { "break"     , HI_CONTROL },
  { "break"     , HI_CONTROL },
  { "catch"     , HI_CONTROL },
  { "case"      , HI_CONTROL },
  { "continue"  , HI_CONTROL },
  { "debugger"  , HI_CONTROL },
  { "default"   , HI_CONTROL },
  { "delete"    , HI_CONTROL },
  { "do"        , HI_CONTROL },
  { "else"      , HI_CONTROL },
  { "finally"   , HI_CONTROL },
  { "for"       , HI_CONTROL },
  { "function"  , HI_CONTROL },
  { "if"        , HI_CONTROL },
  { "in"        , HI_CONTROL },
  { "instanceof", HI_CONTROL },
  { "new"       , HI_VARTYPE },
  { "return"    , HI_CONTROL },
  { "switch"    , HI_CONTROL },
  { "throw"     , HI_CONTROL },
  { "try"       , HI_CONTROL },
  { "typeof"    , HI_VARTYPE },
  { "var"       , HI_VARTYPE },
  { "void"      , HI_VARTYPE },
  { "while"     , HI_CONTROL },
  { "with"      , HI_CONTROL },

  // Keywords in strict mode:
  { "implements", HI_CONTROL },
  { "interface" , HI_CONTROL },
  { "let"       , HI_VARTYPE },
  { "package"   , HI_DEFINE  },
  { "private"   , HI_CONTROL },
  { "protected" , HI_CONTROL },
  { "public"    , HI_CONTROL },
  { "static"    , HI_VARTYPE },
  { "yield"     , HI_CONTROL },

  // Constants:
  { "false", HI_CONST },
  { "null" , HI_CONST },
  { "true" , HI_CONST },

  // Global variables and functions:
  { "arguments"         , HI_VARTYPE },
  { "Array"             , HI_VARTYPE },
  { "Boolean"           , HI_VARTYPE },
  { "Date"              , HI_CONTROL },
  { "decodeURI"         , HI_CONTROL },
  { "decodeURIComponent", HI_CONTROL },
  { "encodeURI"         , HI_CONTROL },
  { "encodeURIComponent", HI_CONTROL },
  { "Error"             , HI_VARTYPE },
  { "eval"              , HI_CONTROL },
  { "EvalError"         , HI_CONTROL },
  { "Function"          , HI_CONTROL },
  { "Infinity"          , HI_CONST   },
  { "isFinite"          , HI_CONTROL },
  { "isNaN"             , HI_CONTROL },
  { "JSON"              , HI_CONTROL },
  { "Math"              , HI_CONTROL },
  { "NaN"               , HI_CONST   },
  { "Number"            , HI_VARTYPE },
  { "Object"            , HI_VARTYPE },
  { "parseFloat"        , HI_CONTROL },
  { "parseInt"          , HI_CONTROL },
  { "RangeError"        , HI_VARTYPE },
  { "ReferenceError"    , HI_VARTYPE },
  { "RegExp"            , HI_CONTROL },
  { "String"            , HI_VARTYPE },
  { "SyntaxError"       , HI_VARTYPE },
  { "TypeError"         , HI_VARTYPE },
  { "undefined"         , HI_CONST   },
  { "URIError"          , HI_VARTYPE },

  { ""        , 0 },
}

func (m *Highlight_HTML) Run_Range_Get_Initial_State( st CrsPos ) Hi_State {

  var initial Hi_State = St_In_None

  if( m.Get_Initial_State( st, m.JS_edges, m.CS_edges ) ) {
    initial = St_JS_None
  } else if( m.Get_Initial_State( st, m.CS_edges, m.JS_edges ) ) {
    initial = St_CS_None
  }
  return initial
}

func (m *Highlight_HTML) Get_Initial_State( st CrsPos,
                                            edges_1, edges_2 Vector[Edges] ) bool {
  var found_containing_edges bool = false

  for k:=0; k<edges_1.Len(); k++ {
    if( !found_containing_edges ) {
      var edges Edges = edges_1.Get(k)
      if( edges.contains( st ) ) {
        found_containing_edges = true
      }
    } else {
      // Since a change was made at st, all the following edges
      // have been invalidated, so remove all following elements
      edges_1.Remove( k )
      k-- //< Since the current element was just removed, stay on k
    }
  }
  if( found_containing_edges ) {
    // Remove all CS_edges past st:
    for k:=0; k<edges_2.Len(); k++ {
      var edges Edges = edges_2.Get(k)
      if( edges.ge( st ) ) {
        edges_2.Remove( k )
        k-- //< Since the current element was just removed, stay on k
      }
    }
  }
  return found_containing_edges
}

func (m *Highlight_HTML) JS_State( state Hi_State ) bool {

  return state == St_JS_None ||
         state == St_JS_Define ||
         state == St_JS_SingleQuote ||
         state == St_JS_DoubleQuote ||
         state == St_JS_C_Comment ||
         state == St_JS_CPP_Comment ||
         state == St_JS_NumberBeg ||
         state == St_JS_NumberDec ||
         state == St_JS_NumberHex ||
         state == St_JS_NumberFraction ||
         state == St_JS_NumberExponent ||
         state == St_JS_NumberTypeSpec
}

func (m *Highlight_HTML) Run_State( l, p int ) (int,int) {

  switch( m.state ) {
  case St_In_None         : l,p = m.Hi_In_None         ( l,p )
  case St_XML_Comment     : l,p = m.Hi_XML_Comment     ( l,p )
  case St_CloseTag        : l,p = m.Hi_CloseTag        ( l,p )
  case St_NumberBeg       : l,p = m.Hi_NumberBeg       ( l,p )
  case St_NumberHex       : l,p = m.Hi_NumberHex       ( l,p )
  case St_NumberDec       : l,p = m.Hi_NumberDec       ( l,p )
  case St_NumberExponent  : l,p = m.Hi_NumberExponent  ( l,p )
  case St_NumberFraction  : l,p = m.Hi_NumberFraction  ( l,p )
  case St_NumberTypeSpec  : l,p = m.Hi_NumberTypeSpec  ( l,p )
  case St_OpenTag_ElemName: l,p = m.Hi_OpenTag_ElemName( l,p )
  case St_OpenTag_AttrName: l,p = m.Hi_OpenTag_AttrName( l,p )
  case St_OpenTag_AttrVal : l,p = m.Hi_OpenTag_AttrVal ( l,p )
  case St_SingleQuote     : l,p = m.Hi_SingleQuote     ( l,p )
  case St_DoubleQuote     : l,p = m.Hi_DoubleQuote     ( l,p )

  case St_JS_None          : l,p = m.Hi_JS_None       ( l,p )
  case St_JS_Define        : l,p = m.Hi_JS_Define     ( l,p )
  case St_JS_SingleQuote   : l,p = m.Hi_SingleQuote   ( l,p )
  case St_JS_DoubleQuote   : l,p = m.Hi_DoubleQuote   ( l,p )
  case St_JS_C_Comment     : l,p = m.Hi_C_Comment     ( l,p )
  case St_JS_CPP_Comment   : l,p = m.Hi_JS_CPP_Comment( l,p )
  case St_JS_NumberBeg     : l,p = m.Hi_NumberBeg     ( l,p )
  case St_JS_NumberDec     : l,p = m.Hi_NumberDec     ( l,p )
  case St_JS_NumberHex     : l,p = m.Hi_NumberHex     ( l,p )
  case St_JS_NumberFraction: l,p = m.Hi_NumberFraction( l,p )
  case St_JS_NumberExponent: l,p = m.Hi_NumberExponent( l,p )
  case St_JS_NumberTypeSpec: l,p = m.Hi_NumberTypeSpec( l,p )

  case St_CS_None          : l,p = m.Hi_CS_None    ( l,p )
  case St_CS_C_Comment     : l,p = m.Hi_C_Comment  ( l,p )
  case St_CS_SingleQuote   : l,p = m.Hi_SingleQuote( l,p )
  case St_CS_DoubleQuote   : l,p = m.Hi_DoubleQuote( l,p )

  default:
    m.state = St_In_None
  }
  return l,p
}

func (m *Highlight_HTML) Hi_In_None( l, p int ) (int,int) {

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1 is ahead of c2: (c3,c2,c1,c0)
      var c3 rune = 0; if( 2<p ) { c3 = m.p_fb.GetR( l, p-3 ) }
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c1=='<' && c0!='!' && c0!='/') {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE )
        m.state = St_OpenTag_ElemName

      } else if( c1=='<' && c0=='/') {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE ) //< '<'
        m.p_fb.SetSyntaxStyle( l, p  , HI_DEFINE ) //< '/'
        p++ // Move past '/'
        m.state = St_CloseTag

      } else if( c3=='<' && c2=='!' && c1=='-' && c0=='-') {
        m.p_fb.SetSyntaxStyle( l, p-3, HI_COMMENT ) //< '<'
        m.p_fb.SetSyntaxStyle( l, p-2, HI_COMMENT ) //< '!'
        m.p_fb.SetSyntaxStyle( l, p-1, HI_COMMENT ) //< '-'
        m.p_fb.SetSyntaxStyle( l, p  , HI_COMMENT ) //< '-'
        p++ // Move past '-'
        m.state = St_XML_Comment

      } else if( c3=='<' && c2=='!' && c1=='D' && c0=='O') {
        // <!DOCTYPE html>
        m.p_fb.SetSyntaxStyle( l, p-3, HI_DEFINE ) //< '<'
        m.p_fb.SetSyntaxStyle( l, p-2, HI_DEFINE ) //< '!'
        p-- // Move back to 'D'
        m.state = St_OpenTag_ElemName

      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
        m.state  = St_NumberBeg
        m.numXSt = St_In_None

      } else {
        ; //< No syntax highlighting on content outside of <>tags
      }
      if( St_In_None != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_XML_Comment( l, p int ) (int,int) {

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )

      // c0 is ahead of c1 is ahead of c2: (c2,c1,c0)
      var c2 rune = 0; if( 1<p ) { c2 = m.p_fb.GetR( l, p-2 ) }
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c2=='-' && c1=='-' && c0=='>' ) {
        p++ // Move past '>'
        m.state = St_In_None
      }
      if( St_XML_Comment != m.state ) { return l,p }
    }
    p = 0
  }
  return l,p
}

func (m *Highlight_HTML) Hi_CloseTag( l, p int ) (int,int) {

  found_elem_name := false

  for ; l<m.p_fb.NumLines(); l++ {
    var p_fl *FLine = m.p_fb.GetLP( l )
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++ // Move past '>'
        m.state = St_In_None

      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( !found_elem_name ) {
        // Returns non-zero if p_fl has HTTP tag at p:
        tag_len := m.Has_HTTP_Tag_At( p_fl, p )
        if( 0<tag_len ) {
          found_elem_name = true
          for k:=0; k<tag_len; k++ {
            m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
            p++
          }
          p--
        } else if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else { //( found_elem_name )
        if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      }
      if( St_CloseTag != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_NumberBeg( l, p int ) (int,int) {

  m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

  var c1 rune = m.p_fb.GetR( l, p )
  p++ //< Move past first digit

  old_state := m.state
  if( St_JS_NumberBeg == old_state ) { m.state = St_JS_NumberDec
  } else                             { m.state = St_NumberDec
  }
  LL := m.p_fb.LineLen( l )
  if( '0' == c1 && (p+1)<LL ) {
    var c0 rune = m.p_fb.GetR( l, p )
    if( 'x' == c0 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

      if( St_JS_NumberBeg == old_state ) { m.state = St_JS_NumberHex
      } else                             { m.state = St_NumberHex
      }
      p++ //< Move past 'x'
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_NumberHex( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.numXSt
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsXDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.numXSt
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_NumberDec( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.numXSt
  } else {
    var c1 rune = m.p_fb.GetR( l, p )

    if( '.'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      if( St_JS_NumberDec == m.state ) { m.state = St_JS_NumberFraction
      } else                           { m.state = St_NumberFraction
      }
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )

      if( St_JS_NumberDec == m.state ) { m.state = St_JS_NumberExponent
      } else                           { m.state = St_NumberExponent
      }
      p++
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p )
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
          p++
        }
      }
    } else if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else if( c1=='L' || c1=='F' || c1=='U' ) {
      if( St_JS_NumberDec == m.state ) { m.state = St_JS_NumberTypeSpec
      } else                           { m.state = St_NumberTypeSpec
      }
    } else {
      m.state = m.numXSt
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_NumberExponent( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.numXSt
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else {
      m.state = m.numXSt
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_NumberFraction( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )
  if( LL <= p ) { m.state = m.numXSt
  } else {
    var c1 rune = m.p_fb.GetR( l, p )
    if( IsDigit(c1) ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      p++
    } else if( 'e'==c1 || 'E'==c1 ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      if( St_JS_NumberFraction == m.state ) { m.state = St_JS_NumberExponent
      } else                                { m.state = St_NumberExponent
      }
      p++
      if( p<LL ) {
        var c0 rune = m.p_fb.GetR( l, p )
        if( '+' == c0 || '-' == c0 ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
          p++
        }
      }
    } else {
      m.state = m.numXSt
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_NumberTypeSpec( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )
  if( p < LL ) {
    var c0 rune = m.p_fb.GetR( l, p )

    if( c0=='L' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      p++
      m.state = m.numXSt

    } else if( c0=='F' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      p++
      m.state = m.numXSt

    } else if( c0=='U' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
      if( p<LL ) {
        var c1 rune = m.p_fb.GetR( l, p )
        if( c1=='L' ) { // UL
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
          if( p<LL ) {
            var c2 rune = m.p_fb.GetR( l, p )
            if( c2=='L' ) { // ULL
              m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE ); p++
            }
          }
        }
      }
      m.state = m.numXSt
    }
  }
  return l,p
}

func (m *Highlight_HTML) Hi_OpenTag_ElemName( l, p int ) (int,int) {

  m.OpenTag_was_script = false
  m.OpenTag_was_style  = false
  found_elem_name := false

  for ; l<m.p_fb.NumLines(); l++ {
    var p_fl *FLine = m.p_fb.GetLP( l )
    LL := m.p_fb.LineLen( l )
    for ; p<LL; p++ {
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++ // Move past '>'
        m.state = St_In_None
        if( m.OpenTag_was_script ) {
          m.state = St_JS_None
          m.JS_edges.Push( Edges{ CrsPos{l,p}, CrsPos{0,0} } )
        } else if( m.OpenTag_was_style ) {
          m.state = St_CS_None
          m.CS_edges.Push( Edges{ CrsPos{l,p}, CrsPos{0,0} } )
        }
      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( !found_elem_name ) {
        // Returns non-zero if lp has HTTP tag at p:
        tag_len := m.Has_HTTP_Tag_At( p_fl, p )
        if( 0<tag_len ) {
          if( p_fl.has_at_ci("script", p) ) {
            m.OpenTag_was_script = true
          } else if( p_fl.has_at("style", p) ) {
            m.OpenTag_was_style = true
          }
          found_elem_name = true
          for k:=0; k<tag_len; k++ {
            m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
            p++
          }
          p--
        } else if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else { //( found_elem_name )
        if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
          p++ //< Move past white space
          m.state = St_OpenTag_AttrName
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      }
      if( St_OpenTag_ElemName != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_OpenTag_AttrName( l, p int ) (int,int) {

  found_attr_name := false
  past__attr_name := false

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )
    for ; p<LL; p++ {
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++ // Move past '>'
        if       ( m.OpenTag_was_style  ) { m.state = St_CS_None
        } else if( m.OpenTag_was_script ) { m.state = St_JS_None
        } else                            { m.state = St_In_None
        }
      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( !found_attr_name ) {
        if( IsXML_Ident( c0 ) ) {
          found_attr_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else if( found_attr_name && !past__attr_name ) {
        if( IsXML_Ident( c0 ) ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else if( c0==' ' || c0=='\t' ) {
          past__attr_name = true
        //m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )
        } else if( c0=='=' ) {
          past__attr_name = true
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
          p++ //< Move past '='
          m.state = St_OpenTag_AttrVal
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      } else { //( found_attr_name && past__attr_name )
        if( c0=='=' ) {
          m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
          p++ //< Move past '='
          m.state = St_OpenTag_AttrVal
        } else if( c0==' ' || c0=='\t' ) {
        //m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )
        } else {
        //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
        }
      }
      if( St_OpenTag_AttrName != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_OpenTag_AttrVal( l, p int ) (int,int) {

  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )
    for ; p<LL; p++ {
      var c0 rune = m.p_fb.GetR( l, p )

      if( c0=='>' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
        p++ // Move past '>'
        if       ( m.OpenTag_was_style  ) { m.state = St_CS_None
        } else if( m.OpenTag_was_script ) { m.state = St_JS_None
        } else                            { m.state = St_In_None
        }
      } else if( c0=='/' || c0=='?' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )

      } else if( c0=='\'' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++// Move past '\''
        m.state = St_SingleQuote
        m.qtXSt = St_OpenTag_AttrName

      } else if( c0=='"' ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = St_DoubleQuote
        m.qtXSt = St_OpenTag_AttrName

      } else if( IsDigit( c0 ) ) {
        m.state = St_NumberBeg
        m.numXSt = St_OpenTag_AttrName

      } else if( c0==' ' || c0=='\t' ) {
      //m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
      } else {
      //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( St_OpenTag_AttrVal != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_SingleQuote( l, p int ) (int,int) {

  exit := false
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='\'') ||
          (c1!='\\' && c0=='\'') ||
          (c1=='\\' && c0=='\'' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ // Move past '\''
        m.state = m.qtXSt
        exit = true
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else                     { slash_escaped = false
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( exit ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_DoubleQuote( l, p int ) (int,int) {

  exit := false
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )
    slash_escaped := false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( (c1==0    && c0=='"') ||
          (c1!='\\' && c0=='"') ||
          (c1=='\\' && c0=='"' && slash_escaped) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = m.qtXSt
        exit = true
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else                     { slash_escaped = false
        }
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( exit ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func JS_OneVarType( c0 rune ) bool {
  return c0=='&' ||
         c0=='.' || c0=='*' ||
         c0=='[' || c0==']'
}

func JS_OneControl( c0 rune ) bool {
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

func JS_TwoControl( c1, c0 rune ) bool {
  return (c1=='=' && c0=='=') ||
         (c1=='&' && c0=='&') ||
         (c1=='|' && c0=='|') ||
         (c1=='|' && c0=='=') ||
         (c1=='&' && c0=='=') ||
         (c1=='!' && c0=='=') ||
         (c1=='+' && c0=='=') ||
         (c1=='-' && c0=='=')
}

func (m *Highlight_HTML) Hi_JS_None( l, p int ) (int,int) {

  for ; l<m.p_fb.NumLines(); l++ {
    var p_fl *FLine = m.p_fb.GetLP( l )
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if       ( c1=='/' && c0 == '/' ) { p--; m.state = St_JS_CPP_Comment
      } else if( c1=='/' && c0 == '*' ) {
        p--
        m.state = St_JS_C_Comment
        m.ccXSt = St_JS_None

      } else if( c0 == '#' ) { m.state = St_JS_Define
      } else if( c0 == '\'') {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = St_JS_SingleQuote
        m.qtXSt = St_JS_None

      } else if( c0 == '"') {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = St_JS_DoubleQuote
        m.qtXSt = St_JS_None

      } else if( !IsIdent(c1) && IsDigit(c0) ) {
        m.state = St_NumberBeg
        m.numXSt = St_JS_None

      } else if( (c1==':' && c0==':') ||
                 (c1=='-' && c0=='>') ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_VARTYPE )
        m.p_fb.SetSyntaxStyle( l, p  , HI_VARTYPE )

      } else if( c1=='<' && c0=='/' && p+7<LL ) {
        if( p_fl.has_at_ci("</script", p-1) ) {
          m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE )
          m.p_fb.SetSyntaxStyle( l, p  , HI_DEFINE )
          p++ // Move past '/'
          m.state = St_CloseTag
          if( 0<m.JS_edges.Len() ) { //< Should always be true
            var edges Edges = m.JS_edges.Get( m.JS_edges.Len()-1 )
            edges.end.crsLine = l
            edges.end.crsChar = p-1
            m.JS_edges.Set( m.JS_edges.Len()-1, edges )
          }
        }
      } else if( JS_TwoControl( c1, c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p-1, HI_CONTROL )
        m.p_fb.SetSyntaxStyle( l, p  , HI_CONTROL )

      } else if( JS_OneVarType( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )

      } else if( JS_OneControl( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )

      } else if( c0 < 32 || 126 < c0 ) {
      //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( St_JS_None != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_JS_Define( l, p int ) (int,int) {

    LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    // c0 is ahead of c1: (c1,c0)
    var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
    var c0 rune =                     m.p_fb.GetR( l, p )

    if( c1=='/' && c0=='/' ) {
      m.p_fb.SetSyntaxStyle( l, p-1, HI_COMMENT )
      m.p_fb.SetSyntaxStyle( l, p  , HI_COMMENT )
      p++
      m.state = St_JS_CPP_Comment

    } else if( c1=='/' && c0=='*' ) {
      m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
      p++
      m.state = St_JS_C_Comment
      m.ccXSt = St_JS_None
    } else {
      m.p_fb.SetSyntaxStyle( l, p, HI_DEFINE )
    }
    if( St_JS_Define != m.state ) { return l,p }
  }
  p=0; l++
  m.state = St_JS_None
  return l,p
}

func (m *Highlight_HTML) Hi_C_Comment( l, p int ) (int,int) {

  exit := false
  for ; l<m.p_fb.NumLines(); l++ {
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )

      if( c1=='*' && c0=='/' ) {
        p++ //< Move past '/'
        m.state = m.ccXSt
        exit = true
      }
      if( exit ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Hi_JS_CPP_Comment( l, p int ) (int,int) {

  LL := m.p_fb.LineLen( l )

  for ; p<LL; p++ {
    m.p_fb.SetSyntaxStyle( l, p, HI_COMMENT )
  }
  l++
  p=0
  m.state = St_JS_None
  return l,p
}

func CS_OneVarType( c0 rune ) bool {
  return c0=='*' || c0=='#'
}
func CS_OneControl( c0 rune ) bool {
  return c0=='.' || c0=='-' || c0==',' ||
         c0==':' || c0==';' ||
         c0=='{' || c0=='}'
}

func (m *Highlight_HTML) Hi_CS_None( l, p int ) (int,int) {

  for ; l<m.p_fb.NumLines(); l++ {
    var p_fl *FLine = m.p_fb.GetLP( l )
    LL := m.p_fb.LineLen( l )

    for ; p<LL; p++ {
      m.p_fb.ClearSyntaxStyles( l, p )

      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = m.p_fb.GetR( l, p-1 ) }
      var c0 rune =                     m.p_fb.GetR( l, p )

      if( c1=='/' && c0 == '*' ) {
        p--
        m.state = St_CS_C_Comment
        m.ccXSt = St_CS_None

      } else if( c0 == '\'') {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = St_CS_SingleQuote
        m.qtXSt = St_CS_None

      } else if( c0 == '"') {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++ //< Move past '\"'
        m.state = St_CS_DoubleQuote
        m.qtXSt = St_CS_None

      } else if( !IsIdent( c1 ) && IsDigit( c0 ) ) {
        // FIXME: For CSS, the following extensions should be highlighted:
        //        px, pt, %, in, em
        m.state = St_NumberBeg
        m.numXSt = St_CS_None

      } else if( c1=='<' && c0=='/' && p+6<LL ) {
        if( p_fl.has_at_ci("</style", p-1) ) {
          m.p_fb.SetSyntaxStyle( l, p-1, HI_DEFINE )
          m.p_fb.SetSyntaxStyle( l, p  , HI_DEFINE )
          p++ // Move past '/'
          m.state = St_CloseTag
          if( 0<m.CS_edges.Len() ) { //< Should always be true
            var edges Edges = m.CS_edges.Get( m.CS_edges.Len()-1 )
            edges.end.crsLine = l
            edges.end.crsChar = p-1
            m.CS_edges.Set( m.CS_edges.Len()-1, edges )
          }
        }
      } else if( CS_OneVarType( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_VARTYPE )

      } else if( CS_OneControl( c0 ) ) {
        m.p_fb.SetSyntaxStyle( l, p, HI_CONTROL )

      } else if( c0 < 32 || 126 < c0 ) {
      //m.p_fb.SetSyntaxStyle( l, p, HI_NONASCII )
      }
      if( St_CS_None != m.state ) { return l,p }
    }
    p = 0
  }
  m.state = St_Done
  return l,p
}

func (m *Highlight_HTML) Has_HTTP_Tag_At( lp *FLine, pos int ) int {

  if( IsXML_Ident( lp.GetR( pos ) ) ) {
    num_HTML_Tags := len(HTML_Tags)

    for k:=0; k<num_HTML_Tags; k++ {
      tag := HTML_Tags[k]

      if( lp.has_at_ci(tag, pos) ) {
        return len(tag)
      }
    }
  }
  return 0
}

