
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

