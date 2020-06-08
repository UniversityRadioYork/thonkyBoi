package main

import (
	"encoding/json"
	"github.com/UniversityRadioYork/myradio-go"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	newsOnJukebox    = true
	studioRedSource  = 1
	studioBlueSource = 2
	jukeboxSource    = 3
	obSource         = 4
	wsSource         = 5
	offAirSource     = 8
	logFile          = "thonkyBoi.log"
)

type webStudioData struct {
	Payload struct {
		Connections []struct {
			Timeslotid int `json:"timeslotid"`
		} `json:"connections"`
	} `json:"payload"`
}

func checkOB(timeslotID uint64) bool {
	return false
}

func checkWS(timeslotID uint64) bool {
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

	for _, val := range wsData.Payload.Connections {
		if val.Timeslotid == int(timeslotID) {
			return true
		}
	}
	return false
}

func checkTimeSoon(t time.Time) bool {
	return t.Add(time.Duration(-59) * time.Minute).Before(time.Now())
}

func main() {

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	log.Println("Software Startup for Upcoming Transition")

	log.Println("Starting API Session and Getting Data")

	session, err := myradio.NewSession("*****") // Timelord Key
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

	log.Println("Started API Session and Got Data")

	currentSel := selInfo.Studio

	if currentSel == offAirSource {
		// Off Air
		log.Println("Currently Off-Air - No SEL Commands to Issue\n")
		return
	}

	if selInfo.Lock != 0 {
		// Selector Locked
		log.Println("Selector Locked - Can't Issue SEL Commands\n")
		return
	}

	if !checkTimeSoon(timeslotInfo.Current.EndTime.Local()) && timeslotInfo.Current.Id != 0 {
		// Multi-Hour Show, let them do news (can be WS)
		log.Println("Multi-Hour Show Continuation - No SEL Commands to Issue\n")
		return
	}

	var commands [3]int

	/*
	   At 59:45, leave alone, unless jukebox either does a news, or going into a show
	*/

	log.Println("Starting Decisioning Process")

	if checkOB(timeslotInfo.Current.Id) {
		commands[0] = wsSource
	} else if currentSel == jukeboxSource {
		if checkTimeSoon(timeslotInfo.Next.StartTime.Local()) || newsOnJukebox {
			commands[0] = wsSource
		}
	}

	/*
		On the hour, leave alone, unless show going into jukebox or WS
	*/

	if currentSel == studioRedSource || currentSel == studioBlueSource {
		if timeslotInfo.Next.Id != 0 && !checkOB(timeslotInfo.Next.Id) && !checkWS(timeslotInfo.Next.Id) {
			commands[1] = 0
		} else {
			commands[1] = wsSource
		}
	}

	/*
		At 02:02, only change if needed (into/out of jukebox)
	*/

	var needToCheck bool

	if timeslotInfo.Next.Id == 0 || (timeslotInfo.Current.Id == 0 && !checkTimeSoon(timeslotInfo.Next.StartTime.Local()) && newsOnJukebox) { // Jukebox Next (either timeslot or hour)
		commands[2] = jukeboxSource
	} else if checkOB(timeslotInfo.Next.Id) { //OB Next
		commands[2] = obSource
	} else if !checkWS(timeslotInfo.Next.Id) && !(timeslotInfo.Current.Id == 0 && !checkTimeSoon(timeslotInfo.Next.StartTime.Local())) { // Not WS and Not Jukebox continuation
		needToCheck = true
	}

	log.Println("Finished Decisioning Process\n")

	log.Println("Upcoming SEL Commands")
	t := [3]string{"59:45", "00:00", "02:02"}
	c := [6]string{"No Action", "", "", "SEL 3 (Jukebox)", "SEL 4 (OB-Line)", "SEL 5 (WebStudio/AutoNews)"}
	for k, v := range commands {
		log.Println(t[k], c[v])
	}
	log.Println("Studio Check", needToCheck, "\n")

	goTime := 29 - time.Now().Second()

	if goTime < 0 {
		log.Println("SEL Decision Unclear at -31 Seconds - Possible Failure - Will Exit and Not Issue SEL Commands\n")
		return
	}

	time.Sleep(time.Duration(goTime) * time.Second)

	log.Println("Decisioning is OK at -31 seconds for SEL Commands")

	time.Sleep(16 * time.Second)
	log.Println("Executing 59:45")
	if commands[0] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[0]))
		exec.Command("sel", strconv.Itoa(commands[0]))
	}

	time.Sleep(15 * time.Second)
	log.Println("Executing 00:00")
	if commands[1] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[1]))
		exec.Command("sel", strconv.Itoa(commands[1]))
	}

	time.Sleep(122 * time.Second)
	log.Println("Executing 02:02")
	if commands[2] != 0 {
		log.Printf("Exec: sel %s\n", strconv.Itoa(commands[2]))
		exec.Command("sel", strconv.Itoa(commands[2]))
	}

	if needToCheck {
		selInfo, err = session.GetSelectorInfo()
		if err != nil {
			log.Println("Error Getting Post-Hour Selector Data - Will Assume Show Live and Exit")
			log.Fatal(err)
		}
		log.Println("Checking for Live Show")
		if selInfo.Studio == wsSource {
			log.Println("No Live Show - Exec: sel ", jukeboxSource)
			exec.Command("sel", strconv.Itoa(jukeboxSource))
		} else {
			log.Printf("Live Show - SEL %v\n", selInfo.Studio)
		}
	} else {
		log.Println("Skipping Studio Check - Not Required")
	}

	log.Println("System Shutdown\n")

}
