package s3api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/pb/s3_pb"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb"
	. "github.com/chrislusf/seaweedfs/weed/s3api/s3_constants"
	"github.com/chrislusf/seaweedfs/weed/s3api/s3err"
	"github.com/chrislusf/seaweedfs/weed/security"
	"github.com/chrislusf/seaweedfs/weed/util"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type S3ApiServerOption struct {
	Filer                     pb.ServerAddress
	Port                      int
	Config                    string
	DomainName                string
	BucketsPath               string
	GrpcDialOption            grpc.DialOption
	AllowEmptyFolder          bool
	AllowDeleteBucketNotEmpty bool
	LocalFilerSocket          *string
}

type S3ApiServer struct {
	s3_pb.UnimplementedSeaweedS3Server
	option         *S3ApiServerOption
	iam            *IdentityAccessManagement
	cb             *CircuitBreaker
	randomClientId int32
	filerGuard     *security.Guard
	client         *http.Client
}

func NewS3ApiServer(router *mux.Router, option *S3ApiServerOption) (s3ApiServer *S3ApiServer, err error) {
	v := util.GetViper()
	signingKey := v.GetString("jwt.filer_signing.key")
	v.SetDefault("jwt.filer_signing.expires_after_seconds", 10)
	expiresAfterSec := v.GetInt("jwt.filer_signing.expires_after_seconds")

	readSigningKey := v.GetString("jwt.filer_signing.read.key")
	v.SetDefault("jwt.filer_signing.read.expires_after_seconds", 60)
	readExpiresAfterSec := v.GetInt("jwt.filer_signing.read.expires_after_seconds")

	s3ApiServer = &S3ApiServer{
		option:         option,
		iam:            NewIdentityAccessManagement(option),
		randomClientId: util.RandomInt32(),
		filerGuard:     security.NewGuard([]string{}, signingKey, expiresAfterSec, readSigningKey, readExpiresAfterSec),
		cb:             NewCircuitBreaker(option),
	}
	if option.LocalFilerSocket == nil || *option.LocalFilerSocket == "" {
		s3ApiServer.client = &http.Client{Transport: &http.Transport{
			MaxIdleConns:        1024,
			MaxIdleConnsPerHost: 1024,
		}}
	} else {
		s3ApiServer.client = &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", *option.LocalFilerSocket)
				},
			},
		}
	}

	s3ApiServer.registerRouter(router)

	go s3ApiServer.subscribeMetaEvents("s3", filer.DirectoryEtcRoot, time.Now().UnixNano())
	return s3ApiServer, nil
}

