package main
import (
    "bytes"
    "fmt" 
    "net"  
    "reflect"
    //"errors"
    "time"
    "strings"
    //"strconv"
    "net/http"
    "io/ioutil"

    "crypto/hmac"
    "crypto/sha1"

    "encoding/base64"
    "encoding/json"

    "github.com/ugorji/go/codec"
)

type Dimension struct {
    Name string `json:"name"`
    Value string `json:"value"`
}

type MetricDatum struct {
    MetricName string `json:"metric_name"`
    Unit string `json:"unit"`
    Value string `json:"value"`
    DimensionList []Dimension `json:"dimension_list"`
}

type MetricData struct {
    Version string `json:"version"`
    Namespace string `json:"namespace"`
    MetricDatumList []MetricDatum `json:"metric_datum_list"`
}

func postMetricDataToWatch(reqBody []byte) {
    accessID := "CWXYQX8Z1LDPZ5DJ14NG"
    accessKey := "k/MOjnPk2yrjNecndc212kNocLX1yPVaq4PTNg=="
    action := "PutMetricData"
    
    date, auth, _ := generateAuth(accessID, accessKey, action) 

    url := "http://cloudwatch-apix.lenovows.com/v1/actions/PutMetricData"

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("date", date)
    req.Header.Set("Authorization", auth)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func getMetricType(metricName string) string {
    metricType := "default"
    metricTypeList := []string{"disk", "memory", "network", "cpu"}
    for i := 0; i < len(metricTypeList); i++ {
        if strings.Contains(metricName, metricTypeList[i]) {
            metricType = metricTypeList[i]
            break;
        }
    }
    return metricType
}

func getTargetUnit(unit string) string {
    switch unit {
    case "%":
        unit = "Percent"
    case "B/s":
        unit = "Bytes/Second"
    case "packet/s":
        unit = "Count/Second"
    case "request/s":
        unit = "Count/Second"
    default:
        unit = "Count/Second"
    }
    return unit
}

func generateAuth(accessID, accessKey, action string) (string, string, error) {
    date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
    toSign := "GET\n"
    toSign += "\n"
    toSign += "application/json\n"
    toSign += date + "\n"
    toSign += "/" + action
    hash := hmac.New(sha1.New, []byte(accessKey))
    hash.Write([]byte(toSign))
    signed := base64.StdEncoding.EncodeToString(hash.Sum(nil))
    auth := "LWS " + accessID + ":" + signed 
    return date, auth, nil
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
    mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
    mh.RawToString = true 
    var v map[string]interface{} 

    for {
        _, _, err := ser.ReadFromUDP(p)
        if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        buf := bytes.Buffer{}
        buf.Write(p)
        dec := codec.NewDecoder(&buf, &mh)
        err = dec.Decode(&v) 
        if err != nil {
            fmt.Printf("Some error %v", err)
            continue;
        }

        var (
            metricName, unit string
            metricValue float64
            resourceId  string
            projectId string
        )

        if value, ok := v["counter_name"]; ok {
            metricName = value.(string)
        } else {
            fmt.Printf("Invalid metric, metric does not exist\n")
            continue
        }

        if value, ok := v["counter_unit"]; ok {
            unit = value.(string)
        } else {
            fmt.Printf("Unit does not exist\n")
            continue
        }

        if value, ok := v["counter_volume"]; ok {
            metricValue = value.(float64)
        } else {
            fmt.Printf("data point does not exist\n")
            continue
        }

        if value, ok := v["resource_id"]; ok {
            resourceId = value.(string)
        } else {
            fmt.Printf("Resource_id does not exist\n")
            continue
        }

        if value, ok := v["project_id"]; ok {
            projectId = value.(string)
        } else {
            fmt.Printf("Project_id does not exist \n")
            continue
        }

        unit = getTargetUnit(unit)
        metricType := getMetricType(metricName)

        //date, auth, _ := generateAuth(accessID, accessKey, "PutMetricData")
        //reqBody := `{"version":"2010-08-01","namespace":"lcp","metric_datum_list":[{"metric_name":"` + metricName + `","unit":"` + unit + `","value":` + strconv.FormatFloat(metricValue, 'f', -1, 64) + `,"dimension_list":[{"name":"resource_id","value":"` + resourceId + `"},{"name":"project_id","value":"` + projectId + `"},{"name":"metric_type","value":"` + metricType + `"}]}]}`
        dimensionList := []map[string]string{
            {"name": "resource_id", "value": resourceId},
            {"name": "project_id", "value": projectId},
            {"name": "metric_type", "value": metricType},
        }
        metricDatum := []map[string]interface{} {
            {
                "metric_name": metricName,
                "unit": unit,
                "value":metricValue,
                "dimension_list": dimensionList,
            },
        }
        metricData := map[string]interface{} {
            "version": "2010-08-01",
            "namespace": "lcp",
            "metric_datum_list": metricDatum,
        }
        //fmt.Println(reqBody)
        fmt.Println(metricData)

        
        //fmt.Printf("v is %s", v["counter_name"])
        b, err := json.Marshal(metricData)
        if err != nil {
            fmt.Printf("Some error when marshalling for josn %v", err)
        }
        fmt.Printf("The type of json is %s\n", reflect.TypeOf(b))
        fmt.Printf("Read a message from %s \n", b)
        go postMetricDataToWatch(b)
        //go sendResponse(ser, remoteaddr)
    }
}
