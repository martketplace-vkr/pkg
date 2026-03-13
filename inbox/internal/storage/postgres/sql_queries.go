package postgres

const (
	queryCreateEvent = `
		insert into inbox (
		                    topic, 
		                    payload, 
		                    trace_carrier, 
		                    key,
							event_type,
		                    created_at,
		                    idempotency_key,
		                    error,
		                    state
        ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning id;
`
	queryInsertBatch = `
	   insert into inbox (
						  topic,
						  payload,
						  trace_carrier,
						  key,
						  event_type,
						  created_at,
	                      idempotency_key,
	                      error,
	                      attempts,
	                      state
      ) values (
				unnest($1::text[]),
				unnest($2::jsonb[]),
		        unnest($3::jsonb[]),
		        unnest($4::text[]),
		        unnest($5::text[]),
		        unnest($6::timestamp[]),
                unnest($7::text[]),
                unnest($8::text[]),
                unnest($9::int[]),
                unnest($10::state[])
	  )	returning id;
`

	// now we are using CAS algorithm (Compare and swap)
	// select events with state 'new' and swap it to 'processing' if the state is still 'new'
	// otherwise we will not select that field.
	queryFetchUnprocessedEventsWithOrder = `
			WITH head AS (
				SELECT DISTINCT ON (i.key)
					i.id,
					i.key,
					i.created_at
				FROM inbox i
				WHERE i.processed_at IS NULL
				  AND i.attempts < $2
				  and state <> 'done' and state <> 'error'
				ORDER BY i.key, i.created_at, i.id
			),
			locked AS (
				SELECT
					i.id,
					i.topic,
					i.payload,
					i.trace_carrier,
					i.created_at,
					i.processed_at,
					i.key,
					i.event_type,
					i.idempotency_key,
					i.error,
					i.attempts
				FROM inbox i
				JOIN head h ON h.id = i.id
				WHERE i.state = 'new'
				  AND (i.locked_until IS NULL OR i.locked_until <= now())
				FOR UPDATE SKIP LOCKED
				LIMIT $3
			),
			updated AS (
				UPDATE inbox i
				SET state = 'processing',
					locked_until = $1
				FROM locked
				WHERE i.id = locked.id
					  and i.processed_at IS NULL
					  AND i.state = 'new'
					  AND (i.locked_until IS NULL OR i.locked_until <= now())
					  AND i.attempts < $2
				RETURNING
					i.id,
					i.topic,
					i.payload,
					i.trace_carrier,
					i.created_at,
					i.processed_at,
					i.key,
					i.event_type,
					i.idempotency_key,
					i.state
			)
			SELECT *
			FROM updated
			ORDER BY created_at;
`

	queryFetchUnprocessedEventsUnordered = `
 			with locked as (
				select
					id,
					topic,
					payload,
					trace_carrier,
					created_at,
					processed_at,
					"key",
					event_type,
					idempotency_key,
					error,
					attempts
				from inbox
				where processed_at is null
				  and state = 'new'
				  and (locked_until is null or locked_until <= now())
				  and attempts < $2
				order by created_at
				limit $3
				for update skip locked
			), updated as (
				update inbox i
				set state = 'processing',
					locked_until = $1
				from locked
				where i.id = locked.id
					  and i.state = 'new'
					  AND (i.locked_until IS NULL OR i.locked_until <= now())
				returning
					i.id,
					i.topic,
					i.payload,
					i.trace_carrier,
					i.created_at,
					i.processed_at,
					i."key",
					i.event_type,
					i.idempotency_key,
					i.state
			) select * from updated
			  order by updated.created_at
	`

	queryFetchExpiredEvents = `
		select id
			 , topic
			 , payload
			 , trace_carrier
			 , created_at
			 , processed_at
			 , key
			 , event_type
			 , idempotency_key
			 , error
		from inbox
		where processed_at < now() - $1 * interval '1 day'
		order by id
		limit $2 for update skip locked;
`

	queryDeleteEventsByIDs = `
		delete from inbox
		where id = any($1);
`

	queryMarkEventAsPublished = `
	update inbox
	set processed_at = now(),
	    locked_until = null,
	    state = $2,
	    attempts = attempts + 1
	where id = $1;`

	queryMarkEventsAsFailed = `
		update inbox
		SET locked_until = NULL,
		    error = $2,
		    state = 'error',
		    attempts = attempts + 1
		WHERE id = $1
	`
	queryMarkEventsAsSkipped = `
		update inbox
		SET locked_until = NULL,
		    state = 'new'
-- 		update only if id is in the list AND previous state is processing
		WHERE id = any($1) and state = 'processing'
	`
)
