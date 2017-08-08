{
  v1:: {
    local apiVersion = {apiVersion: "archon.kubeup.com/v1"},
    user:: {
      local kind = {kind: "User"},
      // new(name, userSpec, userLabels={}):: apiVersion + kind + self.mixin.metadata.name(name) + self.mixin.spec.template.spec.containers(containers) + self.mixin.spec.template.metadata.labels(podLabels),
      new(name): apiVersion + kind + self.mixin.metadata.name(name),
      mixin:: {
        metadata:: {
          local __metadataMixin(metadata) = {metadata+: metadata},
          mixinInstance(metadata):: __metadataMixin(metadata),
          // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
          annotations(annotations):: __metadataMixin({annotations+: annotations}),
          // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
          clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
          // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
          deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
          // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
          finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
          // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
          //
          // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
          //
          // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
          generateName(generateName):: __metadataMixin({generateName: generateName}),
          // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
          //
          // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
          initializers:: {
            local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
            mixinInstance(initializers):: __initializersMixin(initializers),
            // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
            pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
            pendingType:: hidden.meta.v1.initializer,
            // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
            result:: {
              local __resultMixin(result) = __initializersMixin({result+: result}),
              mixinInstance(result):: __resultMixin(result),
              // Suggested HTTP return code for this status, 0 if not set.
              code(code):: __resultMixin({code: code}),
              // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
              details:: {
                local __detailsMixin(details) = __resultMixin({details+: details}),
                mixinInstance(details):: __detailsMixin(details),
                // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                causesType:: hidden.meta.v1.statusCause,
                // The group attribute of the resource associated with the status StatusReason.
                group(group):: __detailsMixin({group: group}),
                // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                name(name):: __detailsMixin({name: name}),
                // If specified, the time in seconds before the operation should be retried.
                retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                uid(uid):: __detailsMixin({uid: uid}),
              },
              detailsType:: hidden.meta.v1.statusDetails,
              // A human-readable description of the status of this operation.
              message(message):: __resultMixin({message: message}),
              // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
              reason(reason):: __resultMixin({reason: reason}),
            },
            resultType:: hidden.meta.v1.status,
          },
          initializersType:: hidden.meta.v1.initializers,
          // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
          labels(labels):: __metadataMixin({labels+: labels}),
          // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: __metadataMixin({name: name}),
          // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
          //
          // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
          namespace(namespace):: __metadataMixin({namespace: namespace}),
        },
        metadataType:: hidden.meta.v1.objectMeta,
        spec:: {
          local __specMixin(spec) = {spec+: spec},
          mixinInstance(spec):: __specMixin(spec),
          name(name):: __specMixin({name: name}),
          passwordHash(passwordHash):: __specMixin({passwordHash: passwordHash}),
          sshAuthorizedKeys(sshAuthorizedKeys):: if std.type(sshAuthorizedKeys) == "array" then __specMixin({sshAuthorizedKeys+: sshAuthorizedKeys}) else __specMixin({sshAuthorizedKeys: [sshAuthorizedKeys]}),
          sudo(sudo):: __specMixin({sudo: sudo}),
          shell(shell):: __specMixin({shell: shell}),
        },
        specType:: hidden.archon.v1.userSpec,
      }
    },
    network:: {
      local kind = {kind: "Network"},
      new(name): apiVersion + kind + self.mixin.metadata.name(name),
      mixin:: {
        metadata:: {
          local __metadataMixin(metadata) = {metadata+: metadata},
          mixinInstance(metadata):: __metadataMixin(metadata),
          // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
          annotations(annotations):: __metadataMixin({annotations+: annotations}),
          // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
          clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
          // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
          deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
          // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
          finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
          // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
          //
          // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
          //
          // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
          generateName(generateName):: __metadataMixin({generateName: generateName}),
          // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
          //
          // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
          initializers:: {
            local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
            mixinInstance(initializers):: __initializersMixin(initializers),
            // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
            pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
            pendingType:: hidden.meta.v1.initializer,
            // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
            result:: {
              local __resultMixin(result) = __initializersMixin({result+: result}),
              mixinInstance(result):: __resultMixin(result),
              // Suggested HTTP return code for this status, 0 if not set.
              code(code):: __resultMixin({code: code}),
              // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
              details:: {
                local __detailsMixin(details) = __resultMixin({details+: details}),
                mixinInstance(details):: __detailsMixin(details),
                // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                causesType:: hidden.meta.v1.statusCause,
                // The group attribute of the resource associated with the status StatusReason.
                group(group):: __detailsMixin({group: group}),
                // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                name(name):: __detailsMixin({name: name}),
                // If specified, the time in seconds before the operation should be retried.
                retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                uid(uid):: __detailsMixin({uid: uid}),
              },
              detailsType:: hidden.meta.v1.statusDetails,
              // A human-readable description of the status of this operation.
              message(message):: __resultMixin({message: message}),
              // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
              reason(reason):: __resultMixin({reason: reason}),
            },
            resultType:: hidden.meta.v1.status,
          },
          initializersType:: hidden.meta.v1.initializers,
          // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
          labels(labels):: __metadataMixin({labels+: labels}),
          // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: __metadataMixin({name: name}),
          // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
          //
          // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
          namespace(namespace):: __metadataMixin({namespace: namespace}),
        },
        metadataType:: hidden.meta.v1.objectMeta,
        spec:: {
          local __specMixin(spec) = {spec+: spec},
          mixinInstance(spec):: __specMixin(spec),
          region(region):: __specMixin({region: region}),
          zone(zone):: __specMixin({zone: zone}),
          subnet(subnet):: __specMixin({subnet: subnet}),
        },
        specType:: hidden.archon.v1.networkSpec,
        status:: {
          local __statusMixin(status) = {status+: status},
          mixinInstance(status):: __statusMixin(status),
          phase(phase):: __statusMixin({phase: phase}),
        },
        statusType:: hidden.archon.v1.networkStatus,
      }
    },
    instanceGroup:: {
      local kind = {kind: "InstanceGroup"},
      new(name): apiVersion + kind + self.mixin.metadata.name(name),
      mixin:: {
        metadata:: {
          local __metadataMixin(metadata) = {metadata+: metadata},
          mixinInstance(metadata):: __metadataMixin(metadata),
          // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
          annotations(annotations):: __metadataMixin({annotations+: annotations}),
          // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
          clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
          // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
          deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
          // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
          finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
          // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
          //
          // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
          //
          // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
          generateName(generateName):: __metadataMixin({generateName: generateName}),
          // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
          //
          // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
          initializers:: {
            local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
            mixinInstance(initializers):: __initializersMixin(initializers),
            // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
            pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
            pendingType:: hidden.meta.v1.initializer,
            // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
            result:: {
              local __resultMixin(result) = __initializersMixin({result+: result}),
              mixinInstance(result):: __resultMixin(result),
              // Suggested HTTP return code for this status, 0 if not set.
              code(code):: __resultMixin({code: code}),
              // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
              details:: {
                local __detailsMixin(details) = __resultMixin({details+: details}),
                mixinInstance(details):: __detailsMixin(details),
                // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                causesType:: hidden.meta.v1.statusCause,
                // The group attribute of the resource associated with the status StatusReason.
                group(group):: __detailsMixin({group: group}),
                // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                name(name):: __detailsMixin({name: name}),
                // If specified, the time in seconds before the operation should be retried.
                retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                uid(uid):: __detailsMixin({uid: uid}),
              },
              detailsType:: hidden.meta.v1.statusDetails,
              // A human-readable description of the status of this operation.
              message(message):: __resultMixin({message: message}),
              // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
              reason(reason):: __resultMixin({reason: reason}),
            },
            resultType:: hidden.meta.v1.status,
          },
          initializersType:: hidden.meta.v1.initializers,
          // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
          labels(labels):: __metadataMixin({labels+: labels}),
          // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: __metadataMixin({name: name}),
          // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
          //
          // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
          namespace(namespace):: __metadataMixin({namespace: namespace}),
        },
        metadataType:: hidden.meta.v1.objectMeta,
        spec:: {
          local __specMixin(spec) = {spec+: spec},
          mixinInstance(spec):: __specMixin(spec),
          replicas(replicas):: __specMixin({replicas: replicas}),
          provisionPolicy(provisionPolicy):: __specMixin({provisionPolicy: provisionPolicy}),
          minReadySeconds(minReadySeconds):: __specMixin({minReadySeconds: minReadySeconds}),
          selector:: {
            local __selectorMixin(selector) = __specMixin({selector+: selector}),
            mixinInstance(selector):: __selectorMixin(selector),
            // matchExpressions is a list of label selector requirements. The requirements are ANDed.
            matchExpressions(matchExpressions):: if std.type(matchExpressions) == "array" then __selectorMixin({matchExpressions+: matchExpressions}) else __selectorMixin({matchExpressions+: [matchExpressions]}),
            matchExpressionsType:: hidden.meta.v1.labelSelectorRequirement,
            // matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
            matchLabels(matchLabels):: __selectorMixin({matchLabels+: matchLabels}),
          },
          selectorType:: hidden.meta.v1.labelSelector,
          reversedInstanceSelector:: {
            local __reversedInstanceSelector(reversedInstanceSelector) = __specMixin({reversedInstanceSelector+: reversedInstanceSelector}),
            mixinInstance(reversedInstanceSelector):: __reversedInstanceSelector(reversedInstanceSelector),
            // matchExpressions is a list of label selector requirements. The requirements are ANDed.
            matchExpressions(matchExpressions):: if std.type(matchExpressions) == "array" then __reversedInstanceSelector({matchExpressions+: matchExpressions}) else __reversedInstanceSelector({matchExpressions+: [matchExpressions]}),
            matchExpressionsType:: hidden.meta.v1.labelSelectorRequirement,
            // matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
            matchLabels(matchLabels):: __reversedInstanceSelector({matchLabels+: matchLabels}),
          },
          reversedInstanceSelectorType:: hidden.meta.v1.labelSelector,
          template:: {
            local __templateMixin(template) = __specMixin({template+: template}),
            mixinInstance(template):: __templateMixin(template),
            // Standard object's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
            metadata:: {
              local __metadataMixin(metadata) = __templateMixin({metadata+: metadata}),
              mixinInstance(metadata):: __metadataMixin(metadata),
              // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
              annotations(annotations):: __metadataMixin({annotations+: annotations}),
              // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
              clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
              // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
              deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
              // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
              finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
              // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
              //
              // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
              //
              // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
              generateName(generateName):: __metadataMixin({generateName: generateName}),
              // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
              //
              // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
              initializers:: {
                local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
                mixinInstance(initializers):: __initializersMixin(initializers),
                // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
                pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
                pendingType:: hidden.meta.v1.initializer,
                // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
                result:: {
                  local __resultMixin(result) = __initializersMixin({result+: result}),
                  mixinInstance(result):: __resultMixin(result),
                  // Suggested HTTP return code for this status, 0 if not set.
                  code(code):: __resultMixin({code: code}),
                  // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
                  details:: {
                    local __detailsMixin(details) = __resultMixin({details+: details}),
                    mixinInstance(details):: __detailsMixin(details),
                    // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                    causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                    causesType:: hidden.meta.v1.statusCause,
                    // The group attribute of the resource associated with the status StatusReason.
                    group(group):: __detailsMixin({group: group}),
                    // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                    name(name):: __detailsMixin({name: name}),
                    // If specified, the time in seconds before the operation should be retried.
                    retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                    // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                    uid(uid):: __detailsMixin({uid: uid}),
                  },
                  detailsType:: hidden.meta.v1.statusDetails,
                  // A human-readable description of the status of this operation.
                  message(message):: __resultMixin({message: message}),
                  // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
                  reason(reason):: __resultMixin({reason: reason}),
                },
                resultType:: hidden.meta.v1.status,
              },
              initializersType:: hidden.meta.v1.initializers,
              // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
              labels(labels):: __metadataMixin({labels+: labels}),
              // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
              name(name):: __metadataMixin({name: name}),
              // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
              //
              // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
              namespace(namespace):: __metadataMixin({namespace: namespace}),
            },
            metadataType:: hidden.meta.v1.objectMeta,
            // Specification of the desired behavior of the pod. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
            spec:: {
              local __specMixin(spec) = __templateMixin({spec+: spec}),
              mixinInstance(spec):: __specMixin(spec),
              os(os):: __specMixin({os: os}),
              image(image):: __specMixin({image: image}),
              instanceType(instanceType):: __specMixin({instanceType: instanceType}),
              networkName(networkName):: __specMixin({networkName: networkName}),
              reclaimPolicy(reclaimPolicy):: __specMixin({reclaimPolicy: reclaimPolicy}),
              files(files):: if std.type(files) == "array" then __specMixin({files+: files}) else __specMixin({files+: [files]}),
              filesType:: hidden.archon.v1.fileSpec,
              secrets(secrets):: if std.type(secrets) == "array" then __specMixin({secrets+: secrets}) else __specMixin({secrets+: [secrets]}),
              secretsType:: hidden.core.v1.localObjectReference,
              configs(configs):: if std.type(configs) == "array" then __specMixin({configs+: configs}) else __specMixin({configs+: [configs]}),
              configsType:: hidden.archon.v1.configSpec,
              users(users):: if std.type(users) == "array" then __specMixin({users+: users}) else __specMixin({users+: [users]}),
              usersType:: hidden.core.v1.localObjectReference,
              hostname(hostname):: __specMixin({hostname: hostname}),
              reservedInstanceRef:: {
                local __reservedInstanceRef(reservedInstanceRef) = __specMixin({reservedInstanceRef+: reservedInstanceRef}),
                mixinInstance(reservedInstanceRef):: __reservedInstanceRef(reservedInstanceRef),
                name(name):: __reservedInstanceRef({name: name}),
              },
              reservedInstanceRefType:: hidden.core.v1.localObjectReference,
            },
            specType:: hidden.archon.v1.instanceSpec,
            secrets(secrets):: if std.type(secrets) == "array" then __templateMixin({secrets+: secrets}) else __templateMixin({secrets+: [secrets]}),
            secretsType:: hidden.core.v1.secret,
          },
          templateType:: hidden.archon.v1.instanceTemplateSpec,
        },
        specType:: hidden.archon.v1.instanceGroupSpec,
      },
    },
    instance:: {
      local kind = {kind: "Instance"},
      new(name): apiVersion + kind + self.mixin.metadata.name(name),
      mixin:: {
        metadata:: {
          local __metadataMixin(metadata) = {metadata+: metadata},
          mixinInstance(metadata):: __metadataMixin(metadata),
          // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
          annotations(annotations):: __metadataMixin({annotations+: annotations}),
          // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
          clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
          // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
          deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
          // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
          finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
          // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
          //
          // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
          //
          // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
          generateName(generateName):: __metadataMixin({generateName: generateName}),
          // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
          //
          // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
          initializers:: {
            local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
            mixinInstance(initializers):: __initializersMixin(initializers),
            // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
            pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
            pendingType:: hidden.meta.v1.initializer,
            // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
            result:: {
              local __resultMixin(result) = __initializersMixin({result+: result}),
              mixinInstance(result):: __resultMixin(result),
              // Suggested HTTP return code for this status, 0 if not set.
              code(code):: __resultMixin({code: code}),
              // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
              details:: {
                local __detailsMixin(details) = __resultMixin({details+: details}),
                mixinInstance(details):: __detailsMixin(details),
                // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                causesType:: hidden.meta.v1.statusCause,
                // The group attribute of the resource associated with the status StatusReason.
                group(group):: __detailsMixin({group: group}),
                // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                name(name):: __detailsMixin({name: name}),
                // If specified, the time in seconds before the operation should be retried.
                retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                uid(uid):: __detailsMixin({uid: uid}),
              },
              detailsType:: hidden.meta.v1.statusDetails,
              // A human-readable description of the status of this operation.
              message(message):: __resultMixin({message: message}),
              // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
              reason(reason):: __resultMixin({reason: reason}),
            },
            resultType:: hidden.meta.v1.status,
          },
          initializersType:: hidden.meta.v1.initializers,
          // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
          labels(labels):: __metadataMixin({labels+: labels}),
          // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: __metadataMixin({name: name}),
          // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
          //
          // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
          namespace(namespace):: __metadataMixin({namespace: namespace}),
        },
        metadataType:: hidden.meta.v1.objectMeta,
        spec:: {
          local __specMixin(spec) = {spec+: spec},
          mixinInstance(spec):: __specMixin(spec),
          os(os):: __specMixin({os: os}),
          image(image):: __specMixin({image: image}),
          instanceType(instanceType):: __specMixin({instanceType: instanceType}),
          networkName(networkName):: __specMixin({networkName: networkName}),
          reclaimPolicy(reclaimPolicy):: __specMixin({reclaimPolicy: reclaimPolicy}),
          files(files):: if std.type(files) == "array" then __specMixin({files+: files}) else __specMixin({files+: [files]}),
          filesType:: hidden.archon.v1.fileSpec,
          secrets(secrets):: if std.type(secrets) == "array" then __specMixin({secrets+: secrets}) else __specMixin({secrets+: [secrets]}),
          secretsType:: hidden.core.v1.localObjectReference,
          configs(configs):: if std.type(configs) == "array" then __specMixin({configs+: configs}) else __specMixin({configs+: [configs]}),
          configsType:: hidden.archon.v1.configSpec,
          users(users):: if std.type(users) == "array" then __specMixin({users+: users}) else __specMixin({users+: [users]}),
          usersType:: hidden.core.v1.localObjectReference,
          hostname(hostname):: __specMixin({hostname: hostname}),
          reservedInstanceRef:: {
            local __reservedInstanceRef(reservedInstanceRef) = __specMixin({reservedInstanceRef+: reservedInstanceRef}),
            mixinInstance(reservedInstanceRef):: __reservedInstanceRef(reservedInstanceRef),
            name(name):: __reservedInstanceRef({name: name}),
          },
          reservedInstanceRefType:: hidden.core.v1.localObjectReference,
        },
        specType:: hidden.archon.v1.instanceSpec,
        status:: {
          local __statusMixin(status) = {status+: status},
          mixinInstance(status):: __statusMixin(status),
          phase(phase):: __statusMixin({phase: phase}),
          conditions(conditions):: if std.type(conditions) == "array" then __statusMixin({conditions+: conditions}) else __statusMixin({conditions+: [conditions]}),
          conditionsType:: hidden.archon.v1.instanceCondition,
          privateIP(privateIP):: __statusMixin({privateIP: privateIP}),
          publicIP(publicIP):: __statusMixin({publicIP: publicIP}),
          instanceID(instanceID):: __statusMixin({instanceID: instanceID}),
          creationTimestamp(creationTimestamp):: __statusMixin({creationTimestamp: creationTimestamp}),
        },
        statusType:: hidden.archon.v1.instanceStatus,
      }
    },
  },
  local hidden = {
    meta:: {
      v1:: {
        local apiVersion = {apiVersion: "meta/v1"},
        // APIGroup contains the name, the supported versions, and the preferred version of a group.
        apiGroup:: {
          new():: {},
          // name is the name of the group.
          name(name):: {name: name},
          // a map of client CIDR to server address that is serving this group. This is to help clients reach servers in the most network-efficient way possible. Clients can use the appropriate server address as per the CIDR that they match. In case of multiple matches, clients should use the longest matching CIDR. The server returns only those CIDRs that it thinks that the client can match. For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP. Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.
          serverAddressByClientCidrs(serverAddressByClientCidrs):: if std.type(serverAddressByClientCidrs) == "array" then {serverAddressByClientCIDRs+: serverAddressByClientCidrs} else {serverAddressByClientCIDRs+: [serverAddressByClientCidrs]},
          serverAddressByClientCidrsType:: hidden.meta.v1.serverAddressByClientCidr,
          // versions are the versions supported in this group.
          versions(versions):: if std.type(versions) == "array" then {versions+: versions} else {versions+: [versions]},
          versionsType:: hidden.meta.v1.groupVersionForDiscovery,
          mixin:: {
            // preferredVersion is the version preferred by the API server, which probably is the storage version.
            preferredVersion:: {
              local __preferredVersionMixin(preferredVersion) = {preferredVersion+: preferredVersion},
              mixinInstance(preferredVersion):: __preferredVersionMixin(preferredVersion),
              // groupVersion specifies the API group and version in the form "group/version"
              groupVersion(groupVersion):: __preferredVersionMixin({groupVersion: groupVersion}),
              // version specifies the version in the form of "version". This is to save the clients the trouble of splitting the GroupVersion.
              version(version):: __preferredVersionMixin({version: version}),
            },
            preferredVersionType:: hidden.meta.v1.groupVersionForDiscovery,
          },
        },
        // APIGroupList is a list of APIGroup, to allow clients to discover the API at /apis.
        apiGroupList:: {
          new():: {},
          // groups is a list of APIGroup.
          groups(groups):: if std.type(groups) == "array" then {groups+: groups} else {groups+: [groups]},
          groupsType:: hidden.meta.v1.apiGroup,
          mixin:: {
          },
        },
        // APIResource specifies the name of a resource and whether it is namespaced.
        apiResource:: {
          new():: {},
          // name is the plural name of the resource.
          name(name):: {name: name},
          // namespaced indicates if a resource is namespaced or not.
          namespaced(namespaced):: {namespaced: namespaced},
          // shortNames is a list of suggested short names of the resource.
          shortNames(shortNames):: if std.type(shortNames) == "array" then {shortNames+: shortNames} else {shortNames+: [shortNames]},
          // singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely. The singularName is more correct for reporting status on a single item and both singular and plural are allowed from the kubectl CLI interface.
          singularName(singularName):: {singularName: singularName},
          // verbs is a list of supported kube verbs (this includes get, list, watch, create, update, patch, delete, deletecollection, and proxy)
          verbs(verbs):: if std.type(verbs) == "array" then {verbs+: verbs} else {verbs+: [verbs]},
          mixin:: {
          },
        },
        // APIResourceList is a list of APIResource, it is used to expose the name of the resources supported in a specific group and version, and if the resource is namespaced.
        apiResourceList:: {
          new():: {},
          // groupVersion is the group and version this APIResourceList is for.
          groupVersion(groupVersion):: {groupVersion: groupVersion},
          // resources contains the name of the resources and if they are namespaced.
          resources(resources):: if std.type(resources) == "array" then {resources+: resources} else {resources+: [resources]},
          resourcesType:: hidden.meta.v1.apiResource,
          mixin:: {
          },
        },
        // APIVersions lists the versions that are available, to allow clients to discover the API at /api, which is the root path of the legacy v1 API.
        apiVersions:: {
          new():: {},
          // a map of client CIDR to server address that is serving this group. This is to help clients reach servers in the most network-efficient way possible. Clients can use the appropriate server address as per the CIDR that they match. In case of multiple matches, clients should use the longest matching CIDR. The server returns only those CIDRs that it thinks that the client can match. For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP. Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.
          serverAddressByClientCidrs(serverAddressByClientCidrs):: if std.type(serverAddressByClientCidrs) == "array" then {serverAddressByClientCIDRs+: serverAddressByClientCidrs} else {serverAddressByClientCIDRs+: [serverAddressByClientCidrs]},
          serverAddressByClientCidrsType:: hidden.meta.v1.serverAddressByClientCidr,
          // versions are the api versions that are available.
          versions(versions):: if std.type(versions) == "array" then {versions+: versions} else {versions+: [versions]},
          mixin:: {
          },
        },
        // DeleteOptions may be provided when deleting an API object.
        deleteOptions:: {
          new():: {},
          // The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately.
          gracePeriodSeconds(gracePeriodSeconds):: {gracePeriodSeconds: gracePeriodSeconds},
          // Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both.
          orphanDependents(orphanDependents):: {orphanDependents: orphanDependents},
          // Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy.
          propagationPolicy(propagationPolicy):: {propagationPolicy: propagationPolicy},
          mixin:: {
            // Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status will be returned.
            preconditions:: {
              local __preconditionsMixin(preconditions) = {preconditions+: preconditions},
              mixinInstance(preconditions):: __preconditionsMixin(preconditions),
              // Specifies the target UID.
              uid(uid):: __preconditionsMixin({uid: uid}),
            },
            preconditionsType:: hidden.meta.v1.preconditions,
          },
        },
        // GroupVersion contains the "group/version" and "version" string of a version. It is made a struct to keep extensibility.
        groupVersionForDiscovery:: {
          new():: {},
          // groupVersion specifies the API group and version in the form "group/version"
          groupVersion(groupVersion):: {groupVersion: groupVersion},
          // version specifies the version in the form of "version". This is to save the clients the trouble of splitting the GroupVersion.
          version(version):: {version: version},
          mixin:: {
          },
        },
        // Initializer is information about an initializer that has not yet completed.
        initializer:: {
          new():: {},
          // name of the process that is responsible for initializing this object.
          name(name):: {name: name},
          mixin:: {
          },
        },
        // Initializers tracks the progress of initialization.
        initializers:: {
          new():: {},
          // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
          pending(pending):: if std.type(pending) == "array" then {pending+: pending} else {pending+: [pending]},
          pendingType:: hidden.meta.v1.initializer,
          mixin:: {
            // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
            result:: {
              local __resultMixin(result) = {result+: result},
              mixinInstance(result):: __resultMixin(result),
              // Suggested HTTP return code for this status, 0 if not set.
              code(code):: __resultMixin({code: code}),
              // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
              details:: {
                local __detailsMixin(details) = __resultMixin({details+: details}),
                mixinInstance(details):: __detailsMixin(details),
                // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                causesType:: hidden.meta.v1.statusCause,
                // The group attribute of the resource associated with the status StatusReason.
                group(group):: __detailsMixin({group: group}),
                // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                name(name):: __detailsMixin({name: name}),
                // If specified, the time in seconds before the operation should be retried.
                retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                uid(uid):: __detailsMixin({uid: uid}),
              },
              detailsType:: hidden.meta.v1.statusDetails,
              // A human-readable description of the status of this operation.
              message(message):: __resultMixin({message: message}),
              // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
              reason(reason):: __resultMixin({reason: reason}),
            },
            resultType:: hidden.meta.v1.status,
          },
        },
        // A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
        labelSelector:: {
          new():: {},
          // matchExpressions is a list of label selector requirements. The requirements are ANDed.
          matchExpressions(matchExpressions):: if std.type(matchExpressions) == "array" then {matchExpressions+: matchExpressions} else {matchExpressions+: [matchExpressions]},
          matchExpressionsType:: hidden.meta.v1.labelSelectorRequirement,
          // matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
          matchLabels(matchLabels):: {matchLabels+: matchLabels},
          mixin:: {
          },
        },
        // A label selector requirement is a selector that contains values, a key, and an operator that relates the key and values.
        labelSelectorRequirement:: {
          new():: {},
          // key is the label key that the selector applies to.
          key(key):: {key: key},
          // operator represents a key's relationship to a set of values. Valid operators ard In, NotIn, Exists and DoesNotExist.
          operator(operator):: {operator: operator},
          // values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
          values(values):: if std.type(values) == "array" then {values+: values} else {values+: [values]},
          mixin:: {
          },
        },
        // ListMeta describes metadata that synthetic resources must have, including lists and various status objects. A resource may have only one of {ObjectMeta, ListMeta}.
        listMeta:: {
          new():: {},
          // String that identifies the server's internal version of this object that can be used by clients to determine when objects have changed. Value must be treated as opaque by clients and passed unmodified back to the server. Populated by the system. Read-only. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#concurrency-control-and-consistency
          resourceVersion(resourceVersion):: {resourceVersion: resourceVersion},
          // SelfLink is a URL representing this object. Populated by the system. Read-only.
          selfLink(selfLink):: {selfLink: selfLink},
          mixin:: {
          },
        },
        // ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.
        objectMeta:: {
          new():: {},
          // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
          annotations(annotations):: {annotations+: annotations},
          // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
          clusterName(clusterName):: {clusterName: clusterName},
          // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
          deletionGracePeriodSeconds(deletionGracePeriodSeconds):: {deletionGracePeriodSeconds: deletionGracePeriodSeconds},
          // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
          finalizers(finalizers):: if std.type(finalizers) == "array" then {finalizers+: finalizers} else {finalizers+: [finalizers]},
          // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
          //
          // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
          //
          // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
          generateName(generateName):: {generateName: generateName},
          // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
          labels(labels):: {labels+: labels},
          // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: {name: name},
          // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
          //
          // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
          namespace(namespace):: {namespace: namespace},
          mixin:: {
            // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
            //
            // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
            initializers:: {
              local __initializersMixin(initializers) = {initializers+: initializers},
              mixinInstance(initializers):: __initializersMixin(initializers),
              // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
              pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
              pendingType:: hidden.meta.v1.initializer,
              // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
              result:: {
                local __resultMixin(result) = __initializersMixin({result+: result}),
                mixinInstance(result):: __resultMixin(result),
                // Suggested HTTP return code for this status, 0 if not set.
                code(code):: __resultMixin({code: code}),
                // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
                details:: {
                  local __detailsMixin(details) = __resultMixin({details+: details}),
                  mixinInstance(details):: __detailsMixin(details),
                  // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                  causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                  causesType:: hidden.meta.v1.statusCause,
                  // The group attribute of the resource associated with the status StatusReason.
                  group(group):: __detailsMixin({group: group}),
                  // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                  name(name):: __detailsMixin({name: name}),
                  // If specified, the time in seconds before the operation should be retried.
                  retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                  // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                  uid(uid):: __detailsMixin({uid: uid}),
                },
                detailsType:: hidden.meta.v1.statusDetails,
                // A human-readable description of the status of this operation.
                message(message):: __resultMixin({message: message}),
                // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
                reason(reason):: __resultMixin({reason: reason}),
              },
              resultType:: hidden.meta.v1.status,
            },
            initializersType:: hidden.meta.v1.initializers,
          },
        },
        // OwnerReference contains enough information to let you identify an owning object. Currently, an owning object must be in the same namespace, so there is no namespace field.
        ownerReference:: {
          new():: {},
          // If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
          blockOwnerDeletion(blockOwnerDeletion):: {blockOwnerDeletion: blockOwnerDeletion},
          // If true, this reference points to the managing controller.
          controller(controller):: {controller: controller},
          // Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names
          name(name):: {name: name},
          // UID of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#uids
          uid(uid):: {uid: uid},
          mixin:: {
          },
        },
        // Patch is provided to give a concrete name and type to the Kubernetes PATCH request body.
        patch:: {
          new():: {},
          mixin:: {
          },
        },
        // Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
        preconditions:: {
          new():: {},
          // Specifies the target UID.
          uid(uid):: {uid: uid},
          mixin:: {
          },
        },
        // ServerAddressByClientCIDR helps the client to determine the server address that they should use, depending on the clientCIDR that they match.
        serverAddressByClientCidr:: {
          new():: {},
          // The CIDR with which clients can match their IP to figure out the server address that they should use.
          clientCidr(clientCidr):: {clientCIDR: clientCidr},
          // Address of this server, suitable for a client that matches the above CIDR. This can be a hostname, hostname:port, IP or IP:port.
          serverAddress(serverAddress):: {serverAddress: serverAddress},
          mixin:: {
          },
        },
        // Status is a return value for calls that don't return other objects.
        status:: {
          new():: {},
          // Suggested HTTP return code for this status, 0 if not set.
          code(code):: {code: code},
          // A human-readable description of the status of this operation.
          message(message):: {message: message},
          // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
          reason(reason):: {reason: reason},
          mixin:: {
            // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
            details:: {
              local __detailsMixin(details) = {details+: details},
              mixinInstance(details):: __detailsMixin(details),
              // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
              causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
              causesType:: hidden.meta.v1.statusCause,
              // The group attribute of the resource associated with the status StatusReason.
              group(group):: __detailsMixin({group: group}),
              // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
              name(name):: __detailsMixin({name: name}),
              // If specified, the time in seconds before the operation should be retried.
              retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
              // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
              uid(uid):: __detailsMixin({uid: uid}),
            },
            detailsType:: hidden.meta.v1.statusDetails,
          },
        },
        // StatusCause provides more information about an api.Status failure, including cases when multiple errors are encountered.
        statusCause:: {
          new():: {},
          // The field of the resource that has caused this error, as named by its JSON serialization. May include dot and postfix notation for nested attributes. Arrays are zero-indexed.  Fields may appear more than once in an array of causes due to fields having multiple errors. Optional.
          //
          // Examples:
          //   "name" - the field "name" on the current resource
          //   "items[0].name" - the field "name" on the first array entry in "items"
          field(field):: {field: field},
          // A human-readable description of the cause of the error.  This field may be presented as-is to a reader.
          message(message):: {message: message},
          // A machine-readable description of the cause of the error. If this value is empty there is no information available.
          reason(reason):: {reason: reason},
          mixin:: {
          },
        },
        // StatusDetails is a set of additional properties that MAY be set by the server to provide additional information about a response. The Reason field of a Status object defines what attributes will be set. Clients must ignore fields that do not match the defined type of each attribute, and should assume that any attribute may be empty, invalid, or under defined.
        statusDetails:: {
          new():: {},
          // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
          causes(causes):: if std.type(causes) == "array" then {causes+: causes} else {causes+: [causes]},
          causesType:: hidden.meta.v1.statusCause,
          // The group attribute of the resource associated with the status StatusReason.
          group(group):: {group: group},
          // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
          name(name):: {name: name},
          // If specified, the time in seconds before the operation should be retried.
          retryAfterSeconds(retryAfterSeconds):: {retryAfterSeconds: retryAfterSeconds},
          // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
          uid(uid):: {uid: uid},
          mixin:: {
          },
        },
        //
        time:: {
          new():: {},
          mixin:: {
          },
        },
        // Event represents a single event to a watched resource.
        watchEvent:: {
          new():: {},
          //
          type(type):: {type: type},
          mixin:: {
          },
        },
      },
    },
    archon:: {
      v1:: {
        local apiVersion = {apiVersion: "archon.kubeup.com/v1"},
        userSpec:: {
          new():: {},
          name(name):: {name: name},
          passwordHash(passwordHash):: {passwordHash: passwordHash},
          sshAuthorizedKeys(sshAuthorizedKeys):: if std.type(sshAuthorizedKeys) == "array" then {sshAuthorizedKeys+: sshAuthorizedKeys} else {sshAuthorizedKeys: [sshAuthorizedKeys]},
          sudo(sudo):: {sudo: sudo},
          shell(shell):: {shell: shell},
          mixin:: {
          },
        },
        networkSpec:: {
          new():: {},
          region(region):: {region: region},
          zone(zone):: {zone: zone},
          subnet(subnet):: {subnet: subnet},
          mixin:: {
          },
        },
        networkStatus:: {
          new():: {},
          phase(phase):: {phase: phase},
          mixin:: {
          },
        },
        instanceGroupSpec:: {
          new():: {},
          replicas(replicas):: {replicas: replicas},
          provisionPolicy(provisionPolicy):: {provisionPolicy: provisionPolicy},
          minReadySeconds(minReadySeconds):: {minReadySeconds: minReadySeconds},
          mixin:: {
            selector:: {
              local __selectorMixin(selector) = {selector+: selector},
              mixinInstance(selector):: __selectorMixin(selector),
              // matchExpressions is a list of label selector requirements. The requirements are ANDed.
              matchExpressions(matchExpressions):: if std.type(matchExpressions) == "array" then __selectorMixin({matchExpressions+: matchExpressions}) else __selectorMixin({matchExpressions+: [matchExpressions]}),
              matchExpressionsType:: hidden.meta.v1.labelSelectorRequirement,
              // matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
              matchLabels(matchLabels):: __selectorMixin({matchLabels+: matchLabels}),
            },
            selectorType:: hidden.meta.v1.labelSelector,
            reversedInstanceSelector:: {
              local __reversedInstanceSelector(reversedInstanceSelector) = {reversedInstanceSelector+: reversedInstanceSelector},
              mixinInstance(reversedInstanceSelector):: __reversedInstanceSelector(reversedInstanceSelector),
              // matchExpressions is a list of label selector requirements. The requirements are ANDed.
              matchExpressions(matchExpressions):: if std.type(matchExpressions) == "array" then __reversedInstanceSelector({matchExpressions+: matchExpressions}) else __reversedInstanceSelector({matchExpressions+: [matchExpressions]}),
              matchExpressionsType:: hidden.meta.v1.labelSelectorRequirement,
              // matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
              matchLabels(matchLabels):: __reversedInstanceSelector({matchLabels+: matchLabels}),
            },
            reversedInstanceSelectorType:: hidden.meta.v1.labelSelector,
            template:: {
              local __templateMixin(template) = {template+: template},
              mixinInstance(template):: __templateMixin(template),
              // Standard object's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
              metadata:: {
                local __metadataMixin(metadata) = __templateMixin({metadata+: metadata}),
                mixinInstance(metadata):: __metadataMixin(metadata),
                // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
                annotations(annotations):: __metadataMixin({annotations+: annotations}),
                // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
                clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
                // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
                deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
                // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
                finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
                // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
                //
                // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
                //
                // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
                generateName(generateName):: __metadataMixin({generateName: generateName}),
                // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
                //
                // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
                initializers:: {
                  local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
                  mixinInstance(initializers):: __initializersMixin(initializers),
                  // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
                  pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
                  pendingType:: hidden.meta.v1.initializer,
                  // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
                  result:: {
                    local __resultMixin(result) = __initializersMixin({result+: result}),
                    mixinInstance(result):: __resultMixin(result),
                    // Suggested HTTP return code for this status, 0 if not set.
                    code(code):: __resultMixin({code: code}),
                    // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
                    details:: {
                      local __detailsMixin(details) = __resultMixin({details+: details}),
                      mixinInstance(details):: __detailsMixin(details),
                      // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                      causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                      causesType:: hidden.meta.v1.statusCause,
                      // The group attribute of the resource associated with the status StatusReason.
                      group(group):: __detailsMixin({group: group}),
                      // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                      name(name):: __detailsMixin({name: name}),
                      // If specified, the time in seconds before the operation should be retried.
                      retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                      // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                      uid(uid):: __detailsMixin({uid: uid}),
                    },
                    detailsType:: hidden.meta.v1.statusDetails,
                    // A human-readable description of the status of this operation.
                    message(message):: __resultMixin({message: message}),
                    // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
                    reason(reason):: __resultMixin({reason: reason}),
                  },
                  resultType:: hidden.meta.v1.status,
                },
                initializersType:: hidden.meta.v1.initializers,
                // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
                labels(labels):: __metadataMixin({labels+: labels}),
                // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
                name(name):: __metadataMixin({name: name}),
                // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
                //
                // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
                namespace(namespace):: __metadataMixin({namespace: namespace}),
              },
              metadataType:: hidden.meta.v1.objectMeta,
              // Specification of the desired behavior of the pod. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
              spec:: {
                local __specMixin(spec) = __templateMixin({spec+: spec}),
                mixinInstance(spec):: __specMixin(spec),
                os(os):: __specMixin({os: os}),
                image(image):: __specMixin({image: image}),
                instanceType(instanceType):: __specMixin({instanceType: instanceType}),
                networkName(networkName):: __specMixin({networkName: networkName}),
                reclaimPolicy(reclaimPolicy):: __specMixin({reclaimPolicy: reclaimPolicy}),
                files(files):: if std.type(files) == "array" then __specMixin({files+: files}) else __specMixin({files+: [files]}),
                filesType:: hidden.archon.v1.fileSpec,
                secrets(secrets):: if std.type(secrets) == "array" then __specMixin({secrets+: secrets}) else __specMixin({secrets+: [secrets]}),
                secretsType:: hidden.core.v1.localObjectReference,
                configs(configs):: if std.type(configs) == "array" then __specMixin({configs+: configs}) else __specMixin({configs+: [configs]}),
                configsType:: hidden.archon.v1.configSpec,
                users(users):: if std.type(users) == "array" then __specMixin({users+: users}) else __specMixin({users+: [users]}),
                usersType:: hidden.core.v1.localObjectReference,
                hostname(hostname):: __specMixin({hostname: hostname}),
                reservedInstanceRef:: {
                  local __reservedInstanceRef(reservedInstanceRef) = __specMixin({reservedInstanceRef+: reservedInstanceRef}),
                  mixinInstance(reservedInstanceRef):: __reservedInstanceRef(reservedInstanceRef),
                  name(name):: __reservedInstanceRef({name: name}),
                },
                reservedInstanceRefType:: hidden.core.v1.localObjectReference,
              },
              specType:: hidden.archon.v1.instanceSpec,
              secrets(secrets):: if std.type(secrets) == "array" then __templateMixin({secrets+: secrets}) else __templateMixin({secrets+: [secrets]}),
              secretsType:: hidden.core.v1.secret,
            },
            templateType:: hidden.archon.v1.instanceTemplateSpec,
          },
        },
        fileSpec:: {
          new():: {},
          name(name):: {name: name},
          encoding(encoding):: {encoding: encoding},
          content(content):: {content: content},
          template(template):: {template: template},
          owner(owner):: {owner: owner},
          userID(userID):: {userID: userID},
          groupID(groupID):: {groupID: groupID},
          filesystem(filesystem):: {filesystem: filesystem},
          path(path):: {path: path},
          permissions(permissions):: {permissions: permissions},
          mixin:: {
          },
        },
        configSpec:: {
          new():: {},
          name(name):: {name: name},
          data(data):: {data: data},
          mixin:: {
          },
        },
        instanceTemplateSpec:: {
          new():: {},
          secrets(secrets):: if std.type(secrets) == "array" then {secrets+: secrets} else {secrets+: [secrets]},
          secretsType:: hidden.core.v1.secret,
          mixin:: {
            metadata:: {
              local __metadataMixin(metadata) = {metadata+: metadata},
              mixinInstance(metadata):: __metadataMixin(metadata),
              // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
              annotations(annotations):: __metadataMixin({annotations+: annotations}),
              // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
              clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
              // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
              deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
              // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
              finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
              // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
              //
              // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
              //
              // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
              generateName(generateName):: __metadataMixin({generateName: generateName}),
              // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
              //
              // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
              initializers:: {
                local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
                mixinInstance(initializers):: __initializersMixin(initializers),
                // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
                pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
                pendingType:: hidden.meta.v1.initializer,
                // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
                result:: {
                  local __resultMixin(result) = __initializersMixin({result+: result}),
                  mixinInstance(result):: __resultMixin(result),
                  // Suggested HTTP return code for this status, 0 if not set.
                  code(code):: __resultMixin({code: code}),
                  // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
                  details:: {
                    local __detailsMixin(details) = __resultMixin({details+: details}),
                    mixinInstance(details):: __detailsMixin(details),
                    // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                    causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                    causesType:: hidden.meta.v1.statusCause,
                    // The group attribute of the resource associated with the status StatusReason.
                    group(group):: __detailsMixin({group: group}),
                    // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                    name(name):: __detailsMixin({name: name}),
                    // If specified, the time in seconds before the operation should be retried.
                    retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                    // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                    uid(uid):: __detailsMixin({uid: uid}),
                  },
                  detailsType:: hidden.meta.v1.statusDetails,
                  // A human-readable description of the status of this operation.
                  message(message):: __resultMixin({message: message}),
                  // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
                  reason(reason):: __resultMixin({reason: reason}),
                },
                resultType:: hidden.meta.v1.status,
              },
              initializersType:: hidden.meta.v1.initializers,
              // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
              labels(labels):: __metadataMixin({labels+: labels}),
              // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
              name(name):: __metadataMixin({name: name}),
              // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
              //
              // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
              namespace(namespace):: __metadataMixin({namespace: namespace}),
            },
            metadataType:: hidden.meta.v1.objectMeta,
            // Specification of the desired behavior of the pod. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
            spec:: {
              local __specMixin(spec) = {spec+: spec},
              mixinInstance(spec):: __specMixin(spec),
              os(os):: __specMixin({os: os}),
              image(image):: __specMixin({image: image}),
              instanceType(instanceType):: __specMixin({instanceType: instanceType}),
              networkName(networkName):: __specMixin({networkName: networkName}),
              reclaimPolicy(reclaimPolicy):: __specMixin({reclaimPolicy: reclaimPolicy}),
              files(files):: if std.type(files) == "array" then __specMixin({files+: files}) else __specMixin({files+: [files]}),
              filesType:: hidden.archon.v1.fileSpec,
              secrets(secrets):: if std.type(secrets) == "array" then __specMixin({secrets+: secrets}) else __specMixin({secrets+: [secrets]}),
              secretsType:: hidden.core.v1.localObjectReference,
              configs(configs):: if std.type(configs) == "array" then __specMixin({configs+: configs}) else __specMixin({configs+: [configs]}),
              configsType:: hidden.archon.v1.configSpec,
              users(users):: if std.type(users) == "array" then __specMixin({users+: users}) else __specMixin({users+: [users]}),
              usersType:: hidden.core.v1.localObjectReference,
              hostname(hostname):: __specMixin({hostname: hostname}),
              reservedInstanceRef:: {
                local __reservedInstanceRef(reservedInstanceRef) = __specMixin({reservedInstanceRef+: reservedInstanceRef}),
                mixinInstance(reservedInstanceRef):: __reservedInstanceRef(reservedInstanceRef),
                name(name):: __reservedInstanceRef({name: name}),
              },
              reservedInstanceRefType:: hidden.core.v1.localObjectReference,
            },
            specType:: hidden.archon.v1.instanceSpec,
            },
        },
        instanceSpec:: {
          new():: {},
          os(os):: {os: os},
          image(image):: {image: image},
          instanceType(instanceType):: {instanceType: instanceType},
          networkName(networkName):: {networkName: networkName},
          reclaimPolicy(reclaimPolicy):: {reclaimPolicy: reclaimPolicy},
          files(files):: if std.type(files) == "array" then {files+: files} else {files+: [files]},
          filesType:: hidden.archon.v1.fileSpec,
          secrets(secrets):: if std.type(secrets) == "array" then {secrets+: secrets} else {secrets+: [secrets]},
          secretsType:: hidden.core.v1.localObjectReference,
          configs(configs):: if std.type(configs) == "array" then {configs+: configs} else {configs+: [configs]},
          configsType:: hidden.archon.v1.configSpec,
          users(users):: if std.type(users) == "array" then {users+: users} else {users+: [users]},
          usersType:: hidden.core.v1.localObjectReference,
          hostname(hostname):: {hostname: hostname},
          mixin:: {
            reservedInstanceRef:: {
              local __reservedInstanceRef(reservedInstanceRef) = {reservedInstanceRef+: reservedInstanceRef},
              mixinInstance(reservedInstanceRef):: __reservedInstanceRef(reservedInstanceRef),
              name(name):: __reservedInstanceRef({name: name}),
            },
            reservedInstanceRefType:: hidden.core.v1.localObjectReference,
          },
        },
      },
    },
    core:: {
      v1:: {
        local apiVersion = {apiVersion: "v1"},
        localObjectReference:: {
          new():: {},
          // Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
          name(name):: {name: name},
          mixin:: {
          },
        },
        // Secret holds secret data of a certain type. The total bytes of the values in the Data field must be less than MaxSecretSize bytes.
        secret:: {
          local kind = {kind: "Secret"},
          new():: apiVersion + kind,
          // Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
          data(data):: {data+: data},
          // stringData allows specifying non-binary secret data in string form. It is provided as a write-only convenience method. All keys and values are merged into the data field on write, overwriting any existing values. It is never output when reading from the API.
          stringData(stringData):: {stringData+: stringData},
          // Used to facilitate programmatic handling of secret data.
          type(type):: {type: type},
          mixin:: {
            // Standard object's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
            metadata:: {
              local __metadataMixin(metadata) = {metadata+: metadata},
              mixinInstance(metadata):: __metadataMixin(metadata),
              // Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations
              annotations(annotations):: __metadataMixin({annotations+: annotations}),
              // The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
              clusterName(clusterName):: __metadataMixin({clusterName: clusterName}),
              // Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
              deletionGracePeriodSeconds(deletionGracePeriodSeconds):: __metadataMixin({deletionGracePeriodSeconds: deletionGracePeriodSeconds}),
              // Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed.
              finalizers(finalizers):: if std.type(finalizers) == "array" then __metadataMixin({finalizers+: finalizers}) else __metadataMixin({finalizers+: [finalizers]}),
              // GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
              //
              // If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).
              //
              // Applied only if Name is not specified. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency
              generateName(generateName):: __metadataMixin({generateName: generateName}),
              // An initializer is a controller which enforces some system invariant at object creation time. This field is a list of initializers that have not yet acted on this object. If nil or empty, this object has been completely initialized. Otherwise, the object is considered uninitialized and is hidden (in list/watch and get calls) from clients that haven't explicitly asked to observe uninitialized objects.
              //
              // When an object is created, the system will populate this list with the current set of initializers. Only privileged users may set or modify this list. Once it is empty, it may not be modified further by any user.
              initializers:: {
                local __initializersMixin(initializers) = __metadataMixin({initializers+: initializers}),
                mixinInstance(initializers):: __initializersMixin(initializers),
                // Pending is a list of initializers that must execute in order before this object is visible. When the last pending initializer is removed, and no failing result is set, the initializers struct will be set to nil and the object is considered as initialized and visible to all clients.
                pending(pending):: if std.type(pending) == "array" then __initializersMixin({pending+: pending}) else __initializersMixin({pending+: [pending]}),
                pendingType:: hidden.meta.v1.initializer,
                // If result is set with the Failure field, the object will be persisted to storage and then deleted, ensuring that other clients can observe the deletion.
                result:: {
                  local __resultMixin(result) = __initializersMixin({result+: result}),
                  mixinInstance(result):: __resultMixin(result),
                  // Suggested HTTP return code for this status, 0 if not set.
                  code(code):: __resultMixin({code: code}),
                  // Extended data associated with the reason.  Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type.
                  details:: {
                    local __detailsMixin(details) = __resultMixin({details+: details}),
                    mixinInstance(details):: __detailsMixin(details),
                    // The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes.
                    causes(causes):: if std.type(causes) == "array" then __detailsMixin({causes+: causes}) else __detailsMixin({causes+: [causes]}),
                    causesType:: hidden.meta.v1.statusCause,
                    // The group attribute of the resource associated with the status StatusReason.
                    group(group):: __detailsMixin({group: group}),
                    // The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described).
                    name(name):: __detailsMixin({name: name}),
                    // If specified, the time in seconds before the operation should be retried.
                    retryAfterSeconds(retryAfterSeconds):: __detailsMixin({retryAfterSeconds: retryAfterSeconds}),
                    // UID of the resource. (when there is a single resource which can be described). More info: http://kubernetes.io/docs/user-guide/identifiers#uids
                    uid(uid):: __detailsMixin({uid: uid}),
                  },
                  detailsType:: hidden.meta.v1.statusDetails,
                  // A human-readable description of the status of this operation.
                  message(message):: __resultMixin({message: message}),
                  // A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it.
                  reason(reason):: __resultMixin({reason: reason}),
                },
                resultType:: hidden.meta.v1.status,
              },
              initializersType:: hidden.meta.v1.initializers,
              // Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
              labels(labels):: __metadataMixin({labels+: labels}),
              // Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
              name(name):: __metadataMixin({name: name}),
              // Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
              //
              // Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
              namespace(namespace):: __metadataMixin({namespace: namespace}),
            },
            metadataType:: hidden.meta.v1.objectMeta,
          },
        },
      },
    },
  },
}
