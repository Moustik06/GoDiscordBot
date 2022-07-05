package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	
	UserID = make(map[string]int)
)

func ConnectToDiscord(){
	discord,err := discordgo.New("Bot " + Token)
	if err != nil{
		fmt.Println("Erreur pendant la création de session")
		return
	}
	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentGuildMessages

	err = discord.Open()
	if err != nil{
		fmt.Println("Erreur de connexion")
		return
	}
	fmt.Println("Lancement du bot")
	discord.UpdateGameStatus(0, "Use ! for commands")

	sc := make(chan os.Signal,1)
	signal.Notify(sc,syscall.SIGINT,syscall.SIGTERM,os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate){
	if message.Author.ID == session.State.User.ID{
		return
	}

	UserID[message.Author.ID] = UserID[message.Author.ID] + 1

	if message.Content == "test"{
		session.ChannelMessageSend(message.ChannelID,"Accepté")

	}

	if message.Content == "score"{
		score := UserID[message.Author.ID]
		session.ChannelMessageSend(message.ChannelID,"Votre score est de : " + strconv.Itoa(score))
	}

	if message.Content == "leaderboard"{
		trieClassement(UserID,session,message)
	}

	if strings.Split(message.Content, " ")[0] == "!fetch" {
		fetchInfoFromUser(session,message,strings.Split(message.Content, " ")[1])
	}
}

func trieClassement(dico map[string]int,s *discordgo.Session,m *discordgo.MessageCreate){
	keys := make([]string, 0, len(dico))
  
    for key := range dico {
        keys = append(keys, key)
    }
  
    sort.SliceStable(keys, func(i, j int) bool{
        return dico[keys[i]] > dico[keys[j]]
    })
  
    for _, k := range keys{
        fmt.Println(k, dico[k])
		User, err := s.User(k)
		if err != nil{
			fmt.Println("Erreur de récupération")
			return
		}
		s.ChannelMessageSend(m.ChannelID,"Le score de "+ User.Username +" est de " + strconv.Itoa(dico[k]))
    }
}

func fetchInfoFromUser(s *discordgo.Session, m *discordgo.MessageCreate,user string){
	User,err := s.User(user)
	if err != nil{
		fmt.Println("Erreur de récupération")
		return
	}

	s.ChannelMessageSend(m.ChannelID,"Username : " + User.Username + "\n" + "Discriminator : " + User.Discriminator + "\n" + User.AvatarURL("1024") + "\n" +User.BannerURL("1024"))

}