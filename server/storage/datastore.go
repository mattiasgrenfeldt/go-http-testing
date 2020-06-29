package storage

import "sync"

// Datastore contains all the state for the basic HTTP server.
type Datastore struct {
	// maps username to userdata
	database      map[string]UserData
	databaseMutex *sync.RWMutex
	// maps session token to username
	session      map[string]string
	sessionMutex *sync.RWMutex
}

// New creates a new Datastore
func New() Datastore {
	return Datastore{
		database:      make(map[string]UserData),
		databaseMutex: &sync.RWMutex{},
		session:       make(map[string]string),
		sessionMutex:  &sync.RWMutex{}}
}

// ReadDatabase reads userdata from the database and returns it based on username. If no entry exists then the bool is false.
func (ds *Datastore) ReadDatabase(username string) (UserData, bool) {
	ds.databaseMutex.RLock()
	userData, exists := ds.database[username]
	ds.databaseMutex.RUnlock()
	return userData, exists
}

// WriteDatabase stores userData in the database for username
func (ds *Datastore) WriteDatabase(username string, userData UserData) {
	ds.databaseMutex.Lock()
	ds.database[username] = userData
	ds.databaseMutex.Unlock()
}

// ReadSession reads username from the session container and returns it based on the token. If no entry exists then the bool is false.
func (ds *Datastore) ReadSession(token string) (string, bool) {
	ds.sessionMutex.RLock()
	username, exists := ds.session[token]
	ds.sessionMutex.RUnlock()
	return username, exists
}

// WriteSession stores username in the session container under the specific token
func (ds *Datastore) WriteSession(token string, username string) {
	ds.sessionMutex.Lock()
	ds.session[token] = username
	ds.sessionMutex.Unlock()
}

// DeleteSession deletes the session associated with token
func (ds *Datastore) DeleteSession(token string) {
	ds.sessionMutex.Lock()
	delete(ds.session, token)
	ds.sessionMutex.Unlock()
}
