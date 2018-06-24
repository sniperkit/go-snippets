while read x1 y1 x2 y2 || [ -n "$y2" ]; do
  if [ $x1 -eq $x2 ]; then
    if [ $y1 -eq $y2 ]; then
      echo "here"
    elif [ $y1 -lt $y2 ]; then
      echo "N"
    else
      echo "S"
    fi
  elif [ $y1 -eq $y2 ]; then
    if [ $x1 -lt $x2 ]; then
      echo "E"
    else
      echo "W"
    fi
  elif [ $x1 -lt $x2 ]; then
    if [ $y1 -lt $y2 ]; then
      echo "NE"
    else
      echo "SE"
    fi
  else
    if [ $y1 -lt $y2 ]; then
      echo "NW"
    else
      echo "SW"
    fi
  fi
done <$1
