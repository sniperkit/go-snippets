import sys

if __name__ == '__main__':
  for path in sys.argv[1:]:
    with open(path) as f:
      with open(path + '.csv', 'w') as out:
        out.write('000\t ')
        for line in f:
          out.write(line.strip() + ' ')

        out.write('\n')
          
