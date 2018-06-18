
# Tweetcode Golang [![Build Status](https://travis-ci.org/Tweetcode/TweetcodeGolang.svg?branch=master)](https://travis-ci.org/Tweetcode/TweetcodeGolang)

![Tweetcode](https://raw.githubusercontent.com/Tweetcode/TweetcodeGolang/assets/assets/gophers_crowd.png)

Tweetcode Golang allows to post code snippets to [@GolangSnippets twitter](https://twitter.com/GolangSnippets) via Travis CI.
Follow Contribution Flow below to tweet your snippet.

## Contribution Flow

### Setup local repo (fist time only)

1. Fork this repo and clone it
2. Run `git config --local include.path ../.gitconfig`

### Add your snippet

1. Create a new branch based on master
2. Create new file with code snippet in `snippets` directory:
   1. It must have .go extension
   2. Give it meaningful name
   3. Make sure it's `UTF-8`-encoded file
3. Commit your changes
   1. Commit message must follow specific format (template provided)
4. Run `make` to preview generated tweet
   1. If there is no errors it'll print generated tweet text and name of
      generated image
4. Push your changes
6. Open Pull Request
   1. In case of failed CI build - fix locally, squash and push again
7. After PR is being merged to master, new snippet will be automatically
   published to [@GolangSnippets twitter](https://twitter.com/GolangSnippets)

## License

MIT
