package s3api

import (
	"encoding/xml"
	"errors"
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
)

// Status represents lifecycle configuration status
type ruleStatus string

// Supported status types
const (
	Enabled  ruleStatus = "Enabled"
	Disabled ruleStatus = "Disabled"
)

// Lifecycle - Configuration for bucket lifecycle.
type Lifecycle struct {
	XMLName xml.Name `xml:"LifecycleConfiguration"`
	Rules   []Rule   `xml:"Rule"`
}

// Rule - a rule for lifecycle configuration.
type Rule struct {
	XMLName                        xml.Name                       `xml:"Rule"`
	ID                             string                         `xml:"ID,omitempty"`
	Status                         ruleStatus                     `xml:"Status"`
	Filter                         Filter                         `xml:"Filter,omitempty"`
	Expiration                     Expiration                     `xml:"Expiration,omitempty"`
	Transition                     []Transition                   `xml:"Transition,omitempty"`
	AbortIncompleteMultipartUpload AbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload,omitempty"`
	NoncurrentVersionExpiration    NoncurrentVersionExpiration    `xml:"NoncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransition    []NoncurrentVersionTransition  `xml:"NoncurrentVersionTransition,omitempty"`
}

// AbortIncompleteMultipartUpload - Specifies the days since the initiation of an incomplete
// multipart upload that Amazon S3 will wait before permanently removing all parts of the
// upload.
type AbortIncompleteMultipartUpload struct {
	XMLName             xml.Name `xml:"AbortIncompleteMultipartUpload"`
	DaysAfterInitiation int64    `xml:"DaysAfterInitiation,omitempty"`
}

// NoncurrentVersionExpiration - Specifies when noncurrent object versions expire.
// Upon expiration, Amazon S3 permanently deletes the noncurrent object versions.
type NoncurrentVersionExpiration struct {
	XMLName                 xml.Name `xml:"NoncurrentVersionExpiration"`
	NewerNoncurrentVersions int64    `xml:"NewerNoncurrentVersions,omitempty"`
	NoncurrentDays          int64    `xml:"NoncurrentDays,omitempty"`
}

type NoncurrentVersionTransition struct {
	XMLName                 xml.Name `xml:"NoncurrentVersionTransition"`
	NewerNoncurrentVersions int64    `xml:"NewerNoncurrentVersions,omitempty"`
	NoncurrentDays          int64    `xml:"NoncurrentDays,omitempty"`
	StorageClass            string   `xml:"StorageClass,omitempty"`
}

// Filter - a filter for a lifecycle configuration Rule.
type Filter struct {
	XMLName xml.Name `xml:"Filter"`
	set     bool

	Prefix    string `xml:"Prefix,omitempty"`
	prefixSet bool

	And    And `xml:"And,omitempty"`
	andSet bool

	Tag    Tag `xml:"Tag,omitempty"`
	tagSet bool

	ObjectSizeGreaterThan    int64 `xml:"ObjectSizeGreaterThan,omitempty"`
	objectSizeGreaterThanSet bool

	ObjectSizeLessThan    int64 `xml:"ObjectSizeLessThan,omitempty"`
	objectSizeLessThanSet bool
}

// Prefix holds the prefix xml tag in <Rule> and <Filter>
type Prefix struct {
	string
	set bool
}

// MarshalXML encodes prefix field into an XML form.
func (p Prefix) MarshalXML(e *xml.Encoder, startElement xml.StartElement) error {
	if !p.set {
		return nil
	}
	return e.EncodeElement(p.string, startElement)
}

// MarshalXML encodes prefix field into an XML form.
func (f Filter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.EncodeElement(f.Prefix, xml.StartElement{Name: xml.Name{Local: "Prefix"}}); err != nil {
		return err
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// And - a tag to combine a prefix and multiple tags for lifecycle configuration rule.
type And struct {
	XMLName               xml.Name `xml:"And"`
	Prefix                string   `xml:"Prefix,omitempty"`
	Tags                  []Tag    `xml:"Tag,omitempty"`
	ObjectSizeGreaterThan int      `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    int      `xml:"ObjectSizeLessThan,omitempty"`
}

// Expiration - expiration actions for a rule in lifecycle configuration.
type Expiration struct {
	XMLName      xml.Name           `xml:"Expiration"`
	Days         int                `xml:"Days,omitempty"`
	Date         ExpirationDate     `xml:"Date,omitempty"`
	DeleteMarker ExpireDeleteMarker `xml:"ExpiredObjectDeleteMarker,omitempty"`

	set bool
}

