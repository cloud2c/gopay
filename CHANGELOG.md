# Changelog

## v1.5.123

### Feat!: Require Response Signature Verification by Default

**BREAKING CHANGE.** `NewClientV3` now sets `verifyRequired` to `true`. A client
with neither `SetCert` nor `SetAliPayPublicKey` configured fails on the first
call instead of silently accepting unverified responses.

The old default existed to avoid breaking callers, but it made the failure mode
invisible: `aliPayPublicKey` is populated only by those two methods, and with
neither of them every response passed verification with a `nil` error. A payment
SDK that quietly stops verifying signatures is the wrong side to err on.

Correctly configured clients are unaffected. Callers who genuinely accept
unverified responses — temporary debugging, a sandbox app with no keys yet —
opt out explicitly:

```go
client.SetVerifyRequired(false) // responses are NOT verified
```

The error message now names both entry points rather than only `SetCert`, and
points at the opt-out.

## v1.5.122

Alipay V3 signing and verification, corrected against the official V3 protocol
docs. Three of the four items are cases where a valid response failed, or an
invalid one passed, without saying why.

### Fix: AES Response Ciphertext Stored on the Shared Client

With an AES key set, the V3 response signature is computed over the ciphertext,
but `doPost` decrypts the body in place. The ciphertext used to be held in
`ClientV3.rawBodyForSign` — written by `doPost`, read and cleared by
`autoVerifySignByCert` — with no synchronization.

`ClientV3` is a long-lived shared instance, so this was a data race (`go test
-race` reports two) and it also produced wrong results under concurrency:

- response B overwrites A's ciphertext, so A verifies against B's body
- A clears the field, so B falls back to the plaintext

Both surface as "signature verification failed" on valid payloads: fine at low
volume, random failures under load.

The ciphertext is now attached to the response's own `Request` context, scoped
per response. Supersedes the approach introduced in v1.5.119.

Also stopped swallowing the base64 decode error of the `alipay-signature`
header, which previously reported a missing header as a signature mismatch.

### Fix: Send alipay-root-cert-sn Header in Cert Mode

