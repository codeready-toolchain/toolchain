module github.com/codeready-toolchain/toolchain-common

go 1.15

require (
	github.com/codeready-toolchain/api v0.0.0-20200805071634-c62858ce3204
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emicklei/go-restful v2.12.0+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.7 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/lestrrat-go/jwx v0.9.0
	github.com/magiconair/properties v1.8.1
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/library-go v0.0.0-20191121124438-7c776f7cc17a
	github.com/operator-framework/operator-sdk v0.19.2
	github.com/pkg/errors v0.9.1
	github.com/redhat-cop/operator-utils v0.3.4
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.5.1
	gopkg.in/h2non/gock.v1 v1.0.14
	gopkg.in/square/go-jose.v2 v2.3.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200427153329-656914f816f9 // indirect
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200204173128-addea2498afe // avoids case-insensitive import collision: "github.com/googleapis/gnostic/openapiv2" and "github.com/googleapis/gnostic/OpenAPIv2"
)

replace github.com/codeready-toolchain/api => github.com/xcoulon/api v0.0.0-20200819120629-173f1a6913c5
