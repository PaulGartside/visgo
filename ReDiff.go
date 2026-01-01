
package main

// Returns success or failure
func (m *Diff) Partial_ReDiff() bool {
  ok := false
  DL := m.CrsLine() // Diff line number
  NUM_DLs := m.DI_List_L.Len()
  SIDE_BAND := 50
  DL_st := Max_i( DL - SIDE_BAND, 0 )
  DL_fn := Min_i( DL + SIDE_BAND, NUM_DLs )

  var da DiffArea
  found_diff_area := m.ReDiff_GetDiffArea( DL_st, DL_fn, &da )

  if( !found_diff_area ) {
    m_vis.CmdLineMessage("rediff: DiffArea not found")
    m.PrintCursor( m_vis.CV() )

  } else {
    ok = true

    m.DI_L_ins_idx = m.Remove_From_DI_Lists( da )

    m.RunDiff( da )
  }
  return ok
}

func (m *Diff) ReDiff_GetDiffArea( DL_st, DL_fn int, da *DiffArea ) bool {
  // Diff area is from DL_st up to but not including DL_fn
  found_st := m.ReDiff_FindDiffSt( &DL_st )
  found_fn := false

  if( found_st ) {
    found_fn = true
    if( DL_fn < m.DI_List_L.Len() ) {
      found_fn = m.ReDiff_FindDiffFn( &DL_fn )
    }
  }
  found_diff_area := found_st && found_fn

  if( found_diff_area ) {
    *da = m.DL_st_fn_2_DiffArea( DL_st, DL_fn )

  } else {
    DL := m.CrsLine() // Diff line number

    DL_st_2 := DL // local diff line start
    DL_fn_2 := DL // local diff line finish

    found_diff_area = m.ReDiff_FindDiffSt( &DL_st_2 ) &&
                      m.ReDiff_FindDiffFn( &DL_fn_2 )

    if( found_diff_area ) {
      *da = m.DL_st_fn_2_DiffArea( DL_st_2, DL_fn_2 )
    }
  }
  return found_diff_area
}

func (m *Diff) Remove_From_DI_Lists( da DiffArea ) int {

  DI_list_s_remove_st := m.DiffLine_S( da.ln_s )
  DI_list_l_remove_st := m.DiffLine_L( da.ln_l )

  DI_list_remove_st := Min_i( DI_list_s_remove_st, DI_list_l_remove_st )

  DI_lists_insert_idx := DI_list_remove_st

  DI_list_s_remove_fn := m.DiffLine_S( da.fnl_s() )
  if( m.pvS.p_fb.NumLines() <= da.fnl_s() ) {
    DI_list_s_remove_fn = m.DI_List_S.Len()
  }
  DI_list_l_remove_fn := m.DiffLine_L( da.fnl_l() )
  if( m.pvL.p_fb.NumLines() <= da.fnl_l() ) {
    DI_list_l_remove_fn = m.DI_List_L.Len()
  }

  DI_list_remove_fn := Max_i( DI_list_s_remove_fn, DI_list_l_remove_fn )

//Log.Log("(DI_list_remove_st,DI_list_remove_fn) = ("
//       + (DI_list_remove_st+1)+","+(DI_list_remove_fn+1) +")")

  for k:=DI_list_remove_st; k<DI_list_remove_fn; k++ {
    m.DI_List_S.Remove( DI_lists_insert_idx )
    m.DI_List_L.Remove( DI_lists_insert_idx )
  }
  return DI_lists_insert_idx
}

func (m *Diff) ReDiff_FindDiffSt( p_DL_st *int ) bool {
  found_diff_st := false

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  var cDI_st Diff_Info = cDI_List.Get( *p_DL_st )

  if( DT_SAME == cDI_st.diff_type ) {
    found_diff_st = m.ReDiff_GetSt_Search_4_Diff_Then_Same( p_DL_st )

  } else if( DT_CHANGED  == cDI_st.diff_type ||
             DT_INSERTED == cDI_st.diff_type ||
             DT_DELETED  == cDI_st.diff_type ) {
    found_diff_st = m.ReDiff_GetDiffSt_Search_4_Same( p_DL_st )
  }
  return found_diff_st
}

func (m *Diff) ReDiff_FindDiffFn( p_DL_fn *int ) bool {
  found_diff_fn := false

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  var cDI_fn Diff_Info = cDI_List.Get( *p_DL_fn )

  if( DT_SAME == cDI_fn.diff_type ) {
    found_diff_fn = m.ReDiff_GetFn_Search_4_Diff_Then_Same( p_DL_fn )

  } else if( DT_CHANGED  == cDI_fn.diff_type ||
             DT_INSERTED == cDI_fn.diff_type ||
             DT_DELETED  == cDI_fn.diff_type ) {
    found_diff_fn = m.ReDiff_GetDiffFn_Search_4_Same( p_DL_fn )
  }
  return found_diff_fn
}

