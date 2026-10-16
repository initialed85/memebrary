# \DetectionAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteDetection**](DetectionAPI.md#DeleteDetection) | **Delete** /api/detections/{primaryKey} | 
[**GetDetection**](DetectionAPI.md#GetDetection) | **Get** /api/detections/{primaryKey} | 
[**GetDetections**](DetectionAPI.md#GetDetections) | **Get** /api/detections | 
[**PatchDetection**](DetectionAPI.md#PatchDetection) | **Patch** /api/detections/{primaryKey} | 
[**PostDetections**](DetectionAPI.md#PostDetections) | **Post** /api/detections | 



## DeleteDetection

> DeleteDetection(ctx, primaryKey).Depth(depth).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	primaryKey := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Path parameter primaryKey
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DetectionAPI.DeleteDetection(context.Background(), primaryKey).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DetectionAPI.DeleteDetection``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDetectionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **depth** | **int64** | Query parameter depth | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDetection

> ResponseWithGenericOfDetection GetDetection(ctx, primaryKey).Depth(depth).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	primaryKey := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Path parameter primaryKey
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DetectionAPI.GetDetection(context.Background(), primaryKey).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DetectionAPI.GetDetection``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDetection`: ResponseWithGenericOfDetection
	fmt.Fprintf(os.Stdout, "Response from `DetectionAPI.GetDetection`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDetectionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfDetection**](ResponseWithGenericOfDetection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDetections

> ResponseWithGenericOfDetection GetDetections(ctx).Limit(limit).Offset(offset).Depth(depth).MemeLoad(memeLoad).IdEq(idEq).IdNe(idNe).IdGt(idGt).IdGte(idGte).IdLt(idLt).IdLte(idLte).IdIn(idIn).IdNotin(idNotin).IdContains(idContains).IdNotcontains(idNotcontains).IdLike(idLike).IdNotlike(idNotlike).IdIlike(idIlike).IdNotilike(idNotilike).IdDesc(idDesc).IdAsc(idAsc).CreatedAtEq(createdAtEq).CreatedAtNe(createdAtNe).CreatedAtGt(createdAtGt).CreatedAtGte(createdAtGte).CreatedAtLt(createdAtLt).CreatedAtLte(createdAtLte).CreatedAtIn(createdAtIn).CreatedAtNotin(createdAtNotin).CreatedAtContains(createdAtContains).CreatedAtNotcontains(createdAtNotcontains).CreatedAtLike(createdAtLike).CreatedAtNotlike(createdAtNotlike).CreatedAtIlike(createdAtIlike).CreatedAtNotilike(createdAtNotilike).CreatedAtDesc(createdAtDesc).CreatedAtAsc(createdAtAsc).UpdatedAtEq(updatedAtEq).UpdatedAtNe(updatedAtNe).UpdatedAtGt(updatedAtGt).UpdatedAtGte(updatedAtGte).UpdatedAtLt(updatedAtLt).UpdatedAtLte(updatedAtLte).UpdatedAtIn(updatedAtIn).UpdatedAtNotin(updatedAtNotin).UpdatedAtContains(updatedAtContains).UpdatedAtNotcontains(updatedAtNotcontains).UpdatedAtLike(updatedAtLike).UpdatedAtNotlike(updatedAtNotlike).UpdatedAtIlike(updatedAtIlike).UpdatedAtNotilike(updatedAtNotilike).UpdatedAtDesc(updatedAtDesc).UpdatedAtAsc(updatedAtAsc).DeletedAtEq(deletedAtEq).DeletedAtNe(deletedAtNe).DeletedAtGt(deletedAtGt).DeletedAtGte(deletedAtGte).DeletedAtLt(deletedAtLt).DeletedAtLte(deletedAtLte).DeletedAtIn(deletedAtIn).DeletedAtNotin(deletedAtNotin).DeletedAtContains(deletedAtContains).DeletedAtNotcontains(deletedAtNotcontains).DeletedAtLike(deletedAtLike).DeletedAtNotlike(deletedAtNotlike).DeletedAtIlike(deletedAtIlike).DeletedAtNotilike(deletedAtNotilike).DeletedAtDesc(deletedAtDesc).DeletedAtAsc(deletedAtAsc).LabelEq(labelEq).LabelNe(labelNe).LabelGt(labelGt).LabelGte(labelGte).LabelLt(labelLt).LabelLte(labelLte).LabelIn(labelIn).LabelNotin(labelNotin).LabelContains(labelContains).LabelNotcontains(labelNotcontains).LabelLike(labelLike).LabelNotlike(labelNotlike).LabelIlike(labelIlike).LabelNotilike(labelNotilike).LabelDesc(labelDesc).LabelAsc(labelAsc).ClassIdEq(classIdEq).ClassIdNe(classIdNe).ClassIdGt(classIdGt).ClassIdGte(classIdGte).ClassIdLt(classIdLt).ClassIdLte(classIdLte).ClassIdIn(classIdIn).ClassIdNotin(classIdNotin).ClassIdContains(classIdContains).ClassIdNotcontains(classIdNotcontains).ClassIdDesc(classIdDesc).ClassIdAsc(classIdAsc).ConfidenceEq(confidenceEq).ConfidenceNe(confidenceNe).ConfidenceGt(confidenceGt).ConfidenceGte(confidenceGte).ConfidenceLt(confidenceLt).ConfidenceLte(confidenceLte).ConfidenceIn(confidenceIn).ConfidenceNotin(confidenceNotin).ConfidenceContains(confidenceContains).ConfidenceNotcontains(confidenceNotcontains).ConfidenceDesc(confidenceDesc).ConfidenceAsc(confidenceAsc).X1Eq(x1Eq).X1Ne(x1Ne).X1Gt(x1Gt).X1Gte(x1Gte).X1Lt(x1Lt).X1Lte(x1Lte).X1In(x1In).X1Notin(x1Notin).X1Contains(x1Contains).X1Notcontains(x1Notcontains).X1Desc(x1Desc).X1Asc(x1Asc).Y1Eq(y1Eq).Y1Ne(y1Ne).Y1Gt(y1Gt).Y1Gte(y1Gte).Y1Lt(y1Lt).Y1Lte(y1Lte).Y1In(y1In).Y1Notin(y1Notin).Y1Contains(y1Contains).Y1Notcontains(y1Notcontains).Y1Desc(y1Desc).Y1Asc(y1Asc).X2Eq(x2Eq).X2Ne(x2Ne).X2Gt(x2Gt).X2Gte(x2Gte).X2Lt(x2Lt).X2Lte(x2Lte).X2In(x2In).X2Notin(x2Notin).X2Contains(x2Contains).X2Notcontains(x2Notcontains).X2Desc(x2Desc).X2Asc(x2Asc).Y2Eq(y2Eq).Y2Ne(y2Ne).Y2Gt(y2Gt).Y2Gte(y2Gte).Y2Lt(y2Lt).Y2Lte(y2Lte).Y2In(y2In).Y2Notin(y2Notin).Y2Contains(y2Contains).Y2Notcontains(y2Notcontains).Y2Desc(y2Desc).Y2Asc(y2Asc).MemeIdEq(memeIdEq).MemeIdNe(memeIdNe).MemeIdGt(memeIdGt).MemeIdGte(memeIdGte).MemeIdLt(memeIdLt).MemeIdLte(memeIdLte).MemeIdIn(memeIdIn).MemeIdNotin(memeIdNotin).MemeIdContains(memeIdContains).MemeIdNotcontains(memeIdNotcontains).MemeIdLike(memeIdLike).MemeIdNotlike(memeIdNotlike).MemeIdIlike(memeIdIlike).MemeIdNotilike(memeIdNotilike).MemeIdDesc(memeIdDesc).MemeIdAsc(memeIdAsc).MemeIdObjectContains(memeIdObjectContains).MemeIdObjectNotcontains(memeIdObjectNotcontains).MemeIdObjectDesc(memeIdObjectDesc).MemeIdObjectAsc(memeIdObjectAsc).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	limit := int32(56) // int32 | SQL LIMIT operator (optional)
	offset := int32(56) // int32 | SQL OFFSET operator (optional)
	depth := int32(56) // int32 | Max recursion depth for loading foreign objects; default = 1  (0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc) (optional)
	memeLoad := "memeLoad_example" // string | load the given directly related object, value is ignored (presence of key is sufficient) (optional)
	idEq := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL = comparison (optional)
	idNe := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL != comparison (optional)
	idGt := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL > comparison, may not work with all column types (optional)
	idGte := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL >= comparison, may not work with all column types (optional)
	idLt := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL < comparison, may not work with all column types (optional)
	idLte := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL <= comparison, may not work with all column types (optional)
	idIn := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL IN comparison, permits comma-separated values (optional)
	idNotin := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	idContains := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL @> comparison (optional)
	idNotcontains := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT @> comparison (optional)
	idLike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	idNotlike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	idIlike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	idNotilike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	idDesc := "idDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	idAsc := "idAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	createdAtEq := time.Now() // time.Time | SQL = comparison (optional)
	createdAtNe := time.Now() // time.Time | SQL != comparison (optional)
	createdAtGt := time.Now() // time.Time | SQL > comparison, may not work with all column types (optional)
	createdAtGte := time.Now() // time.Time | SQL >= comparison, may not work with all column types (optional)
	createdAtLt := time.Now() // time.Time | SQL < comparison, may not work with all column types (optional)
	createdAtLte := time.Now() // time.Time | SQL <= comparison, may not work with all column types (optional)
	createdAtIn := time.Now() // time.Time | SQL IN comparison, permits comma-separated values (optional)
	createdAtNotin := time.Now() // time.Time | SQL NOT IN comparison, permits comma-separated values (optional)
	createdAtContains := time.Now() // time.Time | SQL @> comparison (optional)
	createdAtNotcontains := time.Now() // time.Time | SQL NOT @> comparison (optional)
	createdAtLike := time.Now() // time.Time | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	createdAtNotlike := time.Now() // time.Time | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	createdAtIlike := time.Now() // time.Time | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	createdAtNotilike := time.Now() // time.Time | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	createdAtDesc := "createdAtDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	createdAtAsc := "createdAtAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	updatedAtEq := time.Now() // time.Time | SQL = comparison (optional)
	updatedAtNe := time.Now() // time.Time | SQL != comparison (optional)
	updatedAtGt := time.Now() // time.Time | SQL > comparison, may not work with all column types (optional)
	updatedAtGte := time.Now() // time.Time | SQL >= comparison, may not work with all column types (optional)
	updatedAtLt := time.Now() // time.Time | SQL < comparison, may not work with all column types (optional)
	updatedAtLte := time.Now() // time.Time | SQL <= comparison, may not work with all column types (optional)
	updatedAtIn := time.Now() // time.Time | SQL IN comparison, permits comma-separated values (optional)
	updatedAtNotin := time.Now() // time.Time | SQL NOT IN comparison, permits comma-separated values (optional)
	updatedAtContains := time.Now() // time.Time | SQL @> comparison (optional)
	updatedAtNotcontains := time.Now() // time.Time | SQL NOT @> comparison (optional)
	updatedAtLike := time.Now() // time.Time | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	updatedAtNotlike := time.Now() // time.Time | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	updatedAtIlike := time.Now() // time.Time | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	updatedAtNotilike := time.Now() // time.Time | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	updatedAtDesc := "updatedAtDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	updatedAtAsc := "updatedAtAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	deletedAtEq := time.Now() // time.Time | SQL = comparison (optional)
	deletedAtNe := time.Now() // time.Time | SQL != comparison (optional)
	deletedAtGt := time.Now() // time.Time | SQL > comparison, may not work with all column types (optional)
	deletedAtGte := time.Now() // time.Time | SQL >= comparison, may not work with all column types (optional)
	deletedAtLt := time.Now() // time.Time | SQL < comparison, may not work with all column types (optional)
	deletedAtLte := time.Now() // time.Time | SQL <= comparison, may not work with all column types (optional)
	deletedAtIn := time.Now() // time.Time | SQL IN comparison, permits comma-separated values (optional)
	deletedAtNotin := time.Now() // time.Time | SQL NOT IN comparison, permits comma-separated values (optional)
	deletedAtContains := time.Now() // time.Time | SQL @> comparison (optional)
	deletedAtNotcontains := time.Now() // time.Time | SQL NOT @> comparison (optional)
	deletedAtLike := time.Now() // time.Time | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	deletedAtNotlike := time.Now() // time.Time | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	deletedAtIlike := time.Now() // time.Time | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	deletedAtNotilike := time.Now() // time.Time | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	deletedAtDesc := "deletedAtDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	deletedAtAsc := "deletedAtAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	labelEq := "labelEq_example" // string | SQL = comparison (optional)
	labelNe := "labelNe_example" // string | SQL != comparison (optional)
	labelGt := "labelGt_example" // string | SQL > comparison, may not work with all column types (optional)
	labelGte := "labelGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	labelLt := "labelLt_example" // string | SQL < comparison, may not work with all column types (optional)
	labelLte := "labelLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	labelIn := "labelIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	labelNotin := "labelNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	labelContains := "labelContains_example" // string | SQL @> comparison (optional)
	labelNotcontains := "labelNotcontains_example" // string | SQL NOT @> comparison (optional)
	labelLike := "labelLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	labelNotlike := "labelNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	labelIlike := "labelIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	labelNotilike := "labelNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	labelDesc := "labelDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	labelAsc := "labelAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	classIdEq := int64(789) // int64 | SQL = comparison (optional)
	classIdNe := int64(789) // int64 | SQL != comparison (optional)
	classIdGt := int64(789) // int64 | SQL > comparison, may not work with all column types (optional)
	classIdGte := int64(789) // int64 | SQL >= comparison, may not work with all column types (optional)
	classIdLt := int64(789) // int64 | SQL < comparison, may not work with all column types (optional)
	classIdLte := int64(789) // int64 | SQL <= comparison, may not work with all column types (optional)
	classIdIn := int64(789) // int64 | SQL IN comparison, permits comma-separated values (optional)
	classIdNotin := int64(789) // int64 | SQL NOT IN comparison, permits comma-separated values (optional)
	classIdContains := int64(789) // int64 | SQL @> comparison (optional)
	classIdNotcontains := int64(789) // int64 | SQL NOT @> comparison (optional)
	classIdDesc := "classIdDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	classIdAsc := "classIdAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	confidenceEq := float64(1.2) // float64 | SQL = comparison (optional)
	confidenceNe := float64(1.2) // float64 | SQL != comparison (optional)
	confidenceGt := float64(1.2) // float64 | SQL > comparison, may not work with all column types (optional)
	confidenceGte := float64(1.2) // float64 | SQL >= comparison, may not work with all column types (optional)
	confidenceLt := float64(1.2) // float64 | SQL < comparison, may not work with all column types (optional)
	confidenceLte := float64(1.2) // float64 | SQL <= comparison, may not work with all column types (optional)
	confidenceIn := float64(1.2) // float64 | SQL IN comparison, permits comma-separated values (optional)
	confidenceNotin := float64(1.2) // float64 | SQL NOT IN comparison, permits comma-separated values (optional)
	confidenceContains := float64(1.2) // float64 | SQL @> comparison (optional)
	confidenceNotcontains := float64(1.2) // float64 | SQL NOT @> comparison (optional)
	confidenceDesc := "confidenceDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	confidenceAsc := "confidenceAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	x1Eq := float64(1.2) // float64 | SQL = comparison (optional)
	x1Ne := float64(1.2) // float64 | SQL != comparison (optional)
	x1Gt := float64(1.2) // float64 | SQL > comparison, may not work with all column types (optional)
	x1Gte := float64(1.2) // float64 | SQL >= comparison, may not work with all column types (optional)
	x1Lt := float64(1.2) // float64 | SQL < comparison, may not work with all column types (optional)
	x1Lte := float64(1.2) // float64 | SQL <= comparison, may not work with all column types (optional)
	x1In := float64(1.2) // float64 | SQL IN comparison, permits comma-separated values (optional)
	x1Notin := float64(1.2) // float64 | SQL NOT IN comparison, permits comma-separated values (optional)
	x1Contains := float64(1.2) // float64 | SQL @> comparison (optional)
	x1Notcontains := float64(1.2) // float64 | SQL NOT @> comparison (optional)
	x1Desc := "x1Desc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	x1Asc := "x1Asc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	y1Eq := float64(1.2) // float64 | SQL = comparison (optional)
	y1Ne := float64(1.2) // float64 | SQL != comparison (optional)
	y1Gt := float64(1.2) // float64 | SQL > comparison, may not work with all column types (optional)
	y1Gte := float64(1.2) // float64 | SQL >= comparison, may not work with all column types (optional)
	y1Lt := float64(1.2) // float64 | SQL < comparison, may not work with all column types (optional)
	y1Lte := float64(1.2) // float64 | SQL <= comparison, may not work with all column types (optional)
	y1In := float64(1.2) // float64 | SQL IN comparison, permits comma-separated values (optional)
	y1Notin := float64(1.2) // float64 | SQL NOT IN comparison, permits comma-separated values (optional)
	y1Contains := float64(1.2) // float64 | SQL @> comparison (optional)
	y1Notcontains := float64(1.2) // float64 | SQL NOT @> comparison (optional)
	y1Desc := "y1Desc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	y1Asc := "y1Asc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	x2Eq := float64(1.2) // float64 | SQL = comparison (optional)
	x2Ne := float64(1.2) // float64 | SQL != comparison (optional)
	x2Gt := float64(1.2) // float64 | SQL > comparison, may not work with all column types (optional)
	x2Gte := float64(1.2) // float64 | SQL >= comparison, may not work with all column types (optional)
	x2Lt := float64(1.2) // float64 | SQL < comparison, may not work with all column types (optional)
	x2Lte := float64(1.2) // float64 | SQL <= comparison, may not work with all column types (optional)
	x2In := float64(1.2) // float64 | SQL IN comparison, permits comma-separated values (optional)
	x2Notin := float64(1.2) // float64 | SQL NOT IN comparison, permits comma-separated values (optional)
	x2Contains := float64(1.2) // float64 | SQL @> comparison (optional)
	x2Notcontains := float64(1.2) // float64 | SQL NOT @> comparison (optional)
	x2Desc := "x2Desc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	x2Asc := "x2Asc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	y2Eq := float64(1.2) // float64 | SQL = comparison (optional)
	y2Ne := float64(1.2) // float64 | SQL != comparison (optional)
	y2Gt := float64(1.2) // float64 | SQL > comparison, may not work with all column types (optional)
	y2Gte := float64(1.2) // float64 | SQL >= comparison, may not work with all column types (optional)
	y2Lt := float64(1.2) // float64 | SQL < comparison, may not work with all column types (optional)
	y2Lte := float64(1.2) // float64 | SQL <= comparison, may not work with all column types (optional)
	y2In := float64(1.2) // float64 | SQL IN comparison, permits comma-separated values (optional)
	y2Notin := float64(1.2) // float64 | SQL NOT IN comparison, permits comma-separated values (optional)
	y2Contains := float64(1.2) // float64 | SQL @> comparison (optional)
	y2Notcontains := float64(1.2) // float64 | SQL NOT @> comparison (optional)
	y2Desc := "y2Desc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	y2Asc := "y2Asc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	memeIdEq := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL = comparison (optional)
	memeIdNe := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL != comparison (optional)
	memeIdGt := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL > comparison, may not work with all column types (optional)
	memeIdGte := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL >= comparison, may not work with all column types (optional)
	memeIdLt := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL < comparison, may not work with all column types (optional)
	memeIdLte := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL <= comparison, may not work with all column types (optional)
	memeIdIn := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL IN comparison, permits comma-separated values (optional)
	memeIdNotin := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	memeIdContains := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL @> comparison (optional)
	memeIdNotcontains := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT @> comparison (optional)
	memeIdLike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	memeIdNotlike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	memeIdIlike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	memeIdNotilike := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	memeIdDesc := "memeIdDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	memeIdAsc := "memeIdAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	memeIdObjectContains := TODO // interface{} | SQL @> comparison (optional)
	memeIdObjectNotcontains := TODO // interface{} | SQL NOT @> comparison (optional)
	memeIdObjectDesc := "memeIdObjectDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	memeIdObjectAsc := "memeIdObjectAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DetectionAPI.GetDetections(context.Background()).Limit(limit).Offset(offset).Depth(depth).MemeLoad(memeLoad).IdEq(idEq).IdNe(idNe).IdGt(idGt).IdGte(idGte).IdLt(idLt).IdLte(idLte).IdIn(idIn).IdNotin(idNotin).IdContains(idContains).IdNotcontains(idNotcontains).IdLike(idLike).IdNotlike(idNotlike).IdIlike(idIlike).IdNotilike(idNotilike).IdDesc(idDesc).IdAsc(idAsc).CreatedAtEq(createdAtEq).CreatedAtNe(createdAtNe).CreatedAtGt(createdAtGt).CreatedAtGte(createdAtGte).CreatedAtLt(createdAtLt).CreatedAtLte(createdAtLte).CreatedAtIn(createdAtIn).CreatedAtNotin(createdAtNotin).CreatedAtContains(createdAtContains).CreatedAtNotcontains(createdAtNotcontains).CreatedAtLike(createdAtLike).CreatedAtNotlike(createdAtNotlike).CreatedAtIlike(createdAtIlike).CreatedAtNotilike(createdAtNotilike).CreatedAtDesc(createdAtDesc).CreatedAtAsc(createdAtAsc).UpdatedAtEq(updatedAtEq).UpdatedAtNe(updatedAtNe).UpdatedAtGt(updatedAtGt).UpdatedAtGte(updatedAtGte).UpdatedAtLt(updatedAtLt).UpdatedAtLte(updatedAtLte).UpdatedAtIn(updatedAtIn).UpdatedAtNotin(updatedAtNotin).UpdatedAtContains(updatedAtContains).UpdatedAtNotcontains(updatedAtNotcontains).UpdatedAtLike(updatedAtLike).UpdatedAtNotlike(updatedAtNotlike).UpdatedAtIlike(updatedAtIlike).UpdatedAtNotilike(updatedAtNotilike).UpdatedAtDesc(updatedAtDesc).UpdatedAtAsc(updatedAtAsc).DeletedAtEq(deletedAtEq).DeletedAtNe(deletedAtNe).DeletedAtGt(deletedAtGt).DeletedAtGte(deletedAtGte).DeletedAtLt(deletedAtLt).DeletedAtLte(deletedAtLte).DeletedAtIn(deletedAtIn).DeletedAtNotin(deletedAtNotin).DeletedAtContains(deletedAtContains).DeletedAtNotcontains(deletedAtNotcontains).DeletedAtLike(deletedAtLike).DeletedAtNotlike(deletedAtNotlike).DeletedAtIlike(deletedAtIlike).DeletedAtNotilike(deletedAtNotilike).DeletedAtDesc(deletedAtDesc).DeletedAtAsc(deletedAtAsc).LabelEq(labelEq).LabelNe(labelNe).LabelGt(labelGt).LabelGte(labelGte).LabelLt(labelLt).LabelLte(labelLte).LabelIn(labelIn).LabelNotin(labelNotin).LabelContains(labelContains).LabelNotcontains(labelNotcontains).LabelLike(labelLike).LabelNotlike(labelNotlike).LabelIlike(labelIlike).LabelNotilike(labelNotilike).LabelDesc(labelDesc).LabelAsc(labelAsc).ClassIdEq(classIdEq).ClassIdNe(classIdNe).ClassIdGt(classIdGt).ClassIdGte(classIdGte).ClassIdLt(classIdLt).ClassIdLte(classIdLte).ClassIdIn(classIdIn).ClassIdNotin(classIdNotin).ClassIdContains(classIdContains).ClassIdNotcontains(classIdNotcontains).ClassIdDesc(classIdDesc).ClassIdAsc(classIdAsc).ConfidenceEq(confidenceEq).ConfidenceNe(confidenceNe).ConfidenceGt(confidenceGt).ConfidenceGte(confidenceGte).ConfidenceLt(confidenceLt).ConfidenceLte(confidenceLte).ConfidenceIn(confidenceIn).ConfidenceNotin(confidenceNotin).ConfidenceContains(confidenceContains).ConfidenceNotcontains(confidenceNotcontains).ConfidenceDesc(confidenceDesc).ConfidenceAsc(confidenceAsc).X1Eq(x1Eq).X1Ne(x1Ne).X1Gt(x1Gt).X1Gte(x1Gte).X1Lt(x1Lt).X1Lte(x1Lte).X1In(x1In).X1Notin(x1Notin).X1Contains(x1Contains).X1Notcontains(x1Notcontains).X1Desc(x1Desc).X1Asc(x1Asc).Y1Eq(y1Eq).Y1Ne(y1Ne).Y1Gt(y1Gt).Y1Gte(y1Gte).Y1Lt(y1Lt).Y1Lte(y1Lte).Y1In(y1In).Y1Notin(y1Notin).Y1Contains(y1Contains).Y1Notcontains(y1Notcontains).Y1Desc(y1Desc).Y1Asc(y1Asc).X2Eq(x2Eq).X2Ne(x2Ne).X2Gt(x2Gt).X2Gte(x2Gte).X2Lt(x2Lt).X2Lte(x2Lte).X2In(x2In).X2Notin(x2Notin).X2Contains(x2Contains).X2Notcontains(x2Notcontains).X2Desc(x2Desc).X2Asc(x2Asc).Y2Eq(y2Eq).Y2Ne(y2Ne).Y2Gt(y2Gt).Y2Gte(y2Gte).Y2Lt(y2Lt).Y2Lte(y2Lte).Y2In(y2In).Y2Notin(y2Notin).Y2Contains(y2Contains).Y2Notcontains(y2Notcontains).Y2Desc(y2Desc).Y2Asc(y2Asc).MemeIdEq(memeIdEq).MemeIdNe(memeIdNe).MemeIdGt(memeIdGt).MemeIdGte(memeIdGte).MemeIdLt(memeIdLt).MemeIdLte(memeIdLte).MemeIdIn(memeIdIn).MemeIdNotin(memeIdNotin).MemeIdContains(memeIdContains).MemeIdNotcontains(memeIdNotcontains).MemeIdLike(memeIdLike).MemeIdNotlike(memeIdNotlike).MemeIdIlike(memeIdIlike).MemeIdNotilike(memeIdNotilike).MemeIdDesc(memeIdDesc).MemeIdAsc(memeIdAsc).MemeIdObjectContains(memeIdObjectContains).MemeIdObjectNotcontains(memeIdObjectNotcontains).MemeIdObjectDesc(memeIdObjectDesc).MemeIdObjectAsc(memeIdObjectAsc).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DetectionAPI.GetDetections``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDetections`: ResponseWithGenericOfDetection
	fmt.Fprintf(os.Stdout, "Response from `DetectionAPI.GetDetections`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetDetectionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | SQL LIMIT operator | 
 **offset** | **int32** | SQL OFFSET operator | 
 **depth** | **int32** | Max recursion depth for loading foreign objects; default &#x3D; 1  (0 &#x3D; recurse until graph cycle detected, 1 &#x3D; this object only, 2 &#x3D; this object + neighbours, 3 &#x3D; this object + neighbours + their neighbours... etc) | 
 **memeLoad** | **string** | load the given directly related object, value is ignored (presence of key is sufficient) | 
 **idEq** | **string** | SQL &#x3D; comparison | 
 **idNe** | **string** | SQL !&#x3D; comparison | 
 **idGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **idGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **idLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **idLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **idIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **idNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **idContains** | **string** | SQL @&gt; comparison | 
 **idNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **idLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **idNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **idIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **idNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **idDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **idAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **createdAtEq** | **time.Time** | SQL &#x3D; comparison | 
 **createdAtNe** | **time.Time** | SQL !&#x3D; comparison | 
 **createdAtGt** | **time.Time** | SQL &gt; comparison, may not work with all column types | 
 **createdAtGte** | **time.Time** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **createdAtLt** | **time.Time** | SQL &lt; comparison, may not work with all column types | 
 **createdAtLte** | **time.Time** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **createdAtIn** | **time.Time** | SQL IN comparison, permits comma-separated values | 
 **createdAtNotin** | **time.Time** | SQL NOT IN comparison, permits comma-separated values | 
 **createdAtContains** | **time.Time** | SQL @&gt; comparison | 
 **createdAtNotcontains** | **time.Time** | SQL NOT @&gt; comparison | 
 **createdAtLike** | **time.Time** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **createdAtNotlike** | **time.Time** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **createdAtIlike** | **time.Time** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **createdAtNotilike** | **time.Time** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **createdAtDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **createdAtAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **updatedAtEq** | **time.Time** | SQL &#x3D; comparison | 
 **updatedAtNe** | **time.Time** | SQL !&#x3D; comparison | 
 **updatedAtGt** | **time.Time** | SQL &gt; comparison, may not work with all column types | 
 **updatedAtGte** | **time.Time** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **updatedAtLt** | **time.Time** | SQL &lt; comparison, may not work with all column types | 
 **updatedAtLte** | **time.Time** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **updatedAtIn** | **time.Time** | SQL IN comparison, permits comma-separated values | 
 **updatedAtNotin** | **time.Time** | SQL NOT IN comparison, permits comma-separated values | 
 **updatedAtContains** | **time.Time** | SQL @&gt; comparison | 
 **updatedAtNotcontains** | **time.Time** | SQL NOT @&gt; comparison | 
 **updatedAtLike** | **time.Time** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **updatedAtNotlike** | **time.Time** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **updatedAtIlike** | **time.Time** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **updatedAtNotilike** | **time.Time** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **updatedAtDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **updatedAtAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **deletedAtEq** | **time.Time** | SQL &#x3D; comparison | 
 **deletedAtNe** | **time.Time** | SQL !&#x3D; comparison | 
 **deletedAtGt** | **time.Time** | SQL &gt; comparison, may not work with all column types | 
 **deletedAtGte** | **time.Time** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **deletedAtLt** | **time.Time** | SQL &lt; comparison, may not work with all column types | 
 **deletedAtLte** | **time.Time** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **deletedAtIn** | **time.Time** | SQL IN comparison, permits comma-separated values | 
 **deletedAtNotin** | **time.Time** | SQL NOT IN comparison, permits comma-separated values | 
 **deletedAtContains** | **time.Time** | SQL @&gt; comparison | 
 **deletedAtNotcontains** | **time.Time** | SQL NOT @&gt; comparison | 
 **deletedAtLike** | **time.Time** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **deletedAtNotlike** | **time.Time** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **deletedAtIlike** | **time.Time** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **deletedAtNotilike** | **time.Time** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **deletedAtDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **deletedAtAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **labelEq** | **string** | SQL &#x3D; comparison | 
 **labelNe** | **string** | SQL !&#x3D; comparison | 
 **labelGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **labelGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **labelLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **labelLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **labelIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **labelNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **labelContains** | **string** | SQL @&gt; comparison | 
 **labelNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **labelLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **labelNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **labelIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **labelNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **labelDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **labelAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **classIdEq** | **int64** | SQL &#x3D; comparison | 
 **classIdNe** | **int64** | SQL !&#x3D; comparison | 
 **classIdGt** | **int64** | SQL &gt; comparison, may not work with all column types | 
 **classIdGte** | **int64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **classIdLt** | **int64** | SQL &lt; comparison, may not work with all column types | 
 **classIdLte** | **int64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **classIdIn** | **int64** | SQL IN comparison, permits comma-separated values | 
 **classIdNotin** | **int64** | SQL NOT IN comparison, permits comma-separated values | 
 **classIdContains** | **int64** | SQL @&gt; comparison | 
 **classIdNotcontains** | **int64** | SQL NOT @&gt; comparison | 
 **classIdDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **classIdAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **confidenceEq** | **float64** | SQL &#x3D; comparison | 
 **confidenceNe** | **float64** | SQL !&#x3D; comparison | 
 **confidenceGt** | **float64** | SQL &gt; comparison, may not work with all column types | 
 **confidenceGte** | **float64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **confidenceLt** | **float64** | SQL &lt; comparison, may not work with all column types | 
 **confidenceLte** | **float64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **confidenceIn** | **float64** | SQL IN comparison, permits comma-separated values | 
 **confidenceNotin** | **float64** | SQL NOT IN comparison, permits comma-separated values | 
 **confidenceContains** | **float64** | SQL @&gt; comparison | 
 **confidenceNotcontains** | **float64** | SQL NOT @&gt; comparison | 
 **confidenceDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **confidenceAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **x1Eq** | **float64** | SQL &#x3D; comparison | 
 **x1Ne** | **float64** | SQL !&#x3D; comparison | 
 **x1Gt** | **float64** | SQL &gt; comparison, may not work with all column types | 
 **x1Gte** | **float64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **x1Lt** | **float64** | SQL &lt; comparison, may not work with all column types | 
 **x1Lte** | **float64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **x1In** | **float64** | SQL IN comparison, permits comma-separated values | 
 **x1Notin** | **float64** | SQL NOT IN comparison, permits comma-separated values | 
 **x1Contains** | **float64** | SQL @&gt; comparison | 
 **x1Notcontains** | **float64** | SQL NOT @&gt; comparison | 
 **x1Desc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **x1Asc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **y1Eq** | **float64** | SQL &#x3D; comparison | 
 **y1Ne** | **float64** | SQL !&#x3D; comparison | 
 **y1Gt** | **float64** | SQL &gt; comparison, may not work with all column types | 
 **y1Gte** | **float64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **y1Lt** | **float64** | SQL &lt; comparison, may not work with all column types | 
 **y1Lte** | **float64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **y1In** | **float64** | SQL IN comparison, permits comma-separated values | 
 **y1Notin** | **float64** | SQL NOT IN comparison, permits comma-separated values | 
 **y1Contains** | **float64** | SQL @&gt; comparison | 
 **y1Notcontains** | **float64** | SQL NOT @&gt; comparison | 
 **y1Desc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **y1Asc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **x2Eq** | **float64** | SQL &#x3D; comparison | 
 **x2Ne** | **float64** | SQL !&#x3D; comparison | 
 **x2Gt** | **float64** | SQL &gt; comparison, may not work with all column types | 
 **x2Gte** | **float64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **x2Lt** | **float64** | SQL &lt; comparison, may not work with all column types | 
 **x2Lte** | **float64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **x2In** | **float64** | SQL IN comparison, permits comma-separated values | 
 **x2Notin** | **float64** | SQL NOT IN comparison, permits comma-separated values | 
 **x2Contains** | **float64** | SQL @&gt; comparison | 
 **x2Notcontains** | **float64** | SQL NOT @&gt; comparison | 
 **x2Desc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **x2Asc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **y2Eq** | **float64** | SQL &#x3D; comparison | 
 **y2Ne** | **float64** | SQL !&#x3D; comparison | 
 **y2Gt** | **float64** | SQL &gt; comparison, may not work with all column types | 
 **y2Gte** | **float64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **y2Lt** | **float64** | SQL &lt; comparison, may not work with all column types | 
 **y2Lte** | **float64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **y2In** | **float64** | SQL IN comparison, permits comma-separated values | 
 **y2Notin** | **float64** | SQL NOT IN comparison, permits comma-separated values | 
 **y2Contains** | **float64** | SQL @&gt; comparison | 
 **y2Notcontains** | **float64** | SQL NOT @&gt; comparison | 
 **y2Desc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **y2Asc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **memeIdEq** | **string** | SQL &#x3D; comparison | 
 **memeIdNe** | **string** | SQL !&#x3D; comparison | 
 **memeIdGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **memeIdGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **memeIdLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **memeIdLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **memeIdIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **memeIdNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **memeIdContains** | **string** | SQL @&gt; comparison | 
 **memeIdNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **memeIdLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **memeIdNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **memeIdIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **memeIdNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **memeIdDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **memeIdAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **memeIdObjectContains** | [**interface{}**](interface{}.md) | SQL @&gt; comparison | 
 **memeIdObjectNotcontains** | [**interface{}**](interface{}.md) | SQL NOT @&gt; comparison | 
 **memeIdObjectDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **memeIdObjectAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 

### Return type

[**ResponseWithGenericOfDetection**](ResponseWithGenericOfDetection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchDetection

> ResponseWithGenericOfDetection PatchDetection(ctx, primaryKey).Detection(detection).Depth(depth).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	primaryKey := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Path parameter primaryKey
	detection := *openapiclient.NewDetection() // Detection | 
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DetectionAPI.PatchDetection(context.Background(), primaryKey).Detection(detection).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DetectionAPI.PatchDetection``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PatchDetection`: ResponseWithGenericOfDetection
	fmt.Fprintf(os.Stdout, "Response from `DetectionAPI.PatchDetection`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchDetectionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **detection** | [**Detection**](Detection.md) |  | 
 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfDetection**](ResponseWithGenericOfDetection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostDetections

> ResponseWithGenericOfDetection PostDetections(ctx).Detection(detection).Depth(depth).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	detection := []openapiclient.Detection{*openapiclient.NewDetection()} // []Detection | 
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DetectionAPI.PostDetections(context.Background()).Detection(detection).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DetectionAPI.PostDetections``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PostDetections`: ResponseWithGenericOfDetection
	fmt.Fprintf(os.Stdout, "Response from `DetectionAPI.PostDetections`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostDetectionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **detection** | [**[]Detection**](Detection.md) |  | 
 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfDetection**](ResponseWithGenericOfDetection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

