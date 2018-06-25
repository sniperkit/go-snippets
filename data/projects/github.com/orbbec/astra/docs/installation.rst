.. |sdkname| replace:: Astra

****************
2.1 Installation
****************

System Requirements
===================
Windows
-------
- Microsoft® Windows® 10/8.1/8/7 (32 or 64-bit)
- 150 MB hard disk space
- USB 2.0 Port
- Microsoft® Visual Studio 2013 Community Edition and up

Mac OS X
^^^^^^^^
- Mac® OS X® 10.8.5 or higher
- 150 MB hard disk space
- USB 2.0 Port

..
  Linux
  ^^^^^
  - GNOME or KDE desktop
  - 150 MB hard disk space
  -

Android
-------
- Version 4.4 (Lollipop) and up
- ARM
- USB 2 OTG

Sensor Requirements
===================
- Astra

Installing the SDK
======================
In order to begin installing the |sdkname| SDK, we'll need to grab the latest copy. You can do so `here <https://orbbec3d.com/develop/>`_. Once you've downloaded the SDK archive, extract it to a folder on your disk that you'll remember.

Building the Samples
========================
Windows
-------
#. In an Explorer window, locate |sdkname|-samples.sln under /samples/vs2013 or /samples/vs2015 inside the folder where the SDK was extracted.
#. Open |sdkname|-samples.sln in Visual Studio 2013 or Visual Studio 2015, as appropriate.
#. In Visual Studio, build the solution by either choosing the menu item **BUILD -> Build Solution** or using the **Ctrl+Shift+B** shortcut.
#. Sit tight and in under a minute on most computers the samples should be ready to run.
#. To start a sample in Visual Studio, set any of the sample projects as the **Startup Project** by right-clicking the project and selecting "Set as StartUp Project". Then, start debugging either through the menu item **DEBUG -> Start Debugging** or by using the **F5** function key.

OS X
-------
The sample build system uses CMake 3.2+.

Getting CMake
^^^^^^^^^^^^^
- Option 1: You can download it from https://cmake.org/download/
- Option 2: If you are a homebrew user (https://brew.sh), you can install the latest CMake version by running the following in a terminal::

    $ brew install cmake

Compiling the samples
^^^^^^^^^^^^^^^^^^^^^
#. Open a terminal and change directories to the samples/ directory within the SDK folder.
#. Type ./build_samples.sh
#. Now the samples binaries are build to samples/build/bin/ and you can run then in the terminal. For example::

    $ cd build/bin
    $ ./SimpleStreamViewer-SFML
