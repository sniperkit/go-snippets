package main
import (
    "bytes"
    "fmt" 
    "net"  
    "github.com/ugorji/go/codec"
    //"github.com/ugorji/go-msgpack"
    "encoding/json"
    "reflect"
)

type CloudSample struct {
    CounterName string `codec:"counter_name" json:"counter_name"`
    UserId string `codec:"user_id"`
    MessageSignature string `codec:"message_signature"`
    Timestamp string `codec:"timestamp"`
    ResourceId string `codec:"resource_id"`
    MessageId string `codec:"message_id"`
    Source string `codec:"source"`
    CounterUnit string `codec:"counter_unit"`
    CounterVolume float32 `codec:"counter_volume"`
    ProjectId string `codec:"project_id"`
    *CloudResourceMetadata `codec:"resource_metadata"`
    CounterType string `codec:"counter_type"`
}

type CloudResourceMetadata struct {
    Status string `codec:"status"`
    EphemeralGb int `codec:"ephemeral_gb"`
    DisplayName string `codec:"display_name"`
    Name string `codec:"name"`
    DiskGb int `codec:"disk_gb"`
    KernelId string `codec:"kernel_id"`
    Image map[string]interface{} `codec:"image"`
    RamdiskId string `codec:"ramdisk_id"`
    Vcpus int `codec:"ramdisk_id"`
    MemoryMb int `codec:"memory_mb"`
    InstanceType string `codec:"instance_type"`
    Host string `codec:"host"`
    RootGb string `codec:"root_gb"`
    ImageRef string `codec:"image_ref"`
    Flavor map[string]interface{} `codec:"flavor"`
    AZ string `codec:"OS-EXT-AZ: availability_zone"`
    ImageRefUrl string `codec:"image_ref_url"`
}


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}


func main() {
    p := make([]byte, 4096)
    addr := net.UDPAddr{
        Port: 12234,
        IP: net.ParseIP("0.0.0.0"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    var mh codec.MsgpackHandle
    //mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
    mh.RawToString = false
    //var v interface{} 
    var v CloudSample 
    for {
        buf := bytes.Buffer{}
        //_,remoteaddr, err := ser.ReadFromUDP(p)
        _,_, err := ser.ReadFromUDP(p)
        if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        //fmt.Printf("Read a message from %v %v \n", remoteaddr, p)
        //dec := msgpack.NewDecoder(p)
        buf.Write(p)
        dec := codec.NewDecoder(&buf, &mh)
        err = dec.Decode(&v) 
        //out, err := dec.DecodeBytes()
        if err != nil {
            fmt.Printf("Some error %v", err)
            continue;
        }
        //fmt.Printf("Read a message2 from %s \n", v.CounterName)
        //fmt.Printf("Read a message2 from %v \n", v.ResourceMetadata)
        //fmt.Println(reflect.TypeOf(v.ResourceMetadata))
        //var v2 *CloudResourceMetadata = v.ResourceMetadata
        //fmt.Printf("Read a message2 from %s \n", v2.Flavor)
        //fmt.Println(reflect.TypeOf(v2.Flavor))
        b, err2 := json.Marshal(v)
        fmt.Printf("err2 is %v \n ", err2)
        fmt.Printf("b is %s \n", reflect.TypeOf(b))
        res := string(b)
        fmt.Printf("Read a message from %s \n", res)
        //go sendResponse(ser, remoteaddr)
    }
}
