Name:           tmate-server
Version:        1.8
Release:        1%{?dist}

Summary:        Instant terminal sharing
License:        MIT
Url:            http://tmate.io

BuildRequires:  git-core 
BuildRequires:  build-essential 
BuildRequires:  pkg-config 
BuildRequires:  libtool 
BuildRequires:  libevent-dev 
BuildRequires:  libncurses-dev 
BuildRequires:  zlib1g-dev 
BuildRequires:  automake 
BuildRequires:  libssh-dev 
BuildRequires:  cmake 
BuildRequires:  ruby

Source0:        https://github.com/nviennot/tmate-slave/archive/1.8.zip

%description
Tmate is a fork of tmux providing an instant pairing solution.

%prep
%setup -q -n %{name}-%{version}

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

* Sat Jan 24 2015 - Kevin Mulvey <kmulvey@linux.com> - 1.8-1
- The big bang.
