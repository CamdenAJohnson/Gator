Gator is a RSSFeed CLI.

Gator uses postgresql to store users and feeds.

List of commands:
- register: Register a user. Usage: "register <\name>"
- Login: Login as the given user. Usage: "login <\name>"
- reset: Clears the postgresql database. Usage: "reset"
- users: Prints a list of all users within the postgresql database. Usage: "users"
- agg: Scraps feeds and stores them into the database on a given interval. Usage: "agg <\time>" 
- addfeed: Adds a feed to be scraped during agg. Usage: "addfeed <\name> <\url>"
- feeds: Prints a list of feeds stored within the database. Usage: "feeds"
- follow: Lets the current logged in user to follow a feed for scrapping. Usage: "Follow <\url>"
- following: Prints a list of the feeds the current logged in user is following. Usage: "following" 
- unfollow: Lets the current logged in user unfollow a feed. Usage: "unfollow <\url>"
- browse: Prints a list of scraped feeds. Usage: "browse <\limit>"