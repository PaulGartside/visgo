
package main

import (
//  "fmt"
)

func (m *Vis) VSplitWindow() {

//NoDiff_CV(m)

  var p_cv *FileView = m.CV()
  var cv_tp Tile_Pos = p_cv.GetTilePos()

  // Make sure current view can be vertically split:
  if( m.num_wins < MAX_WINS &&
      ( cv_tp == TP_FULL ||
        cv_tp == TP_TOP__HALF ||
        cv_tp == TP_BOT__HALF ||
        cv_tp == TP_LEFT_HALF ||
        cv_tp == TP_RITE_HALF ||
        cv_tp == TP_TOP__LEFT_QTR ||
        cv_tp == TP_BOT__LEFT_QTR ||
        cv_tp == TP_TOP__RITE_QTR ||
        cv_tp == TP_BOT__RITE_QTR ) ) {
    // New window will be m.num_wins.
    // Duplicate file hist of current window into new window.
    m.file_hist[m.num_wins].CopyP( &m.file_hist[m.win] )

    // Copy current view context into new view
    var p_nv *FileView = m.GetView_Win( m.num_wins )

    p_nv.Set_Context( p_cv )

    // Make new window the current window:
    m.win = m.num_wins
    m.num_wins++

    // Set the new tile positions of the old view p_cv, and the new view p_nv:
    if( cv_tp == TP_FULL ) {
       p_cv.SetTilePos( TP_LEFT_HALF )
       p_nv.SetTilePos( TP_RITE_HALF )
    } else if( cv_tp == TP_TOP__HALF ) {
      p_cv.SetTilePos( TP_TOP__LEFT_QTR )
      p_nv.SetTilePos( TP_TOP__RITE_QTR )
    } else if( cv_tp == TP_BOT__HALF ) {
      p_cv.SetTilePos( TP_BOT__LEFT_QTR )
      p_nv.SetTilePos( TP_BOT__RITE_QTR )
    } else if( cv_tp == TP_LEFT_HALF ) {
      p_cv.SetTilePos( TP_LEFT_QTR )
      p_nv.SetTilePos( TP_LEFT_CTR__QTR )
    } else if( cv_tp == TP_RITE_HALF ) {
      p_cv.SetTilePos( TP_RITE_CTR__QTR )
      p_nv.SetTilePos( TP_RITE_QTR )
    } else if( cv_tp == TP_TOP__LEFT_QTR ) {
      p_cv.SetTilePos( TP_TOP__LEFT_8TH )
      p_nv.SetTilePos( TP_TOP__LEFT_CTR_8TH )
    } else if( cv_tp == TP_BOT__LEFT_QTR ) {
      p_cv.SetTilePos( TP_BOT__LEFT_8TH )
      p_nv.SetTilePos( TP_BOT__LEFT_CTR_8TH )
    } else if( cv_tp == TP_TOP__RITE_QTR ) {
      p_cv.SetTilePos( TP_TOP__RITE_CTR_8TH )
      p_nv.SetTilePos( TP_TOP__RITE_8TH )
    } else { //( cv_tp == TP_BOT__RITE_QTR )
      p_cv.SetTilePos( TP_BOT__RITE_CTR_8TH )
      p_nv.SetTilePos( TP_BOT__RITE_8TH )
    }
  } else if( m.num_wins+1 < MAX_WINS &&
             ( cv_tp == TP_LEFT_TWO_THIRDS ||
               cv_tp == TP_RITE_TWO_THIRDS ) ) {
    m.file_hist[m.num_wins].CopyP( &m.file_hist[m.win] )

    // Copy current view context into new view
    var p_nv *FileView = m.GetView_Win( m.num_wins )

    p_nv.Set_Context( p_cv )

    // Make new window the current window:
    m.num_wins += 1

    // Set the new tile positions.
    if( cv_tp == TP_LEFT_TWO_THIRDS ) {
      p_cv.SetTilePos( TP_LEFT_THIRD )
      p_nv.SetTilePos( TP_CTR__THIRD )
    } else { //( cv_tp == TP_RITE_TWO_THIRDS )
      p_cv.SetTilePos( TP_CTR__THIRD )
      p_nv.SetTilePos( TP_RITE_THIRD )
    }
  }
  m.UpdateViews( false )
}

func (m *Vis) HSplitWindow() {

//NoDiff_CV(m)

  var p_cv *FileView = m.CV()
  var cv_tp Tile_Pos = p_cv.GetTilePos()

  // Make sure current view can be horizontally split:
  if( m.num_wins < MAX_WINS  &&
      ( cv_tp == TP_FULL ||
        cv_tp == TP_LEFT_HALF ||
        cv_tp == TP_RITE_HALF ||
        cv_tp == TP_LEFT_QTR  ||
        cv_tp == TP_RITE_QTR  ||
        cv_tp == TP_LEFT_CTR__QTR ||
        cv_tp == TP_RITE_CTR__QTR ) ) {
    // New window will be m.num_wins.
    // Duplicate file hist of current window into new window.
    m.file_hist[m.num_wins].CopyP( &m.file_hist[m.win] )

    // Copy current view context into new view
    var p_nv *FileView = m.GetView_Win( m.num_wins )

    p_nv.Set_Context( p_cv )

    // Make new window the current window:
    m.win = m.num_wins
    m.num_wins++

    // Set the new tile positions of the old view cv, and the new view nv:
    if( cv_tp == TP_FULL ) {
      p_cv.SetTilePos( TP_TOP__HALF )
      p_nv.SetTilePos( TP_BOT__HALF )
    } else if( cv_tp == TP_LEFT_HALF ) {
      p_cv.SetTilePos( TP_TOP__LEFT_QTR )
      p_nv.SetTilePos( TP_BOT__LEFT_QTR )
    } else if( cv_tp == TP_RITE_HALF ) {
      p_cv.SetTilePos( TP_TOP__RITE_QTR )
      p_nv.SetTilePos( TP_BOT__RITE_QTR )
    } else if( cv_tp == TP_LEFT_QTR ) {
      p_cv.SetTilePos( TP_TOP__LEFT_8TH )
      p_nv.SetTilePos( TP_BOT__LEFT_8TH )
    } else if( cv_tp == TP_RITE_QTR ) {
      p_cv.SetTilePos( TP_TOP__RITE_8TH )
      p_nv.SetTilePos( TP_BOT__RITE_8TH )
    } else if( cv_tp == TP_LEFT_CTR__QTR ) {
      p_cv.SetTilePos( TP_TOP__LEFT_CTR_8TH )
      p_nv.SetTilePos( TP_BOT__LEFT_CTR_8TH )
    } else { //( cv_tp == TP_RITE_CTR__QTR )
      p_cv.SetTilePos( TP_TOP__RITE_CTR_8TH )
      p_nv.SetTilePos( TP_BOT__RITE_CTR_8TH )
    }
  }
  m.UpdateViews( false )
}

