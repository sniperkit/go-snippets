from setuptools import setup
from subprocess import check_output

with open('VERSION.txt', 'r') as content_file:
    version = content_file.read()

    setup(
        name='protobuf-qtypes',
        version=version[1:],
        description='qtypes data structures',
        url='http://github.com/piotrkowalczuk/qtypes',
        author='Piotr Kowalczuk',
        author_email='p.kowalczuk.priv@gmail.com',
        license='MIT',
        packages=['qtypes'],
        install_requires=[
            'protobuf'
        ],
        zip_safe=False,
        keywords=['protobuf', 'data-structures'],
        download_url='https://github.com/piotrkowalczuk/qtypes/archive/%s.tar.gz' % version.rstrip(),
      )