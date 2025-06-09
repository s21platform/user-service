# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/user.proto](#api_user-proto)
    - [CreatePostIn](#-CreatePostIn)
    - [CreatePostOut](#-CreatePostOut)
    - [CreateUserIn](#-CreateUserIn)
    - [CreateUserOut](#-CreateUserOut)
    - [EmptyFriends](#-EmptyFriends)
    - [GetCountFriendsOut](#-GetCountFriendsOut)
    - [GetLoginByUUIDIn](#-GetLoginByUUIDIn)
    - [GetLoginByUUIDOut](#-GetLoginByUUIDOut)
    - [GetOs](#-GetOs)
    - [GetPeerFollowIn](#-GetPeerFollowIn)
    - [GetPeerFollowOut](#-GetPeerFollowOut)
    - [GetPostsByIdsIn](#-GetPostsByIdsIn)
    - [GetPostsByIdsOut](#-GetPostsByIdsOut)
    - [GetUserByLoginIn](#-GetUserByLoginIn)
    - [GetUserByLoginOut](#-GetUserByLoginOut)
    - [GetUserInfoByUUIDIn](#-GetUserInfoByUUIDIn)
    - [GetUserInfoByUUIDOut](#-GetUserInfoByUUIDOut)
    - [GetUserWithOffsetIn](#-GetUserWithOffsetIn)
    - [GetUserWithOffsetOut](#-GetUserWithOffsetOut)
    - [GetUserWithOffsetOutAll](#-GetUserWithOffsetOutAll)
    - [GetUsersByUUIDIn](#-GetUsersByUUIDIn)
    - [GetUsersByUUIDOut](#-GetUsersByUUIDOut)
    - [GetWhoFollowPeerIn](#-GetWhoFollowPeerIn)
    - [GetWhoFollowPeerOut](#-GetWhoFollowPeerOut)
    - [IsUserExistByUUIDIn](#-IsUserExistByUUIDIn)
    - [IsUserExistByUUIDOut](#-IsUserExistByUUIDOut)
    - [Peer](#-Peer)
    - [PostInfo](#-PostInfo)
    - [RemoveFriendsIn](#-RemoveFriendsIn)
    - [RemoveFriendsOut](#-RemoveFriendsOut)
    - [SetFriendsIn](#-SetFriendsIn)
    - [SetFriendsOut](#-SetFriendsOut)
    - [UpdateProfileIn](#-UpdateProfileIn)
    - [UpdateProfileOut](#-UpdateProfileOut)
    - [User](#-User)
    - [UserCreatedMessage](#-UserCreatedMessage)
    - [UserInfoMin](#-UserInfoMin)
    - [UserNicknameUpdated](#-UserNicknameUpdated)
    - [UsersUUID](#-UsersUUID)
  
    - [UserService](#-UserService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/user.proto



<a name="-CreatePostIn"></a>

### CreatePostIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| content | [string](#string) |  |  |






<a name="-CreatePostOut"></a>

### CreatePostOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| post_uuid | [string](#string) |  |  |






<a name="-CreateUserIn"></a>

### CreateUserIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  |  |






<a name="-CreateUserOut"></a>

### CreateUserOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nickname | [string](#string) |  |  |
| user_uuid | [string](#string) |  |  |






<a name="-EmptyFriends"></a>

### EmptyFriends







<a name="-GetCountFriendsOut"></a>

### GetCountFriendsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription | [int64](#int64) |  |  |
| subscribers | [int64](#int64) |  |  |






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






<a name="-GetPeerFollowIn"></a>

### GetPeerFollowIn
Request for subscription


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | Peer uuid |






<a name="-GetPeerFollowOut"></a>

### GetPeerFollowOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription | [Peer](#Peer) | repeated |  |






<a name="-GetPostsByIdsIn"></a>

### GetPostsByIdsIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| post_uuids | [string](#string) | repeated |  |






<a name="-GetPostsByIdsOut"></a>

### GetPostsByIdsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| posts | [PostInfo](#PostInfo) | repeated |  |






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






<a name="-GetWhoFollowPeerIn"></a>

### GetWhoFollowPeerIn
Request for subscribers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | Peer uuid |






<a name="-GetWhoFollowPeerOut"></a>

### GetWhoFollowPeerOut
Response subscribers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscribers | [Peer](#Peer) | repeated | Result of the operation |






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






<a name="-Peer"></a>

### Peer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | Peer uuid |






<a name="-PostInfo"></a>

### PostInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| post_uuid | [int64](#int64) |  |  |
| nickname | [string](#string) |  |  |
| full_name | [string](#string) |  |  |
| avatar_link | [string](#string) |  |  |
| content | [string](#string) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| is_edited | [bool](#bool) |  |  |






<a name="-RemoveFriendsIn"></a>

### RemoveFriendsIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer | [string](#string) |  |  |






<a name="-RemoveFriendsOut"></a>

### RemoveFriendsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="-SetFriendsIn"></a>

### SetFriendsIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer | [string](#string) |  |  |






<a name="-SetFriendsOut"></a>

### SetFriendsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






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






<a name="-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nickname | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| avatar_link | [string](#string) |  |  |
| name | [string](#string) |  |  |
| surname | [string](#string) |  |  |






<a name="-UserCreatedMessage"></a>

### UserCreatedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_uuid | [string](#string) |  |  |






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






<a name="-UserNicknameUpdated"></a>

### UserNicknameUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_uuid | [string](#string) |  |  |
| nickname | [string](#string) |  |  |






<a name="-UsersUUID"></a>

### UsersUUID
Message for UsersUUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |





 

 

 


<a name="-UserService"></a>

### UserService
Service for user info

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetUserByLogin | [.GetUserByLoginIn](#GetUserByLoginIn) | [.GetUserByLoginOut](#GetUserByLoginOut) |  |
| IsUserExistByUUID | [.IsUserExistByUUIDIn](#IsUserExistByUUIDIn) | [.IsUserExistByUUIDOut](#IsUserExistByUUIDOut) |  |
| GetUserInfoByUUID | [.GetUserInfoByUUIDIn](#GetUserInfoByUUIDIn) | [.GetUserInfoByUUIDOut](#GetUserInfoByUUIDOut) |  |
| GetLoginByUUID | [.GetLoginByUUIDIn](#GetLoginByUUIDIn) | [.GetLoginByUUIDOut](#GetLoginByUUIDOut) |  |
| GetUserWithOffset | [.GetUserWithOffsetIn](#GetUserWithOffsetIn) | [.GetUserWithOffsetOut](#GetUserWithOffsetOut) |  |
| UpdateProfile | [.UpdateProfileIn](#UpdateProfileIn) | [.UpdateProfileOut](#UpdateProfileOut) |  |
| GetUsersByUUID | [.GetUsersByUUIDIn](#GetUsersByUUIDIn) | [.GetUsersByUUIDOut](#GetUsersByUUIDOut) |  |
| CreateUser | [.CreateUserIn](#CreateUserIn) | [.CreateUserOut](#CreateUserOut) |  |
| SetFriends | [.SetFriendsIn](#SetFriendsIn) | [.SetFriendsOut](#SetFriendsOut) |  |
| RemoveFriends | [.RemoveFriendsIn](#RemoveFriendsIn) | [.RemoveFriendsOut](#RemoveFriendsOut) |  |
| GetCountFriends | [.EmptyFriends](#EmptyFriends) | [.GetCountFriendsOut](#GetCountFriendsOut) |  |
| GetPeerFollow | [.GetPeerFollowIn](#GetPeerFollowIn) | [.GetPeerFollowOut](#GetPeerFollowOut) |  |
| GetWhoFollowPeer | [.GetWhoFollowPeerIn](#GetWhoFollowPeerIn) | [.GetWhoFollowPeerOut](#GetWhoFollowPeerOut) |  |
| CreatePost | [.CreatePostIn](#CreatePostIn) | [.CreatePostOut](#CreatePostOut) |  |
| GetPostsByIds | [.GetPostsByIdsIn](#GetPostsByIdsIn) | [.GetPostsByIdsOut](#GetPostsByIdsOut) |  |

 



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