func (m *Vis) GoToNextWindow() {

  if( 1 < m.num_wins ) {

    var win_old int = m.win

    m.win++
    m.win = m.win % m.num_wins

    var pV2     *FileView = m.GetView_Win( m.win )
    var pV2_old *FileView = m.GetView_Win( win_old )

    pV2_old.PrintBorders()
    pV2    .PrintBorders()

  //Console::Update()

    pV2.PrintCursor()
  }
}

func (m *Vis) GoToNextWindow_l() {

  if( 1 < m.num_wins ) {
    var win_old int = m.win

    // If next view to go to was not found, dont do anything, just return
    // If next view to go to is found, m.win will be updated to new value
    if( m.GoToNextWindow_l_Find() ) {
      var pV     *FileView = m.GetView_Win( m.win )
      var pV_old *FileView = m.GetView_Win( win_old )

      pV_old.PrintBorders()
      pV    .PrintBorders()

    //Console::Update()

      pV.PrintCursor()
    }
  }
}

func (m *Vis) GoToNextWindow_h() {

  if( 1 < m.num_wins ) {
    var win_old int = m.win

    // If next view to go to was not found, dont do anything, just return
    // If next view to go to is found, m.win will be updated to new value
    if( m.GoToNextWindow_h_Find() ) {
      var pV     *FileView = m.GetView_Win( m.win   )
      var pV_old *FileView = m.GetView_Win( win_old )

      pV_old.PrintBorders()
      pV    .PrintBorders()

    //Console::Update()

      pV.PrintCursor()
    }
  }
}

func (m *Vis) GoToNextWindow_jk() {

  if( 1 < m.num_wins ) {
    var win_old int = m.win

    // If next view to go to was not found, dont do anything, just return
    // If next view to go to is found, m.win will be updated to new value
    if( m.GoToNextWindow_jk_Find() ) {
      var pV     *FileView = m.GetView_Win( m.win   )
      var pV_old *FileView = m.GetView_Win( win_old )

      pV_old.PrintBorders()
      pV    .PrintBorders()

    //Console::Update()

      pV.PrintCursor()
    }
  }
}

func (m *Vis) FlipWindows() {

  if( 1 < m.num_wins ) {
    var split_horizontally bool = false

    for k:=0; !split_horizontally && k<m.num_wins; k++ {
      // pV is View of displayed window k
      var pV *FileView = m.GetView_Win( k )

      split_horizontally = pV.GetTilePos() == TP_TOP__HALF ||
                           pV.GetTilePos() == TP_BOT__HALF
    }
    for k:=0; k<m.num_wins; k++ {
      // pV is View of displayed window k
      var pV  *FileView = m.GetView_Win( k )
      var OTP Tile_Pos = pV.GetTilePos(); // Old tile position

      // New tile position:
      var NTP Tile_Pos
      if( split_horizontally ) { NTP = FlipWindows_Vertically( OTP )
      } else                   { NTP = FlipWindows_Horizontally( OTP )
      }
      if( NTP != TP_NONE ) { pV.SetTilePos( NTP )
      }
    }
    m.UpdateViews( false )
  }
}

