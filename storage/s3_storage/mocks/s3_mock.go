// Code generated by http://github.com/gojuno/minimock (v3.3.6). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/m6yf/bcwork/storage/s3_storage.S3 -o s3_mock.go -n S3Mock -p mocks

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gojuno/minimock/v3"
)

// S3Mock implements s3storage.S3
type S3Mock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetObjectInput          func(bucket string, key string) (ba1 []byte, err error)
	inspectFuncGetObjectInput   func(bucket string, key string)
	afterGetObjectInputCounter  uint64
	beforeGetObjectInputCounter uint64
	GetObjectInputMock          mS3MockGetObjectInput

	funcListS3Objects          func(bucket string, prefix string) (lp1 *s3.ListObjectsV2Output, err error)
	inspectFuncListS3Objects   func(bucket string, prefix string)
	afterListS3ObjectsCounter  uint64
	beforeListS3ObjectsCounter uint64
	ListS3ObjectsMock          mS3MockListS3Objects
}

// NewS3Mock returns a mock for s3storage.S3
func NewS3Mock(t minimock.Tester) *S3Mock {
	m := &S3Mock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetObjectInputMock = mS3MockGetObjectInput{mock: m}
	m.GetObjectInputMock.callArgs = []*S3MockGetObjectInputParams{}

	m.ListS3ObjectsMock = mS3MockListS3Objects{mock: m}
	m.ListS3ObjectsMock.callArgs = []*S3MockListS3ObjectsParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mS3MockGetObjectInput struct {
	mock               *S3Mock
	defaultExpectation *S3MockGetObjectInputExpectation
	expectations       []*S3MockGetObjectInputExpectation

	callArgs []*S3MockGetObjectInputParams
	mutex    sync.RWMutex
}

// S3MockGetObjectInputExpectation specifies expectation struct of the S3.GetObjectInput
type S3MockGetObjectInputExpectation struct {
	mock    *S3Mock
	params  *S3MockGetObjectInputParams
	results *S3MockGetObjectInputResults
	Counter uint64
}

// S3MockGetObjectInputParams contains parameters of the S3.GetObjectInput
type S3MockGetObjectInputParams struct {
	bucket string
	key    string
}

// S3MockGetObjectInputResults contains results of the S3.GetObjectInput
type S3MockGetObjectInputResults struct {
	ba1 []byte
	err error
}

// Expect sets up expected params for S3.GetObjectInput
func (mmGetObjectInput *mS3MockGetObjectInput) Expect(bucket string, key string) *mS3MockGetObjectInput {
	if mmGetObjectInput.mock.funcGetObjectInput != nil {
		mmGetObjectInput.mock.t.Fatalf("S3Mock.GetObjectInput mock is already set by Set")
	}

	if mmGetObjectInput.defaultExpectation == nil {
		mmGetObjectInput.defaultExpectation = &S3MockGetObjectInputExpectation{}
	}

	mmGetObjectInput.defaultExpectation.params = &S3MockGetObjectInputParams{bucket, key}
	for _, e := range mmGetObjectInput.expectations {
		if minimock.Equal(e.params, mmGetObjectInput.defaultExpectation.params) {
			mmGetObjectInput.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetObjectInput.defaultExpectation.params)
		}
	}

	return mmGetObjectInput
}

// Inspect accepts an inspector function that has same arguments as the S3.GetObjectInput
func (mmGetObjectInput *mS3MockGetObjectInput) Inspect(f func(bucket string, key string)) *mS3MockGetObjectInput {
	if mmGetObjectInput.mock.inspectFuncGetObjectInput != nil {
		mmGetObjectInput.mock.t.Fatalf("Inspect function is already set for S3Mock.GetObjectInput")
	}

	mmGetObjectInput.mock.inspectFuncGetObjectInput = f

	return mmGetObjectInput
}

// Return sets up results that will be returned by S3.GetObjectInput
func (mmGetObjectInput *mS3MockGetObjectInput) Return(ba1 []byte, err error) *S3Mock {
	if mmGetObjectInput.mock.funcGetObjectInput != nil {
		mmGetObjectInput.mock.t.Fatalf("S3Mock.GetObjectInput mock is already set by Set")
	}

	if mmGetObjectInput.defaultExpectation == nil {
		mmGetObjectInput.defaultExpectation = &S3MockGetObjectInputExpectation{mock: mmGetObjectInput.mock}
	}
	mmGetObjectInput.defaultExpectation.results = &S3MockGetObjectInputResults{ba1, err}
	return mmGetObjectInput.mock
}

