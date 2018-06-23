#!/usr/bin/env python
# -*- coding: utf-8 -*-
import subprocess

def hist_bin(x):
  if x < 1.5:
    return 0
  elif x < 2:
    return 3
  elif x < 3:
    return 5
  elif x < 6:
    return 7
  else:
    return 9


uptime = subprocess.check_output(['uptime'])
up = uptime.split(':')[-1].strip().split(' ')[::-1]
loadav = ','.join([str(hist_bin(float(x))) for x in up])
print('{%s}|font="Spark Dot-line Medium" color=#48A9A6 size=11' % loadav)
