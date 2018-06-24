program age_distribution
  implicit none
  character(32) :: arg
  character(10) :: str
  integer :: stat,n
  call get_command_argument(1, arg)
  open (11,file=trim(arg),action='read')
  do
    read (11,'(a)',iostat=stat) str
    if (stat /= 0) exit
    if (str /= '') then
      read (str,'(i10)') n
      if ((n .lt. 0) .or. (n .gt. 100)) then
        write (*,'(a)') "This program is for humans"
      else if (n .lt. 3) then
        write (*,'(a)') "Still in Mama's arms"
      else if (n .lt. 5) then
        write (*,'(a)') "Preschool Maniac"
      else if (n .lt. 12) then
        write (*,'(a)') "Elementary school"
      else if (n .lt. 15) then
        write (*,'(a)') "Middle school"
      else if (n .lt. 19) then
        write (*,'(a)') "High school"
      else if (n .lt. 23) then
        write (*,'(a)') "College"
      else if (n .lt. 66) then
        write (*,'(a)') "Working for the man"
      else
        write (*,'(a)') "The Golden Years"
      end if
    end if
  end do
  close (11)
end program age_distribution
