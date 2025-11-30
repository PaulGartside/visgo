
package main

import (
//"fmt"
)

type Highlight_IDL struct {
  hi_code Highlight_Code
}

func (m *Highlight_IDL) Init( p_fb *FileBuf ) {
  m.hi_code.p_fb = p_fb
}

func (m *Highlight_IDL) Run_Range( st CrsPos, fn int ) {

  m.hi_code.state = m.hi_code.Hi_In_None

  l := st.crsLine
  p := st.crsChar

  for nil != m.hi_code.state && l<fn {
    l,p = m.hi_code.state( l, p )
  }
  m.Find_Styles_Keys_In_Range( st, fn )
}

func (m *Highlight_IDL) Find_Styles_Keys_In_Range( st CrsPos, fn int ) {

  Hi_FindKey_In_Range( m.hi_code.p_fb, HiPairs_IDL[:], st, fn )
}

var HiPairs_IDL = [...]HiKeyVal {
  { "abstract"   , HI_VARTYPE },
  { "any"        , HI_VARTYPE },
  { "attribute"  , HI_VARTYPE },
  { "boolean"    , HI_VARTYPE },
  { "case"       , HI_CONTROL },
  { "char"       , HI_VARTYPE },
  { "component"  , HI_CONTROL },
  { "const"      , HI_VARTYPE },
  { "consumes"   , HI_CONTROL },
  { "context"    , HI_CONTROL },
  { "custom"     , HI_CONTROL },
  { "default"    , HI_CONTROL },
  { "double"     , HI_VARTYPE },
  { "exception"  , HI_CONTROL },
  { "emits"      , HI_CONTROL },
  { "enum"       , HI_VARTYPE },
  { "eventtype"  , HI_CONTROL },
  { "factory"    , HI_CONTROL },
  { "FALSE"      , HI_CONST   },
  { "finder"     , HI_CONTROL },
  { "fixed"      , HI_CONTROL },
  { "float"      , HI_CONTROL },
  { "getraises"  , HI_CONTROL },
  { "home"       , HI_CONTROL },
  { "import"     , HI_DEFINE  },
  { "in"         , HI_CONTROL },
  { "inout"      , HI_CONTROL },
  { "interface"  , HI_VARTYPE },
  { "local"      , HI_VARTYPE },
  { "long"       , HI_VARTYPE },
  { "module"     , HI_VARTYPE },
  { "multiple"   , HI_CONTROL },
  { "native"     , HI_VARTYPE },
  { "Object"     , HI_VARTYPE },
  { "octet"      , HI_VARTYPE },
  { "oneway"     , HI_CONTROL },
  { "out"        , HI_CONTROL },
  { "primarykey" , HI_VARTYPE },
  { "private"    , HI_CONTROL },
  { "provides"   , HI_CONTROL },
  { "public"     , HI_CONTROL },
  { "publishes"  , HI_CONTROL },
  { "raises"     , HI_CONTROL },
  { "readonly"   , HI_VARTYPE },
  { "setraises"  , HI_CONTROL },
  { "sequence"   , HI_CONTROL },
  { "short"      , HI_VARTYPE },
  { "string"     , HI_VARTYPE },
  { "struct"     , HI_VARTYPE },
  { "supports"   , HI_CONTROL },
  { "switch"     , HI_CONTROL },
  { "TRUE"       , HI_CONST   },
  { "truncatable", HI_CONTROL },
  { "typedef"    , HI_VARTYPE },
  { "typeid"     , HI_VARTYPE },
  { "typeprefix" , HI_VARTYPE },
  { "unsigned"   , HI_VARTYPE },
  { "union"      , HI_VARTYPE },
  { "uses"       , HI_CONTROL },
  { "ValueBase"  , HI_VARTYPE },
  { "valuetype"  , HI_VARTYPE },
  { "void"       , HI_VARTYPE },
  { "wchar"      , HI_VARTYPE },
  { "wstring"    , HI_VARTYPE },
  { ""           , 0 },
}

