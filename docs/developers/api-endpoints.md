# API Endpoints

List of available API endpoints.


# Application healthcheck

This endpoint returns the health status of the API application.

| Method        | Path           | Produces          |
| :------------ |:---------------| :-----------------|
| GET           | /api/v1/health | 	application/json |

#### Request Sample
   
```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/health'
```

#### Response Sample
```json
{
  "status": true
}
```

# Latest version

The endpoint returns the latest StatusBay version

| Method        | Path            | Produces          |
| :------------ |:----------------| :-----------------|
| GET           | /api/v1/version |  application/json |

#### Request Sample
   
```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/version'
```

#### Response Sample
```json
{
  "LatestVersion": "1.0.1",
  "LatestReleaseDate": 12345,
  "Outdated": true,
  "Notifications": []
}
```

# Metrics

This endpoint returns metrics query data from a given provider.

| Method        | Path                       | Produces          |
| :------------ |:---------------------------| :-----------------|
| GET           | /api/v1/application/metric | 	application/json |

#### Parameters

- **provider** - the provider of the metrics response. 
- **query** - the metric query.
- **from** - the metric data points start time (unix time).
- **to** - the metric data points end time (unix time).

#### Request Sample

```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/application/metric?provider=datadog&query=sum:web.http.2xx\{*\}.as_count()&from=1577734291&to=1577741491'
```

#### Response Sample
```json
[
  {
    "Metric": "web.http.2xx",
    "Points": [
      [
        1577734320000,
        5615
      ],
      [
        1577734350000,
        5090
      ],
      [
        1577734380000,
        4843
      ],
      [
        1577734410000,
        4997
      ],
      [
        1577734440000,
        4716
      ],
      ...
    ]
  }
]

```

# Alerts

This endpoint returns alerts for a given provider.

| Method        | Path                       | Produces          |
| :------------ |:---------------------------| :-----------------|
| GET           | /api/v1/application/alerts | 	application/json |

#### Parameters

- **provider** - the provider of the alerts response. 
- **tags** - comma separated tags to filter by.
- **from** - the alerts data points start time (unix time).
- **to** - the alerts data points end time (unix time).

#### Request Sample 

```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/application/alerts?provider=pingdom&tags=production,virginia&from=1577734291&to=1577741491'
```

#### Response Sample 
```json
[
    {
        "ID": 5681742,
        "URL": "https://my.pingdom.com/app/reports/uptime#check=5681742",
        "Name": "foo",
        "Periods": [
            {
                "Status": "up",
                "StartUnix": 1578126835,
                "EndUnix": 1578127565
            },
            {
                "Status": "down",
                "StartUnix": 1578127565,
                "EndUnix": 1578128635
            }
        ]
    },
    ...
]

```

# Applications

This endpoint returns a list of applications.

| Method        | Path                            | Produces          |
| :------------ |:--------------------------------| :-----------------|
| GET           | /api/v1/kubernetes/applications | application/json  |

#### Parameters

- **limit**  `(default: 20)` - number of records in the response.
- **offset** `(default: 0)` - the number of records you wish to skip before selecting records.
- **cluster** `(default: "" -> all)` - filter applications by cluster name. (for multiple clusters, use comma separated string )
- **namespace** `(default: "" -> all)` - filter application by namespace. (for multiple namespaces use comma separated string)
- **status** `(default: "" -> all)` - filter application by status. (for multiple statuses use comma separated string)
- **name** `(default: "" -> all)` - filter deployments for a specific application name.
- **orderby** `(default: "time")` - order the records response.
- **sortdirection** `(default: "desc")` - sort direction of the response records.
- **from** `(default: "0")` - filter applications by range of time, start time - unix time.
- **to** `(default: "0")` - filter applications by range of time, end time - unix time.
- **distinct** `(default: "false")` - enable uniqueness on returned results.

#### Request Sample

```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/kubernetes/applications'
```

#### Response Sample
```json
{
  "Records": [
    {
      "Name": "foo",
      "Status": "running",
      "Cluster": "cluster1",
      "Namespace": "staging",
      "DeployBy": "test@example.com",
      "Time": 1580045074
    },
    {
      "Name": "foo 2",
      "Status": "successful",
      "Cluster": "cluster1",
      "Namespace": "staging",
      "DeployBy": "foo2@example.com",
      "Time": 1580045027
    },
    {
      "Name": "foo 3",
      "Status": "timeout",
      "Cluster": "cluster2",
      "Namespace": "default",
      "DeployBy": "root@example.com",
      "Time": 1580044962
    },
    ...
  ],
  "Count": 5000
}

```

# Applications Unique Values

This endpoint returns a unique column values

| Method        | Path                                            | Produces          |
| :------------ |:------------------------------------------------| :-----------------|
| GET           | /api/v1/kubernetes/applications/values/{column} | application/json  |

