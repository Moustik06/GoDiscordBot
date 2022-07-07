package main

import (
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	//"time"
	"github.com/bwmarrin/discordgo"
	//"github.com/patrickmn/go-cache"
)

var (
	UserID = make(map[string]int)
)

func ConnectToDiscord() {
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Panic("Erreur pendant la création de session")
		return
	}
	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentGuildMessages

	err = discord.Open()
	if err != nil {
		log.Panic("Erreur de connexion")
	}
	log.Println("Lancement du bot")
	discord.UpdateGameStatus(0, "Use ! for commands")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	UserID[message.Author.ID] += 1

	if message.Content == "test" {
		session.ChannelMessageSend(message.ChannelID, "Accepté")

	}

	if message.Content == "score" {
		score := UserID[message.Author.ID]
		session.ChannelMessageSend(message.ChannelID, "Votre score est de : "+strconv.Itoa(score))
	}

	if message.Content == "leaderboard" {
		trieClassement(UserID, session, message)
	}

	if strings.Split(message.Content, " ")[0] == "!fetch" {
		fetchInfoFromUser(session, message, strings.Split(message.Content, " ")[1])
	}
}

func trieClassement(dico map[string]int, s *discordgo.Session, m *discordgo.MessageCreate) {
	keys := make([]string, 0, len(dico))
	for key := range dico {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return dico[keys[i]] > dico[keys[j]]
	})

	for _, k := range keys {
		log.Println(k, dico[k])
		User, err := s.User(k)
		if err != nil {
			log.Panic("Erreur de récupération")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Le score de "+User.Username+" est de "+strconv.Itoa(dico[k]))
	}
}

func fetchInfoFromUser(s *discordgo.Session, m *discordgo.MessageCreate, u string) {
	user, err := s.User(u)

	if err != nil {
		log.Panic("Erreur de récupération")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Username : "+user.Username+"\n"+"Discriminator : "+user.Discriminator+"\n"+user.AvatarURL("1024")+"\n"+user.BannerURL("1024"))

}

/*func leaderboardJour(s *discordgo.Session, m *discordgo.MessageCreate){
	// On regarde le nombre de msg en 24h


}*/
