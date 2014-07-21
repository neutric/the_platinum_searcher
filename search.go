package the_platinum_searcher

type Searcher struct {
	Root, Pattern string
	Option        *Option
}

func (s *Searcher) Search() error {
	pattern, err := s.pattern()
	if err != nil {
		return err
	}
	grep := make(chan *GrepParams, s.Option.Proc)
	match := make(chan *PrintParams, s.Option.Proc)
	done := make(chan bool)
	go s.find(grep, pattern)
	go s.grep(grep, match)
	go s.print(match, done)
	<-done
	return nil
}

func (s *Searcher) pattern() (*Pattern, error) {
	fileRegexp := s.Option.FileSearchRegexp
	if s.Option.FilesWithRegexp != "" {
		fileRegexp = s.Option.FilesWithRegexp
	}
	return NewPattern(
		s.Pattern,
		fileRegexp,
		s.Option.SmartCase,
		s.Option.IgnoreCase,
		s.Option.Regexp,
	)
}

func (s *Searcher) find(out chan *GrepParams, pattern *Pattern) {
	finder := Finder{out, s.Option}
	finder.Find(s.Root, pattern)
}

func (s *Searcher) grep(in chan *GrepParams, out chan *PrintParams) {
	grepper := Grepper{in, out, s.Option}
	grepper.ConcurrentGrep()
}

func (s *Searcher) print(in chan *PrintParams, done chan bool) {
	printer := NewPrinter(in, done, s.Option)
	printer.Print()
}