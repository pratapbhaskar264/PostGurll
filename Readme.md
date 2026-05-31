# Will be available soon 

[Postgurrll Proxy] ════════ M2M Connection ════════> [Target Server]
   (Protocol Level)                                     (Protocol Level)
   Outbound Headers:                                    Inbound Headers:
   • Authorization: Bearer XYZ                          • Set-Cookie: session_id=123
   • Content-Type: application/json                     • Cache-Control: no-cache
          │                                                    │
          │                                                    │ (Captured by respo.Header)
          ▼                                                    ▼
 [Postgurrll Proxy] ══════ Telemetry Delivery ══════> [Your Postman Client]
                                                        (JSON Body Level)
                                                        {
                                                          "latency_ms": 145,
                                                          "headers": {
                                                             "Set-Cookie": ["session_id=123"],
                                                             "Cache-Control": ["no-cache"]
                                                          }
                                                        }