// MarshalXML encodes expiration field into an XML form.
func (e Expiration) MarshalXML(enc *xml.Encoder, startElement xml.StartElement) error {
	if !e.set {
		return nil
	}
	type expirationWrapper Expiration
	return enc.EncodeElement(expirationWrapper(e), startElement)
}

// ExpireDeleteMarker represents value of ExpiredObjectDeleteMarker field in Expiration XML element.
type ExpireDeleteMarker struct {
	val bool
	set bool
}

// MarshalXML encodes delete marker boolean into an XML form.
func (b ExpireDeleteMarker) MarshalXML(e *xml.Encoder, startElement xml.StartElement) error {
	if !b.set {
		return nil
	}
	return e.EncodeElement(b.val, startElement)
}

// ExpirationDate is a embedded type containing time.Time to unmarshal
// Date in Expiration
type ExpirationDate struct {
	time.Time
}

// MarshalXML encodes expiration date if it is non-zero and encodes
// empty string otherwise
func (eDate ExpirationDate) MarshalXML(e *xml.Encoder, startElement xml.StartElement) error {
	if eDate.Time.IsZero() {
		return nil
	}
	return e.EncodeElement(eDate.Format(time.RFC3339), startElement)
}

// Transition - transition actions for a rule in lifecycle configuration.
type Transition struct {
	XMLName      xml.Name  `xml:"Transition"`
	Days         int       `xml:"Days,omitempty"`
	Date         time.Time `xml:"Date,omitempty"`
	StorageClass string    `xml:"StorageClass,omitempty"`

	set bool
}

// MarshalXML encodes transition field into an XML form.
func (t Transition) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if !t.set {
		return nil
	}
	type transitionWrapper Transition
	return enc.EncodeElement(transitionWrapper(t), start)
}

// TransitionDays is a type alias to unmarshal Days in Transition
type TransitionDays int

// GetBucketLifecycleRule return rules of specified bucket
func GetBucketLifecycleRule(bucket string, fc *filer.FilerConf) (rules []Rule) {
	return rules
}

