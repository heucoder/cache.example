package main

type PeerPicker interface {
	PickPeer(addr string) (PeerGetter, bool) //addr->Getter
}

type PeerGetter interface { //Get Val
	Get(group string, key string) ([]byte, error)
}
