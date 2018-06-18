import (
  "fmt"
  "sync"
  "sync/atomic"
  "time"
  "io"
  "os"
)
type DataFile interface{
  Read()(rsn int64, d Data, err error)
  Write(d Data)(wsn int64, err error)
  Rsn() int64
  Wsn() int64
  Datalen() uint32
}

type Data []byte

type myDataFile struct{
  f *os.File
  fmutex sync.RWMutex
  woffset int64
  roffset int64
  wmutex sync.Mutex
  rmutex sync.Mutex
  dataLen uint32
}

func NewDataFile(path string, dataLen uint32)(DataFile, error){
  f,err:=os.Create(path)
  if err!=nil{
    return nil, err
  }
  if dataLen == 0{
    return nil, errors.New("Invalid data length!")
  }
  df := &myDataFile{f: f, dataLen: dataLen}
  return df,nil
}

func(df *myDataFile) Read()(rsn int64, d Data, err error){
  var offset int64
  df.rmutex.Lock()
  offset = df.roffset
  df.roffset += int64(df.dataLen)
  df.rmutex.Unlock()

  rsn = offset / int64(df.dataLen)
  df.fmutex.RLock()
  defer df.fmutex.RUnlock()
  bytes := make([]byte, df.dataLen)
  _, err = df.f.ReadAt(bytes, offset)
  if err != nil{
    return
  }
  d = bytes
  return
}

// better func Read()

func(df *myDataFile) Read()(rsn int64, d Data, err error){
  // 省略若干代码
  
  rsn = offset / int64(df.dataLen)
  bytes := make([]byte, df.dataLen)
  for {
    df.fmutex.RLock()
    _,err = df.f.ReadAt(bytes,offset)
    if err!=nil{
      if err == io.EOF{
        df.fmutex.RUlock()
        continue
      }
      df.fmutex.RUnLock()
      return
    }
    d = bytes
    df.fmutex.RUnlock()
    return
  }
}

func (df *myDataFile) Write(d Data)(wsn int64, err error){
  //读取并更新偏移量
  var offset int64
  df.wmutex.Lock()
  offset = df.woffset
  df.woffset += int64(df.dataLen)
  df.wmutex.Unlock()

  //写入一个数据块
  wsn = offset / int64(df.dataLen)
  var bytes []byte
  if len(d) > int(df.dataLen){
     bytes = d[0:df.dataLen]
  }else{
    bytes = d
  }
  df.fmutex.Lock()
  df.fmutex.Unlock()
  _,err = df.f.Write(bytes)
  return
}

func (df *myDataFile) Rsn() int64{
  df.rmutex.Lock()
  defer df.rmutex.Unlock()
  return df.roffset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64{
  df.rmutex.Lock()
  defer df.wmutex.Unlock()
  return df.woffset / int64(df.dataLen)
}

// 加入条件变量的Read
func (df *myDataFile) Read()(rsn int64, d Data, err error){
  //省略若干代码

  //读取一个数据块
  rsn = offset / int64(df.dataLen)
  bytes := make([]byte, df.dataLen)
  df.fmutex.RLock()
  defer df.fmutex.RUnlock()
  for {
    _,err = df.f.ReadAt(bytes, offset)
    if err!=nil{
      if err == io.EOF{
        df.rcond.Wait()
        continue
      }
      return
    }
    d = bytes
    return
  }
}

// 加入条件变量的Write
func (df *myDataFile) Write(d Data)(wsn int64, err error){

  //省略若干语句
  var bytes []byte
  df.fmutex.Lock()
  defer df.fmutex.Unlock()
  _,err = df.f.Write(bytes)
  df.rcond.Signal()
  return
}





































