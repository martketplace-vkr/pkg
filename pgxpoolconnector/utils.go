package pgxpoolconnctor

import "strings"

func getDbNameFromDSN(dsn string) string {
	dsnParts := strings.Split(dsn, " ")
	for _, part := range dsnParts {
		if strings.HasPrefix(part, "dbname=") {
			dbNameParts := strings.Split(part, "=")
			return dbNameParts[1]
		}
	}

	return ""
}