func (m *Vis) Quit_JoinTiles_LEFT_HALF() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_RITE_HALF         ) { v.SetTilePos( TP_FULL ); break
    } else if( TP == TP_TOP__RITE_QTR     ) { v.SetTilePos( TP_TOP__HALF )
    } else if( TP == TP_BOT__RITE_QTR     ) { v.SetTilePos( TP_BOT__HALF )
    } else if( TP == TP_RITE_QTR          ) { v.SetTilePos( TP_RITE_HALF )
    } else if( TP == TP_RITE_CTR__QTR     ) { v.SetTilePos( TP_LEFT_HALF )
    } else if( TP == TP_TOP__RITE_8TH     ) { v.SetTilePos( TP_TOP__RITE_QTR )
    } else if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR )
    } else if( TP == TP_BOT__RITE_8TH     ) { v.SetTilePos( TP_BOT__RITE_QTR )
    } else if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_RITE_HALF() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_LEFT_HALF         ) { v.SetTilePos( TP_FULL ); break
    } else if( TP == TP_TOP__LEFT_QTR     ) { v.SetTilePos( TP_TOP__HALF )
    } else if( TP == TP_BOT__LEFT_QTR     ) { v.SetTilePos( TP_BOT__HALF )
    } else if( TP == TP_LEFT_QTR          ) { v.SetTilePos( TP_LEFT_HALF )
    } else if( TP == TP_LEFT_CTR__QTR     ) { v.SetTilePos( TP_RITE_HALF )
    } else if( TP == TP_TOP__LEFT_8TH     ) { v.SetTilePos( TP_TOP__LEFT_QTR )
    } else if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR )
    } else if( TP == TP_BOT__LEFT_8TH     ) { v.SetTilePos( TP_BOT__LEFT_QTR )
    } else if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__HALF() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_BOT__HALF         ) { v.SetTilePos( TP_FULL ); break
    } else if( TP == TP_BOT__LEFT_QTR     ) { v.SetTilePos( TP_LEFT_HALF )
    } else if( TP == TP_BOT__RITE_QTR     ) { v.SetTilePos( TP_RITE_HALF )
    } else if( TP == TP_BOT__LEFT_8TH     ) { v.SetTilePos( TP_LEFT_QTR )
    } else if( TP == TP_BOT__RITE_8TH     ) { v.SetTilePos( TP_RITE_QTR )
    } else if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR )
    } else if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__HALF() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_TOP__HALF         ) { v.SetTilePos( TP_FULL ); break
    } else if( TP == TP_TOP__LEFT_QTR     ) { v.SetTilePos( TP_LEFT_HALF )
    } else if( TP == TP_TOP__RITE_QTR     ) { v.SetTilePos( TP_RITE_HALF )
    } else if( TP == TP_TOP__LEFT_8TH     ) { v.SetTilePos( TP_LEFT_QTR )
    } else if( TP == TP_TOP__RITE_8TH     ) { v.SetTilePos( TP_RITE_QTR )
    } else if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR )
    } else if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__LEFT_QTR() {

  if( m.Have_BOT__HALF() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_TOP__RITE_QTR     ) { v.SetTilePos( TP_TOP__HALF ); break
      } else if( TP == TP_TOP__RITE_8TH     ) { v.SetTilePos( TP_TOP__RITE_QTR )
      } else if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR )
      }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_BOT__LEFT_QTR     ) { v.SetTilePos( TP_LEFT_HALF ); break
      } else if( TP == TP_BOT__LEFT_8TH     ) { v.SetTilePos( TP_LEFT_QTR )
      } else if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR )
      }
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__RITE_QTR() {

  if( m.Have_BOT__HALF() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_TOP__LEFT_QTR     ) { v.SetTilePos( TP_TOP__HALF ); break
      } else if( TP == TP_TOP__LEFT_8TH     ) { v.SetTilePos( TP_TOP__LEFT_QTR )
      } else if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR )
      }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_BOT__RITE_QTR     ) { v.SetTilePos( TP_RITE_HALF ); break
      } else if( TP == TP_BOT__RITE_8TH     ) { v.SetTilePos( TP_RITE_QTR )
      } else if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR )
      }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__LEFT_QTR() {

  if( m.Have_TOP__HALF() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_BOT__RITE_QTR     ) { v.SetTilePos( TP_BOT__HALF ); break
      } else if( TP == TP_BOT__RITE_8TH     ) { v.SetTilePos( TP_BOT__RITE_QTR )
      } else if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR )
      }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_TOP__LEFT_QTR     ) { v.SetTilePos( TP_LEFT_HALF ); break
      } else if( TP == TP_TOP__LEFT_8TH     ) { v.SetTilePos( TP_LEFT_QTR )
      } else if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR )
      }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__RITE_QTR() {

  if( m.Have_TOP__HALF() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_BOT__LEFT_QTR     ) { v.SetTilePos( TP_BOT__HALF ); break
      } else if( TP == TP_BOT__LEFT_8TH     ) { v.SetTilePos( TP_BOT__LEFT_QTR )
      } else if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR )
      }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if       ( TP == TP_TOP__RITE_QTR     ) { v.SetTilePos( TP_RITE_HALF ); break
      } else if( TP == TP_TOP__RITE_8TH     ) { v.SetTilePos( TP_RITE_QTR )
      } else if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR )
      }
    }
  }
}

