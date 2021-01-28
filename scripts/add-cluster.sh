#!/usr/bin/env bash

set -e

user_help () {
    echo "Creates ToolchainCluster"
    echo "options:"
    echo "-t, --type            joining cluster type (host or member)"
    echo "-tn, --type-name      the type name of the joining cluster (host, member or e2e)"
    echo "-mn, --member-ns      namespace where member-operator is running"
    echo "-hn, --host-ns        namespace where host-operator is running"
    echo "-s,  --single-cluster running both operators on single cluster"
    echo "-ms, --multi-ns       use this flag for subsequent member-operators running in a different namespaces on the same cluster to avoid resource naming collisions (mainly for testing purposes)"
    echo "-kc, --kube-config    kubeconfig for managing multiple clusters"
    echo "-sc, --sandbox-config sandbox config file for managing Dev Sandbox instance"
    echo "-le, --lets-encrypt   use let's encrypt certificate"
    exit 0
}

login_to_cluster() {
    if [[ ${SINGLE_CLUSTER} != "true" ]]; then
      if [[ -z ${KUBECONFIG} ]] && [[ -z ${SANDBOX_CONFIG} ]]; then
        echo "Please specify the path to kube config file using the parameter --kube-config"
        echo "or specify SA tokens to be used when reaching operators using the parameters --host-token and --member-token"
      elif [[ -n ${KUBECONFIG} ]]; then
        oc config use-context "$1-admin"
      else
        REGISTER_SERVER_API=`yq -r .$1.serverAPI ${SANDBOX_CONFIG}`
        REGISTER_SA_TOKEN=`yq -r .$1.tokens.registerCluster ${SANDBOX_CONFIG}`
        OC_ADDITIONAL_PARAMS="--token=${REGISTER_SA_TOKEN} --server=${REGISTER_SERVER_API}"
      fi
    fi
}

create_service_account() {
# we need to delete the bindings since we cannot change the roleRef of the existing bindings
if [[ -n `oc get rolebinding ${SA_NAME} 2>/dev/null` ]]; then
    oc delete rolebinding ${SA_NAME} -n ${OPERATOR_NS} ${OC_ADDITIONAL_PARAMS}
fi
cat <<EOF | oc apply ${OC_ADDITIONAL_PARAMS} -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
rules:
- apiGroups:
  - toolchain.dev.openshift.com
  resources:
  - "bannedusers"
  - "changetierrequests"
  - "hostoperatorconfigs"
  - "idlers"
  - "masteruserrecords"
  - "memberstatuses"
  - "notifications"
  - "nstemplatesets"
  - "nstemplatetiers"
  - "registrationservices"
  - "templateupdaterequests"
  - "tiertemplates"
  - "toolchainclusters"
  - "toolchainstatuses"
  - "useraccounts"
  - "usersignups"
  verbs:
  - "*"
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
subjects:
- kind: ServiceAccount
  name: ${SA_NAME}
roleRef:
  kind: Role
  name: ${SA_NAME}
  apiGroup: rbac.authorization.k8s.io
EOF
}

create_service_account_e2e() {
ROLE_NAME=`oc get Roles -n ${OPERATOR_NS} -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' ${OC_ADDITIONAL_PARAMS} | grep "^toolchain-${JOINING_CLUSTER_TYPE}-operator\.v"`
if [[ -z ${ROLE_NAME} ]]; then
    echo "Role that would have a prefix 'toolchain-${JOINING_CLUSTER_TYPE}-operator.v' wasn't found - available roles are:"
    echo `oc get Roles -n ${OPERATOR_NS} ${OC_ADDITIONAL_PARAMS}`
    exit 1
fi
echo "using Role ${ROLE_NAME}"
CLUSTER_ROLE_NAME=`oc get ClusterRoles -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' ${OC_ADDITIONAL_PARAMS} | grep "^toolchain-${JOINING_CLUSTER_TYPE}-operator\.v" -m 1`
if [[ -z ${CLUSTER_ROLE_NAME} ]]; then
    echo "ClusterRole that would have a prefix 'toolchain-${JOINING_CLUSTER_TYPE}-operator.v' wasn't found - available ClusterRoles are:"
    echo `oc get ClusterRoles ${OC_ADDITIONAL_PARAMS}`
    exit 1
fi
echo "using ClusterRole ${CLUSTER_ROLE_NAME}"
cat <<EOF | oc apply ${OC_ADDITIONAL_PARAMS} -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
subjects:
- kind: ServiceAccount
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
roleRef:
  kind: Role
  name: ${ROLE_NAME}
  apiGroup: rbac.authorization.k8s.io
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${SA_NAME}-${OPERATOR_NS}
subjects:
- kind: ServiceAccount
  name: ${SA_NAME}
  namespace: ${OPERATOR_NS}
roleRef:
  kind: ClusterRole
  name: ${CLUSTER_ROLE_NAME}
  apiGroup: rbac.authorization.k8s.io
EOF
}

