### 1. Protocol Documentation

Before making changes to message handling, **must** read:
`doc/message_stream.md`
This document tracks the current state of Bilibili's live protocol and our internal mapping.

### 2. Implementing New Events

When adding new Danmaku event types:

1. **Define the struct** in the relevant message package.
2. **Update the Mock Server:** You must update `pkg/app/mock.go` to include the new event. This ensures real-time UI testing and gRPC downstream verification work without a live connection.
3. **Update Protobuf:** If the event needs to be exposed via gRPC, update the `.proto` files and regenerate the code.
