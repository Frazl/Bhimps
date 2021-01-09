import os
import datetime

import discord
import db
import re
from dotenv import load_dotenv
import sqlite3 as sl
from difflib import SequenceMatcher
from random import choice


load_dotenv()
TOKEN = os.getenv('DISCORD_TOKEN')
GUILD = os.getenv('DISCORD_GUILD')

con = sl.connect('bhimps.db', detect_types=sl.PARSE_DECLTYPES)


client = discord.Client()
global bot_channel
bot_channel = None
FUNNY_STUFF = ['MY WIFE', 'THE M1', 'REURH', 'WOOF WOOF', 'MY HUNNY BUNNY UWU']


@client.event
async def on_ready():
    print('We have logged in as {0.user}'.format(client))
    for channel in client.get_all_channels():
        if channel.name == "bhimp-zone":
            global bot_channel
            bot_channel = channel


async def get_message_author(channel_id, message_id):
    uid = None
    for channel in client.get_all_channels():
        if channel.id == channel_id:
            message = await channel.fetch_message(message_id)
            uid = message.author.id
            break
    return uid


async def handle_score_reaction(payload, amount, is_removal=False):
    uid = await get_message_author(payload.channel_id, payload.message_id)
    if(payload.user_id == uid):
        return
    message_author = await client.fetch_user(uid)
    reaction_author = await client.fetch_user(payload.user_id)
    if(message_author.id == client.user.id and amount < 0):
        message_author = reaction_author
        reaction_author = client.user
        message = "What do you take me for %s? Do you know %s?" % (
            message_author.display_name,
            choice(FUNNY_STUFF)
            )
        embed = discord.Embed(title="Bou Bart Ba Bhimp", description=message)
        await bot_channel.send(embed=embed)
    user_score = db.get_user_score(con, message_author.id)
    user_score += amount
    db.update_user_score(con, message_author.id, user_score)
    if is_removal:
        amount *= -1
        color = 0x0e6b0e
        title = "Score Increased"
        if amount > 0:
            color = 0xe51937
            amount = "+%s" % amount
            title = "Score Decreased"
        description = "%s removed his %s for %s \n %s now has %s score" % (
            reaction_author.display_name,
            amount,
            message_author.display_name,
            message_author.display_name,
            user_score
            )
        embed = discord.Embed(
            title=title,
            description=description,
            color=color
            )
        await bot_channel.send(embed=embed)
    else:
        color = 0xe51937
        title = "Score Decreased"
        if amount > 0:
            amount = "+%s" % amount
            color = 0x0e6b0e
            title = "Score Increased"
        description = "%s %s'd %s \n %s now has %s score" % (
            reaction_author.display_name,
            amount,
            message_author.display_name,
            message_author.display_name,
            user_score
            )
        embed = discord.Embed(
            title=title,
            description=description,
            color=color)
        await bot_channel.send(embed=embed)


async def build_scores_embed(scores):
    embed = discord.Embed(title="Scores")
    for score in scores:
        user = await client.fetch_user(score[0])
        embed.add_field(name=user.display_name, value=score[1], inline=True)
    return embed


async def handle_help_message(message):
    output = '!scores \n'
    output += '!scores - \n'
    output += '!score @mention_user \n'
    await message.channel.send(output)


async def handle_scores_message(message):
    scores = db.get_scores(con, 10)
    embed = await build_scores_embed(scores)
    await message.channel.send(embed=embed)


async def handle_scoresminus_message(message):
    scores = db.get_scores(con, 10, False)
    embed = await build_scores_embed(scores)
    await message.channel.send(embed=embed)


async def handle_score_message(message):
    user = message.mentions[0]
    score = db.get_user_score(con, user.id)
    description = '%s has a score of %s' % (user.display_name, score)
    embed = discord.Embed(title="Score", description=description)
    await message.channel.send(embed=embed)


async def handle_other_messages(message):
    for word in message.content.lower().strip().split():
        word = re.sub('[^0-9a-zA-Z]+', '', word)
        if SequenceMatcher(None, "pavlov", word).ratio() >= 0.5:
            warning_msg = "You said a naughty word... Increasing the timer..."
            await message.channel.send(warning_msg)
            pavlov_time = db.get_timer(con, "pavlov")
            pavlov_time += datetime.timedelta(hours=48)
            time = pavlov_time.strftime('%d/%m/%y %H:%M')
            db.set_timer(con, "pavlov", pavlov_time)
            await message.channel.send("The timer is now at %s" % time)


@client.event
async def on_raw_reaction_add(payload):
    if payload.emoji.name == "minustwo":
        await handle_score_reaction(payload, -2)
    elif payload.emoji.name == "minusone":
        await handle_score_reaction(payload, -1)
    elif payload.emoji.name == "plustwo":
        await handle_score_reaction(payload, +2)
    elif payload.emoji.name == "plusone":
        await handle_score_reaction(payload, +1)


@client.event
async def on_raw_reaction_remove(payload):
    if payload.emoji.name == "minustwo":
        await handle_score_reaction(payload, +2, True)
    if payload.emoji.name == "minusone":
        await handle_score_reaction(payload, +1, True)
    elif payload.emoji.name == "plustwo":
        await handle_score_reaction(payload, -2, True)
    elif payload.emoji.name == "plusone":
        await handle_score_reaction(payload, -1, True)


@client.event
async def on_message(message):
    if message.author == client.user:
        return
    if message.content.lower().startswith("!help"):
        await handle_help_message(message)
    if message.content.lower().startswith("!scores"):
        await handle_scores_message(message)
    if message.content.lower().startswith("!scores -"):
        await handle_scoresminus_message(message)
    contains_mention = len(message.mentions) > 0
    if message.content.lower().startswith("!score") and contains_mention:
        await handle_score_message(message)
    else:
        await handle_other_messages(message)

client.run(TOKEN)