func (m *Diff) DL_st_fn_2_DiffArea( DL_st, DL_fn int ) DiffArea {
  var da DiffArea

  da.ln_s = m.DI_List_S.Get( DL_st ).line_num
  da.ln_l = m.DI_List_L.Get( DL_st ).line_num

  if( DT_DELETED == m.DI_List_S.Get( DL_st ).diff_type ) {
    da.ln_s += 1
  }
  if( DT_DELETED == m.DI_List_L.Get( DL_st ).diff_type ) {
    da.ln_l += 1
  }

  if( DL_fn < m.DI_List_L.Len() ) {
    da.nlines_s = m.DI_List_S.Get( DL_fn ).line_num - da.ln_s
    da.nlines_l = m.DI_List_L.Get( DL_fn ).line_num - da.ln_l
  } else {
    // Need the extra -1 here to avoid a crash.
    // Not sure why it is needed.
  //da.nlines_s = m.pfS.NumLines() - da.ln_s - 1
  //da.nlines_l = m.pfL.NumLines() - da.ln_l - 1
    da.nlines_s = m.pfS.NumLines() - da.ln_s
    da.nlines_l = m.pfL.NumLines() - da.ln_l
  }
  return da
}

func (m *Diff) ReDiff_GetSt_Search_4_Diff_Then_Same( p_DL_st *int ) bool {

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  // Search up for CHANGED, INSERTED or DELETED and then for SAME
  found := false
  L := *p_DL_st
  for ; !found && 0<=L; L-- {
    var di Diff_Info = cDI_List.Get( L )
    if( DT_CHANGED  == di.diff_type ||
        DT_INSERTED == di.diff_type ||
        DT_DELETED  == di.diff_type ) {
      found = true
    }
  }
  if( found ) {
    found = false
    for ; !found && 0<=L; L-- {
      var di Diff_Info = cDI_List.Get( L )
      if( DT_SAME == di.diff_type ) {
        found = true
        *p_DL_st = L+1 // Diff area starts on first diff after first same
      }
    }
  }
  if( !found && L < 0 ) {
    found = true
    *p_DL_st = 0 // Diff area starts at beginning of file
  }
  return found
}

func (m *Diff) ReDiff_GetDiffSt_Search_4_Same( p_DL_st *int ) bool {

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  // Search up for SAME
  found := false
  L := *p_DL_st
  for ; !found && 0<=L; L-- {
    var di Diff_Info = cDI_List.Get( L )
    if( DT_SAME == di.diff_type ) {
      found = true
      *p_DL_st = L+1 // Diff area starts on first diff after first same
    }
  }
  if( !found && L < 0 ) {
    found = true
    *p_DL_st = 0 // Diff area starts at beginning of file
  }
  return found
}

func (m *Diff) ReDiff_GetFn_Search_4_Diff_Then_Same( p_DL_fn *int ) bool {

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  // Search down for CHANGED, INSERTED or DELETED and then for SAME
  found := false
  L := *p_DL_fn
  for ; !found && L<cDI_List.Len(); L++ {
    var di Diff_Info = cDI_List.Get( L )
    if( DT_CHANGED  == di.diff_type ||
        DT_INSERTED == di.diff_type ||
        DT_DELETED  == di.diff_type ) {
      found = true
    }
  }
  if( found ) {
    found = false
    for ; !found && L<cDI_List.Len(); L++ {
      var di Diff_Info = cDI_List.Get( L )
      if( DT_SAME == di.diff_type ) {
        found = true
        *p_DL_fn = L
      }
    }
  }
  if( !found && cDI_List.Len() < L ) {
    found = true
    *p_DL_fn = cDI_List.Len(); // Diff area ends at end of file
  }
  return found
}

func (m *Diff) ReDiff_GetDiffFn_Search_4_Same( p_DL_fn *int ) bool {

  var cDI_List *Vector[Diff_Info] = m.View_2_DI_List_C( m_vis.CV() )

  // Search down for SAME
  found := false
  L := *p_DL_fn
  for ; !found && L<cDI_List.Len(); L++ {
    var di Diff_Info = cDI_List.Get( L )
    if( DT_SAME == di.diff_type ) {
      found = true
      *p_DL_fn = L
    }
  }
  if( !found && cDI_List.Len() < L ) {
    found = true
    *p_DL_fn = cDI_List.Len() // Diff area ends at end of file
  }
  return found
}

