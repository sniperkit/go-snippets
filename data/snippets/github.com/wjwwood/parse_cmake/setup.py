from setuptools import setup, find_packages

setup(
    name='parse_cmake',
    version='0.4.1',
    author='Issac Trotts, William Woodall',
    author_email='itrotts@willowgarage.com, william@osrfoundation.org',
    url='http://github.com/wjwwood/parse_cmake',
    description='Parser for CMakeLists.txt files',
    packages=find_packages(),
    install_requires=['pyPEG2'],
    tests_require=['nose', 'flake8'],
    test_suite='nose.collector',
    classifiers=[
        'Environment :: Console',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License',
        'Programming Language :: Python',
    ],
    license='Apache 2.0',
    zip_safe=False,
    include_package_data=True,
    entry_points={
        'console_scripts': [
            'cmake_pprint = parse_cmake.cmake_pprint:main',
        ]
    }
)
