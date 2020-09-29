cmd
  - cli         <- entrypoint for CLI
  - nats        <- entrypoint for NATS client
app
  - app.go      <- DI container for deps. Ops for interaction layer dangles off this.
  - ops.go      <- all handlers for the interaction layer. Most of the domain logic lives in here. Break up as necessary.
  - policies.go <- logic for checking permissions
mapper          <- functions that map input types to app DTOs
  - protobuf.go
store
  - reader.go   <- interface for read operations. Accepts various arg shapes and returns DTO types defined in app.
  - writer.go   <- interface for write operations. Accepts various arg shapes and returns DTO types defined in app.
  - postgres    <- concrete reader/writer store types for postgres. Holds logic for creating DB connections and logic for querying


