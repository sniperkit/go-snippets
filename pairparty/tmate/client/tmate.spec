Name:           tmate
Version:        1.8.10
Release:        1%{?dist}

Summary:        Instant terminal sharing
License:        MIT
Url:            http://tmate.io

BuildRequires:  autoconf
BuildRequires:  libtool
BuildRequires:  pkgconfig
BuildRequires:  ruby
BuildRequires:  libevent-devel
BuildRequires:  openssl-devel
BuildRequires:  ncurses-devel
BuildRequires:  zlib-devel
BuildRequires:  libssh-devel >= 0.6.0
BuildRequires:  msgpack-devel >= 0.5.8

Source0:        https://github.com/nviennot/tmate/archive/1.8.10.zip

Patch0:         tmate-1.8.10_use_system_libs.patch

%description
Tmate is a fork of tmux providing an instant pairing solution.

%prep
%setup -q -n %{name}-%{version}

%patch0 -p1 -b .tmate-1.8.10_use_system_libs.patch

%build
./autogen.sh
%configure
make %{?_smp_mflags}

%install
make DESTDIR=%{buildroot} install

%files
%defattr(-,root,root)
%doc CHANGES FAQ README-tmux README.md
%{_bindir}/tmate
%{_mandir}/man1/tmate.1*

%changelog

* Sat Jan 24 2015 - Kevin Mulvey <kmulvey@linux.com> - 1.8.10-1
- The big bang.