if [[ $# -lt 2 ]]
then
    user_help
fi

while test $# -gt 0; do
       case "$1" in
            -h|--help)
                user_help
                ;;
            -t|--type)
                shift
                JOINING_CLUSTER_TYPE=$1
                shift
                ;;
            -tn|--type-name)
                shift
                JOINING_CLUSTER_TYPE_NAME=$1
                shift
                ;;
            -mn|--member-ns)
                shift
                MEMBER_OPERATOR_NS=$1
                shift
                ;;
            -hn|--host-ns)
                shift
                HOST_OPERATOR_NS=$1
                shift
                ;;
            -kc|--kube-config)
                shift
                KUBECONFIG=$1
                shift
                ;;
            -sc|--sandbox-config)
                shift
                SANDBOX_CONFIG=$1
                shift
                ;;
            -s|--single-cluster)
                SINGLE_CLUSTER=true
                shift
                ;;
            -ms|--multi-ns)
                MULTI_NS=true
                shift
                ;;
            -le|--lets-encrypt)
                LETS_ENCRYPT=true
                shift
                ;;
            *)
               echo "$1 is not a recognized flag!"
               user_help
               exit -1
               ;;
      esac
done

CLUSTER_JOIN_TO="host"

if [[ -n ${SANDBOX_CONFIG} ]]; then
    HOST_OPERATOR_NS=`yq -r .host.namespace ${SANDBOX_CONFIG}`
    MEMBER_OPERATOR_NS=`yq -r .member.namespace ${SANDBOX_CONFIG}`
fi

# We need this to configurable to work with dynamic namespaces from end to end tests
OPERATOR_NS=${MEMBER_OPERATOR_NS}
CLUSTER_JOIN_TO_OPERATOR_NS=${HOST_OPERATOR_NS}
if [[ ${JOINING_CLUSTER_TYPE} == "host" ]]; then
  CLUSTER_JOIN_TO="member"
  OPERATOR_NS=${HOST_OPERATOR_NS}
  CLUSTER_JOIN_TO_OPERATOR_NS=${MEMBER_OPERATOR_NS}
fi
JOINING_CLUSTER_TYPE_NAME=${JOINING_CLUSTER_TYPE_NAME:-${JOINING_CLUSTER_TYPE}}

# This is using default values i.e. toolchain-member-operator or toolchain-host-operator for local setup
if [[ ${OPERATOR_NS} == "" &&  ${CLUSTER_JOIN_TO_OPERATOR_NS} == "" ]]; then
  OPERATOR_NS=toolchain-${JOINING_CLUSTER_TYPE}-operator
  CLUSTER_JOIN_TO_OPERATOR_NS=toolchain-${CLUSTER_JOIN_TO}-operator
fi

echo ${OPERATOR_NS}
echo ${CLUSTER_JOIN_TO_OPERATOR_NS}

login_to_cluster ${JOINING_CLUSTER_TYPE}

if [[ ${JOINING_CLUSTER_TYPE_NAME} != "e2e" ]]; then
    SA_NAME="toolchaincluster-${JOINING_CLUSTER_TYPE}"
    if [[ ${MULTI_NS} == "true" ]]; then
      SA_NAME="${SA_NAME}-${OPERATOR_NS}"
    fi
    create_service_account
else
    SA_NAME="e2e-service-account"
    if [[ ${MULTI_NS} == "true" ]]; then
      SA_NAME="${SA_NAME}-${OPERATOR_NS}"
    fi
    create_service_account_e2e
fi

