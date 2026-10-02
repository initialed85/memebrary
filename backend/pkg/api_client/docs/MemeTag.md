# MemeTag

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**MemeId** | Pointer to **string** |  | [optional] 
**MemeIdObject** | Pointer to [**Meme**](Meme.md) |  | [optional] 
**TagId** | Pointer to **string** |  | [optional] 
**TagIdObject** | Pointer to [**Tag**](Tag.md) |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewMemeTag

`func NewMemeTag() *MemeTag`

NewMemeTag instantiates a new MemeTag object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemeTagWithDefaults

`func NewMemeTagWithDefaults() *MemeTag`

NewMemeTagWithDefaults instantiates a new MemeTag object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *MemeTag) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *MemeTag) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *MemeTag) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *MemeTag) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDeletedAt

`func (o *MemeTag) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *MemeTag) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *MemeTag) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *MemeTag) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.

### GetId

`func (o *MemeTag) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MemeTag) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MemeTag) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *MemeTag) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMemeId

`func (o *MemeTag) GetMemeId() string`

GetMemeId returns the MemeId field if non-nil, zero value otherwise.

### GetMemeIdOk

`func (o *MemeTag) GetMemeIdOk() (*string, bool)`

GetMemeIdOk returns a tuple with the MemeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemeId

`func (o *MemeTag) SetMemeId(v string)`

SetMemeId sets MemeId field to given value.

### HasMemeId

`func (o *MemeTag) HasMemeId() bool`

HasMemeId returns a boolean if a field has been set.

### GetMemeIdObject

`func (o *MemeTag) GetMemeIdObject() Meme`

GetMemeIdObject returns the MemeIdObject field if non-nil, zero value otherwise.

### GetMemeIdObjectOk

`func (o *MemeTag) GetMemeIdObjectOk() (*Meme, bool)`

GetMemeIdObjectOk returns a tuple with the MemeIdObject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemeIdObject

`func (o *MemeTag) SetMemeIdObject(v Meme)`

SetMemeIdObject sets MemeIdObject field to given value.

### HasMemeIdObject

`func (o *MemeTag) HasMemeIdObject() bool`

HasMemeIdObject returns a boolean if a field has been set.

### GetTagId

`func (o *MemeTag) GetTagId() string`

GetTagId returns the TagId field if non-nil, zero value otherwise.

### GetTagIdOk

`func (o *MemeTag) GetTagIdOk() (*string, bool)`

GetTagIdOk returns a tuple with the TagId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagId

`func (o *MemeTag) SetTagId(v string)`

SetTagId sets TagId field to given value.

### HasTagId

`func (o *MemeTag) HasTagId() bool`

HasTagId returns a boolean if a field has been set.

### GetTagIdObject

`func (o *MemeTag) GetTagIdObject() Tag`

GetTagIdObject returns the TagIdObject field if non-nil, zero value otherwise.

### GetTagIdObjectOk

`func (o *MemeTag) GetTagIdObjectOk() (*Tag, bool)`

GetTagIdObjectOk returns a tuple with the TagIdObject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagIdObject

`func (o *MemeTag) SetTagIdObject(v Tag)`

SetTagIdObject sets TagIdObject field to given value.

### HasTagIdObject

`func (o *MemeTag) HasTagIdObject() bool`

HasTagIdObject returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *MemeTag) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *MemeTag) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *MemeTag) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *MemeTag) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


