package pester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAsynGetCallBack(t *testing.T)  {
	c := New()
	c.Concurrency = 3
	c.KeepLog = true
	c.Backoff = ExponentialBackoff
	c.MaxRetries = 3

	url := "http://192.168.20.14/user/getPublicKey"

	stopChan := make(chan string)

	callBack := func (resp *http.Response, err error, logs string, transmissionParams ...interface{} ) {
		if err != nil {
			t.Errorf("err:%v\nlogs:%v, \ntransmissionParams:%v", err, logs, transmissionParams)
			stopChan <- ""
			return
		}
		if nil == resp || nil == resp.Body{
			t.Errorf("response nil, logs:%s,  transmissionParams:%v", logs, transmissionParams)
			stopChan <- ""
			return
		}
		body,er := ioutil.ReadAll(resp.Body)
		t.Logf("resp: %v, err:%v, logs:%s, transmissionParams:%v\n", string(body), er, logs, transmissionParams)
		stopChan <- ""
	}

	err := c.GetAsyn(url, 50*time.Millisecond, callBack, 111)
	if err != nil {
		t.Errorf(c.LogString())
		return
	}

	<-stopChan
}


var postStopChan = make(chan string)

func asynPostCallBack(resp *http.Response, err error, logs string,  transmissionParams ...interface{} ) {
	if err != nil {
		fmt.Printf("err:%v\nlogs:%v, \ntransmissionParams:%v", err, logs, transmissionParams)
		postStopChan <- ""
		return
	}
	if nil == resp || nil == resp.Body {
		fmt.Printf("response nil, logs:%s,  transmissionParams:%v", logs, transmissionParams)
		postStopChan <- ""
		return
	}
	body,er := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp: %v, err:%v, logs:%s, transmissionParams:%v\n", string(body), er, logs, transmissionParams)
	postStopChan <- ""
}

func TestAsynPostCallBack(t *testing.T) {
	c := New()
	c.KeepLog = true
	c.Backoff = DefaultBackoff
	c.MaxRetries = 5

	url := "http://192.168.20.14:8081/abTest/v1/experiment/name/test2/value"

	bodyParam := map[string]interface{}{"slotId":"1000"}
	body, _ := json.Marshal(bodyParam)
	err := c.PostAsyn(url, body,  100*time.Millisecond, asynPostCallBack, 111, 112,"sss")
	if err != nil {
		t.Errorf(c.LogString())
		return
	}
	t.Logf("waiting for callback")
	<-postStopChan
}
