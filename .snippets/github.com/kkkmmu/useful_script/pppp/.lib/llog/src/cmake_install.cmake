# Install script for directory: /home/leeweop/kkkmmu/pppp/lib/llog/src

# Set the install prefix
if(NOT DEFINED CMAKE_INSTALL_PREFIX)
  set(CMAKE_INSTALL_PREFIX "/usr/local")
endif()
string(REGEX REPLACE "/$" "" CMAKE_INSTALL_PREFIX "${CMAKE_INSTALL_PREFIX}")

# Set the install configuration name.
if(NOT DEFINED CMAKE_INSTALL_CONFIG_NAME)
  if(BUILD_TYPE)
    string(REGEX REPLACE "^[^A-Za-z0-9_]+" ""
           CMAKE_INSTALL_CONFIG_NAME "${BUILD_TYPE}")
  else()
    set(CMAKE_INSTALL_CONFIG_NAME "")
  endif()
  message(STATUS "Install configuration: \"${CMAKE_INSTALL_CONFIG_NAME}\"")
endif()

# Set the component getting installed.
if(NOT CMAKE_INSTALL_COMPONENT)
  if(COMPONENT)
    message(STATUS "Install component: \"${COMPONENT}\"")
    set(CMAKE_INSTALL_COMPONENT "${COMPONENT}")
  else()
    set(CMAKE_INSTALL_COMPONENT)
  endif()
endif()

# Install shared libraries without execute permission?
if(NOT DEFINED CMAKE_INSTALL_SO_NO_EXE)
  set(CMAKE_INSTALL_SO_NO_EXE "0")
endif()

if("${CMAKE_INSTALL_COMPONENT}" STREQUAL "Unspecified" OR NOT CMAKE_INSTALL_COMPONENT)
  foreach(file
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1.1"
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1"
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so"
      )
    if(EXISTS "${file}" AND
       NOT IS_SYMLINK "${file}")
      file(RPATH_CHECK
           FILE "${file}"
           RPATH "")
    endif()
  endforeach()
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1.1;/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1;/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
file(INSTALL DESTINATION "/home/leeweop/kkkmmu/pppp/.lib/llog/lib" TYPE SHARED_LIBRARY FILES
    "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1.1"
    "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1"
    "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so"
    )
  foreach(file
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1.1"
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so.1"
      "$ENV{DESTDIR}/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.so"
      )
    if(EXISTS "${file}" AND
       NOT IS_SYMLINK "${file}")
      if(CMAKE_INSTALL_DO_STRIP)
        execute_process(COMMAND "/usr/bin/strip" "${file}")
      endif()
    endif()
  endforeach()
endif()

if("${CMAKE_INSTALL_COMPONENT}" STREQUAL "Unspecified" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.a")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
file(INSTALL DESTINATION "/home/leeweop/kkkmmu/pppp/.lib/llog/lib" TYPE STATIC_LIBRARY FILES "/home/leeweop/kkkmmu/pppp/.lib/llog/lib/libllog.a")
endif()

if("${CMAKE_INSTALL_COMPONENT}" STREQUAL "Unspecified" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/home/leeweop/kkkmmu/pppp/.lib/llog/include/llog/llog.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
file(INSTALL DESTINATION "/home/leeweop/kkkmmu/pppp/.lib/llog/include/llog" TYPE FILE FILES "/home/leeweop/kkkmmu/pppp/lib/llog/src/llog.h")
endif()

