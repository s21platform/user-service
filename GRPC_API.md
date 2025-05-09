# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/user.proto](#api_user-proto)
    - [GetLoginByUUIDIn](#-GetLoginByUUIDIn)
    - [GetLoginByUUIDOut](#-GetLoginByUUIDOut)
    - [GetOs](#-GetOs)
    - [GetUserByLoginIn](#-GetUserByLoginIn)
    - [GetUserByLoginOut](#-GetUserByLoginOut)
    - [GetUserInfoByUUIDIn](#-GetUserInfoByUUIDIn)
    - [GetUserInfoByUUIDOut](#-GetUserInfoByUUIDOut)
    - [GetUserWithOffsetIn](#-GetUserWithOffsetIn)
    - [GetUserWithOffsetOut](#-GetUserWithOffsetOut)
    - [GetUserWithOffsetOutAll](#-GetUserWithOffsetOutAll)
    - [GetUsersByUUIDIn](#-GetUsersByUUIDIn)
    - [GetUsersByUUIDOut](#-GetUsersByUUIDOut)
    - [IsUserExistByUUIDIn](#-IsUserExistByUUIDIn)
    - [IsUserExistByUUIDOut](#-IsUserExistByUUIDOut)
    - [UpdateProfileFormIn](#-UpdateProfileFormIn)
    - [UpdateProfileFormOut](#-UpdateProfileFormOut)
    - [UpdateProfileIn](#-UpdateProfileIn)
    - [UpdateProfileOut](#-UpdateProfileOut)
    - [UpdateProfileTestIn](#-UpdateProfileTestIn)
    - [UpdateProfileTestOut](#-UpdateProfileTestOut)
    - [User](#-User)
    - [UserInfoMin](#-UserInfoMin)
    - [UsersUUID](#-UsersUUID)
  
    - [UserService](#-UserService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/user.proto



<a name="-GetLoginByUUIDIn"></a>

### GetLoginByUUIDIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="-GetLoginByUUIDOut"></a>

### GetLoginByUUIDOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |






<a name="-GetOs"></a>

### GetOs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| label | [string](#string) |  |  |






<a name="-GetUserByLoginIn"></a>

### GetUserByLoginIn
Data in request or getting uuid by login. If User doesnt exist - user will be creating


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  | Email of target user |






<a name="-GetUserByLoginOut"></a>

### GetUserByLoginOut
Message for response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | UUID of user |
| isNewUser | [bool](#bool) |  | Flag for indicate of new user |






<a name="-GetUserInfoByUUIDIn"></a>

### GetUserInfoByUUIDIn
Request data fo getting user info (for initiator page)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | UUID for target user |






<a name="-GetUserInfoByUUIDOut"></a>

### GetUserInfoByUUIDOut
Response data for initiator page


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nickname | [string](#string) |  |  |
| avatar | [string](#string) |  |  |
| name | [string](#string) | optional |  |
| surname | [string](#string) | optional |  |
| birthdate | [string](#string) | optional |  |
| phone | [string](#string) | optional |  |
| city | [string](#string) | optional |  |
| telegram | [string](#string) | optional |  |
| git | [string](#string) | optional |  |
| os | [GetOs](#GetOs) | optional |  |
| work | [string](#string) | optional |  |
| university | [string](#string) | optional |  |
| skills | [string](#string) | repeated |  |
| hobbies | [string](#string) | repeated |  |
| uuid | [string](#string) | optional |  |






<a name="-GetUserWithOffsetIn"></a>

### GetUserWithOffsetIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int64](#int64) |  |  |
| offset | [int64](#int64) |  |  |
| nickname | [string](#string) |  |  |






<a name="-GetUserWithOffsetOut"></a>

### GetUserWithOffsetOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#User) | repeated |  |
| total | [int64](#int64) |  |  |






<a name="-GetUserWithOffsetOutAll"></a>

### GetUserWithOffsetOutAll



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [GetUserInfoByUUIDOut](#GetUserInfoByUUIDOut) | repeated |  |
| total | [int64](#int64) |  |  |






<a name="-GetUsersByUUIDIn"></a>

### GetUsersByUUIDIn
Request message for getting multiple users by their UUIDs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users_uuid | [UsersUUID](#UsersUUID) | repeated |  |






<a name="-GetUsersByUUIDOut"></a>

### GetUsersByUUIDOut
Response message containing minimal user information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users_info | [UserInfoMin](#UserInfoMin) | repeated |  |






<a name="-IsUserExistByUUIDIn"></a>

### IsUserExistByUUIDIn
Message for request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | UUID for target user |






<a name="-IsUserExistByUUIDOut"></a>

### IsUserExistByUUIDOut
Message for response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| isExist | [bool](#bool) |  | Flag of indicate user exist |






<a name="-UpdateProfileFormIn"></a>

### UpdateProfileFormIn







<a name="-UpdateProfileFormOut"></a>

### UpdateProfileFormOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="-UpdateProfileIn"></a>

### UpdateProfileIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| birthday | [string](#string) |  |  |
| telegram | [string](#string) |  |  |
| github | [string](#string) |  |  |
| os_id | [int64](#int64) |  |  |






<a name="-UpdateProfileOut"></a>

### UpdateProfileOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [bool](#bool) |  |  |






<a name="-UpdateProfileTestIn"></a>

### UpdateProfileTestIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="-UpdateProfileTestOut"></a>

### UpdateProfileTestOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nickname | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| avatar_link | [string](#string) |  |  |
| name | [string](#string) |  |  |
| surname | [string](#string) |  |  |






<a name="-UserInfoMin"></a>

### UserInfoMin
Min user information structure


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| login | [string](#string) |  |  |
| last_avatar | [string](#string) |  |  |
| name | [string](#string) |  |  |
| surname | [string](#string) |  |  |






<a name="-UsersUUID"></a>

### UsersUUID
Message for UsersUUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |





 

 

 


<a name="-UserService"></a>

### UserService
Service for friends

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetUserByLogin | [.GetUserByLoginIn](#GetUserByLoginIn) | [.GetUserByLoginOut](#GetUserByLoginOut) | Add friends method |
| IsUserExistByUUID | [.IsUserExistByUUIDIn](#IsUserExistByUUIDIn) | [.IsUserExistByUUIDOut](#IsUserExistByUUIDOut) |  |
| GetUserInfoByUUID | [.GetUserInfoByUUIDIn](#GetUserInfoByUUIDIn) | [.GetUserInfoByUUIDOut](#GetUserInfoByUUIDOut) |  |
| GetLoginByUUID | [.GetLoginByUUIDIn](#GetLoginByUUIDIn) | [.GetLoginByUUIDOut](#GetLoginByUUIDOut) |  |
| GetUserWithOffset | [.GetUserWithOffsetIn](#GetUserWithOffsetIn) | [.GetUserWithOffsetOut](#GetUserWithOffsetOut) |  |
| UpdateProfile | [.UpdateProfileIn](#UpdateProfileIn) | [.UpdateProfileOut](#UpdateProfileOut) |  |
| GetUsersByUUID | [.GetUsersByUUIDIn](#GetUsersByUUIDIn) | [.GetUsersByUUIDOut](#GetUsersByUUIDOut) |  |
| UpdateProfileTest | [.UpdateProfileTestIn](#UpdateProfileTestIn) | [.UpdateProfileTestOut](#UpdateProfileTestOut) |  |
| UpdateProfileForm | [.UpdateProfileFormIn](#UpdateProfileFormIn) | [.UpdateProfileFormOut](#UpdateProfileFormOut) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

