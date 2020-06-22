package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

var ticker *time.Ticker = time.NewTicker(1 * time.Second)

//Downloader struct containing info regarding a executable command
type Downloader struct {
	URL             string            `json:"url"`
	Started         *time.Time        `json:"started"`
	Completed       *time.Time        `json:"completed"`
	Running         bool              `json:"running"`
	Done            bool              `json:"done"`
	FindCmd         *cmd.Cmd          `json:"-"`
	StatusChan      <-chan cmd.Status `json:"-"`
	Ticker          *time.Ticker      `json:"-"`
	StdOut          []string          `json:"std_out"`
	StdErr          []string          `json:"std_err"`
	Error           bool              `json:"error"`
	Filename        *string           `json:"filename"`
	AudioETA        *string           `json:"audio_eta"`
	VideoETA        *string           `json:"video_eta"`
	AudioDownloaded bool              `json:"audio_ready"`
	VideoDownloaded bool              `json:"video_ready"`
	AudioStarted    bool              `json:"-"`
	VideoStarted    bool              `json:"-"`
	FullStdErr      []string          `json:"-"`
}

func extractETA(str string) *string {
	i := strings.LastIndex(str, "ETA: ")
	s := str[i+5:]
	return &s
}

func (d *Downloader) handleOut(list []string) {
	for _, str := range list {
		if strings.HasPrefix(str, "INFO:") || strings.HasPrefix(str, "WARNING:") || strings.HasPrefix(str, "DEBUG:") {
			if strings.Index(str, "Outfile:") >= 0 && strings.Index(str, "audio") >= 0 {
				d.AudioStarted = true
			} else if strings.Index(str, "Outfile:") >= 0 && strings.Index(str, "audio") < 0 {
				d.AudioDownloaded = true
				s := "0:00:00"
				d.AudioETA = &s
				d.VideoStarted = true
			} else if strings.Index(str, "INFO: Merge") >= 0 {
				d.AudioDownloaded = true
				d.VideoDownloaded = true
				s := "0:00:00"
				d.AudioETA = &s
				d.VideoETA = &s
			}
			// if d.AudioStarted {
			// 	fmt.Print("\n")
			// }
			// fmt.Println(str)
			d.StdErr = append(d.StdErr, str)
		} else if strings.HasPrefix(str, "\r") {
			d.AudioStarted = true
			i := strings.LastIndex(str, "\r")
			statusText := str[i+1:]
			if !strings.HasPrefix(d.StdErr[len(d.StdErr)-1], "[") {
				d.StdErr = append(d.StdErr, statusText)
			} else {
				d.StdErr[len(d.StdErr)-1] = statusText
			}
			// fmt.Print("\r", statusText)
			if d.VideoStarted {
				d.VideoETA = extractETA(statusText)
			} else {
				d.AudioETA = extractETA(statusText)
			}
		} else {
			//fmt.Print("\r", str)
			d.StdErr = append(d.StdErr, str)
		}
	}
}

// Max returns the larger of x or y.
func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

// Start to start download
func (d *Downloader) Start() {
	d.FindCmd = cmd.NewCmd("/usr/local/bin/svtplay-dl", d.URL, "--force")
	d.StatusChan = d.FindCmd.Start() // non-blocking
	d.Running = true
	started := time.Now()
	d.Started = &started
	d.Ticker = time.NewTicker(500 * time.Millisecond)
	d.StdErr = make([]string, 0, 10)

	// Print last line of stdout every 2s
	go func() {
		oldNOut := 0
		oldNErr := 0
		//var prevStderr string = ""
		for range d.Ticker.C {
			status := d.FindCmd.Status()
			nOut := len(status.Stdout)
			nErr := len(status.Stderr)
			diffOut := nOut - oldNOut
			//diffErr := len(status.Stderr) - oldNOut

			if len(status.Stderr)-oldNErr > 0 {
				//fmt.Println("nErr:", nErr, "diffErr", len(status.Stderr)-oldNOut, "oldNErr:", oldNErr, "nErr-1", nErr-1, "len", len(status.Stderr))
				d.handleOut(status.Stderr[Max(oldNErr, 0):])
			} else if nErr > 0 && strings.HasPrefix(status.Stderr[len(status.Stderr)-1], "\r") {
				//fmt.Println("--", "nErr:", nErr, "diffErr", len(status.Stderr)-oldNOut, "oldNErr:", oldNErr, "nErr-1", nErr-1)
				d.handleOut([]string{status.Stderr[nErr-1]})
			} else if nErr > 0 {
				//d.handleOut([]string{status.Stderr[nErr-1]})
			}

			if diffOut > 0 {
				str := status.Stdout[nOut]
				fmt.Println(str)
			}

			oldNOut = nOut
			oldNErr = nErr
			d.StdOut = status.Stdout

			if status.Complete {
				break
			}
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Hour)
		d.FindCmd.Stop()
		d.Error = true
		t := time.Now()
		d.Completed = &t
	}()

	// Check if command is done
	select {
	//case finalStatus := <-d.StatusChan:
	case <-d.StatusChan:
		//fmt.Println("statusChan")
		//fmt.Println(finalStatus)
		t := time.Now()
		d.Completed = &t
		d.Done = true
		d.FullStdErr = d.FindCmd.Status().Stderr
	default:
		// no, still running
		//fmt.Println("still running")
		go func() {
			// Block waiting for command to exit, be stopped, or be killed
			//finalStatus := <-d.StatusChan
			<-d.StatusChan
			//fmt.Println(finalStatus)
			d.Done = true
			d.Running = false
			t := time.Now()
			d.Completed = &t
			d.FullStdErr = d.FindCmd.Status().Stderr
		}()
	}
}

