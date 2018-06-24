program endian
  implicit none
  if (ichar(transfer(1, 'a')) /= 0) then
    write (*,'(a)') "LittleEndian"
  else
    write (*,'(a)') "BigEndian"
  end if
end program endian
