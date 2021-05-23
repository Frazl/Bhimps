package main

import (
	"bhimp/datamanager"
	"database/sql"
	"fmt"
	"log"
	"mime"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var db *sql.DB
var GUILD_ID, BOT_CHANNEL_ID string

func main() {
	log.Println("Beginning database setup")
	db = datamanager.Setup()
	dg, err := discordgo.New("Bot " + "BOT_TOKEN")
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageReactionRemove)
	GUILD_ID, BOT_CHANNEL_ID = getBotChannel(dg)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

// Checks all messages
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!scores" {
		userScores := datamanager.GetUserScores(db, 10, true)
		sendUserScoresEmbed(s, userScores)
	} else if m.Content == "!scores -" {
		userScores := datamanager.GetUserScores(db, 10, false)
		sendUserScoresEmbed(s, userScores)
	} else if m.Content == "!messages" {
		messageScores := datamanager.GetMessageScores(db, 10, true)
		sendMessageScoreEmbeds(s, messageScores)
	} else if m.Content == "!messages -" {
		messageScores := datamanager.GetMessageScores(db, 10, false)
		sendMessageScoreEmbeds(s, messageScores)
	} else if len(m.Mentions) > 0 && strings.HasPrefix(m.Content, "!score") {
		mentionee, _ := strconv.Atoi(m.Mentions[0].ID)
		score := datamanager.GetUserScore(db, mentionee)
		var us = datamanager.UserScore{mentionee, score}
		userScores := []datamanager.UserScore{us}
		sendUserScoresEmbed(s, userScores)
	}
}

// Handles the addition of message reactions
func messageReactionAdd(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(mra.MessageReaction.ChannelID, mra.MessageReaction.MessageID)

	// If message reactor same as author then don't do anything
	if msg.Author.ID == mra.MessageReaction.UserID {
		return
	}

	msgID, _ := strconv.Atoi(msg.ID)
	chID, _ := strconv.Atoi(mra.MessageReaction.ChannelID)
	if err != nil {
		log.Println("Error handling reaction to message", err)
	}
	emojiName := mra.MessageReaction.Emoji.Name
	msgAuthorId, _ := strconv.Atoi(msg.Author.ID)

	botUser, _ := s.User("@me")
	if msg.Author.ID == botUser.ID && (emojiName == "minustwo" || emojiName == "minusone") {
		msgAuthorId, _ = strconv.Atoi(mra.MessageReaction.UserID)
		s.ChannelMessageSend(BOT_CHANNEL_ID, "Stop hitting yourself...")
	}

	responseEmbed := new(discordgo.MessageEmbed)
	if emojiName == "plustwo" {
		datamanager.ModifyUserScore(db, msgAuthorId, 2)
		datamanager.ModifyMessageScore(db, chID, msgID, 2)
		responseEmbed.Color = 0x0e6b0e
	} else if emojiName == "plusone" {
		datamanager.ModifyUserScore(db, msgAuthorId, 1)
		datamanager.ModifyMessageScore(db, chID, msgID, 1)
		responseEmbed.Color = 0x0e6b0e
	} else if emojiName == "minustwo" {
		datamanager.ModifyUserScore(db, msgAuthorId, -2)
		datamanager.ModifyMessageScore(db, chID, msgID, -2)
		responseEmbed.Color = 0xe51937
	} else if emojiName == "minusone" {
		datamanager.ModifyUserScore(db, msgAuthorId, -1)
		datamanager.ModifyMessageScore(db, chID, msgID, -1)
		responseEmbed.Color = 0xe51937
	} else {
		return
	}
	responseEmbed.Title = "Reaction Add"
	// convert id back to string
	messageAuthorIDString := strconv.Itoa(msgAuthorId)
	reactorGuildMember := getDisplayName(s, mra.MessageReaction.GuildID, mra.MessageReaction.UserID)
	messageAuthorGuildMember := getDisplayName(s, mra.MessageReaction.GuildID, messageAuthorIDString)
	responseEmbed.Description = reactorGuildMember + " added a " + mra.MessageReaction.Emoji.MessageFormat() + " to " + messageAuthorGuildMember + "'s message!"
	responseEmbed.Description += "\n"
	authorScore := strconv.Itoa(datamanager.GetUserScore(db, msgAuthorId))
	responseEmbed.Description += messageAuthorGuildMember + " now has a score of " + authorScore
	s.ChannelMessageSendEmbed(BOT_CHANNEL_ID, responseEmbed)
}

