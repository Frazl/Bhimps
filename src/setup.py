import sqlite3 as sl
import discord
import os
from dotenv import load_dotenv

con = sl.connect('bhimps.db')

with con:
    con.execute("""
        CREATE TABLE USERSCORES (
            id INTEGER NOT NULL PRIMARY KEY,
            score INTEGER
        );
    """)
    con.execute("""
        CREATE TABLE TIMERS (
            id text NOT NULL PRIMARY KEY,
            time timestamp
        );
    """)

load_dotenv()
TOKEN = os.getenv('DISCORD_TOKEN')
GUILD = os.getenv('DISCORD_GUILD')

client = discord.Client()


@client.event
async def on_ready():
    print('We have logged in as {0.user}'.format(client))
    for guild in client.guilds:
        if guild.name == GUILD:
            data = []
            async for member in guild.fetch_members():
                data.append((member.id, 0))
            sql = 'INSERT INTO USERSCORES (id, score) values(?, ?)'
            with con:
                con.executemany(sql, data)
            exit(0)


client.run(TOKEN)
