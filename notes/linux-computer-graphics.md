---
layout: page
title: "Linux: Computer Graphics (X11/XCB, Wayland, OpenGL, Vulkan, etc)"
description: "Notes about drawing stuff on screen under Linux."
category : notes
published: true
---

## OpenGL / OpenGL ES

* [OpenGL ES official doc. & specs | Khronos OpenGL ES Registry](https://www.khronos.org/registry/OpenGL/index_es.php)
* [OpenGL® ES 3.0 Reference Pages (API doc.)](https://www.khronos.org/registry/OpenGL-Refpages/es3.0/)
* [Book: OpenGL® ES 3.0 Programming Guide](http://www.opengles-book.com)
    - [Samples source code @ danginsburg/opengles3-book](https://github.com/danginsburg/opengles3-book/)

### glxinfo : find out which versions of OpenGL are supported

```bash
$ glxinfo |grep -i version
```

```
server glx version string: 1.4
client glx version string: 1.4
GLX version: 1.4
    Version: 17.3.1
    Max core profile version: 3.3
    Max compat profile version: 3.0
    Max GLES1 profile version: 1.1
    Max GLES[23] profile version: 3.0
OpenGL core profile version string: 3.3 (Core Profile) Mesa 17.3.1
OpenGL core profile shading language version string: 3.30
OpenGL version string: 3.0 Mesa 17.3.1
OpenGL shading language version string: 1.30
OpenGL ES profile version string: OpenGL ES 3.0 Mesa 17.3.1
OpenGL ES profile shading language version string: OpenGL ES GLSL ES 3.00
```

## X11

### XCB

* [xcb.freedesktop.org](https://xcb.freedesktop.org/)
* [XCB tutorial](https://xcb.freedesktop.org/tutorial/)
* <https://www.x.org/releases/current/doc/>

* [XCB XRender API @ freedesktop.org](https://xcb.freedesktop.org/manual/group__XCB__Render__API.html)
* <https://github.com/conformal/spectrwm/blob/master/spectrwm.c>
* <https://en.wikipedia.org/wiki/X_Rendering_Extension>
* [libXrender.txt - 2002-11-06 – Keith Packard](https://www.x.org/releases/current/doc/libXrender/libXrender.txt)



__EOF__
