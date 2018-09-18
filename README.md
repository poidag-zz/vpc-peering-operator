# VPC Peering Operator

The VPC Peering operator for Kubernetes provides a way to natively define a vpc peering as a Kubernetes object and handles the lifecycle around the Peering and Routing for a VPC.

The premise of this operator is to serve as a self service tool to allow users running in a multi tenant cluster to manage peerings to other AWS VPC's for consumption of their resources.

## CustomResourceDefinitions

The Operator acts on the following [custom resource definitions (CRDs)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/):

- **`VpcPeering`**, which defines a desired VPC Peering.
  The Operator Creates a VPC Peering request upon creation of a VpcPeering CRD. A configurable wait timeout is defined to wait for an accept from another account.

| Parameter                 | Description                                   | Default |
| ------------------------- | --------------------------------------------- | ------- |
| `Spec.PeerOwnerId`        | The account ID owning the VPC to be peered to | `nil`   |
| `Spec.PeerVpcId`          | The VPC ID of the VPC to peer to              | `nil`   |
| `Spec.PeerCIDR`           | The CIDR of the VPC to peer to                | `nil`   |
| `Spec.PeerRegion`         | The region the peer vpc exists within         | `nil`   |
| `Spec.AllowDNSResolution` | The region the peer vpc exists within         | `true`  |
| `Spec.SourceVpcId`        | The VPC ID the operator is running within     | `nil`   |
| `Status.PeeringId`        | The Peering connection ID once created        | `nil`   |

An example is shown in `example/cr.yaml`

## Installation

The Nodes running the Operator require an IAM Instance profile to be associated with the following policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "ec2:CreateRoute",
        "ec2:DeleteVpcPeeringConnection",
        "ec2:DeleteRoute",
        "ec2:CreateVpcPeeringConnection"
      ],
      "Resource": [
        "arn:aws:ec2:*:*:vpc-peering-connection/*",
        "arn:aws:ec2:*:*:route-table/*",
        "arn:aws:ec2:*:*:vpc/*"
      ]
    },
    {
      "Sid": "VisualEditor1",
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeVpcPeeringConnections",
        "ec2:DescribeRouteTables"
      ],
      "Resource": "*"
    }
  ]
}
```

Install the Operator inside a cluster by running the following command:

```sh
kubectl apply -f deploy/
```

> Note: make sure to adapt the namespace in the ClusterRoleBinding if deploying in another namespace than the default namespace.

Create an instance of a VPC Peering CRD

> Note: make sure to adapt the values of the CR defined in `example/cr.yaml`.

```sh
kubectl apply -f example/cr.yaml
```

## Removal

To remove the operator, first delete any custom resources you created in each namespace (please note this will remove routes in routetables associated with the VPC and delete peering connections).

```sh
for n in $(kubectl get namespaces -o jsonpath={..metadata.name}); do
  kubectl delete --all --namespace=$n vpcpeering
done
```

After a couple of minutes you can go ahead and remove the operator itself.

```sh
kubectl delete -f bundle.yaml
```

## Configuration

Configuration is built through environment variables currently defined in `deploy/operator.yaml`
Below are the available configuration options

| Environment Variable   | Description                                                      | Default                |
| ---------------------- | ---------------------------------------------------------------- | ---------------------- |
| `MANAGE_ROUTES`        | Maintain routes in VPC route tables for the peering              | `true`                 |
| `OPERATOR_NAME`        | The name of the operator                                         | `vpc-peering-operator` |
| `WATCH_ALL_NAMESPACES` | Override the SDK and listen to events in all namespaces          | `false`                |
| `POLLER_RETRIES`       | The amount of retries for waiting for a peering to become active | `5`                    |
| `POLLER_WAIT_SECONDS`  | The number of seconds to wait between retries                    | `60`                   |
| `WATCH_NAMESPACE`      | The namespace to watch for CRD events                            | `metadata.namespace`   |