// Set uses given function f to mock the S3.GetObjectInput method
func (mmGetObjectInput *mS3MockGetObjectInput) Set(f func(bucket string, key string) (ba1 []byte, err error)) *S3Mock {
	if mmGetObjectInput.defaultExpectation != nil {
		mmGetObjectInput.mock.t.Fatalf("Default expectation is already set for the S3.GetObjectInput method")
	}

	if len(mmGetObjectInput.expectations) > 0 {
		mmGetObjectInput.mock.t.Fatalf("Some expectations are already set for the S3.GetObjectInput method")
	}

	mmGetObjectInput.mock.funcGetObjectInput = f
	return mmGetObjectInput.mock
}

// When sets expectation for the S3.GetObjectInput which will trigger the result defined by the following
// Then helper
func (mmGetObjectInput *mS3MockGetObjectInput) When(bucket string, key string) *S3MockGetObjectInputExpectation {
	if mmGetObjectInput.mock.funcGetObjectInput != nil {
		mmGetObjectInput.mock.t.Fatalf("S3Mock.GetObjectInput mock is already set by Set")
	}

	expectation := &S3MockGetObjectInputExpectation{
		mock:   mmGetObjectInput.mock,
		params: &S3MockGetObjectInputParams{bucket, key},
	}
	mmGetObjectInput.expectations = append(mmGetObjectInput.expectations, expectation)
	return expectation
}

// Then sets up S3.GetObjectInput return parameters for the expectation previously defined by the When method
func (e *S3MockGetObjectInputExpectation) Then(ba1 []byte, err error) *S3Mock {
	e.results = &S3MockGetObjectInputResults{ba1, err}
	return e.mock
}

// GetObjectInput implements s3storage.S3
func (mmGetObjectInput *S3Mock) GetObjectInput(bucket string, key string) (ba1 []byte, err error) {
	mm_atomic.AddUint64(&mmGetObjectInput.beforeGetObjectInputCounter, 1)
	defer mm_atomic.AddUint64(&mmGetObjectInput.afterGetObjectInputCounter, 1)

	if mmGetObjectInput.inspectFuncGetObjectInput != nil {
		mmGetObjectInput.inspectFuncGetObjectInput(bucket, key)
	}

	mm_params := S3MockGetObjectInputParams{bucket, key}

	// Record call args
	mmGetObjectInput.GetObjectInputMock.mutex.Lock()
	mmGetObjectInput.GetObjectInputMock.callArgs = append(mmGetObjectInput.GetObjectInputMock.callArgs, &mm_params)
	mmGetObjectInput.GetObjectInputMock.mutex.Unlock()

	for _, e := range mmGetObjectInput.GetObjectInputMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.ba1, e.results.err
		}
	}

	if mmGetObjectInput.GetObjectInputMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetObjectInput.GetObjectInputMock.defaultExpectation.Counter, 1)
		mm_want := mmGetObjectInput.GetObjectInputMock.defaultExpectation.params
		mm_got := S3MockGetObjectInputParams{bucket, key}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetObjectInput.t.Errorf("S3Mock.GetObjectInput got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetObjectInput.GetObjectInputMock.defaultExpectation.results
		if mm_results == nil {
			mmGetObjectInput.t.Fatal("No results are set for the S3Mock.GetObjectInput")
		}
		return (*mm_results).ba1, (*mm_results).err
	}
	if mmGetObjectInput.funcGetObjectInput != nil {
		return mmGetObjectInput.funcGetObjectInput(bucket, key)
	}
	mmGetObjectInput.t.Fatalf("Unexpected call to S3Mock.GetObjectInput. %v %v", bucket, key)
	return
}

// GetObjectInputAfterCounter returns a count of finished S3Mock.GetObjectInput invocations
func (mmGetObjectInput *S3Mock) GetObjectInputAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetObjectInput.afterGetObjectInputCounter)
}

// GetObjectInputBeforeCounter returns a count of S3Mock.GetObjectInput invocations
func (mmGetObjectInput *S3Mock) GetObjectInputBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetObjectInput.beforeGetObjectInputCounter)
}

// Calls returns a list of arguments used in each call to S3Mock.GetObjectInput.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetObjectInput *mS3MockGetObjectInput) Calls() []*S3MockGetObjectInputParams {
	mmGetObjectInput.mutex.RLock()

	argCopy := make([]*S3MockGetObjectInputParams, len(mmGetObjectInput.callArgs))
	copy(argCopy, mmGetObjectInput.callArgs)

	mmGetObjectInput.mutex.RUnlock()

	return argCopy
}

