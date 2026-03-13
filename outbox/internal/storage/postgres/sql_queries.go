package postgres

const (
	queryCreateEvent = `
		insert into outbox (
		                    topic, 
		                    payload, 
		                    trace_carrier, 
		                    key, 
		                    created_at,
		                    event_type,
		                    idempotency_key,
		                    attempts,
		                    last_error,
		                    locked_until,
		                    state::state
        ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		returning id;
`
	queryInsertBatch = `
	   insert into outbox (
						  topic,
						  payload,
						  trace_carrier,
						  key,
						  created_at,
						  event_type,
	                      idempotency_key,
						   attempts,
						   last_error,
	                       state
      ) values (
				unnest($1::text[]),
				unnest($2::jsonb[]),
		        unnest($3::jsonb[]),
		        unnest($4::text[]),
		        unnest($5::timestamp[]),
                unnest($6::text[]),
                unnest($7::text[]),
                unnest($8::int[]),
                unnest($9::text[]),
                unnest($10::state[])
	  )	returning id;
`

	queryFetchUnprocessedEventsUnordered = `
        with locked as (
				select
					id,
					topic,
					payload,
					trace_carrier,
					created_at,
					sent_at,
					"key",
					idempotency_key,
					attempts,
					last_error,
					event_type
				from outbox
				where sent_at is null
					  and (locked_until is null or locked_until <= now())
					  and attempts < $3
 					  and state = 'new'
				order by created_at asc
				limit $2
				for update skip locked
		) update outbox o
          set locked_until = $1,
				state = 'processing'
		from locked
		where o.id = locked.id
		  and o.sent_at IS NULL
		  and o.state = 'new'
		  and (o.locked_until IS NULL OR o.locked_until <= now())
		  and o.attempts < $2
		returning
			o.id,
			o.topic,
			o.payload,
			o.trace_carrier,
			o.created_at,
			o.sent_at,
			o."key",
			o.idempotency_key,
			o.attempts,
			o.last_error,
			o.event_type;
`
	queryMarkEventAsPublished = `
	update outbox
	set sent_at = now(),
	    partition = $2,
	    "offset" = $3,
	    locked_until = null,
	    state = 'done',
		attempts = attempts + 1
	where id = $1;`

	queryMarkEventsAsSkipped = `
		update outbox
		SET locked_until = NULL,
		    state = 'new'
-- 		update only if id is in the list AND previous state is processing
		WHERE id = any($1) and state = 'processing'
	`

	queryUnlockEvents = `
		UPDATE outbox
		SET locked_until = NULL,
		    last_error = $2,
		    state = 'error',
			attempts = attempts + 1
		WHERE id = $1;
	`
	queryFetchUnprocessedEventsWithOrder = `
		with head as (
			select distinct on (o."key")
				o.id,
				o."key",
				o.created_at
			from outbox o
			where o.sent_at is null
			AND o.attempts < $3
			AND state <> 'done' AND state <> 'error'
			ORDER BY o."key", o.created_at, o.id
		),
		locked as (
			select
				o.id,
				o.topic,
				o.payload,
				o.trace_carrier,
				o.created_at,
				o.sent_at,
				o."key",
				o.idempotency_key,
				o.event_type,
				o.last_error,
				o.attempts
			from outbox o
			join head h on h.id = o.id
			where o.state = 'new'
			and (o.locked_until is null or o.locked_until <= now())
			for update skip locked
			LIMIT $2
		),
		updated AS (
			update outbox o
			set state = 'processing',
				locked_until = $1
			from locked
			where o.id = locked.id
			and o.sent_at is null
			and o.state = 'new'
			and (o.locked_until is null or o.locked_until <= now())
			and o.attempts < $3
			returning
				o.id,
				o.topic,
				o.payload,
				o.trace_carrier,
				o.created_at,
				o.sent_at,
				o."key",
				o.idempotency_key,
				o.attempts,
				o.last_error,
				o.event_type,
				o.state
		)
		select *
		from updated
		order by created_at;
	`
)
