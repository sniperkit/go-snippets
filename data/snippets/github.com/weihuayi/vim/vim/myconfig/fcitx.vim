
let g:input_toggle = 1

" 在退出插入模式时，自动切换输入法为英文输入法，便于使用 vim 命令
function! Fcitx2Leave()
   let s:input_status = system("fcitx-remote")
   let g:input_toggle = s:input_status " 记录当前的输入法状态
   if s:input_status == 2
      let l:a = system("fcitx-remote -c")
   endif
endfunction

function! Fcitx2Enter()
   let s:input_status = system("fcitx-remote")
   " 如果当前处输入法不是中文, 而且记录离开前的输入法为中文, 则切换到中文输入法
   if s:input_status != 2 && g:input_toggle == 2 
      let l:a = system("fcitx-remote -o")
   endif
endfunction
set timeoutlen=150
autocmd InsertLeave * call Fcitx2Leave()
autocmd InsertEnter * call Fcitx2Enter()
