package m7Logger

import (
	"fmt"
	"strconv"

	"log"
	"os"
	"sync"
	"time"
)

type Logs struct {
	LogWarn   *log.Logger
	LogInfo   *log.Logger
	LogErr    *log.Logger
	LogTraf   *log.Logger
	LogVerb   *log.Logger
	LogHverb  *log.Logger
	Filename  string
	MaxSize   int
	Backup    int
	lock      sync.Mutex
	Fp        *os.File
	Err       error
	ModName   string
	DbgLvl    int
	FileStack []string
	count     int
}

var filename = ""

var logObjs []*Logs

func InitLogger(mod, fname string, maxSz, numFiles int) (*Logs, int) {
	pd := Logs{}
	pd.ModName = mod
	pd.Filename = fname
	pd.MaxSize = maxSz
	pd.Backup = numFiles
	pd.Fp, pd.Err = os.OpenFile(pd.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if pd.Err != nil {
		log.Fatal(pd.Err)
		return &pd, 1
	}
	pd.LogErr = log.New(pd.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	pd.LogWarn = log.New(pd.Fp, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	pd.LogInfo = log.New(pd.Fp, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	pd.LogTraf = log.New(pd.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	pd.LogVerb = log.New(pd.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	pd.LogHverb = log.New(pd.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logObjs = append(logObjs, &pd)
	//	pd.FileSize()
	return &pd, 0
}

//Default struct
func Default() (*Logs, error) {
	d := Logs{}
	d.Filename = "log.log"
	d.MaxSize = 1024
	d.Backup = 3
	d.ModName = ""
	d.DbgLvl = 0
	d.Fp, d.Err = os.OpenFile(d.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if d.Err != nil {
		log.Fatal(d.Err)
		return &d, d.Err
	}
	d.LogErr = log.New(d.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogWarn = log.New(d.Fp, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogInfo = log.New(d.Fp, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogTraf = log.New(d.Fp, "TRAF: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogVerb = log.New(d.Fp, "VERB: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogHverb = log.New(d.Fp, "HVERB: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &d, nil

}

func (d *Logs) FileSize() {
	fi, err := d.Fp.Stat()
	if err != nil {
		fmt.Println("unable to find the file info...")
	}
	if fi.Size() > int64(d.MaxSize) {
		d.rotate()
	}
	d.LogErr = log.New(d.Fp, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogWarn = log.New(d.Fp, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogInfo = log.New(d.Fp, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogTraf = log.New(d.Fp, "TRAF: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogVerb = log.New(d.Fp, "VERB: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	d.LogHverb = log.New(d.Fp, "HVERB: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	if len(d.FileStack) > d.Backup {
		d.delete()
	}
}

func (d *Logs) rotate() {
	err := d.newFile()
	if err != nil {
		fmt.Println(err)
	}
	d.Fp, err = os.Create(d.Filename)
	d.Fp, d.Err = os.OpenFile(d.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if d.Err != nil {
		fmt.Println("error in file opening in file")
	}
}

func (d *Logs) newFile() (err error) {

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.Fp != nil {
		err = d.Fp.Close()
		d.Fp = nil
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(d.Filename)
	if err == nil {
		filename = d.Filename + strconv.Itoa(d.count) + "." + time.Now().Format(time.RFC3339)

		err = os.Rename(d.Filename, filename)
		if err != nil {
			return err
		}
		d.FileStack = append(d.FileStack, filename)
		fmt.Println(" ")
		fmt.Println("*************file size exceed backup created ***********************")
		fmt.Println(d.FileStack)
		fmt.Println("********************************************************************")
	}
	d.count += 1
	//d.Fp, err = os.Create(d.Filename)
	return err

}

func (d *Logs) delete() {

	os.Remove(d.FileStack[0])
	d.FileStack = d.FileStack[1:]
	fmt.Println(" ")
	fmt.Println(">>>>>>>>>>>>>>> one file deleted max backup file limit exceed >>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(d.FileStack)

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

}
