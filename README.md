# hot

## Installation

```
$ go get github.com/naoty/hot
```

## Usage

`hot` command lists files in order of the commited count.

```
$ cd src/github.com/naoty/Timepiece
$ hot
24: README.md
17: Sources/NSDate+Timepiece.swift
15: Tests/NSDate+TimepieceTests.swift
10: Timepiece.xcodeproj/project.pbxproj
9: Timepiece.podspec
7: Sources/Duration.swift
7: Tests/DurationTests.swift
7: Tests/Int+TimepieceTests.swift
6: Sources/Int+Timepiece.swift
4: Sources/NSDateComponents+Timepiece.swift
2: .travis.yml
2: Tests/NSTimeInterval+TimepieceTests.swift
2: Timepiece.xcodeproj/xcshareddata/xcschemes/Timepiece OSX.xcscheme
2: Sources/NSTimeInterval+Timepiece.swift
1: Timepiece.xcodeproj/xcshareddata/xcschemes/Timepiece iOS.xcscheme
1: Timepiece.xcworkspace/contents.xcworkspacedata
1: .gitignore
1: LICENSE
1: Sources/NSCalendar+Timepiece.swift
1: Sources/NSCalendarUnit+Timepiece.swift
1: Sources/String+Timepiece.swift
1: Tests/NSCalendarUnit+TimepieceTests.swift
1: Tests/String+TimepieceTests.swift
1: Timepiece.playground/Contents.swift
1: Timepiece.playground/Sources/SupportCode.swift
1: Timepiece.playground/contents.xcplayground
1: Timepiece.playground/playground.xcworkspace/contents.xcworkspacedata
1: Timepiece.xcodeproj/Timepiece-Info.plist
1: Timepiece.xcodeproj/TimepieceTests-Info.plist
1: Timepiece.xcodeproj/project.xcworkspace/contents.xcworkspacedata
```

`--number` or `-n` flag enables to limit the number of files.

```
$ hot -n 5
24: README.md
17: Sources/NSDate+Timepiece.swift
15: Tests/NSDate+TimepieceTests.swift
10: Timepiece.xcodeproj/project.pbxproj
9: Timepiece.podspec
```

Given a pattern, only matched files are listed.

```
$ hot -n 5 '**/*.swift'
17: Sources/NSDate+Timepiece.swift
15: Tests/NSDate+TimepieceTests.swift
7: Tests/Int+TimepieceTests.swift
7: Sources/Duration.swift
7: Tests/DurationTests.swift
```

## Author

[naoty](https://github.com/naoty)

