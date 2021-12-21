package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/icza/gowut/gwu"
)

const (
	background  = "#000034"
	background1 = "#08275e"
	background2 = "#404f8c"
	foreground  = "#ffffff"
	width       = "600"
	list_height = "400"
	font_size   = "40"
)

type serverResponse struct {
}

type sessHandler struct{}

func connect(ip string) {

}

type Config struct {
	Port    int
	Address string
}

var (
	config Config
)

type Status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Favicon string
}

func style(c gwu.Comp) {
	c.Style().SetBackground(background)
	c.Style().SetColor(foreground)
}
func (h sessHandler) Created(s gwu.Session) {
	win := gwu.NewWindow("app", "Ping Utility")
	win.Style().SetBackground(background)
	win.Style().SetColor(foreground)
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.AddHeadHTML(`<body style = "background-color: #000034;">`)
	tabpanel := gwu.NewTabPanel()
	tabpanel.Style().SetWidth(width)
	tabpanel.SetHAlign(gwu.HACenter)

	//build input
	inputpanel := gwu.NewPanel()
	horizpanel := gwu.NewHorizontalPanel()
	inputpanel.SetHAlign(gwu.HACenter)
	inputpanel.Style().SetColor(foreground)

	//address part
	address_panel := gwu.NewPanel()
	style(address_panel)
	label := gwu.NewLabel("Enter server address:")
	address_panel.Add(label)
	inputfield := gwu.NewTextBox("")
	address_panel.Add(inputfield)
	horizpanel.Add(address_panel)

	//port part
	port_panel := gwu.NewPanel()
	style(port_panel)
	port_label := gwu.NewLabel("Enter server port:")
	port_panel.Add(port_label)
	portinputfield := gwu.NewTextBox("25565")
	port_panel.Add(portinputfield)
	horizpanel.Add(port_panel)

	inputpanel.Add(horizpanel)
	connectbutton := gwu.NewButton("Connect")
	resultpanel := gwu.NewPanel()
	connectbutton.AddEHandlerFunc(func(e gwu.Event) {
		text := inputfield.Text()
		inputfield.SetText("")
		port := portinputfield.Text()
		portinputfield.SetText("25565")

		address := text + ":" + port

		response, delay, err := bot.PingAndListTimeout(address, 5*time.Second)
		if err != nil {
			fmt.Println("Fail:", err)
			e.MarkDirty(inputfield, portinputfield)
			return
		}

		var decode Status
		err = json.Unmarshal(response, &decode)
		if err != nil {
			fmt.Println("Fail unmarshaling: ", err)
		}
		//fmt.Println(decode)
		fmt.Println("Delay: ", delay)

		//build result tab
		tabpanel.Remove(resultpanel)
		resultpanel = gwu.NewPanel()
		style(resultpanel)
		resultpanel.SetHAlign(gwu.HACenter)

		//favicon
		favicon := gwu.NewImage("Favicon", decode.Favicon)
		resultpanel.Add(favicon)

		//blank space
		resultpanel.Add(gwu.NewLabel(""))

		//motd
		resultpanel.Add(gwu.NewLabel(decode.Description.ClearString()))

		//version string
		resultpanel.Add(gwu.NewLabel("Version: " + chat.Text(decode.Version.Name).ClearString()))

		//protocol
		resultpanel.Add(gwu.NewLabel("Protocol: " + strconv.Itoa(decode.Version.Protocol)))

		//online
		resultpanel.Add(gwu.NewLabel("Online Players / Max Players: " + strconv.Itoa(decode.Players.Online) + "/" + strconv.Itoa(decode.Players.Max)))

		//blank space

		//playerlist
		var players []string
		playerlist := gwu.NewListBox(players)
		//playerlist.Style().SetWidth("50")
		uuids := ""
		for index, player := range decode.Players.Sample {
			index = index
			players = append(players, player.Name)
			uuids = uuids + "\n" + player.ID.String()
		}
		playerlist.SetValues(players)
		playerlist.SetToolTip(uuids)
		if len(players) > 0 {
			resultpanel.Add(playerlist)
		}

		tabpanel.AddString("Result", resultpanel)
		tabpanel.SetSelected(1)
		e.MarkDirty(tabpanel)

	}, gwu.ETypeClick)

	inputpanel.Add(connectbutton)

	//add stuff
	tabpanel.AddString("Input", inputpanel)
	title := gwu.NewLabel("Minecraft Ping Utility")
	subtitle := gwu.NewLabel("Utility for geting info about Minecraft servers")
	title.Style().SetFontSize("30")
	win.Add(title)
	win.Add(subtitle)
	win.Add(tabpanel)
	s.AddWin(win)
}
func (h sessHandler) Removed(s gwu.Session) {}

func main() {
	var config Config
	config.Address = "localhost"
	config.Port = 8085

	config_file, err := ioutil.ReadFile("config.json")
	fmt.Println("File read error: ", err)
	if err == nil {
		err2 := json.Unmarshal(config_file, &config)
		if err2 != nil {
			fmt.Println("Error while unmarshalling: ", err2)
		}
	}
	fmt.Println("Address: " + config.Address)
	fmt.Println("Port: " + strconv.Itoa(config.Port))

	fmt.Println("Starting")

	server := gwu.NewServer("", config.Address+":"+strconv.Itoa(config.Port))
	server.AddSessCreatorName("app", "Ping Utility")
	server.AddSHandler(sessHandler{})
	server.Start()
}
