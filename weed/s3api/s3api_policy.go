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
	XMLName         xml.Name `xml:"Expiration"`
	Days            int      `xml:"Days,omitempty"`
	daysSet         bool
	Date            ExpirationDate `xml:"Date,omitempty"`
	dateSet         bool
	DeleteMarker    bool `xml:"ExpiredObjectDeleteMarker,omitempty"`
	deleteMarkerSet bool

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
				DeleteMarker: rule.Expiration.DeleteMarker,
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
	NotImplemet   RuleErrMsg = "Element is not implement"
)

func checkRulesNumer(rules []Rule) bool {
	return len(rules) <= 1000
}

// checkRuleFields verify that whether all fields of rule are correct
func (r Rule) checkFields() error {
	// check ID
	if !r.checkId() {
		return errors.New(string(InvalidId))
	}

	// check Status
	if !r.checkStatus() {
		return errors.New(string(InvalidStatus))
	}

	// check Filter
	if err := r.checkFilter(); err != nil {
		return err
	}

	// check Transitions
	if err := r.checkTransitions(); err != nil {
		return err
	}

	// check Expiration
	if err := r.checkExpiration(); err != nil {
		return err
	}

	// check AbortIncompleteMultipartUpload
	if err := r.checkAbortIncompleteMultipartUpload(); err != nil {
		return err
	}

	// check NoncurrentVersionExpiration
	if err := r.checkNoncurrentVersionExpiration(); err != nil {
		return err
	}

	// check NoncurrentVersionTransition
	if err := r.checkNoncurrentVersionTransition(); err != nil {
		return err
	}

	// check action
	if !r.checkAction() {
		return errors.New(string(NoAction))
	}
	return nil
}

func (r Rule) checkId() bool {
	return len(r.ID) <= 255
}

func (r Rule) checkStatus() bool {
	return r.Status == Enabled || r.Status == Disabled
}

func (r Rule) checkFilter() error {
	pf := &r.Filter
	return pf.checkFields()
}

func (r Rule) checkTransitions() error {
	if len(r.Transition) == 0 {
		return nil
	} else {
		return errors.New(string(NotImplemet))
	}

	//var ppt *Transition // previous pt
	//for i, t := range r.Transition {
	//pt := &t
	//if err := pt.checkFields(); err != nil {
	//return err
	//} else {
	// TODO date & days
	//if i > 1 {
	//if !pt.checkTimeFormat(ppt) {
	//return errors.New(string(InconsistentTimeFormat))
	//}
	//}
	// TODO check storageClass
	//}
	//ppt = pt
	//}
	//return nil
}

func (r Rule) checkExpiration() error {
	return r.Expiration.checkFields()
}

func (r Rule) checkAbortIncompleteMultipartUpload() error {
	if r.AbortIncompleteMultipartUpload.DaysAfterInitiation == 0 {
		return nil
	} else {
		return errors.New(string(NotImplemet))
	}
}

func (r Rule) checkNoncurrentVersionExpiration() error {
	if r.NoncurrentVersionExpiration.checkEmpty() {
		return nil
	} else {
		return errors.New(string(NotImplemet))
	}
}

func (r Rule) checkNoncurrentVersionTransition() error {
	if len(r.NoncurrentVersionTransition) == 0 {
		return nil
	} else {
		return errors.New(string(NotImplemet))
	}
}

func (r Rule) checkAction() bool {
	// TODO: It should use 'return r.Expiration.set || r.Transition.set' if Transition field
	// is implemented
	return r.Expiration.set
}

type FilterErrMsg string

// Filter error message
const (
	InvalidObjectSize             FilterErrMsg = "Invalid object size"
	InvalidObjectSizeRelationship FilterErrMsg = "Invalid object relationship"
	InvalidPrefix                 FilterErrMsg = "Invalid prefix, prefix should end with \\"
	InvalidFieldsCoexist          FilterErrMsg = "Specified more than one element in Filter"
	InvalidTag                    FilterErrMsg = "Invalid Tag"
	DuplicateTagKey               FilterErrMsg = "Duplicate Tag Keys"
	InvalidAnd                    FilterErrMsg = "And Field must includes more than one element"
)

