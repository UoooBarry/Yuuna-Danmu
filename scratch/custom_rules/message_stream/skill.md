### 1. Protocol Documentation

Before making changes to message handling, **must** read:
`doc/message_stream.md`
This document tracks the current state of Bilibili's live protocol and our internal mapping.

### 2. Implementing New Events

When adding new Danmaku event types:

1. **Define the struct** in the relevant message package `pkg/live/type`.
2. **Update the parser** in `pkg/app/parser.go` parse the raw data to the correct struct.
3. **Dispatch the event** in `pkg/app/app.go` dispatch the event to the correct handler.
4. **Handle the frontend** in `pkg/ui` handle the event in the correct UI.
5. **Update Protobuf:** New events need to be exposed via gRPC, gRPC files are located in `/api/grpc`. Update the `.proto` files and regenerate the code.
6. **Update the Mock Server:** You must update `pkg/app/mock.go` to include the new event. This ensures real-time UI testing and gRPC downstream verification work without a live connection.
