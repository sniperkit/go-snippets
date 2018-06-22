# sudo apt-get install python-rope python-flake8 
git submodule add http://github.com/tpope/vim-fugitive.git vim/bundle/fugitive
git submodule add https://github.com/msanders/snipmate.vim.git vim/bundle/snipmate
git submodule add https://github.com/tpope/vim-surround.git vim/bundle/surround
git submodule add https://github.com/tpope/vim-git.git vim/bundle/git
git submodule add https://github.com/ervandew/supertab.git vim/bundle/supertab
git submodule add https://github.com/fholgado/minibufexpl.vim.git vim/bundle/minibufexpl
git submodule add https://github.com/wincent/Command-T.git vim/bundle/command-t
git submodule add https://github.com/mileszs/ack.vim.git vim/bundle/ack
git submodule add https://github.com/sjl/gundo.vim.git vim/bundle/gundo
git submodule add https://github.com/fs111/pydoc.vim.git vim/bundle/pydoc
git submodule add https://github.com/alfredodeza/pytest.vim.git vim/bundle/py.test
git submodule add https://github.com/reinh/vim-makegreen vim/bundle/makegreen
git submodule add https://github.com/vim-scripts/TaskList.vim.git vim/bundle/tasklist
git submodule add https://github.com/vim-scripts/The-NERD-tree.git vim/bundle/nerdtree
git submodule add https://github.com/jistr/vim-nerdtree-tabs.git vim/bundle/vim-nerdtree-tabs
git submodule add https://github.com/python-rope/ropevim.git vim/bundle/ropevim
git submodule add https://github.com/nvie/vim-flake8.git vim/bundle/vim-flake8
git submodule add https://github.com/tomasr/molokai.git vim/bundle/molokai
git submodule add https://github.com/scrooloose/syntastic.git vim/bundle/syntastic
git submodule add https://github.com/vim-scripts/vim-cmake-project.git vim/bundle/vim-cmake-project
git submodule add https://github.com/altercation/vim-colors-solarized.git vim/bundle/vim-colors-solarized
git submodule add https://github.com/plasticboy/vim-markdown.git vim/bundle/vim-markdown
git submodule add https://github.com/vim-scripts/SQLComplete.vim.git vim/bundle/sqlcomplete
git submodule add https://github.com/xuhdev/vim-latex-live-preview.git vim/bundle/vim-latex-live-preview
git submodule init
git submodule update
git submodule foreach git submodule init
git submodule foreach git submodule update