echo "Getting ${JOINING_CLUSTER_TYPE} SA token"
SA_SECRET=`oc get sa ${SA_NAME} -n ${OPERATOR_NS} -o json ${OC_ADDITIONAL_PARAMS} | jq -r .secrets[].name | grep token`
SA_TOKEN=`oc get secret ${SA_SECRET} -n ${OPERATOR_NS}  -o json ${OC_ADDITIONAL_PARAMS} | jq -r '.data["token"]' | base64 --decode`
if [[ ${LETS_ENCRYPT} == "true" ]]; then
    SA_CA_CRT=`curl https://letsencrypt.org/certs/lets-encrypt-r3.pem | base64 -w 0`
else
    SA_CA_CRT=`oc get secret ${SA_SECRET} -n ${OPERATOR_NS} -o json ${OC_ADDITIONAL_PARAMS} | jq -r '.data["ca.crt"]'`
fi

if [[ -n ${SANDBOX_CONFIG} ]]; then
    API_ENDPOINT=`yq -r .${JOINING_CLUSTER_TYPE}.serverAPI ${SANDBOX_CONFIG}`
    JOINING_CLUSTER_NAME=`yq -r .${JOINING_CLUSTER_TYPE}.serverName ${SANDBOX_CONFIG}`

    login_to_cluster ${CLUSTER_JOIN_TO}

    CLUSTER_JOIN_TO_NAME=`yq -r .${CLUSTER_JOIN_TO}.serverName ${SANDBOX_CONFIG}`
else
    API_ENDPOINT=`oc get infrastructure cluster -o jsonpath='{.status.apiServerURL}' ${OC_ADDITIONAL_PARAMS}`
    JOINING_CLUSTER_NAME=`oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}' ${OC_ADDITIONAL_PARAMS}`

    login_to_cluster ${CLUSTER_JOIN_TO}

    CLUSTER_JOIN_TO_NAME=`oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}' ${OC_ADDITIONAL_PARAMS}`
fi

SECRET_NAME=${SA_NAME}-${JOINING_CLUSTER_NAME}
if [[ -n `oc get secret -n ${CLUSTER_JOIN_TO_OPERATOR_NS} ${OC_ADDITIONAL_PARAMS} | grep ${SECRET_NAME}` ]]; then
    oc delete secret ${SECRET_NAME} -n ${CLUSTER_JOIN_TO_OPERATOR_NS} ${OC_ADDITIONAL_PARAMS}
fi
oc create secret generic ${SECRET_NAME} --from-literal=token="${SA_TOKEN}" --from-literal=ca.crt="${SA_CA_CRT}" -n ${CLUSTER_JOIN_TO_OPERATOR_NS} ${OC_ADDITIONAL_PARAMS}

if [[ ${MULTI_NS} != "true" ]]; then
    TOOLCHAINCLUSTER_NAME=${JOINING_CLUSTER_TYPE_NAME}-${JOINING_CLUSTER_NAME}
    OWNER_CLUSTER_NAME=${CLUSTER_JOIN_TO}-${CLUSTER_JOIN_TO_NAME}
else
    TOOLCHAINCLUSTER_NAME=${JOINING_CLUSTER_TYPE_NAME}-${JOINING_CLUSTER_NAME}-${OPERATOR_NS}
    OWNER_CLUSTER_NAME=${CLUSTER_JOIN_TO}-${CLUSTER_JOIN_TO_NAME}-${CLUSTER_JOIN_TO_OPERATOR_NS}
fi

TOOLCHAINCLUSTER_CRD="apiVersion: toolchain.dev.openshift.com/v1alpha1
kind: ToolchainCluster
metadata:
  name: ${TOOLCHAINCLUSTER_NAME}
  namespace: ${CLUSTER_JOIN_TO_OPERATOR_NS}
  labels:
    type: ${JOINING_CLUSTER_TYPE_NAME}
    namespace: ${OPERATOR_NS}
    ownerClusterName: ${OWNER_CLUSTER_NAME}
spec:
  apiEndpoint: ${API_ENDPOINT}
  caBundle: ${SA_CA_CRT}
  secretRef:
    name: ${SECRET_NAME}
"

echo "Creating ToolchainCluster representation of ${JOINING_CLUSTER_TYPE} in ${CLUSTER_JOIN_TO}:"
echo ${TOOLCHAINCLUSTER_CRD}

cat <<EOF | oc apply ${OC_ADDITIONAL_PARAMS} -f -
${TOOLCHAINCLUSTER_CRD}
EOF
