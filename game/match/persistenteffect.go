package match

var seqID = 1

type PersistentHandlerFunc func(ctx2 *Context, exit func())

type PersistentEffect struct {
	exit   func()
	effect PersistentHandlerFunc
}

func (match *Match) ApplyPersistentEffect(f PersistentHandlerFunc) {
	key := seqID

	fx := PersistentEffect{
		exit:   func() { match.RemovePersistentEffect(key) },
		effect: f,
	}

	match.persistentEffects[key] = fx

	seqID++
}

func (match *Match) RemovePersistentEffect(id int) {
	_, ok := match.persistentEffects[id]
	if ok {
		delete(match.persistentEffects, id)
	}
}
