Astra SDK v0.5.0
Copyright (c) 2015-2016 Orbbec
https://www.orbbec3d.com

For help and support, visit us at https://3dclub.orbbec3d.com.

What's New
==========

v0.5.0 2016/04/26

This release cleans up the API and library organization a bit. There are a few breaking changes in this release
from v0.4 but they should be relatively quick to update existing code.

* BREAKING: Library names:
  * astra -> astra_core
  * astraul -> astra
  * astra_api -> astra_core_api
* BREAKING: C++ API stylistic changes: standardize on namespace::ClassName::method_name
* BREAKING: Besides the predictable stylistic changes, a few C++ method renames:
  * DepthStream/ColorStream/etc horizontalFieldOfView -> hFov, verticalFieldOfView -> vFov
  * DepthFrame/ColorFrame/etc resolutionX/resolutionY -> width/height
  * DepthFrame/ColorFrame/etc numberOfPixels -> length
* BREAKING: C++ header filenames renamed according to namespace.hpp or Class.hpp.
  * Main header to include: <Astra/Astra.h> & <AstraUL/AstraUL.h> -> just <astra/astra.hpp>.
  * Don't need to explicitly include astra_core.hpp.
* BREAKING: astra::Astra::{initialize(),terminate()} -> astra::{initialize(),terminate()}
* Fix: Cycling start/stop on a stream multiple times no longer crashes
* Fix: SXGA depth & color support now work
* Enhancement: (ALL) Add const as appropriate in the C++ API, makes passing around references to frame types possible
* Enhancement: (ALL) StreamSet and StreamReader have improved copy-semantics and default ctor,
                     allowing simpler storage as a class field.
* Enhancement: (ALL) Add MultiSensorViewer-SFML sample demonstrating multi-sensor support
* Enhancement: (ALL) Removed the dependency on OpenCV and reimplemented the necessary functionality internally.
* Enhancement: (Windows) Add VS2015 support
* Enhancement: (OSX) libusb, SFML are now distributed with the SDK, and rpaths are automatically setup.
                     No more brew installing! (except maybe CMake)
* Enhancement: (OSX) Add build_samples.sh script in samples/, see ./build_samples.sh -h for help
* Enhancement: (OSX) Building the samples also copies the Astra SDK and SFML binaries to the lib dir,
                     allowing simpler copy-only deployment: Just copy bin/ and lib/.

v0.4.0 2015/10/14

* Add official support for Win64 and OS X 10.8+
* Updated SFML to 2.3.2
* Added features to SimpleDepthViewer-SFML and SimpleStreamViewer-SFML: pausing,
  overlay color on depth, display depth data under the mouse in text overlay. (See keyboard shortcut section below.)
* Minor internal bug fixes

v0.3.0 2015/09/14

* Rename to Astra SDK.
* Rename Sensor to StreamSet in C++ API.
* Various bug fixes and internal enhancements.
* New samples:  SimpleStreamViewer-SFML, SimpleColorViewer-SFML, ColorReaderEvent, ColorReaderPoll
* Samples have improved performance.
* Add IR stream, mirrored depth, and registered depth support.
* VS2013 samples solution no longer requires copying files - compiles and runs out of the box.
* StreamReader start() and stop() are functional now. See SimpleStreamViewer-SFML.
* Add initial getting started documentation

v0.2.1 2015/07/06 Updated Android and Windows drivers for new sensor USB IDs. Add Android test app .apk.

v0.2.0 2015/07/03 First version ready for external use.

Pre-built samples
=================

Pre-built samples are included in the bin/ directory.
Simply plug in your sensor and then run any of the executable files in the bin/ directory.

We recommend starting with SimpleStreamViewer-SFML and SimpleHandViewer-SFML.
In the hand viewer, wave left and right at the sensor a few times to start hand tracking.

OS-specific instructions:

OS X:
=======
Requires:
* OS X 10.8+
* Xcode 6.2+

Windows:
=======
If you want to run the pre-compiled samples and don't have Visual Studio installed,
you must install the Visual C++ Redistributable Packages.
For VS2013: https://www.microsoft.com/en-us/download/details.aspx?id=40784
For VS2015: https://www.microsoft.com/en-us/download/details.aspx?id=48145

Sample keyboard shortcuts
=========================

SimpleStreamViewer-SFML:

* F - toggle fullscreen
* R - toggle registered depth
* M - toggle mirrored streams
* I - enable IR (RGB mode)
* G - enable IR (Gray16 mode)
* C - enable color
* P - toggle pausing the streams
* O - toggle overlay color stream on depth stream

SimpleDepthViewer-SFML:

* F - toggle fullscreen
* R - toggle registered depth
* M - toggle mirrored streams
* P - toggle pausing the streams
* Space bar - toggle text overlay with the depth data under the mouse cursor

Building the samples from source
================================

OS X
====

The sample build system uses CMake 3.2+.

Getting CMake
=============

Option 1:
You can download it from https://cmake.org/download/

Option 2:
If you are a homebrew user (https://brew.sh), you can install the latest
CMake version by running the following in a terminal:
$ brew install cmake

Compiling Samples
=================

Once CMake has been properly installed, compiling the included samples
can be accomplished by enter the following at a terminal prompt:

$ cd /path/to/sdk
$ cd samples
$ ./build_samples.sh

> The samples will be compiled into samples/build/bin/

For additional sample build options e.g. creating an Xcode project or compiling in debug configuration:
$ ./build_samples.sh -h

Windows
=======

Requirements:
* Visual Studio 2013 or later. The VS 2015 Community version is a free download:
https://www.visualstudio.com/en-us/products/visual-studio-community-vs.aspx
* Windows 7 or later. Tested on Windows 10, Windows 8.1, Windows 7.

The provided Visual Studio 2013/2015 solution is already configured to run out of the box.

Steps:
1) Open astra-samples.sln found either in samples/vs2013 or samples/vs2015 folder.
2) Once in Visual Studio, change the Solution Configuration to Release.
3) Build the solution.

> The samples will be compiled into samples/{vs{2013,2015}/bin/{Release,RelWithDebInfo,Debug}
  depending on your version of Visual Studio and the chosen build configuration.

Tips
====

You can exit samples by pressing Control-C. They will catch this signal and exit cleanly.
To exit samples with a GUI window, press Control-C, escape, or simply close the window.

When you start development, we highly recommend using the C++ API (or a higher level language wrapper.)
The C API is better suited for binding to other languages, or in environments where C is preferred.

Documentation
=============

Preliminary documentation in HTML format can be found in the sdk/docs directory.

Known Issues
============

* There is no error message if no sensor is found or plugged in. (But, it doesn't crash!)
  Hot-plugging (plugging in or removing the sensor while an app is running) is not currently supported.

* If a sample crashes or you forcibly terminate an app that uses Astra SDK before it is given a chance to
  properly shutdown the Astra sensor, the sensor driver may exhibit strange behavior.
  To FIX, simply unplug and re-plug the sensor into the USB port. (This is primarily an OpenNI issue.)

* If any samples appear to be running slow on your computer, make sure you are building the samples in
  Release configuration. This is the default if you are using the OSX build_samples.sh script.
