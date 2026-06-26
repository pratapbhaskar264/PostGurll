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




                                                        ---
When you pass that empty string directly into Go’s http.NewRequest(method, url, body), the function is designed to be highly resilient. Under the hood, Go checks if the method string is empty, and if it is, it automatically defaults to "GET".\



