# 0004 — PostgreSQL on Neon instead of SQLite on a volume

v1 kept SQLite on a volume attached to one Fly machine. That couples every
write to one machine: a new machine loses the data, a deploy that replaces the
machine needs the volume to move, and there is no backup unless we build one.

v2 uses PostgreSQL on Neon, reached through `pgx`. The service becomes
stateless: any number of machines can serve, a replaced machine loses nothing,
and backup, restore and branching come with the database. Neon's free tier
covers the expected usage, and Fly's volume disappears from the deploy.

**Consequences:** every query is a network round trip (single-digit
milliseconds in the same region), the deployment needs `DATABASE_URL`, and a
suspended Neon project adds latency to the first request. Migrations must be
additive and run before the new release starts, which `cmd/migrate` does as the
Fly release command.

**Rejected:** SQLite with Litestream (more moving parts, restore drills on us),
Fly Postgres (we operate it), Turso/libSQL (a second dialect to learn for no
gain at this size).
