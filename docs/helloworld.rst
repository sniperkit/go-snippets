.. |sdkname| replace:: Astra

***************
2.2 Hello World
***************
*Time Required: ~10 minutes*

Want to get a taste of our SDK before diving deeper? Then let's get our hands dirty and write some code!

By the end of this tutorial you should be familiar with:

- Proper initialization and termination of the SDK
- Reading data from the sensor
- Examining the depth information provided by the Astra's depth camera

Before We Begin
===============
If you skipped over the section where we install the SDK and build the sample applications provided by the SDK, make sure you've at least downloaded and extracted |sdkname| to a folder you can easily access.

Getting Down to Business!
=========================
Our first step will be to set up a skeleton application as a starting point for progressively adding new functionality.

#. Using your favorite IDE, set up a new console application project and create a new source file called "main.cpp".
#. Copy the following into your main.cpp file:

.. code-block:: c++
   :linenos:

   #include <astra/astra.hpp>

   #include <cstdio>
   #include <iostream>

   int main(int argc, char** argv)
   {
      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 1 - astra.hpp must be included in all applications. It is the core of |sdkname| and is required for all C++ based |sdkname| applications.
- Lines 9-10 - We'll use `std::cin.get() <http://en.cppreference.com/w/cpp/io/basic_istream/get>`_ to make sure we have an opportunity to see our handiwork before our application closes its window.

Initializing and Terminating |sdkname|
======================================
To prepare |sdkname| to do our bidding, we must first initialize |sdkname|, which is unsurprisingly done via the ``initialize`` function. When we're ready to end our session with the SDK, we then need to give |sdkname| an opportunity to cleanly shutdown. This is accomplished by calling the ``terminate`` function, which will, among other things, properly turn off attached devices.

Add the two new lines below:

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 3,7

   int main(int argc, char** argv)
   {
      astra::initialize();

      // what will go here? you'll find out soon!

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

Trust But Verify
================
Before we get ahead of ourselves, let's take a moment to make sure that everything is as we expect it. Compile and run the application. The application should start up, print out a series of diagnostic messages to the console, and then patiently wait for you to press the "Enter" key. Once pressed, the application should gracefully exit.

.. note::

   |sdkname| by default logs a fair amount of diagnostic information to the console. If you do run into an issue, this can be a great place to start looking for answers.

Next up: Talking to Astra.

Connecting to the Astra
=======================
Now that we know how to properly initialize and terminate |sdkname|, it's time to actually communicate with the Astra sensor. For this, we use the ``StreamSet`` class, which broadly encapsulates the idea of a group of related data sources (think: video and audio from a 2D video camera). For now, however, it's sufficient to think of a streamSet as a physical device like the Astra, and the ``StreamSet`` class, your portal to its functionality.

Between our initialization and termination bookends, let's declare a ``StreamSet`` variable.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 5

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

Now, sure, this seems like a small addition from our previous step, but this line is more significant than it appears. Just by declaring and constructing a ``StreamSet`` object, you are instructing |sdkname| to start the process of connecting to the first available Astra sensor it can locate. Cool, right?

.. note:

   |sdkname| provides an additional constructor that will allow you to connect to a specific Astra sensor.

Now that we're connected, we're ready to do what we came here to do - see through the eyes of the Astra!

Retrieving Sensor Data
=========================
Time to put our ``StreamSet`` object to good use and get some data. To do this, we'll need to read one of the streams that the Astra is providing. Streams contain the data coming from our camera packaged in packets of data called "frames". |sdkname| currently supports a number of types of streams, including depth, color, hand, and point streams.

In order to access streams from the Astra and get to the frames, we'll need a ``StreamReader`` to tap into one of the streams. For the purposes of our application,we're going to focus on the depth stream. This stream gives us the distances of anything that our camera sees in pixels, and those pixels are packaged in a frame.

1. First, let's create a ``StreamReader`` using our ``StreamSet``.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 6

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;
      astra::StreamReader reader = streamSet.create_reader();

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 11 - Creates a ``StreamReader``

2. Next we start the depth stream using the ``StreamReader`` that we created in the previous step. Starting the depth stream tells |sdkname| that we're interested in getting depth data from our ``StreamSet``.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 8

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;
      astra::StreamReader reader = streamSet.create_reader();

      reader.stream<astra::DepthStream>().start();

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 13 - Starts the depth stream

3. With our depth stream stared, let's pull the latest depth frame from our depth stream. To do this, we'll need to first retrieve the latest ``Frame`` through our ``StreamReader``, then call ``get<T>`` to get the depth frame data from our frame.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 10,11

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;
      astra::StreamReader reader = streamSet.create_reader();

      reader.stream<astra::DepthStream>().start();

      astra::Frame frame = reader.get_latest_frame();
      const auto depthFrame = frame.get<astra::DepthFrame>();

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 15 - Retrieves the latest frame
- Line 16 - Gets the depth frame from the latest frame

4. The only remaining task is to print some data from the depth frame that we just retrieved.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 13,14,16-20

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;
      astra::StreamReader reader = streamSet.create_reader();

      reader.stream<astra::DepthStream>().start();

      astra::Frame frame = reader.get_latest_frame();
      const auto depthFrame = frame.get<astra::DepthFrame>();

      const int frameIndex = depthFrame.frame_index();
      const short pixelValue = depthFrame.data()[0];

      std::cout << std::endl
                << "Depth frameIndex: " << frameIndex
                << " pixelValue: " << pixelValue
                << std::endl
                << std::endl;

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 18 - Gets a copy of the frame index from our depth frame
- Line 19 - Gets a copy of the value within the first pixel of our depth frame's data
- Line 21-25 - Prints the two aforementioned values to the console
- Line 27-28 - Pauses execution so we can soak in our success

You can go ahead and run your application now to test that everything works. Just like before, a console window should pop up and display diagnostic information. Then, you should see a line with the frame data that we retrieved. Press enter when you're done.

You just retrieved your first frame from |sdkname|! There's one more task before you graduate from our |sdkname| crash course, and that's working with a sequence of frames.

Consuming a StreamSet stream
============================
Now that you know how to create a ``StreamReader`` and get a frame from it, you're ready to work with a stream of frames. To do this, we only need to make a small change and loop over our call to the ``StreamReader``'s ``get_latest_frame`` function. In this particular case, we're going to get the first 100 frames from our depth stream and print the value of each frame's first pixel to the console.

The following code is highly similar to the code from our last example, except we've added a ``do while`` loop around our frame processing code, in addition to some variables that store the number of times we've looped and the maximum number of frames we want to process.

.. code-block:: c++
   :linenos:
   :lineno-start: 6
   :emphasize-lines: 10,11

   int main(int argc, char** argv)
   {
      astra::initialize();

      astra::StreamSet streamSet;
      astra::StreamReader reader = streamSet.create_reader();

      reader.stream<astra::DepthStream>().start();

      const int maxFramesToProcess = 100;
      int count = 0;

      do {
         astra::frame frame = reader.get_latest_frame();
         const auto depthFrame = frame.get<astra::depthframe>();

         const int frameIndex = depthFrame.frame_index();
         const short pixelValue = depthFrame.data()[0];

         std::cout << std::endl
                   << "Depth frameIndex: " << frameIndex
                   << " pixelValue: " << pixelValue
                   << std::endl
                   << std::endl;

         count++;
      } while (count < maxFramesToProcess);

      std::cout << "Press any key to continue...";
      std::cin.get();

      astra::terminate();

      std::cout << "hit enter to exit program" << std::endl;
      std::cin.get();

      return 0;
   }

- Line 15 - Stores the maximum number of frames we're going to process in the loop
- Line 16 - Sentinel to count the number of frames that we've processed
- Line 18-32 - The frame processing loop

Compile and run. While the program is running and the Astra is focused on you, move around a bit and watch the data values on the frames change.

Achievement get! You've just made your first |sdkname| application! If you haven't had your fill of fun with |sdkname| yet, continue on to our Simple Depth Viewer Tutorial.