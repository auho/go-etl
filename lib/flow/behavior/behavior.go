package behavior

type Behavior interface {
}

type BehaviorVoid interface {
	Do(map[string]interface{})
}

type BehaviorInsert interface {
	GetKeys() []string
	GetFields() []string
	Do(map[string]interface{}) [][]interface{}
}

type BehaviorUpdate interface {
	GetFields() []string
	Do(map[string]interface{}) map[string]interface{}
}

type Insert struct {
	keys        []string
	trickInsert TrickInsert
}

func NewInsert(keys []string, ti TrickInsert) *Insert {
	t := &Insert{}
	t.keys = keys
	t.trickInsert = ti

	return t
}

func (ti *Insert) GetFields() []string {
	return ti.keys
}

func (ti *Insert) GetKeys() []string {
	return ti.trickInsert.GetKeys()
}

func (ti *Insert) Do(item map[string]interface{}) [][]interface{} {
	if item == nil {
		return nil
	}

	contents := make([]string, 0)
	for _, key := range ti.keys {
		contents = append(contents, key)
	}

	return ti.trickInsert.Insert(contents)
}

type Update struct {
	keys        []string
	trickUpdate TrickUpdate
}

func NewUpdate(keys []string, tu TrickUpdate) *Update {
	t := &Update{}
	t.keys = keys
	t.trickUpdate = tu

	return t
}

func (tu *Update) GetFields() []string {
	return tu.keys
}

func (tu *Update) Do(item map[string]interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := make([]string, 0)
	for _, key := range tu.keys {
		contents = append(contents, key)
	}

	return tu.trickUpdate.Update(contents)
}