The [official signing rules](https://opendocs.alipay.com/open-v3/054q58) require
the Alipay root certificate serial number as its own header in certificate mode:

```
Authorization: ${signAlgorithm} ${authString},sign=${signature}
alipay-root-cert-sn: ${alipayRootCertSn}
```

V3 never sent it — `AliPayRootCertSN` was parsed by `SetCert` and then only used
by the V1-style gateway payloads in `payment_api.go`. Every V3 REST call in cert
mode went out missing a required header. Added to all six `do*` methods; key mode
leaves the field empty so the header is skipped.

`app_cert_sn` is *not* a header; it belongs inside `authString`, which
`authorization()` already did correctly.

### Fix: Report Cert Serial Mismatch Instead of a Signature Failure

Cert mode returns the serial Alipay signed with in the `alipay-sn` response
header, and the [verification rules](https://opendocs.alipay.com/open-v3/054d0z)
require merchants to compare it with their local certificate. Without the check,
a certificate rotation looks like "signature verification failed" on a valid
response, sending people to inspect keys and payloads when the fix is to replace
the certificate. Skipped when either side is empty (key mode has neither).

### Feat: Verify Responses in Key Mode

Per [接口加签方式](https://opendocs.alipay.com/open-v3/05419m) the two signing
modes are key mode (app private key, app public key, Alipay public key) and
certificate mode; one APPID may configure only one, and everything except
fund-transfer scenarios may use key mode.

V3 could not verify responses in key mode at all: `SetCert` was the only entry
point that populated `aliPayPublicKey`, so `autoVerifySignByCert` returned nil
and responses went unverified — silently.

```go
client, _ := alipay.NewClientV3(appId, privateKey, isProd)
client.SetAliPayPublicKey([]byte(alipayPublicKey)) // PEM or raw base64
client.SetVerifyRequired(true)                     // fail instead of skipping
```

`SetVerifyRequired` defaults to `false` so existing callers are unaffected;
`Clone` carries both settings over. Recommended for payment integrations.

New `alipay/v3/cert_mode_test.go` and `raw_body_test.go` cover the root cert
header in both modes, the `alipay-sn` mismatch error, the new setter and the
per-response ciphertext. Run with `-race`.

## v1.5.121

### Feat: Sync Alipay APIs from upstream

Added 4 new API modules with 25 new Alipay APIs:
- **Ad API** (`ad_api.go`): Conversion data upload, ad report query, promotion page management, task ad query
- **Fee API** (`fee_api.go`): Special fee rate application
- **Risk API** (`risk_api.go`): Consumer complaint processing, marketing risk identification, industry risk identification, content risk detection
- **Subscription API** (`subscription_api.go`): Product/price/customer/subscription CRUD operations

### Chore: Replace import paths from go-pay to cloud2c

Updated all documentation examples to use `github.com/cloud2c/gopay`.

## v1.5.120

### Feat: Add WithoutAES and Clone methods for Alipay V3 Client

Different Alipay APIs have different AES encryption requirements:
- Face verification/OCR APIs require AES encryption for `biz_content`
- Payment APIs (TradeCreate, TradePay, etc.) do NOT support content encryption, setting AES key causes parameter errors

Added thread-safe methods to handle mixed encryption scenarios:
- **`WithoutAES()`**: Returns a new Client instance without AES encryption, original client unaffected
- **`Clone()`**: Returns an independent Client copy without inheriting AES config

#### Usage

```go
// Option A: Separate clients (recommended)
faceClient := alipay.NewClientV3(...).SetAESKey(aesKey)
payClient := alipay.NewClientV3(...)

// Option B: Thread-safe temporary switch
client.SetAESKey(aesKey)
client.FaceVerificationInitialize(...)
client.WithoutAES().TradeCreate(...)  // original client still has AES
```

#### Files Changed

- `alipay/v3/client.go`

## v1.5.119

### Fix: Alipay V3 Signature Verification Failure with AES Encryption

When `aesKey` is set for content encryption, Alipay signs the **ciphertext** response, but the SDK was verifying the signature against the **decrypted plaintext**, causing `crypto/rsa: verification error`.

#### What Changed

- **`client.go`**: Added `rawBodyForSign` field to preserve the original response body before decryption.
- **`request.go`**: Save the raw ciphertext before AES decryption; removed dependency on the `alipay-content-encrypt` response header (some APIs return encrypted content without this header, causing decryption to be skipped).
- **`sign.go`**: `autoVerifySignByCert` now uses the raw ciphertext for signature verification when AES encryption is enabled, instead of the decrypted plaintext.

#### Correct Verification Flow

```
Alipay returns: ciphertext body + signature (signed over ciphertext)
    ↓
doPost: save ciphertext → rawBodyForSign, decrypt → plaintext body
    ↓
autoVerifySignByCert: verify signature using rawBodyForSign (ciphertext) ✅
    ↓
API method: json.Unmarshal(plaintext body) to parse business data
```

#### Files Changed

- `alipay/v3/client.go`
- `alipay/v3/request.go`
- `alipay/v3/sign.go`

## v1.5.118

### Refactor: Unified `do*` Method Signatures & Auto-Encryption

**Breaking Change** — All internal `do*` methods no longer accept `authorization` parameter; authorization is computed automatically.

#### Core Changes

- **`doPost` auto-encryption**: When `aesKey` is set, `doPost` automatically encrypts the request body and signs the ciphertext (encrypt-then-sign). No need to call `doPostWithEncrypt` anymore.
- **Removed `doPostWithEncrypt`**: Logic merged into `doPost`.
- **Unified `authorization` method**: Merged `authorizationWithEncryptBody` into `authorization` with an `encryptedBody` parameter.
- **`doProdPostFile` fully internalized**: File separation, `data` field encoding, and signing are now handled internally. Callers only pass the original `BodyMap` (including file fields) — no more `tempFile`/`signMap` boilerplate.

#### Method Signature Changes

| Method | Before | After |
|--------|--------|-------|
| `doPost` | `(ctx, bm, uri, authorization, aat)` | `(ctx, bm, uri, aat)` |
| `doGet` | `(ctx, uri, authorization, aat)` | `(ctx, uri, aat)` |
| `doPatch` | `(ctx, bm, uri, authorization, aat)` | `(ctx, bm, uri, aat)` |
| `doPut` | `(ctx, bm, uri, authorization, aat)` | `(ctx, bm, uri, aat)` |
| `doDelete` | `(ctx, bm, uri, authorization, aat)` | `(ctx, bm, uri, aat)` |
| `doProdPostFile` | `(ctx, bm, uri, authorization, aat)` or `(ctx, bm, uri, aat, signBm)` | `(ctx, bm, uri, aat)` |

#### Files Changed

- `alipay/v3/request.go` — `doPost` auto-encryption; `doProdPostFile` internalized; removed `doPostWithEncrypt`
- `alipay/v3/sign.go` — Unified `authorization` method
- `alipay/v3/face_verify_api.go` — Removed all if/else encryption branches
- 15 API files — Removed `authorization` declarations and parameter from `do*` calls
- 5 file upload APIs — Removed `tempFile`/`signMap`/`data` boilerplate (~28 lines each)

#### Migration

No external API changes. All public method signatures remain the same. This is only an internal refactoring.
