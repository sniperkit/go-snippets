# fabic.net personal weblog

[![Join the chat at https://gitter.im/fabic-github-com/Lobby](https://badges.gitter.im/fabic-github-com/Lobby.svg)](https://gitter.im/fabic-github-com/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Dude, this is your Github Pages / Jekyll-based weblog for <http://fabic.net>.

Using Jekyll theme [Laynon (by poole)](https://github.com/poole/lanyon/).

## TODO

* <http://www.davidverhasselt.com/about/> inspiration
* <http://raghibm.com/blog/> inspiration
* Integrate [Pocket feed](https://getpocket.com/privacy_controls) ?
  'tis password protected, so you'd have to write some server-side JS somewhere...
* Integrate w/ Diigo somehow ? if their seemingly old RESTful API is still functional.

## Jekyll - Liquid

* <https://jekyllrb.com/docs/variables/>
* <https://shopify.github.io/liquid/basics/introduction/> Liquid language syntax and functions.
* <https://devhints.io/jekyll>
* <https://gist.github.com/smutnyleszek/9803727>
* [Liquid for Designers (Shopify/liquid)](https://github.com/Shopify/liquid/wiki/Liquid-for-Designers)
* [Jekyll & Liquid cheat sheet](https://learn.cloudcannon.com/jekyll-cheat-sheet/)
* [Github: Basic writing and formatting](https://help.github.com/articles/basic-writing-and-formatting-syntax/)

## Installation

### Install Ruby through RVM (optional)

This will install RVM under `~/.rvm/`. the `--ruby` switch wil have Ruby installed too.
Note that it would have messed up with various shell startup scripts (like `.bashrc`, `.bash_profile`, `.profile`, etc),
if not for `--ignore-dotfiles`. Note that it _may_ fallback to compiling Ruby
from source if it can't find binaries for your Linux distribution :-/

    $ gpg --keyserver hkp://keys.gnupg.net \
          --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB

    $ curl -sSL https://get.rvm.io | \
        bash -s stable --ignore-dotfiles --version latest --ruby --with-gems=""

Setup your `$PATH` with `~/.rvm/bin`, and source `~/.rvm/scripts/rvm`.

* See [rvm.io](http://rvm.io/)
* See [ruby-lang.org/en/documentation/installation/](https://www.ruby-lang.org/en/documentation/installation/)

_**todo:** write further RVM related instructions._

### Install Jekyll and dependencies

This will have `bundler` installed under `~/.gem/ruby/2.X.Y/bin/` :

    $ gem install --user-install bundler

This will fetch and install this "project's" dependencies :

    $ bundler

### Run a local instance of Jekyll

    $ bundle exec jekyll clean

    $ bundle exec jekyll serve --incremental --trace --watch --drafts --unpublished --host 0.0.0.0 --port 3939

#### Run with the `jekyll-serve.sh` script

    $ ./jekyll-serve.sh

<hr/>

### Update dependencies

Once in a while run :

    $ bundle update

### Writing plugins

The `_plugins/` directory may contain plugins.

__Note__ that Github won't process these though.

* [Githug/Linguist](https://github.com/github/linguist) does the syntaxt highlighting. See [](https://github.com/github/linguist/blob/master/lib/linguist/languages.yml) for supported languages.  See [Github “basic writing and formatting”](https://help.github.com/articles/basic-writing-and-formatting-syntax/).

## Lanyon

    Lanyon is an unassuming [Jekyll](http://jekyllrb.com) theme that places content first by tucking away navigation in a hidden drawer. It's based on [Poole](http://getpoole.com), the Jekyll butler.

    ![Lanyon](https://f.cloud.github.com/assets/98681/1825266/be03f014-71b0-11e3-9539-876e61530e24.png)
    ![Lanyon with open sidebar](https://f.cloud.github.com/assets/98681/1825267/be04a914-71b0-11e3-966f-8afe9894c729.png)

## Contents

- [Usage](#usage)
- [Options](#options)
  - [Sidebar menu](#sidebar-menu)
  - [Themes](#themes)
  - [Reverse layout](#reverse-layout)
- [Development](#development)
- [Author](#author)
- [License](#license)


## Usage

Lanyon is a theme built on top of [Poole](https://github.com/poole/poole), which provides a fully furnished Jekyll setup—just download and start the Jekyll server. See [the Poole usage guidelines](https://github.com/poole/poole#usage) for how to install and use Jekyll.


## Options

Lanyon includes some customizable options, typically applied via classes on the `<body>` element.


### Sidebar menu

Create a list of nav links in the sidebar by assigning each Jekyll page the correct layout in the page's [front-matter](http://jekyllrb.com/docs/frontmatter/).

```
---
layout: page
title: About
---
```

**Why require a specific layout?** Jekyll will return *all* pages, including the `atom.xml`, and with an alphabetical sort order. To ensure the first link is *Home*, we exclude the `index.html` page from this list by specifying the `page` layout.


### Themes

Lanyon ships with eight optional themes based on the [base16 color scheme](https://github.com/chriskempson/base16). Apply a theme to change the color scheme (mostly applies to sidebar and links).

![Lanyon with red theme](https://f.cloud.github.com/assets/98681/1825270/be065110-71b0-11e3-9ed8-9b8de753a4af.png)
![Lanyon with red theme and open sidebar](https://f.cloud.github.com/assets/98681/1825269/be05ec20-71b0-11e3-91ea-a9138ef07186.png)

There are eight themes available at this time.

![Available theme classes](https://f.cloud.github.com/assets/98681/1817044/e5b0ec06-6f68-11e3-83d7-acd1942797a1.png)

To use a theme, add any one of the available theme classes to the `<body>` element in the `default.html` layout, like so:

```html
<body class="theme-base-08">
  ...
</body>
```

To create your own theme, look to the Themes section of [included CSS file](https://github.com/poole/lanyon/blob/master/public/css/lanyon.css). Copy any existing theme (they're only a few lines of CSS), rename it, and change the provided colors.


### Reverse layout

![Lanyon with reverse layout](https://f.cloud.github.com/assets/98681/1825265/be03f2e4-71b0-11e3-89f1-360705524495.png)
![Lanyon with reverse layout and open sidebar](https://f.cloud.github.com/assets/98681/1825268/be056174-71b0-11e3-88c8-5055bca4307f.png)

Reverse the page orientation with a single class.

```html
<body class="layout-reverse">
  ...
</body>
```


### Sidebar overlay instead of push

Make the sidebar overlap the viewport content with a single class:

```html
<body class="sidebar-overlay">
  ...
</body>
```

This will keep the content stationary and slide in the sidebar over the side content. It also adds a `box-shadow` based outline to the toggle for contrast against backgrounds, as well as a `box-shadow` on the sidebar for depth.

It's also available for a reversed layout when you add both classes:

```html
<body class="layout-reverse sidebar-overlay">
  ...
</body>
```

### Sidebar open on page load

Show an open sidebar on page load by modifying the `<input>` tag within the `sidebar.html` layout to add the `checked` boolean attribute:

```html
<input type="checkbox" class="sidebar-checkbox" id="sidebar-checkbox" checked>
```

Using Liquid you can also conditionally show the sidebar open on a per-page basis. For example, here's how you could have it open on the homepage only:

```html
<input type="checkbox" class="sidebar-checkbox" id="sidebar-checkbox" {% if page.title =="Home" %}checked{% endif %}>
```

## Development

Lanyon has two branches, but only one is used for active development.

- `master` for development.  **All pull requests should be to submitted against `master`.**
- `gh-pages` for our hosted site, which includes our analytics tracking code. **Please avoid using this branch.**


## Author

**Mark Otto**
- <https://github.com/mdo>
- <https://twitter.com/mdo>


## License

Open sourced under the [MIT license](LICENSE.md).

<3
