package main

import (
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/mdp"
	"net/url"
	"encoding/json"
)

func main() {
	balanced.Username = config.BalancedUsername
	config.Prod()
	router.Start()
	gophers := 10

	for i := 0; i < gophers; i++ {
		worker := mdp.NewWorker("tcp://localhost:5555", config.UserService, false)
		defer worker.Close()
		go route(worker)
	}

	select {}
}

type Resp struct {
	Body       []byte `json:"body"`
	StatusCode int    `json:"status_code"`
}

type ServiceReq struct {
	Method string
	Path   string
	Body   []byte
}

func route(worker mdp.Worker) {
	for reply := [][]byte{}; ; {
		request := worker.Recv(reply)
		if len(request) == 0 {
			break
		}
		var req *ServiceReq
		json.Unmarshal(request[0], &req)

		//route to function based on the path and method
		route, pathParams, _ := router.FindRoute(req.Path)
		routeMap := route.Dest.(map[string]interface{})
		handler := routeMap[req.Method].(func(map[string]interface{}, []byte) ([]byte, error, int))

		//add url params to params var
		params := make(map[string]interface{})
		for key, value := range pathParams {
			params[key] = value
		}
		//add url query params
		uri, _ := url.Parse(req.Path)
		values := uri.Query()
		for key, value := range values {
			params[key] = value
		}

		//run handler and do standard http stuff(write JSON, return err, set status code)
		jsonData, err, statusCode := handler(params, req.Body)
		if err != nil {
			resp := &Resp{[]byte(`{"description":"` + err.Error() + `"}`), statusCode}
			d, _ := json.Marshal(resp)
			reply = [][]byte{d}
			continue
		}
		resp := &Resp{jsonData, statusCode}
		d, _ := json.Marshal(resp)
		reply = [][]byte{d}
	}
}
