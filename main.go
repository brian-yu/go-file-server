package skeleton

import (
	"flag"
	"fmt"
	"github.com/61c-teach/sp19-proj5-alternate-dev/skeleton/userlib"
	_ "github.com/61c-teach/sp19-proj5-alternate-dev/skeleton/userlib"
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
	response, err := userlib.ReadFile(r.URL.Path[1:])
	if err != nil {
		http.Error(w, userlib.FILEERRORMSG, userlib.FILEERRORCODE)
		return
	}

	// This will automagically set the right content type for the
	// reply as well.
	w.WriteHeader(userlib.SUCCESSCODE)
	w.Write(response)
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(getCacheStatus()))
}

// The structure used for responding to file requests.
// It contains the file contents (if there is any)
// or the error returned when accessing the file.
// You can add fields in here so long as you do not 
// change or remove the existing fields.
type fileResponse struct {
	responseData []byte
	responseError error
}

// To request files from the cache, we send a message that 
// requests the file and provides a channel for the return
// information.
type fileRequest struct {
	filename string
	response chan(*fileResponse)
}

var port int
var capacity int
var capacityString string = "Cache status:  # of entries %v\ntotal bytes occupied by entries %v\nmax allowed capacity %v\n"
var timeout int

// These names must not change or you will not pass any of the autograder tests!
var fileChan chan(*fileRequest) = make(chan(*fileRequest))
var cacheCapacityChan chan(chan(string)) = make(chan(chan(string)))
var cacheCloseChan chan(chan(string)) = make(chan(chan(string)))

// A wrapper function that does the actual getting of the file
func getFile(filename string) (response *fileResponse) {
	// You need to add sanity checking here: The requested file
	// should be made relative (strip out leading "/" characters,
	// then have a "./" put on the start, and if there is ever the
	// string "/../", replace it with "/", the string "\/" should
	// be replaced with "/", and finally any instances of "//" (or
	// more) should be replaced by a single "/".
	// Hint: A replacement may lead to needing to do more replacements!

	// You should also return a file not found error if after `timeout`
	// seconds if there is no response from the cache, so you will
	// need to modify the end of the function as well.
	
	request := fileRequest{filename, make(chan(*fileResponse))}
	fileChan <- &request
	return <- request.response
}

func getCacheStatus() (response string) {
	responseChan := make(chan(string))
	cacheCapacityChan <- responseChan
	return <- responseChan
}


type cacheEntry struct {
	request string
	response *fileResponse

	// You may want to add other stuff here...
}

// This function is where you need to do all the work...
// Basically, you need to...



// 1:  Create a map to store all the cache entries.  Your cache
// entry stru

// 2:  Go into a continual select loop.

// Hint, you are going to want another channel 

func operateCache() {
	currentCapacity := 0
	currentCapacity = currentCapacity
	fileMap := make(map[string] *fileResponse)
	fileMap = fileMap
	for {
		select {
		case fileReq := <- fileChan:
			fileReq = fileReq
			// Handle a file request here.

		case cacheReq := <- cacheCapacityChan:
			cacheReq = cacheReq
			// Handle a cache capacity request here.

		case <- cacheCloseChan:
			// We want to exit the cache.
			// Make sure you clean up all of your cache state or you will fail most all of the tests!
			return
		}

	}

}


func main(){
	flag.IntVar(&port, "p", 8080, "Port to listen for HTTP requests (default port 8080).")
	flag.IntVar(&capacity, "c", 100000, "Number of bytes to allow in the cache.")
	flag.IntVar(&timeout, "t", 2, "Default timeout (in seconds) to wait before returning an error")
	flag.Parse()
	fmt.Printf("Server starting, port %v\n", port)
	serverString := fmt.Sprintf(":%v", port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/cache/", cacheHandler)

	// Start up the cache logic...
	go operateCache()
	
	log.Fatal(http.ListenAndServe(serverString, nil))
}