func (m *Vis) Quit_JoinTiles_LEFT_QTR() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_LEFT_CTR__QTR     ) { v.SetTilePos( TP_LEFT_HALF ); break
    } else if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR )
    } else if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_RITE_QTR() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_RITE_CTR__QTR     ) { v.SetTilePos( TP_RITE_HALF ); break
    } else if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR )
    } else if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_LEFT_CTR__QTR() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_LEFT_QTR      ) { v.SetTilePos( TP_LEFT_HALF ); break
    } else if( TP == TP_TOP__LEFT_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR )
    } else if( TP == TP_BOT__LEFT_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_RITE_CTR__QTR() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_RITE_QTR      ) { v.SetTilePos( TP_RITE_HALF ); break
    } else if( TP == TP_TOP__RITE_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR )
    } else if( TP == TP_BOT__RITE_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR )
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__LEFT_8TH() {

  if( m.Have_BOT__HALF() || m.Have_BOT__LEFT_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__LEFT_8TH ) { v.SetTilePos( TP_LEFT_QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__RITE_8TH() {

  if( m.Have_BOT__HALF() || m.Have_BOT__RITE_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__RITE_8TH ) { v.SetTilePos( TP_RITE_QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__LEFT_CTR_8TH() {

  if( m.Have_BOT__HALF() || m.Have_BOT__LEFT_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__LEFT_8TH ) { v.SetTilePos( TP_TOP__LEFT_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_TOP__RITE_CTR_8TH() {

  if( m.Have_BOT__HALF() || m.Have_BOT__RITE_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__RITE_8TH ) { v.SetTilePos( TP_TOP__RITE_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__LEFT_8TH() {

  if( m.Have_TOP__HALF() || m.Have_TOP__LEFT_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__LEFT_CTR_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__LEFT_8TH ) { v.SetTilePos( TP_LEFT_QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__RITE_8TH() {

  if( m.Have_TOP__HALF() || m.Have_TOP__RITE_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__RITE_CTR_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__RITE_8TH ) { v.SetTilePos( TP_RITE_QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__LEFT_CTR_8TH() {

  if( m.Have_TOP__HALF() || m.Have_TOP__LEFT_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__LEFT_8TH ) { v.SetTilePos( TP_BOT__LEFT_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__LEFT_CTR_8TH ) { v.SetTilePos( TP_LEFT_CTR__QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_BOT__RITE_CTR_8TH() {

  if( m.Have_TOP__HALF() || m.Have_TOP__RITE_QTR() ) {

    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_BOT__RITE_8TH ) { v.SetTilePos( TP_BOT__RITE_QTR ); break; }
    }
  } else {
    for k:=0; k<m.num_wins; k++ {

      var v *FileView = m.GetView_Win( k )
      var TP Tile_Pos = v.GetTilePos()

      if( TP == TP_TOP__RITE_CTR_8TH ) { v.SetTilePos( TP_RITE_CTR__QTR ); break; }
    }
  }
}

func (m *Vis) Quit_JoinTiles_LEFT_THIRD() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_CTR__THIRD      ) { v.SetTilePos( TP_LEFT_TWO_THIRDS )
    } else if( TP == TP_RITE_TWO_THIRDS ) { v.SetTilePos( TP_FULL )
    }
  }
}

func (m *Vis) Quit_JoinTiles_CTR__THIRD() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_RITE_THIRD ) { v.SetTilePos( TP_RITE_TWO_THIRDS ); }
  }
}

func (m *Vis) Quit_JoinTiles_RITE_THIRD() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if       ( TP == TP_CTR__THIRD      ) { v.SetTilePos( TP_RITE_TWO_THIRDS )
    } else if( TP == TP_LEFT_TWO_THIRDS ) { v.SetTilePos( TP_FULL )
    }
  }
}

func (m *Vis) Quit_JoinTiles_LEFT_TWO_THIRDS() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_RITE_THIRD ) { v.SetTilePos( TP_FULL ); }
  }
}

func (m *Vis) Quit_JoinTiles_RITE_TWO_THIRDS() {

  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_LEFT_THIRD ) { v.SetTilePos( TP_FULL ); }
  }
}

func (m *Vis) Have_BOT__HALF() bool {
  // Diff occupies bottom half:
// FIXME:
//if( m.vis.InDiffMode() ) {
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_BOT__LEFT_QTR && TP_l == TP_BOT__RITE_QTR )
//   || ( TP_s == TP_BOT__RITE_QTR && TP_l == TP_BOT__LEFT_QTR ) )
//  {
//    return true
//  }
//}
  // A view occupies bottom half:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_BOT__HALF ) { return true; }
  }
  return false
}

func (m *Vis) Have_TOP__HALF() bool {
  // Diff occupies top half:
// FIXME:
//if( m.vis.InDiffMode() )
//{
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_TOP__LEFT_QTR && TP_l == TP_TOP__RITE_QTR )
//   || ( TP_s == TP_TOP__RITE_QTR && TP_l == TP_TOP__LEFT_QTR ) )
//  {
//    return true
//  }
//}
  // A view occupies top half:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_TOP__HALF ) { return true; }
  }
  return false
}

func (m *Vis) Have_BOT__LEFT_QTR() bool {
  // Diff occupies bottom left quarter:
// FIXME:
//if( m.vis.InDiffMode() )
//{
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_BOT__LEFT_8TH     && TP_l == TP_BOT__LEFT_CTR_8TH )
//   || ( TP_s == TP_BOT__LEFT_CTR_8TH && TP_l == TP_BOT__LEFT_8TH ) )
//  {
//    return true
//  }
//}
  // A view occupies bottom left quarter:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_BOT__LEFT_QTR ) { return true; }
  }
  return false
}

func (m *Vis) Have_TOP__LEFT_QTR() bool {
  // Diff occupies top left quarter:
// FIXME:
//if( m.vis.InDiffMode() )
//{
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_TOP__LEFT_8TH     && TP_l == TP_TOP__LEFT_CTR_8TH )
//   || ( TP_s == TP_TOP__LEFT_CTR_8TH && TP_l == TP_TOP__LEFT_8TH ) )
//  {
//    return true
//  }
//}
  // A view occupies top left quarter:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_TOP__LEFT_QTR ) { return true; }
  }
  return false
}

func (m *Vis) Have_BOT__RITE_QTR() bool {
  // Diff occupies bottom right quarter:
// FIXME:
//if( m.vis.InDiffMode() )
//{
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_BOT__RITE_8TH     && TP_l == TP_BOT__RITE_CTR_8TH )
//   || ( TP_s == TP_BOT__RITE_CTR_8TH && TP_l == TP_BOT__RITE_8TH ) )
//  {
//    return true
//  }
//}
  // A view occupies bottom right quarter:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_BOT__RITE_QTR ) { return true; }
  }
  return false
}

func (m *Vis) Have_TOP__RITE_QTR() bool {
  // Diff occupies top right quarter:
// FIXME:
//if( m.vis.InDiffMode() )
//{
//  const Tile_Pos TP_s = m.diff.GetViewShort()->GetTilePos()
//  const Tile_Pos TP_l = m.diff.GetViewLong ()->GetTilePos()
//
//  if( ( TP_s == TP_TOP__RITE_8TH     && TP_l == TP_TOP__RITE_CTR_8TH )
//   || ( TP_s == TP_TOP__RITE_CTR_8TH && TP_l == TP_TOP__RITE_8TH ) )
//  {
//    return true
//  }
//}
  // A view occupies top right quarter:
  for k:=0; k<m.num_wins; k++ {

    var v *FileView = m.GetView_Win( k )
    var TP Tile_Pos = v.GetTilePos()

    if( TP == TP_TOP__RITE_QTR ) { return true; }
  }
  return false
}

