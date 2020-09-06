package packet

type Type uint8

//透传类型
const (
	TypeNone Type = iota

	//基础
	TypeConnect
	TypeConnectAck
	TypeDisconnect
	TypeHeartBeak
	TypePing
	TypePong
	TypeError

	//透传
	TypeTransfer = iota + 10
	TypeOnline //上线
	TypeOffline //掉线

)


