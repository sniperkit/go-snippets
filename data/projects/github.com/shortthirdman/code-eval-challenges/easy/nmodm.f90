program nmodm
  implicit none
  character(32) :: arg
  character(10) :: str
  integer :: stat,n,m,sep
  call get_command_argument(1, arg)
  open (11,file=trim(arg),action='read')
  do
    read (11,'(a)',iostat=stat) str
    if (stat /= 0) exit
    if (str /= '') then
      sep = scan(str,',')
      read (str(1:sep-1),'(i10)') n
      read (str(sep+1:),'(i10)') m
      write (*,'(i0)') n - (n / m) * m
    end if
  end do
  close (11)
end program nmodm
