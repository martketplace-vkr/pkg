package postgres

const (
	queryTryToLockEvent = `select pg_try_advisory_lock(hashtext($1))`
	queryUnlockEvent    = `select pg_advisory_unlock(hashtext($1))`
)
