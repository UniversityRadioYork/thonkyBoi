package main

import (
	"encoding/json"
	"github.com/UniversityRadioYork/myradio-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	studioRedSource  = 1
	studioBlueSource = 2
	jukeboxSource    = 3
	obSource         = 4
	wsSource         = 5
	offAirSource     = 8
	configFile       = "config.json"
)

type wsconnection struct {
	Timeslotid    int  `json:"timeslotid"`
	AutoNewsStart bool `json:"autoNewsBeginning"`
	AutoNewsEnd   bool `json:"autoNewsEnd"`
}

type wspayload struct {
	Connections []wsconnection `json:"connections"`
}

type webStudioData struct {
	Payload wspayload `json:"payload"`
}

type configAutoNews struct {
	TimeslotID    int  `json:"timeslotID"`
	AutoNewsStart bool `json:"autoNewsStart"`
	AutoNewsEnd   bool `json:"autoNewsEnd"`
}

type thonkyConfigBoi struct {
	APIKey           string           `json:"apiKey"`
	NewsOnJukebox    bool             `json:"newsOnJukebox"`
	OBShows          []int            `json:"obShows"`
	AutonewsRequests []configAutoNews `json:"autonewsRequests"`
}

func checkManualNews(timeslotID uint64, part string, wsData webStudioData, config thonkyConfigBoi) bool {
	var toReturn bool
	if part == "START" {
		for _, val := range wsData.Payload.Connections {
			if val.Timeslotid == int(timeslotID) {
				toReturn = !val.AutoNewsStart
			}
		}
		for _, val := range config.AutonewsRequests {
			if val.TimeslotID == int(timeslotID) {
				toReturn = !val.AutoNewsStart
			}
		}
	} else if part == "END" {
		for _, val := range wsData.Payload.Connections {
			if val.Timeslotid == int(timeslotID) {
				toReturn = !val.AutoNewsEnd
			}
		}
		for _, val := range config.AutonewsRequests {
			if val.TimeslotID == int(timeslotID) {
				toReturn = !val.AutoNewsEnd
			}
		}
	}
	return toReturn
}

// Is this time coming up soon
func checkTimeSoon(t time.Time) bool {
	return t.Add(time.Duration(-2) * time.Minute).Before(time.Now())
}

func checkOB(timeslotID uint64, config thonkyConfigBoi) bool {
	for _, val := range config.OBShows {
		if val == int(timeslotID) {
			return true
		}
	}
	return false
}

// Is a timeslotID registered for WS
func checkWS(timeslotID uint64, wsData webStudioData) bool {
	for _, val := range wsData.Payload.Connections {
		if val.Timeslotid == int(timeslotID) {
			return true
		}
	}
	return false
}

// Decisioning is where the decisions get made and the core logic is
func Decisioning(timeslotInfo *myradio.CurrentAndNext, wsData webStudioData, currentSel int, config thonkyConfigBoi) ([3]int, bool) {
	// This stuff below has nice names...that's all I have to say

	var jukeboxNext bool

	jukeboxNext = (checkTimeSoon(timeslotInfo.Next.StartTime.Local()) && timeslotInfo.Next.Id == 0) ||
		(!checkTimeSoon(timeslotInfo.Current.EndTime.Local()) && timeslotInfo.Current.Id == 0)

	var obNext bool

	obNext = (checkTimeSoon(timeslotInfo.Next.StartTime.Local()) && checkOB(timeslotInfo.Next.Id, config)) ||
		(!checkTimeSoon(timeslotInfo.Current.EndTime.Local()) && checkOB(timeslotInfo.Current.Id, config))

	var wsNext bool

	wsNext = (checkTimeSoon(timeslotInfo.Next.StartTime.Local()) && checkWS(timeslotInfo.Next.Id, wsData)) ||
		(!checkTimeSoon(timeslotInfo.Current.EndTime.Local()) && checkWS(timeslotInfo.Current.Id, wsData))

	var manualNews [2]bool

	// Middle really isn't a thing we need to worry about, because of catch-all above, except jukebox

	if currentSel == jukeboxSource && jukeboxNext {
		manualNews = [2]bool{!config.NewsOnJukebox, !config.NewsOnJukebox}
	} else if (currentSel == studioRedSource || currentSel == studioBlueSource) && (!jukeboxNext && !obNext && !wsNext) {
		manualNews = [2]bool{true, true}
	} else {
		if checkManualNews(timeslotInfo.Current.Id, "END", wsData, config) {
			manualNews[0] = true
		}

		if checkManualNews(timeslotInfo.Next.Id, "START", wsData, config) {
			manualNews[1] = true
		}
	}

	var commands [3]int

	/*
		59:45 Transition
	*/

	if (currentSel == jukeboxSource || currentSel == obSource) && (manualNews == [2]bool{false, false}) {
		commands[0] = wsSource
	}

	/*
		00:00 Transition
	*/

	if wsNext {
		if manualNews != [2]bool{true, false} {
			commands[1] = wsSource
		}

	} else if obNext {
		if manualNews == [2]bool{false, false} {
			commands[1] = wsSource
		} else if manualNews[1] {
			commands[1] = obSource
		}
	} else if jukeboxNext {
		if !manualNews[0] {
			commands[1] = wsSource
		}
	}

	/*
		02:02 Transition and Studio Check
	*/

	var needToCheck bool

	if jukeboxNext {
		commands[2] = jukeboxSource
	} else if obNext {
		commands[2] = obSource
	} else if wsNext {
		commands[2] = wsSource
	} else {
		needToCheck = true
	}

	return commands, needToCheck
}

