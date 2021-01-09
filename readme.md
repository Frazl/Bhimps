# Bhimp(s)

Bhimp is a python discord bot developed using the discord.py API wrapper. 
It has some basic functionality and was mostly hacked together in a day.
The code lives on github so my friends and I can extend the bot easily to add `funny stuff`. 

## Features

### User Scores 

The bot tracks users scores. The users score is determined by by certain reactions other users can attach to messages sent in discord.
- `:plus_two:` --> Increases user scores by 2
- `:plus_one:` --> Increases user scores by 1
- `:minus_two:` --> Decreases user scores by 2
- `:minus_one:` --> Decreases user scores by 1
---
- The bot accounts for users adding or removing these reactions. So if a user attaches a `:minus_two:` but then removes it, the bot accounts for this and increments the affected user by two upon removal. 
- The bot also disallows users from incrementing or decrementing their own messages.
- The bot also has special functionality when it gets `:minus_two:`'d.

### Timers

In our discord server the word `pavlov` is forbidden. Anytime anyone asks whether people want to play pavlov the timer is increased by 48 hours.

Basic functionality within the bot is contained to set and update timers. 

The current implementation of the pavlov timer is quite poor, and often leads to false positives of the mentioning of the word. It's kinda funny though. 

### Commands
`!scores` displays the scores in a descending format, excluding any users with 0 score.
`!scores -` displays the scores in a ascending format, excluding any users with 0 score.
`!scores @mention_user` displays the score the mentioned user.

## Requirements
- python3
- pip3
- discord API bot permissions and a bot account

### Installation
- Install the requirements. `pip install -r requirements.txt`
- Fill in the entries to the .env file. 
- Setup the database through `python3 setup.py`
- Initialize the pavlov timer `python3 other.py`
- Begin running the proccess `python3 main.py`
