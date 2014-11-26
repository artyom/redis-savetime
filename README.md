**redis-savetime** — command to show duration of redis RDB saves. Looks for
`Background saving started` and `Background saving terminated` lines in log and
calculates time difference between them. Tested with RDB-only persistence
enabled.

Install:

	go get -v github.com/artyom/redis-savetime

Run:

	redis-savetime /var/log/redis.log
