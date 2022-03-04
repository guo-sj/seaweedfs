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
	XMLName xml.Name `xml:"Transition"`
	Days    int      `xml:"Days,omitempty"`
	daysSet bool

	Date         time.Time `xml:"Date,omitempty"`
	dateSet      bool
	StorageClass string `xml:"StorageClass,omitempty"`

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
		return errors.New(string(InvalidNumber))
	}

	locConf := &filer_pb.FilerConf_PathConf{
		LocationPrefix: "/buckets/" + bucket + "/",
		BucketRules:    make([]*filer_pb.FilerConf_PathConf_BucketRule, len(rules)),
	}
	for i, rule := range rules {
		// check all fields of rule
		if err = rule.checkFields(); err != nil {
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
			//Transition: &filer_pb.FilerConf_PathConf_Transition{
			//Days: int64(rule.Transition.Days),
			// Assignment of Date is done below
			//StorageClass: rule.Transition.StorageClass,
			//},
		}
		if date := rule.Expiration.Date; !date.IsZero() {
			locConf.BucketRules[i].Expiration.Date = date.String()
		}
		//if date := rule.Transition.Date; !date.IsZero() {
		//locConf.BucketRules[i].Transition.Date = date.String()
		//}
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
	InvalidNumber RuleErrMsg = "Invalid number of rule"
	InvalidId     RuleErrMsg = "Invalid id of rule"
	InvalidStatus RuleErrMsg = "Invalid status of rule"
	NoAction      RuleErrMsg = "No Actions existing in rule"
)

func checkRulesNumer(rules []Rule) bool {
	return len(rules) <= 1000
}

// checkRuleFields verify that whether all fields of rule are correct
func (rule *Rule) checkFields() error {
	// check ID
	if !rule.checkId() {
		return errors.New(string(InvalidId))
	}

	// check Status
	if !rule.checkStatus() {
		return errors.New(string(InvalidStatus))
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
		return errors.New(string(NoAction))
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
	return pf.checkFields()
}

func (rule *Rule) checkTransitions() error {
	if len(rule.Transition) == 0 {
		return nil
	}

	var ppt *Transition // previous pt
	for i, t := range rule.Transition {
		pt := &t
		if err := pt.checkFields(); err != nil {
			return err
		} else {
			// TODO date & days
			if i > 1 {
				if !pt.checkTimeFormat(ppt) {
					return errors.New(string(InconsistentTimeFormat))
				}
			}
			// TODO check storageClass
		}
		ppt = pt
	}
	return nil
}

func (rule *Rule) checkExpiration() error {
	// TODO
	return nil
}

func (rule *Rule) checkAbortIncompleteMultipartUpload() error {
	// TODO
	return nil
}

func (rule *Rule) checkNoncurrentVersionExpiration() error {
	// TODO
	return nil
}

func (rule *Rule) checkNoncurrentVersionTransition() error {
	// TODO
	return nil
}

func (rule *Rule) checkAction() bool {
	// TODO
	return false
}

type FilterErrMsg string

// Filter error message
const (
	InvalidObjectSize             FilterErrMsg = "Invalid object size"
	InvalidObjectSizeRelationship FilterErrMsg = "Invalid object relationship"
	InvalidFieldsCoexist          FilterErrMsg = "Specified more than one element in Filter"
	InvalidTag                    FilterErrMsg = "Invalid Tag"
	DuplicateTagKey               FilterErrMsg = "Duplicate Tag Keys"
	InvalidAnd                    FilterErrMsg = "And Field must includes more than one element"
)

func (filter *Filter) checkFields() error {
	if filter.checkEmpty() {
		return nil
	}
	if !strings.HasSuffix(filter.Prefix, "/") {
		filter.Prefix += "/"
	}
	if err := filter.checkObjectSize(); err != nil {
		return err
	}
	if err := filter.checkAnd(); err != nil {
		return err
	}
	if filter.checkCoexist() {
		return errors.New(string(InvalidFieldsCoexist))
	}
	return nil
}

func (filter *Filter) checkEmpty() bool {
	var sum int
	if filter.Prefix != "" {
		filter.prefixSet = true
		sum++
	}
	if pa := &filter.And; !pa.checkEmpty() {
		filter.andSet = true
		sum++
	}
	if pt := &filter.Tag; !pt.checkEmpty() {
		filter.tagSet = true
		sum++
	}
	if filter.ObjectSizeGreaterThan != 0 {
		filter.objectSizeGreaterThanSet = true
		sum++
	}
	if filter.ObjectSizeLessThan != 0 {
		filter.objectSizeLessThanSet = true
		sum++
	}
	return sum == 0
}

func (filter *Filter) checkObjectSize() error {
	sg, sl := filter.ObjectSizeGreaterThan, filter.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(string(InvalidObjectSize))
	}
	if sg > sl {
		return errors.New(string(InvalidObjectSizeRelationship))
	}
	return nil
}

func (filter *Filter) checkAnd() error {
	pa := &filter.And
	return pa.checkFields()
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
	return sum > 1
}

func (tag *Tag) checkEmpty() bool {
	return len(tag.Key) > 0
}

func (and *And) checkFields() error {
	if and.checkEmpty() {
		return nil
	}
	if and.checkNumber() < 2 {
		return errors.New(string(InvalidAnd))
	}
	if and.Prefix != "" && !strings.HasSuffix(and.Prefix, "/") {
		and.Prefix += "/"
	}
	sg, sl := and.ObjectSizeGreaterThan, and.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(string(InvalidObjectSize))
	}
	if sg > sl {
		return errors.New(string(InvalidObjectSizeRelationship))
	}
	var s []string
	for _, tag := range and.Tags {
		if tag.checkEmpty() {
			if !checkDuplication(s, tag.Key) {
				s = append(s, tag.Key)
			} else {
				return errors.New(string(DuplicateTagKey))
			}
		}
	}
	return nil
}

func (and *And) checkEmpty() bool {
	if and.Prefix == "" && and.ObjectSizeGreaterThan == 0 && and.ObjectSizeLessThan == 0 &&
		len(and.Tags) == 0 {
		return true
	}
	return false
}

func (and *And) checkNumber() (num int) {
	if and.Prefix != "" {
		num++
	}
	if and.ObjectSizeGreaterThan > 0 {
		num++
	}
	if and.ObjectSizeLessThan > 0 {
		num++
	}
	if len(and.Tags) > 0 {
		num++
	}
	return num
}

func checkDuplication(str []string, t string) bool {
	for _, s := range str {
		if s == t {
			return true
		}
	}
	return false
}

type TransitionErrMsg string

// Transition error message
const (
	InvalidElements        TransitionErrMsg = "Invalid elements in Transition action "
	InconsistentTimeFormat TransitionErrMsg = "Mixed Date and Days in Transition action"
)

func (transition *Transition) checkFields() error {
	if !transition.checkElements() {
		return errors.New(string(InvalidElements))
	}
	// TODO
	return nil
}

func (transition *Transition) checkElements() bool {
	if transition.StorageClass != "" {
		if !transition.Date.IsZero() && transition.Days == 0 {
			transition.dateSet = true
			return true
		}
		if transition.Days != 0 && transition.Date.IsZero() {
			transition.daysSet = true
			return true
		}
	}
	return false
}

func (transition *Transition) checkTimeFormat(pt *Transition) bool {
	return transition.dateSet == pt.dateSet && transition.daysSet == pt.daysSet
}
