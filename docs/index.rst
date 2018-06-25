.. Astra documentation master file, created by
   sphinx-quickstart on Tue Jul 14 22:36:35 2015.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

.. |sdkname| replace:: Astra
.. |numberofsections| replace:: three

The |sdkname| SDK
====================================
Welcome to the book on the |sdkname| SDK! The |sdkname| SDK's exposes to the developer a comprehensive set of functionality to utilize all of the Astra's capabilities through a simple yet powerful API.

This book will be a vital resource to any developer who wants to:

- Quickly get up to speed with |sdkname| SDK development
- Understand the high-level concepts with the SDK
- Find detailed information on specific SDK functionality or features

"The |sdkname| SDK" begins below with an introduction, followed by these |numberofsections| sections:

.. toctree::
   :titlesonly:
   :maxdepth: 1
   :includehidden:

   Getting Started <gettingstarted>
   Concepts <concepts>
   Retrieving Stream Data <retrievingstreamdata>

Once you've read the introduction, jump over to the :doc:`Getting Started <gettingstarted>` section to explore an introductory project.

Introduction
============
The design goals for the SDK are straightforward and lofty - to create the world's best 3D sensor development experience. To facilitate this, the best ideas from many other modern device and sensor SDKs, along with many original ideas from Orbbec's development team, were combined to put emphasis on developer creativity, remove needless complexity, and generally make |sdkname| SDK development a joyful, unintrusive process.

Low Ceremony Design
===================
The designers of the SDK are also experienced 3D sensor developers themselves, and recognize the features that are most commonly used when developing 3D sensor applications. To this end, the |sdkname| SDK reduces the boilerplate code required to obtain sensor data by exposing high-level stream types in addition to low-level stream types through a consistent, easy-to-understand API. This can be best shown with some example code:

.. code-block:: c++
   :linenos:

   astra::Astra::initialize();

   astra::StreamSet streamSet; //By default, a StreamSet will address the Astra
   astra::StreamReader reader = streamSet.create_reader();

   // Low-Level Streams
   reader.stream<astra::depthstream>().start();
   reader.stream<astra::ColorStream>().start();

   //High-Level Streams
   reader.stream<astra::handstream>().start();
   reader.stream<astra::pointstream>().start();

   astra::frame frame = reader.get_latest_frame();

   //Low-Level Streams
   const auto depthFrame = frame.get<astra::depthframe>();
   const auto colorFrame = frame.get<astra::colorframe>();

   //High-Level Streams
   const auto handFrame = frame.get<astra::handframe>();
   const auto pointFrame = frame.get<astra::pointframe>();

   astra::Astra::terminate();

First, take note where we create our ``StreamSet`` object on line 1. Without any additional configuration, it will default to addressing the Astra sensor. In the lines below, you can see that, regardless of the stream type, the same API calls are made to interact with each stream.

If you're an advanced 3D sensor developer and worried that this SDK is designed for simple applications only, rest assured that the |sdkname| SDK also handles less-common cases, often simply by passing different parameters to object constructors. In short, no matter how complicated your application might be, the |sdkname| SDK can handle it.

Cross-Everything
================
The |sdkname| SDK exposes its functionality through a core C API and also provides a modern C++11 API. From this foundation, the SDK is able to support a number of languages and platforms:

Currently Supported
-------------------
- Languages: C, C++11
- Platforms: Windows, OSX, Linux, Android

Planned Support
---------------
- Languages: Java, C#
- Frameworks: Processing, OpenFrameworks, Unity, Cinder

Mobile
======
Mobile devices are becoming more powerful and smaller each passing year, but are comparatively resource contrained to desktop PCs. With this in mind, the |sdkname| SDK has been carefully engineered to minimize CPU and memory overhead to keep mobile devices running smoothly while economizing battery usage.

Extensibility
=============
Perhaps the most powerful feature of the |sdkname| SDK is its plugin layer. Plugins afford advanced developers a method to extend the SDK's functionality with the same consistency and level of support as the rest of the platform. In fact, all of the "stock" stream types included with the SDK are actually supported through plugins themselves.

Plugins aren't limited to adding additional stream type support, though. Imagine adding support to address a network of sensors through the same consistent API as a single sensor. The possibilities are truly innumerable, and will allow exciting first and third party additions as the SDK matures.

.. Indices and tables
   ==================
   * :ref:`genindex`
   * :ref:`modindex`
   * :ref:`search`
