import filecmp
filecmp.cmp('file0.txt', 'file1.txt')
filecmp.cmp('file0.txt', 'file00.txt')
filecmp.dircmp('dira', 'dirb').diff_files
filecmp.dircmp('dira', 'dirb').same_files
filecmp.dircmp('dira', 'dirb').report()
import os
dirb = set(os.listdir('/tmp/dirb'))
dira = set(os.listdir('/tmp/dira'))
dira-dirb
dirb-dira
