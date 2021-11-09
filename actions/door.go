package actions

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type OpenDoorResponse struct {
	Status  string `json:"status"`
	Result  string `json:"result"`
	Data    string `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
type OpenDoorBody struct {
	HouseHostId string `json:"houseHostId"`
	PeopleId    string `json:"peopleId"`
	RoleType    string `json:"roleType"`
}

type OpenDoorRequest struct {
	Token string       `json:"token"`
	Body  OpenDoorBody `json:"body"`
}

func OpenTheDoor(req OpenDoorRequest) {
	client := resty.New()
	// 重试次数
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		// 状态码异常
		if response.StatusCode() > 200 {
			return true
		}
		var r OpenDoorResponse
		err = json.Unmarshal(response.Body(), &r)
		if err != nil {
			return true
		}
		// 业务状态码异常
		return r.Code != "000000"
	}).SetRetryCount(2)

	var res OpenDoorResponse

	resp, err := client.R().
		SetQueryParam("token", req.Token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(OpenDoorBody{
			HouseHostId: req.Body.HouseHostId,
			PeopleId:    req.Body.PeopleId,
			RoleType:    req.Body.HouseHostId,
		}).
		SetResult(&res).
		Post("https://pabaspmj.szxhdz.com:18000/xhapp/service/iacs/info/house/host/commandByHouseHostId")
	if err != nil {
		fmt.Printf("发送开门请求失败")
	}
	fmt.Printf("原始请求结果 resp %v \n", resp)
	fmt.Printf("解析请求结果Message = %s ", res.Message)
}
