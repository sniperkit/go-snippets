set -ex
printenv
/usr/local/bin/opencv_version
python --version
find / -name "Python.h"
echo -e 'from numpy  import *\nx = arange(15).reshape(3, 5)\nprint x' | python
echo -e 'import cv2\nprint cv2.__version__' | python