// Handles the removal of message reactions
func messageReactionRemove(s *discordgo.Session, mrr *discordgo.MessageReactionRemove) {
	msg, err := s.ChannelMessage(mrr.MessageReaction.ChannelID, mrr.MessageReaction.MessageID)

	// If message reactor same as author then don't do anything
	if msg.Author.ID == mrr.MessageReaction.UserID {
		return
	}

	msgID, _ := strconv.Atoi(msg.ID)
	chID, _ := strconv.Atoi(mrr.MessageReaction.ChannelID)
	if err != nil {
		log.Println("Error handling reaction removal to message", err)
	}
	emojiName := mrr.MessageReaction.Emoji.Name
	msgAuthorId, _ := strconv.Atoi(msg.Author.ID)

	responseEmbed := new(discordgo.MessageEmbed)
	if emojiName == "plustwo" {
		datamanager.ModifyUserScore(db, msgAuthorId, -2)
		datamanager.ModifyMessageScore(db, chID, msgID, -2)
		responseEmbed.Color = 0xe51937
	} else if emojiName == "plusone" {
		datamanager.ModifyUserScore(db, msgAuthorId, -1)
		datamanager.ModifyMessageScore(db, chID, msgID, -1)
		responseEmbed.Color = 0xe51937
	} else if emojiName == "minustwo" {
		datamanager.ModifyUserScore(db, msgAuthorId, 2)
		datamanager.ModifyMessageScore(db, chID, msgID, 2)
		responseEmbed.Color = 0x0e6b0e
	} else if emojiName == "minusone" {
		datamanager.ModifyUserScore(db, msgAuthorId, 1)
		datamanager.ModifyMessageScore(db, chID, msgID, 1)
		responseEmbed.Color = 0x0e6b0e
	} else {
		return
	}
	responseEmbed.Title = "Reaction Removed"
	reactorGuildMember := getDisplayName(s, mrr.MessageReaction.GuildID, mrr.MessageReaction.UserID)
	messageAuthorGuildMember := getDisplayName(s, mrr.MessageReaction.GuildID, msg.Author.ID)
	responseEmbed.Description = reactorGuildMember + " removed a " + mrr.MessageReaction.Emoji.MessageFormat() + " to " + messageAuthorGuildMember + "'s message!"
	responseEmbed.Description += "\n"
	authorScore := strconv.Itoa(datamanager.GetUserScore(db, msgAuthorId))
	responseEmbed.Description += messageAuthorGuildMember + " now has a score of " + authorScore
	s.ChannelMessageSendEmbed(BOT_CHANNEL_ID, responseEmbed)
}

// Helper functions

// Attempts to get the nickname of a user of a given guild.
// Falls back to their username if none set.
func getDisplayName(s *discordgo.Session, guildID string, userID string) string {
	gm, err := s.GuildMember(guildID, userID)
	if err != nil {
		log.Println(err)
		log.Println(guildID, userID)
		log.Fatalln("Could not get display name for user in guild")
	}
	if gm.Nick != "" {
		return gm.Nick
	} else {
		return gm.User.Username
	}
}

