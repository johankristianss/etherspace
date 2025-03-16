package dht

const (
	MSG_PING_REQ int = iota
	MSG_PING_RESP
	MSG_FIND_CONTACTS_REQ
	MSG_FIND_CONTACTS_RESP
	MSG_PUT_REQ
	MSG_PUT_RESP
	MSG_GET_REQ
	MSG_GET_RESP
)

func TypeToString(t int) string {
	switch t {
	case MSG_PING_REQ:
		return "MSG_PING_REQ"
	case MSG_PING_RESP:
		return "MSG_PING_RESP"
	case MSG_FIND_CONTACTS_REQ:
		return "MSG_FIND_CONTACTS_REQ"
	case MSG_FIND_CONTACTS_RESP:
		return "MSG_FIND_CONTACTS_RESP"
	case MSG_PUT_REQ:
		return "MSG_PUT_REQ"
	case MSG_PUT_RESP:
		return "MSG_PUT_RESP"
	case MSG_GET_REQ:
		return "MSG_GET_REQ"
	case MSG_GET_RESP:
		return "MSG_GET_RESP"
	}
	return "UNKNOWN"
}
