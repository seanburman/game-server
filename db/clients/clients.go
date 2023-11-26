package db

type Clients struct {
	Mongo mongoClient
}

func NewClients() *Clients {
	mc := NewMongoClient(DatabaseOptions{
		Databases{
			Game: "game",
		},
		Collections{
			User: "users",
		},
	})
	return &Clients{
		Mongo: *mc,
	}
}
