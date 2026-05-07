# QA Test Matrix Template

## Metadata

- Feature area:
- Build/version:
- Device(s):
- OS version(s):
- Tester:
- Date:

## Matrix

| ID | Scenario | Preconditions | Steps | Expected Result | Actual Result | Status | Notes |
|---|---|---|---|---|---|---|---|
| QA-001 | Basic payload round-trip | App build installed | Encode then decode short ASCII payload | Decoded payload equals input |  | Not Run |  |
| QA-002 | Empty payload handling | App build installed | Encode/decode empty payload | System handles empty payload gracefully |  | Not Run |  |
| QA-003 | Max payload boundary | App build installed | Encode payload near configured max | Success at boundary without crash |  | Not Run |  |
| QA-004 | Invalid frame rejection | App build installed | Feed malformed frame to decoder | Decoder returns error and does not crash |  | Not Run |  |
| QA-005 | Cross-platform compatibility | iOS + Android builds available | Encode on iOS, decode on Android (and reverse) | Payload preserved across platforms |  | Not Run |  |

## Defect Log

| Defect ID | Related Test ID | Severity | Summary | Owner | Status |
|---|---|---|---|---|---|
|  |  |  |  |  |  |
