# Detection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClassId** | Pointer to **int64** |  | [optional] 
**Confidence** | Pointer to **float64** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Label** | Pointer to **string** |  | [optional] 
**MemeId** | Pointer to **string** |  | [optional] 
**MemeIdObject** | Pointer to [**Meme**](Meme.md) |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**X1** | Pointer to **float64** |  | [optional] 
**X2** | Pointer to **float64** |  | [optional] 
**Y1** | Pointer to **float64** |  | [optional] 
**Y2** | Pointer to **float64** |  | [optional] 

## Methods

### NewDetection

`func NewDetection() *Detection`

NewDetection instantiates a new Detection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDetectionWithDefaults

`func NewDetectionWithDefaults() *Detection`

NewDetectionWithDefaults instantiates a new Detection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClassId

`func (o *Detection) GetClassId() int64`

GetClassId returns the ClassId field if non-nil, zero value otherwise.

### GetClassIdOk

`func (o *Detection) GetClassIdOk() (*int64, bool)`

GetClassIdOk returns a tuple with the ClassId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassId

`func (o *Detection) SetClassId(v int64)`

SetClassId sets ClassId field to given value.

### HasClassId

`func (o *Detection) HasClassId() bool`

HasClassId returns a boolean if a field has been set.

### GetConfidence

`func (o *Detection) GetConfidence() float64`

GetConfidence returns the Confidence field if non-nil, zero value otherwise.

### GetConfidenceOk

`func (o *Detection) GetConfidenceOk() (*float64, bool)`

GetConfidenceOk returns a tuple with the Confidence field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfidence

`func (o *Detection) SetConfidence(v float64)`

SetConfidence sets Confidence field to given value.

### HasConfidence

`func (o *Detection) HasConfidence() bool`

HasConfidence returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Detection) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Detection) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Detection) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Detection) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDeletedAt

`func (o *Detection) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *Detection) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *Detection) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *Detection) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.

### GetId

`func (o *Detection) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Detection) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Detection) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Detection) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLabel

`func (o *Detection) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *Detection) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *Detection) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *Detection) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetMemeId

`func (o *Detection) GetMemeId() string`

GetMemeId returns the MemeId field if non-nil, zero value otherwise.

### GetMemeIdOk

`func (o *Detection) GetMemeIdOk() (*string, bool)`

GetMemeIdOk returns a tuple with the MemeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemeId

`func (o *Detection) SetMemeId(v string)`

SetMemeId sets MemeId field to given value.

### HasMemeId

`func (o *Detection) HasMemeId() bool`

HasMemeId returns a boolean if a field has been set.

### GetMemeIdObject

`func (o *Detection) GetMemeIdObject() Meme`

GetMemeIdObject returns the MemeIdObject field if non-nil, zero value otherwise.

### GetMemeIdObjectOk

`func (o *Detection) GetMemeIdObjectOk() (*Meme, bool)`

GetMemeIdObjectOk returns a tuple with the MemeIdObject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemeIdObject

`func (o *Detection) SetMemeIdObject(v Meme)`

SetMemeIdObject sets MemeIdObject field to given value.

### HasMemeIdObject

`func (o *Detection) HasMemeIdObject() bool`

HasMemeIdObject returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Detection) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Detection) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Detection) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Detection) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetX1

`func (o *Detection) GetX1() float64`

GetX1 returns the X1 field if non-nil, zero value otherwise.

### GetX1Ok

`func (o *Detection) GetX1Ok() (*float64, bool)`

GetX1Ok returns a tuple with the X1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetX1

`func (o *Detection) SetX1(v float64)`

SetX1 sets X1 field to given value.

### HasX1

`func (o *Detection) HasX1() bool`

HasX1 returns a boolean if a field has been set.

### GetX2

`func (o *Detection) GetX2() float64`

GetX2 returns the X2 field if non-nil, zero value otherwise.

### GetX2Ok

`func (o *Detection) GetX2Ok() (*float64, bool)`

GetX2Ok returns a tuple with the X2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetX2

`func (o *Detection) SetX2(v float64)`

SetX2 sets X2 field to given value.

### HasX2

`func (o *Detection) HasX2() bool`

HasX2 returns a boolean if a field has been set.

### GetY1

`func (o *Detection) GetY1() float64`

GetY1 returns the Y1 field if non-nil, zero value otherwise.

### GetY1Ok

`func (o *Detection) GetY1Ok() (*float64, bool)`

GetY1Ok returns a tuple with the Y1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetY1

`func (o *Detection) SetY1(v float64)`

SetY1 sets Y1 field to given value.

### HasY1

`func (o *Detection) HasY1() bool`

HasY1 returns a boolean if a field has been set.

### GetY2

`func (o *Detection) GetY2() float64`

GetY2 returns the Y2 field if non-nil, zero value otherwise.

### GetY2Ok

`func (o *Detection) GetY2Ok() (*float64, bool)`

GetY2Ok returns a tuple with the Y2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetY2

`func (o *Detection) SetY2(v float64)`

SetY2 sets Y2 field to given value.

### HasY2

`func (o *Detection) HasY2() bool`

HasY2 returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


