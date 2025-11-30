
package main

import (
//"fmt"
)

type Highlight_CPP struct {
  hi_code Highlight_Code
}

func (m *Highlight_CPP) Init( p_fb *FileBuf ) {
  m.hi_code.Init( p_fb )
}

func (m *Highlight_CPP) Run_Range( st CrsPos, fn int ) {

  m.hi_code.state = m.hi_code.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.hi_code.state && l<fn {
    l,p = m.hi_code.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_CPP) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.hi_code.p_fb, HiPairs_CPP[:], st, fn )
}

var HiPairs_CPP = [...]HiKeyVal {
  { "if"                 , HI_CONTROL },
  { "else"               , HI_CONTROL },
  { "for"                , HI_CONTROL },
  { "while"              , HI_CONTROL },
  { "do"                 , HI_CONTROL },
  { "return"             , HI_CONTROL },
  { "switch"             , HI_CONTROL },
  { "case"               , HI_CONTROL },
  { "break"              , HI_CONTROL },
  { "default"            , HI_CONTROL },
  { "continue"           , HI_CONTROL },
  { "template"           , HI_CONTROL },
  { "public"             , HI_CONTROL },
  { "protected"          , HI_CONTROL },
  { "private"            , HI_CONTROL },
  { "typedef"            , HI_CONTROL },
  { "delete"             , HI_CONTROL },
  { "operator"           , HI_CONTROL },
  { "sizeof"             , HI_CONTROL },
  { "using"              , HI_CONTROL },
  { "namespace"          , HI_CONTROL },
  { "goto"               , HI_CONTROL },
  { "friend"             , HI_CONTROL },
  { "try"                , HI_CONTROL },
  { "catch"              , HI_CONTROL },
  { "throw"              , HI_CONTROL },
  { "and"                , HI_CONTROL },
  { "or"                 , HI_CONTROL },
  { "not"                , HI_CONTROL },
  { "new"                , HI_CONTROL },
  { "const_cast"         , HI_CONTROL },
  { "static_cast"        , HI_CONTROL },
  { "dynamic_cast"       , HI_CONTROL },
  { "reinterpret_cast"   , HI_CONTROL },
  { "override"           , HI_CONTROL },

  // Types
  { "auto"               , HI_VARTYPE },
  { "int"                , HI_VARTYPE },
  { "long"               , HI_VARTYPE },
  { "void"               , HI_VARTYPE },
  { "this"               , HI_VARTYPE },
  { "bool"               , HI_VARTYPE },
  { "char"               , HI_VARTYPE },
  { "const"              , HI_VARTYPE },
  { "constexpr"          , HI_VARTYPE },
  { "short"              , HI_VARTYPE },
  { "float"              , HI_VARTYPE },
  { "double"             , HI_VARTYPE },
  { "signed"             , HI_VARTYPE },
  { "unsigned"           , HI_VARTYPE },
  { "extern"             , HI_VARTYPE },
  { "static"             , HI_VARTYPE },
  { "enum"               , HI_VARTYPE },
  { "uint8_t"            , HI_VARTYPE },
  { "uint16_t"           , HI_VARTYPE },
  { "uint32_t"           , HI_VARTYPE },
  { "uint64_t"           , HI_VARTYPE },
  { "size_t"             , HI_VARTYPE },
  { "int8_t"             , HI_VARTYPE },
  { "int16_t"            , HI_VARTYPE },
  { "int32_t"            , HI_VARTYPE },
  { "int64_t"            , HI_VARTYPE },
  { "float32_t"          , HI_VARTYPE },
  { "float64_t"          , HI_VARTYPE },
  { "FILE"               , HI_VARTYPE },
  { "DIR"                , HI_VARTYPE },
  { "class"              , HI_VARTYPE },
  { "struct"             , HI_VARTYPE },
  { "union"              , HI_VARTYPE },
  { "typename"           , HI_VARTYPE },
  { "virtual"            , HI_VARTYPE },
  { "inline"             , HI_VARTYPE },
  { "explicit"           , HI_VARTYPE },

  // Constants
  { "true"               , HI_CONST   },
  { "false"              , HI_CONST   },
  { "NULL"               , HI_CONST   },
  { "nullptr"            , HI_CONST   },

  { "__FUNCTION__"       , HI_DEFINE  },
  { "__PRETTY_FUNCTION__", HI_DEFINE  },
  { "__TIMESTAMP__"      , HI_DEFINE  },
  { "__FILE__"           , HI_DEFINE  },
  { "__func__"           , HI_DEFINE  },
  { "__LINE__"           , HI_DEFINE  },
  { ""                   , 0 },
}

