package api

import (
	"github.com/beevik/etree"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	wsdiscovery "github.com/use-go/onvif/ws-discovery"
	"net/http"
	"onvif-gf-demos/app/model"
	"onvif-gf-demos/app/service"
	"path"
	"regexp"
	"strings"
)

var ONVIF = &ONVIFApi{}

type ONVIFApi struct {}

func (o *ONVIFApi)PostMethod(r *ghttp.Request) {
	serviceName := r.GetString("service")
	methodName := r.GetString("method")
	username := r.GetHeader("username")
	pass := r.GetHeader("password")
	xaddr := r.GetHeader("xaddr")
	acceptedData := r.GetBody()

	message, err := service.CallNecessaryMethod(serviceName, methodName, string(acceptedData), username, pass, xaddr)
	if err != nil {
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.WriteXml(err.Error())
	} else {
		r.Response.WriteStatus(http.StatusOK, message)
	}
}
func (o *ONVIFApi)Discovery(r *ghttp.Request) {
	r.Header.Add("Access-Control-Allow-Origin", "*")
	r.Header.Add("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	interfaceName := g.Cfg().GetString("other.interface")
	devices := wsdiscovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})

	var result []*model.DiscoveryRes
	for _, j := range devices {
		res := &model.DiscoveryRes{}
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			r.Response.WriteHeader(http.StatusBadRequest)
			r.Response.WriteXml(err.Error())
		} else {
			endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
			scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")
			flag := false
			for _, xaddr := range endpoints {
				xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
				if res.URL != "" {
					flag = true
					break
				}
				res.URL = xaddr
			}
			if flag {
				break
			}
			for _, scope := range scopes {
				re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
				if match := re.FindStringSubmatch(scope.Text()); len(match) > 0 {
					res.Name = path.Base(match[0])
				}
				manu := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/manufacturer\/[A-Za-z0-9-]+`)
				if manufacturerSlice := manu.FindStringSubmatch(scope.Text()); len(manufacturerSlice) > 0 {
					res.Manufacturer = path.Base(manufacturerSlice[0])
				}
				vn := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/VideoSourceNumber\/[A-Za-z0-9-]+`)
				if videoSourceNumberSlice := vn.FindStringSubmatch(scope.Text()); len(videoSourceNumberSlice) > 0 {
					res.VideoSourceNumber = path.Base(videoSourceNumberSlice[0])
				}
			}
		}
		result = append(result, res)
	}
	r.Response.WriteStatus(http.StatusOK, result)
}