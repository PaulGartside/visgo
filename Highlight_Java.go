
package main

import (
//"fmt"
)

type Highlight_Java struct {
  hi_code Highlight_Code
}

func (m *Highlight_Java) Init( p_fb *FileBuf ) {
  m.hi_code.Init( p_fb )
}

func (m *Highlight_Java) Run_Range( st CrsPos, fn int ) {

  m.hi_code.state = m.hi_code.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.hi_code.state && l<fn {
    l,p = m.hi_code.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_Java) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.hi_code.p_fb, HiPairs_Java[:], st, fn )
}

var HiPairs_Java = [...]HiKeyVal {
  { "abstract"           , HI_VARTYPE },
  { "boolean"            , HI_VARTYPE },
  { "Boolean"            , HI_VARTYPE },
  { "break"              , HI_CONTROL },
  { "byte"               , HI_VARTYPE },
  { "Byte"               , HI_VARTYPE },
  { "case"               , HI_CONTROL },
  { "catch"              , HI_CONTROL },
  { "char"               , HI_VARTYPE },
  { "Character"          , HI_VARTYPE },
  { "class"              , HI_VARTYPE },
  { "const"              , HI_VARTYPE },
  { "continue"           , HI_CONTROL },
  { "default"            , HI_CONTROL },
  { "do"                 , HI_CONTROL },
  { "double"             , HI_VARTYPE },
  { "Double"             , HI_VARTYPE },
  { "else"               , HI_CONTROL },
  { "enum"               , HI_VARTYPE },
  { "extends"            , HI_VARTYPE },
  { "final"              , HI_VARTYPE },
  { "float"              , HI_VARTYPE },
  { "Float"              , HI_VARTYPE },
  { "finally"            , HI_CONTROL },
  { "for"                , HI_CONTROL },
  { "goto"               , HI_CONTROL },
  { "if"                 , HI_CONTROL },
  { "implements"         , HI_VARTYPE },
  { "import"             , HI_DEFINE  },
  { "instanceof"         , HI_VARTYPE },
  { "int"                , HI_VARTYPE },
  { "Integer"            , HI_VARTYPE },
  { "interface"          , HI_VARTYPE },
  { "Iterator"           , HI_VARTYPE },
  { "long"               , HI_VARTYPE },
  { "Long"               , HI_VARTYPE },
  { "main"               , HI_DEFINE  },
  { "native"             , HI_VARTYPE },
  { "new"                , HI_VARTYPE },
  { "package"            , HI_DEFINE  },
  { "private"            , HI_CONTROL },
  { "protected"          , HI_CONTROL },
  { "public"             , HI_CONTROL },
  { "return"             , HI_CONTROL },
  { "short"              , HI_VARTYPE },
  { "Short"              , HI_VARTYPE },
  { "static"             , HI_VARTYPE },
  { "strictfp"           , HI_VARTYPE },
  { "String"             , HI_VARTYPE },
  { "System"             , HI_DEFINE  },
  { "super"              , HI_VARTYPE },
  { "switch"             , HI_CONTROL },
  { "synchronized"       , HI_CONTROL },
  { "this"               , HI_VARTYPE },
  { "throw"              , HI_CONTROL },
  { "throws"             , HI_CONTROL },
  { "transient"          , HI_VARTYPE },
  { "try"                , HI_CONTROL },
  { "void"               , HI_VARTYPE },
  { "Void"               , HI_VARTYPE },
  { "volatile"           , HI_VARTYPE },
  { "while"              , HI_CONTROL },
  { "virtual"            , HI_VARTYPE },
  { "true"               , HI_CONST   },
  { "false"              , HI_CONST   },
  { "null"               , HI_CONST   },
  { "@Deprecated"        , HI_DEFINE  },
  { "@Override"          , HI_DEFINE  },
  { "@SuppressWarnings"  , HI_DEFINE  },
  { ""                   , 0 },
}

