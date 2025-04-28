
package main

type CrsPos struct {
  crsLine int
  crsChar int
}

//type HighlightType byte

// HI for Highlight
const (
  HI_STAR      byte = 0x01 // Search pattern
  HI_STAR_IN_F = 0x02 // Search pattern in file
  HI_COMMENT   = 0x04 // Comment
  HI_DEFINE    = 0x08 // #define
  HI_CONST     = 0x10 // C constant, '...' or "...", true, false
  HI_CONTROL   = 0x20 // C flow control, if, else, etc.
  HI_VARTYPE   = 0x40 // C variable type, char, int, etc.
  HI_NONASCII  = 0x80 // Non-ascii character
)

type HiKeyVal struct {
  key string
  val byte
}

