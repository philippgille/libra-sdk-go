Changelog
=========

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

vNext
-----

- Added: Decoding of an account state and resource into proper objects (issue [#1](https://github.com/philippgille/libra-sdk-go/issues/1))
  - New struct `libra.AccountState` contains the blob (raw bytes) as well as a decoded `AccountResource`
  - New struct `libra.AccountResource` contains the account's balance, auth key, sent and received events, as well as sequence number
  - New function: `libra.FromAccountStateBlob(accountStateBlob []byte) (AccountState, error)` decodes an account state blob into a `libra.AccountState`
  - New function: `libra.FromAccountResourceBlob(accountResourceBlob []byte) (AccountResource, error)` decodes an account resource blob into a `libra.AccountResource`
- Improved: `Client.GetAccountState(accountAddr string)` now returns an object of the newly added `libra.AccountState`

### Breaking Changes

- The return type of `Client.GetAccountState(accountAddr string)` was changed from `([]byte, error)` to `(AccountState, error)` (for issue [#1](https://github.com/philippgille/libra-sdk-go/issues/1))

v0.1.0 (2019-06-23)
---------------------

Initial release.

### Features:

- Get account state as slice of bytes
- Send raw transaction
