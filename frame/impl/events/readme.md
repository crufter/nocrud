This implementation of the event interface (see /frame/interfaces) uses the Trigga pub/sub messaging server (see github.com/opesun/trigga).

This implementation doesn't panic with a nil connection, but returns errors on operations.