package telegram

import (
	"log"
	"math/rand"
	"net/url"
	"remindbot/lib/e"
	"remindbot/storage"
	"strconv"
	"strings"
	"time"
)

//messages
const (
	msgHelp = `Hi, I'm your daily reminder. I will send you a link per day, so you do not forget to actually read the articles you save. You can either send me a page to save or pick a random one already saved.
	
	I can also provide you with a list of all saved links you currently have. 
	
	Here is a list of all commands: 
	/rnd gets you a random link 📨
	/all gets you all links you have got 📃

	p.s keep in mind that a link is instantly deleted when you retreive it 
	
	Send me a valid page url(that has http:// or https:// in it) to get started if you have not yet`

	msgHello = "Greetings 👋 \n \n" + msgHelp

	msgUnknownCommand = "❓ Ops, I dont know such command. Try this /help"

	msgNothingSaved  = "❌ Ops, seems like you have no pages saved. Send me a valid url to get started or read /help"
	msgSaved         = "Your link has been saved 👍"
	msgAlreadyExists = "🤔 Seems like you have already saved the same page. Type /all to see pages you currently have"
)

//commands
const (
	RndCmd   = "/rnd"
	AllCmd   = "/all"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	log.Printf("got new command '%s' from '%s'", text, username)
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case AllCmd:
		return p.sendAll(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand, "false")
	}
}

func (p *Processor) savePage(chatID int, pURL string, username string) error {
	page := &storage.Page{
		URL:       pURL,
		UserName:  username,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	IfExists, err := p.storage.IfExists(page)
	if err != nil {
		return e.Wrap("Failes to check for existence", err)
	}
	if IfExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists, "false")
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("failed to save a message", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved, "false"); err != nil {
		return e.Wrap("failed to save a message", err)
	}

	go remind(p, p.storage, username, chatID)

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil {
		return p.tg.SendMessage(chatID, msgNothingSaved, "false")
	}

	if err := p.tg.SendMessage(chatID, page.URL, "false"); err != nil {
		return e.Wrap("failed to send message", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendAll(chatID int, username string) error {
	pages, err := p.storage.PickAll(username)
	if err != nil {
		return p.tg.SendMessage(chatID, msgNothingSaved, "true")
	}

	URLs := ""
	for i, p := range pages {
		URLs += strconv.Itoa(i+1) + ") " + p.URL + "\n" + "📅 Saved at: " + p.CreatedAt + "\n\n"
	}

	if err := p.tg.SendMessage(chatID, URLs, "true"); err != nil {
		return e.Wrap("failed to get a list", err)
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp, "false")
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello, "false")
}

func isAddCmd(t string) bool {
	return isUrl(t)
}

func isUrl(t string) bool {
	u, err := url.Parse(t)

	return err == nil && u.Host != ""
}

//reminder
func remind(p *Processor, s storage.Storage, username string, ChatID int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn((24-5)+1) + 5

	time.Sleep(time.Duration(n) * time.Hour)
	page, _ := s.PickRandom(username)
	resp := "📬 Look what I brought you today: \n\n" + page.URL

	p.tg.SendMessage(ChatID, resp, "false")
	p.storage.Remove(page)
}
