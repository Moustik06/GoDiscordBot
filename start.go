package main

import (
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
)

var (
	UserID 			  = make(map[string]int)
	c      				*cache.Cache
	defaultExpiration = time.Hour * 24
)

func cacheUser(){	
	c.Set("UserID", UserID, defaultExpiration)
	
	foo,found := c.Get("UserID")
	if found{
		log.Println(foo)
	}
}
func cacheExpired(user *map[string]int){
	foo,found := c.Get("UserID")
		if !found || foo == nil{
			log.Println("Vide du cache")
			c.Delete("UserID")
			UserID = make(map[string]int)
			c.Set("UserID", user, defaultExpiration)
		}
} 
func ConnectToDiscord() {

	cacheUser()
	
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

	defer discord.Close()
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	time.AfterFunc(time.Duration(30)*time.Second, func() {leaderboardJour(session,message)})
	
	UserID[message.Author.ID]++
	switch message.Content {
	case "test":
		//cacheExpired(&UserID)
		leaderboardJour(session,message)

	case "score":
		score := UserID[message.Author.ID]
		session.ChannelMessageSend(message.ChannelID, "Votre score est de : "+strconv.Itoa(score))

	case "leaderboard":
		trieClassement(UserID, session, message)
	}

	if strings.Split(message.Content, " ")[0] == "!fetch" {
		fetchInfoFromUser(session, message, strings.Split(message.Content, " ")[1])
	}
}

func trieClassement(dico map[string]int, s *discordgo.Session, m *discordgo.MessageCreate) string{
	var tmp string
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
		}
		if len(strings.Split(tmp, "$")) <=4{
			tmp = tmp + "$" + User.Username + " " + strconv.Itoa(dico[k])
		}
		//s.ChannelMessageSend(m.ChannelID, "Le score de "+User.Username+" est de "+strconv.Itoa(dico[k]))
	}
	
	return tmp
}

func fetchInfoFromUser(s *discordgo.Session, m *discordgo.MessageCreate, u string) {
	user, err := s.User(u)

	if err != nil {
		log.Panic("Erreur de récupération")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Username : "+user.Username+"\n"+"Discriminator : "+user.Discriminator+"\n"+user.AvatarURL("1024")+"\n"+user.BannerURL("1024"))
}

func leaderboardJour(s *discordgo.Session, m *discordgo.MessageCreate){
	tmp := trieClassement(UserID,s,m)
	test := strings.Split(tmp,"$")
	embed := new(discordgo.MessageEmbed)
	embed.Title = "Classement des 5 meilleurs du jour !"
	embed.Type = "rich"

	for i,j := range test {
		embed.Description += strconv.Itoa(i) + " - " + j + "\n"
	}
	s.ChannelMessageSendEmbed(m.ChannelID,embed)
	cacheExpired(&UserID)
	time.AfterFunc(time.Second*30, func() {leaderboardJour(s,m)})
}