// PutBucketLifecycleRule set rules of specified bucket
func PutBucketLifecycleRule(bucket string, fc *filer.FilerConf, rules []Rule) (err error) {
	//check number of rules
	if !checkRulesNumer(rules) {
		return errors.New(invalidNumber)
	}

	locConf := &filer_pb.FilerConf_PathConf{
		LocationPrefix: "/buckets/" + bucket + "/",
		BucketRules:    make([]*filer_pb.FilerConf_PathConf_BucketRule, len(rules)),
	}
	for i, rule := range rules {
		// check all fields of rule
		if err = rule.checkRuleFields(); err != nil {
			return err
		}
		locConf.BucketRules[i] = &filer_pb.FilerConf_PathConf_BucketRule{
			Id:     rule.ID,
			Status: string(rule.Status),
			Filter: &filer_pb.FilerConf_PathConf_Filter{
				Prefix: rule.Filter.Prefix,
				And: &filer_pb.FilerConf_PathConf_And{
					Prefix: rule.Filter.And.Prefix,
					// Assignment of Tags is done in loop
				},
				Tag: &filer_pb.FilerConf_PathConf_Tag{
					Key:   rule.Filter.Tag.Key,
					Value: rule.Filter.Tag.Value,
				},
			},
			Expiration: &filer_pb.FilerConf_PathConf_Expiration{
				Days: int64(rule.Expiration.Days),
				// Assignment of Date is done below
				DeleteMarker: rule.Expiration.DeleteMarker.val,
			},
			Transition: &filer_pb.FilerConf_PathConf_Transition{
				Days: int64(rule.Transition.Days),
				// Assignment of Date is done below
				StorageClass: rule.Transition.StorageClass,
			},
		}
		if date := rule.Expiration.Date; !date.IsZero() {
			locConf.BucketRules[i].Expiration.Date = date.String()
		}
		if date := rule.Transition.Date; !date.IsZero() {
			locConf.BucketRules[i].Transition.Date = date.String()
		}
		locConf.BucketRules[i].Filter.And.Tags = make([]*filer_pb.FilerConf_PathConf_Tag, len(rule.Filter.And.Tags))
		for j, tag := range rule.Filter.And.Tags {
			locConf.BucketRules[i].Filter.And.Tags[j] = &filer_pb.FilerConf_PathConf_Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	fc.AddLocationConf(locConf)
	return
}

type RuleErrMsg string

// Rule error message
const (
	invalidNumer RuleErrMsg = "Invalid number of rule"
	invalidId    RuleErrMsg = "Invalid id of rule"
	noAction     RuleErrMsg = "No Actions existing in rule"
)

func checkRulesNumer(rules []Rule) bool {
	return len(rules) <= 1000
}

// checkRuleFields verify that whether all fields of rule are correct
func (rule *Rule) checkFields() error {
	// check ID
	if !rule.checkId() {
		return errors.New(invalidId)
	}

	// check Status
	if !rule.checkStatus() {
		return errors.New(InvalidStatus)
	}

	// check Filter
	if err := rule.checkFilter(); err != nil {
		return err
	}

	// check Transitions
	if err := rule.checkTransitions(); err != nil {
		return err
	}

	// check Expiration
	if err := rule.checkExpiration(); err != nil {
		return err
	}

	// check AbortIncompleteMultipartUpload
	if err := rule.checkAbortIncompleteMultipartUpload(); err != nil {
		return err
	}

	// check NoncurrentVersionExpiration
	if err := rule.checkNoncurrentVersionExpiration(); err != nil {
		return err
	}

	// check NoncurrentVersionTransition
	if err := rule.checkNoncurrentVersionTransition(); err != nil {
		return err
	}

	// check action
	if !rule.checkAction() {
		return errors.New(NoAction)
	}
	return nil
}

func (rule *Rule) checkId() bool {
	return len(rule.ID) <= 255
}

func (rule *Rule) checkStatus() bool {
	return rule.Status == Enabled || rule.Status == Disabled
}

func (rule *Rule) checkFilter() error {
	pf := &rule.Filter
	return pf.checkFilterFields()
}

func (rule *Rule) checkTransitions() error {
	// TODO
}

func (rule *Rule) checkExpiration() error {
	// TODO
}

func (rule *Rule) checkAbortIncompleteMultipartUpload() error {
	// TODO
}

func (rule *Rule) checkNoncurrentVersionExpiration() error {
	// TODO
}

func (rule *Rule) checkTransitions() error {
	// TODO
}

func (rule *Rule) checkAction() error {
	// TODO
}

type FilterErrMsg string

// Filter error message
const (
	InvalidObjectSize             FilterErrMsg = "Invalid object size"
	InvalidObjectSizeRelationship FilterErrMsg = "Invalid object relationship"
	InvalidFieldsCoexist          FilterErrMsg = "Specified more than one element in Filter"
	InvalidTag                    FilterErrMsg = "Invalid Tag"
	DuplicateTagKey               FilterErrMsg = "Duplicate Tag Keys"
)

func (filter *Filter) checkFields() (err error) {
	if filter.Prefix != "" {
		filter.prefixSet = true
		if !strings.HasSuffix(filter.Prefix, "/") {
			filter.Prefix += "/"
		}
	}
	if err = filter.checkObjectSize(); err != nil {
		return err
	}
	if err = filter.checkTag(); err != nil {
		return err
	}
	if !filter.checkCoexist() {
		return errors.New(InvalidFieldsCoexist)
	}
	if err = filter.checkAnd(); err != nil {
		return err
	}

	return err
}

func (filter *Filter) checkObjectSize() error {
	sg, sl := filter.ObjectSizeGreaterThan, filter.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(InvalidObjectSize)
	}
	if sg > sl {
		return errors.New(InvalidObjectSizeRelationship)
	}
	if sg != 0 {
		filter.objectSizeGreaterThanSet = true
	}
	if sl != 0 {
		filter.objectSizeLessThanSet = true
	}
	return nil
}

func (filter *Filter) checkTag() error {
	pt := &filter.Tag
	if pt.check() {
		filter.tagSet = true
	}
	return nil
}

func (filter *Filter) checkCoexist() bool {
	var sum int
	if filter.prefixSet {
		sum++
	}
	if filter.objectSizeGreaterThanSet {
		sum++
	}
	if filter.objectSizeLessThanSet {
		sum++
	}
	if filter.tagSet {
		sum++
	}
	if filter.andSet {
		sum++
	}
	return sum < 2
}

func (filter *Filter) checkAnd() error {
	pa := &filter.And
	return pa.checkFields()
}

func (tag *Tag) check() bool {
	return len(tag.Key) > 0
}

func (and *And) checkFields() error {
	if and.Prefix != "" && !strings.HasSuffix(and.Prefix, "/") {
		and.Prefix += "/"
	}
	sg, sl := and.ObjectSizeGreaterThan, and.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(InvalidObjectSize)
	}
	if sg > sl {
		return errors.New(InvalidObjectSizeRelationship)
	}
	// TODO
	var s []string
	for _, tag := range and.Tags {
		if tag.check() {
			if !checkDuplication(s, tag.Key) {
				s = append(s, tag.Key)
			} else {
				return errors.New(DuplicateTagKey)
			}
		}
	}
	return nil
}

func checkDuplication(str []string, t string) bool {
	for _, s := range str {
		if s == t {
			return true
		}
	}
	return false
}