// MinimockGetObjectInputDone returns true if the count of the GetObjectInput invocations corresponds
// the number of defined expectations
func (m *S3Mock) MinimockGetObjectInputDone() bool {
	for _, e := range m.GetObjectInputMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetObjectInputMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetObjectInputCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetObjectInput != nil && mm_atomic.LoadUint64(&m.afterGetObjectInputCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetObjectInputInspect logs each unmet expectation
func (m *S3Mock) MinimockGetObjectInputInspect() {
	for _, e := range m.GetObjectInputMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to S3Mock.GetObjectInput with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetObjectInputMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetObjectInputCounter) < 1 {
		if m.GetObjectInputMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to S3Mock.GetObjectInput")
		} else {
			m.t.Errorf("Expected call to S3Mock.GetObjectInput with params: %#v", *m.GetObjectInputMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetObjectInput != nil && mm_atomic.LoadUint64(&m.afterGetObjectInputCounter) < 1 {
		m.t.Error("Expected call to S3Mock.GetObjectInput")
	}
}

type mS3MockListS3Objects struct {
	mock               *S3Mock
	defaultExpectation *S3MockListS3ObjectsExpectation
	expectations       []*S3MockListS3ObjectsExpectation

	callArgs []*S3MockListS3ObjectsParams
	mutex    sync.RWMutex
}

// S3MockListS3ObjectsExpectation specifies expectation struct of the S3.ListS3Objects
type S3MockListS3ObjectsExpectation struct {
	mock    *S3Mock
	params  *S3MockListS3ObjectsParams
	results *S3MockListS3ObjectsResults
	Counter uint64
}

// S3MockListS3ObjectsParams contains parameters of the S3.ListS3Objects
type S3MockListS3ObjectsParams struct {
	bucket string
	prefix string
}

// S3MockListS3ObjectsResults contains results of the S3.ListS3Objects
type S3MockListS3ObjectsResults struct {
	lp1 *s3.ListObjectsV2Output
	err error
}

// Expect sets up expected params for S3.ListS3Objects
func (mmListS3Objects *mS3MockListS3Objects) Expect(bucket string, prefix string) *mS3MockListS3Objects {
	if mmListS3Objects.mock.funcListS3Objects != nil {
		mmListS3Objects.mock.t.Fatalf("S3Mock.ListS3Objects mock is already set by Set")
	}

	if mmListS3Objects.defaultExpectation == nil {
		mmListS3Objects.defaultExpectation = &S3MockListS3ObjectsExpectation{}
	}

	mmListS3Objects.defaultExpectation.params = &S3MockListS3ObjectsParams{bucket, prefix}
	for _, e := range mmListS3Objects.expectations {
		if minimock.Equal(e.params, mmListS3Objects.defaultExpectation.params) {
			mmListS3Objects.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmListS3Objects.defaultExpectation.params)
		}
	}

	return mmListS3Objects
}

// Inspect accepts an inspector function that has same arguments as the S3.ListS3Objects
func (mmListS3Objects *mS3MockListS3Objects) Inspect(f func(bucket string, prefix string)) *mS3MockListS3Objects {
	if mmListS3Objects.mock.inspectFuncListS3Objects != nil {
		mmListS3Objects.mock.t.Fatalf("Inspect function is already set for S3Mock.ListS3Objects")
	}

	mmListS3Objects.mock.inspectFuncListS3Objects = f

	return mmListS3Objects
}

// Return sets up results that will be returned by S3.ListS3Objects
func (mmListS3Objects *mS3MockListS3Objects) Return(lp1 *s3.ListObjectsV2Output, err error) *S3Mock {
	if mmListS3Objects.mock.funcListS3Objects != nil {
		mmListS3Objects.mock.t.Fatalf("S3Mock.ListS3Objects mock is already set by Set")
	}

	if mmListS3Objects.defaultExpectation == nil {
		mmListS3Objects.defaultExpectation = &S3MockListS3ObjectsExpectation{mock: mmListS3Objects.mock}
	}
	mmListS3Objects.defaultExpectation.results = &S3MockListS3ObjectsResults{lp1, err}
	return mmListS3Objects.mock
}

// Set uses given function f to mock the S3.ListS3Objects method
func (mmListS3Objects *mS3MockListS3Objects) Set(f func(bucket string, prefix string) (lp1 *s3.ListObjectsV2Output, err error)) *S3Mock {
	if mmListS3Objects.defaultExpectation != nil {
		mmListS3Objects.mock.t.Fatalf("Default expectation is already set for the S3.ListS3Objects method")
	}

	if len(mmListS3Objects.expectations) > 0 {
		mmListS3Objects.mock.t.Fatalf("Some expectations are already set for the S3.ListS3Objects method")
	}

	mmListS3Objects.mock.funcListS3Objects = f
	return mmListS3Objects.mock
}

// When sets expectation for the S3.ListS3Objects which will trigger the result defined by the following
// Then helper
func (mmListS3Objects *mS3MockListS3Objects) When(bucket string, prefix string) *S3MockListS3ObjectsExpectation {
	if mmListS3Objects.mock.funcListS3Objects != nil {
		mmListS3Objects.mock.t.Fatalf("S3Mock.ListS3Objects mock is already set by Set")
	}

	expectation := &S3MockListS3ObjectsExpectation{
		mock:   mmListS3Objects.mock,
		params: &S3MockListS3ObjectsParams{bucket, prefix},
	}
	mmListS3Objects.expectations = append(mmListS3Objects.expectations, expectation)
	return expectation
}

// Then sets up S3.ListS3Objects return parameters for the expectation previously defined by the When method
func (e *S3MockListS3ObjectsExpectation) Then(lp1 *s3.ListObjectsV2Output, err error) *S3Mock {
	e.results = &S3MockListS3ObjectsResults{lp1, err}
	return e.mock
}

// ListS3Objects implements s3storage.S3
func (mmListS3Objects *S3Mock) ListS3Objects(bucket string, prefix string) (lp1 *s3.ListObjectsV2Output, err error) {
	mm_atomic.AddUint64(&mmListS3Objects.beforeListS3ObjectsCounter, 1)
	defer mm_atomic.AddUint64(&mmListS3Objects.afterListS3ObjectsCounter, 1)

	if mmListS3Objects.inspectFuncListS3Objects != nil {
		mmListS3Objects.inspectFuncListS3Objects(bucket, prefix)
	}

	mm_params := S3MockListS3ObjectsParams{bucket, prefix}

	// Record call args
	mmListS3Objects.ListS3ObjectsMock.mutex.Lock()
	mmListS3Objects.ListS3ObjectsMock.callArgs = append(mmListS3Objects.ListS3ObjectsMock.callArgs, &mm_params)
	mmListS3Objects.ListS3ObjectsMock.mutex.Unlock()

	for _, e := range mmListS3Objects.ListS3ObjectsMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.lp1, e.results.err
		}
	}

	if mmListS3Objects.ListS3ObjectsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmListS3Objects.ListS3ObjectsMock.defaultExpectation.Counter, 1)
		mm_want := mmListS3Objects.ListS3ObjectsMock.defaultExpectation.params
		mm_got := S3MockListS3ObjectsParams{bucket, prefix}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmListS3Objects.t.Errorf("S3Mock.ListS3Objects got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmListS3Objects.ListS3ObjectsMock.defaultExpectation.results
		if mm_results == nil {
			mmListS3Objects.t.Fatal("No results are set for the S3Mock.ListS3Objects")
		}
		return (*mm_results).lp1, (*mm_results).err
	}
	if mmListS3Objects.funcListS3Objects != nil {
		return mmListS3Objects.funcListS3Objects(bucket, prefix)
	}
	mmListS3Objects.t.Fatalf("Unexpected call to S3Mock.ListS3Objects. %v %v", bucket, prefix)
	return
}

