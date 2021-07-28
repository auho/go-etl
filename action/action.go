package action

type Action interface {
	Start()
	Done()
	Close()
	GetFields() []string
	Receive([]map[string]interface{})
}
