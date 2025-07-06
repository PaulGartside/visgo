
package main

type Highlight_Base interface {
  Init( p_fb *FileBuf )
  Run_Range( st CrsPos, fn int )
}

type HiStateFunc func( l, p int ) (int, int)

func Hi_FindKey_In_Range( p_fb *FileBuf, HiPairs []HiKeyVal, st CrsPos, fn int ) {

  if( nil == p_fb ) {
    return
  }
  var NUM_LINES int = p_fb.NumLines()

  for l:=st.crsLine; l<=fn && l<NUM_LINES; l++ {
    var lp *FLine = p_fb.GetLP( l )
    var lr RLine = lp.runes
    var sr BLine = lp.styles
    var LL int = lp.Len()

    var st_pos int = True_1_else_2( st.crsLine==l, st.crsChar, 0 )
    var fn_pos int = LLM1( LL )

    for p:=st_pos; p<=fn_pos && p<LL; p++ {
      var key_st bool = 0==sr.GetB(p) && line_start_or_prev_C_non_ident( lr, p )

      for h:=0; key_st && ""!=HiPairs[h].key; h++ {
        var matches bool = true
        key     := HiPairs[h].key
        HI_TYPE := HiPairs[h].val
        KEY_LEN := len( key )

        for k:=0; matches && (p+k)<LL && k<KEY_LEN; k++ {
          if( 0!=sr.GetB(p+k) || rune(key[k]) != lr.GetR(p+k) ) { matches = false
          } else {
            if( k+1 == KEY_LEN ) { // Found pattern
              matches = line_end_or_non_ident( lr, LL, p+k )
              if( matches ) {
                for m:=p; m<p+KEY_LEN; m++ { p_fb.SetSyntaxStyle( l, m, HI_TYPE ) }
                // Increment p one less than KEY_LEN, because p
                // will be incremented again by the for loop
                p += KEY_LEN-1
                // Set key_st to false here to break out of h for loop
                key_st = false
              }
            }
          }
        }
      }
    }
  }
}

// This shows one way to re-use class methods in Go:
//
func Hi_In_SingleQuote_CPP_Go(
     l, p int, p_fb *FileBuf, Hi_In_None HiStateFunc ) (
     int,int,HiStateFunc) {

  var state HiStateFunc = nil
  p_fb.SetSyntaxStyle( l, p, HI_CONST )
  p++
  for ; l<p_fb.NumLines(); l++ {
    LL := p_fb.LineLen( l )

    var slash_escaped bool = false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p ) { c1 = p_fb.GetR( l, p-1 ) }
      var c0 rune =                     p_fb.GetR( l, p )

      if( (c1==0    && c0=='\'') ||
          (c1!='\\' && c0=='\'') ||
          (c1=='\\' && c0=='\'' && slash_escaped) ) {
        // End of single quote:
        p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++
        state = Hi_In_None
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else                     { slash_escaped = false
        }
        p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( nil != state ) { return l,p, state }
    }
    p = 0
  }
  return l,p, state
}

// This shows one way to re-use class methods in Go:
//
func Hi_In_DoubleQuote_CPP_Go(
     l, p int, p_fb *FileBuf, Hi_In_None HiStateFunc ) (
     int,int,HiStateFunc) {

  var state HiStateFunc = nil
  p_fb.SetSyntaxStyle( l, p, HI_CONST )
  p++
  for ; l<p_fb.NumLines(); l++ {
    LL := p_fb.LineLen( l )

    var slash_escaped bool = false
    for ; p<LL; p++ {
      // c0 is ahead of c1: (c1,c0)
      var c1 rune = 0; if( 0<p) { c1 = p_fb.GetR( l, p-1 ) }
      var c0 rune =                    p_fb.GetR( l, p )

      if( (c1==0    && c0=='"') ||
          (c1!='\\' && c0=='"') ||
          (c1=='\\' && c0=='"' && slash_escaped) ) {
        p_fb.SetSyntaxStyle( l, p, HI_CONST )
        p++
        state = Hi_In_None
      } else {
        if( c1=='\\' && c0=='\\' ) { slash_escaped = !slash_escaped
        } else {                     slash_escaped = false
        }
        p_fb.SetSyntaxStyle( l, p, HI_CONST )
      }
      if( nil != state ) { return l,p, state }
    }
    p = 0
  }
  return l,p, state
}