// Formats and sends user scores as embeds.
func sendUserScoresEmbed(s *discordgo.Session, userscores []datamanager.UserScore) {
	responseEmbed := new(discordgo.MessageEmbed)
	responseEmbed.Title = "Scoreboard"
	responseEmbed.Fields = make([]*discordgo.MessageEmbedField, len(userscores))
	if len(userscores) > 0 {
		scoreboardLeaderID := strconv.Itoa(userscores[0].ID)
		scoreboardLeader, _ := s.User(scoreboardLeaderID)
		scoreboardLeaderImg := scoreboardLeader.AvatarURL("")
		thumbnail := new(discordgo.MessageEmbedThumbnail)
		thumbnail.URL = scoreboardLeaderImg
		responseEmbed.Thumbnail = thumbnail
	}
	for i, us := range userscores {
		mev := new(discordgo.MessageEmbedField)
		userID := strconv.Itoa(us.ID)
		mev.Name = getDisplayName(s, GUILD_ID, userID)
		mev.Value = strconv.Itoa(us.Score)
		mev.Inline = true
		responseEmbed.Fields[i] = mev
	}
	_, err := s.ChannelMessageSendEmbed(BOT_CHANNEL_ID, responseEmbed)
	if err != nil {
		log.Fatalln(err)
	}
}

// Formats and sends message scores as embeds.
func sendMessageScoreEmbeds(s *discordgo.Session, messagescores []datamanager.MessageScore) {
	for i, msg := range messagescores {
		responseEmbed := new(discordgo.MessageEmbed)
		responseEmbed.Title = "Message #" + strconv.Itoa(i+1) + " (" + strconv.Itoa(msg.Score) + ")"
		originalMessage, err := s.ChannelMessage(msg.ChannelID, msg.MessageID)
		if err != nil {
			log.Printf("Failed to get message cid: %s id: %s \n", msg.ChannelID, msg.MessageID)
			continue
		}

		// Add user's display name and content of their message
		authorDisplayName := getDisplayName(s, GUILD_ID, originalMessage.Author.ID)
		mea := discordgo.MessageEmbedAuthor{Name: authorDisplayName}
		responseEmbed.Author = &mea
		responseEmbed.Description = originalMessage.Content

		// Add user's image thumbnail
		authorUser, _ := s.User(originalMessage.Author.ID)
		authorAvatar := authorUser.AvatarURL("")
		thumbnail := new(discordgo.MessageEmbedThumbnail)
		thumbnail.URL = authorAvatar
		responseEmbed.Thumbnail = thumbnail

		// If there are image attachments then add them...
		if len(originalMessage.Attachments) > 0 {
			fileType := mime.TypeByExtension(filepath.Ext(originalMessage.Attachments[0].URL))
			if strings.HasPrefix(fileType, "image") {
				mei := discordgo.MessageEmbedImage{
					URL:      originalMessage.Attachments[0].URL,
					ProxyURL: originalMessage.Attachments[0].ProxyURL,
					Width:    originalMessage.Attachments[0].Width,
					Height:   originalMessage.Attachments[0].Height}
				responseEmbed.Image = &mei
				responseEmbed.Description += "\n\n contained image: " + originalMessage.Attachments[0].URL
			} else if strings.HasPrefix(fileType, "video") {
				// Discord does not let bots embed videos...
				responseEmbed.Description += "\n\n contained video: " + originalMessage.Attachments[0].URL
			}
		}
		_, err = s.ChannelMessageSendEmbed(BOT_CHANNEL_ID, responseEmbed)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

// Determines the correct channel to post bot messages and the correct guild.
// Sets global variables that are used by other functions.
// Relies on the environment variables:
// BHIMP_GUILD - The ID of the guild.
// BHIMP_CHANNEL - The ID of the channel within the guild to post messages.
func getBotChannel(s *discordgo.Session) (string, string) {
	ugs, _ := s.UserGuilds(20, "", "")
	for _, ug := range ugs {
		if ug.Name == os.Getenv("CREATORBOT_GUILD") {
			gid := ug.ID
			log.Println("Guild: " + ug.Name)
			channels, _ := s.GuildChannels(gid)
			for _, ch := range channels {
				if ch.Name == os.Getenv("CREATORBOT_CHANNEL") {
					log.Println("Bot Channel: " + ch.Name)
					return gid, ch.ID
				}
			}

		}
	}
	log.Fatalln("Couldn't determine the bot guild and channel")
	return "", ""
}
