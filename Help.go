
package main

import (
//"bytes"
)

var HELP_STR string =
"-------------------------\n" +
" VI-Simplified (vis) Help\n" +
"-------------------------\n" +
"Quit:\n" +
"  :q   - quit current file, exit vis\n" +
"\n" +
"Basic movement:\n" +
"               \n" +
"      /|\\      \n" +
"       |       \n" +
"       k       \n" +
"  <--h   l-->  \n" +
"       j       \n" +
"       |       \n" +
"      \\|/      \n" +
"\n" +
"Movement:\n" +
"  j    - Move one space down\n" +
"  k    - Move one space up\n" +
"  h    - Move one space left\n" +
"  l    - Move one space right\n" +
"  0    - Go to start of line\n" +
"  $    - Go to end   of line\n" +
"  g0   - Go to left  column of window\n" +
"  g$   - Go to right column of window\n" +
"  gf   - Go to filename under cursor or visually highlighed area\n" +
"  gp   - Go to pattern visually highlighted (VISUAL mode only)\n" +
"  H    - Go to top    of window\n" +
"  L    - Go to bottom of window\n" +
"  M    - Go to middle line in window\n" +
"  gg   - Go to start of file\n" +
"  G    - Go to last line of file\n" +
"  b    - Go to previous word\n" +
"  w    - Go to next word\n" +
"  e    - Go to end of next word\n" +
"  B    - Page up\n" +
"  F    - Page down\n" +
"  f    - Search forward on line for next char typed\n" +
"  ;    - Repeat last 'f' search\n" +
"  z<CR>- Move current cursor line to top of window>\n" +
"  zt   - Move cursor line to top of window\n" +
"  zz   - Move current cursor line to middle of window>\n" +
"  zb   - Move cursor line to bottom of window\n" +
"  %    - Jump to opposite {, [, (, }, ] or ) bracket\n" +
"  {    - Jump to enclosing { bracket\n" +
"  }    - Jump to enclosing } bracket\n" +
"  WW   - Move to next editing window\n" +
"  WR   - Flip windows\n" +
"\n" +
"Editing:\n" +
"  x    - Delete char under cursor or visually highlighted area\n" +
"  cp   - Delete from cursor to end of highlighed search area and enter insert mode (NOT IMPLEMENTED)\n" +
"  cw   - Delete from cursor to end of word and enter insert mode\n" +
"  c$   - Delete from cursor to end of line and enter insert mode\n" +
"  d    - Delete visually highlighted area (VISUAL mode only)\n" +
"  dd   - Delete current line\n" +
"  dp   - Delete from cursor to end of highlighed search area (NOT IMPLEMENTED)\n" +
"  dw   - Delete from cursor to end of word\n" +
"  D    - Delete from cursor to end of line.\n" +
"  ~    - Toggle case\n" +
"  .    - Repeat last command\n" +
"  i    - Enter insert mode\n" +
"  a    - Move forward one char and enter insert mode\n" +
"  A    - Move past end of line and enter insert mode\n" +
"  o    - Begin a new line below current line\n" +
"  O    - Begin a new line above current line\n" +
" <ESC> - Leave insert mode\n" +
"  yy   - Yank current line\n" +
"  yw   - Yank from cursor to end of word into paste buffer\n" +
"  y    - When in VISUAL non-block mode yank highlighed area\n" +
"  y    - When in VISUAL     block mode yank highlighed block\n" +
"  Y    - When in VISUAL non-block mode yank highlighed lines\n" +
"  Y    - When in VISUAL     block mode yank highlighed block\n" +
"  r    - When in VISUAL non-block mode yank and erase highlighed area\n" +
"  r    - When in VISUAL     block mode yank and erase highlighed block\n" +
"  R    - When in VISUAL non-block mode yank and erase highlighed lines\n" +
"  R    - When in VISUAL     block mode yank and erase highlighed block\n" +
"  p    - Paste buffer to line below (in front of cursor)\n" +
"  P    - Paste buffer to line above (behind cursor)\n" +
"  J    - Join next line with current line\n" +
"  r    - Replace white space in front of cursor with paste buffer\n" +
"  R    - Replace following characters with those typed\n" +
"  u    - Undo previous change\n" +
"  U    - Undo all changes\n" +
"  m    - Execute map\n" +
"  Q    - Built in map for: '.j0'\n" +
"  s    - Delete char under cursor or visually highlighted area and enter insert mode\n" +
"  v    - Enter VISUAL non-block mode\n" +
"  V    - Enter VISUAL     block mode\n" +
"\n" +
"Search:\n" +
"  *    - Highlight and search for identifier under cursor\n" +
"  /    - Enter search pattern a prompt\n" +
"  n    - Go to next highlighted search item or next diff\n" +
"  N    - Go to previous highlighted search item or prev diff\n" +
"  gp   - Go to pattern visually highlighted (VISUAL mode only)\n" +
"\n" +
"Command Prompt:\n" +
"  :    - Go to command prompt, exit map mode\n" +
"  :b # - Switch to previous buffer\n" +
"  :b N - Switch to buffer N\n" +
"  :be  - Switch to the buffer editor\n" +
"  :bm  - Switch to the message buffer\n" +
"  :cd  - Change working directory to that of current file\n" +
"  :cs1 - Switch to color scheme 1\n" +
"  :cs2 - Switch to color scheme 1\n" +
"  :cs3 - Switch to color scheme 1\n" +
"  :detab=tab_size = Remove tabs. Tabs are tab_size\n" +
"  :diff- Enter diff mode\n" +
"  :nodiff- Exit diff mode\n" +
"  :hi  - Re-syntax-highlight file\n" +
"  :help- Go to help buffer\n" +
"  :e   - Re-read current file\n" +
"  :e filename - Edit filename\n" +
"  :map - Enter map mode to map a command\n" +
"  :n   - Go to next buffer\n" +
"  :pwd - Display current working directory\n" +
"  :q   - quit current file\n" +
"  :qa  - quit all files\n" +
"  :re  - Re-draw entire screen\n" +
"  :run - If in shell buffer, run shell command\n" +
"  :se  - Switch to the search editor\n" +
"  :sh  - Enter shell buffer\n" +
"  :shell-Enter shell buffer\n" +
"  :showmap - Show current map\n" +
"  :syn=syntax_type - Highlight current file with syntax_type\n" +
"  :sp  - Split window horizonally into 2 editing windows\n" +
"  :ts=       - Show tabs size\n" +
"  :ts=<size> - Set tabs to size for current file\n" +
"  :vsp - Split window vertically into 2 editing windows\n" +
"  :3sp - Split window vertically into 3 editing windows\n" +
"  :w   - Write current file\n" +
"  :wq  - Write and quit current file\n" +
"  :w filename - Write current file to filename\n" +
"\n" +
"Alphabetic list of commands:\n" +
"  a    - Move forward one char and enter insert mode\n" +
"  A    - Move past end of line and enter insert mode\n" +
"  b    - Go to previous word\n" +
"  B    - Page up\n" +
"  cw   - Change word (Leaves you in INSERT mode)\n" +
"  c$   - Delete from cursor to end of line and enter insert mode\n" +
"  cp   - Delete from cursor to end of highlighed search area and enter insert mode (NOT IMPLEMENTED)\n" +
"  dd   - Delete current line into paste buffer\n" +
"  dp   - Delete from cursor to end of highlighed search area (NOT IMPLEMENTED)\n" +
"  dw   - Delete word\n" +
"  D    - Delete from cursor to end of line\n" +
"  e    - Go to end of next word\n" +
"  f    - Search forward on line for next char typed\n" +
"  F    - Page down\n" +
"  gg   - Go to start of file\n" +
"  g0   - Go to left  column of window\n" +
"  g$   - Go to right column of window\n" +
"  gf   - Go to filename under cursor or visually highlighed\n" +
"  gp   - Go to pattern visually highlighted (VISUAL mode only)\n" +
"  G    - Go to last line of file\n" +
"  h    - Move one space left\n" +
"  H    - Go to top    of window\n" +
"  i    - Insert before cursor (Leaves you in INSERT mode)\n" +
"  j    - Move one space down\n" +
"  J    - Join next line with current line\n" +
"  k    - Move one space up\n" +
"  l    - Move one space right\n" +
"  L    - Go to bottom of window\n" +
"  m    - Execute map\n" +
"  M    - Go to middle line in window\n" +
"  n    - Go to next highlighted search item or next diff\n" +
"  N    - Go to previous highlighted search item or prev diff\n" +
"  o    - Begin a new line below current line (Leaves you in INSERT mode)\n" +
"  O    - Begin a new line above current line (Leaves you in INSERT mode)\n" +
"  p    - Paste paste buffer to line below or in front of cursor\n" +
"  P    - Paste paste buffer to line above or behind cursor\n" +
"  Q    - Built in map for: '.j0'\n" +
"  r    - Replace are in front of cursor with paste buffer (NOT IMPLEMENTED)\n" +
"  R    - Replace following characters with those typed (Leaves you in INSERT mode)\n" +
"  s    - Delete char under cursor or visually highlighted area and enter insert mode\n" +
"  u    - Undo previous change\n" +
"  U    - Undo all changes\n" +
"  v    - Enter VISUAL non-block mode\n" +
"  V    - Enter VISUAL     block mode\n" +
"  w    - Go to next word\n" +
"  WW   - Go to next window\n" +
"  WR   - Flip windows\n" +
"  x    - Delete char under cursor or visually highlighted area\n" +
"  yy   - Yank current line into paste buffer\n" +
"  yw   - Yank from cursor to end of word into paste buffer\n" +
"  y    - When in VISUAL non-block mode yank highlighed area\n" +
"  y    - When in VISUAL     block mode yank highlighed block\n" +
"  Y    - When in VISUAL non-block mode yank highlighed lines\n" +
"  Y    - When in VISUAL     block mode yank highlighed block\n" +
"  zt   - Reposition cursor line to top of window\n" +
"  z<CR>- Move current cursor line to top of window\n" +
"  zz   - Move current cursor line to middle of window\n" +
"  zb   - Move cursor line to bottom of window\n" +
"\n" +
"Special character list of commands:\n" +
" <ESC> - Leave insert mode\n" +
"  ~    - Toggle case\n" +
"  $    - Go to end of line\n" +
"  %    - Jump to opposite {, [, (, }, ] or ) bracket\n" +
"  {    - Jump to enclosing { bracket\n" +
"  }    - Jump to enclosing } bracket\n" +
"  *    - Highlight and search for identifier under cursor\n" +
"  ;    - Repeat last 'f' search\n" +
"  :    - Go to command prompt\n" +
"  .    - Repeat last command\n" +
"  /    - Enter search pattern a prompt\n" +
"\n" +
"Alphabetic list of letters not mapped to commands:\n" +
"  C, E, I, K, q, S, t, T, X, Z\n" +
"\n"

