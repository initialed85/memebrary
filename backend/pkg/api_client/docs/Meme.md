# Meme

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AiWorkerClaimedUntil** | Pointer to **time.Time** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**DescriptionGenerated** | Pointer to **int64** |  | [optional] 
**DescriptionStatus** | Pointer to **string** |  | [optional] 
**Filename** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**MetadataVersion** | Pointer to **int64** |  | [optional] 
**MimeType** | Pointer to **string** |  | [optional] 
**OriginalName** | Pointer to **string** |  | [optional] 
**ReferencedByMemeTagMemeIdObjects** | Pointer to [**[]MemeTag**](MemeTag.md) |  | [optional] 
**Size** | Pointer to **int64** |  | [optional] 
**SortOrder** | Pointer to **int64** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewMeme

`func NewMeme() *Meme`

NewMeme instantiates a new Meme object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemeWithDefaults

`func NewMemeWithDefaults() *Meme`

NewMemeWithDefaults instantiates a new Meme object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAiWorkerClaimedUntil

`func (o *Meme) GetAiWorkerClaimedUntil() time.Time`

GetAiWorkerClaimedUntil returns the AiWorkerClaimedUntil field if non-nil, zero value otherwise.

### GetAiWorkerClaimedUntilOk

`func (o *Meme) GetAiWorkerClaimedUntilOk() (*time.Time, bool)`

GetAiWorkerClaimedUntilOk returns a tuple with the AiWorkerClaimedUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAiWorkerClaimedUntil

`func (o *Meme) SetAiWorkerClaimedUntil(v time.Time)`

SetAiWorkerClaimedUntil sets AiWorkerClaimedUntil field to given value.

### HasAiWorkerClaimedUntil

`func (o *Meme) HasAiWorkerClaimedUntil() bool`

HasAiWorkerClaimedUntil returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Meme) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Meme) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Meme) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Meme) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDeletedAt

`func (o *Meme) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *Meme) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *Meme) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *Meme) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.

### GetDescription

`func (o *Meme) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Meme) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Meme) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Meme) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetDescriptionGenerated

`func (o *Meme) GetDescriptionGenerated() int64`

GetDescriptionGenerated returns the DescriptionGenerated field if non-nil, zero value otherwise.

### GetDescriptionGeneratedOk

`func (o *Meme) GetDescriptionGeneratedOk() (*int64, bool)`

GetDescriptionGeneratedOk returns a tuple with the DescriptionGenerated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptionGenerated

`func (o *Meme) SetDescriptionGenerated(v int64)`

SetDescriptionGenerated sets DescriptionGenerated field to given value.

### HasDescriptionGenerated

`func (o *Meme) HasDescriptionGenerated() bool`

HasDescriptionGenerated returns a boolean if a field has been set.

### GetDescriptionStatus

`func (o *Meme) GetDescriptionStatus() string`

GetDescriptionStatus returns the DescriptionStatus field if non-nil, zero value otherwise.

### GetDescriptionStatusOk

`func (o *Meme) GetDescriptionStatusOk() (*string, bool)`

GetDescriptionStatusOk returns a tuple with the DescriptionStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptionStatus

`func (o *Meme) SetDescriptionStatus(v string)`

SetDescriptionStatus sets DescriptionStatus field to given value.

### HasDescriptionStatus

`func (o *Meme) HasDescriptionStatus() bool`

HasDescriptionStatus returns a boolean if a field has been set.

### GetFilename

`func (o *Meme) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *Meme) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *Meme) SetFilename(v string)`

SetFilename sets Filename field to given value.

### HasFilename

`func (o *Meme) HasFilename() bool`

HasFilename returns a boolean if a field has been set.

### GetId

`func (o *Meme) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Meme) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Meme) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Meme) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMetadataVersion

`func (o *Meme) GetMetadataVersion() int64`

GetMetadataVersion returns the MetadataVersion field if non-nil, zero value otherwise.

### GetMetadataVersionOk

`func (o *Meme) GetMetadataVersionOk() (*int64, bool)`

GetMetadataVersionOk returns a tuple with the MetadataVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadataVersion

`func (o *Meme) SetMetadataVersion(v int64)`

SetMetadataVersion sets MetadataVersion field to given value.

### HasMetadataVersion

`func (o *Meme) HasMetadataVersion() bool`

HasMetadataVersion returns a boolean if a field has been set.

### GetMimeType

`func (o *Meme) GetMimeType() string`

GetMimeType returns the MimeType field if non-nil, zero value otherwise.

### GetMimeTypeOk

`func (o *Meme) GetMimeTypeOk() (*string, bool)`

GetMimeTypeOk returns a tuple with the MimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMimeType

`func (o *Meme) SetMimeType(v string)`

SetMimeType sets MimeType field to given value.

### HasMimeType

`func (o *Meme) HasMimeType() bool`

HasMimeType returns a boolean if a field has been set.

### GetOriginalName

`func (o *Meme) GetOriginalName() string`

GetOriginalName returns the OriginalName field if non-nil, zero value otherwise.

### GetOriginalNameOk

`func (o *Meme) GetOriginalNameOk() (*string, bool)`

GetOriginalNameOk returns a tuple with the OriginalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalName

`func (o *Meme) SetOriginalName(v string)`

SetOriginalName sets OriginalName field to given value.

### HasOriginalName

`func (o *Meme) HasOriginalName() bool`

HasOriginalName returns a boolean if a field has been set.

### GetReferencedByMemeTagMemeIdObjects

`func (o *Meme) GetReferencedByMemeTagMemeIdObjects() []MemeTag`

GetReferencedByMemeTagMemeIdObjects returns the ReferencedByMemeTagMemeIdObjects field if non-nil, zero value otherwise.

### GetReferencedByMemeTagMemeIdObjectsOk

`func (o *Meme) GetReferencedByMemeTagMemeIdObjectsOk() (*[]MemeTag, bool)`

GetReferencedByMemeTagMemeIdObjectsOk returns a tuple with the ReferencedByMemeTagMemeIdObjects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReferencedByMemeTagMemeIdObjects

`func (o *Meme) SetReferencedByMemeTagMemeIdObjects(v []MemeTag)`

SetReferencedByMemeTagMemeIdObjects sets ReferencedByMemeTagMemeIdObjects field to given value.

### HasReferencedByMemeTagMemeIdObjects

`func (o *Meme) HasReferencedByMemeTagMemeIdObjects() bool`

HasReferencedByMemeTagMemeIdObjects returns a boolean if a field has been set.

### GetSize

`func (o *Meme) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Meme) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Meme) SetSize(v int64)`

SetSize sets Size field to given value.

### HasSize

`func (o *Meme) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetSortOrder

`func (o *Meme) GetSortOrder() int64`

GetSortOrder returns the SortOrder field if non-nil, zero value otherwise.

### GetSortOrderOk

`func (o *Meme) GetSortOrderOk() (*int64, bool)`

GetSortOrderOk returns a tuple with the SortOrder field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSortOrder

`func (o *Meme) SetSortOrder(v int64)`

SetSortOrder sets SortOrder field to given value.

### HasSortOrder

`func (o *Meme) HasSortOrder() bool`

HasSortOrder returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Meme) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Meme) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Meme) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Meme) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