func (s3a *S3ApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()

	// Readiness Probe
	apiRouter.Methods("GET").Path("/status").HandlerFunc(s3a.StatusHandler)

	var routers []*mux.Router
	if s3a.option.DomainName != "" {
		domainNames := strings.Split(s3a.option.DomainName, ",")
		for _, domainName := range domainNames {
			routers = append(routers, apiRouter.Host(
				fmt.Sprintf("%s.%s:%d", "{bucket:.+}", domainName, s3a.option.Port)).Subrouter())
			routers = append(routers, apiRouter.Host(
				fmt.Sprintf("%s.%s", "{bucket:.+}", domainName)).Subrouter())
		}
	}
	routers = append(routers, apiRouter.PathPrefix("/{bucket}").Subrouter())

	for _, bucket := range routers {

		// each case should follow the next rule:
		// - requesting object with query must precede any other methods
		// - requesting object must precede any methods with buckets
		// - requesting bucket with query must precede raw methods with buckets
		// - requesting bucket must be processed in the end

		// objects with query

		// CopyObjectPart
		bucket.Methods("PUT").Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", `.*?(\/|%2F).*?`).HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.CopyObjectPartHandler, ACTION_WRITE)), "PUT")).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// PutObjectPart
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectPartHandler, ACTION_WRITE)), "PUT")).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// CompleteMultipartUpload
		bucket.Methods("POST").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.CompleteMultipartUploadHandler, ACTION_WRITE)), "POST")).Queries("uploadId", "{uploadId:.*}")
		// NewMultipartUpload
		bucket.Methods("POST").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.NewMultipartUploadHandler, ACTION_WRITE)), "POST")).Queries("uploads", "")
		// AbortMultipartUpload
		bucket.Methods("DELETE").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.AbortMultipartUploadHandler, ACTION_WRITE)), "DELETE")).Queries("uploadId", "{uploadId:.*}")
		// ListObjectParts
		bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.ListObjectPartsHandler, ACTION_READ)), "GET")).Queries("uploadId", "{uploadId:.*}")
		// ListMultipartUploads
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.ListMultipartUploadsHandler, ACTION_READ)), "GET")).Queries("uploads", "")

		// GetObjectTagging
		bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetObjectTaggingHandler, ACTION_READ)), "GET")).Queries("tagging", "")
		// PutObjectTagging
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectTaggingHandler, ACTION_TAGGING)), "PUT")).Queries("tagging", "")
		// DeleteObjectTagging
		bucket.Methods("DELETE").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteObjectTaggingHandler, ACTION_TAGGING)), "DELETE")).Queries("tagging", "")

		// PutObjectACL
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectAclHandler, ACTION_WRITE)), "PUT")).Queries("acl", "")
		// PutObjectRetention
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectRetentionHandler, ACTION_WRITE)), "PUT")).Queries("retention", "")
		// PutObjectLegalHold
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectLegalHoldHandler, ACTION_WRITE)), "PUT")).Queries("legal-hold", "")
		// PutObjectLockConfiguration
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectLockConfigurationHandler, ACTION_WRITE)), "PUT")).Queries("object-lock", "")

		// GetObjectACL
		bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetObjectAclHandler, ACTION_READ)), "GET")).Queries("acl", "")

		// objects with query

		// raw objects

		// HeadObject
		bucket.Methods("HEAD").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.HeadObjectHandler, ACTION_READ)), "GET"))

		// GetObject, but directory listing is not supported
		bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetObjectHandler, ACTION_READ)), "GET"))

		// CopyObject
		bucket.Methods("PUT").Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.CopyObjectHandler, ACTION_WRITE)), "COPY"))
		// PutObject
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutObjectHandler, ACTION_WRITE)), "PUT"))
		// DeleteObject
		bucket.Methods("DELETE").Path("/{object:.+}").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteObjectHandler, ACTION_WRITE)), "DELETE"))

		// raw objects

		// buckets with query

		// DeleteMultipleObjects
		bucket.Methods("POST").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteMultipleObjectsHandler, ACTION_WRITE)), "DELETE")).Queries("delete", "")

		// GetBucketACL
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketAclHandler, ACTION_READ)), "GET")).Queries("acl", "")
		// PutBucketACL
		bucket.Methods("PUT").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutBucketAclHandler, ACTION_WRITE)), "PUT")).Queries("acl", "")

		// GetBucketPolicy
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketPolicyHandler, ACTION_READ)), "GET")).Queries("policy", "")
		// PutBucketPolicy
		bucket.Methods("PUT").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutBucketPolicyHandler, ACTION_WRITE)), "PUT")).Queries("policy", "")
		// DeleteBucketPolicy
		bucket.Methods("DELETE").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteBucketPolicyHandler, ACTION_WRITE)), "DELETE")).Queries("policy", "")

		// GetBucketCors
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketCorsHandler, ACTION_READ)), "GET")).Queries("cors", "")
		// PutBucketCors
		bucket.Methods("PUT").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutBucketCorsHandler, ACTION_WRITE)), "PUT")).Queries("cors", "")
		// DeleteBucketCors
		bucket.Methods("DELETE").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteBucketCorsHandler, ACTION_WRITE)), "DELETE")).Queries("cors", "")

		// GetBucketLifecycleConfiguration
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketLifecycleConfigurationHandler, ACTION_READ)), "GET")).Queries("lifecycle", "")
		// PutBucketLifecycleConfiguration
		bucket.Methods("PUT").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PutBucketLifecycleConfigurationHandler, ACTION_WRITE)), "PUT")).Queries("lifecycle", "")
		// DeleteBucketLifecycleConfiguration
		bucket.Methods("DELETE").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteBucketLifecycleHandler, ACTION_WRITE)), "DELETE")).Queries("lifecycle", "")

		// GetBucketLocation
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketLocationHandler, ACTION_READ)), "GET")).Queries("location", "")

		// GetBucketRequestPayment
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.GetBucketRequestPaymentHandler, ACTION_READ)), "GET")).Queries("requestPayment", "")

		// ListObjectsV2
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.ListObjectsV2Handler, ACTION_LIST)), "LIST")).Queries("list-type", "2")

		// buckets with query

		// raw buckets

		// PostPolicy
		bucket.Methods("POST").HeadersRegexp("Content-Type", "multipart/form-data*").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.PostPolicyBucketHandler, ACTION_WRITE)), "POST"))

		// HeadBucket
		bucket.Methods("HEAD").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.HeadBucketHandler, ACTION_READ)), "GET"))

		// PutBucket
		bucket.Methods("PUT").HandlerFunc(track(s3a.PutBucketHandler, "PUT"))
		// DeleteBucket
		bucket.Methods("DELETE").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.DeleteBucketHandler, ACTION_WRITE)), "DELETE"))

		// ListObjectsV1 (Legacy)
		bucket.Methods("GET").HandlerFunc(track(s3a.iam.Auth(s3a.cb.Limit(s3a.ListObjectsV1Handler, ACTION_LIST)), "LIST"))

		// raw buckets

	}

	// ListBuckets
	apiRouter.Methods("GET").Path("/").HandlerFunc(track(s3a.ListBucketsHandler, "LIST"))

	// Redirect iam request from port 29000 to 8111
	apiRouter.Methods("POST").Path("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		// verify whether the request type is "iam" or "s3". For a typical Authorization:
		// "AWS4-HMAC-SHA256 Credential=197FSAQ7HHTA48X64O3A/20220420/test1/iam/aws4_request, SignedHeaders=content-type;
		// host;x-amz-date, Signature=6757dc6b3d7534d67e17842760310e99ee695408497f6edc4fdb84770c252dc8",
		// the request type is "iam"
		s := strings.Split(r.Header["Authorization"][0], "Credential=")
		if len(s) < 2 {
			return false
		}
		s = strings.Split(s[1], ",")
		if len(s) < 2 {
			return false
		}
		s = strings.Split(s[0], "/")
		if len(s) < 5 {
			return false
		}
		return s[3] == "iam"
	}).HandlerFunc(s3a.RedirectPortHandler)

	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(s3err.NotFoundHandler)

}

// RedirectPortHandler redirects s3 request of port 29000 to port 8111
func (s3a *S3ApiServer) RedirectPortHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("RedirectPortHandler")

	// get new url
	cli := &http.Client{}
	reqUrl := "http://" + strings.TrimRight(req.Host, strconv.Itoa(s3a.option.Port)) + "8111" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		reqUrl += "?" + req.URL.RawQuery
	}

	// get new http request
	body, err := io.ReadAll(req.Body)
	req2, err := http.NewRequest(req.Method, reqUrl, strings.NewReader(string(body)))
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to make a new http request: %v", err))
		return
	}
	req2.Header = req.Header // get Header of original request

	// redirect request
	rep2, err := cli.Do(req2)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to redirect request: %v", err))
		return
	}
	defer rep2.Body.Close()

	// get response
	repBody, err := io.ReadAll(rep2.Body)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to get response: %v", err))
		return
	}

	// get status code
	w.WriteHeader(rep2.StatusCode)

	for k, v := range rep2.Header {
		w.Header().Set(k, v[0])
	}
	io.WriteString(w, string(repBody[:]))
}
