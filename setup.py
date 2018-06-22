import sys
from setuptools import setup, find_packages

if sys.version_info < (3, 5, 0):
    typing = ["typing"]
else:
    typing = []

setup(
    name="snippet_ranger",
    description="Part of source{d}'s stack for machine learning on source code. Provides API and "
                "tools to train and use models for ecosystem exploratory snippet mining. "
                "It can help you to learn new libraries faster and speed up coding speed. "
                "The module allows you to train and use hierarchical topic model on top of "
                "babelfish UAST for any library you want.",
    version="0.0.1-alpha",
    license="Apache 2.0",
    author="source{d}",
    author_email="machine-learning@sourced.tech",
    url="https://github.com/src-d/snippet-ranger",
    download_url="https://github.com/src-d/snippet-ranger",
    packages=find_packages(exclude=("snippet_ranger.tests",)),
    keywords=["machine learning on source code", "github", "topic modeling",
              "hierarchical topic modeling", "exploratory code search"],
    entry_points={
        "console_scripts": ["snippet_ranger=snippet_ranger.__main__:main"],
    },
    install_requires=["pandas>=0.20",
                      "ast2vec>=0.2.6-alpha"] + typing,
    package_data={"": ["LICENSE", "README.md"]},
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Environment :: Console",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: POSIX",
        "Programming Language :: Python :: 3.4",
        "Programming Language :: Python :: 3.5",
        "Programming Language :: Python :: 3.6",
        "Topic :: Software Development :: Libraries"
    ]
)
