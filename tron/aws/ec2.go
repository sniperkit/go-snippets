package aws

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/smartystreets/go-aws-auth"
)

const (
	ebEnvIDTagKey = "elasticbeanstalk:environment-id"
)

var (
	myInstanceID = ""
)

type InstanceStateType struct {
	Code int    `xml:"code"`
	Name string `xml:"name"`
}
type PlacementResponseType struct {
	AvailabilityZone string `xml:"availabilityZone"`
	GroupName        string `xml:"groupName"`
	Tenancy          string `xml:"tenancy"`
}
type GroupItemType struct {
	GroupID   string `xml:"groupId"`
	GroupName string `xml:"groupName"`
}
type EbsInstanceBlockDeviceMappingResponseType struct {
	VolumeID            string    `xml:"volumeId"`
	Status              string    `xml:"status"`
	AttachTime          time.Time `xml:"attachTime"`
	DeleteOnTermination bool      `xml:"deleteOnTermination"`
}
type InstanceBlockDeviceMappingResponseItemType struct {
	DeviceName string                                    `xml:"deviceName"`
	EBS        EbsInstanceBlockDeviceMappingResponseType `xml:"ebs"`
}
type ResourceTagSetItemType struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}
type TagSetItemType struct {
	ResourceID   string `xml:"resourceId"`
	ResourceType string `xml:"resourceType"`
	Key          string `xml:"key"`
	Value        string `xml:"value"`
}
type InstanceNetworkInterfaceAttachmentType struct {
	AttachmentID        string    `xml:"attachmentID"`
	DeviceIndex         int       `xml:"deviceIndex"`
	Status              string    `xml:"status"`
	AttachTime          time.Time `xml:"attachTime"`
	DeleteOnTermination bool      `xml:"deleteOnTermination"`
}
type InstanceNetworkInterfaceAssociationType struct {
	PublicIP      string `xml:"publicIp"`
	PublicDNSName string `xml:"publicDnsName"`
	IPOwnerID     string `xml:"ipOwnerId"`
}
type InstancePrivateIpAddressesSetItemType struct {
	PrivateIPAddress string                                  `xml:"privateIpAddress"`
	PrivateDNSName   string                                  `xml:"privateDnsName"`
	Primary          bool                                    `xml:"primary"`
	Association      InstanceNetworkInterfaceAssociationType `xml:"association"`
}
type InstanceNetworkInterfaceSetItemType struct {
	NetworkInterfaceID    string                                  `xml:"networkInterfaceId"`
	SubnetID              string                                  `xml:"subnetId"`
	VPCID                 string                                  `xml:"vpcId"`
	description           string                                  `xml:"description"`
	OwnerID               string                                  `xml:"ownerId"`
	Status                string                                  `xml:"status"`
	MacAddress            string                                  `xml:"macAddress"`
	PrivateIPAddress      string                                  `xml:"privateIpAddress"`
	PrivateDNSName        string                                  `xml:"privateDnsName"`
	SourceDestCheck       string                                  `xml:"sourceDestCheck"`
	GroupSet              []GroupItemType                         `xml:"groupSet>item"`
	Attachment            InstanceNetworkInterfaceAttachmentType  `xml:"attachment"`
	Association           InstanceNetworkInterfaceAssociationType `xml:"association"`
	PrivateIPAddressesSet []InstancePrivateIpAddressesSetItemType `xml:"privateIpAddressesSet>item"`
}
type IamInstanceProfileResponseType struct {
	ARN string `xml:"arn"`
	ID  string `xml:"id"`
}
type RunningInstancesItemType struct {
	InstanceID            string                                       `xml:"instanceId"`
	ImageID               string                                       `xml:"imageId"`
	InstanceState         InstanceStateType                            `xml:"instanceState"`
	PrivateDNSName        string                                       `xml:"privateDnsName"`
	DNSName               string                                       `xml:"dnsName"`
	KeyName               string                                       `xml:"keyName"`
	AMILaunchIndex        int                                          `xml:"amiLaunchIndex"`
	InstanceType          string                                       `xml:"instanceType"`
	LaunchTime            time.Time                                    `xml:"launchTime"`
	Placement             PlacementResponseType                        `xml:"placement"`
	SubnetID              string                                       `xml:"subnetId"`
	VPCID                 string                                       `xml:"vpcId"`
	PrivateIPAddress      string                                       `xml:"privateIpAddress"`
	IPAddress             string                                       `xml:"ipAddress"`
	SourceDestCheck       bool                                         `xml:"sourceDestCheck"`
	GroupSet              []GroupItemType                              `xml:"groupSet>item"`
	Architecture          string                                       `xml:"architecture"`
	RootDeviceType        string                                       `xml:"rootDeviceType"`
	RootDeviceName        string                                       `xml:"rootDeviceName"`
	BlockDeviceMapping    []InstanceBlockDeviceMappingResponseItemType `xml:"blockDeviceMapping>item"`
	InstanceLifecycle     string                                       `xml:"instanceLifecycle"`
	SpotInstanceRequestId string                                       `xml:"spotInstanceRequestId"`
	VirtualizationType    string                                       `xml:"virtualizationType"`
	ClientToken           string                                       `xml:"clientToken"`
	TagSet                []ResourceTagSetItemType                     `xml:"tagSet>item"`
	Hypervisor            string                                       `xml:"hypervisor"`
	NetworkInterfaceSet   []InstanceNetworkInterfaceSetItemType        `xml:"networkInterfaceSet>item"`
	IAMInstanceProfile    IamInstanceProfileResponseType               `xml:"iamInstanceProfile"`
	EBSOptimized          bool                                         `xml:"ebsOptimized"`
	SRIOVNetSupport       string                                       `xml:"sriovNetSupport"`
}
type ReservationInfoType struct {
	ReservationID string                     `xml:"reservationId"`
	OwnerID       string                     `xml:"ownerId"`
	InstancesSet  []RunningInstancesItemType `xml:"instancesSet>item"`
	RequesterID   string                     `xml:"requesterId"`
}

