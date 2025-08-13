package runner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"maps"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/domain"
	"github.com/IsmailCLN/tapir/internal/httpclient"
	"github.com/IsmailCLN/tapir/internal/sharedcontext"
)

// Options controls the concurrent runner.
type Options struct {
	// Concurrency is the number of worker goroutines.
	// If <= 0, runtime.NumCPU() is used.
	Concurrency int
}

// RunConcurrent executes all requests across the given suites in parallel,
// but respects per-suite dependencies declared via TestRequest.DependsOn.
// Each expectation result is streamed on the returned channel as soon as evaluated.
func RunConcurrent(ctx context.Context, suites []domain.TestSuite, opts Options) <-chan Result {
	out := make(chan Result)

	shared := sharedcontext.New()
	assert.SetSharedContext(shared)

	type job struct {
		SuiteName string
		Req       domain.TestRequest
	}
	type done struct {
		SuiteName string
		ReqName   string
	}

	jobs := make(chan job)
	doneCh := make(chan done)

	// workers
	n := opts.Concurrency
	if n <= 0 {
		n = runtime.NumCPU()
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for jb := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				results := runRequest(ctx, jb.SuiteName, jb.Req, shared)
				for _, r := range results {
					select {
					case out <- r:
					case <-ctx.Done():
						return
					}
				}
				// notify scheduler this request is finished
				select {
				case doneCh <- done{SuiteName: jb.SuiteName, ReqName: jb.Req.Name}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// scheduler: build per-suite DAG and feed ready jobs
	go func() {
		defer close(jobs)

		// Build graphs
		type graph struct {
			indeg    map[string]int
			children map[string][]string
			reqs     map[string]domain.TestRequest
			total    int
			sent     int
			done     int
		}

		graphs := make(map[string]*graph)

		for _, s := range suites {
			g := &graph{
				indeg:    make(map[string]int),
				children: make(map[string][]string),
				reqs:     make(map[string]domain.TestRequest),
			}
			for _, r := range s.Requests {
				g.reqs[r.Name] = r
			}
			// init indegrees
			for name := range g.reqs {
				g.indeg[name] = 0
			}
			// add edges
			for _, r := range s.Requests {
				for _, dep := range r.DependsOn {
					if _, ok := g.reqs[dep]; !ok {
						// Unknown dependency: emit a configuration error result but proceed.
						select {
						case out <- Result{
							Suite:    s.Name,
							Request:  r.Name,
							Passed:   false,
							Err:      fmt.Errorf("depends_on references unknown request %q", dep),
							TestName: "depends_on",
						}:
						case <-ctx.Done():
							return
						}
						// don't increase indegree (so it can still run)
						continue
					}
					g.indeg[r.Name]++
					g.children[dep] = append(g.children[dep], r.Name)
				}
			}
			g.total = len(g.reqs)
			graphs[s.Name] = g
		}

		// queue initial ready jobs
		var activeSuites int
		for sName, g := range graphs {
			for name, deg := range g.indeg {
				if deg == 0 {
					select {
					case jobs <- job{SuiteName: sName, Req: g.reqs[name]}:
						g.sent++
						activeSuites++
					case <-ctx.Done():
						return
					}
				}
			}
		}

		// react to completions and release dependents
		totalSent := 0
		for _, g := range graphs {
			totalSent += g.sent
		}

		for {
			if totalSent == 0 {
				// nothing sent initially (e.g., cycles or empty), break
				break
			}
			select {
			case d := <-doneCh:
				// mark done
				g := graphs[d.SuiteName]
				g.done++
				// release children
				for _, child := range g.children[d.ReqName] {
					g.indeg[child]--
					if g.indeg[child] == 0 {
						select {
						case jobs <- job{SuiteName: d.SuiteName, Req: g.reqs[child]}:
							g.sent++
							totalSent++
						case <-ctx.Done():
							return
						}
					}
				}
				// if finished all from this suite and all suites done, we may continue waiting others
				// Stop when every suite has g.done == g.total
				allDone := true
				for _, gg := range graphs {
					if gg.done < gg.total {
						allDone = false
						break
					}
				}
				if allDone {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// close 'out' when all workers finish
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// runRequest executes a single request and returns one Result per expectation.
func runRequest(ctx context.Context, suite string, r domain.TestRequest, shared *sharedcontext.SharedContext) []Result {
	var results []Result

	// ----- 1. Build request body (string only for now) -----
	var bodyReader io.Reader
	if bodyStr, ok := r.Req.Body.(string); ok && bodyStr != "" {
		bodyReader = strings.NewReader(bodyStr)
	}

	// ----- 2. Construct HTTP request -----
	req, err := http.NewRequest(r.Req.Method, r.Req.URL, bodyReader)
	if err != nil {
		appendRequestErrorResults(&results, suite, r, err)
		return results
	}

	// ----- 3. Apply headers with placeholder substitution -----
	for k, v := range r.Req.Headers {
		if strings.Contains(v, "${token}") {
			if t, ok := shared.Get("token"); ok {
				v = strings.ReplaceAll(v, "${token}", t)
			}
		}
		req.Header.Set(k, v)
	}

	// ----- 4. Send request -----
	resp, err := httpclient.Do(ctx, req)
	if err != nil {
		appendRequestErrorResults(&results, suite, r, err)
		return results
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		appendRequestErrorResults(&results, suite, r, err)
		return results
	}

	// ----- 5. Evaluate expectations -----
	for _, exp := range r.Expect {
		// 5a. Copy user‑provided kwargs
		kwargs := make(map[string]any, len(exp.Kwargs)+2)
		maps.Copy(kwargs, exp.Kwargs)

		// 5b. Inject auto params
		kwargs["status_code"] = resp.StatusCode
		kwargs["headers"] = resp.Header

		f, ok := assert.Get(exp.Type)
		if !ok {
			results = append(results, Result{
				Suite:    suite,
				Request:  r.Name,
				Passed:   false,
				Err:      fmt.Errorf("unknown expectation %s", exp.Type),
				TestName: exp.Type,
			})
			continue
		}

		err := f(bodyBytes, kwargs)
		results = append(results, Result{
			Suite:    suite,
			Request:  r.Name,
			Passed:   err == nil,
			Err:      err,
			TestName: exp.Type,
		})
	}

	return results
}
