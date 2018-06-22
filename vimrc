call pathogen#infect()
call pathogen#helptags()

source ~/.vim/myconfig/general.vim
source ~/.vim/myconfig/fcitx.vim
source ~/.vim/myconfig/program.vim
source ~/.vim/myconfig/markdown.vim
source ~/.vim/myconfig/syntastic.vim

autocmd Filetype tex setl updatetime=1
let g:livepreview_previewer = 'evince'

