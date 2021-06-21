package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/onvif-gf-demos/app/model"

	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/use-go/onvif"
	godevice "github.com/use-go/onvif/device"
	discover "github.com/use-go/onvif/ws-discovery"
)

var Device = new(deviceApi)

type deviceApi struct {}

type Host struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func (d *deviceApi)DeviceDiscovery(r *ghttp.Request)  {
	devices := onvif.GetAvailableDevicesAtSpecificEthernetInterface("以太网")

	_ = r.Response.WriteJsonExit(model.CommonRes{
		Message: "操作成功",
		Code:    0,
		Data: devices,
	})
}

func (d *deviceApi)NewDevice(r *ghttp.Request)  {
	var req *model.NewDeviceReq
	if err := r.Parse(&req); err != nil {
		_ = r.Response.WriteJsonExit(model.CommonRes{
			Message: err.Error(),
			Code: -1,
		})
	}
	fmt.Println(req.IP)
	device, err := onvif.NewDevice(req.IP)
	if err != nil {
		panic(err)
	}
	device.Authenticate(req.UserName, req.Password)
	fmt.Printf("%+v\n", device.GetServices())
	res, err := device.CallMethod(godevice.GetUsers{})
	if err != nil {
		panic(err)
	}
	if res == nil {
		_ = r.Response.WriteJsonExit(model.CommonRes{})
	}
	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Println( "GET USERS:", res.StatusCode, string(bs))

	var hosts []*Host
	devices := discover.SendProbe("以太网", nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			fmt.Println(err.Error())
		} else {
			endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
			scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")
			flag := false
			host := &Host{}
			for _, xaddr := range endpoints {
				xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
				host.URL = xaddr
			}
			if flag {
				break
			}
			for _, scope := range scopes {
				re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
				match := re.FindStringSubmatch(scope.Text())
				host.Name = path.Base(match[0])
			}
			hosts = append(hosts, host)
		}
	}
	_ = r.Response.WriteJsonExit(model.CommonRes{
		Message: "操作成功",
		Code:    0,
		Data: hosts,
	})
}