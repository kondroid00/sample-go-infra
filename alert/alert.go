package alert

type Info struct {
	UserInfo *UserInfo
	Request  *Request
	Tag      map[string]string
	Context  map[string]interface{}
}

type UserInfo struct {
	Email     string
	ID        string
	IPAddress string
	Username  string
}

type Request struct {
	URL         string
	Method      string
	Data        string
	QueryString string
	Cookies     string
	Headers     map[string]string
	Env         map[string]string
}

type Alert interface {
	Send(e error, info *Info)
}
