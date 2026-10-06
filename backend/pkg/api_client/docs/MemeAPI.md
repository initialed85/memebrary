# \MemeAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteMeme**](MemeAPI.md#DeleteMeme) | **Delete** /api/memes/{primaryKey} | 
[**GetMeme**](MemeAPI.md#GetMeme) | **Get** /api/memes/{primaryKey} | 
[**GetMemes**](MemeAPI.md#GetMemes) | **Get** /api/memes | 
[**PatchMeme**](MemeAPI.md#PatchMeme) | **Patch** /api/memes/{primaryKey} | 
[**PostMemes**](MemeAPI.md#PostMemes) | **Post** /api/memes | 
[**PostMemesAiWorkerClaim**](MemeAPI.md#PostMemesAiWorkerClaim) | **Post** /api/memes/{primaryKey}/ai-worker-claim | 



## DeleteMeme

> DeleteMeme(ctx, primaryKey).Depth(depth).Execute()



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
	r, err := apiClient.MemeAPI.DeleteMeme(context.Background(), primaryKey).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.DeleteMeme``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteMemeRequest struct via the builder pattern


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


## GetMeme

> ResponseWithGenericOfMeme GetMeme(ctx, primaryKey).Depth(depth).Execute()



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
	resp, r, err := apiClient.MemeAPI.GetMeme(context.Background(), primaryKey).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.GetMeme``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMeme`: ResponseWithGenericOfMeme
	fmt.Fprintf(os.Stdout, "Response from `MemeAPI.GetMeme`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMemeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfMeme**](ResponseWithGenericOfMeme.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMemes

> ResponseWithGenericOfMeme GetMemes(ctx).Limit(limit).Offset(offset).Depth(depth).ReferencedByMemeTagLoad(referencedByMemeTagLoad).IdEq(idEq).IdNe(idNe).IdGt(idGt).IdGte(idGte).IdLt(idLt).IdLte(idLte).IdIn(idIn).IdNotin(idNotin).IdContains(idContains).IdNotcontains(idNotcontains).IdLike(idLike).IdNotlike(idNotlike).IdIlike(idIlike).IdNotilike(idNotilike).IdDesc(idDesc).IdAsc(idAsc).CreatedAtEq(createdAtEq).CreatedAtNe(createdAtNe).CreatedAtGt(createdAtGt).CreatedAtGte(createdAtGte).CreatedAtLt(createdAtLt).CreatedAtLte(createdAtLte).CreatedAtIn(createdAtIn).CreatedAtNotin(createdAtNotin).CreatedAtContains(createdAtContains).CreatedAtNotcontains(createdAtNotcontains).CreatedAtLike(createdAtLike).CreatedAtNotlike(createdAtNotlike).CreatedAtIlike(createdAtIlike).CreatedAtNotilike(createdAtNotilike).CreatedAtDesc(createdAtDesc).CreatedAtAsc(createdAtAsc).UpdatedAtEq(updatedAtEq).UpdatedAtNe(updatedAtNe).UpdatedAtGt(updatedAtGt).UpdatedAtGte(updatedAtGte).UpdatedAtLt(updatedAtLt).UpdatedAtLte(updatedAtLte).UpdatedAtIn(updatedAtIn).UpdatedAtNotin(updatedAtNotin).UpdatedAtContains(updatedAtContains).UpdatedAtNotcontains(updatedAtNotcontains).UpdatedAtLike(updatedAtLike).UpdatedAtNotlike(updatedAtNotlike).UpdatedAtIlike(updatedAtIlike).UpdatedAtNotilike(updatedAtNotilike).UpdatedAtDesc(updatedAtDesc).UpdatedAtAsc(updatedAtAsc).DeletedAtEq(deletedAtEq).DeletedAtNe(deletedAtNe).DeletedAtGt(deletedAtGt).DeletedAtGte(deletedAtGte).DeletedAtLt(deletedAtLt).DeletedAtLte(deletedAtLte).DeletedAtIn(deletedAtIn).DeletedAtNotin(deletedAtNotin).DeletedAtContains(deletedAtContains).DeletedAtNotcontains(deletedAtNotcontains).DeletedAtLike(deletedAtLike).DeletedAtNotlike(deletedAtNotlike).DeletedAtIlike(deletedAtIlike).DeletedAtNotilike(deletedAtNotilike).DeletedAtDesc(deletedAtDesc).DeletedAtAsc(deletedAtAsc).FilenameEq(filenameEq).FilenameNe(filenameNe).FilenameGt(filenameGt).FilenameGte(filenameGte).FilenameLt(filenameLt).FilenameLte(filenameLte).FilenameIn(filenameIn).FilenameNotin(filenameNotin).FilenameContains(filenameContains).FilenameNotcontains(filenameNotcontains).FilenameLike(filenameLike).FilenameNotlike(filenameNotlike).FilenameIlike(filenameIlike).FilenameNotilike(filenameNotilike).FilenameDesc(filenameDesc).FilenameAsc(filenameAsc).OriginalNameEq(originalNameEq).OriginalNameNe(originalNameNe).OriginalNameGt(originalNameGt).OriginalNameGte(originalNameGte).OriginalNameLt(originalNameLt).OriginalNameLte(originalNameLte).OriginalNameIn(originalNameIn).OriginalNameNotin(originalNameNotin).OriginalNameContains(originalNameContains).OriginalNameNotcontains(originalNameNotcontains).OriginalNameLike(originalNameLike).OriginalNameNotlike(originalNameNotlike).OriginalNameIlike(originalNameIlike).OriginalNameNotilike(originalNameNotilike).OriginalNameDesc(originalNameDesc).OriginalNameAsc(originalNameAsc).MimeTypeEq(mimeTypeEq).MimeTypeNe(mimeTypeNe).MimeTypeGt(mimeTypeGt).MimeTypeGte(mimeTypeGte).MimeTypeLt(mimeTypeLt).MimeTypeLte(mimeTypeLte).MimeTypeIn(mimeTypeIn).MimeTypeNotin(mimeTypeNotin).MimeTypeContains(mimeTypeContains).MimeTypeNotcontains(mimeTypeNotcontains).MimeTypeLike(mimeTypeLike).MimeTypeNotlike(mimeTypeNotlike).MimeTypeIlike(mimeTypeIlike).MimeTypeNotilike(mimeTypeNotilike).MimeTypeDesc(mimeTypeDesc).MimeTypeAsc(mimeTypeAsc).SizeEq(sizeEq).SizeNe(sizeNe).SizeGt(sizeGt).SizeGte(sizeGte).SizeLt(sizeLt).SizeLte(sizeLte).SizeIn(sizeIn).SizeNotin(sizeNotin).SizeContains(sizeContains).SizeNotcontains(sizeNotcontains).SizeDesc(sizeDesc).SizeAsc(sizeAsc).DescriptionEq(descriptionEq).DescriptionNe(descriptionNe).DescriptionGt(descriptionGt).DescriptionGte(descriptionGte).DescriptionLt(descriptionLt).DescriptionLte(descriptionLte).DescriptionIn(descriptionIn).DescriptionNotin(descriptionNotin).DescriptionContains(descriptionContains).DescriptionNotcontains(descriptionNotcontains).DescriptionLike(descriptionLike).DescriptionNotlike(descriptionNotlike).DescriptionIlike(descriptionIlike).DescriptionNotilike(descriptionNotilike).DescriptionDesc(descriptionDesc).DescriptionAsc(descriptionAsc).DescriptionStatusEq(descriptionStatusEq).DescriptionStatusNe(descriptionStatusNe).DescriptionStatusGt(descriptionStatusGt).DescriptionStatusGte(descriptionStatusGte).DescriptionStatusLt(descriptionStatusLt).DescriptionStatusLte(descriptionStatusLte).DescriptionStatusIn(descriptionStatusIn).DescriptionStatusNotin(descriptionStatusNotin).DescriptionStatusContains(descriptionStatusContains).DescriptionStatusNotcontains(descriptionStatusNotcontains).DescriptionStatusLike(descriptionStatusLike).DescriptionStatusNotlike(descriptionStatusNotlike).DescriptionStatusIlike(descriptionStatusIlike).DescriptionStatusNotilike(descriptionStatusNotilike).DescriptionStatusDesc(descriptionStatusDesc).DescriptionStatusAsc(descriptionStatusAsc).DescriptionGeneratedEq(descriptionGeneratedEq).DescriptionGeneratedNe(descriptionGeneratedNe).DescriptionGeneratedGt(descriptionGeneratedGt).DescriptionGeneratedGte(descriptionGeneratedGte).DescriptionGeneratedLt(descriptionGeneratedLt).DescriptionGeneratedLte(descriptionGeneratedLte).DescriptionGeneratedIn(descriptionGeneratedIn).DescriptionGeneratedNotin(descriptionGeneratedNotin).DescriptionGeneratedContains(descriptionGeneratedContains).DescriptionGeneratedNotcontains(descriptionGeneratedNotcontains).DescriptionGeneratedDesc(descriptionGeneratedDesc).DescriptionGeneratedAsc(descriptionGeneratedAsc).SortOrderEq(sortOrderEq).SortOrderNe(sortOrderNe).SortOrderGt(sortOrderGt).SortOrderGte(sortOrderGte).SortOrderLt(sortOrderLt).SortOrderLte(sortOrderLte).SortOrderIn(sortOrderIn).SortOrderNotin(sortOrderNotin).SortOrderContains(sortOrderContains).SortOrderNotcontains(sortOrderNotcontains).SortOrderDesc(sortOrderDesc).SortOrderAsc(sortOrderAsc).MetadataVersionEq(metadataVersionEq).MetadataVersionNe(metadataVersionNe).MetadataVersionGt(metadataVersionGt).MetadataVersionGte(metadataVersionGte).MetadataVersionLt(metadataVersionLt).MetadataVersionLte(metadataVersionLte).MetadataVersionIn(metadataVersionIn).MetadataVersionNotin(metadataVersionNotin).MetadataVersionContains(metadataVersionContains).MetadataVersionNotcontains(metadataVersionNotcontains).MetadataVersionDesc(metadataVersionDesc).MetadataVersionAsc(metadataVersionAsc).AiWorkerClaimedUntilEq(aiWorkerClaimedUntilEq).AiWorkerClaimedUntilNe(aiWorkerClaimedUntilNe).AiWorkerClaimedUntilGt(aiWorkerClaimedUntilGt).AiWorkerClaimedUntilGte(aiWorkerClaimedUntilGte).AiWorkerClaimedUntilLt(aiWorkerClaimedUntilLt).AiWorkerClaimedUntilLte(aiWorkerClaimedUntilLte).AiWorkerClaimedUntilIn(aiWorkerClaimedUntilIn).AiWorkerClaimedUntilNotin(aiWorkerClaimedUntilNotin).AiWorkerClaimedUntilContains(aiWorkerClaimedUntilContains).AiWorkerClaimedUntilNotcontains(aiWorkerClaimedUntilNotcontains).AiWorkerClaimedUntilLike(aiWorkerClaimedUntilLike).AiWorkerClaimedUntilNotlike(aiWorkerClaimedUntilNotlike).AiWorkerClaimedUntilIlike(aiWorkerClaimedUntilIlike).AiWorkerClaimedUntilNotilike(aiWorkerClaimedUntilNotilike).AiWorkerClaimedUntilDesc(aiWorkerClaimedUntilDesc).AiWorkerClaimedUntilAsc(aiWorkerClaimedUntilAsc).ReferencedByMemeTagMemeIdObjectsContains(referencedByMemeTagMemeIdObjectsContains).ReferencedByMemeTagMemeIdObjectsNotcontains(referencedByMemeTagMemeIdObjectsNotcontains).ReferencedByMemeTagMemeIdObjectsDesc(referencedByMemeTagMemeIdObjectsDesc).ReferencedByMemeTagMemeIdObjectsAsc(referencedByMemeTagMemeIdObjectsAsc).Execute()



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
	referencedByMemeTagLoad := "referencedByMemeTagLoad_example" // string | load the given indirectly related objects, value is ignored (presence of key is sufficient) (optional)
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
	filenameEq := "filenameEq_example" // string | SQL = comparison (optional)
	filenameNe := "filenameNe_example" // string | SQL != comparison (optional)
	filenameGt := "filenameGt_example" // string | SQL > comparison, may not work with all column types (optional)
	filenameGte := "filenameGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	filenameLt := "filenameLt_example" // string | SQL < comparison, may not work with all column types (optional)
	filenameLte := "filenameLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	filenameIn := "filenameIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	filenameNotin := "filenameNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	filenameContains := "filenameContains_example" // string | SQL @> comparison (optional)
	filenameNotcontains := "filenameNotcontains_example" // string | SQL NOT @> comparison (optional)
	filenameLike := "filenameLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	filenameNotlike := "filenameNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	filenameIlike := "filenameIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	filenameNotilike := "filenameNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	filenameDesc := "filenameDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	filenameAsc := "filenameAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	originalNameEq := "originalNameEq_example" // string | SQL = comparison (optional)
	originalNameNe := "originalNameNe_example" // string | SQL != comparison (optional)
	originalNameGt := "originalNameGt_example" // string | SQL > comparison, may not work with all column types (optional)
	originalNameGte := "originalNameGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	originalNameLt := "originalNameLt_example" // string | SQL < comparison, may not work with all column types (optional)
	originalNameLte := "originalNameLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	originalNameIn := "originalNameIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	originalNameNotin := "originalNameNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	originalNameContains := "originalNameContains_example" // string | SQL @> comparison (optional)
	originalNameNotcontains := "originalNameNotcontains_example" // string | SQL NOT @> comparison (optional)
	originalNameLike := "originalNameLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	originalNameNotlike := "originalNameNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	originalNameIlike := "originalNameIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	originalNameNotilike := "originalNameNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	originalNameDesc := "originalNameDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	originalNameAsc := "originalNameAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	mimeTypeEq := "mimeTypeEq_example" // string | SQL = comparison (optional)
	mimeTypeNe := "mimeTypeNe_example" // string | SQL != comparison (optional)
	mimeTypeGt := "mimeTypeGt_example" // string | SQL > comparison, may not work with all column types (optional)
	mimeTypeGte := "mimeTypeGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	mimeTypeLt := "mimeTypeLt_example" // string | SQL < comparison, may not work with all column types (optional)
	mimeTypeLte := "mimeTypeLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	mimeTypeIn := "mimeTypeIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	mimeTypeNotin := "mimeTypeNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	mimeTypeContains := "mimeTypeContains_example" // string | SQL @> comparison (optional)
	mimeTypeNotcontains := "mimeTypeNotcontains_example" // string | SQL NOT @> comparison (optional)
	mimeTypeLike := "mimeTypeLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	mimeTypeNotlike := "mimeTypeNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	mimeTypeIlike := "mimeTypeIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	mimeTypeNotilike := "mimeTypeNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	mimeTypeDesc := "mimeTypeDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	mimeTypeAsc := "mimeTypeAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	sizeEq := int64(789) // int64 | SQL = comparison (optional)
	sizeNe := int64(789) // int64 | SQL != comparison (optional)
	sizeGt := int64(789) // int64 | SQL > comparison, may not work with all column types (optional)
	sizeGte := int64(789) // int64 | SQL >= comparison, may not work with all column types (optional)
	sizeLt := int64(789) // int64 | SQL < comparison, may not work with all column types (optional)
	sizeLte := int64(789) // int64 | SQL <= comparison, may not work with all column types (optional)
	sizeIn := int64(789) // int64 | SQL IN comparison, permits comma-separated values (optional)
	sizeNotin := int64(789) // int64 | SQL NOT IN comparison, permits comma-separated values (optional)
	sizeContains := int64(789) // int64 | SQL @> comparison (optional)
	sizeNotcontains := int64(789) // int64 | SQL NOT @> comparison (optional)
	sizeDesc := "sizeDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	sizeAsc := "sizeAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionEq := "descriptionEq_example" // string | SQL = comparison (optional)
	descriptionNe := "descriptionNe_example" // string | SQL != comparison (optional)
	descriptionGt := "descriptionGt_example" // string | SQL > comparison, may not work with all column types (optional)
	descriptionGte := "descriptionGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	descriptionLt := "descriptionLt_example" // string | SQL < comparison, may not work with all column types (optional)
	descriptionLte := "descriptionLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	descriptionIn := "descriptionIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	descriptionNotin := "descriptionNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	descriptionContains := "descriptionContains_example" // string | SQL @> comparison (optional)
	descriptionNotcontains := "descriptionNotcontains_example" // string | SQL NOT @> comparison (optional)
	descriptionLike := "descriptionLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionNotlike := "descriptionNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionIlike := "descriptionIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionNotilike := "descriptionNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionDesc := "descriptionDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionAsc := "descriptionAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionStatusEq := "descriptionStatusEq_example" // string | SQL = comparison (optional)
	descriptionStatusNe := "descriptionStatusNe_example" // string | SQL != comparison (optional)
	descriptionStatusGt := "descriptionStatusGt_example" // string | SQL > comparison, may not work with all column types (optional)
	descriptionStatusGte := "descriptionStatusGte_example" // string | SQL >= comparison, may not work with all column types (optional)
	descriptionStatusLt := "descriptionStatusLt_example" // string | SQL < comparison, may not work with all column types (optional)
	descriptionStatusLte := "descriptionStatusLte_example" // string | SQL <= comparison, may not work with all column types (optional)
	descriptionStatusIn := "descriptionStatusIn_example" // string | SQL IN comparison, permits comma-separated values (optional)
	descriptionStatusNotin := "descriptionStatusNotin_example" // string | SQL NOT IN comparison, permits comma-separated values (optional)
	descriptionStatusContains := "descriptionStatusContains_example" // string | SQL @> comparison (optional)
	descriptionStatusNotcontains := "descriptionStatusNotcontains_example" // string | SQL NOT @> comparison (optional)
	descriptionStatusLike := "descriptionStatusLike_example" // string | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionStatusNotlike := "descriptionStatusNotlike_example" // string | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionStatusIlike := "descriptionStatusIlike_example" // string | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionStatusNotilike := "descriptionStatusNotilike_example" // string | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	descriptionStatusDesc := "descriptionStatusDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionStatusAsc := "descriptionStatusAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionGeneratedEq := int64(789) // int64 | SQL = comparison (optional)
	descriptionGeneratedNe := int64(789) // int64 | SQL != comparison (optional)
	descriptionGeneratedGt := int64(789) // int64 | SQL > comparison, may not work with all column types (optional)
	descriptionGeneratedGte := int64(789) // int64 | SQL >= comparison, may not work with all column types (optional)
	descriptionGeneratedLt := int64(789) // int64 | SQL < comparison, may not work with all column types (optional)
	descriptionGeneratedLte := int64(789) // int64 | SQL <= comparison, may not work with all column types (optional)
	descriptionGeneratedIn := int64(789) // int64 | SQL IN comparison, permits comma-separated values (optional)
	descriptionGeneratedNotin := int64(789) // int64 | SQL NOT IN comparison, permits comma-separated values (optional)
	descriptionGeneratedContains := int64(789) // int64 | SQL @> comparison (optional)
	descriptionGeneratedNotcontains := int64(789) // int64 | SQL NOT @> comparison (optional)
	descriptionGeneratedDesc := "descriptionGeneratedDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	descriptionGeneratedAsc := "descriptionGeneratedAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	sortOrderEq := int64(789) // int64 | SQL = comparison (optional)
	sortOrderNe := int64(789) // int64 | SQL != comparison (optional)
	sortOrderGt := int64(789) // int64 | SQL > comparison, may not work with all column types (optional)
	sortOrderGte := int64(789) // int64 | SQL >= comparison, may not work with all column types (optional)
	sortOrderLt := int64(789) // int64 | SQL < comparison, may not work with all column types (optional)
	sortOrderLte := int64(789) // int64 | SQL <= comparison, may not work with all column types (optional)
	sortOrderIn := int64(789) // int64 | SQL IN comparison, permits comma-separated values (optional)
	sortOrderNotin := int64(789) // int64 | SQL NOT IN comparison, permits comma-separated values (optional)
	sortOrderContains := int64(789) // int64 | SQL @> comparison (optional)
	sortOrderNotcontains := int64(789) // int64 | SQL NOT @> comparison (optional)
	sortOrderDesc := "sortOrderDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	sortOrderAsc := "sortOrderAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	metadataVersionEq := int64(789) // int64 | SQL = comparison (optional)
	metadataVersionNe := int64(789) // int64 | SQL != comparison (optional)
	metadataVersionGt := int64(789) // int64 | SQL > comparison, may not work with all column types (optional)
	metadataVersionGte := int64(789) // int64 | SQL >= comparison, may not work with all column types (optional)
	metadataVersionLt := int64(789) // int64 | SQL < comparison, may not work with all column types (optional)
	metadataVersionLte := int64(789) // int64 | SQL <= comparison, may not work with all column types (optional)
	metadataVersionIn := int64(789) // int64 | SQL IN comparison, permits comma-separated values (optional)
	metadataVersionNotin := int64(789) // int64 | SQL NOT IN comparison, permits comma-separated values (optional)
	metadataVersionContains := int64(789) // int64 | SQL @> comparison (optional)
	metadataVersionNotcontains := int64(789) // int64 | SQL NOT @> comparison (optional)
	metadataVersionDesc := "metadataVersionDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	metadataVersionAsc := "metadataVersionAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	aiWorkerClaimedUntilEq := time.Now() // time.Time | SQL = comparison (optional)
	aiWorkerClaimedUntilNe := time.Now() // time.Time | SQL != comparison (optional)
	aiWorkerClaimedUntilGt := time.Now() // time.Time | SQL > comparison, may not work with all column types (optional)
	aiWorkerClaimedUntilGte := time.Now() // time.Time | SQL >= comparison, may not work with all column types (optional)
	aiWorkerClaimedUntilLt := time.Now() // time.Time | SQL < comparison, may not work with all column types (optional)
	aiWorkerClaimedUntilLte := time.Now() // time.Time | SQL <= comparison, may not work with all column types (optional)
	aiWorkerClaimedUntilIn := time.Now() // time.Time | SQL IN comparison, permits comma-separated values (optional)
	aiWorkerClaimedUntilNotin := time.Now() // time.Time | SQL NOT IN comparison, permits comma-separated values (optional)
	aiWorkerClaimedUntilContains := time.Now() // time.Time | SQL @> comparison (optional)
	aiWorkerClaimedUntilNotcontains := time.Now() // time.Time | SQL NOT @> comparison (optional)
	aiWorkerClaimedUntilLike := time.Now() // time.Time | SQL LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	aiWorkerClaimedUntilNotlike := time.Now() // time.Time | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	aiWorkerClaimedUntilIlike := time.Now() // time.Time | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	aiWorkerClaimedUntilNotilike := time.Now() // time.Time | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % (optional)
	aiWorkerClaimedUntilDesc := "aiWorkerClaimedUntilDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	aiWorkerClaimedUntilAsc := "aiWorkerClaimedUntilAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)
	referencedByMemeTagMemeIdObjectsContains := TODO // interface{} | SQL @> comparison (optional)
	referencedByMemeTagMemeIdObjectsNotcontains := TODO // interface{} | SQL NOT @> comparison (optional)
	referencedByMemeTagMemeIdObjectsDesc := "referencedByMemeTagMemeIdObjectsDesc_example" // string | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) (optional)
	referencedByMemeTagMemeIdObjectsAsc := "referencedByMemeTagMemeIdObjectsAsc_example" // string | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MemeAPI.GetMemes(context.Background()).Limit(limit).Offset(offset).Depth(depth).ReferencedByMemeTagLoad(referencedByMemeTagLoad).IdEq(idEq).IdNe(idNe).IdGt(idGt).IdGte(idGte).IdLt(idLt).IdLte(idLte).IdIn(idIn).IdNotin(idNotin).IdContains(idContains).IdNotcontains(idNotcontains).IdLike(idLike).IdNotlike(idNotlike).IdIlike(idIlike).IdNotilike(idNotilike).IdDesc(idDesc).IdAsc(idAsc).CreatedAtEq(createdAtEq).CreatedAtNe(createdAtNe).CreatedAtGt(createdAtGt).CreatedAtGte(createdAtGte).CreatedAtLt(createdAtLt).CreatedAtLte(createdAtLte).CreatedAtIn(createdAtIn).CreatedAtNotin(createdAtNotin).CreatedAtContains(createdAtContains).CreatedAtNotcontains(createdAtNotcontains).CreatedAtLike(createdAtLike).CreatedAtNotlike(createdAtNotlike).CreatedAtIlike(createdAtIlike).CreatedAtNotilike(createdAtNotilike).CreatedAtDesc(createdAtDesc).CreatedAtAsc(createdAtAsc).UpdatedAtEq(updatedAtEq).UpdatedAtNe(updatedAtNe).UpdatedAtGt(updatedAtGt).UpdatedAtGte(updatedAtGte).UpdatedAtLt(updatedAtLt).UpdatedAtLte(updatedAtLte).UpdatedAtIn(updatedAtIn).UpdatedAtNotin(updatedAtNotin).UpdatedAtContains(updatedAtContains).UpdatedAtNotcontains(updatedAtNotcontains).UpdatedAtLike(updatedAtLike).UpdatedAtNotlike(updatedAtNotlike).UpdatedAtIlike(updatedAtIlike).UpdatedAtNotilike(updatedAtNotilike).UpdatedAtDesc(updatedAtDesc).UpdatedAtAsc(updatedAtAsc).DeletedAtEq(deletedAtEq).DeletedAtNe(deletedAtNe).DeletedAtGt(deletedAtGt).DeletedAtGte(deletedAtGte).DeletedAtLt(deletedAtLt).DeletedAtLte(deletedAtLte).DeletedAtIn(deletedAtIn).DeletedAtNotin(deletedAtNotin).DeletedAtContains(deletedAtContains).DeletedAtNotcontains(deletedAtNotcontains).DeletedAtLike(deletedAtLike).DeletedAtNotlike(deletedAtNotlike).DeletedAtIlike(deletedAtIlike).DeletedAtNotilike(deletedAtNotilike).DeletedAtDesc(deletedAtDesc).DeletedAtAsc(deletedAtAsc).FilenameEq(filenameEq).FilenameNe(filenameNe).FilenameGt(filenameGt).FilenameGte(filenameGte).FilenameLt(filenameLt).FilenameLte(filenameLte).FilenameIn(filenameIn).FilenameNotin(filenameNotin).FilenameContains(filenameContains).FilenameNotcontains(filenameNotcontains).FilenameLike(filenameLike).FilenameNotlike(filenameNotlike).FilenameIlike(filenameIlike).FilenameNotilike(filenameNotilike).FilenameDesc(filenameDesc).FilenameAsc(filenameAsc).OriginalNameEq(originalNameEq).OriginalNameNe(originalNameNe).OriginalNameGt(originalNameGt).OriginalNameGte(originalNameGte).OriginalNameLt(originalNameLt).OriginalNameLte(originalNameLte).OriginalNameIn(originalNameIn).OriginalNameNotin(originalNameNotin).OriginalNameContains(originalNameContains).OriginalNameNotcontains(originalNameNotcontains).OriginalNameLike(originalNameLike).OriginalNameNotlike(originalNameNotlike).OriginalNameIlike(originalNameIlike).OriginalNameNotilike(originalNameNotilike).OriginalNameDesc(originalNameDesc).OriginalNameAsc(originalNameAsc).MimeTypeEq(mimeTypeEq).MimeTypeNe(mimeTypeNe).MimeTypeGt(mimeTypeGt).MimeTypeGte(mimeTypeGte).MimeTypeLt(mimeTypeLt).MimeTypeLte(mimeTypeLte).MimeTypeIn(mimeTypeIn).MimeTypeNotin(mimeTypeNotin).MimeTypeContains(mimeTypeContains).MimeTypeNotcontains(mimeTypeNotcontains).MimeTypeLike(mimeTypeLike).MimeTypeNotlike(mimeTypeNotlike).MimeTypeIlike(mimeTypeIlike).MimeTypeNotilike(mimeTypeNotilike).MimeTypeDesc(mimeTypeDesc).MimeTypeAsc(mimeTypeAsc).SizeEq(sizeEq).SizeNe(sizeNe).SizeGt(sizeGt).SizeGte(sizeGte).SizeLt(sizeLt).SizeLte(sizeLte).SizeIn(sizeIn).SizeNotin(sizeNotin).SizeContains(sizeContains).SizeNotcontains(sizeNotcontains).SizeDesc(sizeDesc).SizeAsc(sizeAsc).DescriptionEq(descriptionEq).DescriptionNe(descriptionNe).DescriptionGt(descriptionGt).DescriptionGte(descriptionGte).DescriptionLt(descriptionLt).DescriptionLte(descriptionLte).DescriptionIn(descriptionIn).DescriptionNotin(descriptionNotin).DescriptionContains(descriptionContains).DescriptionNotcontains(descriptionNotcontains).DescriptionLike(descriptionLike).DescriptionNotlike(descriptionNotlike).DescriptionIlike(descriptionIlike).DescriptionNotilike(descriptionNotilike).DescriptionDesc(descriptionDesc).DescriptionAsc(descriptionAsc).DescriptionStatusEq(descriptionStatusEq).DescriptionStatusNe(descriptionStatusNe).DescriptionStatusGt(descriptionStatusGt).DescriptionStatusGte(descriptionStatusGte).DescriptionStatusLt(descriptionStatusLt).DescriptionStatusLte(descriptionStatusLte).DescriptionStatusIn(descriptionStatusIn).DescriptionStatusNotin(descriptionStatusNotin).DescriptionStatusContains(descriptionStatusContains).DescriptionStatusNotcontains(descriptionStatusNotcontains).DescriptionStatusLike(descriptionStatusLike).DescriptionStatusNotlike(descriptionStatusNotlike).DescriptionStatusIlike(descriptionStatusIlike).DescriptionStatusNotilike(descriptionStatusNotilike).DescriptionStatusDesc(descriptionStatusDesc).DescriptionStatusAsc(descriptionStatusAsc).DescriptionGeneratedEq(descriptionGeneratedEq).DescriptionGeneratedNe(descriptionGeneratedNe).DescriptionGeneratedGt(descriptionGeneratedGt).DescriptionGeneratedGte(descriptionGeneratedGte).DescriptionGeneratedLt(descriptionGeneratedLt).DescriptionGeneratedLte(descriptionGeneratedLte).DescriptionGeneratedIn(descriptionGeneratedIn).DescriptionGeneratedNotin(descriptionGeneratedNotin).DescriptionGeneratedContains(descriptionGeneratedContains).DescriptionGeneratedNotcontains(descriptionGeneratedNotcontains).DescriptionGeneratedDesc(descriptionGeneratedDesc).DescriptionGeneratedAsc(descriptionGeneratedAsc).SortOrderEq(sortOrderEq).SortOrderNe(sortOrderNe).SortOrderGt(sortOrderGt).SortOrderGte(sortOrderGte).SortOrderLt(sortOrderLt).SortOrderLte(sortOrderLte).SortOrderIn(sortOrderIn).SortOrderNotin(sortOrderNotin).SortOrderContains(sortOrderContains).SortOrderNotcontains(sortOrderNotcontains).SortOrderDesc(sortOrderDesc).SortOrderAsc(sortOrderAsc).MetadataVersionEq(metadataVersionEq).MetadataVersionNe(metadataVersionNe).MetadataVersionGt(metadataVersionGt).MetadataVersionGte(metadataVersionGte).MetadataVersionLt(metadataVersionLt).MetadataVersionLte(metadataVersionLte).MetadataVersionIn(metadataVersionIn).MetadataVersionNotin(metadataVersionNotin).MetadataVersionContains(metadataVersionContains).MetadataVersionNotcontains(metadataVersionNotcontains).MetadataVersionDesc(metadataVersionDesc).MetadataVersionAsc(metadataVersionAsc).AiWorkerClaimedUntilEq(aiWorkerClaimedUntilEq).AiWorkerClaimedUntilNe(aiWorkerClaimedUntilNe).AiWorkerClaimedUntilGt(aiWorkerClaimedUntilGt).AiWorkerClaimedUntilGte(aiWorkerClaimedUntilGte).AiWorkerClaimedUntilLt(aiWorkerClaimedUntilLt).AiWorkerClaimedUntilLte(aiWorkerClaimedUntilLte).AiWorkerClaimedUntilIn(aiWorkerClaimedUntilIn).AiWorkerClaimedUntilNotin(aiWorkerClaimedUntilNotin).AiWorkerClaimedUntilContains(aiWorkerClaimedUntilContains).AiWorkerClaimedUntilNotcontains(aiWorkerClaimedUntilNotcontains).AiWorkerClaimedUntilLike(aiWorkerClaimedUntilLike).AiWorkerClaimedUntilNotlike(aiWorkerClaimedUntilNotlike).AiWorkerClaimedUntilIlike(aiWorkerClaimedUntilIlike).AiWorkerClaimedUntilNotilike(aiWorkerClaimedUntilNotilike).AiWorkerClaimedUntilDesc(aiWorkerClaimedUntilDesc).AiWorkerClaimedUntilAsc(aiWorkerClaimedUntilAsc).ReferencedByMemeTagMemeIdObjectsContains(referencedByMemeTagMemeIdObjectsContains).ReferencedByMemeTagMemeIdObjectsNotcontains(referencedByMemeTagMemeIdObjectsNotcontains).ReferencedByMemeTagMemeIdObjectsDesc(referencedByMemeTagMemeIdObjectsDesc).ReferencedByMemeTagMemeIdObjectsAsc(referencedByMemeTagMemeIdObjectsAsc).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.GetMemes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMemes`: ResponseWithGenericOfMeme
	fmt.Fprintf(os.Stdout, "Response from `MemeAPI.GetMemes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetMemesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | SQL LIMIT operator | 
 **offset** | **int32** | SQL OFFSET operator | 
 **depth** | **int32** | Max recursion depth for loading foreign objects; default &#x3D; 1  (0 &#x3D; recurse until graph cycle detected, 1 &#x3D; this object only, 2 &#x3D; this object + neighbours, 3 &#x3D; this object + neighbours + their neighbours... etc) | 
 **referencedByMemeTagLoad** | **string** | load the given indirectly related objects, value is ignored (presence of key is sufficient) | 
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
 **filenameEq** | **string** | SQL &#x3D; comparison | 
 **filenameNe** | **string** | SQL !&#x3D; comparison | 
 **filenameGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **filenameGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **filenameLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **filenameLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **filenameIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **filenameNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **filenameContains** | **string** | SQL @&gt; comparison | 
 **filenameNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **filenameLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **filenameNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **filenameIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **filenameNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **filenameDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **filenameAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **originalNameEq** | **string** | SQL &#x3D; comparison | 
 **originalNameNe** | **string** | SQL !&#x3D; comparison | 
 **originalNameGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **originalNameGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **originalNameLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **originalNameLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **originalNameIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **originalNameNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **originalNameContains** | **string** | SQL @&gt; comparison | 
 **originalNameNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **originalNameLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **originalNameNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **originalNameIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **originalNameNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **originalNameDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **originalNameAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **mimeTypeEq** | **string** | SQL &#x3D; comparison | 
 **mimeTypeNe** | **string** | SQL !&#x3D; comparison | 
 **mimeTypeGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **mimeTypeGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **mimeTypeLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **mimeTypeLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **mimeTypeIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **mimeTypeNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **mimeTypeContains** | **string** | SQL @&gt; comparison | 
 **mimeTypeNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **mimeTypeLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **mimeTypeNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **mimeTypeIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **mimeTypeNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **mimeTypeDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **mimeTypeAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **sizeEq** | **int64** | SQL &#x3D; comparison | 
 **sizeNe** | **int64** | SQL !&#x3D; comparison | 
 **sizeGt** | **int64** | SQL &gt; comparison, may not work with all column types | 
 **sizeGte** | **int64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **sizeLt** | **int64** | SQL &lt; comparison, may not work with all column types | 
 **sizeLte** | **int64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **sizeIn** | **int64** | SQL IN comparison, permits comma-separated values | 
 **sizeNotin** | **int64** | SQL NOT IN comparison, permits comma-separated values | 
 **sizeContains** | **int64** | SQL @&gt; comparison | 
 **sizeNotcontains** | **int64** | SQL NOT @&gt; comparison | 
 **sizeDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **sizeAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **descriptionEq** | **string** | SQL &#x3D; comparison | 
 **descriptionNe** | **string** | SQL !&#x3D; comparison | 
 **descriptionGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **descriptionGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **descriptionLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **descriptionLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **descriptionIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **descriptionNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **descriptionContains** | **string** | SQL @&gt; comparison | 
 **descriptionNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **descriptionLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **descriptionAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **descriptionStatusEq** | **string** | SQL &#x3D; comparison | 
 **descriptionStatusNe** | **string** | SQL !&#x3D; comparison | 
 **descriptionStatusGt** | **string** | SQL &gt; comparison, may not work with all column types | 
 **descriptionStatusGte** | **string** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **descriptionStatusLt** | **string** | SQL &lt; comparison, may not work with all column types | 
 **descriptionStatusLte** | **string** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **descriptionStatusIn** | **string** | SQL IN comparison, permits comma-separated values | 
 **descriptionStatusNotin** | **string** | SQL NOT IN comparison, permits comma-separated values | 
 **descriptionStatusContains** | **string** | SQL @&gt; comparison | 
 **descriptionStatusNotcontains** | **string** | SQL NOT @&gt; comparison | 
 **descriptionStatusLike** | **string** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionStatusNotlike** | **string** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionStatusIlike** | **string** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionStatusNotilike** | **string** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **descriptionStatusDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **descriptionStatusAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **descriptionGeneratedEq** | **int64** | SQL &#x3D; comparison | 
 **descriptionGeneratedNe** | **int64** | SQL !&#x3D; comparison | 
 **descriptionGeneratedGt** | **int64** | SQL &gt; comparison, may not work with all column types | 
 **descriptionGeneratedGte** | **int64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **descriptionGeneratedLt** | **int64** | SQL &lt; comparison, may not work with all column types | 
 **descriptionGeneratedLte** | **int64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **descriptionGeneratedIn** | **int64** | SQL IN comparison, permits comma-separated values | 
 **descriptionGeneratedNotin** | **int64** | SQL NOT IN comparison, permits comma-separated values | 
 **descriptionGeneratedContains** | **int64** | SQL @&gt; comparison | 
 **descriptionGeneratedNotcontains** | **int64** | SQL NOT @&gt; comparison | 
 **descriptionGeneratedDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **descriptionGeneratedAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **sortOrderEq** | **int64** | SQL &#x3D; comparison | 
 **sortOrderNe** | **int64** | SQL !&#x3D; comparison | 
 **sortOrderGt** | **int64** | SQL &gt; comparison, may not work with all column types | 
 **sortOrderGte** | **int64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **sortOrderLt** | **int64** | SQL &lt; comparison, may not work with all column types | 
 **sortOrderLte** | **int64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **sortOrderIn** | **int64** | SQL IN comparison, permits comma-separated values | 
 **sortOrderNotin** | **int64** | SQL NOT IN comparison, permits comma-separated values | 
 **sortOrderContains** | **int64** | SQL @&gt; comparison | 
 **sortOrderNotcontains** | **int64** | SQL NOT @&gt; comparison | 
 **sortOrderDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **sortOrderAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **metadataVersionEq** | **int64** | SQL &#x3D; comparison | 
 **metadataVersionNe** | **int64** | SQL !&#x3D; comparison | 
 **metadataVersionGt** | **int64** | SQL &gt; comparison, may not work with all column types | 
 **metadataVersionGte** | **int64** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **metadataVersionLt** | **int64** | SQL &lt; comparison, may not work with all column types | 
 **metadataVersionLte** | **int64** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **metadataVersionIn** | **int64** | SQL IN comparison, permits comma-separated values | 
 **metadataVersionNotin** | **int64** | SQL NOT IN comparison, permits comma-separated values | 
 **metadataVersionContains** | **int64** | SQL @&gt; comparison | 
 **metadataVersionNotcontains** | **int64** | SQL NOT @&gt; comparison | 
 **metadataVersionDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **metadataVersionAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **aiWorkerClaimedUntilEq** | **time.Time** | SQL &#x3D; comparison | 
 **aiWorkerClaimedUntilNe** | **time.Time** | SQL !&#x3D; comparison | 
 **aiWorkerClaimedUntilGt** | **time.Time** | SQL &gt; comparison, may not work with all column types | 
 **aiWorkerClaimedUntilGte** | **time.Time** | SQL &gt;&#x3D; comparison, may not work with all column types | 
 **aiWorkerClaimedUntilLt** | **time.Time** | SQL &lt; comparison, may not work with all column types | 
 **aiWorkerClaimedUntilLte** | **time.Time** | SQL &lt;&#x3D; comparison, may not work with all column types | 
 **aiWorkerClaimedUntilIn** | **time.Time** | SQL IN comparison, permits comma-separated values | 
 **aiWorkerClaimedUntilNotin** | **time.Time** | SQL NOT IN comparison, permits comma-separated values | 
 **aiWorkerClaimedUntilContains** | **time.Time** | SQL @&gt; comparison | 
 **aiWorkerClaimedUntilNotcontains** | **time.Time** | SQL NOT @&gt; comparison | 
 **aiWorkerClaimedUntilLike** | **time.Time** | SQL LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **aiWorkerClaimedUntilNotlike** | **time.Time** | SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with % | 
 **aiWorkerClaimedUntilIlike** | **time.Time** | SQL ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **aiWorkerClaimedUntilNotilike** | **time.Time** | SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with % | 
 **aiWorkerClaimedUntilDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **aiWorkerClaimedUntilAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 
 **referencedByMemeTagMemeIdObjectsContains** | [**interface{}**](interface{}.md) | SQL @&gt; comparison | 
 **referencedByMemeTagMemeIdObjectsNotcontains** | [**interface{}**](interface{}.md) | SQL NOT @&gt; comparison | 
 **referencedByMemeTagMemeIdObjectsDesc** | **string** | SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient) | 
 **referencedByMemeTagMemeIdObjectsAsc** | **string** | SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient) | 

### Return type

[**ResponseWithGenericOfMeme**](ResponseWithGenericOfMeme.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchMeme

> ResponseWithGenericOfMeme PatchMeme(ctx, primaryKey).Meme(meme).Depth(depth).Execute()



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
	meme := *openapiclient.NewMeme() // Meme | 
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MemeAPI.PatchMeme(context.Background(), primaryKey).Meme(meme).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.PatchMeme``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PatchMeme`: ResponseWithGenericOfMeme
	fmt.Fprintf(os.Stdout, "Response from `MemeAPI.PatchMeme`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchMemeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **meme** | [**Meme**](Meme.md) |  | 
 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfMeme**](ResponseWithGenericOfMeme.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostMemes

> ResponseWithGenericOfMeme PostMemes(ctx).Meme(meme).Depth(depth).Execute()



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
	meme := []openapiclient.Meme{*openapiclient.NewMeme()} // []Meme | 
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MemeAPI.PostMemes(context.Background()).Meme(meme).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.PostMemes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PostMemes`: ResponseWithGenericOfMeme
	fmt.Fprintf(os.Stdout, "Response from `MemeAPI.PostMemes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostMemesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **meme** | [**[]Meme**](Meme.md) |  | 
 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfMeme**](ResponseWithGenericOfMeme.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostMemesAiWorkerClaim

> ResponseWithGenericOfMeme PostMemesAiWorkerClaim(ctx, primaryKey).MemeAiWorkerClaimRequest(memeAiWorkerClaimRequest).Depth(depth).Execute()



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
	memeAiWorkerClaimRequest := *openapiclient.NewMemeAiWorkerClaimRequest() // MemeAiWorkerClaimRequest | 
	depth := int64(789) // int64 | Query parameter depth (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MemeAPI.PostMemesAiWorkerClaim(context.Background(), primaryKey).MemeAiWorkerClaimRequest(memeAiWorkerClaimRequest).Depth(depth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MemeAPI.PostMemesAiWorkerClaim``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PostMemesAiWorkerClaim`: ResponseWithGenericOfMeme
	fmt.Fprintf(os.Stdout, "Response from `MemeAPI.PostMemesAiWorkerClaim`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**primaryKey** | **string** | Path parameter primaryKey | 

### Other Parameters

Other parameters are passed through a pointer to a apiPostMemesAiWorkerClaimRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **memeAiWorkerClaimRequest** | [**MemeAiWorkerClaimRequest**](MemeAiWorkerClaimRequest.md) |  | 
 **depth** | **int64** | Query parameter depth | 

### Return type

[**ResponseWithGenericOfMeme**](ResponseWithGenericOfMeme.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

