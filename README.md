[![Build Status](https://travis-ci.org/Coding-Web-Community/CodingBump.svg?branch=main)](https://travis-ci.org/Coding-Web-Community/CodingBump)
# CodingBump
Disboard clone (for now)

# Architecture

![](https://media.discordapp.net/attachments/760959186745163806/761885595857190972/unknown.png)
## API
- Keeps track of server list, bump times, etc

## Discord Bot
- Discord bot POSTS guild ID's to the API when bumped
- other things about the Discord bot

## Website
- Website GETS recently bumped guilds
- Things about the website

# Contributing

Please make a new branch `{name}/development/{thing you're working on}` and only PR into the `development` branch!

IE: `mathis/development/api` -> PR -> `development`
