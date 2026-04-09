# Changelog

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
