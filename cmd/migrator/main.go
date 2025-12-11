package main

import (
	"github.com/dliluashvili/cowatchit/db"
	"github.com/dliluashvili/cowatchit/db/migrations"
)

func main() {
	dbconnection := db.New(".env.dev")

	migrations.CreateUserTable(dbconnection)
	migrations.CreateSessionTable(dbconnection)
	migrations.CreateRoomTable(dbconnection)
	migrations.CreateRoomUserTable(dbconnection)
	migrations.CreateRoomMessageTable(dbconnection)
}