//func (m *Vis) GoToNextWindow_l_Find() bool {
//
//  var found bool = false; // Found next view to go to
//
//  var p_curr_V *FileView = m.GetView_Win( m.win )
//  var curr_TP Tile_Pos = p_curr_V.GetTilePos()
//
//  if( curr_TP == TP_LEFT_HALF ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_TOP__RITE_QTR     == TP ||
//          TP_BOT__RITE_QTR     == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_HALF ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_QTR == TP ||
//          TP_BOT__LEFT_QTR == TP ||
//          TP_TOP__LEFT_8TH == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_TOP__RITE_QTR     == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_BOT__RITE_QTR     == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_QTR == TP ||
//          TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_BOT__LEFT_QTR == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_8TH == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_LEFT_CTR__QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__RITE_QTR     == TP ||
//          TP_BOT__RITE_QTR     == TP ||
//          TP_TOP__RITE_CTR_8TH == TP ||
//          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_CTR__QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_CTR__QTR     == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__RITE_QTR     == TP ||
//          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF         == TP ||
//          TP_RITE_CTR__QTR     == TP ||
//          TP_BOT__RITE_QTR     == TP ||
//          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++  {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_QTR      == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_QTR == TP ||
//          TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF     == TP ||
//          TP_LEFT_QTR      == TP ||
//          TP_BOT__LEFT_QTR == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  }
//  return found
//}

//func (m *Vis) GoToNextWindow_l_Find() bool {
//
//  var found bool = false; // Found next view to go to
//
//  var curr_TP Tile_Pos = m.GetView_Win( m.win ).GetTilePos()
//
//  if( curr_TP == TP_LEFT_HALF ) {
//    found = m.GoToNextWindow_l_Find_TP_LEFT_HALF()
//  } else if( curr_TP == TP_RITE_HALF ) {
//    found = m.GoToNextWindow_l_Find_TP_RITE_HALF()
//  } else if( curr_TP == TP_TOP__LEFT_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_QTR()
//  } else if( curr_TP == TP_BOT__LEFT_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_QTR()
//  } else if( curr_TP == TP_TOP__RITE_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__RITE_QTR()
//  } else if( curr_TP == TP_BOT__RITE_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__RITE_QTR()
//  } else if( curr_TP == TP_LEFT_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_LEFT_QTR()
//  } else if( curr_TP == TP_RITE_QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_RITE_QTR()
//  } else if( curr_TP == TP_LEFT_CTR__QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_LEFT_CTR__QTR()
//  } else if( curr_TP == TP_RITE_CTR__QTR ) {
//    found = m.GoToNextWindow_l_Find_TP_RITE_CTR__QTR()
//  } else if( curr_TP == TP_TOP__LEFT_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_8TH()
//  } else if( curr_TP == TP_BOT__LEFT_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_8TH()
//  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_CTR_8TH()
//  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_CTR_8TH()
//  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__RITE_CTR_8TH()
//  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__RITE_CTR_8TH()
//  } else if( curr_TP == TP_TOP__RITE_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_TOP__RITE_8TH()
//  } else if( curr_TP == TP_BOT__RITE_8TH ) {
//    found = m.GoToNextWindow_l_Find_TP_BOT__RITE_8TH()
//  }
//  return found
//}

//func (m *Vis) GoToNextWindow_l_Find() bool {
//
//  var found bool = false; // Found next view to go to
//
//  var curr_TP Tile_Pos = m.GetView_Win( m.win ).GetTilePos()
//
//  if       ( curr_TP == TP_LEFT_HALF )         { found = m.GoToNextWindow_l_Find_TP_LEFT_HALF()
//  } else if( curr_TP == TP_RITE_HALF )         { found = m.GoToNextWindow_l_Find_TP_RITE_HALF()
//  } else if( curr_TP == TP_TOP__LEFT_QTR )     { found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_QTR()
//  } else if( curr_TP == TP_BOT__LEFT_QTR )     { found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_QTR()
//  } else if( curr_TP == TP_TOP__RITE_QTR )     { found = m.GoToNextWindow_l_Find_TP_TOP__RITE_QTR()
//  } else if( curr_TP == TP_BOT__RITE_QTR )     { found = m.GoToNextWindow_l_Find_TP_BOT__RITE_QTR()
//  } else if( curr_TP == TP_LEFT_QTR )          { found = m.GoToNextWindow_l_Find_TP_LEFT_QTR()
//  } else if( curr_TP == TP_RITE_QTR )          { found = m.GoToNextWindow_l_Find_TP_RITE_QTR()
//  } else if( curr_TP == TP_LEFT_CTR__QTR )     { found = m.GoToNextWindow_l_Find_TP_LEFT_CTR__QTR()
//  } else if( curr_TP == TP_RITE_CTR__QTR )     { found = m.GoToNextWindow_l_Find_TP_RITE_CTR__QTR()
//  } else if( curr_TP == TP_TOP__LEFT_8TH )     { found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_8TH()
//  } else if( curr_TP == TP_BOT__LEFT_8TH )     { found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_8TH()
//  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) { found = m.GoToNextWindow_l_Find_TP_TOP__LEFT_CTR_8TH()
//  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) { found = m.GoToNextWindow_l_Find_TP_BOT__LEFT_CTR_8TH()
//  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) { found = m.GoToNextWindow_l_Find_TP_TOP__RITE_CTR_8TH()
//  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) { found = m.GoToNextWindow_l_Find_TP_BOT__RITE_CTR_8TH()
//  } else if( curr_TP == TP_TOP__RITE_8TH )     { found = m.GoToNextWindow_l_Find_TP_TOP__RITE_8TH()
//  } else if( curr_TP == TP_BOT__RITE_8TH )     { found = m.GoToNextWindow_l_Find_TP_BOT__RITE_8TH()
//  }
//  return found
//}