func main() {

	/*
		Get Config
	*/

	configFile, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var config thonkyConfigBoi
	json.Unmarshal(byteValue, &config)

	/*
		Start Logging
	*/

	log.Println("Software Startup for Upcoming Transition")
	log.Println("Starting API Session and Getting Data")

	/*
		API Calling Stuff
	*/

	session, err := myradio.NewSession(config.APIKey)
	if err != nil {
		log.Println("Error Starting API Session - Will Exit and Not Issue SEL Commands")
		log.Fatal(err)
	}

	selInfo, err := session.GetSelectorInfo()
	if err != nil {
		log.Println("Error Getting Selector Data - Will Exit and Not Issue SEL Commands")
		log.Fatal(err)
	}

	timeslotInfo, err := session.GetCurrentAndNext()
	if err != nil {
		log.Println("Error Getting Timeslot Data - Will Exit and Not Issue SEL Commands")
		log.Fatal(err)
	}

	var wsData webStudioData
	res, err := http.Get("https://ury.org.uk/webstudio/api/v1/status")
	if err != nil {
		log.Println("Error Requesting WebStudio API - Will Exit and Not Issue SEL Commands")
		log.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&wsData)
	if err != nil {
		log.Println("Error Decoding WebStudio API - Will Exit and Not Issue SEL Commands")
		log.Fatal(err)
	}

	currentSel := selInfo.Studio

	log.Println("Started API Session and Got Data")

	/*
		Catch-All Statements
	*/

	if currentSel == offAirSource {
		// Off Air
		log.Printf("Currently Off-Air - No SEL Commands to Issue\n\n")
		return
	}

	if selInfo.Lock != 0 {
		// Selector Locked
		log.Printf("Selector Locked - Can't Issue SEL Commands\n\n")
		return
	}

	if !checkTimeSoon(timeslotInfo.Current.EndTime.Local()) && timeslotInfo.Current.Id != 0 {
		// Multi-Hour Show, let them do news (can be WS)
		log.Printf("Multi-Hour Show Continuation - No SEL Commands to Issue\n\n")
		return
	}

	/*
		Starting to Create Command Sequence
	*/

	log.Println("Starting Decisioning Process")

	commands, needToCheck := Decisioning(timeslotInfo, wsData, currentSel, config)

	log.Printf("Finished Decisioning Process\n\n")

	/*
		Logging the Proposed Plan, and checking its been decided in time
	*/

	log.Println("Upcoming SEL Commands")
	t := [3]string{"59:45", "00:00", "02:02"}
	c := [6]string{"No Action", "", "", "SEL 3 (Jukebox)", "SEL 4 (OB-Line)", "SEL 5 (WebStudio/AutoNews)"}
	for k, v := range commands {
		log.Println(t[k], c[v])
	}
	log.Printf("Studio Check %v \n\n", needToCheck)

	goTime := 29 - time.Now().Second()

	if goTime < 0 {
		log.Printf("SEL Decision Unclear at -31 Seconds - Possible Failure - Will Exit and Not Issue SEL Commands\n\n")
		return
	}

	time.Sleep(time.Duration(goTime) * time.Second)

	log.Println("Decisioning is OK at -31 seconds for SEL Commands")

	/*
		Run Commands
	*/

	time.Sleep(16 * time.Second)
	log.Println("Executing 59:45")
	if commands[0] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[0]))
		exec.Command("sel", strconv.Itoa(commands[0])).Run()
	}

	time.Sleep(15 * time.Second)
	log.Println("Executing 00:00")
	if commands[1] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[1]))
		exec.Command("sel", strconv.Itoa(commands[1])).Run()
	}

	time.Sleep(122 * time.Second)
	log.Println("Executing 02:02")
	if commands[2] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[2]))
		exec.Command("sel", strconv.Itoa(commands[2])).Run()
	}

	/*
		Executes the Studio Check
	*/

	if needToCheck {
		selInfo, err = session.GetSelectorInfo()
		if err != nil {
			log.Println("Error Getting Post-Hour Selector Data - Will Assume Show Live and Exit")
			log.Fatal(err)
		}
		log.Println("Checking for Live Show")
		if selInfo.Studio == wsSource {
			log.Println("No Live Show - Exec: sel ", jukeboxSource)
			exec.Command("sel", strconv.Itoa(jukeboxSource)).Run()
		} else {
			log.Printf("Live Show - SEL %v\n", selInfo.Studio)
		}
	} else {
		log.Println("Skipping Studio Check - Not Required")
	}

	/*
		End
	*/

	log.Printf("System Shutdown\n\n")

}
