@ECHO OFF
REM Should be called from one of the gen_vs scripts

SETLOCAL EnableExtensions EnableDelayedExpansion

SET ASTRADIR=%~dp0

ECHO INFO: Hold on to your butts.

IF NOT EXIST %SFMLPATH% (
   ECHO WARN: SFML -^> %SFMLPATH% does not exist
)

IF NOT EXIST %PROTOPATH% (
   ECHO WARN: PROTOBUF -^> %PROTOPATH% does not exist
)

IF %ARCH% == x86 GOTO ARCH_X86
IF %ARCH% == x64 GOTO ARCH_X64
ECHO ERROR: Unknown Arch: %ARCH% & GOTO FATAL

:ARCH_X86:
SET BITS=32
SET CMAKE_ARCH_OPT=
GOTO MSVC_VERSION
:ARCH_X64:
SET BITS=64
SET CMAKE_ARCH_OPT=-Ax64

:MSVC_VERSION
IF %MSVC% EQU 12 GOTO MSVC12
IF %MSVC% EQU 14 GOTO MSVC14
ECHO ERROR: Unknown MSVC Version %MSVC% & GOTO FATAL

:MSVC12
SET MSVCNAME=2013
GOTO START
:MSVC14
SET MSVCNAME=2015

:START

SET BUILDPATH=%ASTRADIR%build\vs%MSVC%\%ARCH%
SET GENERATOR=Visual Studio %MSVC% %MSVCNAME%

ECHO ---------------------------------------------------------------------------
ECHO INFO: Compiler - %GENERATOR%
ECHO INFO: Arch - %ARCH%
ECHO INFO: Build Path - %BUILDPATH%
ECHO ---------------------------------------------------------------------------

IF EXIST %BUILDPATH% GOTO FIND_CMAKE

MD %BUILDPATH% && GOTO FIND_CMAKE

ECHO ERROR: Could not create working path: %BUILDPATH% & GOTO FATAL

:FIND_CMAKE

WHERE cmake 2> NUL && GOTO CONFIGURE_CMAKE
ECHO ERROR: CMake executable (cmake.exe) could not be found on %%PATH%% & GOTO FATAL

:CONFIGURE_CMAKE

PUSHD %BUILDPATH%

ECHO INFO: Configuring CMake...

cmake -G"%GENERATOR%"^
      -DSFML_ROOT="%SFMLPATH%"^
      -DPROTOBUF_SRC_ROOT_FOLDER="%PROTOPATH%"^
      %CMAKE_ARCH_OPT%^
      %EXTRAARGS%^
      "%ASTRADIR%" && GOTO COPY_DEPENDENCIES

ECHO ERROR: CMake failed. It happens. Fix the errors, run this script again. & GOTO FATAL

:COPY_DEPENDENCIES

ECHO INFO: Copying runtime dependencies
FOR %%G IN (Debug Release RelWithDebInfo MinSizeRel) DO (
    ECHO INFO: Copying for dependencies for %%G configuration...
    XCOPY %ASTRADIR%samples-aux\thirdparty-bin\vs%MSVCNAME%\copy_to_bin%BITS%_dir\*.* bin\%%G /E /Y /I /Q || GOTO COPY_ERROR
)

GOTO SUCCESS

:COPY_ERROR
ECHO ERROR: Unable to copy runtime dependencies & GOTO FATAL

:FATAL
ECHO FATAL: Something went wrong^^! Check above for details.
POPD
EXIT /B 255

:SUCCESS
POPD
ECHO INFO: Success^^! Done^^!
