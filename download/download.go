package download

import "bytes"
import dbg "fmt"
import "exp/html"
import "log"
import "time"

import "github.com/mewkiz/css"
import "github.com/mewkiz/pkg/htmlutil"
import "github.com/mewkiz/pkg/httputil"

// DefaultTimeout is the default timeout interval, which is used if no timeout
// interval was specified in the config file.
const DefaultTimeout = 1 * time.Minute

// Timeout is the time interval to sleep in between page downloads.
var Timeout time.Duration

type Page struct {
	RawUrl string
	RawSel string
	sel    css.Selector
	nodes  []*html.Node
	dead   bool
}

// NewPage returns a new page based on the provided URL and CSS selector.
func NewPage(rawUrl, rawSel string) (page *Page) {
	page = &Page{
		RawUrl: rawUrl,
		RawSel: rawSel,
	}
	/// ### [ todo ] ###
	///   - how to handle empty selector?
	/// ### [/ todo ] ###
	var err error
	page.sel, err = css.Compile(rawSel)
	if err != nil {
		log.Println(err)
	}
	return page
}

// Equal returns true if page a and b are equal and false otherwise.
func (pageA *Page) Equal(pageB *Page) bool {
	if pageA.RawUrl == pageB.RawUrl && pageA.RawSel == pageB.RawSel {
		return true
	}
	return false
}

// Watch continuously watches the page for changes based on the timeout
// interval, as long as it's existance is justified.
func (page *Page) Watch() {
	for {
		// justify existance.
		if page.dead {
			dbg.Println("page.Watch suicide:", page.RawUrl, page.RawSel)
			return
		}
		dbg.Println("download:", page.RawUrl, page.RawSel)
		err := page.download()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(Timeout)
	}
}

// Kill marks the page as dead, which causes the page watcher to return.
func (page *Page) Kill() {
	page.dead = true
}

// download downloads the page and locates the relevant HTML nodes based on the
// CSS selector.
func (page *Page) download() (err error) {
	buf, err := httputil.GetRaw(page.RawUrl)
	if err != nil {
		return err
	}
	doc, err := html.Parse(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	/// ### [ todo ] ###
	///   - how to handle empty selector?
	/// ### [/ todo ] ###
	if page.RawSel == "" {
		page.nodes = append(page.nodes, doc)
		return nil
	}
	page.nodes = page.sel.MatchAll(doc)
	return nil
}

func (page *Page) String() string {
	w := new(bytes.Buffer)
	for _, node := range page.nodes {
		htmlutil.Render(w, node)
	}
	return string(w.Bytes())
}

func init() {
	httputil.SetClient(httputil.InsecureClient)
}
