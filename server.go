package main

import (
	"flag"
	"fmt"
	"github.com/61c-teach/sp19-proj5-userlib"
	"net/http"
	"log"
	_ "strings"
	_ "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// FIXME This should be using the cache!
	// Note that we will be using userlib.ReadFile we provided to read files on the system.
	// The path to the file is given by r.URL.Path and will be the path to the string.
	// Make sure you properly sanitise it (more described in get file).
	/*** YOUR CODE HERE ***/
	// Reads the file from the disk
	response, err := userlib.ReadFile(workingDir, r.URL.Path[1:])
	if err != nil {
		// If we have an error from the read, will return the generic file error message and set the error code to follow that.
		http.Error(w, userlib.FILEERRORMSG, userlib.FILEERRORCODE)
		return
	}
	w.Header().Set(userlib.ContextType, userlib.GetContentType(r.URL.Path))
	/*** YOUR CODE HERE END ***/

	// This will automagically set the right content type for the
	// reply as well.
	// We need to set the correct header code for a success since we should only succeed at this point.
	w.WriteHeader(userlib.SUCCESSCODE) // Make sure you write the correct header code so that the tests do not fail!
	// Write the data which is given to us by the response.
	w.Write(response)
}

// This function will handle the requests to acquire the cache status.
func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the header of the request to a plain text format since we are just dumping information about the cache.
	w.Header().Set(userlib.ContextType, "text/plain; charset=utf-8")
	// Set the success code to the proper success code since the action should not fail.
	w.WriteHeader(userlib.SUCCESSCODE)
	// Get the cache status string from the getCacheStatus function.
	w.Write([]byte(getCacheStatus()))
}

// This function will handle the requests to clear/restart the cache.
func cacheClearHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the header of the request to a plain text format since we are just dumping information about the cache.
	w.Header().Set(userlib.ContextType, "text/plain; charset=utf-8")
	// Set the success code to the proper success code since the action should not fail.
	w.WriteHeader(userlib.SUCCESSCODE)
	// Get the cache status string from the getCacheStatus function.
	w.Write([]byte(CacheClear()))
}

// The structure used for responding to file requests.
// It contains the file contents (if there is any)
// or the error returned when accessing the file.
type fileResponse struct {
	filename string
	responseData []byte
	responseError error
	responseChan chan *fileResponse
}

// To request files from the cache, we send a message that 
// requests the file and provides a channel for the return
// information.
type fileRequest struct {
	filename string
	response chan *fileResponse
}

// DO NOT CHANGE THESE NAMES OR YOU WILL NOT PASS THE TESTS
// Port of the server to run on
var port int
// Capacity of the cache in Bytes
var capacity int
// Timeout for file reads in Seconds.
var timeout int
// The is the working directory of the server
var workingDir string

// The channel to pass file read requests to. This is how you will get a file from the cache.
var fileChan = make(chan *fileRequest )
// The channel to pass a request to get back the capacity info of the cache.
var cacheCapacityChan = make(chan chan string)
// The channel where a bool passed into it will cause the OperateCache function to be closed and all of the data to be cleared.
var cacheCloseChan = make(chan bool)

// A wrapper function that does the actual getting of the file
func getFile(filename string) (response *fileResponse) {
	// You need to add sanity checking here: The requested file
	// should be made relative (strip out leading "/" characters,
	// then have a "./" put on the start, and if there is ever the
	// string "/../", replace it with "/", the string "\/" should
	// be replaced with "/", and finally any instances of "//" (or
	// more) should be replaced by a single "/".
	// Hint: A replacement may lead to needing to do more replacements!

	// Also if you get a request which is just "/", you should return the file "./index.html"

	// You should also return a file not found error if after `timeout`
	// seconds if there is no response from the cache, so you will
	// need to modify the end of the function as well.

	/*** YOUR CODE HERE ***/

	/*** YOUR CODE HERE END ***/

	// Makes the file request object.
	request := fileRequest{filename, make(chan *fileResponse)}
	// Sends a pointer to the file request object to the fileChan so the cache can process the file request.
	fileChan <- &request
	// Returns the result (from the fileResponse channel)
	return <- request.response
}

func getCacheStatus() (response string) {
	// Make a channel for the response of the Capacity request.
	responseChan := make(chan string)
	// Send the response channel to the capacity request channel.
	cacheCapacityChan <- responseChan
	// Return the reply.
	return <- responseChan
}

func CacheClear() (response string) {
	// Send the response channel to the capacity request channel.
	cacheCloseChan <- true
	// We should only return to here once we are sure the currently open cache will not process any more requests.
	// This is because the close channel is blocking until it pulls the item out of there.
	// Now that the cache should be closed, lets relaunch the cache.
	go operateCache()
	return userlib.CacheCloseMessage
}


type cacheEntry struct {
	filename string
	data []byte

	// You may want to add other stuff here...
}

// This function is where you need to do all the work...
// Basically, you need to...

// 1:  Create a map to store all the cache entries.

// 2:  Go into a continual select loop.

// Hint, you are going to want another channel 

func operateCache() {
	/*** YOUR CODE HERE ***/
	// Init our cache to zero bytes.
	currentCapacity := 0
	currentCapacity = currentCapacity // This is just to prevent Golang from yelling at us about unused variables. Removed this once you use the variable elsewhere.
	// Make a file map (this is just like a hashmap in java) for the cache entries.
	fileMap := make(map[string] *cacheEntry)
	fileMap = fileMap // This is just to prevent Golang from yelling at us about unused variables. Removed this once you use the variable elsewhere.
	for {
		// We want to select what we want to do based on what is in different cache channels.
		select {
		case fileReq := <- fileChan:
			fileReq = fileReq // This is just to prevent Golang from yelling at us about unused variables. Removed this once you use the variable elsewhere.
			// Handle a file request here.

		case cacheReq := <- cacheCapacityChan:
			cacheReq = cacheReq // This is just to prevent Golang from yelling at us about unused variables. Removed this once you use the variable elsewhere.
			// Handle a cache capacity request here.

		case <- cacheCloseChan:
			// We want to exit the cache.
			// Make sure you clean up all of your cache state or you will fail most all of the tests!
			return
		}
	}
	/*** YOUR CODE HERE END ***/
}


func main(){
	// Initialize the arguments when the main function is ran. This is to setup the settings needed by
	// other parts of the file server.
	flag.IntVar(&port, "p", 8080, "Port to listen for HTTP requests (default port 8080).")
	flag.IntVar(&capacity, "c", 100000, "Number of bytes to allow in the cache.")
	flag.IntVar(&timeout, "t", 2, "Default timeout (in seconds) to wait before returning an error.")
	flag.StringVar(&workingDir, "d", "public_html/", "The directory which the files are hosted in.")
	// Parse the args.
	flag.Parse()
	// Say that we are starting the server.
	fmt.Printf("Server starting, port %v, working dir: '%s'\n", port, workingDir)
	serverString := fmt.Sprintf(":%v", port)

	// Set up the service handles for certain pattern requests in the url.
	http.HandleFunc("/", handler)
	http.HandleFunc("/cache/", cacheHandler)
	http.HandleFunc("/cache/clear/", cacheClearHandler)

	// Start up the cache logic...
	go operateCache()

	// This starts the web server and will cause it to continue to listen and respond to web requests.
	log.Fatal(http.ListenAndServe(serverString, nil))
}
