program filesize
  implicit none
  character(32) :: arg
  integer, dimension(13) :: buff
  call get_command_argument(1, arg)
  call stat(trim(arg), buff)
  write (*,'(i0)') buff(8)
end program filesize
