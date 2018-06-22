package gorlim

import "os"
import "time"
import "syscall"
//import "fmt"

type IntegerSocketListener struct {
    event chan int
    pipeName string
}

func GetPushSocketListener() *IntegerSocketListener {
    if pushSocketListener == nil {
    	os.MkdirAll(getPushPipeDir(), 0777)
        pushSocketListener = CreateSocketListener(getPushPipeName())
    }
    return pushSocketListener
}

func CreateSocketListener(pipeName string) *IntegerSocketListener {
	sckListener := &IntegerSocketListener{}
	sckListener.event = make(chan int, 16) // TODO buffer size
	sckListener.pipeName = pipeName
	syscall.Mkfifo(pipeName, 0666)
	subscribeToPushEvent(pipeName, sckListener.event)
	return sckListener
}

func (isl* IntegerSocketListener) GetSocketWriteEvent() <-chan int {
  return isl.event;
}

func (isl* IntegerSocketListener) Free() {
  close(isl.event)
  os.Remove(isl.pipeName)
}

var pushSocketListener *IntegerSocketListener

func getPushPipeName() string {
    return getPushPipeDir() +  "/pushntfifo"
}

func getPushPipeDir() string {
    return os.Getenv("HOME") +  "/syncpipes"
}

func subscribeToPushEvent(pipe string, notify chan<- int) {
  	go func() {
  		f, err := os.Open(pipe)
		if err != nil {
			panic (err)
		}
  		for {
  	   		bytes :=  make([]byte, 16)
  	   		n, _ := f.Read(bytes)
         	if n != 0 {
         		//fmt.Println("Read bytes")
           		repoId := 0
           		for i := 0; i < n - 1; i++ {
             		repoId = repoId * 10 + int(bytes[i] - 48);
           		}
           		notify <- repoId // TODO
         	} else { //ugly
         	//fmt.Println("Before sleep")
            	time.Sleep(time.Second)
            //fmt.Println("After sleep")
        	}
  		}
  	}()
}