func (f Filter) checkFields() error {
	if err := f.checkCoexist(); err != nil {
		return err
	}
	if !strings.HasSuffix(f.Prefix, "/") {
		return errors.New(string(InvalidPrefix))
	}
	if err := f.checkObjectSize(); err != nil {
		return err
	}
	if err := f.checkAnd(); err != nil {
		return err
	}
	return nil
}

func (f Filter) checkCoexist() error {
	var sum int
	if f.Prefix != "" {
		sum++
	}
	if !f.And.checkEmpty() {
		sum++
	}
	if !f.Tag.checkEmpty() {
		sum++
	}
	if f.ObjectSizeGreaterThan != 0 {
		sum++
	}
	if f.ObjectSizeLessThan != 0 {
		sum++
	}
	if sum > 1 {
		return errors.New(string(InvalidFieldsCoexist))
	} else {
		return nil
	}
}

func (f Filter) checkObjectSize() error {
	sg, sl := f.ObjectSizeGreaterThan, f.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(string(InvalidObjectSize))
	}
	if sg > sl {
		return errors.New(string(InvalidObjectSizeRelationship))
	}
	return nil
}

func (f Filter) checkAnd() error {
	return f.checkFields()
}

func (t Tag) checkEmpty() bool {
	return len(t.Key) > 0
}

func (a And) checkFields() error {
	if a.checkEmpty() {
		return nil
	}
	if a.checkNumber() < 2 {
		return errors.New(string(InvalidAnd))
	}
	if a.Prefix != "" && !strings.HasSuffix(a.Prefix, "/") {
		a.Prefix += "/"
	}
	sg, sl := a.ObjectSizeGreaterThan, a.ObjectSizeLessThan
	if sg < 0 || sg > 5*1e+6 || sl < 0 || sl > 5*1e+6 {
		return errors.New(string(InvalidObjectSize))
	}
	if sg > sl {
		return errors.New(string(InvalidObjectSizeRelationship))
	}
	var s []string
	for _, tag := range a.Tags {
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

func (a And) checkEmpty() bool {
	if a.Prefix == "" && a.ObjectSizeGreaterThan == 0 && a.ObjectSizeLessThan == 0 &&
		len(a.Tags) == 0 {
		return true
	}
	return false
}

func (a And) checkNumber() (num int) {
	if a.Prefix != "" {
		num++
	}
	if a.ObjectSizeGreaterThan > 0 {
		num++
	}
	if a.ObjectSizeLessThan > 0 {
		num++
	}
	if len(a.Tags) > 0 {
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

func (t Transition) checkFields() error {
	if !t.checkElements() {
		return errors.New(string(InvalidElements))
	}
	// TODO
	return nil
}

func (t Transition) checkElements() bool {
	if t.StorageClass != "" {
		if !t.Date.IsZero() && t.Days == 0 {
			t.dateSet = true
			return true
		}
		if t.Days != 0 && t.Date.IsZero() {
			t.daysSet = true
			return true
		}
	}
	return false
}

func (t Transition) checkTimeFormat(pt Transition) bool {
	return t.dateSet == pt.dateSet && t.daysSet == pt.daysSet
}

func (nve NoncurrentVersionExpiration) checkEmpty() bool {
	return nve.NewerNoncurrentVersions == 0 &&
		nve.NoncurrentDays == 0
}

type ExpirationErrMsg string

// Expiration error message
const (
	CoexistElementsError ExpirationErrMsg = "Date, Days and DeleteMarker can't specified together"
)

func (e Expiration) checkFields() error {
	if err := e.checkCoexist(); err != nil {
		return err
	}
	return nil
}

func (e Expiration) checkCoexist() error {
	sum := 0
	if e.Days != 0 {
		e.daysSet = true
		sum++
	}
	if !e.Date.IsZero() {
		e.dateSet = true
		sum++
	}
	if e.DeleteMarker {
		e.deleteMarkerSet = true
		sum++
	}
	if sum > 1 {
		return errors.New(string(CoexistElementsError))
	}
	e.set = true
	return nil
}
