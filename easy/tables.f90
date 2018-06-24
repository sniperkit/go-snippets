program tables
  implicit none
  integer :: i,j
  do i=1,12
    do j=1,12
      write (*,'(i4)',advance='NO') i * j
    end do
    print *
  end do
end program tables