#### Parameters

- **column** - database column name. available column names i.e: cluster, namespace, status & deploy_by.

#### Request Sample 

```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/kubernetes/applications/values/cluster'
```

#### Response Sample 
```json
[
    "cluster-1",
    "cluster-2",
    "cluster-3",
    "cluster-4",
    ...
]
```

# Application's Deployment Details

This endpoint returns a specific deployment details.

| Method        | Path                                            | Produces          |
| :------------ |:------------------------------------------------| :-----------------|
| GET           | /api/v1/kubernetes/application/{applyID}        | application/json  |

#### Parameters

- **applyID** - Unique apply ID.

#### Request Sample

```bash
$ curl \
  'http://127.0.0.1:8080/api/v1/kubernetes/application/13f77155e111a9bce2a366f25fc9815d0f825517'
```

#### Response Sample
```json
{
  "Name": "example-deployment",
  "DeployBy": "foo@demo.com",
  "Cluster": "telaviv",
  "Namespace": "staging",
  "Status": "running",
  "Time": 1581574816,
  "Details": {
    "Resources": {
      "Deployments": {
        "deployment1": {
          "MetaData": {
            "Name": "deployment1",
            "Namespace": "default",
            "ClusterName": "",
            "SpecHash": 0,
            "Labels": {
              "app.kubernetes.io/managed-by": "me",
              "app.kubernetes.io/name": "deployment1"
            },
            "DesiredState": 1
          },
          "DeploymentEvents": [
            {
              "Message": "Scaled up replica set deployment1-9ff6b5676 to 1",
              "Time": 1580032133000000000,
              "Action": "",
              "ReportingController": "",
              "MarkDescriptions": []
            }
          ],
          "Metrics": null,
          "Pods": {
            "deployment1-9ff6b5676-c7nml": {
              "Phase": "Running",
              "CreationTimestamp": "0001-01-01T00:00:00Z",
              "Events": [
                {
                  "Message": "Successfully assigned default/deployment1-9ff6b5676-c7nml to minikube",
                  "Time": 1580032133000000000,
                  "Action": "Binding",
                  "ReportingController": "default-scheduler",
                  "MarkDescriptions": []
                },
                {
                  "Message": "Container image \"nginx:latest\" already present on machine",
                  "Time": 1580032135000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": []
                },
                {
                  "Message": "Created container nginx",
                  "Time": 1580032135000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": []
                },
                {
                  "Message": "Started container nginx",
                  "Time": 1580032135000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": []
                }
              ]
            }
          },
          "Replicaset": {
            "deployment1-9ff6b5676": {
              "Events": [
                {
                  "Message": "Created pod: deployment1-9ff6b5676-c7nml",
                  "Time": 1580032133000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": []
                }
              ]
            }
          },
          "Status": {
            "ObservedGeneration": 1,
            "Replicas": 1,
            "UpdatedReplicas": 1,
            "ReadyReplicas": 1,
            "AvailableReplicas": 1,
            "UnavailableReplicas": 0
          }
        },
        ...
      },
      "Daemonsets": {
        "deamonser2": {
          "MetaData": {
            "Name": "deamonser2",
            "Namespace": "default",
            "ClusterName": "",
            "SpecHash": 0,
            "Labels": {
              "app": "fluentd-logging"
            },
            "DesiredState": 0
          },
          "DaemonsetEvents": [
            {
              "Message": "Created pod: deamonser2-ddtmx",
              "Time": 1580032384000000000,
              "Action": "",
              "ReportingController": "",
              "MarkDescriptions": null
            }
          ],
          "Pods": {
            "deamonser2-ddtmx": {
              "Phase": "Running",
              "CreationTimestamp": "0001-01-01T00:00:00Z",
              "Events": [
                {
                  "Message": "ContainerCreating",
                  "Time": 1580032385680265000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": null
                },
                {
                  "Message": "Successfully assigned default/deamonser2-ddtmx to minikube",
                  "Time": 1580032383000000000,
                  "Action": "Binding",
                  "ReportingController": "default-scheduler",
                  "MarkDescriptions": null
                },
                {
                  "Message": "Container image \"nginx:latest\" already present on machine",
                  "Time": 1580032385000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": null
                },
                {
                  "Message": "Created container deamonser2",
                  "Time": 1580032385000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": null
                },
                {
                  "Message": "Started container deamonser2",
                  "Time": 1580032385000000000,
                  "Action": "",
                  "ReportingController": "",
                  "MarkDescriptions": null
                }
              ]
            }
          },
          "Status": {
            "ObservedGeneration": 1,
            "Replicas": 0,
            "UpdatedReplicas": 0,
            "ReadyReplicas": 0,
            "AvailableReplicas": 0,
            "UnavailableReplicas": 0
          }
        },
        ...
      }
      ...
    }
  }
}
```