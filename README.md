# thonkyBoi

### URY Auto-Switcher/Auto-Selector (whatever you wanna call it)

## How it (should) Work

* At 59 minutes past the hour, decides what four actions to do

* These are three selector commands (or do nothing) and yes/no to a studio check

* Checks at -31 seconds it has stuff to do

* Does stuff at XX:59:45, YY:00:00 and YY:02:02 based on the three selector commands (can be do nothing)

* If neccesary, does a check to make sure WS isn't on air

## Instructions

### Development and Testing

* Copy `config.json.example` to `config.json`
* Test with `go test -v`

### Configuartion Usage

* OBShows is an array of timeslot IDs (integers) that will be considered to be a Source 4 OB by the software
* autoNewsRequests has JSON objects of timeslotID, autoNewsStart and autoNewsEnd to control autonews manually for non WS shows
* Updates to `config.json.example` should be reflected in the structs towards the start of `thonkyBoi.go`
* The config file path is a `const` at the start of `thonkyBoi.go`

## Code Layout (so you can find stuff)

* Imports, Constants, Structs
* checkAutonews
* checkTimeSoon
* checkOB
* checkWS
* Decisioning
* main
    * Start Logging
    * MyRadio and WebStudio API Calls
    * Catch-All Statements
    * Determining bools for next sources and autonewses and the 3 transitions and studioCheck (Decisioning)
    * Logging Plan, and -31 Seconds Readiness Check
    * Running the Commands

### API Thoughts

`CurrentAndNext` returns (from `myradio-go`) `myradio.Show` structs, rather than timeslots. This annoyed me when I found this out, because I had (seemingly) coded all of this wrong. But it turns out, the `Show`'s ID in C.a.N. isn't the Show ID, but is instead the timeslot ID. So, this code works, even if the API is questionable.


###### Michael Grace 2020