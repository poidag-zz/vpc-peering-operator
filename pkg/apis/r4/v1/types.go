package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VpcPeeringList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []VpcPeering `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VpcPeering struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              VpcPeeringSpec   `json:"spec"`
	Status            VpcPeeringStatus `json:"status,omitempty"`
}

type VpcPeeringSpec struct {
	PeerOwnerId        string `json:"peer_owner"`
	PeerVpcId          string `json:"peer_vpcid"`
	PeerCIDR           string `json:"peer_cidr"`
	PeerRegion         string `json:"peer_region"`
	AllowDNSResolution bool   `json:"allow_dns_resolution"`
	SourceVpcId        string `json:"source_vpcid"`
}
type VpcPeeringStatus struct {
	Status    string  `json:"status"`
	PeeringId *string `json:"peering_id"`
}
