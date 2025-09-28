package requests

type Request interface {
	Encode() []byte
	Decode() []byte
}

type AccessRequest struct {
	Code       byte
	Identifier byte
}

func NewAccessRequest() AccessRequest {
	return AccessRequest{}
}

func (a *AccessRequest) Encode() []byte {
	// convert code to byte
	// generate identifier
	// temprory length set to 00
	// randomly generate authenticator ( for now rad0 has msg-authenticator disabled so use normal authenticator)
	// encode each avp
	// combine and return
	return []byte{}
}
