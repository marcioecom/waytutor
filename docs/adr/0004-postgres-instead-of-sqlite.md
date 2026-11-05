# 0004. PostgreSQL on Neon instead of SQLite on a volume

v1 kept SQLite on a volume attached to one Fly machine. That couples every write
to one machine. A new machine loses the data, a deploy that replaces the machine
has to move the volume, and there is no backup unless we build one.

v2 uses PostgreSQL on Neon, reached through `pgx`. The service becomes
stateless: any number of machines can serve, a replaced machine loses nothing,
and backup, restore and branching come with the database. Neon's free tier
covers the expected usage, and Fly's volume disappears from the deploy.

**Consequences:** every query is a network round trip of single-digit
milliseconds in the same region, the deployment needs `DATABASE_URL`, and a
suspended Neon project adds latency to the first request. Migrations must be
additive and run before the new release starts. `cmd/migrate` does that as the
Fly release command.

The audio cache follows the same reasoning: production has no volume, so the
cache lives in R2. On a laptop the same code caches on disk, which is why the
cache is an interface and not a direct R2 call.

**Rejected:** SQLite with Litestream, because it adds moving parts and leaves
restore drills to us; Fly Postgres, because we would operate it; and
Turso/libSQL, because it adds a second dialect to learn for no gain at this
size.