func (m *Vis) GoToNextWindow_l_Find() bool {

  var found bool = false; // Found next view to go to

  var curr_TP Tile_Pos = m.GetView_Win( m.win ).GetTilePos()

  if       ( curr_TP == TP_LEFT_HALF )         { found = m.GoToNextWin_l_TP_LEFT_HALF()
  } else if( curr_TP == TP_RITE_HALF )         { found = m.GoToNextWin_l_TP_RITE_HALF()
  } else if( curr_TP == TP_TOP__LEFT_QTR )     { found = m.GoToNextWin_l_TP_TOP__LEFT_QTR()
  } else if( curr_TP == TP_BOT__LEFT_QTR )     { found = m.GoToNextWin_l_TP_BOT__LEFT_QTR()
  } else if( curr_TP == TP_TOP__RITE_QTR )     { found = m.GoToNextWin_l_TP_TOP__RITE_QTR()
  } else if( curr_TP == TP_BOT__RITE_QTR )     { found = m.GoToNextWin_l_TP_BOT__RITE_QTR()
  } else if( curr_TP == TP_LEFT_QTR )          { found = m.GoToNextWin_l_TP_LEFT_QTR()
  } else if( curr_TP == TP_RITE_QTR )          { found = m.GoToNextWin_l_TP_RITE_QTR()
  } else if( curr_TP == TP_LEFT_CTR__QTR )     { found = m.GoToNextWin_l_TP_LEFT_CTR__QTR()
  } else if( curr_TP == TP_RITE_CTR__QTR )     { found = m.GoToNextWin_l_TP_RITE_CTR__QTR()
  } else if( curr_TP == TP_TOP__LEFT_8TH )     { found = m.GoToNextWin_l_TP_TOP__LEFT_8TH()
  } else if( curr_TP == TP_BOT__LEFT_8TH )     { found = m.GoToNextWin_l_TP_BOT__LEFT_8TH()
  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) { found = m.GoToNextWin_l_TP_TOP__LEFT_CTR_8TH()
  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) { found = m.GoToNextWin_l_TP_BOT__LEFT_CTR_8TH()
  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) { found = m.GoToNextWin_l_TP_TOP__RITE_CTR_8TH()
  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) { found = m.GoToNextWin_l_TP_BOT__RITE_CTR_8TH()
  } else if( curr_TP == TP_TOP__RITE_8TH )     { found = m.GoToNextWin_l_TP_TOP__RITE_8TH()
  } else if( curr_TP == TP_BOT__RITE_8TH )     { found = m.GoToNextWin_l_TP_BOT__RITE_8TH()
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_LEFT_HALF() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_TOP__RITE_QTR     == TP ||
        TP_BOT__RITE_QTR     == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_RITE_HALF() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_QTR == TP ||
        TP_BOT__LEFT_QTR == TP ||
        TP_TOP__LEFT_8TH == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_BOT__RITE_QTR     == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_QTR == TP ||
        TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_BOT__LEFT_QTR == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_QTR == TP ||
        TP_BOT__LEFT_QTR == TP ||
        TP_TOP__LEFT_8TH == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_LEFT_CTR__QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_QTR     == TP ||
        TP_BOT__RITE_QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_RITE_CTR__QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_QTR      == TP ||
        TP_TOP__RITE_8TH == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__LEFT_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__LEFT_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_CTR__QTR     == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__LEFT_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__LEFT_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF         == TP ||
        TP_RITE_CTR__QTR     == TP ||
        TP_BOT__RITE_QTR     == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__RITE_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_QTR      == TP ||
        TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__RITE_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++  {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_QTR      == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_TOP__RITE_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_QTR == TP ||
        TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_l_TP_BOT__RITE_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF     == TP ||
        TP_LEFT_QTR      == TP ||
        TP_BOT__LEFT_QTR == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

//func (m *Vis) GoToNextWindow_h_Find() bool {
//
//  var found bool = false; // Found next view to go to
//
//  var p_curr_V *FileView = m.GetView_Win( m.win )
//  var curr_TP Tile_Pos = p_curr_V.GetTilePos()
//
//  if( curr_TP == TP_LEFT_HALF ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_TOP__RITE_QTR == TP ||
//          TP_BOT__RITE_QTR == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_HALF ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_TOP__LEFT_QTR     == TP ||
//          TP_BOT__LEFT_QTR     == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_TOP__RITE_QTR == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_BOT__RITE_QTR == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_BOT__LEFT_QTR     == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_LEFT_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_LEFT_CTR__QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_8TH == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_RITE_CTR__QTR ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_QTR     == TP ||
//          TP_BOT__LEFT_QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_TOP__RITE_QTR == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_HALF     == TP ||
//          TP_BOT__RITE_QTR == TP ||
//          TP_RITE_QTR      == TP ||
//          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_QTR      == TP ||
//          TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_QTR      == TP ||
//          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_TOP__LEFT_QTR     == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_LEFT_HALF         == TP ||
//          TP_BOT__LEFT_QTR     == TP ||
//          TP_LEFT_CTR__QTR     == TP ||
//          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_TOP__RITE_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_CTR__QTR     == TP ||
//          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  } else if( curr_TP == TP_BOT__RITE_8TH ) {
//    for k:=0; !found && k<m.num_wins; k++ {
//      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
//      if( TP_RITE_CTR__QTR     == TP ||
//          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
//    }
//  }
//  return found
//}

func (m *Vis) GoToNextWindow_h_Find() bool {

  var found bool = false; // Found next view to go to

  var p_curr_V *FileView = m.GetView_Win( m.win )
  var curr_TP Tile_Pos = p_curr_V.GetTilePos()

  if       ( curr_TP == TP_LEFT_HALF )         { found = m.GoToNextWin_h_TP_LEFT_HALF()
  } else if( curr_TP == TP_RITE_HALF )         { found = m.GoToNextWin_h_TP_RITE_HALF()
  } else if( curr_TP == TP_TOP__LEFT_QTR )     { found = m.GoToNextWin_h_TP_TOP__LEFT_QTR()
  } else if( curr_TP == TP_BOT__LEFT_QTR )     { found = m.GoToNextWin_h_TP_BOT__LEFT_QTR()
  } else if( curr_TP == TP_TOP__RITE_QTR )     { found = m.GoToNextWin_h_TP_TOP__RITE_QTR()
  } else if( curr_TP == TP_BOT__RITE_QTR )     { found = m.GoToNextWin_h_TP_BOT__RITE_QTR()
  } else if( curr_TP == TP_LEFT_QTR )          { found = m.GoToNextWin_h_TP_LEFT_QTR()
  } else if( curr_TP == TP_RITE_QTR )          { found = m.GoToNextWin_h_TP_RITE_QTR()
  } else if( curr_TP == TP_LEFT_CTR__QTR )     { found = m.GoToNextWin_h_TP_LEFT_CTR__QTR()
  } else if( curr_TP == TP_RITE_CTR__QTR )     { found = m.GoToNextWin_h_TP_RITE_CTR__QTR()
  } else if( curr_TP == TP_TOP__LEFT_8TH )     { found = m.GoToNextWin_h_TP_TOP__LEFT_8TH()
  } else if( curr_TP == TP_BOT__LEFT_8TH )     { found = m.GoToNextWin_h_TP_BOT__LEFT_8TH()
  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) { found = m.GoToNextWin_h_TP_TOP__LEFT_CTR_8TH()
  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) { found = m.GoToNextWin_h_TP_BOT__LEFT_CTR_8TH()
  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) { found = m.GoToNextWin_h_TP_TOP__RITE_CTR_8TH()
  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) { found = m.GoToNextWin_h_TP_BOT__RITE_CTR_8TH()
  } else if( curr_TP == TP_TOP__RITE_8TH )     { found = m.GoToNextWin_h_TP_TOP__RITE_8TH()
  } else if( curr_TP == TP_BOT__RITE_8TH )     { found = m.GoToNextWin_h_TP_BOT__RITE_8TH()
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_LEFT_HALF() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_TOP__RITE_QTR == TP || // FIXME
        TP_BOT__RITE_QTR == TP ||
        TP_TOP__RITE_8TH == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_RITE_HALF() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_TOP__LEFT_QTR     == TP ||
        TP_BOT__LEFT_QTR     == TP || // FIXME
        TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_LEFT_CTR__QTR     == TP ||
        TP_BOT__LEFT_QTR     == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_TOP__RITE_QTR == TP ||
        TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_BOT__RITE_QTR == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_RITE_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_LEFT_QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_TOP__RITE_QTR == TP ||
        TP_BOT__RITE_QTR == TP ||
        TP_TOP__RITE_8TH == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_RITE_CTR__QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_QTR     == TP ||
        TP_BOT__LEFT_QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_LEFT_CTR__QTR() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_8TH == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__RITE_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_CTR__QTR     == TP ||
        TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__RITE_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_CTR__QTR     == TP ||
        TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__RITE_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_LEFT_CTR__QTR     == TP ||
        TP_TOP__LEFT_QTR     == TP ||
        TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__RITE_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_HALF         == TP ||
        TP_LEFT_CTR__QTR     == TP ||
        TP_BOT__LEFT_QTR     == TP ||
        TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__LEFT_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_QTR      == TP ||
        TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__LEFT_CTR_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_LEFT_QTR      == TP ||
        TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_TOP__LEFT_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_TOP__RITE_QTR == TP ||
        TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

func (m *Vis) GoToNextWin_h_TP_BOT__LEFT_8TH() bool {

  var found bool = false; // Found next view to go to

  for k:=0; !found && k<m.num_wins; k++ {
    var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
    if( TP_RITE_HALF     == TP ||
        TP_RITE_QTR      == TP ||
        TP_BOT__RITE_QTR == TP ||
        TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
  }
  return found
}

//func (m *Vis) GoToNextWin_h_() bool {
//
//  var found bool = false; // Found next view to go to
//
//  return found
//}

func (m *Vis)  GoToNextWindow_jk_Find() bool {

  var found bool = false; // Found next view to go to

  var p_curr_V *FileView = m.GetView_Win( m.win )
  var curr_TP Tile_Pos = p_curr_V.GetTilePos()

  if( curr_TP == TP_TOP__HALF ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF         == TP ||
          TP_BOT__LEFT_QTR     == TP ||
          TP_BOT__RITE_QTR     == TP ||
          TP_BOT__LEFT_8TH     == TP ||
          TP_BOT__RITE_8TH     == TP ||
          TP_BOT__LEFT_CTR_8TH == TP ||
          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__HALF ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF         == TP ||
          TP_TOP__LEFT_QTR     == TP ||
          TP_TOP__RITE_QTR     == TP ||
          TP_TOP__LEFT_8TH     == TP ||
          TP_TOP__RITE_8TH     == TP ||
          TP_TOP__LEFT_CTR_8TH == TP ||
          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__LEFT_QTR ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF         == TP ||
          TP_BOT__LEFT_QTR     == TP ||
          TP_BOT__LEFT_8TH     == TP ||
          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__RITE_QTR ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF         == TP ||
          TP_BOT__RITE_QTR     == TP ||
          TP_BOT__RITE_8TH     == TP ||
          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__LEFT_QTR ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF         == TP ||
          TP_TOP__LEFT_QTR     == TP ||
          TP_TOP__LEFT_8TH     == TP ||
          TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__RITE_QTR ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF         == TP ||
          TP_TOP__RITE_QTR     == TP ||
          TP_TOP__RITE_8TH     == TP ||
          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__LEFT_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF     == TP ||
          TP_BOT__LEFT_QTR == TP ||
          TP_BOT__LEFT_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__RITE_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF     == TP ||
          TP_BOT__RITE_QTR == TP ||
          TP_BOT__RITE_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__LEFT_CTR_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF         == TP ||
          TP_BOT__LEFT_QTR     == TP ||
          TP_BOT__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_TOP__RITE_CTR_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_BOT__HALF         == TP ||
          TP_BOT__RITE_QTR     == TP ||
          TP_BOT__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__LEFT_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF     == TP ||
          TP_TOP__LEFT_QTR == TP ||
          TP_TOP__LEFT_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__RITE_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF     == TP ||
          TP_TOP__RITE_QTR == TP ||
          TP_TOP__RITE_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__LEFT_CTR_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF         == TP ||
          TP_TOP__LEFT_QTR     == TP ||
          TP_TOP__LEFT_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  } else if( curr_TP == TP_BOT__RITE_CTR_8TH ) {
    for k:=0; !found && k<m.num_wins; k++ {
      var TP Tile_Pos = m.GetView_Win( k ).GetTilePos()
      if( TP_TOP__HALF         == TP ||
          TP_TOP__RITE_QTR     == TP ||
          TP_TOP__RITE_CTR_8TH == TP ) { m.win = k; found = true; }
    }
  }
  return found
}

func FlipWindows_Horizontally( OTP Tile_Pos ) Tile_Pos {

  var NTP Tile_Pos = TP_NONE

  if       ( OTP == TP_LEFT_HALF         ) { NTP = TP_RITE_HALF
  } else if( OTP == TP_RITE_HALF         ) { NTP = TP_LEFT_HALF
  } else if( OTP == TP_TOP__LEFT_QTR     ) { NTP = TP_TOP__RITE_QTR
  } else if( OTP == TP_TOP__RITE_QTR     ) { NTP = TP_TOP__LEFT_QTR
  } else if( OTP == TP_BOT__LEFT_QTR     ) { NTP = TP_BOT__RITE_QTR
  } else if( OTP == TP_BOT__RITE_QTR     ) { NTP = TP_BOT__LEFT_QTR
  } else if( OTP == TP_LEFT_QTR          ) { NTP = TP_RITE_QTR
  } else if( OTP == TP_RITE_QTR          ) { NTP = TP_LEFT_QTR
  } else if( OTP == TP_LEFT_CTR__QTR     ) { NTP = TP_RITE_CTR__QTR
  } else if( OTP == TP_RITE_CTR__QTR     ) { NTP = TP_LEFT_CTR__QTR
  } else if( OTP == TP_TOP__LEFT_8TH     ) { NTP = TP_TOP__RITE_8TH
  } else if( OTP == TP_TOP__RITE_8TH     ) { NTP = TP_TOP__LEFT_8TH
  } else if( OTP == TP_TOP__LEFT_CTR_8TH ) { NTP = TP_TOP__RITE_CTR_8TH
  } else if( OTP == TP_TOP__RITE_CTR_8TH ) { NTP = TP_TOP__LEFT_CTR_8TH
  } else if( OTP == TP_BOT__LEFT_8TH     ) { NTP = TP_BOT__RITE_8TH
  } else if( OTP == TP_BOT__RITE_8TH     ) { NTP = TP_BOT__LEFT_8TH
  } else if( OTP == TP_BOT__LEFT_CTR_8TH ) { NTP = TP_BOT__RITE_CTR_8TH
  } else if( OTP == TP_BOT__RITE_CTR_8TH ) { NTP = TP_BOT__LEFT_CTR_8TH
  }
  return NTP
}

func FlipWindows_Vertically( OTP Tile_Pos ) Tile_Pos {

  var NTP Tile_Pos = TP_NONE

  if       ( OTP == TP_TOP__HALF         ) { NTP = TP_BOT__HALF
  } else if( OTP == TP_BOT__HALF         ) { NTP = TP_TOP__HALF
  } else if( OTP == TP_TOP__LEFT_QTR     ) { NTP = TP_BOT__LEFT_QTR
  } else if( OTP == TP_TOP__RITE_QTR     ) { NTP = TP_BOT__RITE_QTR
  } else if( OTP == TP_BOT__LEFT_QTR     ) { NTP = TP_TOP__LEFT_QTR
  } else if( OTP == TP_BOT__RITE_QTR     ) { NTP = TP_TOP__RITE_QTR
  } else if( OTP == TP_TOP__LEFT_8TH     ) { NTP = TP_BOT__LEFT_8TH
  } else if( OTP == TP_TOP__RITE_8TH     ) { NTP = TP_BOT__RITE_8TH
  } else if( OTP == TP_TOP__LEFT_CTR_8TH ) { NTP = TP_BOT__LEFT_CTR_8TH
  } else if( OTP == TP_TOP__RITE_CTR_8TH ) { NTP = TP_BOT__RITE_CTR_8TH
  } else if( OTP == TP_BOT__LEFT_8TH     ) { NTP = TP_TOP__LEFT_8TH
  } else if( OTP == TP_BOT__RITE_8TH     ) { NTP = TP_TOP__RITE_8TH
  } else if( OTP == TP_BOT__LEFT_CTR_8TH ) { NTP = TP_TOP__LEFT_CTR_8TH
  } else if( OTP == TP_BOT__RITE_CTR_8TH ) { NTP = TP_TOP__RITE_CTR_8TH
  }
  return NTP
}

