// Copyright 2016 Zhang Peihao <zhangpeihao@gmail.com>

package client

type Client struct {
	// AccessKeyId 阿里云Access Key ID
	AccessKeyId     string
	// AccessKeySecret 阿里云Access Key Secret
	AccessKeySecret string

	httpClient *http.Client
}

func NewClient() (c *Client) {

}

// 处理请求
func (c *Client)Query(req Request, resp interface{}) error {
	if (nil == c.httpClient) {
		return clientError(errors.New("httpClient is nil"))
	}

	if (nil == req) {
		return clientError(errors.New("Request is nil"))
	}

	req.Sign(c.Credentials)
	httpReq, err := req.HttpRequestInstance()
	if (nil != err) {
		return clientError(err)
	}

	//必要头部信息设置
	httpReq.Header.Set("X-SDK-Client", `AliyunLiveGoSDK/` + Version)
	if (req.ResponseFormat() == XMLResponseFormat) {
		httpReq.Header.Set("Content-Type", `application/` + strings.ToLower(XMLResponseFormat))
	}else {
		httpReq.Header.Set("Content-Type", `application/` + strings.ToLower(JSONResponseFormat))
	}

	t0 := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	t1 := time.Now()
	if (nil != err) {
		return clientError(err)
	}

	if c.debug {
		log.Printf("http query %s %d (%v) ", req.String(), httpResp.StatusCode, t1.Sub(t0))
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return clientError(err)
	}

	if c.debug {
		log.Printf("body of response:%s", string(body))
	}

	respUnmarshal := c.responseUnmarshal(req)
	//失败响应处理
	if httpResp.StatusCode >= 400 && httpResp.StatusCode <= 599 {
		errorResponse := ErrorResponse{}
		err = respUnmarshal(body, &errorResponse)
		errorResponse.StatusCode = httpResp.StatusCode
		return &errorResponse
	}

	err = respUnmarshal(body, resp)
	if err != nil {
		return clientError(err)
	}

	if c.debug {
		log.Printf("AliyunLiveGoClient.> decoded response into %#v", resp)
	}

	//if (req.ResponseFormat() == XMLResponseFormat) {
	//	//Xml
	//
	//	//失败响应处理
	//	if httpResp.StatusCode >= 400 && httpResp.StatusCode <= 599 {
	//		errorResponse := ErrorResponse{}
	//		err = xml.NewDecoder(httpResp.Body).Decode(&errorResponse)
	//		xml.Unmarshal(body, &errorResponse)
	//		errorResponse.StatusCode = httpResp.StatusCode
	//		return errorResponse
	//	}
	//	err = xml.NewDecoder(httpResp.Body).Decode(resp)
	//	if err != nil {
	//		return clientError(err)
	//	}
	//}else {
	//	//Json
	//
	//	//失败响应处理
	//	if httpResp.StatusCode >= 400 && httpResp.StatusCode <= 599 {
	//		errorResponse := ErrorResponse{}
	//		err = json.NewDecoder(httpResp.Body).Decode(&errorResponse)
	//		errorResponse.StatusCode = httpResp.StatusCode
	//		return errorResponse
	//	}
	//
	//	err = json.NewDecoder(httpResp.Body).Decode(&errorResponse)
	//	if err != nil {
	//		return clientError(err)
	//	}
	//}

	return nil

}