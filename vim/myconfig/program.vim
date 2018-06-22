
function! Ih()
    let str=expand('%:t:r').'_h'
    call append(0, '#ifndef '.str)
    call append(1, '#define '.str)
    call append(line('.'), '#endif // end of '.str)
endfunction

function! Ins(name)
    let str=a:name
    let i = line('.')
    call append(i, 'namespace '.str.' {')
    call append(i+1, '')
    call append(i+2, '} // end of namespace '.str)
    call cursor(i+1, 0)
endfunction


function! Id()
    let i = line('.')
    call append(i+0, '#ifdef DEBUG')
    call append(i+1, 'std::cout << ""  << std::endl;')
    call append(i+2, '#endif')
    call append(i+3, '')
    call cursor(i+2, 0) 
    execute 'normal =='
    execute 'normal 2f"'
endfunction

function! Ic(class_name)
    let str=a:class_name
    let i = line('.')
    call append(i+1, 'class '.str)
    call append(i+2, '{')
    call append(i+3, 'public:')
    call append(i+4, str.'()')
    call append(i+5, '{')
    call append(i+6, '')
    call append(i+7, '}')
    call append(i+8, 'private:')
    call append(i+9, '')
    call append(i+10, '};')
    call append(i+11, '')
    call cursor(i+5, 0)
    execute '.,/}/normal =='
endfunction

function! Im()
    let i = line('.')
    call append(i+1, '#include <iostream>')
    call append(i+2, '')
    call append(i+3, '')
    call append(i+4, 'int main()')
    call append(i+5, '{')
    call append(i+6, '')
    call append(i+7, 'return 0;')
    call append(i+8, '}')
    call cursor(i+8, 0)
    execute 'normal =='
    call cursor(i+3, 0)
endfunction