// ListS3ObjectsAfterCounter returns a count of finished S3Mock.ListS3Objects invocations
func (mmListS3Objects *S3Mock) ListS3ObjectsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListS3Objects.afterListS3ObjectsCounter)
}

// ListS3ObjectsBeforeCounter returns a count of S3Mock.ListS3Objects invocations
func (mmListS3Objects *S3Mock) ListS3ObjectsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListS3Objects.beforeListS3ObjectsCounter)
}

// Calls returns a list of arguments used in each call to S3Mock.ListS3Objects.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmListS3Objects *mS3MockListS3Objects) Calls() []*S3MockListS3ObjectsParams {
	mmListS3Objects.mutex.RLock()

	argCopy := make([]*S3MockListS3ObjectsParams, len(mmListS3Objects.callArgs))
	copy(argCopy, mmListS3Objects.callArgs)

	mmListS3Objects.mutex.RUnlock()

	return argCopy
}

// MinimockListS3ObjectsDone returns true if the count of the ListS3Objects invocations corresponds
// the number of defined expectations
func (m *S3Mock) MinimockListS3ObjectsDone() bool {
	for _, e := range m.ListS3ObjectsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ListS3ObjectsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterListS3ObjectsCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcListS3Objects != nil && mm_atomic.LoadUint64(&m.afterListS3ObjectsCounter) < 1 {
		return false
	}
	return true
}

// MinimockListS3ObjectsInspect logs each unmet expectation
func (m *S3Mock) MinimockListS3ObjectsInspect() {
	for _, e := range m.ListS3ObjectsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to S3Mock.ListS3Objects with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ListS3ObjectsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterListS3ObjectsCounter) < 1 {
		if m.ListS3ObjectsMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to S3Mock.ListS3Objects")
		} else {
			m.t.Errorf("Expected call to S3Mock.ListS3Objects with params: %#v", *m.ListS3ObjectsMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcListS3Objects != nil && mm_atomic.LoadUint64(&m.afterListS3ObjectsCounter) < 1 {
		m.t.Error("Expected call to S3Mock.ListS3Objects")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *S3Mock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetObjectInputInspect()

			m.MinimockListS3ObjectsInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *S3Mock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *S3Mock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetObjectInputDone() &&
		m.MinimockListS3ObjectsDone()
}