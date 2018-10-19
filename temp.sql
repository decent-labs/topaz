
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name CHARACTER varying(255)
);

CREATE TABLE apps (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name CHARACTER varying(255),
    user_id INT NOT NULL,
    "interval" INTERVAL,
    last_flushed TIMESTAMP
);

CREATE TABLE flushes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    app_id INT NOT NULL,
    directory_hash CHARACTER varying(255) NOT NULL,
    "transaction" CHARACTER varying(255) NOT NULL
);

CREATE TABLE objects (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    hash CHARACTER varying(255) NOT NULL,
    app_id INT NOT NULL,
    flush_id INT
);


-- parker user
INSERT INTO users (name)
VALUES ('parker');

-- adam user
INSERT INTO users (name)
VALUES ('adam');

-- nate user
INSERT INTO users (name)
VALUES ('nate');


-- parker app 1
INSERT INTO apps (user_id, interval, last_flushed_id, created_at)
VALUES('app_parker_1', 'parker_id', '00:00:30', 'flush_parker_2', '2018-10-10 10:00:00')

-- parker app 2
INSERT INTO apps (id, user_id, interval, last_flushed_id, created_at)
VALUES('app_parker_2', 'parker_id', '00:00:45', 'flush_parker_2', '2018-10-10 10:00:00')


-- one flush for parker
INSERT INTO flushes (id, user_id, dir_hash, created_at)
VALUES ('flush_parker_1', 'parker_id', 'dirhash1', '2018-10-18 09:30:00');

-- another flush for parker
INSERT INTO flushes (id, user_id, dir_hash, created_at)
VALUES ('flush_parker_2', 'parker_id', 'dirhash2', '2018-10-18 09:30:30');

-- one flush for nate
INSERT INTO flushes (id, user_id, dir_hash, created_at)
VALUES ('flush_nate_1', 'nate_id', 'dirhash3', '2018-10-18 09:30:00');


-- flushed parker object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_flushed_parker_1', 'parker_id', 'flush_parker_1', 'hash1', '2018-10-18 09:29:51');

-- another flushed parker object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_flushed_parker_2', 'parker_id', 'flush_parker_2', 'hash2', '2018-10-18 09:30:10');

-- unflushed parker object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_unflushed_parker_1', 'parker_id', NULL, 'hash3', '2018-10-18 10:00:01');

-- unflushed parker object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_unflushed_parker_2', 'parker_id', NULL, 'hash4', '2018-10-18 10:00:02');


-- unflusehd adam object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_unflushed_adam_1', 'adam_id', NULL, 'hash5', '2018-10-18 10:00:03');

-- unflushed adam object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_unflushed_adam_2', 'adam_id', NULL, 'hash6', '2018-10-18 10:00:03');


-- flushed nate object
INSERT INTO objects (id, user_id, flush_id, hash, created_at)
VALUES ('obj_flushed_nate_1', 'nate_id', 'flush_nate_1', 'hash7', '2018-10-18 09:29:18');




/*
WITH almost 
     AS (WITH lastflush 
              AS (SELECT *, 
                         Row_number() 
                           OVER( 
                             partition BY user_id 
                             ORDER BY created_at DESC) rn 
                  FROM   flushes) 
         SELECT u.id 
          FROM   users u 
                 LEFT OUTER JOIN lastflush lf 
                              ON u.id = lf.user_id 
          WHERE  ( lf.rn = 1 
                    OR lf.rn IS NULL )) 
SELECT * 
FROM   almost a 
WHERE  EXISTS (SELECT o.id 
               FROM   objects o 
               WHERE  o.user_id = a.id 
                      AND o.flush_id IS NULL); 
*/


	 	select distinct a.id
		 	from apps a
		 		inner join objects o on o.app_id = a.app_id
		 	where ((a.last_flushed is null) or (now() - u.last_flushed >= a.interval))
		 		and (o.flush_id is null)
