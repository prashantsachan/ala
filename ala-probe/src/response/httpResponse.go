package response

import(
    log "github.com/Sirupsen/logrus"
    "net/http"
    "io/ioutil"
    "encoding/json"
    )
const HTTP_FIELD_STATUS = "status"
const HTTP_FIELD_HEADERS = "headers"
const HTTP_FIELD_BODY = "body"
const HTTP_FIELD_VERSION = "version"
//HttpResponse is model class for HTTP response
type HttpResponse struct{
    Resp http.Response
}
// httpRespModel represents an HTTP response as a model
type httpRespModel struct{
    //status is the StatusCode of http reponse
    status int 
    //headers passed with response
    headers http.Header 
    // body is string reprensentation of response body
    body string 
    // version of HTTP protocol used(Proto)
    version string 
}
func (this *httpRespModel) asMap()map[string]interface{}{
    m:= make(map[string]interface{})
    m[HTTP_FIELD_STATUS] = this.status
    m[HTTP_FIELD_HEADERS] = this.headers
    m[HTTP_FIELD_VERSION] = this.version
    //Try to convert this to map assuming it to be json response
    var respData  map[string]interface{}
    jsonErr:=json.Unmarshal([]byte(this.body), &respData)
    if jsonErr !=nil{
        log.WithFields(log.Fields{"module": "httpResponse","error":jsonErr, "value":this.body}).Debug(
            "unable to parse resp body as JSON, passing it as string")
        m[HTTP_FIELD_BODY] = this.body   
    }else{
        m[HTTP_FIELD_BODY] = respData
    }
    return m;
}
func (this *HttpResponse) GetType()string{
    return "HTTP"
}

//AsMap returns HTTP response as a map
// response should contain keys: 'status','headers','body' & 'version'
func (this *HttpResponse) AsMap()map[string]interface{}{
    defer this.Resp.Body.Close()
    respBody,err := ioutil.ReadAll(this.Resp.Body)
    if err !=nil{
        log.WithFields(log.Fields{"module": "httpResponse","error":err}).Info("error reading response body")
        respBody = nil    
    }
    modelResp:= httpRespModel{ this.Resp.StatusCode, this.Resp.Header,string(respBody), this.Resp.Proto}
    return modelResp.asMap()
}