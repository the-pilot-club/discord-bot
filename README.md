<p align="center"><a href="https://thepilotclub.org" target="_blank"><img src="https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png" width="300" alt="The Pilot Club Logo"/></a></p>

# TPC Discord Bot v3

This is v3 of The Pilot Club's discord bot built with discordgo.

The configuration for this bot is hosted within another repository (will be released later). Development config map has been provided for developer convenience.  


## Development Requirements

- Go (https://go.dev/doc/install)


## Installation

To install this project, create an empty directory and run the following command:

```
git clone https://github.com/the-pilot-club/discord-bot.git
```

After cloning the repo, you will need to create an ``env`` file. To use the example env, run the following command:

```
cp example.env .env
```

You will need to create a Discord Application and Bot account to run this bot. For more information, visit the Discord documentation here: https://discord.com/developers/docs/quick-start/getting-started

After you create a bot account, obtain a bot token and paste it into your ``.env`` file under the ``BOT_TOKEN`` variable.

All the other values within the env file can remain the same as they are setup for the development environment.

Once that is complete, run the following command to start the bot:

```
go run cmd/main.go
```

Your bot is now running! 

Should any developer wish to use the TPC Development Discord Server to test their bot, please reach out to any of the developers of TPC and they can invite you.

Happy Coding
