# Retweet contest bot
Conduct retweet contests from the command line.  

## Features
* Track a tweet's retweets periodically.
* Choose one randomly to award a prize.
* Track multiple retweet contests

## Requirements
1. Twitter developer account for credentials
2. Go installation (or download binaries)

## Usage
1. Get your credentials from Twitter Developer and place them in credentials.json
2. Run `go run .` to start the application.
    * The options are described when you run the command
    * Fetch makes a copy of the latest 100 retweeters
    * Choose allows you to pick upto 100 winners for the contest
3. Figure out which tweet you want to track and copy it's id (Available in the url)

## ToDo
- [x] Use access token to get a tweet
- [x] Get retweeters of specific tweet
- [x] Update retweeters list to file
- [x] Choose random retweeter
- [x] Document Code