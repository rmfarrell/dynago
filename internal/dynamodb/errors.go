// Package dynamodb maps information from dynamo itself, such as error tables and formats.
// This is an internal package to allow fast iteration on the abstraction
// without having to commit to the interface.
package dynamodb

import (
	"github.com/rmfarrell/dynago/internal/codes"
)

type ErrorConfig struct {
	AmazonCode     string
	ExpectedStatus int
	MappedError    codes.ErrorCode
}

// This variable is mostly exposed so that we can document how errors are mapped
var MappedErrors = []ErrorConfig{
	{"ConditionalCheckFailedException", 400, codes.ErrorConditionFailed},
	{"InternalFailure", 500, codes.ErrorInternalFailure},
	{"InternalServerError", 500, codes.ErrorInternalFailure},
	{"IncompleteSignature", 400, codes.ErrorAuth},
	{"IncompleteSignatureException", 400, codes.ErrorAuth},
	{"InvalidParameterCombination", 400, codes.ErrorInvalidParameter},
	{"InvalidParameterValue", 400, codes.ErrorInvalidParameter},
	{"InvalidQueryParameter", 400, codes.ErrorInvalidParameter},
	{"InvalidSignatureException", 400, codes.ErrorAuth},
	{"ItemCollectionSizeLimitExceededException", 400, codes.ErrorCollectionSizeExceeded},
	{"MalformedQueryString", 404, codes.ErrorInvalidParameter},
	{"MissingAction", 400, codes.ErrorInvalidParameter},
	{"MissingAuthenticationToken", 403, codes.ErrorAuth},
	{"MissingParameter", 400, codes.ErrorInvalidParameter},
	{"OptInRequired", 403, codes.ErrorAuth},
	{"ProvisionedThroughputExceededException", 400, codes.ErrorThroughputExceeded},
	{"RequestExpired", 400, codes.ErrorAuth},
	{"ResourceInUseException", 400, codes.ErrorResourceInUse},
	{"ResourceNotFoundException", 400, codes.ErrorNotFound},
	{"ServiceUnavailable", 503, codes.ErrorServiceUnavailable},
	{"ServiceUnavailableException", 503, codes.ErrorServiceUnavailable},
	{"ThrottlingException", 400, codes.ErrorThrottling},
	{"UnrecognizedClientException", 400, codes.ErrorAuth},
	{"ValidationError", 400, codes.ErrorInvalidParameter},
	{"ValidationException", 400, codes.ErrorInvalidParameter},

	// DynamoDB Streams
	{"ExpiredIteratorException", 400, codes.ErrorExpiredIterator},
	{"LimitExceededException", 400, codes.ErrorThrottling},
	{"TrimmedDataAccessException", 400, codes.ErrorTrimmedData},
}
