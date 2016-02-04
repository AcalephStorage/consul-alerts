package stress

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/influxdb/influxdb/client/v2"
)

// AbstractTag is a struct that abstractly
// defines a tag
type AbstractTag struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

// AbstractTags is a slice of abstract tags
type AbstractTags []AbstractTag

// Template returns a templated string of tags
func (t AbstractTags) Template() string {
	var buf bytes.Buffer
	for i, tag := range t {
		if i == 0 {
			buf.Write([]byte(fmt.Sprintf("%v=%v-%%v,", tag.Key, tag.Value)))
		} else {
			buf.Write([]byte(fmt.Sprintf("%v=%v,", tag.Key, tag.Value)))
		}
	}

	b := buf.Bytes()
	b = b[0 : len(b)-1]

	return string(b)
}

// AbstractField is a struct that abstractly
// defines a field
type AbstractField struct {
	Key  string `toml:"key"`
	Type string `toml:"type"`
}

// AbstractFields is a slice of abstract fields
type AbstractFields []AbstractField

// Template returns a templated string of fields
func (f AbstractFields) Template() (string, []string) {
	var buf bytes.Buffer
	a := make([]string, len(f))
	for i, field := range f {
		buf.Write([]byte(fmt.Sprintf("%v=%%v,", field.Key)))
		a[i] = field.Type
	}

	b := buf.Bytes()
	b = b[0 : len(b)-1]

	return string(b), a
}

// BasicPointGenerator implements the PointGenerator interface
type BasicPointGenerator struct {
	PointCount  int            `toml:"point_count"`
	Tick        string         `toml:"tick"`
	Jitter      bool           `toml:"jitter"`
	Measurement string         `toml:"measurement"`
	SeriesCount int            `toml:"series_count"`
	Tags        AbstractTags   `toml:"tag"`
	Fields      AbstractFields `toml:"field"`
	StartDate   string         `toml:"start_date"`
	time        time.Time
	mu          sync.Mutex
}

// typeArr accepts a string array of types and
// returns an array of equal length where each
// element of the array is an instance of the type
// expressed in the string array.
func typeArr(a []string) []interface{} {
	i := make([]interface{}, len(a))
	for j, ty := range a {
		var t string
		switch ty {
		case "float64":
			t = fmt.Sprintf("%v", rand.Intn(1000))
		case "int":
			t = fmt.Sprintf("%vi", rand.Intn(1000))
		case "bool":
			b := rand.Intn(2) == 1
			t = fmt.Sprintf("%t", b)
		default:
			t = fmt.Sprintf("%v", rand.Intn(1000))
		}
		i[j] = t
	}

	return i
}

// Template returns a function that returns a pointer to a Pnt.
func (b *BasicPointGenerator) Template() func(i int, t time.Time) *Pnt {
	ts := b.Tags.Template()
	fs, fa := b.Fields.Template()
	tmplt := fmt.Sprintf("%v,%v %v %%v", b.Measurement, ts, fs)

	return func(i int, t time.Time) *Pnt {
		p := &Pnt{}
		arr := []interface{}{i}
		arr = append(arr, typeArr(fa)...)
		arr = append(arr, t.UnixNano())

		str := fmt.Sprintf(tmplt, arr...)
		p.Set([]byte(str))
		return p
	}
}

// Pnt is a struct that implements the Point interface.
type Pnt struct {
	line []byte
}

// Set sets the internal state for a Pnt.
func (p *Pnt) Set(b []byte) {
	p.line = b
}

// Next generates very simple points very
// efficiently.
// TODO: Take this out
func (p *Pnt) Next(i int, t time.Time) {
	p.line = []byte(fmt.Sprintf("a,b=c-%v v=%v", i, i))
}

// Line returns a byte array for a point
// in line protocol format.
func (p Pnt) Line() []byte {
	return p.line
}

// Graphite returns a byte array for a point
// in graphite format.
func (p Pnt) Graphite() []byte {
	// TODO: Implement
	return []byte("")
}

// OpenJSON returns a byte array for a point
// in opentsdb json format
func (p Pnt) OpenJSON() []byte {
	// TODO: Implement
	return []byte("")
}

