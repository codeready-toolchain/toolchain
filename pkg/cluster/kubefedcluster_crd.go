package cluster

const kubeFedClusterCrd = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  name: kubefedclusters.core.kubefed.io
spec:
  additionalPrinterColumns:
  - JSONPath: .metadata.creationTimestamp
    name: age
    type: date
  - JSONPath: .status.conditions[?(@.type=='Ready')].status
    name: ready
    type: string
  group: core.kubefed.io
  names:
    kind: KubeFedCluster
    listKind: KubeFedClusterList
    plural: kubefedclusters
    singular: kubefedcluster
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: KubeFedCluster configures KubeFed to be aware of a Kubernetes cluster
        and encapsulates the details necessary to communicate with the cluster.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: KubeFedClusterSpec defines the desired state of KubeFedCluster
          properties:
            apiEndpoint:
              description: The API endpoint of the member cluster. This can be a hostname,
                hostname:port, IP or IP:port.
              type: string
            caBundle:
              description: CABundle contains the certificate authority information.
              format: byte
              type: string
            disabledTLSValidations:
              description: DisabledTLSValidations defines a list of checks to ignore
                when validating the TLS connection to the member cluster.  This can
                be any of *, SubjectName, or ValidityPeriod. If * is specified, it
                is expected to be the only option in list.
              items:
                type: string
              type: array
            secretRef:
              description: Name of the secret containing the token required to access
                the member cluster. The secret needs to exist in the same namespace
                as the control plane and should have a "token" key.
              properties:
                name:
                  description: Name of a secret within the enclosing namespace
                  type: string
              required:
              - name
              type: object
          required:
          - apiEndpoint
          - secretRef
          type: object
        status:
          description: KubeFedClusterStatus contains information about the current
            status of a cluster updated periodically by cluster controller.
          properties:
            conditions:
              description: Conditions is an array of current cluster conditions.
              items:
                description: ClusterCondition describes current state of a cluster.
                properties:
                  lastProbeTime:
                    description: Last time the condition was checked.
                    format: date-time
                    type: string
                  lastTransitionTime:
                    description: Last time the condition transit from one status to
                      another.
                    format: date-time
                    type: string
                  message:
                    description: Human readable message indicating details about last
                      transition.
                    type: string
                  reason:
                    description: (brief) reason for the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of True, False, Unknown.
                    type: string
                  type:
                    description: Type of cluster condition, Ready or Offline.
                    type: string
                required:
                - lastProbeTime
                - status
                - type
                type: object
              type: array
            region:
              description: Region is the name of the region in which all of the nodes
                in the cluster exist.  e.g. 'us-east1'.
              type: string
            zones:
              description: Zones are the names of availability zones in which the
                nodes of the cluster exist, e.g. 'us-east1-a'.
              items:
                type: string
              type: array
          required:
          - conditions
          type: object
      required:
      - spec
  version: v1beta1
  versions:
  - name: v1beta1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`