var downloaders []Downloader

//DownloadsAll to retrive all downloders
func DownloadsAll() []Downloader {
	return downloaders
}

//AddDownload ..
func AddDownload(url string) Downloader {
	if downloaders == nil {
		downloaders = make([]Downloader, 0, 10)
	}
	download := Downloader{URL: url, Running: false, Error: false, Done: false, AudioDownloaded: false, VideoDownloaded: false}
	downloaders = append(downloaders, download)
	return download
}

func main2() {
	// Start a long-running process, capture stdout and stderr
	//findCmd := cmd.NewCmd("find", "/", "--name", "needle")
}

// package main

// // https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

// import (
// 	"bytes"
// 	//"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"runtime"
// 	"sync"
// )

// // CapturingPassThroughWriter is a writer that remembers
// // data written to it and passes it to w
// type CapturingPassThroughWriter struct {
// 	buf bytes.Buffer
// 	w   io.Writer
// }

// // NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
// func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
// 	return &CapturingPassThroughWriter{
// 		w: w,
// 	}
// }

// // Write writes data to the writer, returns number of bytes written and an error
// func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
// 	w.buf.Write(d)
// 	return w.w.Write(d)
// }

// // Bytes returns bytes written to the writer
// func (w *CapturingPassThroughWriter) Bytes() []byte {
// 	return w.buf.Bytes()
// }

// func main() {
// 	cmd := exec.Command("/usr/local/bin/svtplay-dl", "https://www.svtplay.se/video/24782526/appelkriget")
// 	if runtime.GOOS == "windows" {
// 		cmd = exec.Command("tasklist")
// 	}

// 	var errStdout, errStderr error
// 	stdoutIn, _ := cmd.StdoutPipe()
// 	stderrIn, _ := cmd.StderrPipe()
// 	stdout := NewCapturingPassThroughWriter(os.Stdout)
// 	stderr := NewCapturingPassThroughWriter(os.Stderr)
// 	err := cmd.Start()
// 	if err != nil {
// 		log.Fatalf("cmd.Start() failed with '%s'\n", err)
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	go func() {
// 		_, errStdout = io.Copy(stdout, stdoutIn)
// 		wg.Done()
// 	}()

// 	_, errStderr = io.Copy(stderr, stderrIn)
// 	wg.Wait()

// 	err = cmd.Wait()
// 	if err != nil {
// 		log.Fatalf("cmd.Run() failed with %s\n", err)
// 	}
// 	if errStdout != nil || errStderr != nil {
// 		log.Fatal("failed to capture stdout or stderr\n")
// 	}
// 	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
// 	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
// }
// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"log"
// 	"os/exec"
// 	"bufio"
// )

// func main2() {
// 	cmd := exec.Command("ls", "/Users/marky/Desktop")
// 	//cmd.Dir = entryPath
// 	buf := &bytes.Buffer{}
// 	cmd.Stdout = buf

// 	if err := cmd.Start(); err != nil {
// 		log.Printf("Failed to start cmd: %v", err)
// 		return
// 	}

// 	// Do other stuff while cmd runs in background:
// 	log.Println("Doing other stuff...")

// 	// And when you need to wait for the command to finish:
// 	if err := cmd.Wait(); err != nil {
// 		log.Printf("Cmd returned error: %v", err)
// 		// You may decide to continue or return here...
// 	}

// 	fmt.Println("[OUTPUT:]", buf.String())
// }

// type Data struct {
//     output []byte
//     error  error
// }

// func runCommand(ch chan<- string) {
//     cmd := exec.Command("ls", "-la", "/Users/marky/Desktop")
// 	//data, err := cmd.CombinedOutput()
// 	stderr, _ := cmd.StderrPipe()
//     cmd.Start()

//     scanner := bufio.NewScanner(stderr)
//     scanner.Split(bufio.ScanWords)
//     for scanner.Scan() {
// 		m := scanner.Text()
// 		ch <- m
//         //fmt.Println(m)
// 	}
// 	close(ch)

//     // ch <- Data{
//     //     error:  err,
//     //     output: data,
//     // }
// }

// func main() {
// 	//c := make(chan Data)
// 	c := make(chan string)

//     // This will work in background
//     go runCommand(c)

//     // Do other things here

//     // When everything is done, you can check your background process result
//     res := <-c
//     if res.error != nil {
//         fmt.Println("Failed to execute command: ", res.error)
//     } else {
//         // You will be here, runCommand has finish successfuly
// 		//fmt.Println(string(res.output))
// 		fmt.Println(res)
//     }
// }