// OpenTelnet returns a byte array for a point
// in opentsdb-telnet format
func (p Pnt) OpenTelnet() []byte {
	// TODO: Implement
	return []byte("")
}

// Generate returns a point channel. Implements the
// Generate method for the PointGenerator interface
func (b *BasicPointGenerator) Generate() (<-chan Point, error) {
	// TODO: should be 1.5x batch size
	c := make(chan Point, 15000)

	tmplt := b.Template()

	go func(c chan Point) {
		defer close(c)

		start, err := time.Parse("2006-Jan-02", b.StartDate)
		if err != nil {
			fmt.Println(err)
			return
		}

		b.mu.Lock()
		b.time = start
		b.mu.Unlock()

		tick, err := time.ParseDuration(b.Tick)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < b.PointCount; i++ {
			b.mu.Lock()
			b.time = b.time.Add(tick)
			b.mu.Unlock()

			for j := 0; j < b.SeriesCount; j++ {
				p := tmplt(j, b.time)

				c <- *p
			}
		}
	}(c)

	return c, nil
}

// Time returns the timestamp for the latest points
// that are being generated. Implements the Time method
// for the PointGenerator interface.
func (b *BasicPointGenerator) Time() time.Time {
	defer b.mu.Unlock()
	b.mu.Lock()
	t := b.time
	return t
}

// BasicClient implements the InfluxClient
// interface.
type BasicClient struct {
	Enabled       bool   `toml:"enabled"`
	Address       string `toml:"address"`
	Database      string `toml:"database"`
	Precision     string `toml:"precision"`
	BatchSize     int    `toml:"batch_size"`
	BatchInterval string `toml:"batch_interval"`
	Concurrency   int    `toml:"concurrency"`
	SSL           bool   `toml:"ssl"`
	Format        string `toml:"format"`
}

// Batch groups together points
func (c *BasicClient) Batch(ps <-chan Point, r chan<- response) error {
	if !c.Enabled {
		return nil
	}

	var buf bytes.Buffer
	var wg sync.WaitGroup
	counter := NewConcurrencyLimiter(c.Concurrency)

	interval, err := time.ParseDuration(c.BatchInterval)
	if err != nil {
		return err
	}

	ctr := 0

	for p := range ps {
		b := p.Line()
		ctr++

		buf.Write(b)
		buf.Write([]byte("\n"))

		if ctr%c.BatchSize == 0 && ctr != 0 {
			b := buf.Bytes()

			// Trimming the trailing newline character
			b = b[0 : len(b)-1]

			wg.Add(1)
			counter.Increment()
			go func(byt []byte) {
				defer wg.Done()

				rs, err := c.send(byt)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(interval)

				counter.Decrement()
				r <- rs
			}(b)

			var temp bytes.Buffer
			buf = temp
		}

	}

	wg.Wait()

	return nil
}

// post sends a post request with a payload of points
func post(url string, datatype string, data io.Reader) (*http.Response, error) {
	resp, err := http.Post(url, datatype, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	return resp, nil
}

// Send calls post and returns a response
func (c *BasicClient) send(b []byte) (response, error) {
	instanceURL := fmt.Sprintf("http://%v/write?db=%v&precision=%v", c.Address, c.Database, c.Precision)

	t := NewTimer()
	resp, err := post(instanceURL, "application/x-www-form-urlencoded", bytes.NewBuffer(b))
	t.StopTimer()
	if err != nil {
		return response{Timer: t}, err
	}

	r := response{
		Resp:  resp,
		Time:  time.Now(),
		Timer: t,
	}

	return r, nil
}

// BasicQuery implements the QueryGenerator interface
type BasicQuery struct {
	Template   Query `toml:"template"`
	QueryCount int   `toml:"query_count"`
	time       time.Time
}

// QueryGenerate returns a Query channel
func (q *BasicQuery) QueryGenerate(now func() time.Time) (<-chan Query, error) {
	c := make(chan Query, 0)

	go func(chan Query) {
		defer close(c)

		for i := 0; i < q.QueryCount; i++ {
			c <- Query(fmt.Sprintf(string(q.Template), i))
		}

	}(c)

	return c, nil
}

// SetTime sets the internal state of time
func (q *BasicQuery) SetTime(t time.Time) {
	q.time = t
	return
}

// BasicQueryClient implements the QueryClient interface
type BasicQueryClient struct {
	Enabled       bool   `toml:"enabled"`
	Address       string `toml:"address"`
	Database      string `toml:"database"`
	QueryInterval string `toml:"query_interval"`
	Concurrency   int    `toml:"concurrency"`
	client        client.Client
}

// Init initializes the InfluxDB client
func (b *BasicQueryClient) Init() error {
	cl, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%v", b.Address),
	})

	if err != nil {
		return err
	}

	b.client = cl

	return nil
}

