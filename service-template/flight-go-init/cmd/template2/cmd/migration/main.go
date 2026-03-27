package main

import (
	_ "{{MODULE_NAME}}/migrations"
	_ "{{MODULE_NAME}}/migrations/{{PROVIDER}}"

	migration "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration"
)

func main() {
	migration.DriverMain("application-migration.yml")
}