func InstanceID() (string, error) {
	if myInstanceID != "" {
		return myInstanceID, nil
	}
	b, err := httpGet("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", err
	}
	myInstanceID = b
	return myInstanceID, nil
}

func LocalIPv4() (string, error) {
	b, err := httpGet("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		return "", err
	}
	return b, nil
}

func PublicIPv4() (string, error) {
	b, err := httpGet("http://169.254.169.254/latest/meta-data/public-ipv4")
	if err != nil {
		return "", err
	}
	return b, nil
}

type InstanceIdentityResult struct {
	InstanceID       string    `json:"instanceId"`
	ImageID          string    `json:"imageId"`
	Architecture     string    `json:"architecture"`
	PendingTime      time.Time `json:"pendingTime"`
	InstanceType     string    `json:"instanceType"`
	AccountID        string    `json:"AccountId"`
	Region           string    `json:"region"`
	Version          string    `json:"version"`
	AvailabilityZone string    `json:"availabilityZone"`
	PrivateIP        string    `json:"privateIp"`
}

func InstanceIdentity() (*InstanceIdentityResult, error) {
	resp, err := http.Get("http://169.254.169.254/latest/dynamic/instance-identity/document")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	j := &InstanceIdentityResult{}
	if err := json.NewDecoder(resp.Body).Decode(j); err != nil {
		return nil, err
	}
	return j, nil
}

func detectRegion() (string, error) {
	ii, err := InstanceIdentity()
	if err != nil {
		return "", err
	}
	return ii.Region, nil
}

func detectRegionMust() string {
	r, err := detectRegion()
	if err != nil {
		glog.Fatalf("%v", err)
	}
	return r
}

func EbEnvID() (string, error) {
	id, err := InstanceID()
	if err != nil {
		return "", err
	}
	v := url.Values{}
	v.Set("Action", "DescribeTags")
	v.Set("Version", "2014-06-15")
	v.Set("Filter.1.Name", "resource-id")
	v.Set("Filter.1.Value.1", id)
	v.Set("Filter.2.Name", "key")
	v.Set("Filter.2.Value.1", ebEnvIDTagKey)
	urlStr := "https://ec2." + Region + ".amazonaws.com/?" + v.Encode()
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	awsauth.Sign4(req, Credentials())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	res := struct {
		XMLName   xml.Name         `xml:"DescribeTagsResponse"`
		RequestID string           `xml:"requestId"`
		TagSet    []TagSetItemType `xml:"tagSet>item"`
		NextToken string           `xml:"nextToken"`
	}{}
	err = xml.Unmarshal(b, &res)
	if err != nil {
		return "", err
	}
	m := make(map[string]string)
	for _, t := range res.TagSet {
		m[t.Key] = t.Value
	}
	envID := m[ebEnvIDTagKey]
	if envID == "" {
		return "", fmt.Errorf("no value for key %s in response %s", ebEnvIDTagKey, b)
	}
	return envID, nil
}

// Instances returns IP addresses of all the instances belonging to our
// Elasticbeakstalk environment.
func Instances(nextToken string) (ips []string, nt string, err error) {
	envID, err := EbEnvID()
	if err != nil {
		return nil, "", err
	}
	v := url.Values{}
	v.Set("Action", "DescribeInstances")
	v.Set("Version", "2014-06-15")
	v.Set("Filter.1.Name", "tag:"+ebEnvIDTagKey)
	v.Set("Filter.1.Value.1", envID)
	v.Set("Filter.2.Name", "instance-state-name")
	v.Set("Filter.2.Value.1", "running")
	if nextToken != "" {
		v.Set("nextToken", nextToken)
	}
	urlStr := "https://ec2." + Region + ".amazonaws.com/?" + v.Encode()
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, "", err
	}
	awsauth.Sign4(req, Credentials())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	res := struct {
		XMLName        xml.Name              `xml:"DescribeInstancesResponse"`
		RequestID      string                `xml:"requestId"`
		ReservationSet []ReservationInfoType `xml:"reservationSet>item"`
		NextToken      string                `xml:"nextToken"`
	}{}
	err = xml.Unmarshal(b, &res)
	if err != nil {
		return nil, "", err
	}
	nt = res.NextToken
	for _, r := range res.ReservationSet {
		for _, i := range r.InstancesSet {
			ips = append(ips, i.PrivateIPAddress)
		}
	}
	return ips, nt, nil
}

func httpGet(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
