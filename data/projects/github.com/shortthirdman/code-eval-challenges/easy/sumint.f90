program sumint
  implicit none
  character(32) :: arg
  character(10) :: str
  integer :: stat,i,n
  n = 0
  call get_command_argument(1, arg)
  open (11,file=trim(arg),action='read')
  do
    read (11,'(a)',iostat=stat) str
    if (stat /= 0) exit
    if (str /= '') then
      read (str,'(i10)') i
      n = n + i
    end if
  end do
  write (*,'(i0)') n
  close (11)
end program sumint
