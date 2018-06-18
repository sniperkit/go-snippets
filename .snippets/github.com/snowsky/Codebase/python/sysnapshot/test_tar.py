import os
import glob
import tarfile

def main():
    directories = [ dirname for dirname in glob.glob('*') if not dirname.endswith('.tar.gz') ]

    for directory in directories:
        print('working on {}'.format(directory))
        for root, dirs, filenames in os.walk(directory):
            tar = tarfile.open('{}.tar.gz'.format(directory), 'w:gz')
            for filename in filenames:
                tar.add(os.path.join(root, filename))
            tar.close()

main()
