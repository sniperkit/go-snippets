# Operation Go: A Routine Mission

## Project Overview
**About the App:**

Operation Go tells the story of a young covert operative on their first mission into the field.  The operative encounters multiple obstacles along the way that must be solved using Golang.

The entire app is rendered using Go's `present` tool.

**How To Play:**

Clone the repo:

`git clone https://github.com/gophergala/operation-go.git`

Download Go's present tool:

`go get code.google.com/p/go.tools/cmd/present`

cd into project directory

`present -base=.`

Visit `http://127.0.0.1:3999` in your browser

Enjoy :)

**Who's It For?:** Go developers of all experience levels

**Why Should They Play?:**

1. To test their Go skills and solve programming puzzles
2. It's a fun, interactive way to experience Go through the context of a story
3. Compete with other Go developers to find the most creative solutions

**Why I'm Building It:**

1. To introduce Go to new developers in a fun, exciting way
2. To learn more about Go in the process and push the boundaries of the present package

**What I Learned:**

1. *Present Package*: Go's present package is amazingly versatile.  With a little JavaScript to customize the layout and output, the present package can be an easy, fun way to spin up custom Go editing in the browser
2. *URLs*: By default, the pages have ugly URL extensions like .article.  In order to mask that, I chose to customize `dir.tpml` in the templates folder as that file renders at the base URL.  With the dir template customized to be the homepage, I loaded the level missions through an iframe.  It had th added benefit of being able to keep the app state and music consistent even as the entire HTML of the mission pages changes.
3. *OAuth*: OAuth.io for GitHub signon worked really well.  Eventually, I might swap it out for something custom though.
4. *Statefulness*: Cookies seemed to be the fastest way to store user data and keep the statefulness of the app.  There's not much data changing hands, just what level the player is on.
5. *Syntax Highlighting*: The syntax highlighting was the hardest part.  Eventually, I settled with Highlight.js.  It was the easiest to work with, but unfortunately couldn't style on the fly as new code was edited (the cursor kept jumping around).  I might be able to fix this in later versions though.  Ace editors for some reason would prevent the source code from executing and I couldn't get Code Mirror to render properly.  In the end, I decided just to highlight the code that you can't edit, as that was better than nothing.
* *Validations*: One tricky part to solve was how to validate the missions. Validating in JavaScript would be too easy for tempting eyes to find the answer.  One option would be to ask the user to spin up the present tool and a separate server on another port.  That seemed too cumbersome to maintain though.  Another option I tried was having the present tool take advantage of the `exec` package to run server commands.  That worked well and had the validations running successfully, but it would require some extra work to send a success signal to the app without the user being able to emulate that same success signal.  In the end, I ended up solving the puzzles and encrypting the output.  When a user executes their app, the console output is encrypted and checked against the encrypted success value.  Eventually, I'd like to replace this and probably explore keeping the validations within the present package.
* *Fun*: Building this app was incredibly rewarding.  And doing it in 48hrs was quite a challenge. The app idea was inspired by the JavaScript game [Untrusted] (http://alexnisnevich.github.io/untrusted/).  I thought Untrusted did such a good job with the sci-fi/arcade theme that I wanted to explore a more visual/real-world storyline.  After getting the access and animations setup, making the Golang puzzles themselves was the best part.  I hope Operation Go brings as much knowledge and enjoyment to others as it did for me.

## Team
1. Andy Brewer (Full-stack developer)

# Credits
* Soundtrack: Winter Rain by [hunkE] (https://soundcloud.com/hunke/winter-rain-instrumental)
