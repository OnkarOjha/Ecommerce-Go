package utils

const (
	HTTP_BAD_REQUEST                     int = 400
	HTTP_UNAUTHORIZED                    int = 401
	HTTP_PAYMENT_REQUIRED                int = 402
	HTTP_FORBIDDEN                       int = 403
	HTTP_NOT_FOUND                       int = 404
	HTTP_METHOD_NOT_ALLOWED              int = 405
	HTTP_NOT_ACCEPTABLE                  int = 406
	HTTP_PROXY_AUTHENTICATION_REQUIRED   int = 407
	HTTP_REQUEST_TIMEOUT                 int = 408
	HTTP_CONFLICT                        int = 409
	HTTP_GONE                            int = 410
	HTTP_LENGTH_REQUIRED                 int = 411
	HTTP_PRECONDITION_FAILED             int = 412
	HTTP_PAYLOAD_TOO_LARGE               int = 413
	HTTP_URI_TOO_LONG                    int = 414
	HTTP_UNSUPPORTED_MEDIA_TYPE          int = 415
	HTTP_RANGE_NOT_SATISFIABLE           int = 416
	HTTP_EXPECTATION_FAILED              int = 417
	HTTP_TEAPOT                          int = 418
	HTTP_MISDIRECTED_REQUEST             int = 421
	HTTP_UNPROCESSABLE_ENTITY            int = 422
	HTTP_LOCKED                          int = 423
	HTTP_FAILED_DEPENDENCY               int = 424
	HTTP_UPGRADE_REQUIRED                int = 426
	HTTP_PRECONDITION_REQUIRED           int = 428
	HTTP_TOO_MANY_REQUESTS               int = 429
	HTTP_REQUEST_HEADER_FIELDS_TOO_LARGE int = 431
	HTTP_UNAVAILABLE_FOR_LEGAL_REASONS   int = 451
	HTTP_INTERNAL_SERVER_ERROR           int = 500
	HTTP_NOT_IMPLEMENTED                 int = 501
	HTTP_BAD_GATEWAY                     int = 502
	HTTP_SERVICE_UNAVAILABLE             int = 503
	HTTP_GATEWAY_TIMEOUT                 int = 504
	HTTP_HTTP_VERSION_NOT_SUPPORTED      int = 505
	HTTP_VARIANT_ALSO_NEGOTIATES         int = 506
	HTTP_INSUFFICIENT_STORAGE            int = 507
	HTTP_LOOP_DETECTED                   int = 508
	HTTP_NOT_EXTENDED                    int = 510
	HTTP_NETWORK_AUTHENTICATION_REQUIRED int = 511
	HTTP_OK                              int = 200
	HTTP_NO_CONTENT                      int = 204
)
