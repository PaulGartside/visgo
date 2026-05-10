
package main

import (
//"fmt"
)

const seq_len = 7

var seq_inc [seq_len] byte = [seq_len]byte{  79, 101, 127, 139, 163, 181, 199 }
var seq_mod [seq_len] byte = [seq_len]byte{ 131, 151, 173, 191, 211, 229, 251 }
var seq_val [seq_len] byte 

// Contents of p_in are covered and put into p_out
func Cover_Array( p_in *FileBuf, p_out *Vector[byte], seed byte, key string ) {

  Init_seq_val( seed, key )

  p_out.Clear()
  NUM_LINES := p_in.NumLines()

  for l:=0; l<NUM_LINES; l++ {
    LS := p_in.LineSize( l )

    for p:=0; p<LS; p++ {
      var B byte = p_in.GetB( l, p )
      var C byte = Cover_Byte() ^ B
      p_out.Push( C )
    }
    if( l<NUM_LINES-1 || p_in.Has_LF_at_EOF() ) {
      var C byte = Cover_Byte() ^ '\n'
      p_out.Push( C )
    }
  }
}

func Init_seq_val( seed byte, key string ) {
  // Initialize seq_val:
  key_len := len( key )

  for k:=0; k<seq_len; k++ {
    seq_val[k] = byte( (int(seed) + int(seq_inc[k])) % int(seq_mod[k]) )
  //var sum int = int(seed) + int(seq_inc[k])
  //seq_val[k] = byte(sum % int(seq_mod[k]))
  }
  for k:=0; k<seq_len*key_len; k+=1 {
    k_m := k % seq_len

    seq_val[ k_m ] ^= key[ k % key_len ]
    seq_val[ k_m ] %= seq_mod[ k_m ]
  }
}

func Cover_Byte() byte {
  var cb byte = 0xAA
  for k:=0; k<seq_len; k+=1 {
    seq_val[k] = byte( (int(seq_val[k]) + int(seq_inc[k])) % int(seq_mod[k]) )
  //var sum int = int(seq_val[k]) + int(seq_inc[k])
  //seq_val[k] = byte(sum % int(seq_mod[k]))

    cb ^= seq_val[k]
  }
  return cb
}

