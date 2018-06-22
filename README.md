#Project Clay
### Personal Twitter 'Activity' Feature Rebuilt in Go

For awhile, Twitter's mobile app had an "Activity" section under the Discover tab. Here you could find what people you were following were favoriting and who they started following. It was a great way to discover what the people you followed were interested in and find new content. When Twitter removed this feature, I was disappointed. 

Project Clay is my personal server side service that recreates this removed feature by collecting and deriving this same data and more (unfollows, unfavourites). Given your user ID, api keys and correct OAuth tokens, Project Clay will continuously scrape, store, and analyze Twitter data with a Mongo database. Complete data (user profile, user tweets, favourites, followers, and following) is collected for all users up to two degrees away and just user profiles for the third degree. 

This project was developed from January through May 2015 and was my first project with Go and MongoDB so don't expect the most idiomatic go code. **This repo is a sanitized snapshot of the project**. My API keys and sensitive data have been removed. The front-end for this project is Project Embla, a webpage written with PHP to present Project Clay's data, but PHP is so ugly I would feel bad if I posted it.

Project Clay was a fun, useful experiment, but is no longer under active development.


