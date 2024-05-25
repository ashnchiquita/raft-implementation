package server

type Application struct {
	// Application is a key value store stored in memory
	data map[string]string
}

// CONSTRUCTOR
func NewApplication() *Application {
	return &Application{
		data: make(map[string]string),
	}
}

// METHODS
func (a *Application) Ping() string {
	return "PONG"
}

func (a *Application) Get(key string) (string, bool) {
	elem, ok := a.data[key]
	return elem, ok
}

func (a *Application) Set(key, value string) {
	a.data[key] = value
}

func (a *Application) Strlen(key string) (int, bool) {
	elem, ok := a.Get(key)

	// If the key doesn't exist, elem will be zero value of string (empty string)
	return len(elem), ok
}

func (a *Application) Del(key string) (string, bool) {
	elem, ok := a.Get(key)
	delete(a.data, key)

	// If the key doesn't exist, elem will be zero value of string (empty string)
	return elem, ok
}

func (a *Application) Append(key, value string) {
	elem, ok := a.Get(key)

	if ok {
		a.Set(key, elem+value)
	} else {
		a.Set(key, value)
	}
}

// NICE TO HAVE METHODS
func (a *Application) GetAll() map[string]string {
	return a.data
}

func (a *Application) DelAll() {
	a.data = make(map[string]string)
}
