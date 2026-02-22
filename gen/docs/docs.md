# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [application.proto](#application-proto)
    - [Application](#modelhawk-v1-Application)
  
- [info_service.proto](#info_service-proto)
    - [GiveToolInfoRequest](#modelhawk-v1-GiveToolInfoRequest)
    - [GiveToolInfoResponse](#modelhawk-v1-GiveToolInfoResponse)
    - [ToolArgInfo](#modelhawk-v1-ToolArgInfo)
  
    - [InfoService](#modelhawk-v1-InfoService)
  
- [message.proto](#message-proto)
    - [Message](#modelhawk-v1-Message)
  
- [notify_service.proto](#notify_service-proto)
    - [DidCallToolRequest](#modelhawk-v1-DidCallToolRequest)
    - [DidCallToolRequest.ArgsEntry](#modelhawk-v1-DidCallToolRequest-ArgsEntry)
    - [WillCallToolRequest](#modelhawk-v1-WillCallToolRequest)
    - [WillCallToolRequest.ArgsEntry](#modelhawk-v1-WillCallToolRequest-ArgsEntry)
  
    - [NotifyService](#modelhawk-v1-NotifyService)
  
- [permission_service.proto](#permission_service-proto)
    - [WantsToCallToolRequest](#modelhawk-v1-WantsToCallToolRequest)
    - [WantsToCallToolRequest.ArgsEntry](#modelhawk-v1-WantsToCallToolRequest-ArgsEntry)
    - [WantsToCallToolResponse](#modelhawk-v1-WantsToCallToolResponse)
  
    - [PermissionService](#modelhawk-v1-PermissionService)
  
- [service_status_response.proto](#service_status_response-proto)
    - [ServiceStatusResponse](#modelhawk-v1-ServiceStatusResponse)
  
    - [ServiceStatusResponse.Result](#modelhawk-v1-ServiceStatusResponse-Result)
  
- [Scalar Value Types](#scalar-value-types)



<a name="application-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## application.proto



<a name="modelhawk-v1-Application"></a>

### Application
An Application identifies a particular deployment of an AI app. Examples:
- A single chat session
- A deployment of claude on a particular computer (encompassing all its sessions)

The appropriate scope really depends on how you want to use the security app that is monitoring the AI app.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) | optional |  |





 

 

 

 



<a name="info_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## info_service.proto



<a name="modelhawk-v1-GiveToolInfoRequest"></a>

### GiveToolInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [Application](#modelhawk-v1-Application) | optional |  |
| name | [string](#string) | optional |  |
| desc | [string](#string) | optional |  |
| args | [ToolArgInfo](#modelhawk-v1-ToolArgInfo) | repeated |  |






<a name="modelhawk-v1-GiveToolInfoResponse"></a>

### GiveToolInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [int32](#int32) | optional |  |
| msg | [string](#string) | optional |  |






<a name="modelhawk-v1-ToolArgInfo"></a>

### ToolArgInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| type | [string](#string) | optional |  |
| desc | [string](#string) | optional |  |





 

 

 


<a name="modelhawk-v1-InfoService"></a>

### InfoService
InfoService is a service for telling the security app about configuration in the AI app that&#39;s being monitored.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GiveToolInfo | [GiveToolInfoRequest](#modelhawk-v1-GiveToolInfoRequest) | [ServiceStatusResponse](#modelhawk-v1-ServiceStatusResponse) | Tell the security tool about a tool that is available in a particular context. |

 



<a name="message-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## message.proto



<a name="modelhawk-v1-Message"></a>

### Message
A message sent from or to an AI model.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [string](#string) | optional |  |
| content | [string](#string) | optional |  |





 

 

 

 



<a name="notify_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## notify_service.proto



<a name="modelhawk-v1-DidCallToolRequest"></a>

### DidCallToolRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [Application](#modelhawk-v1-Application) | optional |  |
| tool_name | [string](#string) | optional |  |
| args | [DidCallToolRequest.ArgsEntry](#modelhawk-v1-DidCallToolRequest-ArgsEntry) | repeated |  |
| result | [string](#string) | optional |  |
| last_messages | [Message](#modelhawk-v1-Message) | repeated |  |






<a name="modelhawk-v1-DidCallToolRequest-ArgsEntry"></a>

### DidCallToolRequest.ArgsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="modelhawk-v1-WillCallToolRequest"></a>

### WillCallToolRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [Application](#modelhawk-v1-Application) | optional |  |
| tool_name | [string](#string) | optional |  |
| args | [WillCallToolRequest.ArgsEntry](#modelhawk-v1-WillCallToolRequest-ArgsEntry) | repeated |  |
| last_messages | [Message](#modelhawk-v1-Message) | repeated |  |






<a name="modelhawk-v1-WillCallToolRequest-ArgsEntry"></a>

### WillCallToolRequest.ArgsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 

 

 


<a name="modelhawk-v1-NotifyService"></a>

### NotifyService
NotifyService is a service for telling the security app about stuff that has happened in the AI app that&#39;s being monitored.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| WillCallTool | [WillCallToolRequest](#modelhawk-v1-WillCallToolRequest) | [ServiceStatusResponse](#modelhawk-v1-ServiceStatusResponse) | WillCallTool can be called by the AI app to tell the security app that the AI model will call a tool. |
| DidCallTool | [DidCallToolRequest](#modelhawk-v1-DidCallToolRequest) | [ServiceStatusResponse](#modelhawk-v1-ServiceStatusResponse) | DidCallTool can be called by the AI app to tell the security app that the AI model did call a tool. |

 



<a name="permission_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## permission_service.proto



<a name="modelhawk-v1-WantsToCallToolRequest"></a>

### WantsToCallToolRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [Application](#modelhawk-v1-Application) | optional |  |
| tool_name | [string](#string) | optional |  |
| args | [WantsToCallToolRequest.ArgsEntry](#modelhawk-v1-WantsToCallToolRequest-ArgsEntry) | repeated |  |
| last_messages | [Message](#modelhawk-v1-Message) | repeated |  |






<a name="modelhawk-v1-WantsToCallToolRequest-ArgsEntry"></a>

### WantsToCallToolRequest.ArgsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="modelhawk-v1-WantsToCallToolResponse"></a>

### WantsToCallToolResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permitted | [bool](#bool) | optional |  |
| details | [string](#string) | optional |  |





 

 

 


<a name="modelhawk-v1-PermissionService"></a>

### PermissionService
PermissionService is a service that AI apps can use to ask the security app for permission to do stuff.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| WantsToCallTool | [WantsToCallToolRequest](#modelhawk-v1-WantsToCallToolRequest) | [WantsToCallToolResponse](#modelhawk-v1-WantsToCallToolResponse) | The AI app wants to call a tool. The security app can approve or deny it. |

 



<a name="service_status_response-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service_status_response.proto



<a name="modelhawk-v1-ServiceStatusResponse"></a>

### ServiceStatusResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [ServiceStatusResponse.Result](#modelhawk-v1-ServiceStatusResponse-Result) | optional |  |
| msg | [string](#string) | optional |  |





 


<a name="modelhawk-v1-ServiceStatusResponse-Result"></a>

### ServiceStatusResponse.Result


| Name | Number | Description |
| ---- | ------ | ----------- |
| RESULT_OK | 0 |  |
| RESULT_ERROR | 1 |  |


 

 

 



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

