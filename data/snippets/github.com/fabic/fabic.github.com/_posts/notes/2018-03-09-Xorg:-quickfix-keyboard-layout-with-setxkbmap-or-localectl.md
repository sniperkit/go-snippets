---
layout: note
title: "Xorg: quickfix keyboard layout with setxkbmap or localectl"
tagline: "Fix keyboard layout from the command line using setxkbmap or localectl."
published: true
category : notes
tags : [note, Linux, Xorg, cli]
---


Query current XKB keyboard states with :

```bash
$ setxkbmap -query
```

One-liner using `setxkbmap` :

```bash
$ setxkbmap -model tm2030USB \
            -layout us       \
            -variant ,altgr-intl,alt-intl        \
            -option "" -option grp:shifts_toggle \
            -verbose \
  && echo
  && setxkbmap -query
```

__Note:__ the trailing `,` on the `-variant ,alt-intl ...` argument, means `""` empty
string => use the default variant for the specified layout.


'tis better to use `localectl` nowadays, for it will update a conf. file
located at `/etc/X11/xorg.conf.d/00-keyboard.conf` :

```bash
$ localectl set-x11-keymap us tm2030USB ,altgr-intl,alt-intl grp:shifts_toggle
```

List supported keyboard models, layouts, layout variants and options :

```bash
$ localectl list-x11-keymap-models
$ localectl list-x11-keymap-layouts
$ localectl list-x11-keymap-variants us
$ localectl list-x11-keymap-options
```

References :

* <https://wiki.archlinux.org/index.php/Keyboard_configuration_in_Xorg>
* <https://wiki.archlinux.org/index.php/X_KeyBoard_extension>

__EOF__
