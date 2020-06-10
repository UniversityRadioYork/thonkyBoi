package main

import (
	"github.com/UniversityRadioYork/myradio-go"
	"testing"
	"time"
)

func TestStudiosStudiosNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 0, 0}
	expectedStudioCheck := true
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{{TimeslotID: 1, AutoNewsEnd: true}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosJukeboxAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 5, 3}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 0, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosJukeboxNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 0, 3}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 0, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosOBAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 5, 4}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: true}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosOBNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 0, 4}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosWSAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 5, 5}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2}}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: true}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosWSNoAutonewsFirst(t *testing.T) {
	expectedCommands := [3]int{0, 0, 5}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2, AutoNewsStart: true}}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestStudiosWSNoAutonewsSecond(t *testing.T) {
	expectedCommands := [3]int{0, 5, 5}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 1, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2, AutoNewsStart: false}}}}
	var currentSel int = studioRedSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{}, AutonewsRequests: []configAutoNews{}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxStudiosAutonews(t *testing.T) {
	expectedCommands := [3]int{5, 0, 0}
	expectedStudioCheck := true
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 1, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false, AutoNewsStart: true}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxStudiosNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 0, 0}
	expectedStudioCheck := true
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 1, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsStart: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxJukeboxAutonews(t *testing.T) {
	expectedCommands := [3]int{5, 5, 3}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 0, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: true, OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxJukeboxNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 0, 3}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 0, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 2}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: false, OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxOBAutonews(t *testing.T) {
	expectedCommands := [3]int{5, 5, 4}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: true, OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsEnd: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxOBNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 4, 4}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 2, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: true, OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 2, AutoNewsStart: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxWSAutonews(t *testing.T) {
	expectedCommands := [3]int{5, 5, 5}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 1, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 1, AutoNewsStart: true}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: true, OBShows: []int{2}, AutonewsRequests: []configAutoNews{}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}

func TestJukeboxWSNoAutonews(t *testing.T) {
	expectedCommands := [3]int{0, 5, 5}
	expectedStudioCheck := false
	var timeslotInfo myradio.CurrentAndNext = myradio.CurrentAndNext{
		Current: myradio.Show{Id: 0, EndTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
		Next:    myradio.Show{Id: 1, StartTime: myradio.Time{Time: time.Now().Add(time.Minute)}},
	}
	var wsData webStudioData = webStudioData{Payload: wspayload{Connections: []wsconnection{{Timeslotid: 1}}}}
	var currentSel int = jukeboxSource
	var config thonkyConfigBoi = thonkyConfigBoi{NewsOnJukebox: true, OBShows: []int{2}, AutonewsRequests: []configAutoNews{configAutoNews{TimeslotID: 1, AutoNewsStart: false}}}
	actualCommands, actualStudioCheck := Decisioning(&timeslotInfo, wsData, currentSel, config)

	if (expectedCommands != actualCommands) || (expectedStudioCheck != actualStudioCheck) {
		t.Errorf("Test Failed: Expected: %v, %v, Got: %v, %v", expectedCommands, expectedStudioCheck, actualCommands, actualStudioCheck)
	}
}
