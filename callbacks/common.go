package callbacks

const (
	UrlParamCallbackCommand = "CallbackCommand"
	UrlParamSdkAppid        = "SdkAppid"
	UrlParamContentType     = "contenttype"
	UrlParamClientIP        = "ClientIP"
	UrlParamOptPlatform     = "OptPlatform"
)

const (
	PlatformRESTAPI = "RESTAPI"
	PlatformWeb     = "Web"
	PlatformAndroid = "Android"
	PlatformiOS     = "iOS"
	PlatformWindows = "Windows"
	PlatformMac     = "Mac"
	PlatformUnkown  = "Unkown"
)

type CallbackCommonResponse struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    int    `json:"ErrorCode"`
}

func ResponseOK() CallbackCommonResponse {
	return CallbackCommonResponse{ActionStatus: "OK"}
}

func ResponseFail(code int, reason string) CallbackCommonResponse {
	return CallbackCommonResponse{ActionStatus: "FAIL", ErrorCode: code, ErrorInfo: reason}
}
