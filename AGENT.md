## Style Guide

* **Style:** Rose-Pine as the base style for GUI.

## Project Structure

```text
├── build/             # Binary build artifacts
├── doc/               # Protocol documentation (message_stream.md)
|--- api/
|   ├── grpc/          # gRPC service definitions and generated code
├── pkg/
│   ├── app/           # Core logic and Mock implementations
│   ├── bilibili/      # Bilibili authenticated connection handling
|   ├── live/          # Bilibili Live API connection handling, including data parsers
|   ├── server/        # Server logic
|   └── ui/            # UI logic
├── frontend/          # Svelte source code
└── main.go            # Entry point

```

## Code Standards

* **Safety First:** Handle all pointer references and nil-checks, especially during JSON/Protobuf unmarshaling.
* **Readability:** Use clear naming conventions that reflect the Bilibili API terminology.
* **Scalability:** Keep the message parsing logic decoupled from the UI/gRPC dispatching logic.