// Query runs the query
func (b *BasicQueryClient) Query(cmd Query) (response, error) {
	q := client.Query{
		Command:  string(cmd),
		Database: b.Database,
	}

	t := NewTimer()
	_, err := b.client.Query(q)
	t.StopTimer()

	if err != nil {
		return response{Timer: t}, err
	}

	// Needs actual response type
	r := response{
		Time:  time.Now(),
		Timer: t,
	}

	return r, nil

}

// Exec listens to the query channel an executes queries as they come in
func (b *BasicQueryClient) Exec(qs <-chan Query, r chan<- response) error {
	if !b.Enabled {
		return nil
	}
	var wg sync.WaitGroup
	counter := NewConcurrencyLimiter(b.Concurrency)

	b.Init()

	interval, err := time.ParseDuration(b.QueryInterval)
	if err != nil {
		return err
	}

	for q := range qs {
		wg.Add(1)
		counter.Increment()
		func(q Query) {
			defer wg.Done()
			qr, _ := b.Query(q)
			r <- qr
			time.Sleep(interval)
			counter.Decrement()
		}(q)
	}

	wg.Wait()

	return nil
}

// resetDB will drop an create a new database on an existing
// InfluxDB instance.
func resetDB(c client.Client, database string) error {
	_, err := c.Query(client.Query{
		// Change to DROP DATABASE %s IF EXISTS
		Command: fmt.Sprintf("DROP DATABASE %s", database),
	})
	if err != nil {
		return err
	}

	_, err = c.Query(client.Query{
		Command: fmt.Sprintf("CREATE DATABASE %s", database),
	})
	if err != nil {
		return err
	}

	return nil
}

// BasicProvisioner implements the Provisioner
// interface.
type BasicProvisioner struct {
	Enabled       bool   `toml:"enabled"`
	Address       string `toml:"address"`
	Database      string `toml:"database"`
	ResetDatabase bool   `toml:"reset_database"`
}

// Provision runs the resetDB function.
func (b *BasicProvisioner) Provision() error {
	if !b.Enabled {
		return nil
	}

	cl, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%v", b.Address),
	})

	if err != nil {
		return err
	}

	if b.ResetDatabase {
		resetDB(cl, b.Database)
	}

	return nil
}

// BasicWriteHandler handles write responses.
func BasicWriteHandler(rs <-chan response, wt *Timer) {
	n := 0
	success := 0
	fail := 0

	s := time.Duration(0)

	for t := range rs {

		n++

		if t.Success() {
			success++
		} else {
			fail++
		}

		s += t.Timer.Elapsed()

	}

	if n == 0 {
		return
	}

	fmt.Printf("Total Requests: %v\n", n)
	fmt.Printf("	Success: %v\n", success)
	fmt.Printf("	Fail: %v\n", fail)
	fmt.Printf("Average Response Time: %v\n", s/time.Duration(n))
	fmt.Printf("Points Per Second: %v\n\n", float64(n)*float64(10000)/float64(wt.Elapsed().Seconds()))
}

// BasicReadHandler handles read responses.
func BasicReadHandler(r <-chan response, rt *Timer) {
	n := 0
	s := time.Duration(0)
	for t := range r {
		n++
		s += t.Timer.Elapsed()
	}

	if n == 0 {
		return
	}

	fmt.Printf("Total Queries: %v\n", n)
	fmt.Printf("Average Query Response Time: %v\n\n", s/time.Duration(n))
}
