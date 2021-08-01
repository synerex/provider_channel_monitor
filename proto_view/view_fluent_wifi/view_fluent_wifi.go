package view_fluent_wifi

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"

	view "github.com/synerex/provider_channel_monitor/proto_view"

	fluentd "github.com/synerex/proto_fluentd"

	api "github.com/synerex/synerex_api"

	sxutil "github.com/synerex/synerex_sxutil"
)

var totalTerminals int32
var amps map[string]*AMPM
var (
	temp_data = make(map[string]string)
	temp_time = make(map[string][2]string)
)

type AMPM struct { // for signal info
	AMPname string
	lastTS  int64
	count   int32
}

func init() {
	fmt.Printf("Initial view Fluent-WiFi\n")
	view.AddSubscriber(7, subscribeFluentdSupply)

	amps = make(map[string]*AMPM)

}

const dateFmt = "2006-01-02T15:04:05.999999Z"
const truncatemin = "2006-01-02T15:04"

func outputWiFi(bc, ob, dt string) {

	log.Printf("WiFi %s, %s, %s", bc, ob, dt)

}

type WiFiInfo struct {
	Timestamp float64
	HMAC      string
	AMP       string
	Power     int32
	OID       string
	Flag      string
	Counts    int32
}

func objScanStore(wi *WiFiInfo) {
	var bucketName string
	var objectName string

	tsint, tsfrac := math.Modf(wi.Timestamp)
	tsfracint := int64(tsfrac * 1000000)
	timeobj := time.Unix(int64(tsint), tsfracint)
	wi.AMP = "AMPM18-HZ" + wi.AMP //  add identity

	bucketName = "higashiyama" // fixed .. TODO:

	//		objectName := "sensor_id/year/month/date/hour/min"
	objectName = fmt.Sprintf("WiFi/%s/%4d/%02d/%02d/%02d/%02d", wi.AMP, timeobj.Year(), timeobj.Month(), timeobj.Day(), timeobj.Hour(), timeobj.Minute())

	line := fmt.Sprintf("%s,%s,%s,%d,%s,%s,%d", timeobj.Format(dateFmt), wi.AMP, wi.HMAC, wi.Power, wi.OID, wi.Flag, wi.Counts)
	if _, ok := temp_time[wi.AMP]; !ok {
		temp_str := [2]string{wi.AMP, objectName}
		temp_time[wi.AMP] = temp_str
		log.Print(wi.AMP, temp_time[wi.AMP])
	}
	if temp_time[wi.AMP][0] == wi.AMP && temp_time[wi.AMP][1] != objectName {
		outputWiFi(bucketName, temp_time[wi.AMP][1], temp_data[wi.AMP])
		temp_data[wi.AMP] = ""
		temp_str := [2]string{temp_time[wi.AMP][0], objectName}
		temp_time[wi.AMP] = temp_str
	}

	if temp_data[wi.AMP] != "" {
		temp_data[wi.AMP] = temp_data[wi.AMP] + line + "\n"
	} else {
		temp_data[wi.AMP] = line + "\n"
	}

}

func convertWiFi(dt map[string]interface{}) {
	host := dt["h"].(string)
	wifi := dt["d"].(string)

	wbases := strings.Split(wifi, "\n")
	for _, s := range wbases {
		vals := strings.Split(s, ",")
		//		log.Printf("Split into %d", len(vals))
		if len(vals) < 5 {
			continue
		}
		ts, _ := strconv.ParseFloat(vals[0], 64)
		pow, _ := strconv.ParseInt(vals[2], 10, 32)
		cts, _ := strconv.ParseInt(vals[4], 10, 32)
		totalTerminals++

		stFlag := ""
		if len(vals[3]) > 6 {
			stFlag = vals[3][7:]
		}
		wi := WiFiInfo{
			Timestamp: ts, //vals[0],
			HMAC:      vals[1],
			AMP:       host[9:],
			Power:     int32(pow),
			OID:       vals[3][:6],
			Flag:      stFlag,
			Counts:    int32(cts),
		}

		//        ds.store(&wi)
		objScanStore(&wi)

	}
	log.Printf("WiFi from %s wifi [%d] total count %d\n", host, len(wbases)-1, totalTerminals)
}

func checkAMPM() string {
	tm := time.Now().Unix()
	st := make([]string, 0)

	for n, v := range amps {
		nm := n[10:] // slice from "AMPM18-HZ0XX"
		if tm-v.lastTS < 10 {
			st = append(st, nm)
		}
	}
	sort.Slice(st, func(i, j int) bool {
		ii, _ := strconv.Atoi(st[i])
		jj, _ := strconv.Atoi(st[j])
		return ii < jj
	})

	return fmt.Sprintf("%d/%d %v", len(st), len(amps), st)

}

func jsonDecode(jsonByte []byte) map[string]interface{} {
	var dt map[string]interface{}
	err := json.Unmarshal(jsonByte, &dt)
	if err == nil {
		return dt
	}
	fmt.Println("jsonDecodeErr:", err)
	return nil
}

func base64UnCompress(str string) []byte {
	data, _ := base64.StdEncoding.DecodeString(str)
	dt1, err := zlib.NewReader(bytes.NewReader(data))
	if err == nil {
		buf, err := ioutil.ReadAll(dt1)
		if err == nil {
			return buf
		} else {
			fmt.Println("base64UncompErr:", err)
		}
	} else {
		fmt.Println("base64UncompErr:", err)
	}
	return []byte(" ")
}

// callback for each fluentd Supply
func fluentSupplyCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {
	// check if demand is match with my supply.
	//	log.Println("Got Fluentd Supply callback")

	record := &fluentd.FluentdRecord{}
	if sp.Cdata != nil {	
		err := proto.Unmarshal(sp.Cdata.Entity, record)

		if err == nil {
			//		log.Println("Got record:", record.Tag, record.Time)
			recordStr := *(*string)(unsafe.Pointer(&(record.Record)))
			replaced := strings.Replace(recordStr, "=>", ":", 1)
			dt0 := jsonDecode([]byte(replaced))
			if dt0 != nil {
				buf := base64UnCompress(dt0["m"].(string))
				if len(buf) > 1 {
					dt := jsonDecode(buf)
					if dt != nil {
						if record.Tag == "ampsense.pack.test.signal" {
							//				log.Printf("ID:%v, %v, %v", dt["a"], dt["ts"], dt["g"])
							ampName := dt["a"].(string)
							amp := amps[ampName]
							if amp == nil {
								amps[ampName] = &AMPM{
									AMPname: ampName,
									count:   0,
								}
								amp = amps[ampName]
							}
							amp.lastTS = time.Now().Unix()
							amp.count++
							sxutil.SetNodeStatus(totalTerminals, checkAMPM())
						} else if record.Tag == "ampsense.pack.packet.test" {
							//						log.Printf("packet:%v\n", dt)
							convertWiFi(dt)
						} else { // unknown data.
							log.Printf("UNmarshal Result: %s, %v\n", record.Tag, dt)
						}
					}
				}
			}
			return
		}
	}
	log.Printf("Unmarshal error on View_Pcoutner %s", sp.SupplyName)
}

func subscribeFluentdSupply(client *sxutil.SXServiceClient) {
	//
	log.Printf("Subscribe Fluentd Supply")
	//	ctx := context.Background() //
	//	client.SubscribeSupply(ctx, fluentSupplyCallback)
	sxutil.SimpleSubscribeSupply(client, fluentSupplyCallback) // error prone..
	//	log.Printf("Error on subscribe with fluentd")
}
