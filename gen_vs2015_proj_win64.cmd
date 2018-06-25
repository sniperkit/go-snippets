@ECHO OFF

SETLOCAL EnableExtensions EnableDelayedExpansion

REM Update your dependency directories below as necessary
SET SFMLPATH=C:\libs\SFML-2.3.2\vs2015\x64
SET PROTOPATH=C:\libs\protobuf-2.6.1\vs2015

SET MSVC=14
SET ARCH=x64
SET EXTRAARGS=%*

CALL %~dp0\generate.cmd
