let r = ref 1
in while !r < 100 do
     print_int !r ;
     print_newline() ;
     r := !r+2
   done ;;
